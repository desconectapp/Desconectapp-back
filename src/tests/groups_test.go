package test

import (
	"bytes"
	"encoding/json"
	"gin/controller"
	"gin/router"
	"gin/service"
	"log"
	"net/http/httptest"
	"strconv"
	"testing"
	"net/http"
	"github.com/stretchr/testify/assert"
)

// groups.GET("", router.groupsController.ListGroups)
func TestListGroups(t *testing.T){
router := router.NewRouter()
	r :=  router.SetupRoutes()
	
	token, err := service.NewTestToken(1)
	assert.Equal(t, err, nil, "Error should be nil")

	name := "Karasuno"
	location := "Miyagi"
	members := 4
	
	limit := "9"
	offset := "0"

	req := httptest.NewRequest("GET", "/groups?limit="+limit+"&offset="+offset, nil)
	req.Header.Add("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response controller.PaginatedGroups

	err = json.Unmarshal(w.Body.Bytes(), &response)
	
	assert.Equal(t, err, nil, "Error should be nil")
	log.Println(response)

	assert.LessOrEqual(t, strconv.Itoa(len(response.Groups)), limit)
	assert.Equal(t, false, response.HasMore)

	assert.Equal(t, members, int(response.Groups[0].MembersCount))
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
	
	NewGroup(t, r, "Karasuno", "Miyagi", activityId, members, false, token)
}

// groups.GET("/:groupId", router.groupsController.GetGroup)
func TestGetGroup(t *testing.T) {
	router := router.NewRouter()
	r :=  router.SetupRoutes()
	
	token, err := service.NewTestToken(1)
	assert.Equal(t, err, nil, "Error should be nil")

	groupID := 2
	name := "Nekoma"
	location := "Tokyo"
	members := 2
	

	req := httptest.NewRequest("GET", "/groups/"+strconv.Itoa(groupID), nil)
	req.Header.Add("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response service.GroupWithMembers
	err = json.Unmarshal(w.Body.Bytes(), &response)

	log.Println(response.Members)
	
	assert.Equal(t, err, nil, "Error should be nil")
	assert.Equal(t, name, *response.Name)
	assert.Equal(t, location, *response.Location)
	assert.Equal(t, members, len(response.Members))
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
	
	group := NewGroup(t, r, "Karasuno", "Miyagi", activityId, members, false, token)

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
	
	token, err := service.NewTestToken(23)
	assert.Equal(t, err, nil, "Error should be nil")
	
	name := "Aoba Johsai"
	location := "Miyagi"
	members := 3
	members2 := 4
	name2 := "Argentina"
	location2 := "Buenos Aires"
	
	limit := "6"
	offset := "0"

	req := httptest.NewRequest("GET", "/groups/user?limit="+limit+"&offset="+offset, nil)
	req.Header.Add("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response controller.PaginatedUserGroup
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, err, nil, "Error should be nil")

	assert.LessOrEqual(t, strconv.Itoa(len(response.Groups)), limit)
	assert.Equal(t, false, response.HasMore)

	assert.Equal(t, members2, int(response.Groups[1].MembersCount))
	assert.Equal(t, location2, *response.Groups[1].Location)
	assert.Equal(t, name2, *response.Groups[1].Name)
	assert.Equal(t, members, int(response.Groups[0].MembersCount))
	assert.Equal(t, location, *response.Groups[0].Location)
	assert.Equal(t, name, *response.Groups[0].Name)
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
	
	group := NewGroup(t, r, "Shiratorizawa", "Miyagi", activityId, members, false, token)

	req := httptest.NewRequest("DELETE", "/groups/user-from-group/"+strconv.Itoa(int(group)), nil)
	req.Header.Add("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response controller.GroupIdResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, err, nil, "Error should be nil")
	assert.Equal(t, group, response.GroupId)
}

func TestUpdateGroupDescription(t *testing.T) {
	router := router.NewRouter()
	r :=  router.SetupRoutes()
	
	token, err := service.NewTestToken(32)
	assert.Equal(t, err, nil, "Error should be nil")

	name := "Huntrix"
	location := "Seoul"
	desc := "New Description"
	groupID := 7

	body := NewDescription{
		NewDesc: desc,
	}
	jsonBody, err := json.Marshal(body)
	assert.Equal(t, err, nil, "Error should be nil")

	req := httptest.NewRequest("PUT", "/groups/description/"+strconv.Itoa(int(groupID)), bytes.NewReader(jsonBody))
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	req_2 := httptest.NewRequest("GET", "/groups/"+strconv.Itoa(int(groupID)), nil)
	req_2.Header.Add("Authorization", "Bearer "+token)
	w_2 := httptest.NewRecorder()
	r.ServeHTTP(w_2, req_2)

	var response service.GroupWithMembers
	err = json.Unmarshal(w_2.Body.Bytes(), &response)
	assert.Equal(t, err, nil, "Error should be nil")
	
	assert.Equal(t, name, *response.Name)
	assert.Equal(t, location, *response.Location)
	assert.Equal(t, desc, *response.Description)
	assert.Equal(t, 3, len(response.Members))
}

func TestChangeGroupStatusToTrueThenFalse(t *testing.T) {
	router := router.NewRouter()
	r := router.SetupRoutes()

	token, err := service.NewTestToken(32)
	assert.Equal(t, err, nil, "Error should be nil")

	groupID := 5
	groupName := "Shiratorizawa"
	location := "Miyagi"

	checkStatus := func(status bool) {
		body, _ := json.Marshal(NewStatus{NewStatus: status})
		req := httptest.NewRequest("PUT", "/groups/public/"+strconv.Itoa(int(groupID)), bytes.NewReader(body))
		req.Header.Add("content-type", "application/json")
		req.Header.Add("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		reqGet := httptest.NewRequest("GET", "/groups/"+strconv.Itoa(int(groupID)), nil)
		reqGet.Header.Add("Authorization", "Bearer "+token)
		wGet := httptest.NewRecorder()
		r.ServeHTTP(wGet, reqGet)

		var resp service.GroupWithMembers
		err := json.Unmarshal(wGet.Body.Bytes(), &resp)

		log.Println(resp.Members)

		assert.NoError(t, err)
		assert.Equal(t, groupName, *resp.Name)
		assert.Equal(t, location, *resp.Location)
		assert.Equal(t, status, resp.Public)
		assert.NotEmpty(t, resp.Members)
	}

	checkStatus(true)
	checkStatus(false)
}

func TestGetStatusOpenListNoFilter(t *testing.T) {
	router := router.NewRouter()
	r := router.SetupRoutes()

	token, err := service.NewTestToken(32)
	assert.Equal(t, err, nil, "Error should be nil")

	groupID := 5
	groupName := "Shiratorizawa"
	location := "Miyagi"

	checkStatus := func(status bool) {
		body, _ := json.Marshal(NewStatus{NewStatus: status})
		req := httptest.NewRequest("PUT", "/groups/public/"+strconv.Itoa(int(groupID)), bytes.NewReader(body))
		req.Header.Add("content-type", "application/json")
		req.Header.Add("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		reqGet := httptest.NewRequest("GET", "/groups/"+strconv.Itoa(int(groupID)), nil)
		reqGet.Header.Add("Authorization", "Bearer "+token)
		wGet := httptest.NewRecorder()
		r.ServeHTTP(wGet, reqGet)

		var resp service.GroupWithMembers
		err := json.Unmarshal(wGet.Body.Bytes(), &resp)

		log.Println(resp)

		assert.NoError(t, err)
		assert.Equal(t, groupName, *resp.Name)
		assert.Equal(t, location, *resp.Location)
		assert.Equal(t, status, resp.Public)
		assert.NotEmpty(t, resp.Members)
	}

	checkStatus(true)

	limit := "6"
	offset := "0"
	
	body := ActivityIdStruct{
		ActivityID: 0,
	}
	jsonBody, err := json.Marshal(body)
	assert.Equal(t, err, nil, "Error should be nil")

	req := httptest.NewRequest("GET", "/groups/open?limit="+limit+"&offset="+offset, bytes.NewReader(jsonBody))
	req.Header.Add("content-type", "application/json")	
	req.Header.Add("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response controller.PaginatedOpenGroup

	log.Println(response)

	err = json.Unmarshal(w.Body.Bytes(), &response)
	
	assert.Equal(t, err, nil, "Error should be nil")

	assert.LessOrEqual(t, strconv.Itoa(len(response.Groups)), limit)
	assert.Equal(t, false, response.HasMore)

	assert.Equal(t, location, *response.Groups[0].Location)
	assert.Equal(t, groupName, *response.Groups[0].Name)
	assert.Equal(t, int32(3), response.Groups[0].MemberCount)

	checkStatus(false)
}

func TestGetStatusOpenListWithFilter(t *testing.T) {
	router := router.NewRouter()
	r := router.SetupRoutes()

	token, err := service.NewTestToken(32)
	assert.Equal(t, err, nil, "Error should be nil")

	group := 5
	groupName := "Shiratorizawa"
	location := "Miyagi"

	checkStatus := func(status bool) {
		body, _ := json.Marshal(NewStatus{NewStatus: status})
		req := httptest.NewRequest("PUT", "/groups/public/"+strconv.Itoa(int(group)), bytes.NewReader(body))
		req.Header.Add("content-type", "application/json")
		req.Header.Add("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		reqGet := httptest.NewRequest("GET", "/groups/"+strconv.Itoa(int(group)), nil)
		reqGet.Header.Add("Authorization", "Bearer "+token)
		wGet := httptest.NewRecorder()
		r.ServeHTTP(wGet, reqGet)

		var resp service.GroupWithMembers
		err := json.Unmarshal(wGet.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, groupName, *resp.Name)
		assert.Equal(t, location, *resp.Location)
		assert.Equal(t, status, resp.Public)
		assert.NotEmpty(t, resp.Members)
	}

	checkStatus(true)

	limit := "6"
	offset := "0"
	
	body := ActivityIdStruct{
		ActivityID: 4,
	}
	jsonBody, err := json.Marshal(body)
	assert.Equal(t, err, nil, "Error should be nil")

	req := httptest.NewRequest("GET", "/groups/open?limit="+limit+"&offset="+offset, bytes.NewReader(jsonBody))
	req.Header.Add("content-type", "application/json")	
	req.Header.Add("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response controller.PaginatedOpenGroup

	log.Println(response)

	err = json.Unmarshal(w.Body.Bytes(), &response)
	
	assert.Equal(t, err, nil, "Error should be nil")

	assert.LessOrEqual(t, strconv.Itoa(len(response.Groups)), limit)
	assert.Equal(t, false, response.HasMore)

	assert.Equal(t, location, *response.Groups[0].Location)
	assert.Equal(t, groupName, *response.Groups[0].Name)
	assert.Equal(t, int32(3), response.Groups[0].MemberCount)
}

func TestCreatePublic(t *testing.T) {
	router := router.NewRouter()
	r := router.SetupRoutes()

	token, err := service.NewTestToken(32)
	assert.Equal(t, err, nil, "Error should be nil")

	body := GroupInfo{
		ActivityID:  1,
		Name:        "PublicGroup",
		Location:    "location",
		MembersIds:  []int32{1, 2, 3, 4},
		Public: true,
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
	assert.Equal(t, true, *response.Public)
}