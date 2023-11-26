package repository

import (
	"context"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	"github.com/blazee5/quizmaster-backend/internal/models"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (repo *Repository) Create(ctx context.Context, questionID int, input domain.CreateAnswer) (int, error) {
	var id int

	err := repo.db.QueryRowxContext(ctx, "INSERT INTO answers (question_id, order_id) VALUES ($1, $2) RETURNING id",
		questionID, input.OrderID).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (repo *Repository) GetByID(ctx context.Context, id int) (models.Answer, error) {
	var answer models.Answer

	err := repo.db.QueryRowxContext(ctx, "SELECT * FROM answers WHERE id = $1", id).StructScan(&answer)

	if err != nil {
		return models.Answer{}, err
	}

	return answer, nil
}

func (repo *Repository) GetAnswersByQuestionID(ctx context.Context, id int) ([]models.Answer, error) {
	answers := make([]models.Answer, 0)

	err := repo.db.SelectContext(ctx, &answers, "SELECT * FROM answers WHERE question_id = $1", id)

	if err != nil {
		return nil, err
	}

	return answers, nil
}

func (repo *Repository) Update(ctx context.Context, answerID int, input domain.Answer) error {
	err := repo.db.QueryRowxContext(ctx, `UPDATE answers SET
		text = $1,
		is_correct = $2
		WHERE id = $3`,
		input.Text, input.IsCorrect, answerID).Err()

	if err != nil {
		return err
	}

	return nil
}

func (repo *Repository) Delete(ctx context.Context, answerID int) error {
	err := repo.db.QueryRowxContext(ctx, `DELETE FROM answers WHERE id = $1`, answerID).Err()

	if err != nil {
		return err
	}

	return nil
}

func (repo *Repository) ChangeOrder(ctx context.Context, questionID int, input domain.ChangeAnswerOrder) error {
	tx, err := repo.db.Beginx()

	if err != nil {
		return err
	}

	for _, item := range input.Orders {
		_, err := tx.ExecContext(ctx, "UPDATE answers SET order_id = $1 WHERE id = $2 AND question_id = $3", item.OrderID, item.AnswerID, questionID)

		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}