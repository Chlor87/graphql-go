package main

import (
	"log"
	"net/http"
	"os"

	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/Chlor87/graphql/model"
	"github.com/Chlor87/graphql/repo"
	"github.com/Chlor87/graphql/resolvers"
	"github.com/Chlor87/graphql/util"
)

const (
	defaultPort = "8080"
	dsn         = "user=postgres dbname=postgres password=password host=localhost port=5432 sslmode=disable"
)

var (
	db     *gorm.DB
	schema *graphql.Schema
)

func init() {
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		// Logger: logger.Default.LogMode(logger.Info),
	})
	check(err)
	s, err := util.LoadSchema("./graph")
	check(err)

	todoRepo, err := repo.New[model.Todo](db)
	check(err)

	userRepo, err := repo.New[model.User](db)
	check(err)

	schema = graphql.MustParseSchema(
		s,
		resolvers.NewRoot(todoRepo, userRepo),
		graphql.UseFieldResolvers(),
	)
	check(err)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	http.Handle("/query", &relay.Handler{Schema: schema})

	log.Printf("starting http service on :%s\n", port)
	log.Fatalln(http.ListenAndServe(":"+port, nil))

}

func check(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
