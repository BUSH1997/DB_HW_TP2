package forum

import (
	"github.com/BUSH1997/DB_HW_TP2/app/models"
	"github.com/BUSH1997/DB_HW_TP2/app/tools"
)

type UseCase interface {
	CreateForum(forum models.Forum) (models.Forum, *models.CustomError)
	GetForum(slug string) (models.Forum, *models.CustomError)
	CreateThread(thread models.Thread) (models.Thread, *models.CustomError)
	GetForumUsers(slug string, filter tools.FilterUser) ([]models.User, *models.CustomError)
	GetForumThreads(slug string, filter tools.FilterThread) ([]models.Thread, *models.CustomError)
	CreatePosts(slugOrId string, post []models.Post) ([]models.Post, *models.CustomError)
	CreateVote(slugOrId string, vote models.Vote) (models.Thread, *models.CustomError)
	GetThread(slugOrId string) (models.Thread, *models.CustomError)
	GetPosts(slugOrId string, filter tools.FilterPosts)([]*models.Post, *models.CustomError)
	GetPost(id string, filter tools.FilterOnePost)(models.PostInfo, *models.CustomError)
	UpdateThread(slugOrId string, thread models.Thread) (models.Thread, *models.CustomError)
	UpdatePost(id string, post models.Post) (models.Post, *models.CustomError)
}
