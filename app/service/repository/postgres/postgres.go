package postgres

import (
	"github.com/BUSH1997/DB_HW_TP2/app/models"
	"github.com/jackc/pgx"
)

type StorageServiceDB struct {
	db *pgx.ConnPool
}

func NewStorageServiceDB(db *pgx.ConnPool, err error) (*StorageServiceDB, error) {
	if err != nil {
		return nil, err
	}
	return &StorageServiceDB{
		db: db,
	}, nil
}

func (r *StorageServiceDB) GetStatus() (models.Status, error) {
	var result models.Status
	row := r.db.QueryRow(
		`select * from
		(select count(*) from users) as u,
 		(select count(*) from forum) as f,
		(select count(*) from thread) as t,
		(select count(*) from post) as p;`)

	err := row.Scan(
		&result.User,
		&result.Forum,
		&result.Thread,
		&result.Post)
	if err != nil {
		return models.Status{}, err
	}

	return result, nil
}
func (r *StorageServiceDB) Clear() error {
	_, err := r.db.Exec(`truncate users, forum, thread, post, vote, users_forum;`)
	if err != nil {
		return err
	}

	return nil
}
