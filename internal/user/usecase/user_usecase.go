package usecase

import (
	"database/sql"
	"github.com/SanSanchezzz/tech-db-forum/internal/errors"
	"github.com/SanSanchezzz/tech-db-forum/internal/models"
	"github.com/SanSanchezzz/tech-db-forum/internal/user"
)

type UserUsecase struct {
	UserRepository user.UserRepository
}

func NewUserUsecase(userRepository user.UserRepository) *UserUsecase {
	return &UserUsecase{
		UserRepository: userRepository,
	}
}

func (uc *UserUsecase) Create(nickname, fullname, about, email string) ([]*models.User, *errors.ErrorResponse) {
	users, err := uc.UserRepository.GetByNicknameOrEmail(nickname, email)
	if err != nil {
		return nil, errors.NewErrorResponse(errors.ErrInternal, err)
	}

	if len(users) > 0 {
		return users, nil
	}

	err = uc.UserRepository.Insert(nickname, fullname, about, email)
	if err != nil {
		return nil, errors.NewErrorResponse(errors.ErrInternal, err)
	}

	return nil, nil
}

func (uc *UserUsecase) GetUsersByForum(limit int, desc bool, slug, since *string) ([]*models.User, *errors.ErrorResponse) {
	users, err := uc.UserRepository.SelectUsersByForum(limit, desc, slug, since)

	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, errors.NewErrorResponse(errors.ErrInternal, err)
	default:
		return users, nil
	}
}

func (uc *UserUsecase) GetByNickname(nickname string) (*models.User, *errors.ErrorResponse) {
	user, err := uc.UserRepository.SelectByNickname(nickname)
	switch err {
	case sql.ErrNoRows:
		return nil, errors.NewErrorResponse(errors.ErrNotFoundUser, nil)
	}
	if err != nil {
		return nil, errors.NewErrorResponse(errors.ErrInternal, nil)
	}

	return user, nil
}

func (uc *UserUsecase) CheckUserByEmail(email string) (bool, *errors.ErrorResponse) {
	_, err := uc.UserRepository.SelectUserByEmail(email)

	switch {
	case err == sql.ErrNoRows:
		return false, nil
	case err != nil:
		return false, errors.NewErrorResponse(errors.ErrInternal, nil)
	}

	return true, nil
}

func (uc *UserUsecase) UpdateProfile(user *models.User) *errors.ErrorResponse {
	err := uc.UserRepository.Update(user)

	if err != nil {
		return errors.NewErrorResponse(errors.ErrInternal, nil)
	}

	return nil
}
