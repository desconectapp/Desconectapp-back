package test

import (
	"encoding/json"
	"gin/controller"
	"gin/router"
	"gin/service"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)


// groups.GET("", router.groupsController.ListGroups)
func TestListGroups(t *testing.T){
router := router.NewRouter()
	r :=  router.SetupRoutes()
	
	user1, token := NewUserWithProfile(t, r, "bokuto", CreateProfile{Name: "bokuto", Age: 17, City: "Tokyo", CurrentSituation: "WORKING", Gender: "Male"})
	user2, _ := NewUserWithProfile(t, r, "akaashi", CreateProfile{Name: "akaashi", Age: 16, City: "Tokyo", CurrentSituation: "WORKING", Gender: "Male"})
	user3, _ := NewUserWithProfile(t, r, "konoha", CreateProfile{Name: "konoha", Age: 16, City: "Tokyo", CurrentSituation: "WORKING", Gender: "Male"})
		
	activityId := int32(4)

	members := []int32{user1, user2, user3}

	name := "Fukurodani"
	location := "Tokyo"
	
	NewGroup(t, r, name, location, activityId, members, token)

	limit := "6"
	offset := "0"

	req := httptest.NewRequest("GET", "/groups?limit="+limit+"&offset="+offset, nil)
	req.Header.Add("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response controller.PaginatedGroups

	err := json.Unmarshal(w.Body.Bytes(), &response)
	
	assert.Equal(t, err, nil, "Error should be nil")

	assert.LessOrEqual(t, strconv.Itoa(len(response.Groups)), limit)
	assert.Equal(t, false, response.HasMore)

	assert.Equal(t, len(members), int(response.Groups[0].MembersCount))
	assert.Equal(t, location, *response.Groups[0].Location)
	assert.Equal(t, name, *response.Groups[0].Name)
}

// groups.POST("", router.groupsController.CreateGroup)
func TestPostGroup(t *testing.T) {
	router := router.NewRouter()
	r :=  router.SetupRoutes()
	
	user1, token := NewUser(t, r, "tsukki")
	user2, _ := NewUser(t, r, "hinata")
	user3, _ := NewUser(t, r, "kageyama")
		
	activityId := int32(4)

	members := []int32{user1, user2, user3}
	
	NewGroup(t, r, "Karasuno", "Miyagi", activityId, members, token)
}

// groups.GET("/:groupId", router.groupsController.GetGroup)
func TestGetGroup(t *testing.T) {
	router := router.NewRouter()
	r :=  router.SetupRoutes()
	
	user1, token := NewUserWithProfile(t, r, "kuroo", CreateProfile{Name: "kuroo", Age: 17, City: "Tokyo", CurrentSituation: "WORKING", Gender: "Male"})
	user2, _ := NewUserWithProfile(t, r, "kenma", CreateProfile{Name: "kenma", Age: 16, City: "Tokyo", CurrentSituation: "WORKING", Gender: "Male"})
	user3, _ := NewUserWithProfile(t, r, "lev", CreateProfile{Name: "lev", Age: 15, City: "Tokyo", CurrentSituation: "WORKING", Gender: "Male"})
		
	activityId := int32(4)

	members := []int32{user1, user2, user3}

	name := "Nekoma"
	location := "Tokyo"
	
	group := NewGroup(t, r, name, location, activityId, members, token)

	req := httptest.NewRequest("GET", "/groups/"+strconv.Itoa(int(group)), nil)
	req.Header.Add("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response service.GroupWithMembers
	err := json.Unmarshal(w.Body.Bytes(), &response)
	
	assert.Equal(t, err, nil, "Error should be nil")
	assert.Equal(t, name, *response.Name)
	assert.Equal(t, location, *response.Location)

	assert.NotEqual(t, 0, len(response.Members))
}

// groups.DELETE("/:groupId", router.groupsController.DeleteGroup)
func TestDeleteGroup(t *testing.T) {
	router := router.NewRouter()
	r :=  router.SetupRoutes()
	
	user1, token := NewUser(t, r, "suga")
	user2, _ := NewUser(t, r, "asahi")
	user3, _ := NewUser(t, r, "daichi")
		
	activityId := int32(4)

	members := []int32{user1, user2, user3}
	
	group := NewGroup(t, r, "Karasuno", "Miyagi", activityId, members, token)

	req := httptest.NewRequest("DELETE", "/groups/"+strconv.Itoa(int(group)), nil)
	req.Header.Add("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response controller.GroupIdResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, err, nil, "Error should be nil")
	assert.Equal(t, group, response.GroupId)
}

// groups.GET("/user", router.groupsController.ListUserGroups)

func TestUserGroups(t *testing.T){
router := router.NewRouter()
	r :=  router.SetupRoutes()
	
	user1, token := NewUserWithProfile(t, r, "oikawa", CreateProfile{Name: "oikawa", Age: 17, City: "Miyagi", CurrentSituation: "WORKING", Gender: "Male"})
	user2, _ := NewUserWithProfile(t, r, "iwaizumi", CreateProfile{Name: "iwaizumi", Age: 17, City: "Miyagi", CurrentSituation: "WORKING", Gender: "Male"})
	user3, _ := NewUserWithProfile(t, r, "kunimi", CreateProfile{Name: "kunimi", Age: 15, City: "Miyagi", CurrentSituation: "WORKING", Gender: "Male"})
		
	activityId := int32(4)

	members := []int32{user1, user2, user3}

	name := "Aoba Johsai"
	location := "Miyagi"
	
	NewGroup(t, r, name, location, activityId, members, token)

	user5, _:= NewUserWithProfile(t, r, "juan", CreateProfile{Name: "juan", Age: 26, City: "Buenos Airres", CurrentSituation: "WORKING", Gender: "Male"})
	user6, _ := NewUserWithProfile(t, r, "nacho", CreateProfile{Name: "nacho", Age: 25, City: "Buenos Airres", CurrentSituation: "WORKING", Gender: "Male"})
	user7, _ := NewUserWithProfile(t, r, "miguel", CreateProfile{Name: "miguel", Age: 25, City: "Buenos Airres", CurrentSituation: "WORKING", Gender: "Male"})
	
	members2 := []int32{user1, user5, user6, user7}

	name2 := "Boca Juniors"
	location2 := "Buenos Aires"
	
	NewGroup(t, r, name2, location2, activityId, members2, token)

	limit := "6"
	offset := "0"

	req := httptest.NewRequest("GET", "/groups/user?limit="+limit+"&offset="+offset, nil)
	req.Header.Add("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response controller.PaginatedUserGroup

	err := json.Unmarshal(w.Body.Bytes(), &response)
	
	assert.Equal(t, err, nil, "Error should be nil")

	assert.LessOrEqual(t, strconv.Itoa(len(response.Groups)), limit)
	assert.Equal(t, false, response.HasMore)

	assert.Equal(t, len(members2), int(response.Groups[0].MembersCount))
	assert.Equal(t, location2, *response.Groups[0].Location)
	assert.Equal(t, name2, *response.Groups[0].Name)
	assert.Equal(t, len(members), int(response.Groups[1].MembersCount))
	assert.Equal(t, location, *response.Groups[1].Location)
	assert.Equal(t, name, *response.Groups[1].Name)
}

// groups.DELETE("/user-from-group/:groupId", router.groupsController.ExitGroup)
func TestDeleteUserFromGroup(t *testing.T) {
	router := router.NewRouter()
	r :=  router.SetupRoutes()
	
	user1, token := NewUserWithProfile(t, r, "ushijima", CreateProfile{Name: "ushijima", Age: 17, City: "Miyagi", CurrentSituation: "WORKING", Gender: "Male"})
	user2, _ := NewUserWithProfile(t, r, "tendo", CreateProfile{Name: "tendo", Age: 16, City: "Miyagi", CurrentSituation: "WORKING", Gender: "Male"})
	user3, _ := NewUserWithProfile(t, r, "semi", CreateProfile{Name: "semi", Age: 15, City: "Miyagi", CurrentSituation: "WORKING", Gender: "Male"})
		
	activityId := int32(4)

	members := []int32{user1, user2, user3}
	
	group := NewGroup(t, r, "Shiratorizawa", "Miyagi", activityId, members, token)

	req := httptest.NewRequest("DELETE", "/groups/user-from-group/"+strconv.Itoa(int(group)), nil)
	req.Header.Add("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response controller.GroupIdResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, err, nil, "Error should be nil")
	assert.Equal(t, group, response.GroupId)
}
