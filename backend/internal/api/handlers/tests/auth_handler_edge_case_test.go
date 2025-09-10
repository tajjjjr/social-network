package tests

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/tajjjjr/social-network/backend/internal/api/handlers"
	"github.com/tajjjjr/social-network/backend/internal/models"
	"github.com/tajjjjr/social-network/backend/internal/service"
	"github.com/tajjjjr/social-network/backend/internal/store"
	"github.com/tajjjjr/social-network/backend/pkg/utils"
)

// setupEdgeCaseTestDB creates an in-memory SQLite database for edge case testing
func setupEdgeCaseTestDB(t *testing.T) *sql.DB {
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
			dateofbirth TEXT,
			nickname TEXT,
			aboutme TEXT,
			is_profile_public BOOLEAN,
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

// createEdgeCaseTestUser creates a test user in the database
func createEdgeCaseTestUser(t *testing.T, db *sql.DB) {
	passwordManager := utils.NewPasswordManager(utils.PasswordConfig{})
	hashedPassword, err := passwordManager.HashPassword("TestPassword123!")
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	_, err = db.Exec(`
		INSERT INTO Users (email, password, firstname, lastname, nickname, is_profile_public)
		VALUES (?, ?, ?, ?, ?, ?)
	`, "test@example.com", hashedPassword, "Test", "User", "testuser", true)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
}

// Test expired sessions
func TestExpiredSessions(t *testing.T) {
	db := setupEdgeCaseTestDB(t)
	defer db.Close()

	createEdgeCaseTestUser(t, db)

	authStore := store.NewAuthStore(db)
	authService := service.NewAuthService(authStore)
	authHandler := handlers.NewAuthHandler(authService)

	// Create an expired session manually
	expiredTime := time.Now().Add(-1 * time.Hour) // 1 hour ago
	_, err := db.Exec(`
		INSERT INTO Sessions (id, user_id, created_at, expires_at)
		VALUES (?, ?, ?, ?)
	`, "expired-session-123", 1, time.Now().Add(-2*time.Hour), expiredTime)
	if err != nil {
		t.Fatalf("Failed to create expired session: %v", err)
	}

	// Try to access protected route with expired session
	req := httptest.NewRequest("GET", "/protected", nil)
	req.AddCookie(&http.Cookie{
		Name:  "session_id",
		Value: "expired-session-123",
	})

	rr := httptest.NewRecorder()
	protectedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.RespondJSON(w, http.StatusOK, utils.Response{Message: "Access granted"})
	})

	authMiddleware := authHandler.AuthMiddleware(protectedHandler)
	authMiddleware.ServeHTTP(rr, req)

	// Should deny access for expired session
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code for expired session: got %v want %v",
			status, http.StatusUnauthorized)
	}

	var resp utils.Response
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	if resp.Message != "Invalid or expired session" {
		t.Errorf("expected 'Invalid or expired session', got %q", resp.Message)
	}
}

// Test invalid session IDs
func TestInvalidSessionIDs(t *testing.T) {
	db := setupEdgeCaseTestDB(t)
	defer db.Close()

	createEdgeCaseTestUser(t, db)

	authStore := store.NewAuthStore(db)
	authService := service.NewAuthService(authStore)
	authHandler := handlers.NewAuthHandler(authService)

	invalidSessionIDs := []string{
		"non-existent-session",
		"",
		"invalid-uuid-format",
		"sql-injection'; DROP TABLE Sessions;--",
		"very-long-session-id-that-exceeds-normal-length-" + strings.Repeat("a", 1000),
		"session-with-special-chars-!@#$%^&*()",
		"session\nwith\nnewlines",
		"session\twith\ttabs",
	}

	for _, sessionID := range invalidSessionIDs {
		t.Run("InvalidSession_"+sessionID[:min(20, len(sessionID))], func(t *testing.T) {
			req := httptest.NewRequest("GET", "/protected", nil)
			req.AddCookie(&http.Cookie{
				Name:  "session_id",
				Value: sessionID,
			})

			rr := httptest.NewRecorder()
			protectedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				utils.RespondJSON(w, http.StatusOK, utils.Response{Message: "Access granted"})
			})

			authMiddleware := authHandler.AuthMiddleware(protectedHandler)
			authMiddleware.ServeHTTP(rr, req)

			// Should deny access for invalid session
			if status := rr.Code; status != http.StatusUnauthorized {
				t.Errorf("handler returned wrong status code for invalid session '%s': got %v want %v",
					sessionID, status, http.StatusUnauthorized)
			}
		})
	}
}

