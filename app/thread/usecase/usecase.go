package usecase

import (
	"github.com/BUSH1997/DB_HW_TP2/app/forum"
	"github.com/BUSH1997/DB_HW_TP2/app/models"
	"github.com/BUSH1997/DB_HW_TP2/app/thread"
	"github.com/BUSH1997/DB_HW_TP2/app/tools"
	"github.com/BUSH1997/DB_HW_TP2/app/user"
	"github.com/jackc/pgx"
	"strconv"
)

type UseCase struct {
	Repository      thread.Repository
	RepositoryUser  user.Repository
	RepositoryForum forum.Repository
}

func NewUseCase(repository thread.Repository, userRepository user.Repository, forumRepository forum.Repository) *UseCase {
	return &UseCase{Repository: repository, RepositoryUser: userRepository, RepositoryForum: forumRepository}
}

func (uc *UseCase) CreatePosts(slugOrId string, post []models.Post) ([]models.Post, *models.CustomError) {
	var thread models.Thread
	id, err := strconv.Atoi(slugOrId)
	if err != nil {
		thread, err = uc.Repository.GetThreadBySlug(slugOrId)
		if err != nil {
			if pgErr, ok := err.(pgx.PgError); ok && pgErr.Code == "23505" {
				return nil, &models.CustomError{Message: models.ConflictData}
			}

			if err == pgx.ErrNoRows {
				return nil, &models.CustomError{Message: models.NoUser}
			}

			return nil, &models.CustomError{Message: err.Error()}
		}

	} else {
		thread, err = uc.Repository.GetThreadById(id)
		if err != nil {
			if pgErr, ok := err.(pgx.PgError); ok && pgErr.Code == "23505" {
				return nil, &models.CustomError{Message: models.ConflictData}
			}

			if err == pgx.ErrNoRows {
				return nil, &models.CustomError{Message: models.NoUser}
			}

			return nil, &models.CustomError{Message: err.Error()}
		}
	}
	posts, err := uc.Repository.CreatePosts(int(thread.Id), thread.Forum, post)
	if err != nil {
		if pgErr, ok := err.(pgx.PgError); ok && pgErr.Code == "23505" {
			return nil, &models.CustomError{Message: models.ConflictData}
		}
		if pgErr, ok := err.(pgx.PgError); ok && pgErr.Code == models.PgxBadParentErrorCode {
			return nil, &models.CustomError{Message: models.BadParentPost}
		}
		if pgErr, ok := err.(pgx.PgError); ok && pgErr.Code == "23503" {
			return nil, &models.CustomError{Message: models.NoUser}
		}
		return nil, &models.CustomError{Message: err.Error()}
	}

	return posts, nil
}

func (uc *UseCase) CreateVote(slugOrId string, vote models.Vote) (models.Thread, *models.CustomError) {
	var thread models.Thread
	err := uc.Repository.CreateVoteBySlugOrId(slugOrId, vote)
	if err != nil {
		if pgErr, ok := err.(pgx.PgError); ok && pgErr.Code == "23503" {
			return models.Thread{}, &models.CustomError{Message: models.NoUser}
		}
		if pgErr, ok := err.(pgx.PgError); ok && pgErr.Code == "23505" {
			err = uc.Repository.UpdateVoteBySlugOrId(slugOrId, vote)
			if err != nil {
				return models.Thread{}, &models.CustomError{Message: err.Error()}
			}
			thread, err = uc.Repository.GetThreadBySlugOrId(slugOrId)
			if err != nil {
				return models.Thread{}, &models.CustomError{Message: err.Error()}
			}

			return thread, nil
		}
		return models.Thread{}, &models.CustomError{Message: err.Error()}
	}

	thread, err = uc.Repository.GetThreadBySlugOrId(slugOrId)
	if err != nil {
		return models.Thread{}, &models.CustomError{Message: err.Error()}
	}

	return thread, nil
}

func (uc *UseCase) GetThreadDetails(slugOrId string) (models.Thread, *models.CustomError) {
	thread, err := uc.Repository.GetThreadBySlugOrId(slugOrId)
	if err != nil {
		return models.Thread{}, &models.CustomError{Message: err.Error()}
	}
	return thread, nil
}

func (uc *UseCase) GetPosts(slugOrId string, filter tools.FilterPosts) ([]*models.Post, *models.CustomError) {
	var result []*models.Post
	var err error

	switch filter.Sort {
	case tools.SortParamFlatDefault:
		result, err = uc.Repository.GetPostsFlatSlugOrId(slugOrId, filter)
	case tools.SortParamParentTree:
		result, err = uc.Repository.GetPostsParentTreeSlugOrId(slugOrId, filter)
	case tools.SortParamTree:
		result, err = uc.Repository.GetPostsTreeSlugOrId(slugOrId, filter)
	}
	if err != nil {
		return nil, &models.CustomError{Message: err.Error()}
	}

	if len(result) == 0 {
		_, err := uc.Repository.GetThreadBySlugOrId(slugOrId)
		if err != nil {
			return nil, &models.CustomError{Message: models.NoUser}
		}
		return []*models.Post{}, nil
	}

	return result, nil
}

func (uc *UseCase) UpdateThread(slugOrId string, thread models.Thread) (models.Thread, *models.CustomError) {
	thread, err := uc.Repository.UpdateThread(slugOrId, thread)
	if err != nil {
		if pgErr, ok := err.(pgx.PgError); ok && pgErr.Code == "23505" {
			return models.Thread{}, &models.CustomError{Message: models.ConflictData}
		}
		if err == pgx.ErrNoRows {
			return models.Thread{}, &models.CustomError{Message: models.NoUser}
		}

		return models.Thread{}, &models.CustomError{Message: err.Error()}
	}
	return thread, nil
}

func (uc *UseCase) GetPost(id string, filter tools.FilterOnePost) (models.PostInfo, *models.CustomError) {
	var result models.PostInfo

	idNum, err := strconv.Atoi(id)
	if err != nil {
		return models.PostInfo{}, &models.CustomError{Message: err.Error()}
	}

	post, err := uc.Repository.GetPostById(idNum)
	if err != nil {
		return models.PostInfo{}, &models.CustomError{Message: err.Error()}
	}
	result.Post = post

	if filter.User {
		user, err := uc.RepositoryUser.GetUser(post.Author)
		if err != nil {
			return models.PostInfo{}, &models.CustomError{Message: err.Error()}
		}
		result.Author = &user
	}

	if filter.Thread {
		thread, err := uc.Repository.GetThreadById(int(post.Thread))
		if err != nil {
			return models.PostInfo{}, &models.CustomError{Message: err.Error()}
		}
		result.Thread = &thread
	}

	if filter.Forum {
		forum, err := uc.RepositoryForum.GetForumBySlug(post.Forum)
		if err != nil {
			return models.PostInfo{}, &models.CustomError{Message: err.Error()}
		}
		result.Forum = &forum
	}

	return result, nil
}

func (uc *UseCase) UpdatePost(id string, post models.Post) (models.Post, *models.CustomError) {
	idNum, err := strconv.Atoi(id)
	if err != nil {
		return models.Post{}, &models.CustomError{Message: err.Error()}
	}
	if post.Message == "" {
		post, err = uc.Repository.GetPostById(idNum)
	} else {
		post, err = uc.Repository.UpdatePost(idNum, post)
	}
	if err != nil {
		return models.Post{}, &models.CustomError{Message: err.Error()}
	}

	return post, nil
}
