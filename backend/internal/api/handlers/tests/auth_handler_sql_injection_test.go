package tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/tajjjjr/social-network/backend/internal/api/handlers"
	"github.com/tajjjjr/social-network/backend/internal/service"
	"github.com/tajjjjr/social-network/backend/internal/store"
	"github.com/tajjjjr/social-network/backend/pkg/utils"
)

// setupLoginTestDB creates an in-memory SQLite database for login testing
func setupLoginTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	// Create Users table
	_, err = db.Exec(`
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
		)
	`)
	if err != nil {
		t.Fatalf("Failed to create Users table: %v", err)
	}

	// Create Sessions table
	_, err = db.Exec(`
		CREATE TABLE Sessions (
			id TEXT PRIMARY KEY,
			user_id INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			expires_at DATETIME NOT NULL,
			FOREIGN KEY (user_id) REFERENCES Users(id)
		)
	`)
	if err != nil {
		t.Fatalf("Failed to create Sessions table: %v", err)
	}

	return db
}

// createLoginTestUser creates a test user in the database
func createLoginTestUser(t *testing.T, db *sql.DB) {
	passwordManager := utils.NewPasswordManager(utils.PasswordConfig{})
	hashedPassword, err := passwordManager.HashPassword("TestPassword123!")
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	_, err = db.Exec(`
		INSERT INTO Users (email, password, firstname, lastname, nickname, is_profile_public)
		VALUES (?, ?, ?, ?, ?, ?)
	`, "test@example.com", hashedPassword, "Test", "User", "testuser", 1)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
}

// SQL injection payloads to test against login
var loginSQLInjectionPayloads = []string{
	// Basic SQL injection attempts
	"' OR '1'='1",
	"' OR 1=1--",
	"' OR 1=1#",
	"admin'--",
	"admin'#",
	"' or 1=1#",
	"' or 1=1--",
	"') or '1'='1--",
	"') or ('1'='1--",

	// Union-based injection
	"' UNION SELECT 1,2,3--",
	"' UNION SELECT null,null,null--",
	"' UNION SELECT email,password,id FROM Users--",
	"test@example.com' UNION SELECT 1,'admin','password123'--",

	// Boolean-based blind injection
	"test@example.com' AND 1=1--",
	"test@example.com' AND 1=2--",
	"test@example.com' AND (SELECT COUNT(*) FROM Users)>0--",

	// Stacked queries (dangerous)
	"test@example.com'; DROP TABLE Users;--",
	"test@example.com'; INSERT INTO Users (email,password) VALUES ('hacker','hacked');--",
	"test@example.com'; UPDATE Users SET password='hacked' WHERE email='test@example.com';--",

	// Special characters and escape attempts
	"test@example.com\\",
	"test@example.com\"",
	"test@example.com`",
	"test@example.com'",
	"test@example.com';",
	"test@example.com')",
	"test@example.com'\"",
}

func TestLoginSQLInjection_JSON(t *testing.T) {
	db := setupLoginTestDB(t)
	defer db.Close()

	createLoginTestUser(t, db)

	authStore := store.NewAuthStore(db)
	authService := service.NewAuthService(authStore)
	authHandler := handlers.NewAuthHandler(authService)

	for i, payload := range loginSQLInjectionPayloads {
		t.Run(fmt.Sprintf("JSON_Payload_%d_%s", i, strings.ReplaceAll(payload, "'", "QUOTE")), func(t *testing.T) {
			// Test with malicious email
			loginData := map[string]string{
				"email":    payload,
				"password": "TestPassword123!",
			}

			jsonData, _ := json.Marshal(loginData)
			req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			authHandler.Login(w, req)

			// Should not return 200 (successful login) for injection attempts
			if w.Code == http.StatusOK {
				t.Errorf("SQL injection payload succeeded: %s (Response: %s)", payload, w.Body.String())
			}

			// Test with malicious password
			loginData = map[string]string{
				"email":    "test@example.com",
				"password": payload,
			}

			jsonData, _ = json.Marshal(loginData)
			req = httptest.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			w = httptest.NewRecorder()
			authHandler.Login(w, req)

			// Should not return 200 (successful login) for injection attempts
			if w.Code == http.StatusOK {
				t.Errorf("SQL injection payload succeeded in password field: %s (Response: %s)", payload, w.Body.String())
			}
		})
	}
}

