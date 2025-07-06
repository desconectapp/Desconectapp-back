package service

import (
	"context"
	repository "gin/db/generated"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type GroupWithMembers struct {
	ID          string             `json:"id"`
	Name        *string            `json:"name"`
	Description *string            `json:"description"`
	CreatedAt   pgtype.Timestamptz `json:"created_at"`
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

func (s *GroupsService) CreateGroup(groupParams repository.CreateGroupParams, groupMembersInfo repository.BatchAddUserToGroupParams) (int32, error) {
	id, err := s.queries.CreateGroup(s.ctx, groupParams)
	if err != nil {
		return -1, err
	}
	groupMembersInfo.GroupID = id
	err = s.queries.BatchAddUserToGroup(s.ctx, groupMembersInfo)
	if err != nil {
		return -1, err
	}
	return id, nil
}


func (s *GroupsService) ListGroups(params repository.ListGroupsParams) ([]repository.ListGroupsRow, error) {
	groupsList, err := s.queries.ListGroups(s.ctx, params)
	if err != nil {
		return nil, err
	}
	return groupsList, nil
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
        CreatedAt:   group.CreatedAt,
        Location:    group.Location,
        Icon:        group.Icon,
        Members:     members,
    }
}