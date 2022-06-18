package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

func get_sqlite_conn() *sql.DB {
	db, err := sql.Open("sqlite3", "sqlite.db")

	if err != nil {
		log.Fatal(err)
	}

	return db
}

func main() {
	listen_port := flag.String("port", ":8080", "Listening port")

	flag.Parse()

	db := get_sqlite_conn()
	repository := NewSQLiteRepository(db)

	if err := repository.Migrate(); err != nil {
		log.Fatal(err)
	}

	handler := func(handler func(w http.ResponseWriter, r *http.Request, repository *SQLiteRepository)) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			handler(w, r, repository)
		}
	}

	r := mux.NewRouter()
	r.HandleFunc("/", handler(getCounterHandler)).Methods("GET")
	r.HandleFunc("/", handler(increaseCounterHandler)).Methods("POST")

	serveMux := http.NewServeMux()
	serveMux.Handle("/", r)
	server := &http.Server{
		Addr:    *listen_port,
		Handler: serveMux,
	}

	err := server.ListenAndServe()

	if err != nil {
		log.Println(err)
	}
}
