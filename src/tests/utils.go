package test

import (
	"bytes"
	"encoding/json"
	"gin/controller"
	"gin/service"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type AuthBody struct {
	Email	string	`json:"email"`
	Password	string	`json:"password"`
}

type CreateProfile struct {
	Name             string `json:"name"`
	Age              int32  `json:"age"`
	City             string `json:"city"`
	CurrentSituation string `json:"current_situation"`
	Gender           string `json:"gender"`
}

type ActivityIdStruct struct {
    ActivityID int32 `json:"activity_id"`
}

type ActivityBatchStruct struct {
    ActivityIDBatch []int32 `json:"activity_ids"`
}

func NewUser(t *testing.T, r *gin.Engine, emailStart string) (int32, string) {

	email := emailStart + "@test.com"
	body := AuthBody{
    	Email: email,
  		Password: "password123",
    }

	jsonBody, err := json.Marshal(body)
	assert.Equal(t, err, nil, "Error should be nil")

	req := httptest.NewRequest("POST", "/auth/signup", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response controller.AuthResponse

	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, err, nil, "Error should be nil")

	_, err = service.ValidateSession(response.Token)
	assert.Equal(t, err, nil, "Error should be nil")

	_, err = service.ValidateSession(response.RefreshToken)
	assert.Equal(t, err, nil, "Error should be nil")

	assert.Equal(t, w.Code, http.StatusCreated, "Status code should be 201")

	return response.UserId, response.Token
}

func AddPreference(t *testing.T, r *gin.Engine, activityId int32, token string) {
	body := ActivityIdStruct{
		ActivityID: activityId,
	}
	jsonBody, err := json.Marshal(body)
	assert.Equal(t, err, nil, "Error should be nil")

	req := httptest.NewRequest("POST", "/preferences", bytes.NewReader(jsonBody))
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response controller.ActivityIdResponse

	err = json.Unmarshal(w.Body.Bytes(), &response)
		
	assert.Equal(t, err, nil, "Error should be nil")

	
	assert.Equal(t, w.Code, http.StatusOK, "Status code should be 200")
	assert.Equal(t, activityId, response.ActivityPreferenseID, "The activity ids should match")
}

type GroupInfo struct {
    Name        string  `json:"name"`
    Description string  `json:"description"`
    Location    string  `json:"location"`
    ActivityID  int32   `json:"activity_id"`
    MembersIds  []int32 `json:"members_ids"`
}


func NewGroup(t *testing.T, r *gin.Engine, name string, location string, activityID int32, memberIds []int32, token string) {
	body := GroupInfo{
		ActivityID: activityID,
		Name: name,
		Location: location,
		MembersIds: memberIds,

	}
	jsonBody, err := json.Marshal(body)
	
	assert.Equal(t, err, nil, "Error should be nil")

	req := httptest.NewRequest("POST", "/groups", bytes.NewReader(jsonBody))
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusCreated, "Status code should be 201")
}