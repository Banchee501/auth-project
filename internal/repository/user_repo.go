package repository

import (
	"auth-project/internal/models"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user models.User) (int, error) {
	var id int

	query := `
	INSERT INTO users (email, password_hash)
	VALUES ($1, $2)
	RETURNING id
	`

	err := r.db.QueryRow(context.Background(), query,
		user.Email,
		user.PasswordHash,
	).Scan(&id)

	return id, err
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User

	query := `
	SELECT id, email, password_hash
	FROM users
	WHERE email=$1
	`

	err := r.db.QueryRow(context.Background(), query, email).
		Scan(&user.ID, &user.Email, &user.PasswordHash)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetByID(id int) (*models.User, error) {

	var user models.User

	query := `
	SELECT id, email, password_hash
	FROM users
	WHERE id=$1
	`

	err := r.db.QueryRow(context.Background(), query, id).
		Scan(&user.ID, &user.Email, &user.PasswordHash)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
