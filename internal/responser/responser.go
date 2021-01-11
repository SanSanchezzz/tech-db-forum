package responser

import (
	"github.com/SanSanchezzz/tech-db-forum/internal/errors"
	"github.com/labstack/echo/v4"
)

func RespondWithError(err *errors.ErrorResponse, ctx echo.Context) error {

	return ctx.JSON(err.StatusCode, err.JsonError)
}