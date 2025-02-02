package database

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"forum-go/models"

	"golang.org/x/crypto/bcrypt"
)

// CreateUsersTable creates the users table if it doesn't exist
func CreateUsersTable(db *sql.DB) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			user_id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL UNIQUE COLLATE BINARY,
			email TEXT NOT NULL UNIQUE COLLATE BINARY,
			password_hash TEXT NOT NULL,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			session_id TEXT UNIQUE
		);
    `)
	if err != nil {
		log.Fatalf("Error creating users table: %v", err)
	}
	fmt.Println("Users table created successfully.")

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS movies (
            movie_id INTEGER PRIMARY KEY AUTOINCREMENT,
            title TEXT NOT NULL,
            description TEXT NOT NULL,
            release_date TEXT,  -- SQLite does not have a specific DATE, using TEXT
            image_url TEXT
        );
    `)
	if err != nil {
		log.Fatalf("Error creating movies table: %v", err)
	}
	fmt.Println("Movies table created successfully.")

	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS genres (
            genre_id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL UNIQUE
        );
    `)
	if err != nil {
		log.Fatalf("Error creating genres table: %v", err)
	}
	fmt.Println("Genres table created successfully.")

	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS comments (
            comment_id INTEGER PRIMARY KEY AUTOINCREMENT,
            user_id INTEGER NOT NULL,
            movie_id INTEGER NOT NULL,
            content TEXT NOT NULL,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
            FOREIGN KEY (movie_id) REFERENCES movies(movie_id) ON DELETE CASCADE
        );
    `)
	if err != nil {
		log.Fatalf("Error creating comments table: %v", err)
	}
	fmt.Println("Comments table created successfully.")

	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS movie_genre (
            movie_id INTEGER NOT NULL,
            genre_id INTEGER NOT NULL,
            FOREIGN KEY (movie_id) REFERENCES movies(movie_id) ON DELETE CASCADE,
            FOREIGN KEY (genre_id) REFERENCES genres(genre_id) ON DELETE CASCADE,
            PRIMARY KEY (movie_id, genre_id)
        );
    `)
	if err != nil {
		log.Fatalf("Error creating movie_genre table: %v", err)
	}
	fmt.Println("Movie_genre table created successfully.")

	//Indexes
	_, err = db.Exec("CREATE INDEX IF NOT EXISTS idx_users_username ON users(username)")
	if err != nil {
		log.Fatalf("Failed to create index idx_users_username %v", err)
	}
	_, err = db.Exec("CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)")
	if err != nil {
		log.Fatalf("Failed to create index idx_users_email %v", err)
	}
	_, err = db.Exec("CREATE INDEX IF NOT EXISTS idx_comments_user_id ON comments(user_id)")
	if err != nil {
		log.Fatalf("Failed to create index idx_comments_user_id %v", err)
	}
	_, err = db.Exec("CREATE INDEX IF NOT EXISTS idx_comments_movie_id ON comments(movie_id)")
	if err != nil {
		log.Fatalf("Failed to create index idx_comments_movie_id %v", err)
	}
	_, err = db.Exec("CREATE INDEX IF NOT EXISTS idx_movie_genre_movie_id ON movie_genre(movie_id)")
	if err != nil {
		log.Fatalf("Failed to create index idx_movie_genre_movie_id %v", err)
	}
	_, err = db.Exec("CREATE INDEX IF NOT EXISTS idx_movie_genre_genre_id ON movie_genre(genre_id)")
	if err != nil {
		log.Fatalf("Failed to create index idx_movie_genre_genre_id %v", err)
	}

	fmt.Println("Indexes created successfully")

}

// HashPassword hashes the password
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error hashing password: %w", err)
	}
	return string(hashedPassword), nil
}

//ComparePassword compares the hash with the password
func ComparePassword(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err
}

// Function to get movies with genres
func GetMoviesWithGenres(db *sql.DB) ([]models.MovieWithGenres, error) {
	rows, err := db.Query(`
		SELECT m.movie_id, m.title, m.description, m.release_date, m.image_url,
GROUP_CONCAT(g.name) AS genres
FROM movies m
LEFT JOIN movie_genre mg ON m.movie_id = mg.movie_id
LEFT JOIN genres g ON mg.genre_id = g.genre_id
GROUP BY m.movie_id
	`)
	if err != nil {
		return nil, fmt.Errorf("error querying movies with genres: %w", err)
	}
	defer rows.Close()

	var moviesWithGenres []models.MovieWithGenres
	for rows.Next() {
		var movie models.Movie
		var movieID int
		var genresString sql.NullString
		err := rows.Scan(&movieID, &movie.Title, &movie.Description, &movie.ReleaseDate, &movie.ImageURL, &genresString)
		if err != nil {
			return nil, fmt.Errorf("error scanning movie with genres: %w", err)
		}

		var genres []string
		if genresString.Valid {
			genres = splitString(genresString.String, ",")
		}
		moviesWithGenres = append(moviesWithGenres, models.MovieWithGenres{Movie: movie, Genres: genres})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating movie with genres rows: %w", err)
	}

	return moviesWithGenres, nil
}

func GetGenresWithMovies(db *sql.DB) ([]models.GenreWithMovies, error) {
	rows, err := db.Query(`
        SELECT 
            g.genre_id, g.name,
            GROUP_CONCAT(DISTINCT m.title) AS movies
        FROM genres g
        LEFT JOIN movie_genre mg ON g.genre_id = mg.genre_id
        LEFT JOIN movies m ON mg.movie_id = m.movie_id
        GROUP BY g.genre_id
    `)
	if err != nil {
		return nil, fmt.Errorf("error querying genres with movies: %w", err)
	}
	defer rows.Close()

	var genresWithMovies []models.GenreWithMovies
	for rows.Next() {
		var genre models.Genre
		var genreID int
		var moviesString sql.NullString
		err := rows.Scan(&genreID, &genre.Name, &moviesString)
		if err != nil {
			return nil, fmt.Errorf("error scanning genre with movies: %w", err)
		}

		var movies []string
		if moviesString.Valid {
			movies = splitString(moviesString.String, ",")
		}
		genresWithMovies = append(genresWithMovies, models.GenreWithMovies{Genre: genre, Movies: movies})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating genre with movies rows: %w", err)
	}

	return genresWithMovies, nil
}

// Helper function to split comma-separated string into a slice
func splitString(str string, separator string) []string {
	var result []string
	if str != "" {
		for _, val := range strings.Split(str, separator) {
			result = append(result, strings.TrimSpace(val))
		}
	}
	return result
}

func ValidateSessionID(sessionID string, userID *int) error {
    db, err := sql.Open("sqlite3", "database.db")
    if err != nil {
        return fmt.Errorf("failed to open database: %w", err)
    }
    defer db.Close()
	err = db.QueryRow("SELECT user_id FROM users WHERE session_id = ?", sessionID).Scan(userID)
	if err != nil {
		return fmt.Errorf("invalid session: %w", err)
	}
    return nil
}
