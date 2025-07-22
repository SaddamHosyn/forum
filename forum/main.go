package main

import (
	"forum-go/database"
	"forum-go/server"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	 if database.DB == nil {
		log.Fatal("Database connection is nil after initialization")
	 }
	defer database.DB.Close()

	server.Startserver(database.DB)

}
