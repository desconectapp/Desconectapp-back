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

// TimeSlot represents a time range with start and end hours
type TimeSlot struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

// Schedules represents the weekly schedule structure
type Schedules map[string][]TimeSlot

// CreateActivityRequestInput represents the input structure for creating activity requests
type CreateActivityRequestInput struct {
	UserID             *int32    `json:"user_id"`
	ActivityID         *int32    `json:"activity_id"`
	Description        *string   `json:"description"`
	Longitude          *float64  `json:"longitude"`
	Latitude           *float64  `json:"latitude"`
	SearchRadius       *int32    `json:"search_radius"`
	MaxParticipants    *int32    `json:"max_participants"`
	ParticipantsNeeded *int32    `json:"participants_needed"`
	Schedules          Schedules `json:"schedules"`
}

// parseSchedulesToWeekHours converts the schedules JSON to week_hours array
// Monday starts at 0, Tuesday at 24, Wednesday at 48, etc.
// For each day, hours are added as: day_offset + hour
func parseSchedulesToWeekHours(schedules Schedules) []int32 {

	dayOffsets := map[string]int32{
		"monday":    0,
		"tuesday":   24,
		"wednesday": 48,
		"thursday":  72,
		"friday":    96,
		"saturday":  120,
		"sunday":    144,
	}

	var weekHours []int32

	for day, timeSlots := range schedules {
		if offset, exists := dayOffsets[day]; exists {
			for _, slot := range timeSlots {
				// Add all hours in the time slot range
				for hour := slot.Start; hour < slot.End; hour++ {
					weekHours = append(weekHours, offset+int32(hour))
				}
			}
		}
	}

	return weekHours
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
	if len(activityParams.WeekHours) == 0 {
		ErrorWithStatus(ctx, "Schedules cannot be empty", http.StatusBadRequest)
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
		ErrorWithStatus(ctx, "Could not bind", http.StatusBadRequest)
		return
	}

	// Parse schedules to week_hours
	weekHours := parseSchedulesToWeekHours(input.Schedules)

	// Create the repository params
	activityParams := repository.CreateActivityRequestParams{
		UserID:              input.UserID,
		ActivityID:          input.ActivityID,
		Description:         input.Description,
		WeekHours:           weekHours,
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
