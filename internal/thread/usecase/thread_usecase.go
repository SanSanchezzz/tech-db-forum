package usecase

import (
	"database/sql"
	"github.com/SanSanchezzz/tech-db-forum/internal/errors"
	"github.com/SanSanchezzz/tech-db-forum/internal/models"
	"github.com/SanSanchezzz/tech-db-forum/internal/thread"
	"github.com/SanSanchezzz/tech-db-forum/internal/user"
	"strconv"
)

type ThreadUsecase struct {
	ThreadRepository thread.ThreadRepository
	UserRepository user.UserRepository
}

func NewThreadUsecase(threadRepository thread.ThreadRepository, userRepository user.UserRepository) *ThreadUsecase {
	return &ThreadUsecase{
		ThreadRepository: threadRepository,
		UserRepository: userRepository,
	}
}

func (tu *ThreadUsecase) GetThread(slugOrID string) (*models.Thread, *errors.ErrorResponse) {
	var thread *models.Thread
	var err error

	if id, err := strconv.Atoi(slugOrID); err == nil {
		thread, err = tu.ThreadRepository.SelectByID(uint32(id))
	} else {
		thread, err = tu.ThreadRepository.SelectBySlug(slugOrID)
	}

	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, errors.NewErrorResponse(errors.ErrInternal, err)
	default:
		return thread, nil
	}
}

func (tu *ThreadUsecase) CreateThread(thread *models.Thread) (uint32, *errors.ErrorResponse) {
	id, err := tu.ThreadRepository.Insert(thread)

	if err != nil {
		return 0, errors.NewErrorResponse(errors.ErrInternal, err)
	}

	return id, nil
}

func (tu *ThreadUsecase) UpdateThread(thread *models.Thread) *errors.ErrorResponse {
	err := tu.ThreadRepository.Update(thread)

	if err != nil {
		return errors.NewErrorResponse(errors.ErrInternal, err)
	}

	return nil
}

func (tu *ThreadUsecase) GetThreadsByParam(limit int, desc bool, slug, since *string) ([]*models.Thread, *errors.ErrorResponse) {
	threads, err := tu.ThreadRepository.SelectThreadsByParam(limit, desc, slug, since)
	if err != nil {
		return nil, errors.NewErrorResponse(errors.ErrInternal, err)
	}

	return threads, nil
}

func (tu *ThreadUsecase) VoteThread(voice int, nickname, slugOrID string) (*models.Thread, *errors.ErrorResponse) {
	thread, err := tu.ThreadRepository.InsertVote(voice, nickname, slugOrID)

	if err != nil {
		return nil, errors.NewErrorResponse(errors.ErrNotFoundUser, err)
	}

	return thread, nil
}

