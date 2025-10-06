package service

import (
	"context"
	repository "gin/db/generated"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type GroupWithMembers struct {
	ID          int32                           `json:"id"`
	Name        *string                         `json:"name"`
	Description *string                         `json:"description"`
	Activity    string                          `json:"activity"`
	Icon        *string                         `json:"icon"`
	Location    *string                         `json:"location"`
	AvatarUrl   *string                         `json:"avatar_url"`
	Members     []repository.GetGroupMembersRow `json:"members"`
	Public      bool                            `json:"public"`
	Time        string                          `json:"created_at"`
}

type OpenGroup struct {
	ID           int32   `json:"id"`
	Name         *string `json:"name"`
	Description  *string `json:"description"`
	Location     *string `json:"location"`
	ActivityName string  `json:"activity_name"`
	MemberCount  int32   `json:"member_count"`
	Photo        string  `json:"photo"`
	Time         string  `json:"created_at"`
}

type ActivityFilter struct {
	ActivityId int32 `json:"activity_id"`
	Limit      int32 `json:"limit"`
	Offset     int32 `json:"offset"`
}

type GroupsService struct {
	queries *repository.Queries
	ctx     context.Context
}

func NewGroupsService(conn *pgxpool.Pool) *GroupsService {
	queries := repository.New(conn)
	ctx := context.Background()

	return &GroupsService{
		queries: queries,
		ctx:     ctx,
	}
}

func (s *GroupsService) CreateGroup(groupParams repository.CreateGroupParams) (repository.CreateGroupRow, error) {
	group, err := s.queries.CreateGroup(s.ctx, groupParams)
	if err != nil {
		return repository.CreateGroupRow{}, err
	}
	return group, nil
}

func (s *GroupsService) ListGroups(params repository.ListGroupsParams) ([]repository.ListGroupsRow, error) {
	groupsList, err := s.queries.ListGroups(s.ctx, params)
	if err != nil {
		return nil, err
	}
	return groupsList, nil
}

func (s *GroupsService) ListUserGroups(params repository.ListUserGroupsParams) ([]repository.ListUserGroupsRow, error) {
	groupsList, err := s.queries.ListUserGroups(s.ctx, params)
	if err != nil {
		return nil, err
	}
	return groupsList, nil
}

func (s *GroupsService) ExitGroup(exitParams repository.ExitGroupParams) error {
	err := s.queries.ExitGroup(s.ctx, exitParams)
	return err
}

func (s *GroupsService) JoinGroup(joinParams repository.AddUserToGroupParams) error {
	err := s.queries.AddUserToGroup(s.ctx, joinParams)
	return err
}

func (s *GroupsService) GetGroup(groupId int32) (GroupWithMembers, error) {
	var groupWithMembers GroupWithMembers

	group, err := s.queries.GetGroup(s.ctx, groupId)

	if err != nil {
		return groupWithMembers, err
	}

	members, err := s.queries.GetGroupMembers(s.ctx, groupId)

	if err != nil {
		return groupWithMembers, err
	}

	log.Println(*group.Public)

	return addMembers(group, members), err

}

func addMembers(group repository.GetGroupRow, members []repository.GetGroupMembersRow) GroupWithMembers {
	return GroupWithMembers{
		ID:          group.ID,
		Name:        group.Name,
		Activity:    group.Activity,
		Description: group.Description,
		Location:    group.Location,
		Icon:        group.Icon,
		Members:     members,
		Public:      *group.Public,
		Time:        group.CreatedAt.Time.String(),
		AvatarUrl:   group.AvatarUrl,
	}
}

func (s *GroupsService) DeleteGroup(id int32) (int32, error) {

	groupId, err := s.queries.DeleteGroup(s.ctx, id)

	if err != nil {
		return -1, err
	}

	return groupId, nil
}

func (s *GroupsService) UpdateGroupDescription(params repository.UpdateGroupDescriptiomParams) error {

	err := s.queries.UpdateGroupDescriptiom(s.ctx, params)

	return err
}

func (s *GroupsService) ChangeGroupPublic(params repository.ChangeGroupPublicParams) error {

	err := s.queries.ChangeGroupPublic(s.ctx, params)

	return err
}

func (s *GroupsService) ChangeGroupName(params repository.ChangeGroupNameParams) error {

	err := s.queries.ChangeGroupName(s.ctx, params)

	return err
}

func (s *GroupsService) ChangeGroupLocation(params repository.ChangeGroupLocationParams) error {

	err := s.queries.ChangeGroupLocation(s.ctx, params)

	return err
}

func (s *GroupsService) UpdateGroupAvatar(params repository.UpdateGroupAvatarParams) error {
	err := s.queries.UpdateGroupAvatar(s.ctx, params)
	return err
}

func (s *GroupsService) GetOpenGroups(filter ActivityFilter) ([]OpenGroup, error) {
	var openGroups []OpenGroup
	var err error

	if filter.ActivityId == 0 {
		openGroups, err = s.GetOpenGroupsNoFilter(
			repository.GetOpenGroupsNoFilterParams{
				Limit:  filter.Limit,
				Offset: filter.Offset,
			},
		)
	} else {
		openGroups, err = s.GetPublicOpenGroupsWithFilter(
			repository.GetOpenGroupsWithFilterParams{
				Limit:      filter.Limit,
				Offset:     filter.Offset,
				ActivityID: &filter.ActivityId,
			},
		)
	}

	return openGroups, err
}

func (s *GroupsService) GetOpenGroupsNoFilter(filter repository.GetOpenGroupsNoFilterParams) ([]OpenGroup, error) {
	groups, err := s.queries.GetOpenGroupsNoFilter(s.ctx, filter)

	if err != nil {
		return nil, err
	}

	log.Println(groups)

	var openGroups []OpenGroup

	for _, group := range groups {
		openGroups = append(openGroups, OpenGroup{
			ID:           group.ID,
			Name:         group.Name,
			Location:     group.Location,
			Description:  group.Description,
			ActivityName: group.ActivityName,
			MemberCount:  int32(group.MemberCount),
			Photo:        *group.Icon,
		})
	}
	return openGroups, err
}

func (s *GroupsService) GetPublicOpenGroupsWithFilter(filter repository.GetOpenGroupsWithFilterParams) ([]OpenGroup, error) {
	groups, err := s.queries.GetOpenGroupsWithFilter(s.ctx, filter)

	if err != nil {
		return nil, err
	}

	var openGroups []OpenGroup

	for _, group := range groups {
		openGroups = append(openGroups, OpenGroup{
			ID:           group.ID,
			Name:         group.Name,
			Location:     group.Location,
			Description:  group.Description,
			ActivityName: group.ActivityName,
			MemberCount:  int32(group.MemberCount),
			Photo:        *group.Icon,
		})
	}
	return openGroups, err
}

func (s *GroupsService) GetUserRecommendations(filter repository.GetPreferredGroupsParams) ([]OpenGroup, error) {
	groups, err := s.queries.GetPreferredGroups(s.ctx, filter)

	if err != nil {
		return nil, err
	}

	var openGroups []OpenGroup

	for _, group := range groups {
		openGroups = append(openGroups, OpenGroup{
			ID:           group.ID,
			Name:         group.Name,
			Location:     group.Location,
			Description:  group.Description,
			ActivityName: group.ActivityName,
			MemberCount:  int32(group.MemberCount),
			Photo:        *group.ActivityIcon,
			Time:         group.CreatedAt.Time.String(),
		})
	}
	return openGroups, err
}

