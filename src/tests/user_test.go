package test

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"gin/router"
	"gin/service"
	
	controller "gin/controller"
)

const (
	USERS_NUMBER_FILE = 11
)

func TestGetUserList(t *testing.T) {
	router := router.NewRouter()
	r :=  router.SetupRoutes()

	token, err := service.NewTestToken(1)

	assert.Equal(t, err, nil, "Error should be nil")

	limit := "11"
	offset := "0"

	req := httptest.NewRequest("GET", "/users?limit="+limit+"&offset="+offset, nil)
	req.Header.Add("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response controller.PaginatedUsers

	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, err, nil, "Error should be nil")

	log.Println(response)
	
	assert.Equal(t, w.Code, http.StatusOK, "Status code should be 200")
	assert.Equal(t, USERS_NUMBER_FILE, len(response.Users))
}