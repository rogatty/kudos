package main

import (
	"log"
	"net/http"
	"strconv"
)

func getCounter(url string, repository *SQLiteRepository) string {
	kudos, err := repository.GetByUrl(url)

	if err != nil {
		log.Println(err)
		return "0"
	}

	return strconv.FormatInt(kudos.Counter, 10)
}

func increaseCounter(url string, repository *SQLiteRepository) string {
	kudos, err := repository.IncreaseCounterByUrl(url)

	if err != nil {
		log.Println(err)
		return "0"
	}

	return strconv.FormatInt(kudos.Counter, 10)
}

func getCounterHandler(w http.ResponseWriter, r *http.Request, repository *SQLiteRepository) {
	var counter string
	url := r.URL.Query().Get("url")

	if url != "" {
		counter = getCounter(url, repository)
	} else {
		counter = "0"
	}

	w.Write([]byte(counter))
}

func increaseCounterHandler(w http.ResponseWriter, r *http.Request, repository *SQLiteRepository) {
	url := r.URL.Query().Get("url")

	if url == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), 400)
		return
	}

	counter := increaseCounter(url, repository)

	w.Write([]byte(counter))
}
