package database

import (
	"database/sql"
	"fmt"
	"log"
	"forum-go/models"

)

// InitDB initializes the database and populates tables if empty
func InitDB(db *sql.DB) {
	// Create all tables
	CreateAllTables(db)

	// Populate tables if empty
	err := PopulateData(db)
	if err != nil {
		log.Println("Error populating database:", err)
	}
}

// CreateAllTables ensures all necessary tables are created
func CreateAllTables(database *sql.DB) {
	CreateUsersTable(database)
	CreateMoviesTable(database)
	CreateGenresTable(database)
	CreateCommentsTable(database)
	CreateMovieGenreTable(database)
}

// CreateUsersTable defines the schema for users
func CreateUsersTable(database *sql.DB) {
	createUsersTableQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY,
		username TEXT UNIQUE NOT NULL, 
		email TEXT UNIQUE NOT NULL, 
		password_hash TEXT UNIQUE NOT NULL, -- Hashed password_hash
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := database.Exec(createUsersTableQuery); err != nil {
		fmt.Println("Error creating users table:", err)
		return
	}
	fmt.Println("Users table created successfully or already exists.")
}

// CreateMoviesTable defines the schema for movies
func CreateMoviesTable(database *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS movies (
		id INTEGER PRIMARY KEY,
		title TEXT NOT NULL, 
		description TEXT,
		release_date TEXT,
		image_url TEXT
	);`
	_, err := database.Exec(query)
	if err != nil {
		fmt.Println("Error creating movies table:", err)
	}
}

// CreateGenresTable defines the schema for genres
func CreateGenresTable(database *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS genres (
		id INTEGER PRIMARY KEY,
		name TEXT UNIQUE NOT NULL
	);`
	_, err := database.Exec(query)
	if err != nil {
		fmt.Println("Error creating genres table:", err)
	}
}

// CreateCommentsTable defines the schema for comments
func CreateCommentsTable(database *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS comments (
		id INTEGER PRIMARY KEY,
		user_id INTEGER,
		movie_id INTEGER,
		content TEXT,
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (movie_id) REFERENCES movies(id)
	);`
	_, err := database.Exec(query)
	if err != nil {
		fmt.Println("Error creating comments table:", err)
	}
}

// CreateMovieGenreTable defines the linking table for movies and genres
func CreateMovieGenreTable(database *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS movie_genre (
		movie_id INTEGER,
		genre_id INTEGER,
		FOREIGN KEY (movie_id) REFERENCES movies(id),
		FOREIGN KEY (genre_id) REFERENCES genres(id),
		PRIMARY KEY (movie_id, genre_id)
	);`
	_, err := database.Exec(query)
	if err != nil {
		fmt.Println("Error creating movie_genre table:", err)
	}
}

// PopulateData inserts initial data only if tables are empty
func PopulateData(db *sql.DB) error {
	var count int

	// Check if the users table has data
	row := db.QueryRow("SELECT COUNT(*) FROM users;")
	err := row.Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		fmt.Println("Database already populated.")
		return nil
	}

	// Insert Users
	_, err = db.Exec(`INSERT INTO users (username, email, password_hash) VALUES 
		('Joon', 'joon@example.com', 'P@$$wOrd1'),
		('gigi', 'gigi@example.com', 'MyS3cretPass!'),
		('mayre', 'mayre@example.com', 'FlyT0TheMoon'),
		('sagyngoogle', 'sagyngoogle@example.com', 'SecondStar');`)
	if err != nil {
		return err
	}

	// Insert Movies
	_, err = db.Exec(`INSERT INTO movies (title, description, release_date, image_url) VALUES 
		('Inception', 'A mind-bending thriller', '2010-07-16', 'inception.jpg'),
		('The Matrix', 'A sci-fi classic', '1999-03-31', 'matrix.jpg');`)
	if err != nil {
		return err
	}

	// Insert Genres
	_, err = db.Exec(`INSERT INTO genres (name) VALUES ('Action'), ('Sci-Fi'), ('Drama');`)
	if err != nil {
		return err
	}

	// Insert Comments
	_, err = db.Exec(`INSERT INTO comments (user_id, movie_id, content) VALUES 
		(1, 1, 'Amazing movie!'),
		(2, 2, 'A masterpiece!');`)
	if err != nil {
		return err
	}

	// Insert into movie_genre (linking movies and genres)
	_, err = db.Exec(`INSERT INTO movie_genre (movie_id, genre_id) VALUES 
		(1, 2), -- Inception (Sci-Fi)
		(2, 2), -- The Matrix (Sci-Fi)
		(2, 1); -- The Matrix (Action)
	`)
	if err != nil {
		return err
	}

	fmt.Println("Database populated successfully!")
	return nil
}

// Function to fetch and return movies with their associated genres
func GetMoviesWithGenres(db *sql.DB) ([]models.MovieWithGenres, error) {
	rows, err := db.Query(`
        SELECT m.movie_id, m.title, m.description, m.release_date, m.image_url, g.name
        FROM movies m
        LEFT JOIN movie_genre mg ON m.movie_id = mg.movie_id
        LEFT JOIN genres g ON mg.genre_id = g.genre_id
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	movieMap := make(map[int]models.MovieWithGenres)

	for rows.Next() {
		var movieID int
		var movieTitle, movieDescription, movieReleaseDate, movieImageURL string
		var genreName string
		err := rows.Scan(&movieID, &movieTitle, &movieDescription, &movieReleaseDate, &movieImageURL, &genreName)
		if err != nil {
			return nil, err
		}
		if _, ok := movieMap[movieID]; !ok {
			movieMap[movieID] = models.MovieWithGenres{
				Movie: models.Movie{
					Title:       movieTitle,
					Description: movieDescription,
					ReleaseDate: movieReleaseDate,
					ImageURL:    movieImageURL,
				},
				Genres: []string{},
			}
		}
		if genreName != "" {
			currentMovie := movieMap[movieID]
			currentMovie.Genres = append(currentMovie.Genres, genreName)
			movieMap[movieID] = currentMovie
		}

	}
	var moviesWithGenres []models.MovieWithGenres
	for _, movieWithGenre := range movieMap {
		moviesWithGenres = append(moviesWithGenres, movieWithGenre)
	}

	return moviesWithGenres, nil
}

// Function to fetch and return genres with their associated movies
func GetGenresWithMovies(db *sql.DB) ([]models.GenreWithMovies, error) {
	rows, err := db.Query(`
    SELECT g.genre_id, g.name, m.title
        FROM genres g
        LEFT JOIN movie_genre mg ON g.genre_id = mg.genre_id
        LEFT JOIN movies m ON mg.movie_id = m.movie_id
    `)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	genreMap := make(map[int]models.GenreWithMovies)

	for rows.Next() {
		var genreID int
		var genreName string
		var movieTitle sql.NullString
		err := rows.Scan(&genreID, &genreName, &movieTitle)
		if err != nil {
			return nil, err
		}
		if _, ok := genreMap[genreID]; !ok {
			genreMap[genreID] = models.GenreWithMovies{
				Genre: models.Genre{
					Name: genreName,
				},
				Movies: []string{},
			}
		}
		if movieTitle.Valid {
			currentGenre := genreMap[genreID]
			currentGenre.Movies = append(currentGenre.Movies, movieTitle.String)
			genreMap[genreID] = currentGenre
		}
	}

	var genresWithMovies []models.GenreWithMovies
	for _, genreWithMovie := range genreMap {
		genresWithMovies = append(genresWithMovies, genreWithMovie)
	}

	return genresWithMovies, nil
}
