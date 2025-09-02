package service

import (
	"context"
	repository "gin/db/generated"
	"github.com/jackc/pgx/v5"
)

type GroupWithMembers struct {
	ID          int32             `json:"id"`
	Name        *string            `json:"name"`
	Description *string            `json:"description"`
	Activity    string             `json:"activity"`
	Icon        *string            `json:"icon"`
	Location    *string            `json:"location"`
	Members 	[]repository.GetGroupMembersRow		`json:"members"`
}


type GroupsService struct {
	queries *repository.Queries
	ctx     context.Context
}

func NewGroupsService(conn *pgx.Conn) *GroupsService {
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
	// groupMembersInfo.GroupID = id
	// err = s.queries.BatchAddUserToGroup(s.ctx, groupMembersInfo)
	// if err != nil {
	// 	s.queries.DeleteGroup(s.ctx, id)
	// 	return -1, err
	// }
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

func (s *GroupsService) ExitGroup(exitParams repository.ExitGroupParams) (error) {
	err := s.queries.ExitGroup(s.ctx, exitParams)
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
    }
}

func (s *GroupsService) DeleteGroup(id int32) (int32, error) {

	groupId, err := s.queries.DeleteGroup(s.ctx, id) 

	if err != nil {
		return -1, err
	}

	return groupId, nil
}