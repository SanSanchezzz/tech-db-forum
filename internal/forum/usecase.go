package forum

import (
	"github.com/SanSanchezzz/tech-db-forum/internal/errors"
	"github.com/SanSanchezzz/tech-db-forum/internal/models"
)

type ForumUsecase interface {
	GetForum(slug string) (*models.Forum, *errors.ErrorResponse)
	CreateForum(forum *models.Forum) *errors.ErrorResponse

}
