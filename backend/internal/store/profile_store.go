package store

import (
	"database/sql"

	"github.com/tajjjjr/social-network/backend/internal/models"
)

// ProfileStore handles database operations for profile.
type ProfileStore struct {
	DB *sql.DB
}

// NewProfileStore creates a new ProfileStore.
func NewProfileStore(db *sql.DB) *ProfileStore {
	return &ProfileStore{DB: db}
}

func (ps *ProfileStore) MyProfileDetails(userid int64) (models.ProfileDetails, error) {
	// Initialize an empty ProfileDetails struct
	var profile models.ProfileDetails
	// Prepare the SQL query to fetch user details
	var firstName, lastName, email, nickname, aboutMe, avatar sql.NullString
	var dateOfBirth sql.NullTime
	var isProfilePublic int

	query := `SELECT firstname, lastname, email, nickname, aboutme, dateofbirth, is_profile_public, avatar 
			  FROM Users 
			  WHERE id = ?`

	err := ps.DB.QueryRow(query, userid).Scan(
		&firstName,
		&lastName,
		&email,
		&nickname,
		&aboutMe,
		&dateOfBirth,
		&isProfilePublic,
		&avatar,
	)
	if err != nil {
		return profile, err
	}

	// Convert the retrieved data to strings, handling null values
	profile.FirstName = getStringValue(firstName)
	profile.LastName = getStringValue(lastName)
	profile.Email = getStringValue(email)
	profile.Nickname = getStringValue(nickname)
	profile.About = getStringValue(aboutMe)
	profile.DateOfBirth = dateOfBirth.Time.Format("2006-01-02")
	profile.ProfilePublic = isProfilePublic == 1
	profile.Avatar = getStringValue(avatar)

	return profile, nil
}

func (followStore *ProfileStore) GetFollowersStat(userid int64) (int, int, error) {
	followers := 0
	following := 0
	err := followStore.DB.QueryRow("SELECT COUNT(*) FROM Followers WHERE follower_id = ? AND status = 'accepted'", userid).Scan(&following)
	if err != nil {
		return 0, 0, err
	}
	err = followStore.DB.QueryRow("SELECT COUNT(*) FROM Followers WHERE followee_id = ? AND status = 'accepted'", userid).Scan(&followers)
	if err != nil {
		return 0, 0, err
	}
	return followers, following, nil
}

