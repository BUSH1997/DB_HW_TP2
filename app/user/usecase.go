package user

import "github.com/BUSH1997/DB_HW_TP2/app/models"

type UseCase interface {
	CreateUser(user models.User) ([]models.User, error)
	GetUserProfile(nickname string) (models.User, *models.CustomError)
	UpdateUserProfile(user models.User) (models.User, *models.CustomError)
}
