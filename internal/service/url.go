package service

import (
	"context"
	"errors"
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
	_, err := url.ParseRequestURI(originalURL)
	if err != nil {
		return nil, errors.New("invalid URL")
	}

	shortCode := customCode
	if shortCode == "" {
		shortCode = utils.GenerateShortCode(6)
	}

	urlModel := &model.URL{
		ShortCode:   shortCode,
		OriginalURL: originalURL,
		CreatedAt:   time.Now(),
	}

	err = s.repo.CreateURL(ctx, urlModel)
	if err != nil {
		return nil, err
	}

	_ = s.cache.SetURL(ctx, shortCode, originalURL, 24*time.Hour)

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
