package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Session struct {
	UserID    int
	DeviceID  string
	UserAgent string
	IP        string
}

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
	expiresAt time.Time,
) error {
	query := `
	INSERT INTO refresh_tokens (user_id, token, device_id, user_agent, ip, expires_at)
	VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.Exec(context.Background(), query,
		userID, token, deviceID, userAgent, ip, expiresAt)

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

func (r *RefreshRepository) GetSession(token string) (*Session, error) {
	query := `
	SELECT user_id, device_id, user_agent, ip
	FROM refresh_tokens
	WHERE token=$1 AND expires_at > now()
	`

	var s Session

	err := r.db.QueryRow(context.Background(), query, token).
		Scan(&s.UserID, &s.DeviceID, &s.UserAgent, &s.IP)

	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (r *RefreshRepository) Delete(token string) error {
	_, err := r.db.Exec(context.Background(),
		`DELETE FROM refresh_tokens WHERE token=$1`, token)
	return err
}
