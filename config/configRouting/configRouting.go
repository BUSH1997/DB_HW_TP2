package configRouting

import (
	forumHandler "github.com/BUSH1997/DB_HW_TP2/app/forum/delivery/http"
	serviceHandler "github.com/BUSH1997/DB_HW_TP2/app/service/delivery/http"
	userHandler "github.com/BUSH1997/DB_HW_TP2/app/user/delivery/http"
	"github.com/labstack/echo/v4"
)

type ServerConfigRouting struct {
	ForumHandler forumHandler.ForumHandler
	UserHandler userHandler.UserHandler
	ServiceHandler serviceHandler.ServiceHandler
}

func (cr *ServerConfigRouting) ConfigRouting(router *echo.Echo) {
	router.POST("user/:nickname/create", cr.UserHandler.SignUpUser)
	router.GET("user/:nickname/profile", cr.UserHandler.GetUser)
	router.POST("user/:nickname/profile", cr.UserHandler.UpdateUser)
	router.POST("forum/create", cr.ForumHandler.CreateForum)
	router.GET("forum/:slug/details", cr.ForumHandler.GetForumDetails)
	router.POST("forum/:slug/create", cr.ForumHandler.CreateThread)
	router.GET("forum/:slug/users", cr.ForumHandler.GetUsersForum)
	router.GET("forum/:slug/threads", cr.ForumHandler.GetForumThreads)
	router.POST("thread/:slug_or_id/create", cr.ForumHandler.CreatePosts)
	router.POST("thread/:slug_or_id/vote", cr.ForumHandler.Vote)
	router.GET("thread/:slug_or_id/details", cr.ForumHandler.Details)
	router.GET("thread/:slug_or_id/posts", cr.ForumHandler.GetPosts)
	router.POST("thread/:slug_or_id/details", cr.ForumHandler.UpdateThread)
	router.GET("post/:id/details", cr.ForumHandler.GetOnePost)
	router.POST("post/:id/details", cr.ForumHandler.UpdatePost)
	router.GET("service/status", cr.ServiceHandler.Status)
	router.POST("service/clear", cr.ServiceHandler.Clear)
}
