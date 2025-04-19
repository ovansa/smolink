package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
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
	// 1. Validate URL
	if _, err := url.ParseRequestURI(originalURL); err != nil {
		return nil, errors.New("invalid URL")
	}

	// 2. Handle custom code or generate new one
	var shortCode string
	var err error

	if customCode != "" {
		// Check if custom code already exists
		if existingURL, _ := s.repo.GetURL(ctx, customCode); existingURL != nil {
			return nil, errors.New("custom code already in use")
		}
		shortCode = customCode
	} else {
		// Generate secure random code with retry logic
		for i := 0; i < 3; i++ { // Try 3 times to generate unique code
			shortCode, err = utils.GenerateShortCodeSecure(6)
			if err != nil {
				return nil, fmt.Errorf("failed to generate short code: %w", err)
			}

			// Verify uniqueness
			if existingURL, _ := s.repo.GetURL(ctx, shortCode); existingURL == nil {
				break
			}
		}
	}

	// 3. Create URL record
	urlModel := &model.URL{
		ShortCode:   shortCode,
		OriginalURL: originalURL,
		CreatedAt:   time.Now(),
	}

	if err := s.repo.CreateURL(ctx, urlModel); err != nil {
		return nil, fmt.Errorf("failed to create URL: %w", err)
	}

	// 4. Update cache (log failure but don't fail operation)
	if err := s.cache.SetURL(ctx, shortCode, originalURL, 24*time.Hour); err != nil {
		// s.logger.Error("failed to cache URL",
		// 	"error", err,
		// 	"short_code", shortCode,
		// 	"original_url", originalURL,
		// )
		log.Printf("failed to cache URL: %v", err)
	}

	return urlModel, nil
}

func (s *URLService) ResolveURL(ctx context.Context, shortCode, ip, userAgent string) (string, error) {
	// Try cache
	original, err := s.cache.GetURL(ctx, shortCode)
	if err == nil {
		go s.recordAnalytics(ctx, shortCode, ip, userAgent)
		return original, nil
	}

	// Fallback to DB
	urlModel, err := s.repo.GetURL(ctx, shortCode)
	if err != nil {
		return "", err
	}

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
