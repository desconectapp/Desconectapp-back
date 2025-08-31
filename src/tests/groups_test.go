package test

// import (
// 	"encoding/json"
// 	controller "gin/controller"
// 	"gin/router"
// 	"gin/service"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// )

// // groups.GET("/:groupId", router.groupsController.GetGroup)

// func TestGetGroup(t *testing.T) {
// 	router := router.NewRouter()
// 	r :=  router.SetupRoutes()

// 	userID := int32(3)

// 	token, err := service.NewTestToken(userID)

// 	assert.Equal(t, err, nil, "Error should be nil")

// 	req := httptest.NewRequest("GET", "/users/user", nil)
// 	req.Header.Add("Authorization", "Bearer "+token)
// 	w := httptest.NewRecorder()
// 	r.ServeHTTP(w, req)

// 	var response controller.Profile

// 	err = json.Unmarshal(w.Body.Bytes(), &response)
// 	assert.Equal(t, err, nil, "Error should be nil")
	
// 	assert.Equal(t, w.Code, http.StatusOK, "Status code should be 200")
// 	assert.Equal(t, userID, response.UserID, "The user ids should match")
// }

// groups.GET("", router.groupsController.ListGroups)


// groups.GET("/user", router.groupsController.ListUserGroups)


// groups.POST("", router.groupsController.CreateGroup)

// func TestPostGroup(t *testing.T) {
// 	router := router.NewRouter()
// 	r :=  router.SetupRoutes()
	
// 	_, token := NewUser(t, r, "test_post_preference")
		
// 	activityId := int32(19)
	
// 	NewGroup(t, r, "skater")
// }

// groups.DELETE("/:groupId", router.groupsController.DeleteGroup)


// groups.DELETE("/user-from-group/:groupId", router.groupsController.ExitGroup)