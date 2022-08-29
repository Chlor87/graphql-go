package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/Chlor87/graphql/domain"
	mw "github.com/Chlor87/graphql/middleware"
	"github.com/Chlor87/graphql/resolvers"
	"github.com/Chlor87/graphql/util"
)

const (
	defaultPort = "8080"
	dsn         = "user=postgres dbname=postgres password=password host=localhost port=5432 sslmode=disable"
)

var (
	schema *graphql.Schema
	d      *domain.Domain
)

func init() {
	var err error
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		// uncomment to enable ORM logger
		// Logger: logger.Default.LogMode(logger.Info),
	})
	check(err)
	s, err := util.LoadSchema("./graphql")

	time.Sleep(time.Second)
	check(err)

	d, err = domain.New(db)
	check(err)

	schema = graphql.MustParseSchema(
		s,
		resolvers.NewRoot(d),
		graphql.UseFieldResolvers(),
	)
	check(err)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	mws := mw.Build(
		mw.Log,
		mw.AddUser(d.User),
		mw.AddUserLoader(d.User),
	)

	// serve graphiql at root
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./public/index.html")
	})

	// graphql entrypoint
	http.Handle("/query", mws(&relay.Handler{Schema: schema}))

	log.Printf("starting http service on :%s\n", port)
	log.Fatalln(http.ListenAndServe("localhost:"+port, nil))

}

func check(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
