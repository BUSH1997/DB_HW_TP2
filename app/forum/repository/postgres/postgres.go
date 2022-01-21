package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/BUSH1997/DB_HW_TP2/app/models"
	"github.com/BUSH1997/DB_HW_TP2/app/tools"
	"github.com/jackc/pgx"
	"strconv"
	"strings"
)

type StorageForumDB struct {
	db *pgx.ConnPool
}

const Asc = "asc"
const Desc = "desc"
const Default = ""

func NewStorageForumDB(db *pgx.ConnPool, err error) (*StorageForumDB, error) {
	if err != nil {
		return nil, err
	}
	return &StorageForumDB{
		db: db,
	}, nil
}

func (r *StorageForumDB) CreateForum(forum models.Forum) (models.Forum, error) {
	var nickName string
	err := r.db.QueryRow(`select nickname from users where nickname = $1`, forum.User).Scan(&nickName)
	if err != nil && nickName == "" {
		return models.Forum{}, nil
	}

	if err != nil  {
		return models.Forum{}, err
	}

	if nickName != forum.User {
		forum.User = nickName
	}

	_, err = r.db.Exec(`insert into forum(title, "user", slug) values ($1, $2, $3)`, forum.Title, forum.User, forum.Slug)

	if err != nil {
		return models.Forum{}, err
	}

	return forum, nil
}

func (r *StorageForumDB) GetForum(slug string) (models.Forum, error) {
	var result models.Forum
	err := r.db.QueryRow(`select title, "user", slug, posts, threads from forum where slug=$1`, slug).
		Scan(&result.Title, &result.User, &result.Slug, &result.Posts, &result.Threads)

	if err != nil {
		return models.Forum{}, err
	}
	return result, nil
}

func (r *StorageForumDB) CreateThread(thread models.Thread) (models.Thread, error) {
	var slug sql.NullString

	err := r.db.QueryRow(`insert into thread (title, author, forum, message, slug, created)
		values ($1, $2, coalesce((select slug from forum 
        where slug = $3), $3), $4, coalesce(nullif($5,'')), $6) 
		returning id, title, author, forum, message, slug, created`,
		thread.Title, thread.Author, thread.Forum, thread.Message, thread.Slug, thread.Created).
		Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &slug, &thread.Created)

	if err != nil {
		return models.Thread{}, err
	}
	thread.Slug = slug.String
	return thread, nil
}

