package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

var db gorm.DB

func main() {
	port := ":8080"
	var err error
	db, err = gorm.Open("sqlite3", "todos.db")
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&Todo{})

	r := mux.NewRouter()
	r.HandleFunc("/todos", getTodos).Methods("GET")
	r.HandleFunc("/todos", createTodo).Methods("POST")
	r.HandleFunc("/todos/{id}", getTodo).Methods("GET")
	r.HandleFunc("/todos/{id}", deleteTodo).Methods("DELETE")
	r.HandleFunc("/todos/{id}", updateTodo).Methods("PUT")

	http.Handle("/", r)

	log.Println("Starting server on port", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		panic(err)
	}
}

// Render 404
func renderNotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
}

// Render a given status code
func renderStatus(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
}

// Render an object to JSON with a given status code
func renderJsonWithStatus(w http.ResponseWriter, status int, objects interface{}) {
	w.Header()["Content-Type"] = []string{"application/json"}
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(objects)
}
