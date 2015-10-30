package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
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
}
