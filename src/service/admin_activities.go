package service

import (
	"context"

	repository "gin/db/generated"
	"github.com/jackc/pgx/v5"
)

type AdminActivityService struct {
	queries *repository.Queries
	ctx     context.Context
}

type AdminActivity struct {
	ID        int32  `json:"id"`
	Name      string `json:"name"`
	Icon      *string `json:"icon"`
	Category  string `json:"category"`
	CreatedAt string `json:"created_at"`
}


func toNullCategory(cat *string) repository.NullCategories {
	if cat == nil {
		return repository.NullCategories{Valid: false}
	}
	return repository.NullCategories{
		Categories: repository.Categories(*cat),
		Valid:      true,
	}
}

func NewAdminActivityService(conn *pgx.Conn) *AdminActivityService {
	queries := repository.New(conn)
	ctx := context.Background()

	return &AdminActivityService{
		queries: queries,
		ctx:     ctx,
	}
}

func (s *AdminActivityService) ListActivities(limit, offset int32, name *string, category *string) ([]AdminActivity, error) {
	rows, err := s.queries.AdminListActivities(s.ctx, repository.AdminListActivitiesParams{
		Limit:    limit,
		Offset:   offset,
		Name:     name,
		Category: toNullCategory(category),
	})
	if err != nil {
		return nil, err
	}

	activities := make([]AdminActivity, len(rows))
	for i, r := range rows {
		activities[i] = AdminActivity{
			ID:        r.ID,
			Name:      r.Name,
			Icon:      r.Icon,
			Category:  string(r.Category),
			CreatedAt: r.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		}
	}
	return activities, nil
}

func (s *AdminActivityService) GetActivity(id int32) (*AdminActivity, error) {
	r, err := s.queries.AdminGetActivity(s.ctx, id)
	if err != nil {
		return nil, err
	}

	return &AdminActivity{
		ID:        r.ID,
		Name:      r.Name,
		Icon:      r.Icon,
		Category:  string(r.Category),
		CreatedAt: r.CreatedAt.Time.Format("2006-01-02 15:04:05"),
	}, nil
}

func (s *AdminActivityService) CreateActivity(name, icon string, category string) (*AdminActivity, error) {
	r, err := s.queries.AdminCreateActivity(s.ctx, repository.AdminCreateActivityParams{
		Name:     name,
		Icon:     &icon,
		Category: repository.Categories(category),
	})
	if err != nil {
		return nil, err
	}

	return &AdminActivity{
		ID:        r.ID,
		Name:      r.Name,
		Icon:      r.Icon,
		Category:  string(r.Category),
		CreatedAt: r.CreatedAt.Time.Format("2006-01-02 15:04:05"),
	}, nil
}

func (s *AdminActivityService) UpdateActivity(id int32, name, icon, category string) error {
	return s.queries.AdminUpdateActivity(s.ctx, repository.AdminUpdateActivityParams{
		ID:       id,
		Name:     name,
		Icon:     &icon,
		Category: repository.Categories(category),
	})
}

func (s *AdminActivityService) DeleteActivity(id int32) error {
	return s.queries.AdminDeleteActivity(s.ctx, id)
}

func (s *AdminActivityService) CountActivities(name *string, category *string) (int64, error) {
	return s.queries.AdminCountActivities(s.ctx, repository.AdminCountActivitiesParams{
		Name:     name,
		Category: toNullCategory(category),
	})
}
