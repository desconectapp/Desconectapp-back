package service

import (
	"context"
	repository "gin/db/generated"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CommunityWithMembers struct {
	ID           int32                     `json:"id"`
	Name         *string                    `json:"name"`
	Description  *string                    `json:"description"`
	Activity     string                    `json:"activity"`
	Icon         *string                    `json:"icon"`
	Location     *string                    `json:"location"`
	LocationName *string                    `json:"locationName"`
	AvatarUrl    *string                    `json:"avatarUrl"`
	Members      []repository.GetCommunityMembersRow `json:"members"`
	CreatedAt    string                    `json:"createdAt"`
	IsCurrentUserAdmin bool               `json:"is_current_user_admin"`
	WeekTimeslots      []int32            `json:"week_timeslots"`
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

func (s *CommunitiesService) CreateCommunity(params repository.CreateCommunityParams) (repository.CreateCommunityRow, error) {
	return s.queries.CreateCommunity(s.ctx, params)
}

func (s *CommunitiesService) ListUserCommunities(params repository.ListUserCommunitiesParams) ([]repository.ListUserCommunitiesRow, error) {
	return s.queries.ListUserCommunities(s.ctx, params)
}

func (s *CommunitiesService) GetCommunity(params repository.GetCommunityParams) (CommunityWithMembers, error) {
	var community CommunityWithMembers

	row, err := s.queries.GetCommunity(s.ctx, params)
	if err != nil {
		return community, err
	}

	members, err := s.queries.GetCommunityMembers(s.ctx, params.ID)
	if err != nil {
		return community, err
	}

	return mapCommunityWithMembers(row, members), nil
}

func mapCommunityWithMembers(row repository.GetCommunityRow, members []repository.GetCommunityMembersRow) CommunityWithMembers {
	return CommunityWithMembers{
		ID:           row.ID,
		Name:         row.Name,
		Description:  row.Description,
		Activity:     row.Activity,
		Icon:         row.Icon,
		Location:     row.Location,
		LocationName: row.LocationName,
		AvatarUrl:    row.AvatarUrl,
		Members:      members,
		CreatedAt:    row.CreatedAt.Time.String(),
		IsCurrentUserAdmin: row.IsCurrentUserAdmin,
		WeekTimeslots: row.WeekTimeslots,
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

func (s *CommunitiesService) AddUserToCommunity(params repository.AddUserToCommunityParams) error {
	return s.queries.AddUserToCommunity(s.ctx, params)
}

func (s *CommunitiesService) ExitCommunity(params repository.ExitCommunityParams) error {
	return s.queries.ExitCommunity(s.ctx, params)
}

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
