package service

import (
	"context"
	repository "gin/db/generated"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CommunityWithMembers struct {
	ID           int32                     `json:"id"`
	Name         string                    `json:"name"`
	Description  string                    `json:"description"`
	Activity     string                    `json:"activity"`
	Icon         string                    `json:"icon"`
	Location     string                    `json:"location"`
	LocationName string                    `json:"locationName"`
	AvatarUrl    string                    `json:"avatarUrl"`
	Members      []repository.GetCommunityMembersRow `json:"members"`
	CreatedAt    string                    `json:"createdAt"`
}

type CommunityWithLocation struct {
	ID            int32   `json:"id"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Location      string  `json:"location"`
	LocationName  string  `json:"locationName"`
	WeekTimeslots []int32 `json:"weekTimeslots"`
	ActivityName  string  `json:"activityName"`
	MemberCount   int32   `json:"memberCount"`
	DistanceKM    float64 `json:"distanceKm"`
	AvatarUrl     string  `json:"avatarUrl"`
	Icon          string  `json:"icon"`
}

type CommunitiesService struct {
	queries *repository.Queries
	ctx     context.Context
}

func NewCommunitiesService(conn *pgxpool.Pool) *CommunitiesService {
	queries := repository.New(conn)
	ctx := context.Background()

	return &CommunitiesService{
		queries: queries,
		ctx:     ctx,
	}
}

// CreateCommunity wraps the sqlc CreateCommunity query
func (s *CommunitiesService) CreateCommunity(params repository.CreateCommunityParams) (repository.CreateCommunityRow, error) {
	return s.queries.CreateCommunity(s.ctx, params)
}

// ListUserCommunities returns communities the user belongs to
func (s *CommunitiesService) ListUserCommunities(params repository.ListUserCommunitiesParams) ([]repository.ListUserCommunitiesRow, error) {
	return s.queries.ListUserCommunities(s.ctx, params)
}

// GetCommunity returns a single community with members
func (s *CommunitiesService) GetCommunity(communityID int32) (CommunityWithMembers, error) {
	var community CommunityWithMembers

	row, err := s.queries.GetCommunity(s.ctx, communityID)
	if err != nil {
		return community, err
	}

	members, err := s.queries.GetCommunityMembers(s.ctx, communityID)
	if err != nil {
		return community, err
	}

	return mapCommunityWithMembers(row, members), nil
}

func mapCommunityWithMembers(row repository.GetCommunityRow, members []repository.GetCommunityMembersRow) CommunityWithMembers {
	return CommunityWithMembers{
		ID:           row.ID,
		Name:         *row.Name,
		Description:  *row.Description,
		Activity:     row.Activity,
		Icon:         *row.Icon,
		Location:     *row.Location,
		LocationName: *row.LocationName,
		AvatarUrl:    *row.AvatarUrl,
		Members:      members,
		CreatedAt:    row.CreatedAt.Time.String(),
	}
}

// DeleteCommunity deletes a community by ID
func (s *CommunitiesService) DeleteCommunity(id int32) (int32, error) {
	return s.queries.DeleteCommunity(s.ctx, id)
}

// UpdateCommunity updates community properties
func (s *CommunitiesService) UpdateCommunityDescription(params repository.UpdateCommunityDescriptionParams) error {
	return s.queries.UpdateCommunityDescription(s.ctx, params)
}

func (s *CommunitiesService) ChangeCommunityName(params repository.ChangeCommunityNameParams) error {
	return s.queries.ChangeCommunityName(s.ctx, params)
}

func (s *CommunitiesService) ChangeCommunityLocation(params repository.ChangeCommunityLocationParams) error {
	return s.queries.ChangeCommunityLocation(s.ctx, params)
}

func (s *CommunitiesService) UpdateCommunityAvatar(params repository.UpdateCommunityAvatarParams) error {
	return s.queries.UpdateCommunityAvatar(s.ctx, params)
}

// AddUserToCommunity adds a user
func (s *CommunitiesService) AddUserToCommunity(params repository.AddUserToCommunityParams) error {
	return s.queries.AddUserToCommunity(s.ctx, params)
}

// ExitCommunity removes a user
func (s *CommunitiesService) ExitCommunity(params repository.ExitCommunityParams) error {
	return s.queries.ExitCommunity(s.ctx, params)
}

// GetCommunitiesWithLocation for filtering by distance and activity
func (s *CommunitiesService) GetCommunitiesWithLocation(params repository.GetCommunitiesWithLocationParams) ([]CommunityWithLocation, error) {
	rows, err := s.queries.GetCommunitiesWithLocation(s.ctx, params)
	if err != nil {
		return nil, err
	}

	var result []CommunityWithLocation
	for _, r := range rows {
		result = append(result, CommunityWithLocation{
			ID:            r.ID,
			Name:          *r.Name,
			Description:   *r.Description,
			Location:      *r.Location,
			LocationName:  *r.LocationName,
			WeekTimeslots: r.WeekTimeslots,
			ActivityName:  r.ActivityName,
			MemberCount:   int32(r.MemberCount),
			DistanceKM:    r.DistanceKm,
			AvatarUrl:     *r.AvatarUrl,
			Icon:          *r.Icon,
		})
	}
	return result, nil
}
