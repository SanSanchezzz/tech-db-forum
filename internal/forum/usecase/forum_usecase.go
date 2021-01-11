package usecase

import (
	"database/sql"
	"github.com/SanSanchezzz/tech-db-forum/internal/errors"
	"github.com/SanSanchezzz/tech-db-forum/internal/forum"
	"github.com/SanSanchezzz/tech-db-forum/internal/models"
)

type ForumUsecase struct {
	ForumRepository forum.ForumRepository
}

func NewForumUsecse(forumRepository forum.ForumRepository) *ForumUsecase {
	return &ForumUsecase{
		ForumRepository: forumRepository,
	}
}

func (fu *ForumUsecase) GetForum(slug string) (*models.Forum, *errors.ErrorResponse) {
	forum, err := fu.ForumRepository.SelectBySlug(slug)

	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, errors.NewErrorResponse(errors.ErrInternal, err)
	default:
		return forum, nil
	}
}

func (fu *ForumUsecase) CreateForum(forum *models.Forum) *errors.ErrorResponse {
	err := fu.ForumRepository.Insert(forum)

	if err != nil {
		return errors.NewErrorResponse(errors.ErrInternal, err)
	}

	return nil
}
