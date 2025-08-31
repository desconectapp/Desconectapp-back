package service

import (
	"context"
	"fmt"
	repository "gin/db/generated"
	"math"

	"github.com/jackc/pgx/v5"
)

type MatchingService struct {
	queries *repository.Queries
	ctx     context.Context
}

func NewMatchingService(conn *pgx.Conn) *MatchingService {
	queries := repository.New(conn)
	ctx := context.Background()

	return &MatchingService{
		queries: queries,
		ctx:     ctx,
	}
}

func (s *MatchingService) FindMatches(request repository.ActivityRequest) error {
	matches, err := s.queries.FindPartialMatches(s.ctx, request.ActivityID)
	if err != nil {
		return err
	}

	fmt.Printf("DEBUG: Found %d matches for ActivityID=%v, WeekHours=%v\n", len(matches), request.ActivityID, request.WeekHours)

	for _, match := range matches {
		fmt.Printf("DEBUG: Checking match ID=%d\n", match.ID)
		if isMatchCompatible(request, match) {
			fmt.Printf("DEBUG: Match is compatible\n")
			members, err := s.queries.GetPartialMatchMembers(s.ctx, match.ID)
			if err != nil {
				return err
			}

			fmt.Printf("DEBUG: Current members count: %d, participants needed: %d\n", len(members), *match.ParticipantsNeeded)
			if len(members)+1 >= int(*match.ParticipantsNeeded) && len(members)+1 >= int(*request.ParticipantsNeeded) {
				fmt.Printf("DEBUG: Creating group (enough members)\n")

				s.createGroup(request, match, members)
				return nil
			} else {
				fmt.Printf("DEBUG: Creating combined partial match\n")
				s.createCombinedPartialMatch(request, match)
			}
		} else {
			fmt.Printf("DEBUG: Match is not compatible\n")
		}
	}

	fmt.Printf("DEBUG: Creating unconditional partial match for user %v\n", request.UserID)
	s.createUnconditionalPartialMatch(request)

	return nil

}

func isMatchCompatible(request repository.ActivityRequest, match repository.PartialMatch) bool {
	reqActivityID := int32(0)
	if request.ActivityID != nil {
		reqActivityID = *request.ActivityID
	}
	matchActivityID := int32(0)
	if match.ActivityID != nil {
		matchActivityID = *match.ActivityID
	}

	fmt.Printf("DEBUG COMPATIBILITY: Request ActivityID=%d, Match ActivityID=%d\n", reqActivityID, matchActivityID)

	if reqActivityID != matchActivityID {
		fmt.Printf("DEBUG COMPATIBILITY: ActivityID mismatch!\n")
		return false
	}

	if !hasSharedWeekHours(request.WeekHours, match.WeekHours) {
		fmt.Printf("DEBUG COMPATIBILITY: No shared week hours!\n")
		return false
	}

	if !areLocationsCompatible(request, match) {
		fmt.Printf("DEBUG COMPATIBILITY: Locations not compatible!\n")
		return false
	}

	if !areParticipantCountsCompatible(request, match) {
		fmt.Printf("DEBUG COMPATIBILITY: Participant counts not compatible!\n")
		return false
	}

	fmt.Printf("DEBUG COMPATIBILITY: Match is compatible!\n")
	return true
}

