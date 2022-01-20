package thread

import (
	"github.com/BUSH1997/DB_HW_TP2/app/models"
	"github.com/BUSH1997/DB_HW_TP2/app/tools"
)

type UseCase interface {
	CreatePosts(slugOrId string, post []models.Post) ([]models.Post, *models.CustomError)
	CreateVote(slugOrId string, vote models.Vote) (models.Thread, *models.CustomError)
	GetThreadDetails(slugOrId string) (models.Thread, *models.CustomError)
	GetPosts(slugOrId string, filter tools.FilterPosts)([]*models.Post, *models.CustomError)
	GetPost(id string, filter tools.FilterOnePost)(models.PostInfo, *models.CustomError)
	UpdateThread(slugOrId string, thread models.Thread) (models.Thread, *models.CustomError)
	UpdatePost(id string, post models.Post) (models.Post, *models.CustomError)
}
