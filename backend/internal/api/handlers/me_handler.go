package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/tajjjjr/social-network/backend/internal/models"

	"fmt"
)

func NewMeHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, err := r.Cookie("session_id")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		cookieValue := strings.TrimPrefix(sessionCookie.Value, "session_id=")

		fmt.Println("Session Cookie:", cookieValue)
		var userID int

		errSession := db.QueryRow("SELECT user_id FROM sessions WHERE id = ?", cookieValue).Scan(&userID)
		if errSession != nil {
			fmt.Println("Error retrieving user ID from session:", errSession)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		fmt.Println("User ID from session:", userID)

		var user models.User
		errUser := db.QueryRow(
			"SELECT id, email, avatar, firstname, lastname, dateofbirth, nickname, aboutme, is_profile_public, created_at FROM Users WHERE id = ?",
			userID,
		).Scan(&user.ID, &user.Email, &user.Avatar, &user.FirstName, &user.LastName, &user.DateOfBirth, &user.Nickname, &user.AboutMe, &user.IsProfilePublic, &user.CreatedAt)
		if errUser != nil {
			fmt.Println("Error retrieving user:", errUser)
			_ = json.NewEncoder(w).Encode(map[string]string{"message": "This is the /me endpoint", "error": "User not found"})
			return
		}

		err = json.NewEncoder(w).Encode(user)
		if err != nil {
			fmt.Println("Error encoding json to the client side from new me handler")
		}
	}
}
