package controller

import (
	gen "gin/db/generated"
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