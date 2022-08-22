package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/Chlor87/graphql/graph/generated"
	"github.com/Chlor87/graphql/middleware"
	"github.com/Chlor87/graphql/model"
	"github.com/Chlor87/graphql/repo"
	"github.com/Chlor87/graphql/resolvers"
)

const (
	defaultPort = "8080"
	dsn         = "user=postgres dbname=postgres password=password host=localhost port=5432 sslmode=disable"
)

var (
	db *gorm.DB
)

func init() {
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	check(err)
}

func check(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	todoRepo, err := repo.New[model.Todo](db)
	check(err)

	userRepo, err := repo.New[model.User](db)
	check(err)

	cfg := generated.Config{
		Resolvers: &resolvers.Resolver{
			DB:       db,
			TodoRepo: todoRepo,
			UserRepo: userRepo,
		},
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(cfg))

	mw := middleware.Build(
		middleware.AddUser(userRepo),
		middleware.AddUserLoader(userRepo),
	)

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", mw(srv))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
