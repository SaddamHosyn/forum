package models

import "time"

// User struct matches database schema, JSON tags are removed
type User struct {
    UserID       int
	Username     string
	Email        string
	PasswordHash string
    SessionID    string
}

// Movie struct for database schema, JSON tags are removed
type Movie struct {
	Title       string
	Description string
	ReleaseDate string
	ImageURL    string
}
type Genre struct {
	Name string
}

// MovieGenre struct for database schema, JSON tags are removed
type MovieGenre struct {
	MovieID int
	GenreID int
}

// MovieWithGenres struct for showing the movies with their genres
type MovieWithGenres struct {
	Movie  Movie
	Genres []string
}

// GenreWithMovies struct for showing the genres with their movies
type GenreWithMovies struct {
	Genre  Genre
	Movies []string
}

// Comment struct for database schema, JSON tags are removed
type Comment struct {
	UserID  int
	MovieID int
	Content string
}

type MoviePost struct {
	PostID    int
	UserID    int
	MovieID   int
	PostText  string
	CreatedAt time.Time
}

// MovieRating struct for database schema
type MovieRating struct {
	RatingID  int
	UserID    int
	MovieID   int
	Rating    int
	CreatedAt time.Time
}
