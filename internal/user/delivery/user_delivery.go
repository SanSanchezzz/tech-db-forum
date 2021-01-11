package delivery

import (
	"encoding/json"
	"github.com/SanSanchezzz/tech-db-forum/internal/errors"
	"github.com/SanSanchezzz/tech-db-forum/internal/models"
	"github.com/SanSanchezzz/tech-db-forum/internal/responser"
	"github.com/SanSanchezzz/tech-db-forum/internal/user"
	"github.com/labstack/echo/v4"
	"net/http"
)

type UserHandler struct {
	userUsecase user.UserUsecase
}

func CreateUserHandler(userUsecase user.UserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
	}
}

func (uh *UserHandler) Configure(e *echo.Echo) {
	e.POST("api/user/:nickname/create", uh.HandlerCreateUser())
	e.GET("api/user/:nickname/profile", uh.HandlerGetProfile())
	e.POST("api/user/:nickname/profile", uh.HandlerUpdateProfile())
}

func (uh *UserHandler) HandlerCreateUser() echo.HandlerFunc {
	type Request struct {
		Fullname	string `json:"fullname"`
		About		string `json:"about"`
		Email		string `json:"email"`
	}
	return func(context echo.Context) error {
		nickname := context.Param("nickname")
		req := &Request{}
		_ = context.Bind(req)

		oldUsers, err := uh.userUsecase.Create(nickname, req.Fullname, req.About, req.Email)

		if err != nil {
			return responser.RespondWithError(err, context)
		}

		if oldUsers != nil {
			context.Response().WriteHeader(http.StatusConflict)
			response, _ := json.Marshal(oldUsers)
			_, e := context.Response().Write(response)
			return e
		}

		newUser := &models.User{
			Fullname: req.Fullname,
			About: req.About,
			Email: req.Email,
			Nickname: nickname,
		}

		return context.JSON(http.StatusCreated, newUser)
	}
}

func (uh *UserHandler) HandlerGetProfile() echo.HandlerFunc {
	return func (context echo.Context) error {
		nickname := context.Param("nickname")

		user, err := uh.userUsecase.GetByNickname(nickname)
		if err != nil {
			return responser.RespondWithError(err, context)
		}

		return context.JSON(http.StatusOK, user)
	}
}

func (uh *UserHandler) HandlerUpdateProfile() echo.HandlerFunc {
	type Request struct {
		Fullname	string `json:"fullname"`
		About		string `json:"about"`
		Email		string `json:"email"`
	}
	return func(context echo.Context) error {
		nickname := context.Param("nickname")
		req := &Request{}
		err := context.Bind(req)

		if err != nil {
			return responser.RespondWithError(errors.NewErrorResponse(errors.ErrInternal, nil), context)
		}

		user, e := uh.userUsecase.GetByNickname(nickname)

		if e != nil {
			return responser.RespondWithError(e, context)
		}

		if req.Email == "" && req.About == "" && req.Fullname == "" {
			response, _ := json.Marshal(user)
			_, e := context.Response().Write(response)
			return e
		}

		flag, e := uh.userUsecase.CheckUserByEmail(req.Email)

		if e != nil {
			return responser.RespondWithError(e, context)
		}

		if flag == true {
			return responser.RespondWithError(errors.NewErrorResponse(errors.ErrEmailExists, nil), context)
		}

		if req.Email != "" {
			user.Email = req.Email
		}
		if req.Fullname != "" {
			user.Fullname = req.Fullname
		}
		if req.About != "" {
			user.About = req.About
		}

		e = uh.userUsecase.UpdateProfile(user)

		if e != nil {
			return responser.RespondWithError(e, context)
		}

		response, _ := json.Marshal(user)
		_, err = context.Response().Write(response)
		context.Response().WriteHeader(http.StatusOK)
		return err
	}
}
