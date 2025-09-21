package test

import (
	"gin/router"
	"gin/service"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChatsGetToken(t *testing.T) {
	r := router.NewRouter()
	router := r.SetupRoutes()

	userID, _ := NewUser(t, router, "chattest")

	token, err := service.NewTestToken(userID)
	assert.Equal(t, err, nil, "Error should be nil")

	t.Run("should return supabase token for authenticated user", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/chats/token", nil)
		req.Header.Add("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Status code should be 200")

		responseBody := w.Body.String()
		assert.Contains(t, responseBody, "supabase_token", "Response should contain supabase_token")

		assert.Contains(t, responseBody, `"supabase_token":`, "Response should have supabase_token field")
	})

	t.Run("should return 401 for unauthenticated request", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/chats/token", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code, "Status code should be 401")
	})

	t.Run("should return 401 for invalid token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/chats/token", nil)
		req.Header.Add("Authorization", "Bearer invalid-token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code, "Status code should be 401")
	})

	t.Run("should accept token without Bearer prefix", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/chats/token", nil)
		req.Header.Add("Authorization", token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Status code should be 200")

		responseBody := w.Body.String()
		assert.Contains(t, responseBody, "supabase_token", "Response should contain supabase_token")
	})

}
