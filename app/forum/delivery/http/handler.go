package http

import (
	"github.com/BUSH1997/DB_HW_TP2/app/forum"
	"github.com/BUSH1997/DB_HW_TP2/app/models"
	"github.com/BUSH1997/DB_HW_TP2/app/tools"
	"github.com/labstack/echo/v4"
	"net/http"
)

type ForumHandler struct {
	useCase forum.UseCase
}

func NewForumHandler(useCase forum.UseCase) *ForumHandler {
	return &ForumHandler{
		useCase: useCase,
	}
}

func (fh *ForumHandler) CreateForum(ctx echo.Context) error {
	var newForum models.Forum

	if err := ctx.Bind(&newForum); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	forum, err := fh.useCase.CreateForum(newForum)
	if err != nil {
		if err.Message == models.NoUser {
			return ctx.JSON(http.StatusNotFound, err)
		}
		if err.Message == models.ConflictData {
			return ctx.JSON(http.StatusConflict, forum)
		}
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, forum)
}

func (fh *ForumHandler) GetForumDetails(ctx echo.Context) error {
	slug := ctx.Param("slug")

	forum, err := fh.useCase.GetDetailsForum(slug)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, err)
	}

	return ctx.JSON(http.StatusOK, forum)
}

func (fh *ForumHandler) CreateThread(ctx echo.Context) error {
	var newThread models.Thread

	if err := ctx.Bind(&newThread); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	newThread.Forum = ctx.Param("slug")

	thread, err := fh.useCase.CreateThread(newThread)
	if err != nil {
		if err.Message == models.NoUser {
			return ctx.JSON(http.StatusNotFound, err)
		}
		if err.Message == models.ConflictData {
			return ctx.JSON(http.StatusConflict, thread)
		}
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, thread)
}

func (fh *ForumHandler) GetUsersForum(ctx echo.Context) error {
	slug := ctx.Param("slug")
	filter := tools.ParseQueryFilterUser(ctx)

	users, err := fh.useCase.GetUsersForum(slug, filter)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, err)
	}

	return ctx.JSON(http.StatusOK, users)
}

func (fh *ForumHandler) GetForumThreads(ctx echo.Context) error {
	slug := ctx.Param("slug")
	filter := tools.ParseQueryFilterThread(ctx)

	users, err := fh.useCase.GetForumThreads(slug, filter)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, err)
	}

	return ctx.JSON(http.StatusOK, users)
}

func (fh *ForumHandler) CreatePosts(ctx echo.Context) error {
	var newPosts []models.Post

	if err := ctx.Bind(&newPosts); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	slugOrId := ctx.Param("slug_or_id")
	posts, err := fh.useCase.CreatePosts(slugOrId, newPosts)
	if err != nil {
		if err.Message == models.NoUser {
			return ctx.JSON(http.StatusNotFound, err)
		}
		if err.Message == models.ConflictData {
			return ctx.JSON(http.StatusConflict, err)
		}
		if err.Message == models.BadParentPost {
			return ctx.JSON(http.StatusConflict, err)
		}
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, posts)
}

func (fh *ForumHandler) Vote(ctx echo.Context) error {
	var newVoice models.Vote

	if err := ctx.Bind(&newVoice); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}
	slugOrId := ctx.Param("slug_or_id")

	thread, err := fh.useCase.CreateVote(slugOrId, newVoice)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, err)
	}

	return ctx.JSON(http.StatusOK, thread)
}

func (fh *ForumHandler) Details(ctx echo.Context) error {
	slugOrId := ctx.Param("slug_or_id")
	thread, err := fh.useCase.GetThreadDetails(slugOrId)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, err)
	}

	return ctx.JSON(http.StatusOK, thread)
}

func (fh *ForumHandler) GetPosts (ctx echo.Context) error {
	filter := tools.ParseQueryFilterPost(ctx)
	slugOrId := ctx.Param("slug_or_id")

	posts, err := fh.useCase.GetPosts(slugOrId, filter)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, err)
	}

	return ctx.JSON(http.StatusOK, posts)
}

func (fh *ForumHandler) UpdateThread (ctx echo.Context) error {
	var newThread models.Thread

	if err := ctx.Bind(&newThread); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}
	slugOrId := ctx.Param("slug_or_id")

	thread, err := fh.useCase.UpdateThread(slugOrId, newThread)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, err)
	}

	return ctx.JSON(http.StatusOK, thread)
}

func (fh *ForumHandler) GetOnePost (ctx echo.Context) error {
	id := ctx.Param("id")
	filter := tools.ParseQueryFilterOnePost(ctx)

	post, err := fh.useCase.GetPost(id, filter)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, err)
	}

	return ctx.JSON(http.StatusOK, post)
}

func (fh *ForumHandler) UpdatePost (ctx echo.Context) error {
	var postInfo models.Post

	if err := ctx.Bind(&postInfo); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	id := ctx.Param("id")
	post, err := fh.useCase.UpdatePost(id, postInfo)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, err)
	}

	return ctx.JSON(http.StatusOK, post)
}
