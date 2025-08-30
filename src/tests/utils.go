package test

import (
	"bytes"
	"encoding/json"
	"gin/controller"
	"gin/service"
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
