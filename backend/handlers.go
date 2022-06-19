package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var ErrInvalidUrl = errors.New("invalid url")

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

func getUrlFromRequest(w http.ResponseWriter, r *http.Request, allowUrlPrefix string) (string, error) {
	url := r.URL.Query().Get("url")

	if url == "" {
		http.Error(w, "Missing url param", 400)
		return "", ErrInvalidUrl
	}

	if allowUrlPrefix != "" && !strings.HasPrefix(url, allowUrlPrefix) {
		http.Error(w, "Invalid url param", 400)
		return "", ErrInvalidUrl
	}

	return url, nil
}

func getCounterHandler(w http.ResponseWriter, r *http.Request, repository *SQLiteRepository, allowUrlPrefix string) {
	url, err := getUrlFromRequest(w, r, allowUrlPrefix)

	if err != nil {
		return
	}

	counter := getCounter(url, repository)
	w.Write([]byte(counter))
}

func increaseCounterHandler(w http.ResponseWriter, r *http.Request, repository *SQLiteRepository, allowUrlPrefix string) {
	url, err := getUrlFromRequest(w, r, allowUrlPrefix)

	if err != nil {
		return
	}

	counter := increaseCounter(url, repository)
	w.Write([]byte(counter))
}
