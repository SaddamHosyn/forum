package database

import (
	"fmt"
)

/* type DB struct {
	*sql.DB
} */

func createTables() error {
	// Create users table
	_, err := DB.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            username TEXT UNIQUE NOT NULL,
            email TEXT UNIQUE NOT NULL,
            password_hash TEXT NOT NULL,
            session_token TEXT UNIQUE,
            session_expiry DATETIME,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
        )
    `)
	if err != nil {
		return fmt.Errorf("error creating users table: %v", err)
	}

	// Create categories table
	_, err = DB.Exec(`
        CREATE TABLE IF NOT EXISTS categories (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT UNIQUE NOT NULL,
            emoji TEXT UNIQUE NOT NULL,
            UNIQUE(name,emoji)

        )
    `)
	if err != nil {
		return fmt.Errorf("error creating categories table: %v", err)
	}

	// Create posts table
	_, err = DB.Exec(`
        CREATE TABLE IF NOT EXISTS posts (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            title TEXT NOT NULL UNIQUE,
            content TEXT NOT NULL,
            user_id  INTEGER NOT NULL,
            categories TEXT NOT NULL,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY(user_id) REFERENCES users(id)
        )
    `)
	if err != nil {
		return fmt.Errorf("error creating posts table: %v", err)
	}

	// Create comments table
	_, err = DB.Exec(`
        CREATE TABLE IF NOT EXISTS comments (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            content TEXT NOT NULL,
            user_id INTEGER NOT NULL,
            post_id INTEGER NOT NULL,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY(user_id) REFERENCES users(id),
            FOREIGN KEY(post_id) REFERENCES posts(id)
        )
    `)
	if err != nil {
		return fmt.Errorf("error creating comments table: %v", err)
	}

	// Create posts_categories table
	_, err = DB.Exec(`
        CREATE TABLE IF NOT EXISTS posts_categories (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            post_id INTEGER NOT NULL,
            categories TEXT NOT NULL,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
            UNIQUE (post_id, categories)
        )
    `)
	if err != nil {
		return fmt.Errorf("error creating posts_categories table: %v", err)
	}

	_, err = DB.Exec(`
        CREATE TABLE IF NOT EXISTS sessions (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            user_id INTEGER NOT NULL,
            session_token TEXT NOT NULL UNIQUE,
            session_expiry DATETIME NOT NULL,
            FOREIGN KEY (user_id) REFERENCES users(id)
        )
    `)
	if err != nil {
		return fmt.Errorf("error creating sessions table: %v", err)
	}

	_, err = DB.Exec(`
        CREATE TABLE IF NOT EXISTS votes (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            user_id INTEGER NOT NULL,
            post_id INTEGER,
            comment_id INTEGER,
            vote INTEGER NOT NULL CHECK (vote IN (1, 0, -1)),
            FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
            FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
            FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE,
            UNIQUE (user_id, post_id),
            UNIQUE (user_id, comment_id)
        )   
    `)
	if err != nil {
		return fmt.Errorf("error creating sessions table: %v", err)
	}

	// Enable foreign key support
	_, err = DB.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		return fmt.Errorf("error enabling foreign key support: %v", err)
	}

	return nil
}

func insertCategories() error {
	_, err := DB.Exec(`
    INSERT OR IGNORE INTO categories (name, emoji) VALUES
    ('Action', '💥'), ('Adventure', '🌄'), ('Animation', '🧚'), 
    ('Biography', '📚'), ('Comedy', '😂'), ('Crime', '🕵️'), ('Documentary', '🎥'), 
    ('Drama', '🎭'), ('Fantasy', '🧙'), ('Horror', '👻'), ('Mystery', '🔍'), 
    ('Romance', '❤️'), ('Sci-Fi', '🚀'), ('Thriller', '😱'), ('Western', '🤠')
`)
	if err != nil {
		return fmt.Errorf("error inserting categories: %v", err)
	}
	return nil
}
func insertUsers() error {
	_, err := DB.Exec(`
    INSERT OR IGNORE INTO users (username, email, password_hash) VALUES
    ('admin', 'admin@admin.com', '$2a$10$ryPUUMn0CPeuNh.NpQZOwuyoymt1sdzXrePhSeYArwv9puWlg1mF2'),
    ('Mama', 'mama@yahoo.com', '$2a$10$bfVNqrSBscGyfsGMSyEvaOCRbBbC54I2Lht5XuaBLiZKcdgoIRJQO'),
    ('batman', 'batman@batman.com', '$2a$10$1ZAK4MxQuwCJZGqhpBBzPOMoDDeGob..uwEIIO9YsHpqx8qXPNH8u')
`)

	if err != nil {
		return fmt.Errorf("error inserting users: %v", err)
	}
	return nil
}
func insertPosts() error {
	_, err := DB.Exec(`
    INSERT OR IGNORE INTO posts (title, content, user_id, categories) VALUES
    ('The Thrilling Ride of "Quantum Horizon"', 'Quantum Horizon blends cutting-edge effects with a gripping narrative. The zero-gravity fights are breathtaking.', 2, 'Sci-Fi,Action,Thriller'),
    ('Laughing Through Time: A Hilarious Adventure', 'A refreshing take on time-travel comedy. Clever writing and impeccable timing had me in stitches.', 1, 'Comedy,Adventure,Sci-Fi'),
    ('Whispers of the West: A Haunting Frontier Tale', 'This unconventional Western infuses horror into a frontier setting. Eerie atmosphere keeps viewers on edge.', 3, 'Western,Horror,Mystery'),
    ('Brushstrokes of Genius: A Compelling Artist''s Biography', 'A meticulously crafted documentary about painter Isabella Rossi. Balances interviews with stunning visuals of her work.', 2, 'Documentary,Biography')
`)

	if err != nil {
		return fmt.Errorf("error inserting posts: %v", err)
	}
	return nil
}
