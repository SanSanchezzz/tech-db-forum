package main

import (
	"database/sql"
	"fmt"
	config "github.com/SanSanchezzz/tech-db-forum/congif"
	delivery2 "github.com/SanSanchezzz/tech-db-forum/internal/forum/delivery"
	forumRep "github.com/SanSanchezzz/tech-db-forum/internal/forum/repository"
	usecase2 "github.com/SanSanchezzz/tech-db-forum/internal/forum/usecase"
	delivery4 "github.com/SanSanchezzz/tech-db-forum/internal/post/delivery"
	"github.com/SanSanchezzz/tech-db-forum/internal/post/repository"
	usecase4 "github.com/SanSanchezzz/tech-db-forum/internal/post/usecase"
	delivery5 "github.com/SanSanchezzz/tech-db-forum/internal/service/delivery"
	repository2 "github.com/SanSanchezzz/tech-db-forum/internal/service/repository"
	usecase5 "github.com/SanSanchezzz/tech-db-forum/internal/service/usecase"
	delivery3 "github.com/SanSanchezzz/tech-db-forum/internal/thread/delivery"
	threadRep "github.com/SanSanchezzz/tech-db-forum/internal/thread/repository"
	usecase3 "github.com/SanSanchezzz/tech-db-forum/internal/thread/usecase"
	"github.com/SanSanchezzz/tech-db-forum/internal/user/delivery"
	userRep "github.com/SanSanchezzz/tech-db-forum/internal/user/repository"
	"github.com/SanSanchezzz/tech-db-forum/internal/user/usecase"
	"log"
	"os"

	"github.com/sirupsen/logrus"

	_ "github.com/lib/pq"

	"github.com/labstack/echo/v4"
)

const (
	argsCount = 2
	argsUsage = "Usage: go run main.go $config_file"
	dbName    = "postgres"
)

func main() {
	if len(os.Args) != argsCount {
		fmt.Println(argsUsage)
		return
	}

	configFileName := os.Args[1]

	conf, err := config.LoadConfig(configFileName)
	if err != nil {
		logrus.Fatal(err)
	}

	dbConn, err := sql.Open(dbName, conf.GetDBConn())
	if err != nil {
		logrus.Fatal(err)
	}

	defer func() {
		_ = dbConn.Close()
	}()

	if err := dbConn.Ping(); err != nil {
		logrus.Fatal(err)
	}
	log.Printf("DB connected on %s", conf.GetDBConn())

	e := echo.New()

	userRepository := userRep.NewUserRepository(dbConn)
	forumRepository := forumRep.NewForumRepository(dbConn)
	threadRepository := threadRep.NewThreadRepository(dbConn)
	postRepository := repository.NewPostRepository(dbConn)
	serviceRepository := repository2.NewServiceRepository(dbConn)

	userUsecase := usecase.NewUserUsecase(userRepository)
	forumUsecase := usecase2.NewForumUsecse(forumRepository)
	threadUsecase := usecase3.NewThreadUsecase(threadRepository, userRepository)
	postUsecase := usecase4.NewPostUsecse(postRepository, userRepository, forumRepository, threadRepository)
	serviceUsecase := usecase5.NewServiceUsecase(serviceRepository)


	uh := delivery.CreateUserHandler(userUsecase)
	fh := delivery2.CreateForumHandler(userUsecase, forumUsecase, threadUsecase)
	th := delivery3.CreateThreadHandler(threadUsecase, userUsecase, forumUsecase, postUsecase)
	ph := delivery4.CreatePostHandler(userUsecase, forumUsecase, threadUsecase, postUsecase)
	sh := delivery5.CreateServiceHandler(serviceUsecase)


	uh.Configure(e)
	fh.Configure(e)
	th.Configure(e)
	ph.Configure(e)
	sh.Configure(e)

	e.Logger.Fatal(e.Start(conf.GetServerConn()))
}
