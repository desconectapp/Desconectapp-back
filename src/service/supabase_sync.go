package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	repository "gin/db/generated"
)

type SupabaseSyncService struct {
	queries    *repository.Queries
	ctx        context.Context
	supabaseURL string
	anonKey     string
	httpClient  *http.Client
}

func NewSupabaseSyncService(queries *repository.Queries) *SupabaseSyncService {
	supabaseURL := os.Getenv("SUPABASE_URL")
	// Use service role key for backend operations (bypasses RLS)
	serviceRoleKey := os.Getenv("SUPABASE_SERVICE_ROLE_KEY")
	// Fallback to anon key if service role key not set (for backwards compatibility)
	if serviceRoleKey == "" {
		serviceRoleKey = os.Getenv("SUPABASE_ANON_KEY")
	}

	if supabaseURL == "" || serviceRoleKey == "" {
		// Log warning but don't fail - allows service to work without Supabase
		fmt.Println("WARN: SUPABASE_URL or SUPABASE_SERVICE_ROLE_KEY/SUPABASE_ANON_KEY not set, Supabase sync will be disabled")
	}

	return &SupabaseSyncService{
		queries:     queries,
		ctx:         context.Background(),
		supabaseURL: supabaseURL,
		anonKey:     serviceRoleKey, // Using service role key for backend operations
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// SyncGroupMembers syncs all members of a group to Supabase
// It deletes all existing members for the group and inserts the current ones
func (s *SupabaseSyncService) SyncGroupMembers(groupID int32) error {
	if s.supabaseURL == "" || s.anonKey == "" {
		// Supabase not configured, skip sync
		return nil
	}

	// Get all current members for this group
	members, err := s.queries.GetGroupMembers(s.ctx, groupID)
	if err != nil {
		return fmt.Errorf("failed to get group members: %w", err)
	}

	// Get user UUIDs for the members
	userUUIDs := make([]string, 0, len(members))
	for _, member := range members {
		if member.Uuid.Valid {
			userUUIDs = append(userUUIDs, member.Uuid.String())
		}
	}

	// Delete all existing members for this group in Supabase
	err = s.deleteGroupMembersFromSupabase(groupID)
	if err != nil {
		// Log error but continue - might be first sync
		fmt.Printf("WARN: Failed to delete existing group members from Supabase for group %d: %v\n", groupID, err)
	}

	// Insert current members into Supabase
	if len(userUUIDs) > 0 {
		err = s.insertGroupMembersToSupabase(groupID, userUUIDs)
		if err != nil {
			return fmt.Errorf("failed to insert group members to Supabase: %w", err)
		}
	}

	return nil
}

// DeleteGroupMembers deletes all members for a group from Supabase
func (s *SupabaseSyncService) DeleteGroupMembers(groupID int32) error {
	if s.supabaseURL == "" || s.anonKey == "" {
		return nil
	}

	return s.deleteGroupMembersFromSupabase(groupID)
}

func (s *SupabaseSyncService) deleteGroupMembersFromSupabase(groupID int32) error {
	url := fmt.Sprintf("%s/rest/v1/group_members?group_id=eq.%d", s.supabaseURL, groupID)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create DELETE request: %w", err)
	}

	req.Header.Set("apikey", s.anonKey)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.anonKey))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Prefer", "return=minimal")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute DELETE request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Supabase DELETE failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (s *SupabaseSyncService) insertGroupMembersToSupabase(groupID int32, userUUIDs []string) error {
	url := fmt.Sprintf("%s/rest/v1/group_members", s.supabaseURL)

	// Prepare data for bulk insert
	records := make([]map[string]interface{}, len(userUUIDs))
	for i, uuid := range userUUIDs {
		records[i] = map[string]interface{}{
			"group_id": groupID,
			"user_id": uuid,
		}
	}

	jsonData, err := json.Marshal(records)
	if err != nil {
		return fmt.Errorf("failed to marshal group members: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create POST request: %w", err)
	}

	req.Header.Set("apikey", s.anonKey)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.anonKey))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Prefer", "return=minimal")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute POST request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Supabase POST failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

