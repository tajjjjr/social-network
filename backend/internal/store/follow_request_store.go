package store

import (
	"database/sql"
	"strings"
	"time"

	"github.com/tajjjjr/social-network/backend/internal/models"
)

// FollowRequestStore handles database operations for follow request.
type FollowRequestStore struct {
	DB *sql.DB
}

// NewFollowRequestStore creates a new FollowRequestStore.
func NewFollowRequestStore(db *sql.DB) *FollowRequestStore {
	return &FollowRequestStore{DB: db}
}

// AcceptFollowConnection accepts a follow request by setting status to 'accepted' and updating accepted_at
func (fr *FollowRequestStore) AcceptFollowConnection(followConnectionID int64) error {
	query := "UPDATE Followers SET status = 'accepted', accepted_at = ? WHERE id = ? AND status = 'pending'"
	_, err := fr.DB.Exec(query, time.Now(), followConnectionID)
	return err
}

// RejectFollowConnection rejects a follow request by setting status to 'rejected'
func (fr *FollowRequestStore) RejectFollowConnection(followConnectionID int64) error {
	query := "UPDATE Followers SET status = 'rejected' WHERE id = ? AND status = 'pending'"
	_, err := fr.DB.Exec(query, followConnectionID)
	return err
}

// GetPendingFollowRequests retrieves all pending follow requests for a user
func (fr *FollowRequestStore) GetPendingFollowRequests(followeeID int64) ([]int64, error) {
	query := "SELECT id FROM Followers WHERE followee_id = ? AND status = 'pending'"
	rows, err := fr.DB.Query(query, followeeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requestIDs []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		requestIDs = append(requestIDs, id)
	}

	return requestIDs, rows.Err()
}

func (fr *FollowRequestStore) UserInfo(userID int64) (string, string, error) {
	var firstName, lastName sql.NullString
	query := "SELECT firstname, lastname FROM Users WHERE id = ?"
	err := fr.DB.QueryRow(query, userID).Scan(&firstName, &lastName)
	if err != nil {
		return "", "", err
	}

	name := ""
	if firstName.Valid && lastName.Valid {
		name = firstName.String + " " + lastName.String
	} else if firstName.Valid {
		name = firstName.String
	} else if lastName.Valid {
		name = lastName.String
	} else {
		name = "User"
	}

	return name, name, nil
}

func (fr *FollowRequestStore) RetrieveRequestInfo(requestID int64) (int64, int64, error) {
	var followerID, followeeID int64
	query := "SELECT follower_id, followee_id FROM Followers WHERE id = ?"
	err := fr.DB.QueryRow(query, requestID).Scan(&followerID, &followeeID)
	return followerID, followeeID, err
}

func (followstore *FollowRequestStore) AddtoNotification(follower_id int64, message string) error {
	// Determine notification type based on message content
	notificationType := "follow_request"
	if strings.Contains(message, "accepted") {
		notificationType = "follow_request_accepted"
	} else if strings.Contains(message, "rejected") {
		notificationType = "follow_request_rejected"
	}

	_, err := followstore.DB.Exec(`INSERT INTO Notifications (user_id, type, message) VALUES (?, ?, ?)`,
		follower_id, notificationType, message)
	return err
}

func (fr *FollowRequestStore) FollowRequestCancel(requestID int64) error {
	_, err := fr.DB.Exec("DELETE FROM Followers WHERE id = ? AND status = 'pending'", requestID)
	return err
}

// GetRequestIDByUsers retrieves the request ID for a pending follow request between two users
func (fr *FollowRequestStore) GetRequestIDByUsers(followerID, followeeID int64) (int64, error) {
	var requestID int64
	query := "SELECT id FROM Followers WHERE follower_id = ? AND followee_id = ? AND status = 'pending'"
	err := fr.DB.QueryRow(query, followerID, followeeID).Scan(&requestID)
	return requestID, err
}

func (fr *FollowRequestStore) GetPendingFollowRequest(userid int64) (models.FollowRequestUserResponse, error) {
	var response models.FollowRequestUserResponse

	query := `
        SELECT f.id, f.follower_id, u.firstname, u.lastname, u.avatar, f.requested_at
        FROM Followers f
        JOIN Users u ON f.follower_id = u.id
        WHERE f.followee_id = ? AND f.status = 'pending'
        ORDER BY f.requested_at DESC
    `

	rows, err := fr.DB.Query(query, userid)
	if err != nil {
		return response, err
	}
	defer rows.Close()

	var requests []models.FollowUserRequest
	for rows.Next() {
		var request models.FollowUserRequest
		if err := rows.Scan(
			&request.RequestID,
			&request.FollowerID,
			&request.FirstName,
			&request.LastName,
			&request.Avatar,
			&request.RequestedAt,
		); err != nil {
			return response, err
		}
		requests = append(requests, request)
	}

	if err := rows.Err(); err != nil {
		return response, err
	}

	response.Followers = requests
	return response, nil
}
