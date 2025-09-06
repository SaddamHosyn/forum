package model

import (
	"database/sql"
	"time"
)

var DB *sql.DB

type HomePageData struct {
	Posts      []Post
	Categories []Category
	User       *User
	Success    bool
	IsLoggedIn bool
}
//todo: why pointer to user not to others?

type SubmitPostData struct {
	Categories []Category
	User       *User
}

type ProfileData struct {
	Title      string
	User       *User
	Posts      []*Post
	Categories []Category
	IsLoggedIn bool
	Success    bool
	Error      string
}


type Category struct {
	ID    string
	Name  string
	Emoji string
}

type User struct {
	ID            int
	Username      string
	Email         string
	Password      string
	SessionToken  sql.NullString
	SessionExpiry sql.NullTime
	CreatedAt     time.Time
}
// todo: why nullstring and nulltime?

type Post struct {
	ID         int
	Author     string // Added field
	Title      string
	Content    string
	UserID     int    // Fixed casing
	Categories string // Added field
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Upvotes    int
	Downvotes  int
	Comments   []Comment
}
//todo: why comments are slice of strings?

type Comment struct {
	ID        int
	Content   string
	Author    string // Username from users table
	UserID    int
	PostID    int
	CreatedAt time.Time
	Upvotes   int
	Downvotes int
	TimeAgo   string
}

type Votes struct {
	ID     int
	Vote   int // 1 for upvote, -1 for downvote
	UserID int
	PostID int
}
