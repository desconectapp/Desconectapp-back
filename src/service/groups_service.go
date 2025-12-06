package service

import (
	"context"
	"fmt"
	repository "gin/db/generated"
	"log"
	"strconv"
	"strings"

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
	ID            int32   `json:"id"`
	Name          *string `json:"name"`
	Description   *string `json:"description"`
	Location      *string `json:"location"`
	LocationName  *string `json:"location_name"`
	Coords        *string `json:"coords"`
	AvatarUrl     *string `json:"avatar_url"`
	WeekTimeslots []int32 `json:"week_timeslots"`
	ActivityName  string  `json:"activity_name"`
	MemberCount   int32   `json:"member_count"`
	Photo         string  `json:"photo"`
	Time          string  `json:"created_at"`
}

type ActivityFilter struct {
	ActivityId      int32    `json:"activity_id"`
	ActivityIds     []int32  `json:"activity_ids"`
	MyPreferences   bool     `json:"my_preferences"`
	UserID          *int32   `json:"user_id"`
	Limit           int32    `json:"limit"`
	Offset          int32    `json:"offset"`
	Latitude        *float64 `json:"latitude"`
	Longitude       *float64 `json:"longitude"`
	Radius          *float64 `json:"radius"`
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
	// Get location name from coordinates
	parts := strings.Split(*groupParams.Location, ",")
	// El front las manda como long,lat
	lat := parts[1]
	long := parts[0]
	location, err := getLocationFromCoordinates(lat, long)
	if err != nil {
		log.Println("Error getting location from coordinates:", err)
		return repository.CreateGroupRow{}, err
	}
	groupParams.LocationName = &location
	// Create the group
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
		Location:    group.LocationName,
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

	// Check if location parameters are provided first
	hasLocation := filter.Latitude != nil && filter.Longitude != nil && filter.Radius != nil

	// If location parameters are provided, always use location-based filtering
	if hasLocation {
		if filter.MyPreferences && filter.UserID != nil {
			// Use location + preferences
			return s.GetPreferredGroupsWithLocation(repository.GetPreferredGroupsWithLocationParams{
				UserID:    *filter.UserID,
				Latitude:  *filter.Latitude,
				Longitude: *filter.Longitude,
				Radius:    *filter.Radius,
				Limit:     filter.Limit,
				Offset:    filter.Offset,
			})
		} else {
			// Use location + optional activity filter
			if len(filter.ActivityIds) > 0 {
				log.Printf("DEBUG: Location + multiple activities case. Activities: %v, Location: lat=%f, lng=%f, radius=%f", 
					filter.ActivityIds, *filter.Latitude, *filter.Longitude, *filter.Radius)
				log.Printf("DEBUG: NOTE - Request params are: latitude=%f, longitude=%f", *filter.Latitude, *filter.Longitude)
				
				// For multiple activities with location, we need to get all groups by activities first
				// then filter by location in the application layer (temporary solution)
				allActivityGroups, err := s.GetOpenGroupsByActivities(repository.GetOpenGroupsByActivitiesParams{
					ActivityIds: filter.ActivityIds,
					UserID:      filter.UserID,
					Limit:       filter.Limit * 3, // Get more to account for location filtering
					Offset:      filter.Offset,
				})
				if err != nil {
					log.Printf("DEBUG: Error getting groups by activities: %v", err)
					return nil, err
				}
				log.Printf("DEBUG: Found %d groups for activities %v", len(allActivityGroups), filter.ActivityIds)
				
				// Filter by location manually (swap lat/lng to match SQL query behavior)
				filteredGroups, err := s.filterGroupsByLocation(allActivityGroups, *filter.Longitude, *filter.Latitude, *filter.Radius, int(filter.Limit))
				log.Printf("DEBUG: After location filtering: %d groups", len(filteredGroups))
				return filteredGroups, err
			} else {
				// Single activity or no activity filter
				var activityID *int32
				if filter.ActivityId != 0 {
					activityID = &filter.ActivityId
				}

				openGroups, err = s.GetOpenGroupsWithLocation(
					repository.GetOpenGroupsWithLocationParams{
						Limit:      filter.Limit,
						Offset:     filter.Offset,
						Latitude:   *filter.Latitude,
						Longitude:  *filter.Longitude,
						Radius:     *filter.Radius,
						ActivityID: activityID,
						UserID:     filter.UserID,
					},
				)
			}
		}
	} else {
		// No location - use preference or activity-based filtering
		if filter.MyPreferences && filter.UserID != nil {
			return s.GetUserRecommendations(repository.GetPreferredGroupsParams{
				UserID: *filter.UserID,
				Limit:  filter.Limit,
				Offset: filter.Offset,
			})
		} else if len(filter.ActivityIds) > 0 {
			return s.GetOpenGroupsByActivities(repository.GetOpenGroupsByActivitiesParams{
				ActivityIds: filter.ActivityIds,
				UserID:      filter.UserID,
				Limit:       filter.Limit,
				Offset:      filter.Offset,
			})
		} else if filter.ActivityId == 0 {
			openGroups, err = s.GetAllOpenGroups(
				repository.GetAllOpenGroupsParams{
					Limit:      filter.Limit,
					Offset:     filter.Offset,
					ActivityID: nil,
					UserID:     filter.UserID,
				},
			)
		} else {
			log.Println("Fetching open groups with activity filter:", filter.ActivityId)
			openGroups, err = s.GetPublicOpenGroupsWithFilter(
				repository.GetOpenGroupsWithFilterParams{
					Limit:      filter.Limit,
					Offset:     filter.Offset,
					ActivityID: &filter.ActivityId,
					UserID:     filter.UserID,
				},
			)
		}
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
			ID:            group.ID,
			Name:          group.Name,
			Location:      group.LocationName,
			LocationName:  group.LocationName,
			Coords:        group.Location,
			AvatarUrl:     group.AvatarUrl,
			WeekTimeslots: group.WeekTimeslots,
			Description:   group.Description,
			ActivityName:  group.ActivityName,
			MemberCount:   int32(group.MemberCount),
			Photo:         *group.Icon,
			Time:          group.CreatedAt.Time.String(),
		})
	}
	return openGroups, err
}

