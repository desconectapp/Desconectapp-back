package test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	models "gin/db/generated"
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

func TestGetUser(t *testing.T) {
	router := router.NewRouter()
	r :=  router.SetupRoutes()

	userID := int32(6)

	token, err := service.NewTestToken(userID)

	assert.Equal(t, err, nil, "Error should be nil")

	req := httptest.NewRequest("GET", "/users/user", nil)
	req.Header.Add("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response models.Profile

	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, err, nil, "Error should be nil")

	log.Println(response)
	
	assert.Equal(t, w.Code, http.StatusOK, "Status code should be 200")
	assert.Equal(t, userID, response.UserID, "The user ids should match")
}

// users.POST("/profile", router.controller.CreateProfile)

func TestUserProfile(t *testing.T) {
	router := router.NewRouter()
	r :=  router.SetupRoutes()

	userID := int32(8)
	token, err := service.NewTestToken(userID)

	age := int32(24)
	name := "Martina"
	city := "Buenos Aires"
	currentSituation :=  "UNEMPLOYED"
	gender := "Female"

	body := CreateProfile{
		Name: name,
		Age: age,
		City: city,
		CurrentSituation: currentSituation,
		Gender: gender,
	}

	jsonBody, err := json.Marshal(body)

	req := httptest.NewRequest("POST", "/users/profile", bytes.NewReader(jsonBody))
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response controller.UserProfileCreatedResponse

	err = json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, err, nil, "Error should be nil")

	assert.Equal(t, w.Code, http.StatusOK, "Status code should be 200")
	assert.Equal(t, userID, response.ProfileUserID, "The user ids should match")

	req_2 := httptest.NewRequest("GET", "/users/user", nil)
	req_2.Header.Add("Authorization", "Bearer "+token)
	w_2 := httptest.NewRecorder()
	r.ServeHTTP(w_2, req_2)

	var response_2 models.Profile

	err = json.Unmarshal(w_2.Body.Bytes(), &response_2)
	assert.Equal(t, err, nil, "Error should be nil")

	
	assert.Equal(t, w_2.Code, http.StatusOK, "Status code should be 200")
	assert.Equal(t, userID, response_2.UserID, "The user ids should match")
	assert.Equal(t, age, response_2.Age, "Age should match")
	assert.Equal(t, name, response_2.Name, "Name should match")
	assert.Equal(t, city, response_2.City, "City should match")
	assert.Equal(t, currentSituation, response_2.CurrentSituation, "Current Situation should match")
	assert.Equal(t, gender, response_2.Gender, "Gender should match")
}