package delivery

import (
	"encoding/json"
	"github.com/SanSanchezzz/tech-db-forum/internal/errors"
	"github.com/SanSanchezzz/tech-db-forum/internal/forum"
	"github.com/SanSanchezzz/tech-db-forum/internal/models"
	"github.com/SanSanchezzz/tech-db-forum/internal/post"
	"github.com/SanSanchezzz/tech-db-forum/internal/responser"
	"github.com/SanSanchezzz/tech-db-forum/internal/thread"
	"github.com/SanSanchezzz/tech-db-forum/internal/user"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
)

type ThreadHandler struct {
	threadUsecase thread.ThreadUsecase
	userUsecase user.UserUsecase
	forumUsecase forum.ForumUsecase
	postUsecase post.PostUsecase
}

func CreateThreadHandler(threadUsecase thread.ThreadUsecase, userUsecase user.UserUsecase, forumUsecase forum.ForumUsecase, postUsecase post.PostUsecase) *ThreadHandler {
	return &ThreadHandler{
		threadUsecase: threadUsecase,
		userUsecase: userUsecase,
		forumUsecase: forumUsecase,
		postUsecase: postUsecase,
	}
}

func (th *ThreadHandler) Configure(e *echo.Echo) {
	e.POST("api/thread/:slug_or_id/create", th.HandlerCreatePost())
	e.POST("api/thread/:slug_or_id/vote", th.HandlerVote())
	e.GET("api/thread/:slug_or_id/details", th.HandlerGetThread())
	e.POST("api/thread/:slug_or_id/details", th.HandlerThreadUpdate())
	e.GET("api/thread/:slug_or_id/posts", th.HandlerGetPosts())
}

func (th *ThreadHandler) HandlerGetPosts() echo.HandlerFunc {
	return func(context echo.Context) error {
		slugOrID := context.Param("slug_or_id")

		thread, err := th.threadUsecase.GetThread(slugOrID)
		if err != nil {
			return responser.RespondWithError(err, context)
		}
		if thread == nil {
			return responser.RespondWithError(errors.NewErrorResponse(errors.ErrForumDoesNotExists, nil), context)
		}

		limit := context.QueryParam("limit")
		since := context.QueryParam("since")
		desc := context.QueryParam("desc")
		sort := context.QueryParam("sort")
		if desc != "true" {
			desc = "false"
		}

		posts, err := th.postUsecase.GetPosts(&limit, &since, &desc, &sort, thread.ID)
		if err != nil {
			return responser.RespondWithError(err, context)
		}

		return context.JSON(http.StatusOK, posts)
	}
}

func (th *ThreadHandler) HandlerGetThread() echo.HandlerFunc {
	return func(context echo.Context) error {
		slugOrID := context.Param("slug_or_id")

		thread, err := th.threadUsecase.GetThread(slugOrID)
		if err != nil {
			return responser.RespondWithError(err, context)
		}
		if thread == nil {
			return responser.RespondWithError(errors.NewErrorResponse(errors.ErrForumDoesNotExists, nil), context)
		}

		return context.JSON(http.StatusOK, thread)
	}
}

func (th *ThreadHandler) HandlerThreadUpdate() echo.HandlerFunc {
	type Request struct {
		Title string `json:"title"`
		Message string `json:"message"`
	}
	return func(context echo.Context) error {
		slugOrID := context.Param("slug_or_id")

		req := &Request{}
		e := context.Bind(req)
		if e != nil {
			return e
		}

		thread, err := th.threadUsecase.GetThread(slugOrID)
		if err != nil {
			return responser.RespondWithError(err, context)
		}
		if thread == nil {
			return responser.RespondWithError(errors.NewErrorResponse(errors.ErrForumDoesNotExists, nil), context)
		}

		if req.Message == "" && req.Title == "" {
			return context.JSON(http.StatusOK, thread)
		}

		if req.Title != "" {
			thread.Title = req.Title
		}

		if req.Message != "" {
			thread.Message = req.Message
		}

		err = th.threadUsecase.UpdateThread(thread)
		if err != nil {
			return responser.RespondWithError(err, context)
		}

		return context.JSON(http.StatusOK, thread)
	}
}

func (th *ThreadHandler) HandlerVote() echo.HandlerFunc {
	type Request struct {
		Nickname string `json:"nickname"`
		Voice int `json:"voice"`
	}
	return func(context echo.Context) error {
		slugOrID := context.Param("slug_or_id")

		req := &Request{}
		context.Bind(req)

		thread, err := th.threadUsecase.VoteThread(req.Voice, req.Nickname, slugOrID)
		if err != nil {
			return responser.RespondWithError(err, context)
		}

		return context.JSON(http.StatusOK, thread)
	}
}

func (th *ThreadHandler) HandlerCreatePost() echo.HandlerFunc {
	return func(context echo.Context) error {
		slugOrID := context.Param("slug_or_id")

		thread, err := th.threadUsecase.GetThread(slugOrID)
		if err != nil {
			return responser.RespondWithError(err, context)
		}
		if thread == nil {
			return responser.RespondWithError(errors.NewErrorResponse(errors.ErrForumDoesNotExists, nil), context)
		}

		result, e := ioutil.ReadAll(context.Request().Body)
		if e != nil {
			return responser.RespondWithError(errors.NewErrorResponse(errors.ErrInternal, nil), context)
		}

		var posts []*models.Post
		e = json.Unmarshal(result, &posts)
		if e != nil {
			return responser.RespondWithError(errors.NewErrorResponse(errors.ErrInternal, nil), context)
		}

		if len(posts) == 0 {
			return context.JSON(http.StatusCreated, posts)
		}

		err = th.postUsecase.CreatePosts(posts, thread)
		if err != nil {
			return responser.RespondWithError(err, context)
		}

		return context.JSON(http.StatusCreated, posts)
	}
}
