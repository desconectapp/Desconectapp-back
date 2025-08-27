package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

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

// users.DELETE("", router.controller.DeleteUser)

func TestDeleteUser(t *testing.T) {
	router := router.NewRouter()
	r :=  router.SetupRoutes()

	userId, token := NewUser(t, r, "test_delete")

	time.Sleep(1 * time.Second)

	req_2 := httptest.NewRequest("DELETE", "/users", nil)
	req_2.Header.Add("Authorization", "Bearer "+token)
	w_2 := httptest.NewRecorder()
	r.ServeHTTP(w_2, req_2)

	var response_2 controller.UserDeletedResponse

	err := json.Unmarshal(w_2.Body.Bytes(), &response_2)
	assert.Equal(t, err, nil, "Error should be nil")

	
	assert.Equal(t, w_2.Code, http.StatusOK, "Status code should be 200")
	assert.Equal(t, userId, response_2.DeletedUserID, "The user ids should match")
}


func TestDeleteNonExistentUser(t *testing.T) {
	router := router.NewRouter()
	r :=  router.SetupRoutes()

	userID := int32(0)
	token, err := service.NewTestToken(userID)
	assert.Equal(t, err, nil, "Error should be nil")

	req_2 := httptest.NewRequest("DELETE", "/users", nil)
	req_2.Header.Add("Authorization", "Bearer "+token)
	w_2 := httptest.NewRecorder()
	r.ServeHTTP(w_2, req_2)

	assert.Equal(t, w_2.Code, http.StatusNotFound, "Status code should be 200")
}

// users.GET("/user", router.controller.GetUser)

func TestGetUser(t *testing.T) {
	router := router.NewRouter()
	r :=  router.SetupRoutes()

	userID := int32(3)

	token, err := service.NewTestToken(userID)

	assert.Equal(t, err, nil, "Error should be nil")

	req := httptest.NewRequest("GET", "/users/user", nil)
	req.Header.Add("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response controller.Profile

	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, err, nil, "Error should be nil")
	
	assert.Equal(t, w.Code, http.StatusOK, "Status code should be 200")
	assert.Equal(t, userID, response.UserID, "The user ids should match")
}

func TestGetNonExistentUser(t *testing.T) {
	router := router.NewRouter()
	r :=  router.SetupRoutes()

	userID := int32(0)
	token, err := service.NewTestToken(userID)
	assert.Equal(t, err, nil, "Error should be nil")

	req := httptest.NewRequest("GET", "/users/user", nil)
	req.Header.Add("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusNotFound, "Status code should be 404")
}

// users.POST("/profile", router.controller.CreateProfile)

func TestCreateUserProfile(t *testing.T) {
	router := router.NewRouter()
	r :=  router.SetupRoutes()

	userID, token := NewUser(t, r, "maomao")

	age := int32(24)
	name := "Maomao"
	city := "Japan"
	currentSituation :=  "WORKING"
	gender := "Female"

	body := CreateProfile{
		Name: name,
		Age: age,
		City: city,
		CurrentSituation: currentSituation,
		Gender: gender,
	}

	jsonBody, err := json.Marshal(body)
	assert.Equal(t, err, nil, "Error should be nil")

	req := httptest.NewRequest("POST", "/users/profile", bytes.NewReader(jsonBody))
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response controller.Profile

	err = json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, err, nil, "Error should be nil")

	assert.Equal(t, w.Code, http.StatusOK, "Status code should be 200")
	assert.Equal(t, userID, response.UserID, "The user ids should match")
	assert.Equal(t, age, response.Age, "Age should match")
	assert.Equal(t, name, response.Name, "Name should match")
	assert.Equal(t, city, response.City, "City should match")
	assert.Equal(t, currentSituation, response.CurrentSituation, "Current Situation should match")
	assert.Equal(t, gender, response.Gender, "Gender should match")
}

func TestCreateUserProfileBadBind(t *testing.T) {
	router := router.NewRouter()
	r :=  router.SetupRoutes()

	_, token := NewUser(t, r, "invalid_bind")

	body := AuthBody{
		Email: "invalid@test.com",
		Password: "lets_fail",
	}

	jsonBody, err := json.Marshal(body)
	assert.Equal(t, err, nil, "Error should be nil")

	req := httptest.NewRequest("POST", "/users/profile", bytes.NewReader(jsonBody))
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusBadRequest, "Status code should be 400")

}

func TestCreateUserProfileInvalidAge(t *testing.T) {
	router := router.NewRouter()
	r :=  router.SetupRoutes()

	
	_, token := NewUser(t, r, "invalid_age")

	age := int32(240)
	name := "Maomao"
	city := "Japan"
	currentSituation :=  "WORKING"
	gender := "Female"

	body := CreateProfile{
		Name: name,
		Age: age,
		City: city,
		CurrentSituation: currentSituation,
		Gender: gender,
	}

	jsonBody, err := json.Marshal(body)
	assert.Equal(t, err, nil, "Error should be nil")

	req := httptest.NewRequest("POST", "/users/profile", bytes.NewReader(jsonBody))
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusBadRequest, "Status code should be 400")

}