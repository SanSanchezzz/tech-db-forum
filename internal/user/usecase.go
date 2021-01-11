package user

import (
	"github.com/SanSanchezzz/tech-db-forum/internal/errors"
	"github.com/SanSanchezzz/tech-db-forum/internal/models"
)

type UserUsecase interface {
	Create(nickname, fullname, about, email string) ([]*models.User, *errors.ErrorResponse)
	GetByNickname(nickname string) (*models.User, *errors.ErrorResponse)
	CheckUserByEmail(email string) (bool, *errors.ErrorResponse)
	UpdateProfile(user *models.User) *errors.ErrorResponse
	GetUsersByForum(limit int, desc bool, slug, since *string) ([]*models.User, *errors.ErrorResponse)
}
