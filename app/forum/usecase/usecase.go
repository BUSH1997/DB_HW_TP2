package usecase

import (
	"github.com/BUSH1997/DB_HW_TP2/app/forum"
	"github.com/BUSH1997/DB_HW_TP2/app/models"
	"github.com/BUSH1997/DB_HW_TP2/app/tools"
	"github.com/BUSH1997/DB_HW_TP2/app/user"
	"github.com/jackc/pgx"
	"strconv"
)

type UseCase struct {
	forumRepo forum.Repository
	userRepo user.Repository
}

func NewUseCase(forumRepo forum.Repository, userRepo user.Repository) *UseCase {
	return &UseCase{
		forumRepo: forumRepo,
		userRepo: userRepo,
	}
}

func (uc *UseCase) CreateForum(forumGet models.Forum) (models.Forum, *models.CustomError) {
	forum, err := uc.forumRepo.AddForum(forumGet)
	if err != nil {
		if pgErr, ok := err.(pgx.PgError); ok && pgErr.Code == models.PgxNoFoundFieldErrorCode {
			return models.Forum{}, &models.CustomError{Message: models.NoUser}
		}
		if pgErr, ok := err.(pgx.PgError); ok && pgErr.Code == models.PgxUniqErrorCode {
			forum, err = uc.forumRepo.GetForumBySlug(forumGet.Slug)
			if err != nil {
				return models.Forum{}, &models.CustomError{Message: err.Error()}
			}
			return forum, &models.CustomError{Message: models.ConflictData}
		}
		return models.Forum{}, &models.CustomError{Message: err.Error()}
	}

	return forum, nil
}

func (uc *UseCase) GetDetailsForum(slug string) (models.Forum, *models.CustomError) {
	forum, err := uc.forumRepo.GetDetailsForum(slug)
	if err == pgx.ErrNoRows {
		return models.Forum{}, &models.CustomError{Message: models.NoSlug}
	}
	return forum, nil
}

func (uc *UseCase) CreateThread(threadGet models.Thread) (models.Thread, *models.CustomError) {
	var randomSlug bool
	if threadGet.Slug == "" {
		randomSlug = true
	}
	thread, err := uc.forumRepo.AddThread(threadGet)
	if err != nil {
		if pgErr, ok := err.(pgx.PgError); ok && pgErr.Code == models.PgxNoFoundFieldErrorCode {
			return models.Thread{}, &models.CustomError{Message: models.NoUser}
		}
		if pgErr, ok := err.(pgx.PgError); ok && pgErr.Code == models.PgxUniqErrorCode {
			thread, err = uc.forumRepo.GetThreadBySlug(threadGet.Slug)
			if err != nil {
				return models.Thread{}, &models.CustomError{Message: err.Error()}
			}
			return thread, &models.CustomError{Message: models.ConflictData}
		}
		return models.Thread{}, &models.CustomError{Message: err.Error()}
	}

	if randomSlug == true {
		thread.Slug = ""
	}
	return thread, nil
}

func (uc *UseCase) GetUsersForum(slug string, filter tools.FilterUser) ([]models.User, *models.CustomError) {
	users, err := uc.forumRepo.GetUsersForum(slug, filter)
	if users == nil {
		_, err = uc.forumRepo.GetForumBySlug(slug)
		if err != nil {
			return nil, &models.CustomError{Message: models.NoSlug}
		}
		return []models.User{}, nil
	}
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, &models.CustomError{Message: models.NoSlug}
		}
		return nil, &models.CustomError{Message: err.Error()}
	}

	return users, nil
}

func (uc *UseCase) GetForumThreads(slug string, filter tools.FilterThread) ([]models.Thread, *models.CustomError) {
	threads, err := uc.forumRepo.GetForumThreads(slug, filter)
	if threads == nil {
		_, err := uc.forumRepo.GetForumBySlug(slug)
		if err != nil {
			return nil, &models.CustomError{Message: err.Error()}
		}
		return []models.Thread{}, nil
	}
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, &models.CustomError{Message: models.NoSlug}
		}
		return nil, &models.CustomError{Message: err.Error()}
	}

	return threads, nil
}

