package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"forum-go/database"
	. "forum-go/models"

	_ "github.com/mattn/go-sqlite3"
)

const sessionCookieName = "forum_session"
const userCookieName = "forum_user"
const sessionExpiry = 30 * time.Minute


func main() {
	// Connect to SQLite database
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	//check if the SETUP_DONE variable is set.
	if os.Getenv("SETUP_DONE") != "true" {
		// if it is not set setup database.
		fmt.Print("Doing initial setup...")
		setupDatabase(db)
		// set the environment variable so we don't have to setup the database every time.
		os.Setenv("SETUP_DONE", "true")
		fmt.Println("Initial setup done")
	}
	// Start the server
	fmt.Print("Starting server...")
	StartServer(db)
}

func setupDatabase(db *sql.DB) {
	// Create tables if not exist
	database.CreateUsersTable(db)

	// Load genres from JSON file
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

	// Fetch and log movies with genres
	// Fetch and log movies with genres after setup
	moviesWithGenres, err := database.GetMoviesWithGenres(db)
	if err != nil {
		log.Println("Error getting movies with genres:", err)
	}
	fmt.Println("\nMovies with Genres:")
	for _, mwg := range moviesWithGenres {
		if len(mwg.Genres) > 0 { // Add this condition
			fmt.Printf("  - Movie: %s, Genres: %v\n", mwg.Movie.Title, mwg.Genres)
		}
	}


	// Fetch and log genres with movies
	genresWithMovies, err := database.GetGenresWithMovies(db)
	if err != nil {
		log.Println("Error getting genres with movies:", err)
	}
	fmt.Println("\nGenres with Movies:")
	for _, gwm := range genresWithMovies {
		fmt.Printf("  - Genre: %s, Movies: %v\n", gwm.Genre.Name, gwm.Movies)
	}

	fmt.Println("Movie data inserted successfully.")
}

// Function to start a simple HTTP server
func StartServer(db *sql.DB) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Welcome to the forum!")
	})

	http.HandleFunc("/register", registerHandler(db))     // User registration
	http.HandleFunc("/login", loginHandler(db))            // User login
	http.HandleFunc("/comment", authMiddleware(commentHandler(db))) // Handling new comment, secured with auth

	log.Println("Starting server on :8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// Function to handle new user registration
func registerHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		// read the request body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// unmarshall the json
		var newUser User
		err = json.Unmarshal(body, &newUser)
		if err != nil {
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}

		// hash the password
		hashedPassword, err := database.HashPassword(newUser.PasswordHash)
		if err != nil {
			http.Error(w, "Error hashing the password", http.StatusInternalServerError)
			return
		}

		// check if user with the email exist
		var existingEmail string
		err = db.QueryRow("SELECT email FROM users WHERE email = ?", newUser.Email).Scan(&existingEmail)
		if err == nil {
			http.Error(w, "Email already exist", http.StatusConflict)
			return
		}
		// generate a unique session ID
		sessionID := uuid.New().String()

		// insert the user in the database, with session id
		_, err = db.Exec("INSERT INTO users (username, email, password_hash, session_id) VALUES (?, ?, ?, ?)", newUser.Username, newUser.Email, hashedPassword, sessionID)
		if err != nil {
			log.Printf("Failed to insert user %s: %v", newUser.Username, err)
			http.Error(w, "Failed to register user", http.StatusInternalServerError)
			return
		}

		//set session cookie
		setSessionCookie(w, sessionID)

		//fetch user data and cache in cookie.
		err = setUserCookie(db, w, sessionID)
		if err != nil {
			http.Error(w, "Failed to set user cookie", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintln(w, "User Registered Successfully")

	}
}

// Function to handle user login
func loginHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }

        body, err := ioutil.ReadAll(r.Body)
        if err != nil {
            http.Error(w, "Failed to read request body", http.StatusBadRequest)
            return
        }
        defer r.Body.Close()

        var loginUser User
        err = json.Unmarshal(body, &loginUser)
        if err != nil {
            http.Error(w, "Invalid JSON format", http.StatusBadRequest)
            return
        }

		//check if the user exists by email
        var storedHash string
		var userID int
		var sessionID string
        err = db.QueryRow("SELECT user_id, password_hash, session_id FROM users WHERE email = ?", loginUser.Email).Scan(&userID, &storedHash, &sessionID)
        if err != nil {
            http.Error(w, "Invalid email or password", http.StatusUnauthorized)
            return
        }

		// check if the passwords match
        if err := database.ComparePassword(storedHash, loginUser.PasswordHash); err != nil {
            http.Error(w, "Invalid email or password", http.StatusUnauthorized)
            return
        }

		if sessionID == "" {
			sessionID = uuid.New().String()

			_, err = db.Exec("UPDATE users SET session_id = ? WHERE user_id = ?", sessionID, userID)
			if err != nil {
				http.Error(w, "Failed to update session ID", http.StatusInternalServerError)
				return
			}

		}
		//set session cookie
		setSessionCookie(w, sessionID)


        // Fetch user data and cache in cookie.
		err = setUserCookie(db,w, sessionID)
		if err != nil {
			http.Error(w, "Failed to set user cookie", http.StatusInternalServerError)
			return
		}


        w.WriteHeader(http.StatusOK)
        fmt.Fprintln(w, "User Logged In Successfully")
    }
}

