package delivery

import (
	"github.com/SanSanchezzz/tech-db-forum/internal/responser"
	"github.com/SanSanchezzz/tech-db-forum/internal/service"
	"github.com/labstack/echo/v4"
	"net/http"
)

type ServiceHandler struct {
	serviceHandler service.ServiceUsecase
}

func CreateServiceHandler(serviceHandler service.ServiceUsecase) *ServiceHandler {
	return &ServiceHandler{
		serviceHandler: serviceHandler,
	}
}

func (uh *ServiceHandler) Configure(e *echo.Echo) {
	e.GET("api/service/status", uh.HandlerGetStatus())
	e.POST("api/service/clear", uh.HandlerClear())
}

func (uh *ServiceHandler) HandlerGetStatus() echo.HandlerFunc {
	return func(context echo.Context) error {
		service, err := uh.serviceHandler.GetStatus()
		if err != nil {
			return responser.RespondWithError(err, context)
		}

		return context.JSON(http.StatusOK, *service)
	}
}

func (uh *ServiceHandler) HandlerClear() echo.HandlerFunc {
	return func(context echo.Context) error {
		err := uh.serviceHandler.Clear()
		if err != nil {
			return responser.RespondWithError(err, context)
		}

		return context.JSON(http.StatusOK, nil)
	}
}

