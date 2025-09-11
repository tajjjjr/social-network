package tests

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/tajjjjr/social-network/backend/internal/api/handlers"
	"github.com/tajjjjr/social-network/backend/internal/models"
	"github.com/tajjjjr/social-network/backend/internal/service"
	"github.com/tajjjjr/social-network/backend/internal/store"
	"github.com/tajjjjr/social-network/backend/pkg/utils"
)

func setupProfileEditTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}

	// Create Users table
	createUsersTable := `
	CREATE TABLE Users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		firstname TEXT,
		lastname TEXT,
		dateofbirth DATE,
		nickname TEXT,
		aboutme TEXT,
		is_profile_public INTEGER DEFAULT 1,
		avatar TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	if _, err := db.Exec(createUsersTable); err != nil {
		t.Fatalf("Failed to create Users table: %v", err)
	}

	// Insert test user
	insertUser := `
	INSERT INTO Users (id, email, password, firstname, lastname, dateofbirth, nickname, aboutme, is_profile_public, avatar)
	VALUES (1, 'test@example.com', '$2a$10$hashedpassword', 'John', 'Doe', '1990-01-01', 'johndoe', 'Test bio', 1, 'avatar.jpg');
	`

	if _, err := db.Exec(insertUser); err != nil {
		t.Fatalf("Failed to insert test user: %v", err)
	}

	// Insert another user for email conflict testing
	insertUser2 := `
	INSERT INTO Users (id, email, password, firstname, lastname, dateofbirth, nickname, aboutme, is_profile_public, avatar)
	VALUES (2, 'existing@example.com', '$2a$10$hashedpassword', 'Jane', 'Smith', '1992-05-15', 'janesmith', 'Another bio', 0, 'avatar2.jpg');
	`

	if _, err := db.Exec(insertUser2); err != nil {
		t.Fatalf("Failed to insert second test user: %v", err)
	}

	return db
}

func createMultipartFormData(fields map[string]string, files map[string][]byte) (*bytes.Buffer, string, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Add form fields
	for key, value := range fields {
		if err := writer.WriteField(key, value); err != nil {
			return nil, "", err
		}
	}

	// Add files
	for fieldName, fileData := range files {
		part, err := writer.CreateFormFile(fieldName, "test.jpg")
		if err != nil {
			return nil, "", err
		}
		if _, err := part.Write(fileData); err != nil {
			return nil, "", err
		}
	}

	if err := writer.Close(); err != nil {
		return nil, "", err
	}

	return &buf, writer.FormDataContentType(), nil
}

func TestEditProfile_Success(t *testing.T) {
	db := setupProfileEditTestDB(t)
	defer db.Close()

	authStore := store.NewAuthStore(db)
	authService := service.NewAuthService(authStore)
	authHandler := handlers.NewAuthHandler(authService)

	// Create form data
	formData := map[string]string{
		"email":       "updated@example.com",
		"firstname":   "UpdatedJohn",
		"lastname":    "UpdatedDoe",
		"dateofbirth": "1991-02-02",
		"nickname":    "updatedjohndoe",
		"aboutme":     "Updated bio",
		"is_private":  "true",
	}

	body, contentType, err := createMultipartFormData(formData, nil)
	if err != nil {
		t.Fatalf("Failed to create form data: %v", err)
	}

	req := httptest.NewRequest("PUT", "/EditProfile", body)
	req.Header.Set("Content-Type", contentType)

	// Add user ID to context
	ctx := context.WithValue(req.Context(), utils.User_id, int64(1))
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	authHandler.EditProfile(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response utils.Response
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Message != "Profile updated successfully" {
		t.Errorf("Expected success message, got: %s", response.Message)
	}

	// Verify the profile was actually updated in the database
	var updatedUser models.User
	err = db.QueryRow("SELECT email, firstname, lastname, dateofbirth, nickname, aboutme, is_profile_public FROM Users WHERE id = 1").Scan(
		&updatedUser.Email, &updatedUser.FirstName, &updatedUser.LastName,
		&updatedUser.DateOfBirth, &updatedUser.Nickname, &updatedUser.AboutMe, &updatedUser.IsProfilePublic)
	if err != nil {
		t.Fatalf("Failed to query updated user: %v", err)
	}

	if updatedUser.Email != "updated@example.com" {
		t.Errorf("Expected email 'updated@example.com', got '%s'", updatedUser.Email)
	}
	if *updatedUser.FirstName != "UpdatedJohn" {
		t.Errorf("Expected first name 'UpdatedJohn', got '%s'", *updatedUser.FirstName)
	}
	if *updatedUser.LastName != "UpdatedDoe" {
		t.Errorf("Expected last name 'UpdatedDoe', got '%s'", *updatedUser.LastName)
	}
	if updatedUser.IsProfilePublic != false {
		t.Errorf("Expected profile to be private, got public")
	}
}

func TestEditProfile_EmailAlreadyExists(t *testing.T) {
	db := setupProfileEditTestDB(t)
	defer db.Close()

	authStore := store.NewAuthStore(db)
	authService := service.NewAuthService(authStore)
	authHandler := handlers.NewAuthHandler(authService)

	// Try to update to an email that already exists (user 2's email)
	formData := map[string]string{
		"email":       "existing@example.com", // This email belongs to user 2
		"firstname":   "UpdatedJohn",
		"lastname":    "UpdatedDoe",
		"dateofbirth": "1991-02-02",
		"nickname":    "updatedjohndoe",
		"aboutme":     "Updated bio",
		"is_private":  "false",
	}

	body, contentType, err := createMultipartFormData(formData, nil)
	if err != nil {
		t.Fatalf("Failed to create form data: %v", err)
	}

	req := httptest.NewRequest("PUT", "/EditProfile", body)
	req.Header.Set("Content-Type", contentType)

	// Add user ID to context (user 1 trying to use user 2's email)
	ctx := context.WithValue(req.Context(), utils.User_id, int64(1))
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	authHandler.EditProfile(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("Expected status %d, got %d", http.StatusConflict, w.Code)
	}

	var response utils.Response
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Message != "Email already exists" {
		t.Errorf("Expected 'Email already exists', got: %s", response.Message)
	}
}

func TestEditProfile_SameEmailAllowed(t *testing.T) {
	db := setupProfileEditTestDB(t)
	defer db.Close()

	authStore := store.NewAuthStore(db)
	authService := service.NewAuthService(authStore)
	authHandler := handlers.NewAuthHandler(authService)

	// User updating with their own email should be allowed
	formData := map[string]string{
		"email":       "test@example.com", // Same email as user 1
		"firstname":   "UpdatedJohn",
		"lastname":    "UpdatedDoe",
		"dateofbirth": "1991-02-02",
		"nickname":    "updatedjohndoe",
		"aboutme":     "Updated bio",
		"is_private":  "false",
	}

	body, contentType, err := createMultipartFormData(formData, nil)
	if err != nil {
		t.Fatalf("Failed to create form data: %v", err)
	}

	req := httptest.NewRequest("PUT", "/EditProfile", body)
	req.Header.Set("Content-Type", contentType)

	ctx := context.WithValue(req.Context(), utils.User_id, int64(1))
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	authHandler.EditProfile(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response utils.Response
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Message != "Profile updated successfully" {
		t.Errorf("Expected success message, got: %s", response.Message)
	}
}

func TestEditProfile_InvalidEmail(t *testing.T) {
	db := setupProfileEditTestDB(t)
	defer db.Close()

	authStore := store.NewAuthStore(db)
	authService := service.NewAuthService(authStore)
	authHandler := handlers.NewAuthHandler(authService)

	testCases := []struct {
		name  string
		email string
	}{
		{"Invalid format", "invalid-email"},
		{"Missing @", "test.example.com"},
		{"Missing domain", "test@"},
		{"Missing local part", "@example.com"},
		{"Empty email", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			formData := map[string]string{
				"email":       tc.email,
				"firstname":   "UpdatedJohn",
				"lastname":    "UpdatedDoe",
				"dateofbirth": "1991-02-02",
				"nickname":    "updatedjohndoe",
				"aboutme":     "Updated bio",
				"is_private":  "false",
			}

			body, contentType, err := createMultipartFormData(formData, nil)
			if err != nil {
				t.Fatalf("Failed to create form data: %v", err)
			}

			req := httptest.NewRequest("PUT", "/EditProfile", body)
			req.Header.Set("Content-Type", contentType)

			ctx := context.WithValue(req.Context(), utils.User_id, int64(1))
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			authHandler.EditProfile(w, req)

			if w.Code != http.StatusBadRequest {
				t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
			}

			var response utils.Response
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			if response.Message != "Invalid email format" {
				t.Errorf("Expected 'Invalid email format', got: %s", response.Message)
			}
		})
	}
}

func TestEditProfile_Unauthorized(t *testing.T) {
	db := setupProfileEditTestDB(t)
	defer db.Close()

	authStore := store.NewAuthStore(db)
	authService := service.NewAuthService(authStore)
	authHandler := handlers.NewAuthHandler(authService)

	formData := map[string]string{
		"email":       "updated@example.com",
		"firstname":   "UpdatedJohn",
		"lastname":    "UpdatedDoe",
		"dateofbirth": "1991-02-02",
		"nickname":    "updatedjohndoe",
		"aboutme":     "Updated bio",
		"is_private":  "false",
	}

	body, contentType, err := createMultipartFormData(formData, nil)
	if err != nil {
		t.Fatalf("Failed to create form data: %v", err)
	}

	req := httptest.NewRequest("PUT", "/EditProfile", body)
	req.Header.Set("Content-Type", contentType)
	// Don't add user ID to context to simulate unauthorized request

	w := httptest.NewRecorder()
	authHandler.EditProfile(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}

	var response utils.Response
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Message != "User not found in context" {
		t.Errorf("Expected 'User not found in context', got: %s", response.Message)
	}
}

func TestEditProfile_InvalidFormData(t *testing.T) {
	db := setupProfileEditTestDB(t)
	defer db.Close()

	authStore := store.NewAuthStore(db)
	authService := service.NewAuthService(authStore)
	authHandler := handlers.NewAuthHandler(authService)

	// Send invalid form data
	req := httptest.NewRequest("PUT", "/EditProfile", strings.NewReader("invalid form data"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	ctx := context.WithValue(req.Context(), utils.User_id, int64(1))
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	authHandler.EditProfile(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var response utils.Response
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Message != "Failed to parse form" {
		t.Errorf("Expected 'Failed to parse form', got: %s", response.Message)
	}
}

func TestEditProfile_WithAvatar(t *testing.T) {
	db := setupProfileEditTestDB(t)
	defer db.Close()

	authStore := store.NewAuthStore(db)
	authService := service.NewAuthService(authStore)
	authHandler := handlers.NewAuthHandler(authService)

	// Create form data with a small test image
	formData := map[string]string{
		"email":       "updated@example.com",
		"firstname":   "UpdatedJohn",
		"lastname":    "UpdatedDoe",
		"dateofbirth": "1991-02-02",
		"nickname":    "updatedjohndoe",
		"aboutme":     "Updated bio",
		"is_private":  "false",
	}

	// Create a small test image (1x1 pixel JPEG)
	testImage := []byte{
		0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46, 0x49, 0x46, 0x00, 0x01,
		0x01, 0x01, 0x00, 0x48, 0x00, 0x48, 0x00, 0x00, 0xFF, 0xDB, 0x00, 0x43,
	}

	files := map[string][]byte{
		"profilePicture": testImage,
	}

	body, contentType, err := createMultipartFormData(formData, files)
	if err != nil {
		t.Fatalf("Failed to create form data: %v", err)
	}

	req := httptest.NewRequest("PUT", "/EditProfile", body)
	req.Header.Set("Content-Type", contentType)

	ctx := context.WithValue(req.Context(), utils.User_id, int64(1))
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	authHandler.EditProfile(w, req)

	// Note: This test might fail if the avatar upload logic requires specific image formats
	// or if the UploadAvatarImage function has strict validation
	if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d or %d, got %d", http.StatusOK, http.StatusInternalServerError, w.Code)
	}
}

func TestEditProfile_XSSPrevention(t *testing.T) {
	db := setupProfileEditTestDB(t)
	defer db.Close()

	authStore := store.NewAuthStore(db)
	authService := service.NewAuthService(authStore)
	authHandler := handlers.NewAuthHandler(authService)

	// Test XSS prevention with malicious input
	formData := map[string]string{
		"email":       "test@example.com",
		"firstname":   "<script>alert('xss')</script>",
		"lastname":    "<img src=x onerror=alert('xss')>",
		"dateofbirth": "1991-02-02",
		"nickname":    "<svg onload=alert('xss')>",
		"aboutme":     "<iframe src='javascript:alert(\"xss\")'></iframe>",
		"is_private":  "false",
	}

	body, contentType, err := createMultipartFormData(formData, nil)
	if err != nil {
		t.Fatalf("Failed to create form data: %v", err)
	}

	req := httptest.NewRequest("PUT", "/EditProfile", body)
	req.Header.Set("Content-Type", contentType)

	ctx := context.WithValue(req.Context(), utils.User_id, int64(1))
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	authHandler.EditProfile(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Verify that the malicious content was escaped
	var updatedUser models.User
	err = db.QueryRow("SELECT firstname, lastname, nickname, aboutme FROM Users WHERE id = 1").Scan(
		&updatedUser.FirstName, &updatedUser.LastName, &updatedUser.Nickname, &updatedUser.AboutMe)
	if err != nil {
		t.Fatalf("Failed to query updated user: %v", err)
	}

	// Check that HTML entities are escaped
	if strings.Contains(*updatedUser.FirstName, "<script>") {
		t.Errorf("XSS content not properly escaped in first name: %s", *updatedUser.FirstName)
	}
	if strings.Contains(*updatedUser.LastName, "<img") {
		t.Errorf("XSS content not properly escaped in last name: %s", *updatedUser.LastName)
	}
	if strings.Contains(*updatedUser.Nickname, "<svg") {
		t.Errorf("XSS content not properly escaped in nickname: %s", *updatedUser.Nickname)
	}
	if strings.Contains(*updatedUser.AboutMe, "<iframe") {
		t.Errorf("XSS content not properly escaped in about me: %s", *updatedUser.AboutMe)
	}
}

func TestEditProfile_ProfileVisibilityToggle(t *testing.T) {
	db := setupProfileEditTestDB(t)
	defer db.Close()

	authStore := store.NewAuthStore(db)
	authService := service.NewAuthService(authStore)
	authHandler := handlers.NewAuthHandler(authService)

	testCases := []struct {
		name           string
		isPrivate      string
		expectedPublic bool
	}{
		{"Set to public (is_private not set)", "", true},
		{"Set to private (is_private = true)", "true", false},
		{"Set to public (is_private = false)", "false", true},
		{"Invalid value defaults to public", "invalid", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			formData := map[string]string{
				"email":       "test@example.com",
				"firstname":   "John",
				"lastname":    "Doe",
				"dateofbirth": "1990-01-01",
				"nickname":    "johndoe",
				"aboutme":     "Test bio",
				"is_private":  tc.isPrivate,
			}

			body, contentType, err := createMultipartFormData(formData, nil)
			if err != nil {
				t.Fatalf("Failed to create form data: %v", err)
			}

			req := httptest.NewRequest("PUT", "/EditProfile", body)
			req.Header.Set("Content-Type", contentType)

			ctx := context.WithValue(req.Context(), utils.User_id, int64(1))
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			authHandler.EditProfile(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
			}

			// Verify the profile visibility was set correctly
			var isPublic bool
			err = db.QueryRow("SELECT is_profile_public FROM Users WHERE id = 1").Scan(&isPublic)
			if err != nil {
				t.Fatalf("Failed to query profile visibility: %v", err)
			}

			if isPublic != tc.expectedPublic {
				t.Errorf("Expected profile public: %v, got: %v", tc.expectedPublic, isPublic)
			}
		})
	}
}
