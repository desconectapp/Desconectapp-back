package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	repository "gin/db/generated"

	"github.com/gin-gonic/gin"
)

type mockActivitiesService struct {
	listResp   []repository.ActivityRequest
	listErr    error
	createResp repository.ActivityRequest
	createErr  error
}

func (m *mockActivitiesService) ListActivities() ([]repository.ActivityRequest, error) {
	return m.listResp, m.listErr
}

func (m *mockActivitiesService) CreateActivity(params repository.CreateActivityRequestParams) (repository.ActivityRequest, error) {
	return m.createResp, m.createErr
}

func TestListActivities(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &mockActivitiesService{
		listResp: []repository.ActivityRequest{{ID: 1, Activity: "hiking"}},
	}
	ctrl := NewActivitiesControllerWithService(svc)

	r := gin.Default()
	r.GET("/activities", ctrl.ListActivities)

	req := httptest.NewRequest(http.MethodGet, "/activities", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var got []repository.ActivityRequest
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if len(got) != 1 || got[0].Activity != "hiking" {
		t.Fatalf("unexpected response: %v", got)
	}
}

func TestCreateActivity(t *testing.T) {
	gin.SetMode(gin.TestMode)
	expected := repository.ActivityRequest{ID: 1, Activity: "cycling"}
	svc := &mockActivitiesService{createResp: expected}
	ctrl := NewActivitiesControllerWithService(svc)

	r := gin.Default()
	r.POST("/activities", ctrl.CreateActivity)

	body := map[string]interface{}{"activity": "cycling", "day_of_week": "EVERYDAY"}
	payload, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/activities", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var got repository.ActivityRequest
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if got.Activity != expected.Activity {
		t.Fatalf("expected activity %s, got %s", expected.Activity, got.Activity)
	}
}
