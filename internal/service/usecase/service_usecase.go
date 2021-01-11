package usecase

import (
	"github.com/SanSanchezzz/tech-db-forum/internal/errors"
	"github.com/SanSanchezzz/tech-db-forum/internal/models"
	"github.com/SanSanchezzz/tech-db-forum/internal/service"
)

type ServiceUsecase struct {
	ServiceRepository service.ServiceRepository
}

func NewServiceUsecase(serviceRepository service.ServiceRepository) *ServiceUsecase {
	return &ServiceUsecase{
		ServiceRepository: serviceRepository,
	}
}

func (uc *ServiceUsecase) GetStatus() (*models.Service, *errors.ErrorResponse) {
	service := uc.ServiceRepository.GetStatus()
	if service == nil {
		return nil, errors.NewErrorResponse(errors.ErrInternal, nil)
	}

	return service, nil
}

func (uc *ServiceUsecase) Clear() *errors.ErrorResponse {
	err := uc.ServiceRepository.Clear()
	if err != nil {
		return errors.NewErrorResponse(errors.ErrInternal, nil)
	}

	return nil
}
