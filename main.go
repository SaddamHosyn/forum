package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	"forum-go/database"

	_ "github.com/mattn/go-sqlite3"
)

// User struct matches JSON structure
type User struct {
	Username     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
}

// Movie struct for JSON data
type Movie struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ReleaseDate string `json:"release_date"`
	ImageURL    string `json:"image_url"`
}

type Genre struct {
	Name string `json:"name"`
}

// Comment struct for JSON data
type Comment struct {
	UserID  int    `json:"user_id"`
	MovieID int    `json:"movie_id"`
	Content string `json:"content"`
}

// MovieGenre struct for JSON data, this is not directly insertable.
type MovieGenre struct {
    MovieID  int    `json:"movie_id"`
    GenreID int  `json:"genre_id"`
}


func main() {
	// Connect to SQLite database
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create tables if not exist
	database.CreateUsersTable(db)

	// Load users from JSON file
	users, err := loadUsersFromJSON(filepath.Join("database", "users.json"))
	if err != nil {
		log.Fatalf("Error loading users: %v", err)
	}

	// Insert users into the database
	for _, user := range users {
		hashedPassword, err := database.HashPassword(user.PasswordHash)

		if err != nil {
			log.Printf("Failed to hash password for user %s: %v", user.Username, err)
			continue
		}

		// check if user with email exist before creating the user
		var existingEmail string
		err = db.QueryRow("SELECT email FROM users WHERE email = ?", user.Email).Scan(&existingEmail)
		if err == nil {
			log.Printf("User with email %s already exist. skipping user : %s\n", user.Email, user.Username)
			continue // email already exist skip to the next user
		}

		_, err = db.Exec("INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)", user.Username, user.Email, hashedPassword)
		if err != nil {
			log.Printf("Failed to insert user %s: %v", user.Username, err)
		} else {
			log.Printf("Inserted user: %s", user.Username)
		}
	}

	genres, err := loadGenresFromJSON(filepath.Join("database", "genres.json"))
	if err != nil {
		log.Fatalf("Error loading genres: %v", err)
	}
	for _, genre := range genres {
        var existingGenre string
        err = db.QueryRow("SELECT name FROM genres WHERE name = ?", genre.Name).Scan(&existingGenre)
        if err == nil {
            log.Printf("Genre with name %s already exist. skipping genre\n", genre.Name)
            continue
        }
		_, err = db.Exec("INSERT INTO genres (name) VALUES (?)", genre.Name)
		if err != nil {
			log.Printf("Failed to insert genre %s: %v", genre.Name, err)
		} else {
			log.Printf("Inserted genre: %s", genre.Name)
		}
	}

	// Load and insert movies
	movies, err := loadMoviesFromJSON(filepath.Join("database", "movies.json"))
	if err != nil {
		log.Fatalf("Error loading movies: %v", err)
	}
	for _, movie := range movies {
        var existingMovieTitle string
        err = db.QueryRow("SELECT title FROM movies WHERE title = ?", movie.Title).Scan(&existingMovieTitle)
        if err == nil {
            log.Printf("Movie with title %s already exist. skipping movie\n", movie.Title)
            continue
        }
		_, err = db.Exec("INSERT INTO movies (title, description, release_date, image_url) VALUES (?, ?, ?, ?)", movie.Title, movie.Description, movie.ReleaseDate, movie.ImageURL)
		if err != nil {
			log.Printf("Failed to insert movie %s: %v", movie.Title, err)
		} else {
			log.Printf("Inserted movie: %s", movie.Title)
		}
	}


	// Load and insert comments
	comments, err := loadCommentsFromJSON(filepath.Join("database", "comments.json"))
	if err != nil {
		log.Fatalf("Error loading comments: %v", err)
	}
	for _, comment := range comments {
		// Check if the comment already exists based on user_id, movie_id, and content
		var existingCommentID int
		err = db.QueryRow("SELECT comment_id FROM comments WHERE user_id = ? AND movie_id = ? AND content = ?", comment.UserID, comment.MovieID, comment.Content).Scan(&existingCommentID)
		if err == nil {
			log.Printf("Comment already exist skipping comment for user_id = %d, movie_id = %d\n", comment.UserID, comment.MovieID)
			continue
		}

		_, err = db.Exec("INSERT INTO comments (user_id, movie_id, content) VALUES (?, ?, ?)", comment.UserID, comment.MovieID, comment.Content)
		if err != nil {
			log.Printf("Failed to insert comment: %v", err)
		} else {
			log.Printf("Inserted comment: user_id = %d, movie_id = %d", comment.UserID, comment.MovieID)
		}
	}

	// Load and insert movie_genre data
	movieGenres, err := loadMovieGenreFromJSON(filepath.Join("database", "movie_genre.json"))
	if err != nil {
		log.Fatalf("Error loading movie_genre data: %v", err)
	}
	for _, mg := range movieGenres {
        var existingMovieGenreID int
        err = db.QueryRow("SELECT movie_id FROM movie_genre WHERE movie_id = ? AND genre_id = ?", mg.MovieID, mg.GenreID).Scan(&existingMovieGenreID)
        if err == nil {
            log.Printf("movie_genre with movie_id %d and genre_id %d already exist. skipping movie_genre \n", mg.MovieID, mg.GenreID)
            continue
        }
		_, err = db.Exec("INSERT INTO movie_genre (movie_id, genre_id) VALUES (?, ?)", mg.MovieID, mg.GenreID)
		if err != nil {
			log.Printf("Failed to insert movie_genre entry: %v", err)
		} else {
			log.Printf("Inserted movie_genre: movie_id = %d, genre_id = %d", mg.MovieID, mg.GenreID)
		}
	}

	fmt.Println("User data inserted successfully.")
	fmt.Print("Start Server...")
	StartServer()
}

// Function to start a simple HTTP server
func StartServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Welcome to the forum!")
	})

	log.Println("Starting server on :8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// Function to load users from JSON file
func loadUsersFromJSON(filename string) ([]User, error) {
	var users []User
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// Function to load movies from JSON file
func loadMoviesFromJSON(filename string) ([]Movie, error) {
	var movies []Movie
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &movies)
	if err != nil {
		return nil, err
	}
	return movies, nil
}

// Function to load genres from JSON file
func loadGenresFromJSON(filename string) ([]Genre, error) {
    var genres []Genre
    data, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    err = json.Unmarshal(data, &genres)
    if err != nil {
        return nil, err
    }
    return genres, nil
}


// Function to load comments from JSON file
func loadCommentsFromJSON(filename string) ([]Comment, error) {
	var comments []Comment
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &comments)
	if err != nil {
		return nil, err
	}
	return comments, nil
}

// Function to load movie_genre data from JSON file
func loadMovieGenreFromJSON(filename string) ([]MovieGenre, error) {
    var movieGenres []MovieGenre
    data, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    err = json.Unmarshal(data, &movieGenres)
    if err != nil {
        return nil, err
    }
    return movieGenres, nil
}
