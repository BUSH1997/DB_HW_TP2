package usecase

import (
	"github.com/BUSH1997/DB_HW_TP2/app/models"
	"github.com/BUSH1997/DB_HW_TP2/app/service"
)

type UseCase struct {
	repo service.Repository
}

func NewUseCase(repo service.Repository) *UseCase {
	return &UseCase{
		repo: repo,
	}
}

func (uc *UseCase) GetStatus() (models.Status, error)  {
	return uc.repo.GetStatus()
}

func (uc *UseCase) Clear() error {
	return uc.repo.Clear()
}