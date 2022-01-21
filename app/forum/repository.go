package forum

import (
	"github.com/BUSH1997/DB_HW_TP2/app/models"
	"github.com/BUSH1997/DB_HW_TP2/app/tools"
)

type Repository interface {
	CreateForum(forum models.Forum) (models.Forum, error)
	GetForumBySlug(slug string) (models.Forum, error)
	GetForum(slug string) (models.Forum, error)
	CreateThread(thread models.Thread) (models.Thread, error)
	GetForumUsers(slug string, filter tools.FilterUser) ([]models.User, error)
	GetForumThreads(slug string, filter tools.FilterThread) ([]models.Thread, error)
	CreatePosts(threadId int, threadForum string, post []models.Post) ([]models.Post, error)
	GetThreadBySlug(slug string) (models.Thread, error)
	GetThreadById(id int) (models.Thread, error)
	GetThreadBySlugOrId(slugOrId string) (models.Thread, error)
	CreateVoteBySlugOrId(slugOrId string, vote models.Vote) error
	UpdateVoteBySlugOrId(slugOrId string, vote models.Vote) error
	GetPostById(id int)(models.Post, error)
	UpdatePost(id int, post models.Post)(models.Post, error)
	GetPostsFlatSlugOrId(slugOrId string, posts tools.FilterPosts)([]*models.Post, error)
	GetPostsTreeSlugOrId(slugOrId string, posts tools.FilterPosts)([]*models.Post, error)
	GetPostsParentTreeSlugOrId(slugOrId string, posts tools.FilterPosts)([]*models.Post, error)
	UpdateThread(slugOrId string, thread models.Thread) (models.Thread, error)
}
