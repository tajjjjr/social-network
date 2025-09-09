package store

import (
	"database/sql"
	"time"

	"github.com/tajjjjr/social-network/backend/internal/models"
)

// AuthStore handles database operations for authentication.
type AuthStore struct {
	DB *sql.DB
}

// NewAuthStore creates a new AuthStore.
func NewAuthStore(db *sql.DB) *AuthStore {
	return &AuthStore{DB: db}
}

// GetUserByEmail fetches a user by email
func (s *AuthStore) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := s.DB.QueryRow(
		"SELECT id, email, password, firstname, lastname, dateofbirth, avatar, nickname, aboutme, is_profile_public, created_at FROM Users WHERE email = ?",
		email,
	).Scan(&user.ID, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.DateOfBirth, &user.Avatar, &user.Nickname, &user.AboutMe, &user.IsProfilePublic, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// CreateSession creates a new session record and returns the session ID
/*
func (s *UserStore) CreateSession(userID int) (string, error) {
	sessionID := uuid.New().String()
	expiresAt := time.Now().Add(7 * 24 * time.Hour) // 7 days expiry

	_, err := s.db.Exec(
		"INSERT INTO Sessions (id, user_id, expires_at) VALUES (?, ?, ?)",
		sessionID, userID, expiresAt,
	)
	if err != nil {
		return "", err
	}
	return sessionID, nil
}
*/
// CreateSession creates a new session for a user.
func (s *AuthStore) CreateSession(userID int64, sessionID string, expiresAt time.Time) error {
	_, err := s.DB.Exec("INSERT INTO Sessions (user_id, id, created_at, expires_at) VALUES (?, ?, ?, ?)", userID, sessionID, time.Now(), expiresAt)
	return err
}

// -----------------------------------------------

// GetUserIDBySession returns the user ID associated with a valid session ID
func (s *AuthStore) GetUserIDBySession(sessionID string) (int, error) {
	var userID int
	var expiresAt time.Time

	err := s.DB.QueryRow(
		"SELECT user_id, expires_at FROM Sessions WHERE id = ?",
		sessionID,
	).Scan(&userID, &expiresAt)
	if err != nil {
		return 0, err
	}

	if expiresAt.Before(time.Now()) {
		return 0, sql.ErrNoRows // session expired
	}

	return userID, nil
}

// DeleteSession deletes a session by ID
func (s *AuthStore) DeleteSession(sessionID string) error {
	_, err := s.DB.Exec("DELETE FROM Sessions WHERE id = ?", sessionID)
	return err
}

// CreateUser creates a new user in the database
func (s *AuthStore) CreateUser(user *models.User) (int64, error) {
	stmt, err := s.DB.Prepare(`
		INSERT INTO Users (email, password, firstname, lastname, dateofbirth, nickname, aboutme, is_profile_public, avatar, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(
		user.Email,
		user.Password,
		user.FirstName,
		user.LastName,
		user.DateOfBirth,
		user.Nickname,
		user.AboutMe,
		user.IsProfilePublic,
		user.Avatar,
		user.CreatedAt,
	)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// UserExists checks if a user with the given email already exists
func (s *AuthStore) UserExists(email string) (bool, error) {
	var count int
	err := s.DB.QueryRow("SELECT COUNT(*) FROM Users WHERE email = ?", email).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *AuthStore) NewEditEmailExist(email string, userid int64) (bool, error) {
	var count int
	err := s.DB.QueryRow("SELECT COUNT(*) FROM Users WHERE email = ? and id != ?", email, userid).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func GetUserIDFromSession(sessionID string, db *sql.DB) (int64, error) {
	var userID int64
	err := db.QueryRow("SELECT user_id FROM Sessions WHERE id = ?", sessionID).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func (s *AuthStore) EditProfile(user *models.User, userid int64) error {
	var err error
	if *user.Avatar != "no profile photo" {
		_, err = s.DB.Exec("UPDATE Users SET email = ?, firstname = ?, lastname = ?, dateofbirth = ?, nickname = ?, aboutme = ?, is_profile_public = ?, avatar = ? WHERE id = ?", user.Email, user.FirstName, user.LastName, user.DateOfBirth, user.Nickname, user.AboutMe, user.IsProfilePublic, user.Avatar, userid)
	} else {
		_, err = s.DB.Exec("UPDATE Users SET email = ?, firstname = ?, lastname = ?, dateofbirth = ?, nickname = ?, aboutme = ?, is_profile_public = ? WHERE id = ?", user.Email, user.FirstName, user.LastName, user.DateOfBirth, user.Nickname, user.AboutMe, user.IsProfilePublic, userid)
	}
	return err
}

// EditProfile: already matches schema
// CreateUser: already matches schema
// GetUserByEmail: already matches schema
// GetUserIDBySession: already matches schema
// NewEditEmailExist: already matches schema
// UserExists: already matches schema
// All queries in this file now match the Users table schema and model fields.