func (s *GroupsService) GetAllOpenGroups(filter repository.GetAllOpenGroupsParams) ([]OpenGroup, error) {
	groups, err := s.queries.GetAllOpenGroups(s.ctx, filter)

	if err != nil {
		return nil, err
	}

	var openGroups []OpenGroup

	for _, group := range groups {
		openGroups = append(openGroups, OpenGroup{
			ID:            group.ID,
			Name:          group.Name,
			Location:      group.LocationName,
			LocationName:  group.LocationName,
			Coords:        group.Location,
			AvatarUrl:     group.AvatarUrl,
			WeekTimeslots: group.WeekTimeslots,
			Description:   group.Description,
			ActivityName:  group.ActivityName,
			MemberCount:   int32(group.MemberCount),
			Photo:         *group.Icon,
			Time:          group.CreatedAt.Time.String(),
		})
	}
	return openGroups, err
}

func (s *GroupsService) GetPreferredGroupsWithLocation(filter repository.GetPreferredGroupsWithLocationParams) ([]OpenGroup, error) {
	groups, err := s.queries.GetPreferredGroupsWithLocation(s.ctx, filter)

	if err != nil {
		return nil, err
	}

	var openGroups []OpenGroup

	for _, group := range groups {
		openGroups = append(openGroups, OpenGroup{
			ID:            group.ID,
			Name:          group.Name,
			Location:      group.LocationName,
			LocationName:  group.LocationName,
			Coords:        group.Location,
			AvatarUrl:     group.AvatarUrl,
			WeekTimeslots: group.WeekTimeslots,
			Description:   group.Description,
			ActivityName:  group.ActivityName,
			MemberCount:   int32(group.MemberCount),
			Photo:         *group.Icon,
			Time:          group.CreatedAt.Time.String(),
		})
	}
	return openGroups, err
}

