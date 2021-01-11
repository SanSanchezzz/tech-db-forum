package post

import (
	"github.com/SanSanchezzz/tech-db-forum/internal/models"
	"time"
)

type PostRepository interface {

	Insert(posts []*models.Post, thread *models.Thread, now *time.Time) error
	GetPostsDesc(limit, since *string, id uint32) ([]*models.Post, error)
	GetPosts(limit, since *string, id uint32) ([]*models.Post, error)

	GetPostsDescTree(limit, since *string, id uint32) ([]*models.Post, error)
	GetPostsTree(limit, since *string, id uint32) ([]*models.Post, error)

	GetPostsDescParentTree(limit, since *string, id uint32) ([]*models.Post, error)
	GetPostsParentTree(limit, since *string, id uint32) ([]*models.Post, error)
	GetPost(id int) (*models.Post, error)
	Update(post *models.Post) error


}
