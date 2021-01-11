package delivery

import (
	"github.com/SanSanchezzz/tech-db-forum/internal/forum"
	"github.com/SanSanchezzz/tech-db-forum/internal/post"
	"github.com/SanSanchezzz/tech-db-forum/internal/responser"
	"github.com/SanSanchezzz/tech-db-forum/internal/thread"
	"github.com/SanSanchezzz/tech-db-forum/internal/user"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type PostHandler struct {
	userUsecase user.UserUsecase
	forumUsecase forum.ForumUsecase
	threadUsecase thread.ThreadUsecase
	postUsecase post.PostUsecase
}

func CreatePostHandler(userUsecase user.UserUsecase, forumUsecase forum.ForumUsecase, threadUsecase thread.ThreadUsecase, postUsecase post.PostUsecase) *PostHandler {
	return &PostHandler{
		userUsecase: userUsecase,
		forumUsecase: forumUsecase,
		threadUsecase: threadUsecase,
		postUsecase: postUsecase,
	}
}

func (ph *PostHandler) Configure(e *echo.Echo) {
	//e.POST("api/forum/create", fh.HandlerCreateForum())
	e.GET("api/post/:id/details", ph.HandlerGetPost())
	e.POST("api/post/:id/details", ph.HandlerUpdatePost())
	//e.POST("api/forum/:slug/create", fh.HandlerCreateThread())
	//e.GET("api/forum/:slug/threads", fh.HandlerGetThreads())

}

func (ph * PostHandler) HandlerUpdatePost() echo.HandlerFunc {
	type Request struct {
		Message string `json:"message"`
	}
	return func(context echo.Context) error {
		id, _ := strconv.Atoi(context.Param("id"))

		req := &Request{}
		context.Bind(req)

		post, e := ph.postUsecase.GetPost(id)
		if e != nil {
			return responser.RespondWithError(e, context)
		}

		if req.Message == post.Message || req.Message == "" {
			return context.JSON(http.StatusOK, post)
		}

		post.Message = req.Message

		e = ph.postUsecase.UpdatePost(post)
		if e != nil {
			return responser.RespondWithError(e, context)
		}

		return context.JSON(http.StatusOK, post)
	}
}

func (ph * PostHandler) HandlerGetPost() echo.HandlerFunc {
	return func(context echo.Context) error {
		id, _ := strconv.Atoi(context.Param("id"))
		related := context.QueryParam("related")

		resp, e := ph.postUsecase.GetFullPost(id, related)
		if e != nil {
			return responser.RespondWithError(e, context)
		}


		return context.JSON(http.StatusOK, resp)
	}
}

//func (fh *PostHandler) HandlerCreateForum() echo.HandlerFunc {
//	type Request struct {
//		Slug string `json:"slug"`
//		User string `json:"user"`
//		Title string `json:"title"`
//	}
//	return func(context echo.Context) error {
//		req := &Request{}
//		context.Bind(req)
//
//		user, err := fh.userUsecase.GetByNickname(req.User)
//
//		if err != nil {
//			return responser.RespondWithError(err, context)
//		}
//
//		forum, err := fh.forumUsecase.GetForum(req.Slug)
//
//		if err != nil {
//			return responser.RespondWithError(err, context)
//		}
//
//		if forum != nil {
//			return context.JSON(http.StatusConflict, forum)
//		}
//
//		newForum := &models.Forum{
//			Slug: req.Slug,
//			User: user.Nickname,
//			Title: req.Title,
//		}
//
//		req = &Request{
//			User:  newForum.User,
//			Title: newForum.Title,
//			Slug:  newForum.Slug,
//		}
//
//		err = fh.forumUsecase.CreateForum(newForum)
//
//		if err != nil {
//			return responser.RespondWithError(err, context)
//		}
//
//		return context.JSON(http.StatusCreated, req)
//	}
//}
//
//func (fh * ForumHandler) HandlerGetForum() echo.HandlerFunc {
//	return func(context echo.Context) error {
//		slug := context.Param("slug")
//
//		forum, err := fh.forumUsecase.GetForum(slug)
//
//		if err != nil {
//			return responser.RespondWithError(err, context)
//		}
//
//		if forum == nil {
//			return context.JSON(http.StatusNotFound, nil)
//		}
//
//		return context.JSON(http.StatusOK, forum)
//	}
//}
//
//func (fh *ForumHandler) HandlerCreateThread() echo.HandlerFunc {
//	type Request struct {
//		Title string `json:"title"`
//		Nickname string `json:"author"`
//		Forum string `json:"forum"`
//		Message string `json:"message"`
//		Slug string `json:"slug"`
//		Created time.Time `json:"created"`
//	}
//	return func(context echo.Context) error {
//		slug := context.Param("slug")
//
//		req := &Request{}
//		context.Bind(req)
//
//		user, err := fh.userUsecase.GetByNickname(req.Nickname)
//
//		if err != nil {
//			return responser.RespondWithError(err, context)
//		}
//
//		forum, err := fh.forumUsecase.GetForum(slug)
//
//		if err != nil {
//			return responser.RespondWithError(err, context)
//		}
//
//		if forum == nil {
//			return responser.RespondWithError(errors.NewErrorResponse(errors.ErrForumDoesNotExists, nil), context)
//		}
//
//		thread, err := fh.threadUsecase.GetThread(req.Slug)
//
//		if err != nil {
//			return responser.RespondWithError(err, context)
//		}
//
//		if thread != nil {
//			return context.JSON(http.StatusConflict, thread)
//		}
//
//		thread = &models.Thread{
//			Title: req.Title,
//			Nickname: user.Nickname,
//			Forum: forum.Slug,
//			Message: req.Message,
//			Created: req.Created,
//		}
//
//		if forum.Slug != req.Slug {
//			thread.Slug = req.Slug
//		}
//
//		id, err := fh.threadUsecase.CreateThread(thread)
//
//		if err != nil {
//			return responser.RespondWithError(err, context)
//		}
//		thread.ID = id
//
//		return context.JSON(http.StatusCreated, thread)
//	}
//}
//
//func (fh * ForumHandler) HandlerGetThreads() echo.HandlerFunc {
//	return func(context echo.Context) error {
//		slug := context.Param("slug")
//
//		forum, err := fh.forumUsecase.GetForum(slug)
//
//		if err != nil {
//			return responser.RespondWithError(err, context)
//		}
//
//		if forum == nil {
//			return context.JSON(http.StatusNotFound, nil)
//		}
//
//		limit, _ := strconv.Atoi(context.QueryParam("limit"))
//		since := context.QueryParam("since")
//		descParam := context.QueryParam("desc")
//		desc := false
//		if descParam == "true" {
//			desc = true
//		}
//		//print(since)
//
//		threads, err := fh.threadUsecase.GetThreadsByParam(limit, desc, &slug, &since)
//		if err != nil {
//			return responser.RespondWithError(err, context)
//		}
//
//		return context.JSON(http.StatusOK, threads)
//	}
//}
