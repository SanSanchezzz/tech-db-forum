package post

import (
	"github.com/SanSanchezzz/tech-db-forum/internal/errors"
	"github.com/SanSanchezzz/tech-db-forum/internal/models"
)

type PostUsecase interface {
	CreatePosts(posts []*models.Post, thread *models.Thread) *errors.ErrorResponse
	GetPosts(limit, since, desc, sort *string, id uint32) ([]*models.Post, *errors.ErrorResponse)
	GetFullPost(id int, related string) (*models.PostFull, *errors.ErrorResponse)
	GetPost(id int) (*models.Post, *errors.ErrorResponse)
	UpdatePost(post *models.Post) *errors.ErrorResponse


	//GetFullPost()

}
