package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	controller "gin/controller"
	"gin/router"
	"gin/service"
)

func addPreference(t *testing.T, r *gin.Engine, activityId int32, token string) {
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



// preferences.POST("", router.preferencesController.AddPreference)

func TestPostPreference(t *testing.T) {
	router := router.NewRouter()
	r :=  router.SetupRoutes()
	
	userID := int32(6)
	
	token, err := service.NewTestToken(userID)
	assert.Equal(t, err, nil, "Error should be nil")
	
	activityId := int32(19)
	
	addPreference(t, r, activityId, token)
}


// preferences.GET("", router.preferencesController.GetUserPreferences)
func TestGetPreference(t *testing.T) {
	router := router.NewRouter()
	r :=  router.SetupRoutes()

	userID := int32(7)

	token, err := service.NewTestToken(userID)
	assert.Equal(t, err, nil, "Error should be nil")

	activityId := int32(1)

	addPreference(t, r, activityId, token)

	req := httptest.NewRequest("GET", "/preferences", nil)
	req.Header.Add("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response controller.PaginatedPreferences

	err = json.Unmarshal(w.Body.Bytes(), &response)
		
	assert.Equal(t, err, nil, "Error should be nil")

	// log.Println(response)

	assert.Equal(t, activityId, response.Preferences[0].ID, "The activity ids should match")
}

// preferences.DELETE("", router.preferencesController.DeletePreference)

// func TestDeletePreference(t *testing.T) {
// 	router := router.NewRouter()
// 	r :=  router.SetupRoutes()

// 	userID := int32(8)

// 	token, err := service.NewTestToken(userID)
// 	assert.Equal(t, err, nil, "Error should be nil")

// 	activityId := int32(1)

// 	addPreference(t, r, activityId, token)

// 	req := httptest.NewRequest("DELETE", "/preferences", nil)
// 	req.Header.Add("Authorization", "Bearer "+token)
// 	w := httptest.NewRecorder()
// 	r.ServeHTTP(w, req)

// 	var response controller.PaginatedPreferences

// 	err = json.Unmarshal(w.Body.Bytes(), &response)
		
// 	assert.Equal(t, err, nil, "Error should be nil")

// 	// log.Println(response)

// 	assert.Equal(t, activityId, response.Preferences[0].ID, "The activity ids should match")
// }


// preferences.POST("/batch", router.preferencesController.BatchAddUserPreferences)     