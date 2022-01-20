package http

import (
	"github.com/BUSH1997/DB_HW_TP2/app/models"
	"github.com/BUSH1997/DB_HW_TP2/app/user"
	"github.com/labstack/echo/v4"
	"net/http"
)

type UserHandler struct {
	useCase user.UseCase
}

func NewUserHandler(useCase user.UseCase) *UserHandler {
	return &UserHandler{
		useCase: useCase,
	}
}

func (uh *UserHandler) SignUpUser(ctx echo.Context) error {
	var newUser models.User

	if err := ctx.Bind(&newUser); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	newUser.Nickname = ctx.Param("nickname")
	users, err := uh.useCase.CreateUser(newUser)
	if err != nil {
		if users[0].Email == "" {
			return ctx.JSON(http.StatusInternalServerError, err.Error())
		}
		return ctx.JSON(http.StatusConflict, users)
	}

	return ctx.JSON(http.StatusCreated, users[0])
}

func (uh *UserHandler) GetUser(ctx echo.Context) error {
	nickname := ctx.Param("nickname")
	user, err := uh.useCase.GetUserProfile(nickname)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, err)
	}
	return ctx.JSON(http.StatusOK, user)
}

func (uh *UserHandler) UpdateUser(ctx echo.Context) error {
	var UserUpdate models.User

	if err := ctx.Bind(&UserUpdate); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	UserUpdate.Nickname = ctx.Param("nickname")

	user, err := uh.useCase.UpdateUserProfile(UserUpdate)
	if err != nil {
		if err.Message == models.NoUser {
			return ctx.JSON(http.StatusNotFound, err)
		}
		if err.Message == models.ConflictData {
			return ctx.JSON(http.StatusConflict, err)
		}
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, user)
}
