package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

var db gorm.DB

func main() {
	var err error
	db, err = gorm.Open("sqlite3", "todos.db")
	if err != nil {
		panic(err)
	}
}
