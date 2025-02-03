PRAGMA foreign_keys = ON;

-- Users Table (unchanged)
CREATE TABLE users (
    user_id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE COLLATE BINARY,
    email TEXT NOT NULL UNIQUE COLLATE BINARY,
    password_hash TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Movies Table (unchanged)
CREATE TABLE movies (
    movie_id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    release_date TEXT,
    image_url TEXT
);

-- Genres Table (unchanged)
CREATE TABLE genres (
    genre_id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE
);

-- Comments Table (unchanged)
CREATE TABLE comments (
    comment_id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    movie_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (movie_id) REFERENCES movies(movie_id) ON DELETE CASCADE
);

-- Movie_Genre Join Table (unchanged)
CREATE TABLE movie_genre (
    movie_id INTEGER NOT NULL,
    genre_id INTEGER NOT NULL,
    FOREIGN KEY (movie_id) REFERENCES movies(movie_id) ON DELETE CASCADE,
    FOREIGN KEY (genre_id) REFERENCES genres(genre_id) ON DELETE CASCADE,
    PRIMARY KEY (movie_id, genre_id)
);

-- Movie Posts Table
CREATE TABLE movie_posts (
    post_id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    movie_id INTEGER NOT NULL,
    post_text TEXT,  -- Optional text content of the post
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (movie_id) REFERENCES movies(movie_id) ON DELETE CASCADE
);

-- Movie Ratings (Likes and Dislikes) Table
CREATE TABLE movie_ratings (
    rating_id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    movie_id INTEGER NOT NULL,
    rating INTEGER NOT NULL,  -- 1 for like, -1 for dislike, 0 for neutral
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (movie_id) REFERENCES movies(movie_id) ON DELETE CASCADE,
    UNIQUE(user_id, movie_id) -- Ensure a user can only have one rating per movie
);

-- Indexes (Added for the new tables)
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_comments_user_id ON comments(user_id);
CREATE INDEX idx_comments_movie_id ON comments(movie_id);
CREATE INDEX idx_movie_genre_movie_id ON movie_genre(movie_id);
CREATE INDEX idx_movie_genre_genre_id ON movie_genre(genre_id);
CREATE INDEX idx_movie_posts_user_id ON movie_posts(user_id);
CREATE INDEX idx_movie_posts_movie_id ON movie_posts(movie_id);
CREATE INDEX idx_movie_ratings_user_id ON movie_ratings(user_id);
CREATE INDEX idx_movie_ratings_movie_id ON movie_ratings(movie_id);
