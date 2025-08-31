package test

import (
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
	assert.GreaterOrEqual(t, len(groupsResponse.Members), 1, "There should be at least one group")
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
	assert.GreaterOrEqual(t, len(groupsResponse.Members), 1, "There should be at least one group")

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
	assert.GreaterOrEqual(t, len(groupsResponse_2.Members), 1, "There should be at least one group")
}
