package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/weather/{location}", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Weather endpoint reached"))
	}).Methods(http.MethodGet)

	log.Println("Server is starting at port :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
