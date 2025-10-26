package service

import (
	"context"
	"fmt"
	repository "gin/db/generated"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ActivitiesRequestService struct {
	queries         *repository.Queries
	ctx             context.Context
	matchingService *MatchingService
}

func NewActivitiesRequestService(conn *pgxpool.Pool) *ActivitiesRequestService {
	queries := repository.New(conn)
	ctx := context.Background()
	matchingService := NewMatchingService(conn)

	return &ActivitiesRequestService{
		queries:         queries,
		ctx:             ctx,
		matchingService: matchingService,
	}
}

func (s *ActivitiesRequestService) ListActivitiesRequests(params repository.ListActivityRequestsParams) ([]repository.ActivityRequest, error) {
	activities, err := s.queries.ListActivityRequests(s.ctx, params)
	if err != nil {
		return nil, err
	}
	return activities, nil
}

func (s *ActivitiesRequestService) CreateActivityRequest(params repository.CreateActivityRequestParams) (repository.ActivityRequest, error) {
	if params.UserID == nil || params.ActivityID == nil {
		return repository.ActivityRequest{}, fmt.Errorf("UserID and ActivityID cannot be nil")
	}

	// Primero chequeamos si la actividad es nueva (id=-1)
	if *params.ActivityID == -1 {
		newActivity, err := GenerateActivity(s, *params.Description)
		if err != nil {
			return repository.ActivityRequest{}, err
		}
		params.ActivityID = &newActivity.ID
		// Actualizamos la descripcion para que incluya el nombre de la actividad	
		fmt.Printf("Old description: %v\n", *params.Description)
		formattedDescription := fmt.Sprintf("%s %s", *newActivity.Icon, newActivity.Name)
		params.Description = &formattedDescription
		fmt.Printf("Updated description to: %v\n", *params.Description)
	}

	existingActivityRequest, err := s.queries.GetActivityRequestByUserAndActivityID(s.ctx, repository.GetActivityRequestByUserAndActivityIDParams{
		UserID:     params.UserID,
		ActivityID: params.ActivityID,
	})
	if err != nil && err != pgx.ErrNoRows {
		return repository.ActivityRequest{}, err
	}

	if existingActivityRequest.ID != 0 {
		fmt.Printf("Deleting activity request %v", existingActivityRequest)
		s.queries.DeleteActivityRequest(s.ctx, existingActivityRequest.ID)
		s.queries.DeletePartialMatchesByUserAndActivityID(s.ctx, repository.DeletePartialMatchesByUserAndActivityIDParams{
			UserID:     *existingActivityRequest.UserID,
			ActivityID: existingActivityRequest.ActivityID,
		})
	}
	
	// Cambio la descripcion para que sea el nombre de la actividad
	activity, err := s.queries.GetActivityByID(s.ctx, *params.ActivityID)
	if err != nil {
		return repository.ActivityRequest{}, err
	}
	formattedDescription := fmt.Sprintf("%s %s", *activity.Icon, activity.Name)
	params.Description = &formattedDescription

	locationName, err := getLocationFromCoordinates(fmt.Sprintf("%f", *params.Latitude), fmt.Sprintf("%f", *params.Longitude))
	if err != nil {
		return repository.ActivityRequest{}, err
	}
	params.LocationName = &locationName

	// Armamos la ActivityRequest
	activityReq, err := s.queries.CreateActivityRequest(s.ctx, params)
	if err != nil {
		return repository.ActivityRequest{}, err
	}
	err = s.matchingService.FindMatches(activityReq)
	if err != nil {
		return repository.ActivityRequest{}, err
	}
	return activityReq, nil
}

func (s *ActivitiesRequestService) GetActivities(params repository.GetActivitiesParams) ([]repository.Activity, error) {
	activities, err := s.queries.GetActivities(s.ctx, params)

	if err != nil {
		return nil, err
	}

	return activities, nil
}

func (s *ActivitiesRequestService) DeleteActivityRequest(requestId int) error {
	err := s.queries.DeleteActivityRequest(s.ctx, int32(requestId))
	if err != nil {
		return err
	}
	return nil
}

// func (s *ActivitiesRequestService) CreateActivity(name string) (repository.Activity, error) {
// 	activity, err := GenerateActivity(s, name)
// 	if err != nil {
// 		return repository.Activity{}, err
// 	}
// 	return activity, nil
// }