// Function to handle new comments
func commentHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }

        // Read the request body
        body, err := ioutil.ReadAll(r.Body)
        if err != nil {
            http.Error(w, "Failed to read request body", http.StatusBadRequest)
            return
        }
        defer r.Body.Close()

		//get the user from the cookie
		userCookie, err := r.Cookie(userCookieName)
		if err != nil {
			http.Error(w, "User not authenticated", http.StatusUnauthorized)
			return
		}

		var cachedUser User
		err = json.Unmarshal([]byte(userCookie.Value), &cachedUser)
		if err != nil {
			http.Error(w, "Failed to read user from cookie", http.StatusUnauthorized)
			return
		}


        // Unmarshal the JSON into a Comment struct
        var newComment Comment
        err = json.Unmarshal(body, &newComment)
        if err != nil {
            http.Error(w, "Invalid JSON format", http.StatusBadRequest)
            return
        }
        // Check if the user exist, from the cookie
        if cachedUser.UserID == 0 {
			http.Error(w, "Invalid user id", http.StatusBadRequest)
            return
		}

        // Check if the movie exist
        var existingMovieID int
        err = db.QueryRow("SELECT movie_id FROM movies WHERE movie_id = ?", newComment.MovieID).Scan(&existingMovieID)
        if err != nil {
            http.Error(w, "Movie does not exist, can't add comment", http.StatusBadRequest)
            return
        }

		// Check if the comment already exists
		var existingCommentID int
		err = db.QueryRow("SELECT comment_id FROM comments WHERE user_id = ? AND movie_id = ? AND content = ?", cachedUser.UserID, newComment.MovieID, newComment.Content).Scan(&existingCommentID)
		if err == nil {
			http.Error(w, "Comment already exist", http.StatusConflict)
			return
		}


        // Insert the comment into the database
        _, err = db.Exec("INSERT INTO comments (user_id, movie_id, content) VALUES (?, ?, ?)", cachedUser.UserID, newComment.MovieID, newComment.Content)
        if err != nil {
            log.Printf("Failed to insert comment: %v", err)
            http.Error(w, "Failed to add comment", http.StatusInternalServerError)
            return
        }
        w.WriteHeader(http.StatusCreated)
        fmt.Fprintln(w, "Comment added successfully!")
    }
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

// Function to set the session cookie
func setSessionCookie(w http.ResponseWriter, sessionID string) {
	cookie := &http.Cookie{
		Name:     sessionCookieName,
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure: false,
		SameSite: http.SameSiteStrictMode,
		Expires: time.Now().Add(sessionExpiry),
	}
	http.SetCookie(w, cookie)
}

// Function to set the user cookie
func setUserCookie(db *sql.DB, w http.ResponseWriter, sessionID string) error {
	var user User
	err := db.QueryRow("SELECT user_id, username, email FROM users WHERE session_id = ?", sessionID).Scan(&user.UserID, &user.Username, &user.Email)
	if err != nil {
		return err
	}

	userJSON, err := json.Marshal(user)
	if err != nil {
		return err
	}

	cookie := &http.Cookie {
		Name: userCookieName,
		Value: string(userJSON),
		Path: "/",
		HttpOnly: true,
		Secure: false,
		SameSite: http.SameSiteStrictMode,
		Expires: time.Now().Add(sessionExpiry),
	}
	http.SetCookie(w, cookie)
	return nil
}

// Authentication middleware to validate session IDs
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, err := r.Cookie(sessionCookieName)
		if err != nil {
			http.Error(w, "Session not found", http.StatusUnauthorized)
			return
		}

		sessionID := sessionCookie.Value

		// Verify the session ID exists
		var userID int
		err = database.ValidateSessionID(sessionID, &userID)
		if err != nil {
			http.Error(w, "Invalid session", http.StatusUnauthorized)
			return
		}
        // session is valid, continue to the protected handler
		next.ServeHTTP(w, r)
	}
}