func TestLoginSQLInjection_FormData(t *testing.T) {
	db := setupLoginTestDB(t)
	defer db.Close()

	createLoginTestUser(t, db)

	authStore := store.NewAuthStore(db)
	authService := service.NewAuthService(authStore)
	authHandler := handlers.NewAuthHandler(authService)

	for i, payload := range loginSQLInjectionPayloads {
		t.Run(fmt.Sprintf("Form_Payload_%d_%s", i, strings.ReplaceAll(payload, "'", "QUOTE")), func(t *testing.T) {
			// Test with malicious email
			formData := url.Values{}
			formData.Set("email", payload)
			formData.Set("password", "TestPassword123!")

			req := httptest.NewRequest("POST", "/login", strings.NewReader(formData.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			w := httptest.NewRecorder()
			authHandler.Login(w, req)

			// Should not return 200 (successful login) for injection attempts
			if w.Code == http.StatusOK {
				t.Errorf("SQL injection payload succeeded: %s (Response: %s)", payload, w.Body.String())
			}

			// Test with malicious password
			formData = url.Values{}
			formData.Set("email", "test@example.com")
			formData.Set("password", payload)

			req = httptest.NewRequest("POST", "/login", strings.NewReader(formData.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			w = httptest.NewRecorder()
			authHandler.Login(w, req)

			// Should not return 200 (successful login) for injection attempts
			if w.Code == http.StatusOK {
				t.Errorf("SQL injection payload succeeded in password field: %s (Response: %s)", payload, w.Body.String())
			}
		})
	}
}

func TestLoginValidCredentials(t *testing.T) {
	db := setupLoginTestDB(t)
	defer db.Close()

	createLoginTestUser(t, db)

	authStore := store.NewAuthStore(db)
	authService := service.NewAuthService(authStore)
	authHandler := handlers.NewAuthHandler(authService)

	// Test valid login
	loginData := map[string]string{
		"email":    "test@example.com",
		"password": "TestPassword123!",
	}

	jsonData, _ := json.Marshal(loginData)
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	authHandler.Login(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Valid login failed. Expected 200, got %d. Response: %s", w.Code, w.Body.String())
	}

	// Check if session cookie is set
	cookies := w.Result().Cookies()
	sessionCookieFound := false
	for _, cookie := range cookies {
		if cookie.Name == "session_id" && cookie.Value != "" {
			sessionCookieFound = true
			break
		}
	}

	if !sessionCookieFound {
		t.Error("Session cookie not set after valid login")
	}
}

func TestLoginDatabaseIntegrityAfterInjectionAttempts(t *testing.T) {
	db := setupLoginTestDB(t)
	defer db.Close()

	createLoginTestUser(t, db)

	authStore := store.NewAuthStore(db)
	authService := service.NewAuthService(authStore)
	authHandler := handlers.NewAuthHandler(authService)

	// Count users before injection attempts
	var userCountBefore int
	err := db.QueryRow("SELECT COUNT(*) FROM Users").Scan(&userCountBefore)
	if err != nil {
		t.Fatalf("Failed to count users before test: %v", err)
	}

	// Try some dangerous injection payloads
	dangerousPayloads := []string{
		"'; DROP TABLE Users;--",
		"'; DELETE FROM Users;--",
		"'; INSERT INTO Users (email,password) VALUES ('hacker','hacked');--",
	}

	for _, payload := range dangerousPayloads {
		loginData := map[string]string{
			"email":    payload,
			"password": "TestPassword123!",
		}

		jsonData, _ := json.Marshal(loginData)
		req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		authHandler.Login(w, req)

		// Should not succeed
		if w.Code == http.StatusOK {
			t.Errorf("Dangerous SQL injection payload succeeded: %s", payload)
		}
	}

	// Count users after injection attempts
	var userCountAfter int
	err = db.QueryRow("SELECT COUNT(*) FROM Users").Scan(&userCountAfter)
	if err != nil {
		t.Fatalf("Failed to count users after test: %v", err)
	}

	if userCountBefore != userCountAfter {
		t.Errorf("User count changed after injection attempts. Before: %d, After: %d", userCountBefore, userCountAfter)
	}

	// Verify the original test user still exists and can login
	loginData := map[string]string{
		"email":    "test@example.com",
		"password": "TestPassword123!",
	}

	jsonData, _ := json.Marshal(loginData)
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	authHandler.Login(w, req)

	if w.Code != http.StatusOK {
		t.Error("Original test user can no longer login after injection attempts")
	}
}

func setupSignupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	_, err = db.Exec(`
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
		)
	`)
	if err != nil {
		t.Fatalf("Failed to create Users table: %v", err)
	}

	return db
}

// Helper to create multipart form body for signup
func createSignupMultipartForm(fields map[string]string) (*bytes.Buffer, string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for key, val := range fields {
		_ = writer.WriteField(key, val)
	}
	err := writer.Close()
	return body, writer.FormDataContentType(), err
}

func TestSignupHandler_SQLInjectionAttempt(t *testing.T) {
	db := setupSignupTestDB(t)
	defer db.Close()

	authStore := store.NewAuthStore(db)
	authService := service.NewAuthService(authStore)
	authHandler := handlers.NewAuthHandler(authService)

	fields := map[string]string{
		"email":             "' OR 1=1;--",
		"password":          "Password1!",
		"firstName":         "John",
		"lastName":          "Doe",
		"dob":               "2000-01-01",
		"nickname":          "hacker",
		"aboutMe":           "<script>alert(1)</script>",
		"profileVisibility": "public",
	}
	body, contentType, _ := createSignupMultipartForm(fields)
	req := httptest.NewRequest("POST", "/signup", body)
	req.Header.Set("Content-Type", contentType)
	w := httptest.NewRecorder()
	authHandler.Signup(w, req)
	resp := w.Result()
	if resp.StatusCode == http.StatusOK {
		t.Error("SQL injection attempt should not succeed")
	}
}

func TestSignupHandler_XSSPrevention(t *testing.T) {
	db := setupSignupTestDB(t)
	defer db.Close()

	authStore := store.NewAuthStore(db)
	authService := service.NewAuthService(authStore)
	authHandler := handlers.NewAuthHandler(authService)

	fields := map[string]string{
		"email":             "test@example.com",
		"password":          "Password1!",
		"firstName":         "<script>alert('xss')</script>",
		"lastName":          "<img src=x onerror=alert(1)>",
		"dob":               "2000-01-01",
		"nickname":          "&lt;script&gt;",
		"aboutMe":           "Hello <b>world</b> & <script>alert('xss')</script>",
		"profileVisibility": "public",
	}
	body, contentType, err := createSignupMultipartForm(fields)
	if err != nil {
		t.Fatalf("Failed to create multipart form: %v", err)
	}

	req := httptest.NewRequest("POST", "/signup", body)
	req.Header.Set("Content-Type", contentType)
	w := httptest.NewRecorder()
	authHandler.Signup(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var firstName, lastName, nickname, aboutMe string
	err = db.QueryRow("SELECT firstname, lastname, nickname, aboutme FROM Users WHERE email = ?", "test@example.com").
		Scan(&firstName, &lastName, &nickname, &aboutMe)
	if err != nil {
		t.Fatalf("Failed to query user: %v", err)
	}

	expectedFirstName := "&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;"
	expectedLastName := "&lt;img src=x onerror=alert(1)&gt;"
	expectedNickname := "&amp;lt;script&amp;gt;"
	expectedAboutMe := "Hello &lt;b&gt;world&lt;/b&gt; &amp; &lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;"

	if firstName != expectedFirstName {
		t.Errorf("first_name not properly escaped. Expected: %s, Got: %s", expectedFirstName, firstName)
	}
	if lastName != expectedLastName {
		t.Errorf("last_name not properly escaped. Expected: %s, Got: %s", expectedLastName, lastName)
	}
	if nickname != expectedNickname {
		t.Errorf("nickname not properly escaped. Expected: %s, Got: %s", expectedNickname, nickname)
	}
	if aboutMe != expectedAboutMe {
		t.Errorf("about_me not properly escaped. Expected: %s, Got: %s", expectedAboutMe, aboutMe)
	}
}

func TestSignupHandler_SQLInjectionVariants(t *testing.T) {
	db := setupSignupTestDB(t)
	defer db.Close()

	authStore := store.NewAuthStore(db)
	authService := service.NewAuthService(authStore)
	authHandler := handlers.NewAuthHandler(authService)

	injections := []string{
		"' OR '1'='1",
		"' OR 1=1--",
		"' OR 1=1#",
		"admin'--",
		"' UNION SELECT 1,2,3--",
		"' UNION SELECT email,password,id FROM Users--",
		"test@example.com' AND 1=1--",
		"test@example.com' AND (SELECT COUNT(*) FROM Users)>0--",
		"test@example.com'; DROP TABLE Users;--",
		"test@example.com'; INSERT INTO Users (email,password) VALUES ('hacker','hacked');--",
		"test@example.com'; UPDATE Users SET password='hacked' WHERE email='test@example.com';--",
		"test@example.com' AND sqlite_version()--",
		"test@example.com' OR 1=1 LIMIT 1--",
		"test@example.com' OR 1=1 LIMIT 1;--",
		"test@example.com' OR 1=1;--",
		"test@example.com' OR 1=1#",
		"test@example.com' OR 1=1/*",
		"test@example.com' OR 1=1-- -",
		"test@example.com' OR 1=1--+",
		"test@example.com' OR 1=1--%0A",
		"test@example.com' OR 1=1--%0D%0A",
		"test@example.com' OR 1=1--%23",
		"test@example.com' OR 1=1--%3B",
		"test@example.com' OR 1=1--%2F%2A",
		"test@example.com' OR 1=1--%2D%2D",
		"test@example.com' OR 1=1--%20",
		"test@example.com' OR 1=1--%09",
		"test@example.com' OR 1=1--%0B",
		"test@example.com' OR 1=1--%0C",
		"test@example.com' OR 1=1--%0D",
		"test@example.com' OR 1=1--%0A%0D",
		"test@example.com' OR 1=1--%0D%0A",
		"test@example.com' OR 1=1--%0A%0A",
		"test@example.com' OR 1=1--%0D%0D",
		"test@example.com' OR 1=1--%0A%0A%0A",
		"test@example.com' OR 1=1--%0D%0D%0D",
		"test@example.com' OR 1=1--%0A%0D%0A%0D",
		"test@example.com' OR 1=1--%0D%0A%0D%0A",
		"test@example.com' OR 1=1--%0A%0A%0A%0A",
		"test@example.com' OR 1=1--%0D%0D%0D%0D",
		"test@example.com' OR 1=1--%0A%0D%0A%0D%0A",
		"test@example.com' OR 1=1--%0D%0A%0D%0A%0D%0A",
		"test@example.com' OR 1=1--%0A%0A%0A%0A%0A",
		"test@example.com' OR 1=1--%0D%0D%0D%0D%0D",
		"test@example.com' OR 1=1--%0A%0D%0A%0D%0A%0D%0A",
		"test@example.com' OR 1=1--%0D%0A%0D%0A%0D%0A%0D%0A",
		"test@example.com' OR 1=1--%0A%0A%0A%0A%0A%0A",
		"test@example.com' OR 1=1--%0D%0D%0D%0D%0D%0D",
		"test@example.com' OR 1=1--%0A%0D%0A%0D%0A%0D%0A%0D%0A%0D%0A",
		"test@example.com' OR 1=1--%0D%0A%0D%0A%0D%0A%0D%0A%0D%0A",
	}
	for _, inj := range injections {
		fields := map[string]string{
			"email":             inj,
			"password":          "Password1!",
			"firstName":         "John",
			"lastName":          "Doe",
			"dob":               "2000-01-01",
			"nickname":          "hacker",
			"aboutMe":           "test",
			"profileVisibility": "public",
		}
		body, contentType, _ := createSignupMultipartForm(fields)
		req := httptest.NewRequest("POST", "/signup", body)
		req.Header.Set("Content-Type", contentType)
		w := httptest.NewRecorder()
		authHandler.Signup(w, req)
		resp := w.Result()
		if resp.StatusCode == http.StatusOK {
			t.Errorf("SQL injection variant '%s' should not succeed", inj)
		}
	}
}
