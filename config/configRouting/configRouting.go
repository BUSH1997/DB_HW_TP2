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
	router.POST("api/user/:nickname/create", cr.UserHandler.SignUpUser)
	router.GET("api/user/:nickname/profile", cr.UserHandler.GetUser)
	router.POST("api/user/:nickname/profile", cr.UserHandler.UpdateUser)
	router.POST("api/forum/create", cr.ForumHandler.CreateForum)
	router.GET("api/forum/:slug/details", cr.ForumHandler.GetForumDetails)
	router.POST("api/forum/:slug/create", cr.ForumHandler.CreateThread)
	router.GET("api/forum/:slug/users", cr.ForumHandler.GetUsersForum)
	router.GET("api/forum/:slug/threads", cr.ForumHandler.GetForumThreads)
	router.POST("api/thread/:slug_or_id/create", cr.ForumHandler.CreatePosts)
	router.POST("api/thread/:slug_or_id/vote", cr.ForumHandler.Vote)
	router.GET("api/thread/:slug_or_id/details", cr.ForumHandler.Details)
	router.GET("api/thread/:slug_or_id/posts", cr.ForumHandler.GetPosts)
	router.POST("api/thread/:slug_or_id/details", cr.ForumHandler.UpdateThread)
	router.GET("api/post/:id/details", cr.ForumHandler.GetOnePost)
	router.POST("api/post/:id/details", cr.ForumHandler.UpdatePost)
	router.GET("api/service/status", cr.ServiceHandler.Status)
	router.POST("api/service/clear", cr.ServiceHandler.Clear)
}
