package service

import (
	"context"
	"fmt"
	repository "gin/db/generated"

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
	matches, err := s.queries.FindPartialMatches(s.ctx, repository.FindPartialMatchesParams{
		ActivityID: request.ActivityID,
		DayOfWeek:  request.DayOfWeek,
	})
	if err != nil {
		return err
	}

	fmt.Printf("DEBUG: Found %d matches for ActivityID=%v, DayOfWeek=%v\n", len(matches), request.ActivityID, request.DayOfWeek)

	for _, match := range matches {
		fmt.Printf("DEBUG: Checking match ID=%d\n", match.ID)
		if isMatchCompatible(request, match) {
			fmt.Printf("DEBUG: Match is compatible\n")
			members, err := s.queries.GetPartialMatchMembers(s.ctx, match.ID)
			if err != nil {
				return err
			}

			fmt.Printf("DEBUG: Current members count: %d, participants needed: %d\n", len(members), *match.ParticipantsNeeded)
			if len(members)+1 >= int(*match.ParticipantsNeeded) {
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
	fmt.Printf("DEBUG COMPATIBILITY: Request DayOfWeek=%v, Match DayOfWeek=%v\n", request.DayOfWeek, match.DayOfWeek)

	if reqActivityID != matchActivityID {
		fmt.Printf("DEBUG COMPATIBILITY: ActivityID mismatch!\n")
		return false
	}
	if request.DayOfWeek != match.DayOfWeek {
		fmt.Printf("DEBUG COMPATIBILITY: DayOfWeek mismatch!\n")
		return false
	}

	fmt.Printf("DEBUG COMPATIBILITY: Match is compatible!\n")
	return true
}

func (s *MatchingService) createCombinedPartialMatch(request repository.ActivityRequest, match repository.PartialMatch) error {
	newMatch, err := s.queries.CreatePartialMatch(s.ctx, repository.CreatePartialMatchParams{
		ActivityID:         match.ActivityID,
		Description:        match.Description,
		DayOfWeek:          match.DayOfWeek,
		ParticipantsNeeded: match.ParticipantsNeeded,
		MembersCount:       func() *int32 { v := *match.MembersCount + 1; return &v }(),
	})
	if err != nil {
		return err
	}

	members, err := s.queries.GetPartialMatchMembers(s.ctx, match.ID)
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
		err = s.queries.DeletePartialMatchesByUser(s.ctx, member.ID)
		if err != nil {
			return err
		}
	}
	if request.UserID != nil {
		err = s.queries.DeletePartialMatchesByUser(s.ctx, *request.UserID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *MatchingService) createUnconditionalPartialMatch(request repository.ActivityRequest) error {
	newMatch, err := s.queries.CreatePartialMatch(s.ctx, repository.CreatePartialMatchParams{
		ActivityID:         request.ActivityID,
		Description:        request.Description,
		DayOfWeek:          request.DayOfWeek,
		MembersCount:       &[]int32{1}[0],
		ParticipantsNeeded: request.ParticipantsNeeded,
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
