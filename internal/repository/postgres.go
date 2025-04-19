package repository

import (
	"context"
	"smolink/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) CreateURL(ctx context.Context, url *model.URL) error {
	_, err := r.db.Exec(ctx, "INSERT INTO urls (short_code, original_url) VALUES ($1, $2)", url.ShortCode, url.OriginalURL)
	return err
}

func (r *PostgresRepository) GetURL(ctx context.Context, shortCode string) (*model.URL, error) {
	var url model.URL
	err := r.db.QueryRow(ctx, "SELECT id, short_code, original_url, click_count FROM urls WHERE short_code = $1", shortCode).Scan(&url.ID, &url.ShortCode, &url.OriginalURL, &url.ClickCount)
	if err != nil {
		return nil, err
	}
	return &url, nil
}

func (r *PostgresRepository) IncrementClickCount(ctx context.Context, urlID int) error {
	_, err := r.db.Exec(ctx, "UPDATE urls SET click_count = click_count + 1 WHERE id = $1", urlID)
	return err
}

func (r *PostgresRepository) LogAnalytics(ctx context.Context, analytics *model.URLAnalytics) error {
	_, err := r.db.Exec(ctx, "INSERT INTO url_analytics (url_id, ip_address, user_agent) VALUES ($1, $2, $3)", analytics.URLID, analytics.IPAddress, analytics.UserAgent)
	return err
}
