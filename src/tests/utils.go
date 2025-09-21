package test

import (
	"bytes"
	"encoding/json"
	"gin/controller"
	"gin/service"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type AuthBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateProfile struct {
	Name             string `json:"name"`
	Age              int32  `json:"age"`
	City             string `json:"city"`
	CurrentSituation string `json:"current_situation"`
	Gender           string `json:"gender"`
}

type TimeSlot struct {
	Start int `json:"start"`
	End   int `json:"end"`
}
type Schedules map[string][]TimeSlot

type CreateActivityRequestInput struct {
	UserID             *int32    `json:"user_id"`
	ActivityID         *int32    `json:"activity_id"`
	Description        *string   `json:"description"`
	Longitude          *float64  `json:"longitude"`
	Latitude           *float64  `json:"latitude"`
	SearchRadius       *int32    `json:"search_radius"`
	MaxParticipants    *int32    `json:"max_participants"`
	ParticipantsNeeded *int32    `json:"participants_needed"`
	Schedules          Schedules `json:"schedules"`
}

type ActivityRequest struct {
	ID                  int32            `json:"id"`
	UserID              *int32           `json:"user_id"`
	ActivityID          *int32           `json:"activity_id"`
	Description         *string          `json:"description"`
	WeekHours           []int32          `json:"week_hours"`
	ParticipantsNeeded  *int32           `json:"participants_needed"`
	MaximumParticipants *int32           `json:"maximum_participants"`
	Latitude            *float64         `json:"latitude"`
	Longitude           *float64         `json:"longitude"`
	SearchRadius        *int32           `json:"search_radius"`
	CreatedAt           pgtype.Timestamp `json:"created_at"`
	ExpiresAt           pgtype.Timestamp `json:"expires_at"`
}
type ActivityIdStruct struct {
	ActivityID int32 `json:"activity_id"`
}

type ActivityBatchStruct struct {
	ActivityIDBatch []int32 `json:"activity_ids"`
}

type NewDescription struct {
	NewDesc    string `json:"description"`
}

type NewStatus struct {
	NewStatus    bool `json:"status"`
}

func SendActivityRequest(t *testing.T, r *gin.Engine, userId int32, activityId int32, participantsNeeded int32) ActivityRequest {
	int32Ptr := func(i int32) *int32 { return &i }
	float64Ptr := func(f float64) *float64 { return &f }
	strPtr := func(s string) *string { return &s }

	body := CreateActivityRequestInput{
		UserID:             int32Ptr(userId),
		ActivityID:         int32Ptr(activityId),
		Description:        strPtr("Looking for a running buddy"),
		ParticipantsNeeded: int32Ptr(participantsNeeded),
		MaxParticipants:    int32Ptr(5),
		Latitude:           float64Ptr(37.7749),
		Longitude:          float64Ptr(-122.4194),
		SearchRadius:       int32Ptr(10),
		Schedules: Schedules{
			"monday":    {{Start: 9, End: 11}, {Start: 14, End: 16}},
			"wednesday": {{Start: 10, End: 12}},
		},
	}

	jsonBody, err := json.Marshal(body)
	assert.Equal(t, err, nil, "Error should be nil")

	token, err := service.NewTestToken(userId)
	assert.Equal(t, err, nil, "Error should be nil")

	req := httptest.NewRequest("POST", "/activities/request", bytes.NewReader(jsonBody))
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Status code should be 200")
	if w.Code != http.StatusOK {
		t.Fatalf("Unexpected status code: %d, body: %s", w.Code, w.Body.String())
	}

	var response ActivityRequest
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, err, nil, "Error should be nil")

	assert.Equal(t, *body.UserID, *response.UserID, "UserID should match")
	assert.Equal(t, *body.ActivityID, *response.ActivityID, "ActivityID should match")
	assert.Equal(t, *body.Description, *response.Description, "Description should match")
	assert.Equal(t, *body.ParticipantsNeeded, *response.ParticipantsNeeded, "Participants Needed should match")
	assert.Equal(t, *body.MaxParticipants, *response.MaximumParticipants, "Maximum Participants should match")
	assert.Equal(t, *body.Latitude, *response.Latitude, "Latitude should match")
	assert.Equal(t, *body.Longitude, *response.Longitude, "Longitude should match")
	assert.Equal(t, *body.SearchRadius, *response.SearchRadius, "SearchRadius should match")

	return response
}

func NewUser(t *testing.T, r *gin.Engine, emailStart string) (int32, string) {

	email := emailStart + "@test.com"
	body := AuthBody{
		Email:    email,
		Password: "password123",
	}

	jsonBody, err := json.Marshal(body)
	assert.Equal(t, err, nil, "Error should be nil")

	req := httptest.NewRequest("POST", "/auth/signup", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response controller.AuthResponse

	log.Println(response)

	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, err, nil, "Error should be nil")

	_, _, err = service.ValidateSession(response.Token)
	assert.Equal(t, err, nil, "Error should be nil")

	_, _, err = service.ValidateSession(response.RefreshToken)
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
    MembersIds  []int32 `json:"user_ids"`
}


func NewGroup(t *testing.T, r *gin.Engine, name string, location string, activityID int32, memberIds []int32, token string) int32 {
	body := GroupInfo{
		ActivityID: activityID,
		Name: name,
		Location: location,
		MembersIds: memberIds,
		Description: "",

	}
	jsonBody, err := json.Marshal(body)
	
	assert.Equal(t, err, nil, "Error should be nil")

	req := httptest.NewRequest("POST", "/groups", bytes.NewReader(jsonBody))
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response controller.NewGroup

	err = json.Unmarshal(w.Body.Bytes(), &response)
		
	assert.Equal(t, err, nil, "Error should be nil")
	assert.Equal(t, w.Code, http.StatusCreated, "Status code should be 201")
	assert.Equal(t, body.ActivityID, response.ActivityID)
	assert.Equal(t, body.Name, *response.Name)
	assert.Equal(t, body.Location, *response.Location)
	assert.Equal(t, body.MembersIds, response.Members)
	assert.Equal(t, false, *response.Status)

	return response.ID
}

func NewUserWithProfile(t *testing.T, r *gin.Engine, emailStart string, profile CreateProfile) (int32, string) {
	userID, token := NewUser(t, r, emailStart)

	jsonBody, err := json.Marshal(profile)
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
	assert.Equal(t, profile.Age, response.Age, "Age should match")
	assert.Equal(t, profile.Name, response.Name, "Name should match")
	assert.Equal(t, profile.City, response.City, "City should match")
	assert.Equal(t, profile.CurrentSituation, response.CurrentSituation, "Current Situation should match")
	assert.Equal(t, profile.Gender, response.Gender, "Gender should match")

	return userID, token
}