func (uc *UseCase) CreatePosts(slugOrId string, post []models.Post) ([]models.Post, *models.CustomError) {
	var thread models.Thread
	id, err := strconv.Atoi(slugOrId)
	if err != nil {
		thread, err = uc.forumRepo.GetThreadBySlug(slugOrId)
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
		thread, err = uc.forumRepo.GetThreadById(id)
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
	posts, err := uc.forumRepo.CreatePosts(int(thread.Id), thread.Forum, post)
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
	err := uc.forumRepo.CreateVoteBySlugOrId(slugOrId, vote)
	if err != nil {
		if pgErr, ok := err.(pgx.PgError); ok && pgErr.Code == "23503" {
			return models.Thread{}, &models.CustomError{Message: models.NoUser}
		}
		if pgErr, ok := err.(pgx.PgError); ok && pgErr.Code == "23505" {
			err = uc.forumRepo.UpdateVoteBySlugOrId(slugOrId, vote)
			if err != nil {
				return models.Thread{}, &models.CustomError{Message: err.Error()}
			}
			thread, err = uc.forumRepo.GetThreadBySlugOrId(slugOrId)
			if err != nil {
				return models.Thread{}, &models.CustomError{Message: err.Error()}
			}

			return thread, nil
		}
		return models.Thread{}, &models.CustomError{Message: err.Error()}
	}

	thread, err = uc.forumRepo.GetThreadBySlugOrId(slugOrId)
	if err != nil {
		return models.Thread{}, &models.CustomError{Message: err.Error()}
	}

	return thread, nil
}

func (uc *UseCase) GetThreadDetails(slugOrId string) (models.Thread, *models.CustomError) {
	thread, err := uc.forumRepo.GetThreadBySlugOrId(slugOrId)
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
		result, err = uc.forumRepo.GetPostsFlatSlugOrId(slugOrId, filter)
	case tools.SortParamParentTree:
		result, err = uc.forumRepo.GetPostsParentTreeSlugOrId(slugOrId, filter)
	case tools.SortParamTree:
		result, err = uc.forumRepo.GetPostsTreeSlugOrId(slugOrId, filter)
	}
	if err != nil {
		return nil, &models.CustomError{Message: err.Error()}
	}

	if len(result) == 0 {
		_, err := uc.forumRepo.GetThreadBySlugOrId(slugOrId)
		if err != nil {
			return nil, &models.CustomError{Message: models.NoUser}
		}
		return []*models.Post{}, nil
	}

	return result, nil
}

func (uc *UseCase) UpdateThread(slugOrId string, thread models.Thread) (models.Thread, *models.CustomError) {
	thread, err := uc.forumRepo.UpdateThread(slugOrId, thread)
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

	post, err := uc.forumRepo.GetPostById(idNum)
	if err != nil {
		return models.PostInfo{}, &models.CustomError{Message: err.Error()}
	}
	result.Post = post

	if filter.User {
		user, err := uc.userRepo.GetUser(post.Author)
		if err != nil {
			return models.PostInfo{}, &models.CustomError{Message: err.Error()}
		}
		result.Author = &user
	}

	if filter.Thread {
		thread, err := uc.forumRepo.GetThreadById(int(post.Thread))
		if err != nil {
			return models.PostInfo{}, &models.CustomError{Message: err.Error()}
		}
		result.Thread = &thread
	}

	if filter.Forum {
		forum, err := uc.forumRepo.GetForumBySlug(post.Forum)
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
		post, err = uc.forumRepo.GetPostById(idNum)
	} else {
		post, err = uc.forumRepo.UpdatePost(idNum, post)
	}
	if err != nil {
		return models.Post{}, &models.CustomError{Message: err.Error()}
	}

	return post, nil
}

