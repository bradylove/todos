package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"strconv"
	"time"
)

var db gorm.DB

type Todo struct {
	ID          int       `json:"id"`
	Status      string    `json:"status"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func main() {
	var err error
	db, err = gorm.Open("sqlite3", "todos.db")
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&Todo{})

	r := mux.NewRouter()
	r.HandleFunc("/todos", getTodos).Methods("GET")
	r.HandleFunc("/todos/{id}", getTodo).Methods("GET")

	http.Handle("/", r)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func getTodos(w http.ResponseWriter, r *http.Request) {
	var todos []Todo
	db.Find(&todos)

	w.Header()["Content-Type"] = []string{"application/json"}
	err := json.NewEncoder(w).Encode(&todos)
	if err != nil {
		// TODO: Render 500 error
		panic(err)
	}
}

func getTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		// TODO: Return 404
		panic(err)
	}

	var todo Todo
	db.First(&todo, id)
	// TODO check if one was found, if not, render 404

	w.Header()["Content-Type"] = []string{"application/json"}
	err = json.NewEncoder(w).Encode(&todo)
	if err != nil {
		// TODO: Render 500 error
		panic(err)
	}
}
