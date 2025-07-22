package handler

import (
	"database/sql"
	"fmt"
	"forum-go/database"
	"log"
	"net/http"
	"strconv"
)

/*
	func VoteHandler(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			ErrorHandler(w, r, http.StatusMethodNotAllowed)
			return
		}

		// Check if user is logged in
		isLoggedIn, userID := database.CheckUserLoggedIn(r)
		if !isLoggedIn {
			ErrorHandler(w, r, http.StatusUnauthorized)
			return
		}

		// Parse form values
		postID, err := strconv.Atoi(r.FormValue("post_id"))
		if err != nil {
			ErrorHandler(w, r, http.StatusBadRequest)
			return
		}

		voteValue, err := strconv.Atoi(r.FormValue("vote"))
		if err != nil || (voteValue != 1 && voteValue != -1) {
			ErrorHandler(w, r, http.StatusBadRequest)
			return
		}

		// Insert or update vote in database
		err = database.UpdateVote(userID, postID, voteValue)
		if err != nil {
			ErrorHandler(w, r, http.StatusInternalServerError)
			return
		}

		// Redirect back to the post
		http.Redirect(w, r, fmt.Sprintf("/viewpost?id=%d", postID), http.StatusSeeOther)
	}
*/
func VoteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	// Check if user is logged in
	isLoggedIn, userID := database.CheckUserLoggedIn(r)
	if !isLoggedIn {
		ErrorHandler(w, r, http.StatusUnauthorized)
		return
	}

	// Parse form values
	postID, err := strconv.Atoi(r.FormValue("post_id"))
	if err != nil {
		ErrorHandler(w, r, http.StatusBadRequest)
		return
	}

	voteValue, err := strconv.Atoi(r.FormValue("vote"))
	if err != nil || (voteValue != 1 && voteValue != -1) {
		ErrorHandler(w, r, http.StatusBadRequest)
		return
	}

	// Check if the user already voted on this post
	var existingVote int
	err = database.DB.QueryRow("SELECT vote FROM votes WHERE user_id = ? AND post_id = ?", userID, postID).Scan(&existingVote)

	if err == sql.ErrNoRows {
		// User hasn't voted yet → Insert new vote
		_, err = database.DB.Exec("INSERT INTO votes (user_id, post_id, vote) VALUES (?, ?, ?)", userID, postID, voteValue)
		if err != nil {
			ErrorHandler(w, r, http.StatusInternalServerError)
			return
		}
	} else if err == nil {
		if existingVote == voteValue {
			// User clicked the same button → Undo vote (delete it)
			_, err = database.DB.Exec("DELETE FROM votes WHERE user_id = ? AND post_id = ?", userID, postID)
			if err != nil {
				ErrorHandler(w, r, http.StatusInternalServerError)
				return
			}
		} else {
			// User changed their vote → Update it
			_, err = database.DB.Exec("UPDATE votes SET vote = ? WHERE user_id = ? AND post_id = ?", voteValue, userID, postID)
			if err != nil {
				ErrorHandler(w, r, http.StatusInternalServerError)
				return
			}
		}
	} else {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	// Redirect back to the post
	http.Redirect(w, r, fmt.Sprintf("/viewpost?id=%d", postID), http.StatusSeeOther)
}

func VoteCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Println("Invalid method:", r.Method)
		return
	}

	isLoggedIn, userID := database.CheckUserLoggedIn(r)
	if !isLoggedIn {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		log.Println("Unauthorized user tried to vote.")
		return
	}

	commentIDStr := r.FormValue("comment_id")
	voteValueStr := r.FormValue("vote")

	if commentIDStr == "" || voteValueStr == "" {
		http.Error(w, "Missing comment_id or vote", http.StatusBadRequest)
		log.Println("Missing comment_id or vote in request.")
		return
	}

	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		log.Println("Invalid comment ID:", commentIDStr)
		return
	}

	voteValue, err := strconv.Atoi(voteValueStr)
	if err != nil || (voteValue != 1 && voteValue != -1) {
		http.Error(w, "Invalid vote value", http.StatusBadRequest)
		log.Println("Invalid vote value:", voteValueStr)
		return
	}

	// Check if the user already voted on this comment
	var existingVote int
	err = database.DB.QueryRow("SELECT vote FROM votes WHERE user_id = ? AND comment_id = ?", userID, commentID).Scan(&existingVote)

	if err == sql.ErrNoRows {
		// User hasn't voted yet → Insert new vote
		_, err = database.DB.Exec("INSERT INTO votes (user_id, comment_id, vote) VALUES (?, ?, ?)", userID, commentID, voteValue)
		if err != nil {
			http.Error(w, "Error processing vote", http.StatusInternalServerError)
			log.Println("Database error:", err)
			return
		}
		//log.Printf("User %d voted %d on comment %d\n", userID, voteValue, commentID)

	} else if err == nil {
		if existingVote == voteValue {
			// User clicked the same button → Undo vote (delete it)
			_, err = database.DB.Exec("DELETE FROM votes WHERE user_id = ? AND comment_id = ?", userID, commentID)
			if err != nil {
				http.Error(w, "Error removing vote", http.StatusInternalServerError)
				log.Println("Database error (removing vote):", err)
				return
			}
			//log.Printf("User %d removed vote on comment %d\n", userID, commentID)

		} else {
			// User changed their vote → Update it
			_, err = database.DB.Exec("UPDATE votes SET vote = ? WHERE user_id = ? AND comment_id = ?", voteValue, userID, commentID)
			if err != nil {
				http.Error(w, "Error updating vote", http.StatusInternalServerError)
				log.Println("Database error (updating vote):", err)
				return
			}
			//log.Printf("User %d changed vote to %d on comment %d\n", userID, voteValue, commentID)
		}
	} else {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Println("Database error:", err)
		return
	}

	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}
