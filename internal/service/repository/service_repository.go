package repository

import (
	"database/sql"
	"github.com/SanSanchezzz/tech-db-forum/internal/models"
	"github.com/SanSanchezzz/tech-db-forum/internal/service"
)

type ServiceRepository struct {
	dbConn *sql.DB
}

func NewServiceRepository(conn *sql.DB) service.ServiceRepository {
	return &ServiceRepository{
		dbConn: conn,
	}
}

func (ur *ServiceRepository) GetStatus() *models.Service {
	var service models.Service

	err := ur.dbConn.QueryRow(
		`select count(nickname) from users;`,
		).Scan(&service.User)

	err = ur.dbConn.QueryRow(
		`select count( * ), sum(threads), sum(posts) from forums;`,
		).Scan(&service.Forum, &service.Thread, &service.Post)

	if err != sql.ErrNoRows {
		return &service
	}
	if err != nil {
		return nil
	}

	return &service
}

func (ur *ServiceRepository) Clear() error {
	_, err := ur.dbConn.Exec(
		`truncate votes, posts, threads, forums, users restart identity cascade`,
	)

	if err != nil {
		return nil
	}

	return err
}


