package test

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"

	"gin/router"
	"gin/service"

	controller "gin/controller"
)

// // auth.POST("/signup", router.authController.Signup)

func TestSignUp(t *testing.T) {
	router := router.NewRouter()
	r := router.SetupRoutes()

	NewUser(t, r, "test")
}

func TestSignUpEmailAlreadyExists(t *testing.T) {
	router := router.NewRouter()
	r := router.SetupRoutes()

	NewUser(t, r, "test_exists")

	body := AuthBody{
		Email:    "test_exists@test.com",
		Password: "password123",
	}

	jsonBody, err := json.Marshal(body)
	assert.Equal(t, err, nil, "Error should be nil")

	req := httptest.NewRequest("POST", "/auth/signup", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusConflict, "Status code should be 409")
}

// auth.POST("/login", router.authController.Login)

func TestLogIn(t *testing.T) {
	router := router.NewRouter()
	r := router.SetupRoutes()

	body := AuthBody{
		Email:    "martina@example.com",
		Password: "password123",
	}

	jsonBody, err := json.Marshal(body)

	req := httptest.NewRequest("POST", "/auth/login", bytes.NewReader(jsonBody))
	req.Header.Add("content-type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response controller.AuthResponse

	err = json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, err, nil, "Error should be nil")

	_, _, err = service.ValidateSession(response.Token)
	assert.Equal(t, err, nil, "Error should be nil")

	_, _, err = service.ValidateSession(response.RefreshToken)
	assert.Equal(t, err, nil, "Error should be nil")

	assert.Equal(t, w.Code, http.StatusOK, "Status code should be 200")
}

// auth.POST("/refresh", router.authController.Refresh)
func TestRefresh(t *testing.T) {
	router := router.NewRouter()
	r := router.SetupRoutes()

	body := AuthBody{
		Email:    "martina@example.com",
		Password: "password123",
	}

	jsonBody, err := json.Marshal(body)

	req := httptest.NewRequest("POST", "/auth/login", bytes.NewReader(jsonBody))
	req.Header.Add("content-type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response controller.AuthResponse

	err = json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, err, nil, "Error should be nil")

	_, _, err = service.ValidateSession(response.Token)
	assert.Equal(t, err, nil, "Invalid Token")

	_, _, err = service.ValidateSession(response.RefreshToken)
	assert.Equal(t, err, nil, "Invalid Refresh Token")

	assert.Equal(t, w.Code, http.StatusOK, "Status code should be 200")

	time.Sleep(1 * time.Second)

	bodyRefresh := struct {
		RefreshToken string `json:"refresh_token"`
	}{
		RefreshToken: response.RefreshToken,
	}

	jsonBodyRefresh, err := json.Marshal(bodyRefresh)

	req_2 := httptest.NewRequest("POST", "/auth/refresh", bytes.NewReader(jsonBodyRefresh))
	req_2.Header.Add("content-type", "application/json")
	w_2 := httptest.NewRecorder()
	r.ServeHTTP(w_2, req_2)

	var response_2 controller.AuthResponse
	err = json.Unmarshal(w_2.Body.Bytes(), &response_2)
	assert.Equal(t, err, nil, "Error should be nil")

	_, _, err = service.ValidateSession(response_2.Token)
	assert.Equal(t, err, nil, "Invalid Token")

	assert.NotEqual(t, response.Token, response_2.Token)
	assert.Equal(t, w_2.Code, http.StatusOK, "Status code should be 200")
}
