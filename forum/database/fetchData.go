package database

import (
	"database/sql"
	"errors"
	"fmt"
	"forum-go/model"
)

func FetchPosts() ([]model.Post, error) {

	if DB == nil {
		return nil, errors.New("database connection is nil")
	}

	query := `
    SELECT p.id, u.username, p.title, p.content, p.user_id,
			p.categories, p.created_at, p.updated_at,
		    COALESCE(SUM(CASE WHEN v.vote = 1 THEN 1 ELSE 0 END), 0) AS upvotes,
            COALESCE(SUM(CASE WHEN v.vote = -1 THEN 1 ELSE 0 END), 0) AS downvotes
    FROM posts p
    JOIN users u ON p.user_id = u.id
	LEFT JOIN votes v ON p.id = v.post_id
    GROUP BY p.id
    ORDER BY p.created_at ASC
    LIMIT 10
`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying posts: %w", err)
	}
	defer rows.Close()

	var posts []model.Post
	for rows.Next() {
		var p model.Post
		err := rows.Scan(
			&p.ID,
			&p.Author,
			&p.Title,
			&p.Content,
			&p.UserID,
			&p.Categories,
			&p.CreatedAt,
			&p.UpdatedAt,
			&p.Upvotes,
			&p.Downvotes,
		)
		if err != nil {
			return nil, fmt.Errorf("error querying posts: %w", err)
		}
		posts = append(posts, p)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error after scanning posts: %w", err)
	}
	return posts, nil
}

func FetchCategories() ([]model.Category, error) {
	query := "SELECT id, name, emoji FROM categories ORDER BY name ASC"
	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []model.Category
	for rows.Next() {
		var c model.Category
		err := rows.Scan(&c.ID, &c.Name, &c.Emoji)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func FetchCommentsByPostID(postID int) ([]model.Comment, error) {
	query := `
        SELECT c.id, c.content, u.username, c.user_id, c.post_id, c.created_at,
                COALESCE(SUM(CASE WHEN v.vote = 1 THEN 1 ELSE 0 END), 0) AS upvotes,
                COALESCE(SUM(CASE WHEN v.vote = -1 THEN 1 ELSE 0 END), 0) AS downvotes
        FROM comments c
        JOIN users u ON c.user_id = u.id
        LEFT JOIN votes v ON c.id = v.comment_id
        WHERE c.post_id = ?
        GROUP BY c.id
        ORDER BY c.created_at ASC
    `

	rows, err := DB.Query(query, postID)
	if err != nil {
		return nil, fmt.Errorf("error querying comments: %w", err)
	}
	defer rows.Close()

	var comments []model.Comment
	for rows.Next() {
		var c model.Comment
		err := rows.Scan(
			&c.ID,
			&c.Content,
			&c.Author,
			&c.UserID,
			&c.PostID,
			&c.CreatedAt,
			&c.Upvotes,
			&c.Downvotes,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning comment row: %w", err)
		}
		comments = append(comments, c)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error after scanning comments: %w", err)
	}

	return comments, nil
}

func FetchPostByID(postID int) (*model.Post, error) {
	query := `
        SELECT p.id, u.username, p.title, p.content, p.user_id, p.categories, 
               p.created_at, p.updated_at,
			   COALESCE(SUM(CASE WHEN v.vote = 1 THEN 1 ELSE 0 END), 0) AS upvotes,
               COALESCE(SUM(CASE WHEN v.vote = -1 THEN 1 ELSE 0 END), 0) AS downvotes
        FROM posts p
        JOIN users u ON p.user_id = u.id
		LEFT JOIN votes v ON p.id = v.post_id
        WHERE p.id = ?
		GROUP BY p.id
    `

	var post model.Post
	err := DB.QueryRow(query, postID).Scan(
		&post.ID,
		&post.Author,
		&post.Title,
		&post.Content,
		&post.UserID,
		&post.Categories,
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.Upvotes,
		&post.Downvotes,
	)
	if err != nil {
		return nil, err
	}

	// Fetch comments for this post
	comments, err := FetchCommentsByPostID(postID)
	if err != nil {
		return nil, fmt.Errorf("error fetching comments: %w", err)
	}

	post.Comments = comments
	return &post, nil
}

// Add this function to fetch user data by ID
func FetchUserById(userID int) (*model.User, error) {
	var user model.User
	err := DB.QueryRow(
		"SELECT id, username, email, session_token, session_expiry, created_at FROM users WHERE id = ?", userID).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.SessionToken,
		&user.SessionExpiry,
		&user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func FetchPostsByUserID(userID int) ([]*model.Post, error) {
	query := `
        SELECT p.id, u.username, p.title, p.content, p.categories, 
               p.created_at, p.updated_at,
               COALESCE(SUM(CASE WHEN v.vote = 1 THEN 1 ELSE 0 END), 0) AS upvotes,
               COALESCE(SUM(CASE WHEN v.vote = -1 THEN 1 ELSE 0 END), 0) AS downvotes
        FROM posts p
        JOIN users u ON p.user_id = u.id
        LEFT JOIN votes v ON p.id = v.post_id
        WHERE p.user_id = ?
        GROUP BY p.id
        ORDER BY p.created_at DESC
    `

	rows, err := DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("error querying user posts: %w", err)
	}
	defer rows.Close()

	var posts []*model.Post
	for rows.Next() {
		var post model.Post
		err := rows.Scan(
			&post.ID,
			&post.Author,
			&post.Title,
			&post.Content,
			&post.Categories,
			&post.CreatedAt,
			&post.UpdatedAt,
			&post.Upvotes,
			&post.Downvotes,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning post: %w", err)
		}
		posts = append(posts, &post)
	}

	return posts, nil
}

func FetchLikedPostsByUserID(userID int) ([]*model.Post, error) {
	query := `
        SELECT p.id, u.username, p.title, p.content, p.categories, 
               p.created_at, p.updated_at,
               COALESCE(SUM(CASE WHEN v.vote = 1 THEN 1 ELSE 0 END), 0) AS upvotes,
               COALESCE(SUM(CASE WHEN v.vote = -1 THEN 1 ELSE 0 END), 0) AS downvotes
        FROM posts p
        JOIN votes v ON p.id = v.post_id
        JOIN users u ON p.user_id = u.id
        WHERE v.user_id = ? AND v.vote = 1
        GROUP BY p.id
        ORDER BY p.created_at DESC
    `

	rows, err := DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("error querying liked posts: %w", err)
	}
	defer rows.Close()

	var likedPosts []*model.Post
	for rows.Next() {
		var p model.Post
		err := rows.Scan(
			&p.ID,
			&p.Author,
			&p.Title,
			&p.Content,
			&p.Categories,
			&p.CreatedAt,
			&p.UpdatedAt,
			&p.Upvotes,
			&p.Downvotes,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning liked posts: %w", err)
		}
		likedPosts = append(likedPosts, &p)
	}

	return likedPosts, nil
}

func FetchDislikedPostsByUserID(userID int) ([]*model.Post, error) {
	query := `
        SELECT p.id, u.username, p.title, p.content, p.categories, 
               p.created_at, p.updated_at,
               COALESCE(SUM(CASE WHEN v.vote = 1 THEN 1 ELSE 0 END), 0) AS upvotes,
               COALESCE(SUM(CASE WHEN v.vote = -1 THEN 1 ELSE 0 END), 0) AS downvotes
        FROM posts p
        JOIN votes v ON p.id = v.post_id
        JOIN users u ON p.user_id = u.id
        WHERE v.user_id = ? AND v.vote = -1
        GROUP BY p.id
        ORDER BY p.created_at DESC
    `

	rows, err := DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("error querying disliked posts: %w", err)
	}
	defer rows.Close()

	var dislikedPosts []*model.Post
	for rows.Next() {
		var p model.Post
		err := rows.Scan(
			&p.ID,
			&p.Author,
			&p.Title,
			&p.Content,
			&p.Categories,
			&p.CreatedAt,
			&p.UpdatedAt,
			&p.Upvotes,
			&p.Downvotes,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning disliked posts: %w", err)
		}
		dislikedPosts = append(dislikedPosts, &p)
	}

	return dislikedPosts, nil
}