func (r *StorageForumDB) GetForumUsers(slug string, filter tools.FilterUser) ([]models.User, error) {
	query := strings.Builder{}
	query.WriteString("select u.nickname, fullname, about, email from users_forum as u join users on u.nickname = users.nickname ")

	var rows *pgx.Rows
	var err error
	if filter.Since == Default {
		query.WriteString("where u.slug = $1 order by u.nickname collate \"C\" " + filter.Desc + " limit $2")
		rows, err = r.db.Query(query.String(), slug, filter.Limit)
	} else {
		if filter.Desc == Desc {
			query.WriteString("where u.slug = $1 and u.nickname < ($2 collate \"C\") order by u.nickname collate \"C\" desc limit $3")
			rows, err = r.db.Query(query.String(), slug, filter.Since,filter.Limit)
		} else {
			query.WriteString("where u.slug = $1 and u.nickname > ($2 collate \"C\") order by u.nickname collate \"C\" asc limit $3")
			rows, err = r.db.Query(query.String(), slug, filter.Since,filter.Limit)
		}
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err = rows.Scan(&user.Nickname, &user.FullName, &user.About, &user.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return users, nil
}

func (r *StorageForumDB) GetForumThreads(slug string, filter tools.FilterThread) ([]models.Thread, error) {
	var rows *pgx.Rows
	var err error

	if filter.Sort != Asc && filter.Sort != Desc {
		return nil, errors.New("sql attack")
	}

	if filter.Since == Default {
		rows, err = r.db.Query(`select id, title, author, forum, message, votes, slug, created from thread where forum = $1 order by created `+
			filter.Sort +` limit $2`, slug, filter.Limit)
	} else {
		if filter.Sort == tools.SortParamTrue {
			rows, err = r.db.Query(`select id, title, author, forum, message, votes, slug, created from thread where forum = $1 and created <= $3 order by created `+
			filter.Sort +` limit $2`, slug, filter.Limit, filter.Since)
		} else {
			rows, err = r.db.Query(`select id, title, author, forum, message, votes, slug, created from thread where forum = $1 and created >= $3 order by created `+
			filter.Sort +` limit $2`, slug, filter.Limit, filter.Since)
		}
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var nullSlug sql.NullString
	var threads []models.Thread
	for rows.Next() {
		var thread models.Thread
		err = rows.Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message,
			&thread.Votes, &nullSlug, &thread.Created)
		if err != nil {
			return nil, err
		}
		thread.Slug = nullSlug.String
		threads = append(threads, thread)
	}

	return threads, nil
}

func (r *StorageForumDB) GetForumBySlug(slug string) (models.Forum, error) {
	var forumData models.Forum
	err := r.db.QueryRow(`select slug, title, "user", posts, threads from forum where slug=$1`, slug).
		Scan(&forumData.Slug, &forumData.Title, &forumData.User, &forumData.Posts, &forumData.Threads)
	if err != nil {
		return models.Forum{}, err
	}

	return forumData, nil
}

func (r *StorageForumDB) CreatePosts(threadId int, threadForum string, posts []models.Post) ([]models.Post, error) {
	query := `insert into post(parent, author, message, thread, forum) values `
	var values []interface{}
	if len(posts) == 0 {
		query += fmt.Sprintf(`(0, null, null, %d, '%s')`, threadId, threadForum)
	}
	for i, post := range posts {
		value := fmt.Sprintf("($%d, $%d, $%d, $%d, $%d),",
			i*5+1, i*5+2, i*5+3, i*5+4, i*5+5)
		query += value
		values = append(values, post.Parent, post.Author, post.Message, threadId, threadForum)
	}
	query = strings.TrimSuffix(query, ",")
	query += ` returning id, parent, author, message, isEdited, forum, thread, created;`

	rows, err := r.db.Query(query, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []models.Post{}
	if len(posts) != 0 {
		for rows.Next() {
			var post models.Post
			err := rows.Scan(&post.Id, &post.Parent, &post.Author, &post.Message,
				&post.IsEdited, &post.Forum, &post.Thread, &post.Created)
			if err != nil {
				return nil, err
			}
			result = append(result, post)
		}

		if rows.Err() != nil {
			return nil, rows.Err()
		}
	}
	return result, nil
}

func (r *StorageForumDB) GetThreadBySlug(slug string) (models.Thread, error) {
	var result models.Thread
	err := r.db.QueryRow(`select id, title, author, forum, message, votes, slug, created from thread where slug=$1`, slug).
		Scan(&result.Id, &result.Title, &result.Author, &result.Forum, &result.Message, &result.Votes,
		&result.Slug, &result.Created)
	if err != nil {
		return models.Thread{}, err
	}

	return result, nil
}

func (r *StorageForumDB) GetThreadById(id int) (models.Thread, error) {
	var result models.Thread
	row := r.db.QueryRow(`select id, title, author, forum, message, votes, slug, created from thread where id=$1`, id)
	var slug sql.NullString
	err := row.Scan(&result.Id, &result.Title, &result.Author, &result.Forum, &result.Message, &result.Votes,
		&slug, &result.Created)
	if err != nil {
		return models.Thread{}, err
	}
	result.Slug = slug.String
	return result, nil
}

func (r *StorageForumDB) CreateVoteBySlugOrId(slugOrId string, vote models.Vote) error {
	query := strings.Builder{}
	query.WriteString("insert into vote(nickname, voice, thread) values ($1, $2, ")
	id, err := strconv.Atoi(slugOrId)
	if err != nil {
		query.WriteString("(select id from thread where slug = $3));")
		_, err = r.db.Exec(query.String(), vote.NickName, vote.Voice, slugOrId)
	} else {
		query.WriteString("$3);")
		_, err = r.db.Exec(query.String(), vote.NickName, vote.Voice, id)
	}

	if err != nil {
		return err
	}

	return nil
}

func (r *StorageForumDB) UpdateVoteBySlugOrId(slugOrId string, vote models.Vote) error {
	query := strings.Builder{}
	query.WriteString("update vote set voice=$1 where nickname=$2 and ")

	id, err := strconv.Atoi(slugOrId)
	if err != nil {
		query.WriteString("thread=(select id from thread where slug = $3)")
		_, err = r.db.Exec(query.String(), vote.Voice, vote.NickName, slugOrId)
	} else {
		query.WriteString("thread=$3")
		_, err = r.db.Exec(query.String(), vote.Voice, vote.NickName, id)
	}

	if err != nil {
		return err
	}

	return nil
}

func (r *StorageForumDB) GetThreadBySlugOrId(slugOrId string) (models.Thread, error) {
	query := strings.Builder{}
	query.WriteString("select id, title, author, forum, message, votes, slug, created from thread where ")

	var result models.Thread
	var row *pgx.Row
	id, err := strconv.Atoi(slugOrId)
	if err != nil {
		query.WriteString("slug=$1")
		row = r.db.QueryRow(query.String(), slugOrId)
	} else {
		query.WriteString("id=$1")
		row = r.db.QueryRow(query.String(), id)
	}

	var slug sql.NullString
	err = row.Scan(&result.Id, &result.Title, &result.Author, &result.Forum, &result.Message, &result.Votes,
		&slug, &result.Created)
	if err != nil {
		return models.Thread{}, err
	}
	result.Slug = slug.String
	return result, nil
}

func (r *StorageForumDB) GetPostsFlatSlugOrId(slugOrId string, filter tools.FilterPosts) ([]*models.Post, error) {
	query := strings.Builder{}
	query.WriteString("select id, parent, author, message, isEdited, forum, thread, created from post where ")

	var rows *pgx.Rows
	var err error

	id, err := strconv.Atoi(slugOrId)
	if err != nil {
		var tmpId sql.NullInt64
		row := r.db.QueryRow(`select id from thread where slug = $1`, slugOrId)
		err = row.Scan(&tmpId)

		if filter.Since == Default {
			query.WriteString("thread = $1 order by id " + filter.Desc + " limit $2")
			rows, err = r.db.Query(query.String(), tmpId, filter.Limit)
		} else {
			if filter.Desc == Desc {
				query.WriteString("thread = $1 and id < $2 order by id desc limit $3")
				rows, err = r.db.Query(query.String(), tmpId, filter.Since, filter.Limit)
			} else {
				query.WriteString("thread = $1 and id > $2 order by id asc limit $3")
				rows, err = r.db.Query(query.String(), tmpId, filter.Since, filter.Limit)
			}
		}
	} else {
		if filter.Since == Default {
			query.WriteString("thread = $1 order by id " + filter.Desc + " limit $2")
			rows, err = r.db.Query(query.String(), id, filter.Limit)
		} else {
			if filter.Desc == Desc {
				query.WriteString("thread = $1 and id < $2 order by id desc limit $3")
				rows, err = r.db.Query(query.String(), id, filter.Since, filter.Limit)
			} else {
				query.WriteString("thread = $1 and id > $2 order by id asc limit $3")
				rows, err = r.db.Query(query.String(), id, filter.Since, filter.Limit)
			}
		}
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		post := &models.Post{}

		err = rows.Scan(&post.Id, &post.Parent, &post.Author, &post.Message,
			            &post.IsEdited, &post.Forum, &post.Thread, &post.Created)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, err
}

func (r *StorageForumDB) GetPostsTreeSlugOrId(slugOrId string, filter tools.FilterPosts) ([]*models.Post, error) {
	query := strings.Builder{}
	query.WriteString("select id, parent, author, message, isEdited, forum, thread, created from post where ")

	var rows *pgx.Rows
	var err error

	id, err := strconv.Atoi(slugOrId)
	if err != nil {
		var tmpId sql.NullInt64
		row := r.db.QueryRow("select id from thread where slug = $1", slugOrId)
		err = row.Scan(&tmpId)
		if filter.Since == Default {
			query.WriteString("thread = $1 order by paths " + filter.Desc + " , id " + filter.Desc + " limit $2")
			rows, err = r.db.Query(query.String(), tmpId, filter.Limit)
		} else {
			if filter.Desc == Desc {
				query.WriteString("thread = $1 and paths < (select paths from post where id=$2) order by paths desc, id desc limit $3")
				rows, err = r.db.Query(query.String(), tmpId, filter.Since, filter.Limit)
			} else {
				query.WriteString("thread = $1 and paths > (select paths from post where id=$2) order by paths asc, id asc limit $3")
				rows, err = r.db.Query(query.String(), tmpId, filter.Since, filter.Limit)
			}
		}
	} else {
		if filter.Since == Default {
			query.WriteString("thread = $1 order by paths " + filter.Desc + " , id " + filter.Desc + " limit $2")
			rows, err = r.db.Query(query.String(), id, filter.Limit)
		} else {
			if filter.Desc == Desc {
				query.WriteString("thread = $1 and paths < (select paths from post where id=$2) order by paths desc, id desc limit $3")
				rows, err = r.db.Query(query.String(), id, filter.Since, filter.Limit)
			} else {
				query.WriteString("thread = $1 and paths > (select paths from post where id=$2) order by paths asc, id asc limit $3")
				rows, err = r.db.Query(query.String(), id, filter.Since, filter.Limit)
			}
		}
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*models.Post
	for rows.Next() {
		post := &models.Post{}

		err = rows.Scan(
			&post.Id,
			&post.Parent,
			&post.Author,
			&post.Message,
			&post.IsEdited,
			&post.Forum,
			&post.Thread,
			&post.Created)
		if err != nil {
			return nil, err
		}

		result = append(result, post)
	}

	return result, err
}

func (r *StorageForumDB) GetPostsParentTreeSlugOrId(slugOrId string, filter tools.FilterPosts) ([]*models.Post, error) {
	query := strings.Builder{}
	query.WriteString("select id, parent, author, message, isEdited, forum, thread, created from post where ")

	var rows *pgx.Rows
	var err error

	id, err := strconv.Atoi(slugOrId)
	if err != nil {
		var tmpId sql.NullInt64
		row := r.db.QueryRow("select id from thread where slug = $1", slugOrId)
		err = row.Scan(&tmpId)
		if filter.Since == Default {
			if filter.Desc == Desc {
				query.WriteString("paths[1] in (select id from post where thread = $1 and parent = 0 order by id desc limit $2) order by paths[1] desc, paths asc, id asc;")
				rows, err = r.db.Query(query.String(), tmpId, filter.Limit)
			} else {
				query.WriteString("paths[1] in (select id from post where thread = $1 and parent = 0 order by id asc limit $2) order by paths asc, id asc;")
				rows, err = r.db.Query(query.String(), tmpId, filter.Limit)
			}
		} else {
			if filter.Desc == Desc {
				query.WriteString("paths[1] in (select id from post where thread = $1 and parent = 0 and paths[1] < (select paths[1] from post where id = $2) order by id desc limit $3) order by paths[1] desc, paths asc, id asc;")
				rows, err = r.db.Query(query.String(), tmpId, filter.Since, filter.Limit)
			} else {
				query.WriteString("paths[1] in (select id from post where thread = $1 and parent = 0 and paths[1] > (select paths[1] from post where id = $2) order by id asc limit $3) order by paths asc, id asc;")
				rows, err = r.db.Query(query.String(), tmpId, filter.Since, filter.Limit)
			}
		}
	} else {
		if filter.Since == Default {
			if filter.Desc == Desc {
				query.WriteString("paths[1] in (select id from post where thread = $1 and parent = 0 order by id desc limit $2) order by paths[1] desc, paths asc, id asc;")
				rows, err = r.db.Query(query.String(), id, filter.Limit)
			} else {
				query.WriteString("paths[1] in (select id from post where thread = $1 and parent = 0 order by id asc limit $2) order by paths asc, id asc;")
				rows, err = r.db.Query(query.String(), id, filter.Limit)
			}
		} else {
			if filter.Desc == Desc {
				query.WriteString("paths[1] in (select id from post where thread = $1 and parent = 0 and paths[1] < (select paths[1] from post where id = $2) order by id desc limit $3) order by paths[1] desc, paths asc, id asc;")
				rows, err = r.db.Query(query.String(), id, filter.Since, filter.Limit)
			} else {
				query.WriteString("paths[1] in (select id from post where thread = $1 and parent = 0 and paths[1] > (select paths[1] from post where id = $2) order by id asc limit $3) order by paths asc, id asc;")
				rows, err = r.db.Query(query.String(), id, filter.Since, filter.Limit)
			}
		}
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*models.Post
	for rows.Next() {
		post := &models.Post{}

		err = rows.Scan(&post.Id, &post.Parent, &post.Author, &post.Message,
			            &post.IsEdited, &post.Forum, &post.Thread, &post.Created)
		if err != nil {
			return nil, err
		}

		result = append(result, post)
	}

	return result, err
}

func (r *StorageForumDB) UpdateThread(slugOrId string, thread models.Thread) (models.Thread, error) {
	query := strings.Builder{}
	query.WriteString("update thread set title=(case when $1='' then title else $1 end), " +
		                                   "author=(case when $2='' then author else $2 end), " +
		                                   "forum=(case when $3='' then forum else $3 end), " +
		                                   "message=(case when $4='' then message else $4 end) ")

	var row *pgx.Row
	var err error
	id, err := strconv.Atoi(slugOrId)
	if err != nil {
		query.WriteString("where slug=$5 returning id, title, author, forum, message, votes, slug, created")
		row = r.db.QueryRow(query.String(), thread.Title, thread.Author, thread.Forum, thread.Message, slugOrId)
	} else {
		query.WriteString("where id=$5 returning id, title, author, forum, message, votes, slug, created")
		row = r.db.QueryRow(query.String(), thread.Title, thread.Author, thread.Forum, thread.Message, id)
	}

	err = row.Scan(
		&thread.Id,
		&thread.Title,
		&thread.Author,
		&thread.Forum,
		&thread.Message,
		&thread.Votes,
		&thread.Slug,
		&thread.Created)
	if err != nil {
		return models.Thread{}, err
	}

	return thread, nil
}

func (r *StorageForumDB) GetPostById(id int) (models.Post, error) {
	var result models.Post
	row := r.db.QueryRow(`select id, parent, author, message, isEdited, forum, thread, created from post where id=$1`, id)

	err := row.Scan(&result.Id, &result.Parent, &result.Author, &result.Message, &result.IsEdited,
		&result.Forum, &result.Thread, &result.Created)
	if err != nil {
		return models.Post{}, err
	}
	return result, nil
}

func (r *StorageForumDB) UpdatePost(id int, post models.Post) (models.Post, error) {
	err := r.db.QueryRow(`update post set message=$1, isedited = case when message = $1 then isedited else true end 
		where id=$2 returning id, parent, author, message, isedited, forum, thread, created`,
		post.Message, id).
		Scan(&post.Id, &post.Parent, &post.Author, &post.Message,
		&post.IsEdited, &post.Forum, &post.Thread, &post.Created,
		)
	if err != nil {
		return models.Post{}, err
	}

	return post, nil
}
