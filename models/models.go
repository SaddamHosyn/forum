package models


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


// MovieGenre struct for JSON data, this is not directly insertable.
type MovieGenre struct {
	MovieID int `json:"movie_id"`
	GenreID int `json:"genre_id"`
}

// MovieWithGenres struct for showing the movies with their genres
type MovieWithGenres struct {
    Movie      Movie
    Genres     []string
}
// GenreWithMovies struct for showing the genres with their movies
type GenreWithMovies struct {
    Genre  Genre
    Movies []string
}



// Comment struct for JSON data
type Comment struct {
	UserID  int    `json:"user_id"`
	MovieID int    `json:"movie_id"`
	Content string `json:"content"`
}