func FetchPostsByCategory(category string) ([]model.Post, error) {
	query := `
        SELECT p.id, u.username, p.title, p.content, p.user_id,
               p.categories, p.created_at, p.updated_at,
               COALESCE(SUM(CASE WHEN v.vote = 1 THEN 1 ELSE 0 END), 0) AS upvotes,
               COALESCE(SUM(CASE WHEN v.vote = -1 THEN 1 ELSE 0 END), 0) AS downvotes
        FROM posts p
        JOIN users u ON p.user_id = u.id
        LEFT JOIN votes v ON p.id = v.post_id
        WHERE p.categories LIKE ?
        GROUP BY p.id
        ORDER BY p.created_at DESC
    `

	rows, err := DB.Query(query, "%"+category+"%")
	if err != nil {
		return nil, fmt.Errorf("error querying posts by category: %w", err)
	}
	defer rows.Close()

	var posts []model.Post
	for rows.Next() {
		var p model.Post
		err := rows.Scan(
			&p.ID,
			&p.Author,
			&p.Title,
			&p.Content,
			&p.UserID,
			&p.Categories,
			&p.CreatedAt,
			&p.UpdatedAt,
			&p.Upvotes,
			&p.Downvotes,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning post: %w", err)
		}
		posts = append(posts, p)
	}

	return posts, nil
}

func UpdateVote(userID, postID, voteValue int) error {
	// Check if the user has already voted on this post
	var existingVote int
	err := DB.QueryRow("SELECT vote FROM votes WHERE user_id = ? AND post_id = ?", userID, postID).Scan(&existingVote)

	if err == sql.ErrNoRows {
		// User hasn't voted yet → Insert new vote
		_, err = DB.Exec("INSERT INTO votes (user_id, post_id, vote) VALUES (?, ?, ?)", userID, postID, voteValue)
		return err
	} else if err != nil {
		return err // Unexpected database error
	}

	// If the vote is different, update it
	if existingVote != voteValue {
		_, err = DB.Exec("UPDATE votes SET vote = ? WHERE user_id = ? AND post_id = ?", voteValue, userID, postID)
	}
	return err
}
