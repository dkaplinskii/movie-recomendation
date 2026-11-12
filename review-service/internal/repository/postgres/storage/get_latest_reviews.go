package storage

import (
	"context"
	"fmt"

	"studentgit.kata.academy/movie-recommendation-platform-20/review-service/internal/entity"
)

const getLatestReviewsQuery = `
SELECT film_id, review, created_at
FROM reviews
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2
`

func (r *PostgresRepository) GetLatestReviews(ctx context.Context, userID int64, limit int32) ([]entity.Review, error) {
	rows, err := r.pool.Query(
		ctx,
		getLatestReviewsQuery,
		userID,
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf("query reviews: %w", err)
	}
	defer rows.Close()

	reviews := make([]entity.Review, 0)

	for rows.Next() {
		var review entity.Review

		if err := rows.Scan(
			&review.FilmID,
			&review.Review,
			&review.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan review: %w", err)
		}

		reviews = append(reviews, review)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate reviews: %w", err)
	}

	return reviews, nil
}
