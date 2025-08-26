package test

import "github.com/jackc/pgx/v5/pgtype"

type AuthBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateProfile struct {
	Name             string `json:"name"`
	Age              int32  `json:"age"`
	City             string `json:"city"`
	CurrentSituation string `json:"current_situation"`
	Gender           string `json:"gender"`
}

type AddPreference struct {
	ActivityID int32 `json:"activity_id"`
}

type TimeSlot struct {
	Start int `json:"start"`
	End   int `json:"end"`
}
type Schedules map[string][]TimeSlot

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

type ActivityRequest struct {
	ID                  int32            `json:"id"`
	UserID              *int32           `json:"user_id"`
	ActivityID          *int32           `json:"activity_id"`
	Description         *string          `json:"description"`
	WeekHours           []int32          `json:"week_hours"`
	ParticipantsNeeded  *int32           `json:"participants_needed"`
	MaximumParticipants *int32           `json:"maximum_participants"`
	Latitude            *float64         `json:"latitude"`
	Longitude           *float64         `json:"longitude"`
	SearchRadius        *int32           `json:"search_radius"`
	CreatedAt           pgtype.Timestamp `json:"created_at"`
	ExpiresAt           pgtype.Timestamp `json:"expires_at"`
}
