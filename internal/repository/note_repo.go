package repository

import (
	"auth-project/internal/models"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type NoteRepository struct {
	db *pgxpool.Pool
}

func NewNoteRepository(db *pgxpool.Pool) *NoteRepository {
	return &NoteRepository{db: db}
}

func (r *NoteRepository) Create(userID int, text string) error {
	query := `
	INSERT INTO notes (user_id, text)
	VALUES ($1, $2)
	`

	_, err := r.db.Exec(context.Background(), query, userID, text)
	return err
}

func (r *NoteRepository) GetAllByUser(userID int) ([]models.Note, error) {

	query := `
	SELECT id, user_id, text
	FROM notes
	WHERE user_id=$1
	`

	rows, err := r.db.Query(context.Background(), query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []models.Note

	for rows.Next() {
		var n models.Note
		err := rows.Scan(&n.ID, &n.UserID, &n.Text)
		if err != nil {
			return nil, err
		}
		notes = append(notes, n)
	}

	return notes, nil
}

func (r *NoteRepository) Delete(noteID, userID int) error {

	query := `
	DELETE FROM notes
	WHERE id=$1 AND user_id=$2
	`

	res, err := r.db.Exec(context.Background(), query, noteID, userID)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("not found or not allowed")
	}

	return nil
}
