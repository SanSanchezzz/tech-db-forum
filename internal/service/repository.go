package service

import "github.com/SanSanchezzz/tech-db-forum/internal/models"

type ServiceRepository interface {
	GetStatus() *models.Service
	Clear() error
}
