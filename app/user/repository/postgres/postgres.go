package postgres

import (
	"github.com/BUSH1997/DB_HW_TP2/app/models"
	"github.com/jackc/pgx"
)

type StorageUserDB struct {
	db *pgx.ConnPool
}

func NewStorageUserDB(db *pgx.ConnPool, err error) (*StorageUserDB, error) {
	if err != nil {
		return nil, err
	}
	return &StorageUserDB{
		db: db,
	}, nil
}

func (r *StorageUserDB) AddUser(user models.User) (models.User, error) {
	_, err := r.db.Exec(`insert into users(nickname, fullname, about, email) values ($1, $2, $3, $4)`,
		user.Nickname, user.FullName, user.About, user.Email)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (r *StorageUserDB) GetUser(nickname string) (models.User, error) {
	var result models.User
	row := r.db.QueryRow(`select nickname, fullname, about, email 
		from users where nickname=$1`, nickname)

	err := row.Scan(&result.Nickname, &result.FullName, &result.About, &result.Email)
	if err != nil {
		return models.User{}, err
	}
	return result, nil
}

func (r *StorageUserDB) UpdateUser(user models.User) (models.User, error) {
	query := r.db.QueryRow(`update users set 
		fullname=coalesce(nullif($1, ''), fullname), 
		about=coalesce(nullif($2, ''), about),
		email=coalesce(nullif($3, ''), email) 
		where nickname=$4 returning nickname, fullname, about, email`, user.FullName, user.About, user.Email, user.Nickname)

	err := query.Scan(
		&user.Nickname,
		&user.FullName,
		&user.About,
		&user.Email)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (r *StorageUserDB) GetUsersByNicknameOrEmail(nickname string, email string) ([]models.User, error) {
	rows, err := r.db.Query(`select nickname, fullname, about, email 
		from users where nickname=$1 or email=$2`, nickname, email)
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
