package controller

import (
	gen "gin/db/generated"
	"time"
)

type PaginatedUsers struct {
	Users	[]gen.Profile `json:"users"`
	HasMore	bool	`json:"has_more"`
}

type PaginatedGroups struct {
	Groups []gen.ListGroupsRow `json:"groups"`
	HasMore	bool	`json:"has_more"`
}

type PaginatedMembers struct {
	Members []gen.ListUserGroupsRow `json:"members"`
	HasMore bool	`json:"has_more"`
}

type PaginatedPreferences struct {
	Preferences []gen.GetUserPreferencesRow `json:"preferences"`
	HasMore bool	`json:"has_more"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type AuthResponse struct {
	UserId			int32		`json:"user_id"`
	Token            string    `json:"token"`
	RefreshToken     string    `json:"refresh_token"`
	ExpiresAt        time.Time `json:"expires_at"`
	RefreshExpiresAt time.Time `json:"refresh_expires_at"`
}

type UserUpdateInfo struct {
	Name             string  `json:"name"`
	Email            string  `json:"email"`
	Age              int32   `json:"age"`
	City             string  `json:"city"`
	CurrentSituation string  `json:"current_situation"`
	ActivityIDs      []int32 `json:"activity_ids"`
}

type UserDeletedResponse struct {
	DeletedUserID	int32	`json:"deleted"`
}

type UserProfileCreatedResponse struct {
	ProfileUserID	int32	`json:"updated"`
}

type ActivityIdResponse struct {
	ActivityPreferenseID	int32	`json:"activity"`
}