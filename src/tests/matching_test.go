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

func TestTwoCompatibleActivityRequestsMatch(t *testing.T) {
	router := router.NewRouter()
	r := router.SetupRoutes()

	SendActivityRequest(t, r, 1, 1, 2)
	SendActivityRequest(t, r, 2, 1, 2)

	token, err := service.NewTestToken(2)
	assert.Nil(t, err, "Error should be nil")

	req := httptest.NewRequest("GET", "/groups/user", nil)
	req.Header.Add("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code, "Status code should be 200")
	if w.Code != 200 {
		t.Fatalf("Unexpected status code: %d, body: %s", w.Code, w.Body.String())
	}
	var groupsResponse controller.PaginatedMembers
	err = json.Unmarshal(w.Body.Bytes(), &groupsResponse)
	assert.Nil(t, err, "Error unmarshaling response")
	assert.GreaterOrEqual(t, 1, len(groupsResponse.Members), "There should be at least one group")
}

func TestMatchDoesNotDeletePreviousSearch(t *testing.T) {
	router := router.NewRouter()
	r := router.SetupRoutes()

	SendActivityRequest(t, r, 3, 1, 2)
	SendActivityRequest(t, r, 3, 2, 2)

	token, err := service.NewTestToken(2)
	assert.Nil(t, err, "Error should be nil")

	SendActivityRequest(t, r, 4, 1, 2)

	req := httptest.NewRequest("GET", "/groups/user", nil)
	req.Header.Add("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code, "Status code should be 200")
	if w.Code != 200 {
		t.Fatalf("Unexpected status code: %d, body: %s", w.Code, w.Body.String())
	}
	var groupsResponse controller.PaginatedMembers
	err = json.Unmarshal(w.Body.Bytes(), &groupsResponse)
	assert.Nil(t, err, "Error unmarshaling response")
	assert.GreaterOrEqual(t, 1, len(groupsResponse.Members), "There should be at least one group")

	SendActivityRequest(t, r, 4, 2, 2)

	req = httptest.NewRequest("GET", "/groups/user", nil)
	req.Header.Add("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code, "Status code should be 200")
	if w.Code != 200 {
		t.Fatalf("Unexpected status code: %d, body: %s", w.Code, w.Body.String())
	}
	var groupsResponse_2 controller.PaginatedMembers
	err = json.Unmarshal(w.Body.Bytes(), &groupsResponse_2)
	assert.Nil(t, err, "Error unmarshaling response")
	assert.GreaterOrEqual(t, 1, len(groupsResponse_2.Members), "There should be at least one group")
}

func TestUserCantMatchWithTheSameUser(t *testing.T) {
	router := router.NewRouter()
	r := router.SetupRoutes()

	SendActivityRequest(t, r, 10, 1, 2)
	SendActivityRequest(t, r, 10, 1, 2)
	token, err := service.NewTestToken(10)
	assert.Nil(t, err, "Error should be nil")
	req := httptest.NewRequest("GET", "/groups/user", nil)
	req.Header.Add("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code, "Status code should be 200")
	if w.Code != 200 {
		t.Fatalf("Unexpected status code: %d, body: %s", w.Code, w.Body.String())
	}
	var groupsResponse controller.PaginatedMembers
	err = json.Unmarshal(w.Body.Bytes(), &groupsResponse)
	assert.Nil(t, err, "Error unmarshaling response")
	assert.Equal(t, 0, len(groupsResponse.Members), "There should be no groups")
}

func TestRequestsForDifferentActivitiesDontMatch(t *testing.T) {
	router := router.NewRouter()
	r := router.SetupRoutes()

	SendActivityRequest(t, r, 5, 4, 2)
	SendActivityRequest(t, r, 6, 5, 2)

	token, err := service.NewTestToken(11)
	assert.Nil(t, err, "Error should be nil")
	req := httptest.NewRequest("GET", "/groups/user", nil)
	req.Header.Add("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code, "Status code should be 200")
	if w.Code != 200 {
		t.Fatalf("Unexpected status code: %d, body: %s", w.Code, w.Body.String())
	}
	var groupsResponse controller.PaginatedMembers
	err = json.Unmarshal(w.Body.Bytes(), &groupsResponse)
	assert.Nil(t, err, "Error unmarshaling response")
	assert.Equal(t, 0, len(groupsResponse.Members), "There should be no groups")
}

func TestRequestWithoutSchedulesIsNotCreated(t *testing.T) {
	router := router.NewRouter()
	r := router.SetupRoutes()
	int32Ptr := func(i int32) *int32 { return &i }
	float64Ptr := func(f float64) *float64 { return &f }
	strPtr := func(s string) *string { return &s }

	body := CreateActivityRequestInput{
		UserID:             int32Ptr(12),
		ActivityID:         int32Ptr(3),
		Description:        strPtr("Looking for a running buddy"),
		ParticipantsNeeded: int32Ptr(2),
		MaxParticipants:    int32Ptr(5),
		Latitude:           float64Ptr(37.7749),
		Longitude:          float64Ptr(-122.4194),
		SearchRadius:       int32Ptr(10),
		Schedules:          Schedules{},
	}

	jsonBody, err := json.Marshal(body)
	assert.Equal(t, err, nil, "Error should be nil")

	token, err := service.NewTestToken(12)
	assert.Equal(t, err, nil, "Error should be nil")

	req := httptest.NewRequest("POST", "/activities/request", bytes.NewReader(jsonBody))
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code, "Status code should be 400, bad request")
	if w.Code != 400 {
		t.Fatalf("Unexpected status code: %d, body: %s", w.Code, w.Body.String())
	}
}

func TestRequestsTooFarApartDontMatch(t *testing.T) {

	router := router.NewRouter()
	r := router.SetupRoutes()
	int32Ptr := func(i int32) *int32 { return &i }
	float64Ptr := func(f float64) *float64 { return &f }
	strPtr := func(s string) *string { return &s }

	body := CreateActivityRequestInput{
		UserID:             int32Ptr(7),
		ActivityID:         int32Ptr(3),
		Description:        strPtr("Looking for a running buddy"),
		ParticipantsNeeded: int32Ptr(2),
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

	token, err := service.NewTestToken(13)
	assert.Equal(t, err, nil, "Error should be nil")

	req := httptest.NewRequest("POST", "/activities/request", bytes.NewReader(jsonBody))
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
	assert.Equal(t, err, nil, "Error should be nil")

	body = CreateActivityRequestInput{
		UserID:             int32Ptr(8),
		ActivityID:         int32Ptr(3),
		Description:        strPtr("Looking for a running buddy"),
		ParticipantsNeeded: int32Ptr(2),
		MaxParticipants:    int32Ptr(5),
		Latitude:           float64Ptr(89.7749),
		Longitude:          float64Ptr(-172.4194),
		SearchRadius:       int32Ptr(10),
		Schedules: Schedules{
			"monday":    {{Start: 9, End: 11}, {Start: 14, End: 16}},
			"wednesday": {{Start: 10, End: 12}},
		},
	}

	jsonBody, err = json.Marshal(body)
	assert.Equal(t, err, nil, "Error should be nil")
	token, err = service.NewTestToken(14)
	assert.Equal(t, err, nil, "Error should be nil")
	req = httptest.NewRequest("POST", "/activities/request", bytes.NewReader(jsonBody))
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code, "Status code should be 200")
	if w.Code != 200 {
		t.Fatalf("Unexpected status code: %d, body: %s", w.Code, w.Body.String())
	}
	var groupsResponse controller.PaginatedMembers
	err = json.Unmarshal(w.Body.Bytes(), &groupsResponse)
	assert.Nil(t, err, "Error unmarshaling response")
	assert.Equal(t, 0, len(groupsResponse.Members), "There should be no groups")

}