func (ps *ProfileStore) GetNumberOfPosts(userid int64) (int, error) {
	var count int
	err := ps.DB.QueryRow("SELECT COUNT(*) FROM Posts WHERE user_id = ?", userid).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (ps *ProfileStore) GetFollowStatus(user_id, LoggedInUser int64) (string, error) {
	var status string
	query := `SELECT status FROM Followers WHERE follower_id = ? AND followee_id = ?`
	err := ps.DB.QueryRow(query, LoggedInUser, user_id).Scan(&status)
	if err != nil {
		if err == sql.ErrNoRows {
			return "follow", nil // No follow relationship found
		}
		return "", err // Other error
	}

	if status == "accepted" {
		return "following", nil // Already following
	} else if status == "pending" {
		return "pending", nil // Follow request is pending
	}
	return "follow", nil // Follow request was rejected, can follow again
}

func (s *ProfileStore) GetPostsOfUser(id int64) ([]models.Post, error) {
	query := `
		SELECT p.id, p.user_id, p.content, p.image, p.privacy, p.created_at, p.updated_at,
			   u.firstname, u.lastname, u.nickname, u.avatar
		FROM Posts p
		JOIN Users u ON p.user_id = u.id
		WHERE p.user_id = ?`

	rows, err := s.DB.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		var updatedAt sql.NullTime
		var firstName, lastName, nickname, avatar sql.NullString

		if err := rows.Scan(
			&post.ID, &post.UserID, &post.Content, &post.Image, &post.Privacy,
			&post.CreatedAt, &updatedAt,
			&firstName, &lastName, &nickname, &avatar); err != nil {
			return nil, err
		}

		// Set the updated_at field and is_edited flag
		if updatedAt.Valid {
			post.UpdatedAt = &updatedAt.Time
			post.IsEdited = true
		}

		// Set user information
		post.Author = models.User{
			FirstName: &firstName.String,
			LastName:  &lastName.String,
			Nickname:  &nickname.String,
			Avatar:    &avatar.String,
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (followstore *ProfileStore) GetUserFollowers(userid int64) (models.FollowListResponse, error) {
	var followersList models.FollowListResponse
	rows, err := followstore.DB.Query(`
		SELECT u.id, u.firstname, u.lastname, u.avatar 
		FROM Users u 
		INNER JOIN Followers f ON u.id = f.follower_id 
		WHERE f.followee_id = ? AND f.status = 'accepted'`, userid)
	if err != nil {
		return followersList, err
	}
	defer rows.Close()

	for rows.Next() {
		var follower models.FollowUser
		var firstName, lastName, avatar sql.NullString
		err := rows.Scan(&follower.FollowerID, &firstName, &lastName, &avatar)
		if err != nil {
			return followersList, err
		}

		if firstName.Valid {
			follower.FirstName = firstName.String
		}
		if lastName.Valid {
			follower.LastName = lastName.String
		}
		if avatar.Valid {
			follower.Avatar = avatar.String
		}

		followersList.Followers = append(followersList.Followers, follower)
	}

	if err = rows.Err(); err != nil {
		return followersList, err
	}

	return followersList, nil
}

func (followstore *ProfileStore) GetUserFollowees(userid int64) (models.FollowListResponse, error) {
	var followersList models.FollowListResponse
	rows, err := followstore.DB.Query(`
		SELECT u.id, u.firstname, u.lastname, u.avatar 
		FROM Users u 
		INNER JOIN Followers f ON u.id = f.followee_id 
		WHERE f.follower_id = ? AND f.status = 'accepted'`, userid)
	if err != nil {
		return followersList, err
	}
	defer rows.Close()

	for rows.Next() {
		var follower models.FollowUser
		var firstName, lastName, avatar sql.NullString
		err := rows.Scan(&follower.FollowerID, &firstName, &lastName, &avatar)
		if err != nil {
			return followersList, err
		}

		if firstName.Valid {
			follower.FirstName = firstName.String
		}
		if lastName.Valid {
			follower.LastName = lastName.String
		}
		if avatar.Valid {
			follower.Avatar = avatar.String
		}

		followersList.Followers = append(followersList.Followers, follower)
	}

	if err = rows.Err(); err != nil {
		return followersList, err
	}

	return followersList, nil
}

func (pr *ProfileStore) GetUserPostPhotos(userId int64) ([]models.Photo, error) {
	var photos []models.Photo
	rows, err := pr.DB.Query(`
		SELECT p.image
		FROM Posts p
		WHERE p.user_id = ?`, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var photo models.Photo
		var image sql.NullString
		err := rows.Scan(&image)
		if err != nil {
			return nil, err
		}
		if image.Valid {
			photo.Image = image.String
		}
		photos = append(photos, photo)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return photos, nil
}

func (pr *ProfileStore) GetUserCommentPhotos(userId int64) ([]models.Photo, error) {
	var photos []models.Photo
	rows, err := pr.DB.Query(`
		SELECT c.image
		FROM Comments c
		WHERE c.user_id = ?`, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var photo models.Photo
		var image sql.NullString
		err := rows.Scan(&image)
		if err != nil {
			return nil, err
		}
		if image.Valid {
			photo.Image = image.String
		}
		photos = append(photos, photo)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return photos, nil
}

// Helper function to handle null string values
func getStringValue(s sql.NullString) string {
	if s.Valid {
		return s.String
	}
	return ""
}
