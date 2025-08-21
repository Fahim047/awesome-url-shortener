package db

import (
	"context"
)

type TopURL struct {
	ShortKey   string `json:"short_key"`
	ClickCount int    `json:"click_count"`
}

func GetTopURLs(ctx context.Context) ([]TopURL, error) {
	rows, err := Pool.Query(ctx, `
		SELECT short_key, click_count 
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
		if err := rows.Scan(&t.ShortKey, &t.ClickCount); err != nil {
			return nil, err
		}
		results = append(results, t)
	}
	return results, nil
}
