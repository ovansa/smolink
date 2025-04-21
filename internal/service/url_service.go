package service

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"smolink/internal/errors"
	"smolink/internal/model"
	"smolink/internal/repository"
	"smolink/pkg/utils"
	"time"
)

type URLService struct {
	repo  *repository.PostgresRepository
	cache *repository.RedisRepository
}

func NewURLService(repo *repository.PostgresRepository, cache *repository.RedisRepository) *URLService {
	return &URLService{repo: repo, cache: cache}
}

func (s *URLService) ShortenURL(ctx context.Context, originalURL, customCode string) (*model.URL, error) {
	if _, err := url.ParseRequestURI(originalURL); err != nil {
		return nil, errors.ErrInvalidURL
	}

	var shortCode string
	var err error

	if customCode != "" {
		if existingURL, _ := s.repo.GetURL(ctx, customCode); existingURL != nil {
			return nil, errors.ErrCodeInUse
		}
		shortCode = customCode
	} else {
		// Generate secure random code with retry logic
		for i := 0; i < 3; i++ { // Try 3 times to generate unique code
			shortCode, err = utils.GenerateShortCodeSecure(6)
			if err != nil {
				return nil, fmt.Errorf("%w: %v", errors.ErrInternal, err)
			}

			// Verify uniqueness
			if existingURL, _ := s.repo.GetURL(ctx, shortCode); existingURL == nil {
				break
			}
		}
	}

	urlModel := &model.URL{
		ShortCode:   shortCode,
		OriginalURL: originalURL,
		CreatedAt:   time.Now(),
	}

	if err := s.repo.CreateURL(ctx, urlModel); err != nil {
		return nil, fmt.Errorf("%w %v", errors.ErrInternal, err)
	}

	if err := s.cache.SetURL(ctx, shortCode, originalURL, 24*time.Hour); err != nil {
		log.Printf("failed to cache URL: %v", err)
	}

	return urlModel, nil
}

func (s *URLService) ResolveURL(ctx context.Context, shortCode, ip, userAgent string) (string, error) {
	original, err := s.cache.GetURL(ctx, shortCode)
	if err == nil {
		log.Print("Successfully fetched from Cache")
		go s.recordAnalytics(ctx, shortCode, ip, userAgent)
		return original, nil
	}

	// Fallback to DB
	urlModel, err := s.repo.GetURL(ctx, shortCode)
	log.Print("Did not find record from cache. Fetching from DB")
	if err != nil {
		return "", errors.ErrShortCodeNotFound
	}

	log.Print("Successfully fetched from DB")
	_ = s.cache.SetURL(ctx, shortCode, urlModel.OriginalURL, 24*time.Hour)
	go s.recordAnalytics(ctx, shortCode, ip, userAgent)

	return urlModel.OriginalURL, nil
}

func (s *URLService) recordAnalytics(ctx context.Context, shortCode, ip, userAgent string) {
	urlModel, err := s.repo.GetURL(ctx, shortCode)
	if err != nil {
		return
	}

	_ = s.repo.IncrementClickCount(ctx, urlModel.ID)
	_ = s.repo.LogAnalytics(ctx, &model.URLAnalytics{
		URLID:      urlModel.ID,
		IPAddress:  ip,
		UserAgent:  userAgent,
		AccessedAt: time.Now(),
	})
}
