package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"gin/router"
	"gin/service"

	controller "gin/controller"
)


// users.GET("", router.controller.ListUsers)
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
	
	assert.Equal(t, w.Code, http.StatusOK, "Status code should be 200")
	assert.Equal(t, limit, strconv.Itoa(len(response.Users)))
	assert.Equal(t, false, response.HasMore)
}

func TestGetUserListPagination(t *testing.T) {
	router := router.NewRouter()
	r :=  router.SetupRoutes()

	token, err := service.NewTestToken(1)

	assert.Equal(t, err, nil, "Error should be nil")

	limit := "6"
	offset := "0"

	req := httptest.NewRequest("GET", "/users?limit="+limit+"&offset="+offset, nil)
	req.Header.Add("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response controller.PaginatedUsers

	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, err, nil, "Error should be nil")
	
	assert.Equal(t, w.Code, http.StatusOK, "Status code should be 200")
	assert.Equal(t, limit, strconv.Itoa(len(response.Users)))
	assert.Equal(t, true, response.HasMore)

	offset_2 := limit

	req_2 := httptest.NewRequest("GET", "/users?limit="+limit+"&offset="+offset_2, nil)
	req_2.Header.Add("Authorization", "Bearer "+token)
	w_2 := httptest.NewRecorder()
	r.ServeHTTP(w_2, req_2)

	var response_2 controller.PaginatedUsers

	err = json.Unmarshal(w_2.Body.Bytes(), &response_2)
	assert.Equal(t, err, nil, "Error should be nil")

	
	assert.Equal(t, w_2.Code, http.StatusOK, "Status code should be 200")
	assert.LessOrEqual(t, limit, strconv.Itoa(len(response.Users)))
	assert.Equal(t, false, response_2.HasMore)
}

// users.POST("/profile", router.controller.CreateProfile)

// users.DELETE("", router.controller.DeleteUser)

func TestDeleteUser(t *testing.T) {
	router := router.NewRouter()
	r :=  router.SetupRoutes()

	userID := int32(5)

	token, err := service.NewTestToken(userID)

	assert.Equal(t, err, nil, "Error should be nil")

	req := httptest.NewRequest("DELETE", "/users", nil)
	req.Header.Add("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response controller.UserDeletedResponse

	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, err, nil, "Error should be nil")
	
	assert.Equal(t, w.Code, http.StatusOK, "Status code should be 200")
	assert.Equal(t, userID, response.DeletedUserID, "The user ids should match")
}

// users.GET("/user", router.controller.GetUser)
