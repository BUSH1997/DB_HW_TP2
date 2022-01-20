package main

import (
	"fmt"
	forumHandler "github.com/BUSH1997/DB_HW_TP2/app/forum/delivery/http"
	forumRepository "github.com/BUSH1997/DB_HW_TP2/app/forum/repository/postgres"
	forumUC "github.com/BUSH1997/DB_HW_TP2/app/forum/usecase"
	serviceHandler "github.com/BUSH1997/DB_HW_TP2/app/service/delivery/http"
	serviceRepository "github.com/BUSH1997/DB_HW_TP2/app/service/repository/postgres"
	serviceUC "github.com/BUSH1997/DB_HW_TP2/app/service/usecase"
	userHandler "github.com/BUSH1997/DB_HW_TP2/app/user/delivery/http"
	userRepository "github.com/BUSH1997/DB_HW_TP2/app/user/repository/postgres"
	userUC "github.com/BUSH1997/DB_HW_TP2/app/user/usecase"
	"github.com/BUSH1997/DB_HW_TP2/config/configRouting"
	"github.com/jackc/pgx"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/labstack/echo/v4"
	"log"
)

var (
	router = echo.New()
)

func main() {
	dbConn, dbErr := GetPostgres()
	if dbErr != nil {
		panic(dbErr)
	}

	storageUser, err := userRepository.NewStorageUserDB(dbConn, dbErr)
	if err != nil {
		panic(err)
	}

	userHandler := userHandler.NewUserHandler(userUC.NewUseCase(storageUser))

	storageForum, err := forumRepository.NewStorageForumDB(dbConn, dbErr)
	if err != nil {
		panic(err)
	}

	forumHandler := forumHandler.NewForumHandler(forumUC.NewUseCase(storageForum, storageUser))

	storageService, err := serviceRepository.NewStorageServiceDB(dbConn, dbErr)
	if err != nil {
		panic(err)
	}

	serviceHandler := serviceHandler.NewServiceHandler(serviceUC.NewUseCase(storageService))

	serverRouting := configRouting.ServerConfigRouting{
		ForumHandler:   *forumHandler,
		UserHandler: *userHandler,
		ServiceHandler: *serviceHandler,
	}
	serverRouting.ConfigRouting(router)
	if err := router.Start(":5000"); err != nil {
		log.Fatal(err)
	}
}

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
