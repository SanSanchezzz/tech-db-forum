package repository

import (
	"database/sql"
	"github.com/SanSanchezzz/tech-db-forum/internal/forum"
	"github.com/SanSanchezzz/tech-db-forum/internal/models"
)

type ForumRepository struct {
	dbConn *sql.DB
}

func NewForumRepository(conn *sql.DB) forum.ForumRepository{
	return &ForumRepository{
		dbConn: conn,
	}
}

func (fr *ForumRepository) SelectBySlug(slug string) (*models.Forum, error) {
	forum := &models.Forum{}

	err := fr.dbConn.QueryRow(
		`select slug, title, nickname, threads, posts from forums where slug = $1`,
		slug,
	).Scan(
		&forum.Slug,
		&forum.Title,
		&forum.User,
		&forum.Threads,
		&forum.Posts,
	)

	if err != nil {
		return nil, err
	}

	return forum, nil
}

func (fr *ForumRepository) Insert(forum *models.Forum) error {
	_, err := fr.dbConn.Exec(
		`insert into forums (slug, title, nickname) values ($1, $2, $3)`,
		forum.Slug, forum.Title, forum.User,
	)

	if err != nil {
		return err
	}

	return nil
}
