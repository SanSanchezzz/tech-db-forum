package usecase

import (
	"database/sql"
	"github.com/SanSanchezzz/tech-db-forum/internal/errors"
	"github.com/SanSanchezzz/tech-db-forum/internal/forum"
	"github.com/SanSanchezzz/tech-db-forum/internal/models"
	"github.com/SanSanchezzz/tech-db-forum/internal/post"
	"github.com/SanSanchezzz/tech-db-forum/internal/thread"
	"github.com/SanSanchezzz/tech-db-forum/internal/user"
	"github.com/lib/pq"
	"strings"
	"time"
)

type PostUsecase struct {
	PostRepository post.PostRepository
	UserRepository user.UserRepository
	ForumRepository forum.ForumRepository
	ThreadRepository thread.ThreadRepository
}

func NewPostUsecse(postRepository post.PostRepository, userRepository user.UserRepository, forumRepository forum.ForumRepository, threadRepository thread.ThreadRepository) *PostUsecase {
	return &PostUsecase{
		PostRepository: postRepository,
		UserRepository: userRepository,
		ForumRepository: forumRepository,
		ThreadRepository: threadRepository,
	}
}

func (pu *PostUsecase) CreatePosts(posts []*models.Post, thread *models.Thread) *errors.ErrorResponse {
	location, _ := time.LoadLocation("UTC")
	now := time.Now().In(location).Round(time.Microsecond)
	//now := time.Now()
	err := pu.PostRepository.Insert(posts, thread, &now)
	//switch {
	//case err.(*pq.Error).Code == errors.ConstParentNotFound:
	//	return errors.NewErrorResponse(errors.ErrParent, err)
	//case err.(*pq.Error).Code == errors.ConstUserNotFound:
	//	return errors.NewErrorResponse(errors.ErrNotFoundUser, err)
	//case err != nil:
	//	return errors.NewErrorResponse(errors.ErrInternal, err)
	//default:
	//	return nil
	//}
	if err != nil {
		if err.(*pq.Error).Code == errors.ConstParentNotFound {
			return errors.NewErrorResponse(errors.ErrParent, err)
		}
		if err.(*pq.Error).Code == errors.ConstUserNotFound {
			return errors.NewErrorResponse(errors.ErrNotFoundUser, err)
		}
		return errors.NewErrorResponse(errors.ErrInternal, err)
	}

	return nil
}

func (pu *PostUsecase) GetPosts(limit, since, desc, sort *string, id uint32) ([]*models.Post, *errors.ErrorResponse) {
	var posts []*models.Post
	var err error
	switch *sort {
	case "parent_tree":
		if *desc == "true" {
			posts, err = pu.PostRepository.GetPostsDescParentTree(limit, since, id)
		} else {
			posts, err = pu.PostRepository.GetPostsParentTree(limit, since, id)
		}
	case "tree":
		if *desc == "true" {
			posts, err = pu.PostRepository.GetPostsDescTree(limit, since, id)
		} else {
			posts, err = pu.PostRepository.GetPostsTree(limit, since, id)
		}
	default:
		if *desc == "true" {
		 posts, err = pu.PostRepository.GetPostsDesc(limit, since, id)
		} else {
			posts, err = pu.PostRepository.GetPosts(limit, since, id)
		}
	}

	if err != nil {
		return nil, errors.NewErrorResponse(errors.ErrInternal, err)
	}

	return posts, nil
}

func (pu *PostUsecase) GetPost(id int) (*models.Post, *errors.ErrorResponse) {
	post, err := pu.PostRepository.GetPost(id)

	switch {
	case err == sql.ErrNoRows:
		return nil, errors.NewErrorResponse(errors.ErrForumDoesNotExists, err)
	case err != nil:
		return nil, errors.NewErrorResponse(errors.ErrInternal, err)
	}
	return post, nil
}

func (pu *PostUsecase) UpdatePost(post *models.Post) *errors.ErrorResponse {
	err := pu.PostRepository.Update(post)

	switch {
	case err == sql.ErrNoRows:
		return errors.NewErrorResponse(errors.ErrForumDoesNotExists, err)
	case err != nil:
		return errors.NewErrorResponse(errors.ErrInternal, err)
	}
	return nil
}

func (pu *PostUsecase) GetFullPost(id int, related string) (*models.PostFull, *errors.ErrorResponse) {
	post, err := pu.PostRepository.GetPost(id)

	switch {
	case err == sql.ErrNoRows:
		return nil, errors.NewErrorResponse(errors.ErrForumDoesNotExists, err)
	case err != nil:
		return nil, errors.NewErrorResponse(errors.ErrInternal, err)
	}

	postFull := models.PostFull{}
	postFull.Post = post

	for _, arg := range strings.Split(string(related), ",") {
		switch arg {
		case "forum":
			forum, err := pu.ForumRepository.SelectBySlug(post.Forum)
			if err != nil {
				return nil, errors.NewErrorResponse(errors.ErrForumDoesNotExists, err)
			}
			postFull.Forum = forum
		case "thread":
			thread, err := pu.ThreadRepository.SelectByID(post.Thread)
			if err != nil {
				return nil, errors.NewErrorResponse(errors.ErrForumDoesNotExists, err)
			}
			postFull.Thread = thread
		case "user":
			user, err := pu.UserRepository.SelectByNickname(post.Author)
			if err != nil {
				return nil, errors.NewErrorResponse(errors.ErrForumDoesNotExists, err)
			}
			postFull.User = user
		}
	}

	return &postFull, nil
}