func (s *GroupsService) GetOpenGroupsByActivities(filter repository.GetOpenGroupsByActivitiesParams) ([]OpenGroup, error) {
	groups, err := s.queries.GetOpenGroupsByActivities(s.ctx, filter)

	if err != nil {
		return nil, err
	}

	var openGroups []OpenGroup

	for _, group := range groups {
		openGroups = append(openGroups, OpenGroup{
			ID:            group.ID,
			Name:          group.Name,
			Location:      group.LocationName,
			LocationName:  group.LocationName,
			Coords:        group.Location,
			AvatarUrl:     group.AvatarUrl,
			WeekTimeslots: group.WeekTimeslots,
			Description:   group.Description,
			ActivityName:  group.ActivityName,
			MemberCount:   int32(group.MemberCount),
			Photo:         *group.Icon,
			Time:          group.CreatedAt.Time.String(),
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
			ID:            group.ID,
			Name:          group.Name,
			Location:      group.LocationName,
			LocationName:  group.LocationName,
			Coords:        group.Location,
			AvatarUrl:     group.AvatarUrl,
			WeekTimeslots: group.WeekTimeslots,
			Description:   group.Description,
			ActivityName:  group.ActivityName,
			MemberCount:   int32(group.MemberCount),
			Photo:         *group.Icon,
			Time:          group.CreatedAt.Time.String(),
		})
	}
	return openGroups, err
}

func (s *GroupsService) GetOpenGroupsWithLocation(filter repository.GetOpenGroupsWithLocationParams) ([]OpenGroup, error) {
	groups, err := s.queries.GetOpenGroupsWithLocation(s.ctx, filter)

	if err != nil {
		return nil, err
	}

	var openGroups []OpenGroup
	for _, group := range groups {
		openGroups = append(openGroups, OpenGroup{
			ID:            group.ID,
			Name:          group.Name,
			Location:      group.Location,
			LocationName:  group.LocationName,
			AvatarUrl:     group.AvatarUrl,
			WeekTimeslots: group.WeekTimeslots,
			Description:   group.Description,
			ActivityName:  group.ActivityName,
			MemberCount:   int32(group.MemberCount),
			Photo:         *group.Icon,
			Time:          group.CreatedAt.Time.String(),
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
			ID:            group.ID,
			Name:          group.Name,
			Location:      group.LocationName,
			LocationName:  group.LocationName,
			Coords:        group.Location,
			AvatarUrl:     group.AvatarUrl,
			WeekTimeslots: group.WeekTimeslots,
			Description:   group.Description,
			ActivityName:  group.ActivityName,
			MemberCount:   int32(group.MemberCount),
			Photo:         *group.ActivityIcon,
			Time:          group.CreatedAt.Time.String(),
		})
	}
	return openGroups, err
}

func (s *GroupsService) filterGroupsByLocation(groups []OpenGroup, lat, lng, radius float64, limit int) ([]OpenGroup, error) {
	var filteredGroups []OpenGroup
	log.Printf("DEBUG filterGroupsByLocation: Processing %d groups, target lat=%f, lng=%f, radius=%f (NOTE: params may be swapped to match SQL behavior)", len(groups), lat, lng, radius)
	
	for i, group := range groups {
		log.Printf("DEBUG Group %d: ID=%d, Name=%s, Coords=%v", i+1, group.ID, 
			func() string { if group.Name != nil { return *group.Name } else { return "NULL" } }(), 
			func() string { if group.Coords != nil { return *group.Coords } else { return "NULL" } }())
			
		if group.Coords != nil && *group.Coords != "" {
			// Parse coordinates from group.Coords (format: "lng,lat")
			parts := strings.Split(*group.Coords, ",")
			log.Printf("DEBUG: Coords parts: %v", parts)
			if len(parts) == 2 {
				groupLat, err1 := strconv.ParseFloat(parts[1], 64)
				groupLng, err2 := strconv.ParseFloat(parts[0], 64)
				
				if err1 == nil && err2 == nil {
					// Calculate distance using Haversine formula
					distance := calculateDistance(lat, lng, groupLat, groupLng)
					log.Printf("DEBUG: Group %d distance: %f km (threshold: %f km)", group.ID, distance, radius)
					if distance <= radius {
						log.Printf("DEBUG: Group %d INCLUDED (distance %f <= %f)", group.ID, distance, radius)
						filteredGroups = append(filteredGroups, group)
					} else {
						log.Printf("DEBUG: Group %d EXCLUDED (distance %f > %f)", group.ID, distance, radius)
					}
				} else {
					log.Printf("DEBUG: Error parsing coords for group %d: err1=%v, err2=%v", group.ID, err1, err2)
				}
			} else {
				log.Printf("DEBUG: Invalid coords format for group %d: %s", group.ID, *group.Coords)
			}
		} else {
			log.Printf("DEBUG: Group %d has no coordinates", group.ID)
		}
	}
	
	// Limit results
	if len(filteredGroups) > limit {
		filteredGroups = filteredGroups[:limit]
	}
	
	return filteredGroups, nil
}