func (s *MatchingService) createCombinedPartialMatch(request repository.ActivityRequest, match repository.PartialMatch) error {
	intersectionWeekHours := getWeekHoursIntersection(request.WeekHours, match.WeekHours)

	var midLat, midLon *float64
	if request.Latitude != nil && request.Longitude != nil && match.Latitude != nil && match.Longitude != nil {
		lat, lon := calculateMidpoint(*request.Latitude, *request.Longitude, *match.Latitude, *match.Longitude)
		midLat = &lat
		midLon = &lon
	}

	var avgSearchRadius *int32
	if request.SearchRadius != nil && match.SearchRadius != nil {
		avg := (*request.SearchRadius + *match.SearchRadius) / 2
		avgSearchRadius = &avg
	}

	var highestMinParticipants *int32
	reqMin := int32(1)
	if request.ParticipantsNeeded != nil {
		reqMin = *request.ParticipantsNeeded
	}
	matchMin := int32(1)
	if match.ParticipantsNeeded != nil {
		matchMin = *match.ParticipantsNeeded
	}
	if reqMin > matchMin {
		highestMinParticipants = &reqMin
	} else {
		highestMinParticipants = &matchMin
	}

	var lowestMaxParticipants *int32
	reqMax := int32(10)
	if request.MaximumParticipants != nil {
		reqMax = *request.MaximumParticipants
	}
	matchMax := int32(10)
	if match.MaximumParticipants != nil {
		matchMax = *match.MaximumParticipants
	}
	if reqMax < matchMax {
		lowestMaxParticipants = &reqMax
	} else {
		lowestMaxParticipants = &matchMax
	}

	members, err := s.queries.GetPartialMatchMembers(s.ctx, match.ID)
	if err != nil {
		return err
	}
	newMembersCount := int32(len(members) + 1)

	newMatch, err := s.queries.CreatePartialMatch(s.ctx, repository.CreatePartialMatchParams{
		ActivityID:          match.ActivityID,
		Description:         match.Description,
		WeekHours:           intersectionWeekHours,
		ParticipantsNeeded:  highestMinParticipants,
		MaximumParticipants: lowestMaxParticipants,
		MembersCount:        &newMembersCount,
		Latitude:            midLat,
		Longitude:           midLon,
		SearchRadius:        avgSearchRadius,
	})
	if err != nil {
		return err
	}

	memberIDs := make([]int32, len(members))
	for i, member := range members {
		memberIDs[i] = member.ID
	}
	if len(memberIDs) > 0 {
		err = s.queries.BatchAddUsersToPartialMatch(s.ctx, repository.BatchAddUsersToPartialMatchParams{
			PartialMatchID: newMatch.ID,
			UserIds:        memberIDs,
		})
		if err != nil {
			return err
		}
	}

	if request.UserID != nil {
		err = s.queries.AddUserToPartialMatch(s.ctx, repository.AddUserToPartialMatchParams{
			PartialMatchID: newMatch.ID,
			UserID:         *request.UserID,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *MatchingService) createGroup(request repository.ActivityRequest, match repository.PartialMatch, members []repository.GetPartialMatchMembersRow) error {
	groupID, err := s.queries.CreateGroup(s.ctx, repository.CreateGroupParams{
		Name:        match.Description,
		Description: match.Description,
		ActivityID:  *match.ActivityID,
	})
	if err != nil {
		return err
	}

	memberIDs := make([]int32, len(members))
	for i, member := range members {
		memberIDs[i] = member.ID
	}
	if len(memberIDs) > 0 {
		err = s.queries.BatchAddUserToGroup(s.ctx, repository.BatchAddUserToGroupParams{
			GroupID: groupID,
			UserIds: memberIDs,
		})
		if err != nil {
			return err
		}
	}

	if request.UserID != nil {
		err = s.queries.AddUserToGroup(s.ctx, repository.AddUserToGroupParams{
			GroupID: groupID,
			UserID:  *request.UserID,
		})
		if err != nil {
			return err
		}
	}

	for _, member := range members {
		err = s.queries.DeletePartialMatchesByUserAndActivityID(s.ctx, repository.DeletePartialMatchesByUserAndActivityIDParams{
			UserID:     member.ID,
			ActivityID: match.ActivityID,
		})
		if err != nil {
			return err
		}
	}
	if request.UserID != nil {
		err = s.queries.DeletePartialMatchesByUserAndActivityID(s.ctx, repository.DeletePartialMatchesByUserAndActivityIDParams{
			UserID:     *request.UserID,
			ActivityID: request.ActivityID,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *MatchingService) createUnconditionalPartialMatch(request repository.ActivityRequest) error {
	newMatch, err := s.queries.CreatePartialMatch(s.ctx, repository.CreatePartialMatchParams{
		ActivityID:          request.ActivityID,
		Description:         request.Description,
		WeekHours:           request.WeekHours,
		ParticipantsNeeded:  request.ParticipantsNeeded,
		MaximumParticipants: request.MaximumParticipants,
		MembersCount:        &[]int32{1}[0],
		Latitude:            request.Latitude,
		Longitude:           request.Longitude,
		SearchRadius:        request.SearchRadius,
	})
	if err != nil {
		return err
	}

	if request.UserID != nil {
		err = s.queries.AddUserToPartialMatch(s.ctx, repository.AddUserToPartialMatchParams{
			PartialMatchID: newMatch.ID,
			UserID:         *request.UserID,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func hasSharedWeekHours(hours1, hours2 []int32) bool {
	hourSet := make(map[int32]bool)
	for _, hour := range hours1 {
		hourSet[hour] = true
	}

	for _, hour := range hours2 {
		if hourSet[hour] {
			return true
		}
	}

	return false
}

func calculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371

	lat1Rad := lat1 * math.Pi / 180
	lon1Rad := lon1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lon2Rad := lon2 * math.Pi / 180

	deltaLat := lat2Rad - lat1Rad
	deltaLon := lon2Rad - lon1Rad

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}

func areLocationsCompatible(request repository.ActivityRequest, match repository.PartialMatch) bool {
	if request.Latitude == nil || request.Longitude == nil || request.SearchRadius == nil ||
		match.Latitude == nil || match.Longitude == nil || match.SearchRadius == nil {
		return false
	}

	distance := calculateDistance(*request.Latitude, *request.Longitude, *match.Latitude, *match.Longitude)

	requestInMatchRadius := distance <= float64(*match.SearchRadius)

	matchInRequestRadius := distance <= float64(*request.SearchRadius)

	return requestInMatchRadius && matchInRequestRadius
}

func areParticipantCountsCompatible(request repository.ActivityRequest, match repository.PartialMatch) bool {
	reqMinParticipants := int32(1)
	if request.ParticipantsNeeded != nil {
		reqMinParticipants = *request.ParticipantsNeeded
	}

	reqMaxParticipants := int32(10)
	if request.MaximumParticipants != nil {
		reqMaxParticipants = *request.MaximumParticipants
	}

	matchMinParticipants := int32(1)
	if match.ParticipantsNeeded != nil {
		matchMinParticipants = *match.ParticipantsNeeded
	}

	matchMaxParticipants := int32(10)
	if match.MaximumParticipants != nil {
		matchMaxParticipants = *match.MaximumParticipants
	}

	return reqMaxParticipants >= matchMinParticipants && matchMaxParticipants >= reqMinParticipants
}

func getWeekHoursIntersection(hours1, hours2 []int32) []int32 {
	hourSet := make(map[int32]bool)
	for _, hour := range hours1 {
		hourSet[hour] = true
	}

	var intersection []int32
	addedHours := make(map[int32]bool)

	for _, hour := range hours2 {
		if hourSet[hour] && !addedHours[hour] {
			intersection = append(intersection, hour)
			addedHours[hour] = true
		}
	}

	return intersection
}

func calculateMidpoint(lat1, lon1, lat2, lon2 float64) (float64, float64) {
	midLat := (lat1 + lat2) / 2
	midLon := (lon1 + lon2) / 2
	return midLat, midLon
}
