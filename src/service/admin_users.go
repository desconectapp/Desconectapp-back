package service

import (
	"context"

	repository "gin/db/generated"
	"github.com/jackc/pgx/v5"
)

type AdminUserService struct {
	queries *repository.Queries
	ctx     context.Context
}


type AdminUser struct {
	ID               int32  `json:"id"`
	Email            string `json:"email"`
	EmailValidated   bool   `json:"email_validated"`
	Name             string `json:"name"`
	Age              int32  `json:"age"`
	City             string `json:"city"`
	CurrentSituation string `json:"current_situation"`
	Gender           string `json:"gender"`
	ProfileComplete  bool   `json:"profile_complete"`
	CreatedAt        string `json:"created_at"`
	IsSuspended 	 bool	`json:"is_suspended"`
}

func NewAdminUserService(conn *pgx.Conn) *AdminUserService {
	queries := repository.New(conn)
	ctx := context.Background()

	return &AdminUserService{
		queries: queries,
		ctx:     ctx,
	}
}

func (s *AdminUserService) ListUsers(limit, offset int32, email, name *string, validated *bool) ([]AdminUser, error) {
	rows, err := s.queries.AdminListUsers(s.ctx, repository.AdminListUsersParams{
		Limit:          limit,
		Offset:         offset,
		Email:          email,
		Name:           name,
		EmailValidated: validated,
	})
	if err != nil {
		return nil, err
	}

	users := make([]AdminUser, len(rows))
	for i, r := range rows {
		users[i] = AdminUser{
			ID:               r.ID,
			Email:            r.Email,
			EmailValidated:   r.EmailValidated,
			Name:             r.Name,
			Age:              r.Age,
			City:             r.City,
			CurrentSituation: r.CurrentSituation,
			Gender:           r.Gender,
			ProfileComplete:  r.ProfileComplete,
			CreatedAt:        r.CreatedAt.Time.Format("2006-01-02 15:04:05"),
			IsSuspended:      r.IsSuspended,
		}
	}
	return users, nil
}

func (s *AdminUserService) GetUser(id int32) (*AdminUser, error) {
	r, err := s.queries.AdminGetUser(s.ctx, id)
	if err != nil {
		return nil, err
	}

	return &AdminUser{
		ID:               r.ID,
		Email:            r.Email,
		EmailValidated:   r.EmailValidated,
		Name:             r.Name,
		Age:              r.Age,
		City:             r.City,
		CurrentSituation: r.CurrentSituation,
		Gender:           r.Gender,
		ProfileComplete:  r.ProfileComplete,
		CreatedAt:        r.CreatedAt.Time.Format("2006-01-02 15:04:05"),
	}, nil
}

func (s *AdminUserService) CreateUserWithProfile(email, password string, validated bool, name string, age int32, city, situation, gender string) (*AdminUser, error) {
	user, err := s.queries.AdminCreateUser(s.ctx, repository.AdminCreateUserParams{
		Email:          email,
		Password:       password,
		EmailValidated: validated,
	})
	if err != nil {
		return nil, err
	}

	_, err = s.queries.AdminCreateProfile(s.ctx, repository.AdminCreateProfileParams{
		UserID:           user.ID,
		Name:             name,
		Age:              age,
		City:             city,
		CurrentSituation: situation,
		Gender:           gender,
	})
	if err != nil {
		return nil, err
	}

	return s.GetUser(user.ID)
}

func (s *AdminUserService) UpdateUserWithProfile(id int32, email string, validated bool, name string, age int32, city, situation, gender string, complete bool) error {
	err := s.queries.AdminUpdateUser(s.ctx, repository.AdminUpdateUserParams{
		ID:             id,
		Email:          email,
		EmailValidated: validated,
	})
	if err != nil {
		return err
	}

	return s.queries.AdminUpdateProfile(s.ctx, repository.AdminUpdateProfileParams{
		UserID:           id,
		Name:             name,
		Age:              age,
		City:             city,
		CurrentSituation: situation,
		Gender:           gender,
		ProfileComplete:  complete,
	})
}

func (s *AdminUserService) DeleteUser(id int32) error {
	return s.queries.AdminDeleteUser(s.ctx, id)
}

func (s *AdminUserService) CountUsers(email, name *string, validated *bool) (int64, error) {
	return s.queries.AdminCountUsers(s.ctx, repository.AdminCountUsersParams{
		Email:          email,
		Name:           name,
		EmailValidated: validated,
	})
}

func (s *AdminUserService) SuspendUser(id int32) any {
	return s.queries.AdminSuspendUser(s.ctx, id)
}

