package test

import (
	"context"
	"smolink/internal/model"
	"time"
)

func (ta *TestApp) SeedShortURL(shortCode, originalURL string) error {
	return ta.PGRepo.CreateURL(context.Background(), &model.URL{
		ShortCode:   shortCode,
		OriginalURL: originalURL,
		CreatedAt:   time.Now(),
	})
}
