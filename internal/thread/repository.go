package thread

import "github.com/SanSanchezzz/tech-db-forum/internal/models"

type ThreadRepository interface {
	SelectBySlug(slug string) (*models.Thread, error)
	SelectByID(id uint32) (*models.Thread, error)

	Insert(thread *models.Thread) (uint32, error)
	SelectThreadsByParam(limit int, desc bool, slug, since *string) ([]*models.Thread, error)
 	InsertVote(voice int, nickname, slugOrID string) (*models.Thread, error)

	Update(thread *models.Thread) error


}
