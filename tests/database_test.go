package tests

import (
	"database/sql"
	"forum-go/auth"
	"forum-go/database"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
	database.DB = nil
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open in-memory sqlite database: %v", err)
	}

	database.DB = db

	if err := database.InitDB(); err != nil {
		t.Fatalf("InitDB failed on memory DB: %v", err)
	}

	return db
}

func TestDatabaseCategoriesAndPosts(t *testing.T) {
	_ = setupTestDB(t)

	categories, err := database.FetchCategories()
	if err != nil {
		t.Fatalf("FetchCategories failed: %v", err)
	}
	if len(categories) == 0 {
		t.Error("Expected categories to be populated, got 0")
	}

	posts, err := database.FetchPosts()
	if err != nil {
		t.Fatalf("FetchPosts failed: %v", err)
	}
	if len(posts) == 0 {
		t.Error("Expected posts to be populated, got 0")
	}
}

func TestAuthAddUserAndGetInfo(t *testing.T) {
	db := setupTestDB(t)

	testUsername := "testuser99"
	testEmail := "testuser99@example.com"
	testPassword := "TestPass123!"

	err := auth.AddUser(db, testUsername, testEmail, testPassword)
	if err != nil {
		t.Fatalf("AddUser failed: %v", err)
	}

	user, err := auth.GetUserInfo(db, testUsername)
	if err != nil {
		t.Fatalf("GetUserInfo failed: %v", err)
	}

	if user.Username != testUsername {
		t.Errorf("GetUserInfo username mismatch: got %v, want %v", user.Username, testUsername)
	}
	if user.Email != testEmail {
		t.Errorf("GetUserInfo email mismatch: got %v, want %v", user.Email, testEmail)
	}
}
