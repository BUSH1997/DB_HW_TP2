package http

import (
	"github.com/BUSH1997/DB_HW_TP2/app/service"
	"github.com/labstack/echo/v4"
	"net/http"
)

type ServiceHandler struct {
	useCase service.UseCase
}

func NewServiceHandler(useCase service.UseCase) *ServiceHandler {
	return &ServiceHandler{
		useCase: useCase,
	}
}

func (sh *ServiceHandler) Status(ctx echo.Context) error  {
	status, err := sh.useCase.GetStatus()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, status)
}

func (sh *ServiceHandler) Clear(ctx echo.Context) error {
	err := sh.useCase.Clear()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	return ctx.NoContent(http.StatusOK)
}
