package tools

import (
	"fmt"
	"github.com/jackc/pgx"
)

func GetPostgres() (*pgx.ConnPool, error) {
	dsn := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s sslmode=disable",
		"bush", "forum",
		"docker", "localhost",
		"5432")
	db, err := pgx.ParseConnectionString(dsn)
	if err != nil {
		return nil, err
	}
	pool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:     db,
		MaxConnections: 100,
		AfterConnect:   nil,
		AcquireTimeout: 0,
	})
	if err != nil {
		panic(err)
	}

	return pool, nil
}