// Test concurrent logins
func TestConcurrentLogins(t *testing.T) {
	db := setupEdgeCaseTestDB(t)
	defer db.Close()

	createEdgeCaseTestUser(t, db)

	authStore := store.NewAuthStore(db)
	authService := service.NewAuthService(authStore)
	authHandler := handlers.NewAuthHandler(authService)

	// Reduce concurrency for SQLite compatibility
	const numConcurrentLogins = 3
	var wg sync.WaitGroup
	results := make([]int, numConcurrentLogins)
	sessionIDs := make([]string, numConcurrentLogins)

	// Add a mutex to prevent SQLite database lock issues
	var mu sync.Mutex

	// Perform concurrent logins
	for i := 0; i < numConcurrentLogins; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			// Add small delay to reduce database contention
			time.Sleep(time.Duration(index) * 10 * time.Millisecond)

			loginReq := models.LoginRequest{
				Email:    "test@example.com",
				Password: "TestPassword123!",
			}
			body, _ := json.Marshal(loginReq)
			req, err := http.NewRequest("POST", "/login", strings.NewReader(string(body)))
			if err != nil {
				mu.Lock()
				results[index] = http.StatusInternalServerError
				mu.Unlock()
				return
			}
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			authHandler.Login(rr, req)

			mu.Lock()
			results[index] = rr.Code

			// Extract session ID from cookie
			cookies := rr.Result().Cookies()
			for _, cookie := range cookies {
				if cookie.Name == "session_id" {
					sessionIDs[index] = cookie.Value
					break
				}
			}
			mu.Unlock()
		}(i)
	}

	wg.Wait()

	// Verify all logins succeeded
	successCount := 0
	for i, statusCode := range results {
		if statusCode == http.StatusOK {
			successCount++
		} else {
			t.Errorf("Login %d failed with status %d", i, statusCode)
		}
	}

	if successCount != numConcurrentLogins {
		t.Errorf("Expected %d successful logins, got %d", numConcurrentLogins, successCount)
	}

	// Verify all sessions are unique
	sessionSet := make(map[string]bool)
	for i, sessionID := range sessionIDs {
		if sessionID == "" {
			t.Errorf("Login %d did not receive a session ID", i)
			continue
		}
		if sessionSet[sessionID] {
			t.Errorf("Duplicate session ID found: %s", sessionID)
		}
		sessionSet[sessionID] = true
	}

	// Verify all sessions are valid and can access protected routes
	for i, sessionID := range sessionIDs {
		if sessionID == "" {
			continue
		}

		req := httptest.NewRequest("GET", "/protected", nil)
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: sessionID,
		})

		rr := httptest.NewRecorder()
		protectedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			utils.RespondJSON(w, http.StatusOK, utils.Response{Message: "Access granted"})
		})

		authMiddleware := authHandler.AuthMiddleware(protectedHandler)
		authMiddleware.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Session %d (%s) should be valid for protected route access", i, sessionID)
		}
	}
}

// Helper function for min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Test session cleanup after multiple logins (if implemented)
func TestSessionCleanup(t *testing.T) {
	db := setupEdgeCaseTestDB(t)
	defer db.Close()

	createEdgeCaseTestUser(t, db)

	authStore := store.NewAuthStore(db)
	authService := service.NewAuthService(authStore)
	authHandler := handlers.NewAuthHandler(authService)

	// Perform multiple logins for the same user
	sessionIDs := make([]string, 3)
	var err error
	for i := 0; i < 3; i++ {
		loginReq := models.LoginRequest{
			Email:    "test@example.com",
			Password: "TestPassword123!",
		}
		body, _ := json.Marshal(loginReq)
		req, err := http.NewRequest("POST", "/login", strings.NewReader(string(body)))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		authHandler.Login(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("Login %d failed", i)
		}

		// Extract session ID
		cookies := rr.Result().Cookies()
		for _, cookie := range cookies {
			if cookie.Name == "session_id" {
				sessionIDs[i] = cookie.Value
				break
			}
		}
	}

	// Count total sessions for this user
	var sessionCount int
	err = db.QueryRow("SELECT COUNT(*) FROM Sessions WHERE user_id = 1").Scan(&sessionCount)
	if err != nil {
		t.Fatal(err)
	}

	// Note: This test documents current behavior.
	// If session cleanup is implemented, this test should be updated accordingly.
	t.Logf("Total sessions for user after 3 logins: %d", sessionCount)

	// All sessions should still be valid (current behavior)
	for i, sessionID := range sessionIDs {
		if sessionID == "" {
			t.Errorf("Login %d did not receive session ID", i)
			continue
		}

		req := httptest.NewRequest("GET", "/protected", nil)
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: sessionID,
		})

		rr := httptest.NewRecorder()
		protectedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			utils.RespondJSON(w, http.StatusOK, utils.Response{Message: "Access granted"})
		})

		authMiddleware := authHandler.AuthMiddleware(protectedHandler)
		authMiddleware.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Session %d should be valid", i)
		}
	}
}
