package db

import (
	"context"
	"time"
)

type TopURL struct {
	ShortKey   string     `json:"short_key"`
	ClickCount int        `json:"click_count"`
	CreatedAt  time.Time  `json:"created_at"`
	ExpireAt   *time.Time `json:"expire_at,omitempty"`
}

func GetTopURLs(ctx context.Context) ([]TopURL, error) {
	rows, err := Pool.Query(ctx, `
	SELECT short_key, click_count, created_at, expire_at
	FROM url_mappings
	ORDER BY click_count DESC
	LIMIT 10
`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []TopURL
	for rows.Next() {
		var t TopURL
		if err := rows.Scan(&t.ShortKey, &t.ClickCount, &t.CreatedAt, &t.ExpireAt); err != nil {
			return nil, err
		}
		results = append(results, t)
	}

	return results, nil
}
