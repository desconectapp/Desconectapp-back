package controller

import (
	repository "gin/db/generated"
	"gin/service"
	"net/http"
	"strconv"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
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

func NewActivitesController(conn *pgx.Conn) *ActivitiesController {
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

	activities, err := c.service.ListActivitiesRequests(activityParams)

	if err != nil {
		ctx.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic})
		return
	}
	ctx.JSON(http.StatusOK, activities)
}

func (c *ActivitiesController) CreateActivityRequest(ctx *gin.Context) {
	var input CreateActivityRequestInput
	if err := ctx.ShouldBind(&input); err != nil {
		ctx.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic})
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
