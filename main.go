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
	r.HandleFunc("/todos", createTodo).Methods("POST")
	r.HandleFunc("/todos/{id}", getTodo).Methods("GET")
	r.HandleFunc("/todos/{id}", deleteTodo).Methods("DELETE")
	r.HandleFunc("/todos/{id}", updateTodo).Methods("PUT")

	http.Handle("/", r)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func renderNotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
}

func renderStatus(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
}

func renderJsonWithStatus(w http.ResponseWriter, status int, objects interface{}) {
	w.Header()["Content-Type"] = []string{"application/json"}
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(objects)
}

func getTodoId(r *http.Request) (int, error) {
	vars := mux.Vars(r)
	return strconv.Atoi(vars["id"])
}

// Renders a JSON array of all Todos in the database.
func getTodos(w http.ResponseWriter, r *http.Request) {
	var todos []Todo
	db.Find(&todos)

	renderJsonWithStatus(w, http.StatusOK, &todos)
}

// Renders a JSON object of a single Todo with a given ID.
func getTodo(w http.ResponseWriter, r *http.Request) {
	id, err := getTodoId(r)
	if err != nil {
		renderNotFound(w)
		return
	}

	var todo Todo
	if err := db.First(&todo, id).Error; err != nil {
		renderNotFound(w)
		return
	}

	renderJsonWithStatus(w, http.StatusOK, &todo)
}

// Creates and saves a new Todo and returns the new Todo as a JSON object.
func createTodo(w http.ResponseWriter, r *http.Request) {
	var m map[string]string
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		renderStatus(w, http.StatusBadRequest)
		return
	}

	todo := Todo{Description: m["description"], Status: "pending"}
	if err := db.Create(&todo).Error; err != nil {
		renderStatus(w, http.StatusInternalServerError)
		return
	}

	renderJsonWithStatus(w, http.StatusCreated, &todo)
}

// Updates and saves an existing Todo and returns the updated Todo as a JSON
// object.
func updateTodo(w http.ResponseWriter, r *http.Request) {
	id, err := getTodoId(r)
	if err != nil {
		renderNotFound(w)
		return
	}

	var m map[string]string
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		renderStatus(w, http.StatusBadRequest)
		return
	}

	var todo Todo
	if err := db.First(&todo, id).Error; err != nil {
		renderNotFound(w)
		return
	}

	todo.Description = m["description"]
	todo.Status = m["status"]

	if err := db.Save(&todo).Error; err != nil {
		renderStatus(w, http.StatusInternalServerError)
		return
	}

	renderJsonWithStatus(w, http.StatusOK, &todo)
}

// Deletes a Todo and returns the deleted Todo as a JSON object.
func deleteTodo(w http.ResponseWriter, r *http.Request) {
	id, err := getTodoId(r)
	if err != nil {
		renderNotFound(w)
		return
	}

	var todo Todo
	db.First(&todo, id)
	if err := db.First(&todo, id).Error; err != nil {
		renderNotFound(w)
		return
	}

	if err := db.Delete(&todo).Error; err != nil {
		renderStatus(w, http.StatusInternalServerError)
	} else {
		renderJsonWithStatus(w, http.StatusOK, &todo)
	}
}
