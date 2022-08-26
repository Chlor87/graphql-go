package main

import (
	"bytes"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/Chlor87/graphql/model"
	"github.com/Chlor87/graphql/repo"
	"github.com/Chlor87/graphql/resolvers"
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
		Logger: logger.Default.LogMode(logger.Info),
	})
	check(err)
	s, err := readSchema("graph")
	check(err)

	todoRepo, err := repo.New[model.Todo](db)
	check(err)

	schema = graphql.MustParseSchema(s, &resolvers.Root{DB: db, TodoRepo: todoRepo})
	check(err)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	http.Handle("/query", &relay.Handler{Schema: schema})

	// userRepo, err := repo.New[model.User](db)
	// check(err)
	log.Printf("starting http service on :%s\n", port)
	log.Fatalln(http.ListenAndServe(":"+port, nil))

}

func check(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func readSchema(dir string) (res string, err error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}

	var (
		mu  sync.Mutex
		wg  sync.WaitGroup
		tmp bytes.Buffer
	)

	for _, f := range files {
		if f.IsDir() {
			continue
		}
		wg.Add(1)
		go func(f fs.FileInfo) {
			defer wg.Done()
			s, err := ioutil.ReadFile(filepath.Join(dir, f.Name()))
			// TODO: return err to caller
			check(err)
			mu.Lock()
			_, err = tmp.Write(s)
			mu.Unlock()
			check(err)
		}(f)
	}

	wg.Wait()

	return tmp.String(), nil
}
