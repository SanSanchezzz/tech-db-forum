package forum

import "github.com/SanSanchezzz/tech-db-forum/internal/models"

type ForumRepository interface {
	SelectBySlug(slug string) (*models.Forum, error)
	Insert(forum *models.Forum) error

}
