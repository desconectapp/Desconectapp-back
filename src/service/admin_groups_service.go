package service

import (
	"context"
	repository "gin/db/generated"

	"github.com/jackc/pgx/v5"
)

type AdminGroupService struct {
	queries *repository.Queries
	ctx     context.Context
}

type AdminGroup struct {
	ID          int32   `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	Location    *string `json:"location,omitempty"`
	ActivityID  int32   `json:"activity_id"`
	CreatedAt   string  `json:"created_at"`
	MemberCount int64   `json:"member_count"`
}

type AdminGroupMember struct {
	ID    int32   `json:"id"`
	Name  string  `json:"name"`
	Email *string `json:"email,omitempty"`
}

func NewAdminGroupService(conn *pgx.Conn) *AdminGroupService {
	queries := repository.New(conn)
	ctx := context.Background()

	return &AdminGroupService{
		queries: queries,
		ctx:     ctx,
	}
}

func (s *AdminGroupService) ListGroups(
	limit, offset int32,
	name *string,
	sortField, sortOrder string,
) ([]AdminGroup, error) {
	if sortOrder == "DESC" {
		rows, err := s.queries.AdminListGroupsDesc(s.ctx, repository.AdminListGroupsDescParams{
			Limit:  limit,
			Offset: offset,
			// Name:      name,
			// SortField: sortField,
		})
		if err != nil {
			return nil, err
		}
		groups := make([]AdminGroup, len(rows))
		for i, r := range rows {
			var name = ""
			if r.Name != nil {
				name = *r.Name
			}

			groups[i] = AdminGroup{
				ID:          r.ID,
				Name:        name,
				Description: r.Description,
				Location:    r.Location,
				ActivityID:  r.ActivityID,
				CreatedAt:   r.CreatedAt.Time.Format("2006-01-02 15:04:05"),
				MemberCount: r.MemberCount,
			}
		}
		return groups, nil
	}

	rows, err := s.queries.AdminListGroupsAsc(s.ctx, repository.AdminListGroupsAscParams{
		Limit:  limit,
		Offset: offset,
		// Name:      name,
		// SortField: sortField,
	})
	if err != nil {
		return nil, err
	}
	groups := make([]AdminGroup, len(rows))
	for i, r := range rows {
		var name = ""
		if r.Name != nil {
			name = *r.Name
		}

		groups[i] = AdminGroup{
			ID:          r.ID,
			Name:        name,
			Description: r.Description,
			Location:    r.Location,
			ActivityID:  r.ActivityID,
			CreatedAt:   r.CreatedAt.Time.Format("2006-01-02 15:04:05"),
			MemberCount: r.MemberCount,
		}
	}
	return groups, nil
}

func (s *AdminGroupService) GetGroup(id int32) (*AdminGroup, error) {
	r, err := s.queries.AdminGetGroup(s.ctx, id)
	if err != nil {
		return nil, err
	}

	var name = ""
	if r.Name != nil {
		name = *r.Name
	}

	return &AdminGroup{
		ID:          r.ID,
		Name:        name,
		Description: r.Description,
		Location:    r.Location,
		ActivityID:  r.ActivityID,
		CreatedAt:   r.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		MemberCount: r.MemberCount,
	}, nil
}

func (s *AdminGroupService) CreateGroup(name string, description *string, location *string, activityID int32) (*AdminGroup, error) {
	r, err := s.queries.AdminCreateGroup(s.ctx, repository.AdminCreateGroupParams{
		Name:        &name,
		Description: description,
		Location:    location,
		ActivityID:  activityID,
	})
	if err != nil {
		return nil, err
	}

	return &AdminGroup{
		ID:          r.ID,
		Name:        name,
		Description: r.Description,
		Location:    r.Location,
		ActivityID:  r.ActivityID,
		CreatedAt:   r.CreatedAt.Time.Format("2006-01-02 15:04:05"),
	}, nil
}

func (s *AdminGroupService) UpdateGroup(id int32, name string, description *string, location *string, activityID int32) (*repository.Group, error) {
	r, err := s.queries.AdminUpdateGroup(s.ctx, repository.AdminUpdateGroupParams{
		ID:          id,
		Name:        &name,
		Description: description,
		Location:    location,
		ActivityID:  activityID,
	})

	if err != nil {
		return nil, err
	}

	return &r, nil
}

func (s *AdminGroupService) CountGroups() (int64, error) {
	count, err := s.queries.AdminCountGroups(s.ctx)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *AdminGroupService) DeleteGroup(id int32) error {
	return s.queries.AdminDeleteGroup(s.ctx, id)
}

func (s *AdminGroupService) ListGroupMembers(groupID int32) ([]AdminGroupMember, error) {
	println("ListGroupMembers called with groupID:", groupID)
	rows, err := s.queries.AdminListGroupMembers(s.ctx, groupID)
	if err != nil {
		return nil, err
	}
	out := make([]AdminGroupMember, len(rows))
	for i, r := range rows {
		out[i] = AdminGroupMember{
			ID:    r.ID,
			Name:  r.Name,
			Email: &r.Email,
		}
	}
	return out, nil
}

func (s *AdminGroupService) AddGroupMember(groupID, userID int32) error {
	return s.queries.AdminAddGroupMember(s.ctx, repository.AdminAddGroupMemberParams{
		GroupID: groupID,
		UserID:  userID,
	})
}

func (s *AdminGroupService) RemoveGroupMember(groupID, userID int32) error {
	return s.queries.AdminRemoveGroupMember(s.ctx, repository.AdminRemoveGroupMemberParams{
		GroupID: groupID,
		UserID:  userID,
	})
}
