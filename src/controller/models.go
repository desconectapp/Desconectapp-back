package controller

import (
	gen "gin/db/generated"
	s "gin/service"
	"time"

	"github.com/jackc/pgx/pgtype"
)

type PaginatedUsers struct {
	Users   []gen.Profile `json:"users"`
	HasMore bool          `json:"has_more"`
}

type PaginatedGroups struct {
	Groups  []gen.ListGroupsRow `json:"groups"`
	HasMore bool                `json:"has_more"`
}

type PaginatedUserGroup struct {
	Groups  []gen.ListUserGroupsRow `json:"groups"`
	HasMore bool                    `json:"has_more"`
}

type PaginatedPreferences struct {
	Preferences []gen.GetUserPreferencesRow `json:"preferences"`
	HasMore     bool                        `json:"has_more"`
}

type PaginatedOpenGroup struct {
	Groups  []s.OpenGroup `json:"groups"`
	HasMore bool          `json:"has_more"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type AuthResponse struct {
	UserId           int32     `json:"user_id"`
	UserUuid         string    `json:"user_uuid"`
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
	DeletedUserID int32 `json:"deleted"`
}

type UserProfileCreatedResponse struct {
	ProfileUserID int32 `json:"updated"`
}

type ActivityIdResponse struct {
	ActivityPreferenseID int32 `json:"activity"`
}

type GroupIdResponse struct {
	GroupId int32 `json:"group_id"`
}

type NewGroup struct {
	ID           int32   `json:"id"`
	Name         *string `json:"name"`
	Description  *string `json:"description"`
	Location     *string `json:"location"`
	ActivityID   int32   `json:"activity_id"`
	Members      []int32 `json:"members"`
	ActivityName string  `json:"activity_name"`
	ActivityIcon *string `json:"activity_icon"`
	Status       *bool   `json:"status"`
}

type ActivityIdBatchResponse struct {
	ActivityIdBatchIDs []int32 `json:"activities"`
}

type Profile struct {
	UserID           int32            `json:"user_id"`
	Name             string           `json:"name"`
	CreatedAt        pgtype.Timestamp `json:"created_at"`
	Age              int32            `json:"age"`
	City             string           `json:"city"`
	CurrentSituation string           `json:"current_situation"`
	Gender           string           `json:"gender"`
}
