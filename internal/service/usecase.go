package service

import (
	"github.com/SanSanchezzz/tech-db-forum/internal/errors"
	"github.com/SanSanchezzz/tech-db-forum/internal/models"
)

type ServiceUsecase interface {
	GetStatus() (*models.Service, *errors.ErrorResponse)
	Clear() *errors.ErrorResponse
}
