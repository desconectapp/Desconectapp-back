package controller

import (
	repository "gin/db/generated"
	"gin/service"
	"net/http"
	"strconv"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CreateActivityRequestInput represents the input structure for creating activity requests
type CreateActivityRequestInput struct {
	UserID             *int32   `json:"user_id"`
	ActivityID         *int32   `json:"activity_id"`
	Description        *string  `json:"description"`
	Longitude          *float64 `json:"longitude"`
	Latitude           *float64 `json:"latitude"`
	SearchRadius       *int32   `json:"search_radius"`
	MaxParticipants    *int32   `json:"max_participants"`
	ParticipantsNeeded *int32   `json:"participants_needed"`
	Timeslots          []uint16 `json:"timeslots"`
}

// convertTimeslotsToInt32 converts uint16 timeslots to int32 for database storage
func convertTimeslotsToInt32(timeslots []uint16) []int32 {
	result := make([]int32, len(timeslots))
	for i, slot := range timeslots {
		result[i] = int32(slot)
	}
	return result
}

type ActivitiesController struct {
	service *service.ActivitiesRequestService
}

func NewActivitesController(conn *pgxpool.Pool) *ActivitiesController {
	service := service.NewActivitiesRequestService(conn)
	return &ActivitiesController{
		service: service,
	}
}

func (c *ActivitiesController) ListActivitiesRequests(ctx *gin.Context) {
	var activityParams repository.ListActivityRequestsParams

	limit, offset := GetLimmitAndOffset(ctx)

	activityParams.Limit = int32(limit)
	activityParams.Offset = int32(offset)

	// Get user ID from JWT token
	userToken, exists := ctx.Get("userID")
	if !exists {
		ErrorWithStatus(ctx, "User not authenticated", http.StatusUnauthorized)
		return
	}
	userID := userToken.(int32)
	activityParams.UserID = &userID

	activities, err := c.service.ListActivitiesRequests(activityParams)

	if err != nil {
		ctx.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic})
		return
	}
	ctx.JSON(http.StatusOK, activities)
}

func validateActivityRequest(ctx *gin.Context, activityParams repository.CreateActivityRequestParams) bool {
	log.Printf("Validating activity request: %+v", activityParams)

	if len(activityParams.WeekTimeslots) == 0 {
		log.Printf("Validation failed: Timeslots cannot be empty")
		ErrorWithStatus(ctx, "Timeslots cannot be empty", http.StatusBadRequest)
		return false
	}

	if activityParams.ParticipantsNeeded == nil || *activityParams.ParticipantsNeeded < 2 {
		ErrorWithStatus(ctx, "Participants needed must be greater than 1", http.StatusBadRequest)
		return false
	}

	if activityParams.MaximumParticipants == nil || *activityParams.MaximumParticipants < 2 {
		ErrorWithStatus(ctx, "Maximum participants must be greater than 1", http.StatusBadRequest)
		return false
	}

	if *activityParams.ParticipantsNeeded > *activityParams.MaximumParticipants {
		ErrorWithStatus(ctx, "Participants needed must be less than maximum participants", http.StatusBadRequest)
		return false
	}

	if activityParams.Latitude == nil || activityParams.Longitude == nil {
		ErrorWithStatus(ctx, "Latitude and longitude are required", http.StatusBadRequest)
		return false
	}

	if *activityParams.Latitude < -90.0 || *activityParams.Latitude > 90.0 {
		ErrorWithStatus(ctx, "Latitude must be between -90 and 90 degrees", http.StatusBadRequest)
		return false
	}

	if *activityParams.Longitude < -180.0 || *activityParams.Longitude > 180.0 {
		ErrorWithStatus(ctx, "Longitude must be between -180 and 180 degrees", http.StatusBadRequest)
		return false
	}

	if activityParams.SearchRadius == nil || *activityParams.SearchRadius < 1 || *activityParams.SearchRadius > 100 {
		ErrorWithStatus(ctx, "Search radius must be between 1 and 100", http.StatusBadRequest)
		return false
	}

	return true
}

func (c *ActivitiesController) CreateActivityRequest(ctx *gin.Context) {
	var input CreateActivityRequestInput
	if err := ctx.ShouldBind(&input); err != nil {
		log.Printf("Error binding request: %v", err)
		ErrorWithStatus(ctx, "Could not bind", http.StatusBadRequest)
		return
	}

	log.Printf("Received activity request: %+v", input)

	// Convert timeslots to int32 for database storage
	weekTimeslots := convertTimeslotsToInt32(input.Timeslots)

	// Create the repository params
	activityParams := repository.CreateActivityRequestParams{
		UserID:              input.UserID,
		ActivityID:          input.ActivityID,
		Description:         input.Description,
		WeekTimeslots:       weekTimeslots,
		ParticipantsNeeded:  input.ParticipantsNeeded,
		MaximumParticipants: input.MaxParticipants,
		Latitude:            input.Latitude,
		Longitude:           input.Longitude,
		SearchRadius:        input.SearchRadius,
	}

	if !validateActivityRequest(ctx, activityParams) {
		return
	}

	activity, err := c.service.CreateActivityRequest(activityParams)
	if err != nil {
		ctx.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic})
		return
	}

	ctx.JSON(http.StatusOK, activity)
}

func (c *ActivitiesController) GetActivities(ctx *gin.Context) {
	var activityParams repository.GetActivitiesParams

	limit, offset := GetLimmitAndOffset(ctx)
	query := ctx.Query("q")
	activityParams.Search = &query

	activityParams.Limit = int32(limit)
	activityParams.Offset = int32(offset)

	log.Printf("Limit = %d, Offset = %d", activityParams.Limit, activityParams.Offset)

	activities, err := c.service.GetActivities(activityParams)

	if err != nil {
		ctx.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic})
		return
	}

	ctx.JSON(http.StatusOK, activities)

}

func (c *ActivitiesController) DeleteActivityRequest(ctx *gin.Context) {
	requestId := ctx.Param("requestId")
	requestIdInt, err := strconv.Atoi(requestId)
	if err != nil {
		ctx.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic})
		return
	}
	err = c.service.DeleteActivityRequest(requestIdInt)
	if err != nil {
		ctx.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"deleted": requestId,
	})
}
