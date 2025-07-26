package service

import (
	"context"
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

	if len(matches) == 0 {
		newMatch, err := s.queries.CreatePartialMatch(s.ctx, repository.CreatePartialMatchParams{
			ActivityID:         request.ActivityID,
			Description:        request.Description,
			DayOfWeek:          request.DayOfWeek,
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

	for _, match := range matches {
		if isMatchCompatible(request, match) {
			members, err := s.queries.GetPartialMatchMembers(s.ctx, match.ID)
			if err != nil {
				return err
			}

			if len(members)+1 >= int(*match.ParticipantsNeeded) {
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
			} else {
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
			}
		}
	}

	return nil
}

func isMatchCompatible(request repository.ActivityRequest, match repository.PartialMatch) bool {
	if request.ActivityID != match.ActivityID {
		return false
	}
	if request.DayOfWeek != match.DayOfWeek {
		return false
	}

	return true
}
