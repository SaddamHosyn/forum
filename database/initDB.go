package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

var DB *sql.DB

func InitDB() error {

	var err error

	if DB == nil {
		dbPath := os.Getenv("DB_PATH")
		if dbPath == "" {
			dbPath = "reeltalk.db"
		}
		DB, err = sql.Open("sqlite3", dbPath)
		if err != nil {
			return fmt.Errorf("error opening database: %v", err)
		}
		log.Println("Database connection opened successfully")
	}

	// Test the connection
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("error pinging database: %v", err)
	}
	log.Println("Database pinged successfully")

	if err = createTables(); err != nil {
		return fmt.Errorf("error creating tables: %v", err)
	}

	if err = insertCategories(); err != nil {
		return fmt.Errorf("error inserting categories data: %v", err)
	}

	if err = insertUsers(); err != nil {
		return fmt.Errorf("error inserting users data: %v", err)
	}

	if err = insertPosts(); err != nil {
		return fmt.Errorf("error inserting users data: %v", err)
	}

	if err = verifyData(); err != nil {
		return fmt.Errorf("error verifying data: %v", err)
	}

	// Execute schema
	_, err = DB.Exec(`PRAGMA foreign_keys = ON`)
	if err != nil {
		return fmt.Errorf("error enabling foreign keys: %v", err)
	}

	log.Println("Database initialization completed")
	return nil
}

func verifyData() error {

	tables := []string{"categories", "users"}
	for _, table := range tables {
		var count int
		err := DB.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", table)).Scan(&count)
		if err != nil {
			return fmt.Errorf("error counting rows in %s: %v", table, err)
		}
		fmt.Printf("Table %s has %d rows\n", table, count)
	}

	return nil
}
