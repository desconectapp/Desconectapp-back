package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"gin/router"
	"gin/service"
	controller "gin/controller"

)

// preferences.GET("", router.preferencesController.GetUserPreferences)



// preferences.POST("", router.preferencesController.AddPreference)

func TestPostPreference(t *testing.T) {
	router := router.NewRouter()
	r :=  router.SetupRoutes()

	userID := int32(6)

	token, err := service.NewTestToken(userID)
	assert.Equal(t, err, nil, "Error should be nil")

	activityId := int32(19)

	body := AddPreference{
		ActivityID: activityId,
	}
	jsonBody, err := json.Marshal(body)

	req := httptest.NewRequest("POST", "/preferences", bytes.NewReader(jsonBody))
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response controller.ActivityPreferenseRespponse

	err = json.Unmarshal(w.Body.Bytes(), &response)
		
	assert.Equal(t, err, nil, "Error should be nil")

	
	assert.Equal(t, w.Code, http.StatusOK, "Status code should be 200")
	assert.Equal(t, activityId, response.ActivityPreferenseID, "The activity ids should match")
}


// preferences.DELETE("", router.preferencesController.DeletePreference)



// preferences.POST("/batch", router.preferencesController.BatchAddUserPreferences)     