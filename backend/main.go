package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	PORT  = ":1234"
	QUERY = "football"
)

func main() {
	if err := initDb(); err != nil {
		log.Fatal("[DB ERROR] ", err)
	}
	go fetchLatestVideos(QUERY)

	log.Println("Server running on port ", PORT)
	log.Fatal(http.ListenAndServe(PORT, newRouter()))
}

func newRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/ping", pingHandler).Methods(http.MethodGet)
	r.HandleFunc("/", listHandler).Methods(http.MethodGet)
	r.HandleFunc("/search", searchHandler).Methods(http.MethodGet)
	return r
}
