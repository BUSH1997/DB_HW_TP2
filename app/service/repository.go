package service

import "github.com/BUSH1997/DB_HW_TP2/app/models"

type Repository interface {
	GetStatus() (models.Status, error)
	Clear() error
}