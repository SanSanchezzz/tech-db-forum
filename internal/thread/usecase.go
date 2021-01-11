package thread

import (
	"github.com/SanSanchezzz/tech-db-forum/internal/errors"
	"github.com/SanSanchezzz/tech-db-forum/internal/models"
)

type ThreadUsecase interface {
	GetThread(slug string) (*models.Thread, *errors.ErrorResponse)
	CreateThread(thread *models.Thread) (uint32, *errors.ErrorResponse)
	GetThreadsByParam(limit int, desc bool, slug, since *string) ([]*models.Thread, *errors.ErrorResponse)
	VoteThread(voice int, nickname, slugOrID string) (*models.Thread, *errors.ErrorResponse)

	UpdateThread(thread *models.Thread) *errors.ErrorResponse

}
