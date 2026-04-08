package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RefreshRepository struct {
	db *pgxpool.Pool
}

func NewRefreshRepository(db *pgxpool.Pool) *RefreshRepository {
	return &RefreshRepository{db: db}
}

func (r *RefreshRepository) Save(
	userID int,
	token string,
	deviceID string,
	userAgent string,
	ip string,
	exp time.Time,
) error {

	query := `
	INSERT INTO refresh_tokens (user_id, token, device_id, user_agent, ip, expires_at)
	VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.Exec(
		context.Background(),
		query,
		userID,
		token,
		deviceID,
		userAgent,
		ip,
		exp,
	)

	return err
}

func (r *RefreshRepository) Find(token string) (int, error) {
	var userID int

	query := `
	SELECT user_id
	FROM refresh_tokens
	WHERE token=$1 AND expires_at > now()
	`

	err := r.db.QueryRow(context.Background(), query, token).
		Scan(&userID)

	return userID, err
}

func (r *RefreshRepository) Delete(token string) error {
	_, err := r.db.Exec(context.Background(),
		`DELETE FROM refresh_tokens WHERE token=$1`, token)
	return err
}
