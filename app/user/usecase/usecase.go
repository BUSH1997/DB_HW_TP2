package usecase

import (
	"github.com/BUSH1997/DB_HW_TP2/app/models"
	"github.com/BUSH1997/DB_HW_TP2/app/user"
	"github.com/jackc/pgx"
)

type UseCase struct {
	repo user.Repository
}

func NewUseCase(repo user.Repository) *UseCase {
	return &UseCase{
		repo: repo,
	}
}

func (uc *UseCase) CreateUser(user models.User) ([]models.User, error) {
	var resultArray []models.User
	result, err := uc.repo.AddUser(user)
	if err != nil {
		if pgErr, ok := err.(pgx.PgError); ok && pgErr.Code == "23505" {
			result, err1 := uc.repo.GetUsersByNicknameOrEmail(user.Nickname, user.Email)
			if err1 != nil {
				return nil, err1
			}
			return result, err
		}
	}

	resultArray = append(resultArray, result)
	return resultArray, err
}

func (uc *UseCase) GetUserProfile(nickname string) (models.User, *models.CustomError) {
	user, err := uc.repo.GetUser(nickname)
	if err == pgx.ErrNoRows {
		return models.User{}, &models.CustomError{Message: models.NoUser}
	}
	return user, nil
}

func (uc *UseCase) UpdateUserProfile(user models.User) (models.User, *models.CustomError) {
	userNew, err := uc.repo.UpdateUser(user)
	if err != nil {
		if pgErr, ok := err.(pgx.PgError); ok && pgErr.Code == "23505" {
			return models.User{}, &models.CustomError{Message: models.ConflictData}
		}
		if err == pgx.ErrNoRows {
			return models.User{}, &models.CustomError{Message: models.NoUser}
		}

		return models.User{}, &models.CustomError{Message: err.Error()}
	}
	return userNew, nil
}
