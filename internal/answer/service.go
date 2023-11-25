package answer

import (
	"context"
	"github.com/blazee5/quizmaster-backend/internal/domain"
)

type Service interface {
	Create(ctx context.Context, userID, quizID, questionID int, input domain.CreateAnswer) (int, error)
	Update(ctx context.Context, answerID, userID, quizID int, input domain.Answer) error
	Delete(ctx context.Context, answerID, userID, quizID int) error
	ChangeOrder(ctx context.Context, userID, quizID, questionID int, input domain.ChangeAnswerOrder) error
}
