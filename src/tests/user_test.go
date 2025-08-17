package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"gin/router"
	"gin/service"
)

func TestGetUser(t *testing.T) {
	router := router.NewRouter()
	r :=  router.SetupRoutes()

	token, err := service.NewTestToken(1)

	assert.Equal(t, err, nil, "Error should be nil")

	req := httptest.NewRequest("GET", "/users", nil)
	req.Header.Add("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusOK, "Status code should be 200")

}