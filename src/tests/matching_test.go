package test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"gin/controller"
	"gin/router"
	"gin/service"

	"github.com/stretchr/testify/assert"
)

func TestPostActivityRequest(t *testing.T) {
	router := router.NewRouter()
	r := router.SetupRoutes()

	token, err := service.NewTestToken(1)
	assert.Nil(t, err, "Error should be nil")

	body := controller.CreateActivityRequestInput{
		UserID:             int32Ptr(1),
		ActivityID:         int32Ptr(1),
		Description:        strPtr("Looking for a running buddy"),
		ParticipantsNeeded: int32Ptr(2),
		MaxParticipants:    int32Ptr(5),
		Latitude:           float64Ptr(37.7749),
		Longitude:          float64Ptr(-122.4194),
		SearchRadius:       int32Ptr(10),
		Schedules: controller.Schedules{
			"monday":    {{Start: 9, End: 11}, {Start: 14, End: 16}},
			"wednesday": {{Start: 10, End: 12}},
		},
	}
	jsonBody, err := json.Marshal(body)
	assert.Nil(t, err, "Error marshaling body to JSON")

	req := httptest.NewRequest("POST", "/activities/request", bytes.NewBuffer(jsonBody))
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code, "Status code should be 200")
	if w.Code != 200 {
		t.Fatalf("Unexpected status code: %d, body: %s", w.Code, w.Body.String())
	}

	var response ActivityRequest
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err, "Error unmarshaling response")
	assert.Equal(t, *body.UserID, *response.UserID, "UserID should match")
	assert.Equal(t, *body.ActivityID, *response.ActivityID, "ActivityID should match")
	assert.Equal(t, *body.Description, *response.Description, "Description should match")
	assert.Equal(t, *body.ParticipantsNeeded, *response.ParticipantsNeeded, "ParticipantsNeeded should match")
	assert.Equal(t, *body.MaxParticipants, *response.MaximumParticipants, "MaximumParticipants should match")
	assert.Equal(t, *body.Latitude, *response.Latitude, "Latitude should match")
	assert.Equal(t, *body.Longitude, *response.Longitude, "Longitude should match")
	assert.Equal(t, *body.SearchRadius, *response.SearchRadius, "SearchRadius should match")
}

func int32Ptr(i int32) *int32       { return &i }
func float64Ptr(f float64) *float64 { return &f }
func strPtr(s string) *string       { return &s }
