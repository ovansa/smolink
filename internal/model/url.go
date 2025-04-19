package model

import "time"

type URL struct {
	ID          int    `json:"id"`
	ShortCode   string `json:"short_code"`
	OriginalURL string `json:"original_url"`
	ClickCount  int    `json:"click_count"`
	CreatedAt   string `json:"created_at"`
}

type URLAnalytics struct {
	ID         int       `json:"id"`
	URLID      int       `json:"url_id"`
	IPAddress  string    `json:"ip_address"`
	UserAgent  string    `json:"user_agent"`
	AccessedAt time.Time `json:"accessed_at"`
}
