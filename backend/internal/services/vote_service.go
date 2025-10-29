package services

import (
	"context"
	"fmt"
	"time"

	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/repository/interfaces"
	"spark-park-cricket-backend/internal/utils"

	"github.com/google/uuid"
)

// VoteService implements the VoteServiceInterface
type VoteService struct {
	voteRepo     interfaces.VoteRepositoryInterface
	groupService GroupServiceInterface
}

// NewVoteService creates a new vote service
func NewVoteService(voteRepo interfaces.VoteRepositoryInterface, groupService GroupServiceInterface) VoteServiceInterface {
	return &VoteService{
		voteRepo:     voteRepo,
		groupService: groupService,
	}
}

// CreateVote creates a new vote with options
func (s *VoteService) CreateVote(ctx context.Context, req *models.CreateVoteRequest, userID string) (*models.Vote, error) {
	// Validate request
	if err := s.validateCreateVoteRequest(req); err != nil {
		return nil, err
	}

	// Create vote
	vote := &models.Vote{
		ID:          uuid.New().String(),
		Title:       req.Title,
		Description: req.Description,
		Type:        req.Type,
		Status:      models.VoteStatusActive,
		CreatedBy:   userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Create vote in database
	if err := s.voteRepo.CreateVote(ctx, vote); err != nil {
		utils.LogError(err, "Failed to create vote", map[string]interface{}{
			"user_id": userID,
			"title":   req.Title,
		})
		return nil, fmt.Errorf("failed to create vote: %w", err)
	}

	// Create vote options
	options := make([]*models.VoteOption, len(req.Options))
	for i, optionText := range req.Options {
		options[i] = &models.VoteOption{
			ID:          uuid.New().String(),
			VoteID:      vote.ID,
			Text:        optionText,
			Description: "",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
	}

	if err := s.voteRepo.CreateVoteOptions(ctx, options); err != nil {
		// Clean up the vote if options creation fails
		s.voteRepo.DeleteVote(ctx, vote.ID)
		utils.LogError(err, "Failed to create vote options", map[string]interface{}{
			"vote_id": vote.ID,
		})
		return nil, fmt.Errorf("failed to create vote options: %w", err)
	}

	utils.LogInfo("Vote created successfully", map[string]interface{}{
		"vote_id": vote.ID,
		"user_id": userID,
		"title":   vote.Title,
		"type":    vote.Type,
	})

	return vote, nil
}

// GetVote retrieves a vote with its options
func (s *VoteService) GetVote(ctx context.Context, id string) (*models.VoteWithOptions, error) {
	voteWithOptions, err := s.voteRepo.GetVoteWithOptions(ctx, id)
	if err != nil {
		utils.LogError(err, "Failed to get vote", map[string]interface{}{
			"vote_id": id,
		})
		return nil, fmt.Errorf("failed to get vote: %w", err)
	}

	return voteWithOptions, nil
}

// GetVoteWithResults retrieves a vote with results and user's vote
func (s *VoteService) GetVoteWithResults(ctx context.Context, id string, userID string) (*models.VoteWithResults, error) {
	voteWithResults, err := s.voteRepo.GetVoteWithResults(ctx, id, userID)
	if err != nil {
		utils.LogError(err, "Failed to get vote with results", map[string]interface{}{
			"vote_id": id,
			"user_id": userID,
		})
		return nil, fmt.Errorf("failed to get vote with results: %w", err)
	}

	return voteWithResults, nil
}

// UpdateVote updates an existing vote
func (s *VoteService) UpdateVote(ctx context.Context, id string, req *models.UpdateVoteRequest, userID string) (*models.Vote, error) {
	// Get existing vote
	vote, err := s.voteRepo.GetVoteByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check if user is the creator
	if vote.CreatedBy != userID {
		return nil, fmt.Errorf("unauthorized: only vote creator can update vote")
	}

	// Check if vote is still active
	if vote.Status != models.VoteStatusActive {
		return nil, fmt.Errorf("cannot update closed or cancelled vote")
	}

	// Update fields
	if req.Title != nil {
		vote.Title = *req.Title
	}
	if req.Description != nil {
		vote.Description = *req.Description
	}
	if req.Status != nil {
		vote.Status = *req.Status
		if *req.Status == models.VoteStatusClosed {
			now := time.Now()
			vote.ClosedAt = &now
		}
	}
	vote.UpdatedAt = time.Now()

	// Update in database
	if err := s.voteRepo.UpdateVote(ctx, vote); err != nil {
		utils.LogError(err, "Failed to update vote", map[string]interface{}{
			"vote_id": id,
			"user_id": userID,
		})
		return nil, fmt.Errorf("failed to update vote: %w", err)
	}

	utils.LogInfo("Vote updated successfully", map[string]interface{}{
		"vote_id": id,
		"user_id": userID,
	})

	return vote, nil
}

// DeleteVote deletes a vote
func (s *VoteService) DeleteVote(ctx context.Context, id string, userID string) error {
	// Get existing vote
	vote, err := s.voteRepo.GetVoteByID(ctx, id)
	if err != nil {
		return err
	}

	// Check if user is the creator
	if vote.CreatedBy != userID {
		return fmt.Errorf("unauthorized: only vote creator can delete vote")
	}

	// Delete vote (cascade will handle options and user votes)
	if err := s.voteRepo.DeleteVote(ctx, id); err != nil {
		utils.LogError(err, "Failed to delete vote", map[string]interface{}{
			"vote_id": id,
			"user_id": userID,
		})
		return fmt.Errorf("failed to delete vote: %w", err)
	}

	utils.LogInfo("Vote deleted successfully", map[string]interface{}{
		"vote_id": id,
		"user_id": userID,
	})

	return nil
}

// ListVotes lists votes with filters
func (s *VoteService) ListVotes(ctx context.Context, filters *models.VoteFilters) (*models.PaginatedVoteList, error) {
	// Set default pagination if not provided
	if filters.Limit == 0 {
		filters.Limit = 20 // Default page size
	}

	// Get votes with pagination info
	paginatedVotes, err := s.voteRepo.ListVotes(ctx, filters)
	if err != nil {
		utils.LogError(err, "Failed to list votes", nil)
		return nil, fmt.Errorf("failed to list votes: %w", err)
	}

	return paginatedVotes, nil
}

// CastVote allows a user to cast their vote
func (s *VoteService) CastVote(ctx context.Context, voteID string, req *models.VoteRequest, userID string) error {
	// Get vote to validate
	vote, err := s.voteRepo.GetVoteByID(ctx, voteID)
	if err != nil {
		return err
	}

	// Check if vote is active
	if vote.Status != models.VoteStatusActive {
		return fmt.Errorf("voting is only allowed on active votes")
	}

	// Check group membership if vote is assigned to groups
	voteGroups, err := s.groupService.GetVoteGroups(ctx, voteID)
	if err != nil {
		utils.LogError(err, "Failed to get vote groups", map[string]interface{}{
			"vote_id": voteID,
		})
		return fmt.Errorf("failed to check vote group assignments: %w", err)
	}

	// If vote is assigned to groups, check if user is a member of any assigned group
	if len(voteGroups) > 0 {
		userIsMember := false
		for _, group := range voteGroups {
			isMember, err := s.groupService.ValidateGroupAccess(ctx, group.ID, userID)
			if err != nil {
				utils.LogError(err, "Failed to check group membership", map[string]interface{}{
					"group_id": group.ID,
					"user_id":  userID,
				})
				continue
			}
			if isMember {
				userIsMember = true
				break
			}
		}

		if !userIsMember {
			return fmt.Errorf("voting is restricted to group members only")
		}
	}

	// Get vote options to validate selections
	options, err := s.voteRepo.GetVoteOptions(ctx, voteID)
	if err != nil {
		return fmt.Errorf("failed to get vote options: %w", err)
	}

	// Validate selected options
	if err := s.validateVoteRequest(req, options, vote.Type); err != nil {
		return err
	}

	// Check if user has already voted
	existingVote, err := s.voteRepo.GetUserVote(ctx, voteID, userID)
	if err == nil && existingVote != nil {
		// User has already voted - update their vote
		existingVote.SelectedOptions = req.SelectedOptions
		existingVote.VotedAt = time.Now()

		if err := s.voteRepo.UpdateUserVote(ctx, existingVote); err != nil {
			utils.LogError(err, "Failed to update user vote", map[string]interface{}{
				"vote_id": voteID,
				"user_id": userID,
			})
			return fmt.Errorf("failed to update vote: %w", err)
		}

		utils.LogInfo("Vote updated successfully", map[string]interface{}{
			"vote_id": voteID,
			"user_id": userID,
			"options": req.SelectedOptions,
		})

		return nil
	}

	// Create new user vote
	userVote := &models.UserVote{
		ID:              uuid.New().String(),
		VoteID:          voteID,
		UserID:          userID,
		SelectedOptions: req.SelectedOptions,
		VotedAt:         time.Now(),
	}

	if err := s.voteRepo.CreateUserVote(ctx, userVote); err != nil {
		utils.LogError(err, "Failed to cast vote", map[string]interface{}{
			"vote_id": voteID,
			"user_id": userID,
		})
		return fmt.Errorf("failed to cast vote: %w", err)
	}

	utils.LogInfo("Vote cast successfully", map[string]interface{}{
		"vote_id": voteID,
		"user_id": userID,
		"options": req.SelectedOptions,
	})

	return nil
}

// GetUserVote retrieves a user's vote for a specific vote
func (s *VoteService) GetUserVote(ctx context.Context, voteID string, userID string) (*models.UserVote, error) {
	userVote, err := s.voteRepo.GetUserVote(ctx, voteID, userID)
	if err != nil {
		return nil, err
	}

	return userVote, nil
}

// HasUserVoted checks if a user has voted for a specific vote
func (s *VoteService) HasUserVoted(ctx context.Context, voteID string, userID string) (bool, error) {
	hasVoted, err := s.voteRepo.HasUserVoted(ctx, voteID, userID)
	if err != nil {
		return false, err
	}

	return hasVoted, nil
}

// CloseVote closes a vote
func (s *VoteService) CloseVote(ctx context.Context, id string, userID string) error {
	req := &models.UpdateVoteRequest{
		Status: &[]models.VoteStatus{models.VoteStatusClosed}[0],
	}

	_, err := s.UpdateVote(ctx, id, req, userID)
	return err
}

// CancelVote cancels a vote
func (s *VoteService) CancelVote(ctx context.Context, id string, userID string) error {
	req := &models.UpdateVoteRequest{
		Status: &[]models.VoteStatus{models.VoteStatusCancelled}[0],
	}

	_, err := s.UpdateVote(ctx, id, req, userID)
	return err
}

// validateCreateVoteRequest validates the create vote request
func (s *VoteService) validateCreateVoteRequest(req *models.CreateVoteRequest) error {
	if req.Title == "" {
		return fmt.Errorf("title is required")
	}
	if len(req.Title) < 3 || len(req.Title) > 255 {
		return fmt.Errorf("title must be between 3 and 255 characters")
	}

	if req.Description == "" {
		return fmt.Errorf("description is required")
	}
	if len(req.Description) < 10 || len(req.Description) > 1000 {
		return fmt.Errorf("description must be between 10 and 1000 characters")
	}

	if req.Type != models.VoteTypeSingle && req.Type != models.VoteTypeMultiple {
		return fmt.Errorf("invalid vote type: must be 'single' or 'multiple'")
	}

	if len(req.Options) < 2 {
		return fmt.Errorf("at least 2 options are required")
	}
	if len(req.Options) > 20 {
		return fmt.Errorf("maximum 20 options allowed")
	}

	for i, option := range req.Options {
		if option == "" {
			return fmt.Errorf("option %d cannot be empty", i+1)
		}
		if len(option) > 255 {
			return fmt.Errorf("option %d must be less than 255 characters", i+1)
		}
	}

	return nil
}

// validateVoteRequest validates the vote request
func (s *VoteService) validateVoteRequest(req *models.VoteRequest, options []*models.VoteOption, voteType models.VoteType) error {
	if len(req.SelectedOptions) == 0 {
		return fmt.Errorf("at least one option must be selected")
	}

	// Create a map of valid option IDs
	validOptions := make(map[string]bool)
	for _, option := range options {
		validOptions[option.ID] = true
	}

	// Validate selected options
	for _, selectedOption := range req.SelectedOptions {
		if !validOptions[selectedOption] {
			return fmt.Errorf("invalid option selected: %s", selectedOption)
		}
	}

	// Validate vote type constraints
	if voteType == models.VoteTypeSingle && len(req.SelectedOptions) > 1 {
		return fmt.Errorf("only one option can be selected for single choice vote")
	}

	if voteType == models.VoteTypeMultiple && len(req.SelectedOptions) > len(options) {
		return fmt.Errorf("cannot select more options than available")
	}

	return nil
}
