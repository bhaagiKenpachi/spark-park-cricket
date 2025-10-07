package supabase

import (
	"context"
	"fmt"
	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/repository/interfaces"
	"spark-park-cricket-backend/internal/utils"
	"time"

	"github.com/google/uuid"
	"github.com/supabase-community/postgrest-go"
	"github.com/supabase-community/supabase-go"
)

// VoteRepository implements the VoteRepositoryInterface using Supabase
type VoteRepository struct {
	client *supabase.Client
}

// NewVoteRepository creates a new vote repository
func NewVoteRepository(client *supabase.Client) interfaces.VoteRepositoryInterface {
	return &VoteRepository{
		client: client,
	}
}

// CreateVote creates a new vote
func (r *VoteRepository) CreateVote(ctx context.Context, vote *models.Vote) error {
	vote.ID = uuid.New().String()
	vote.CreatedAt = time.Now()
	vote.UpdatedAt = time.Now()

	var result []models.Vote
	_, err := r.client.From("votes").Insert(vote, false, "", "", "").ExecuteTo(&result)
	if err != nil {
		utils.LogError(err, "Failed to create vote", map[string]interface{}{
			"vote_id": vote.ID,
			"title":   vote.Title,
		})
		return fmt.Errorf("failed to create vote: %w", err)
	}

	if len(result) > 0 {
		*vote = result[0]
	}

	return nil
}

// GetVoteByID retrieves a vote by ID
func (r *VoteRepository) GetVoteByID(ctx context.Context, id string) (*models.Vote, error) {
	var votes []models.Vote
	_, err := r.client.From("votes").Select("*", "", false).Eq("id", id).ExecuteTo(&votes)
	if err != nil {
		utils.LogError(err, "Failed to get vote by ID", map[string]interface{}{
			"vote_id": id,
		})
		return nil, fmt.Errorf("failed to get vote: %w", err)
	}

	if len(votes) == 0 {
		return nil, fmt.Errorf("vote not found")
	}

	return &votes[0], nil
}

// GetVoteWithOptions retrieves a vote with its options
func (r *VoteRepository) GetVoteWithOptions(ctx context.Context, id string) (*models.VoteWithOptions, error) {
	// Get vote
	vote, err := r.GetVoteByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Get options
	options, err := r.GetVoteOptions(ctx, id)
	if err != nil {
		return nil, err
	}

	return &models.VoteWithOptions{
		Vote:    vote,
		Options: options,
	}, nil
}

// GetVoteWithResults retrieves a vote with results and user's vote
func (r *VoteRepository) GetVoteWithResults(ctx context.Context, id string, userID string) (*models.VoteWithResults, error) {
	// Get vote with options
	voteWithOptions, err := r.GetVoteWithOptions(ctx, id)
	if err != nil {
		return nil, err
	}

	// Get results
	results, err := r.GetVoteResults(ctx, id)
	if err != nil {
		return nil, err
	}

	// Get user's vote
	var userVote *models.UserVote
	if userID != "" {
		userVote, _ = r.GetUserVote(ctx, id, userID)
	}

	// Get voted users
	votedUsers, err := r.GetVotedUsers(ctx, id)
	if err != nil {
		return nil, err
	}

	// Get total vote count
	totalVotes, err := r.GetTotalVoteCount(ctx, id)
	if err != nil {
		return nil, err
	}

	return &models.VoteWithResults{
		Vote:       voteWithOptions.Vote,
		Options:    voteWithOptions.Options,
		Results:    results,
		UserVote:   userVote,
		TotalVotes: totalVotes,
		VotedUsers: votedUsers,
	}, nil
}

// UpdateVote updates an existing vote
func (r *VoteRepository) UpdateVote(ctx context.Context, vote *models.Vote) error {
	vote.UpdatedAt = time.Now()

	var result []models.Vote
	_, err := r.client.From("votes").Update(vote, "", "").Eq("id", vote.ID).ExecuteTo(&result)
	if err != nil {
		utils.LogError(err, "Failed to update vote", map[string]interface{}{
			"vote_id": vote.ID,
		})
		return fmt.Errorf("failed to update vote: %w", err)
	}

	if len(result) == 0 {
		return fmt.Errorf("vote not found")
	}

	return nil
}

// DeleteVote deletes a vote
func (r *VoteRepository) DeleteVote(ctx context.Context, id string) error {
	_, err := r.client.From("votes").Delete("", "").Eq("id", id).ExecuteTo(nil)
	if err != nil {
		utils.LogError(err, "Failed to delete vote", map[string]interface{}{
			"vote_id": id,
		})
		return fmt.Errorf("failed to delete vote: %w", err)
	}

	return nil
}

// ListVotes lists votes with filters
func (r *VoteRepository) ListVotes(ctx context.Context, filters *models.VoteFilters) ([]*models.Vote, error) {
	query := r.client.From("votes").Select("*", "", false)

	// Add filters
	if filters.Status != nil {
		query = query.Eq("status", string(*filters.Status))
	}

	if filters.CreatedBy != nil {
		query = query.Eq("created_by", *filters.CreatedBy)
	}

	if filters.Type != nil {
		query = query.Eq("type", string(*filters.Type))
	}

	// Add ordering and pagination
	query = query.Order("created_at", &postgrest.OrderOpts{Ascending: false})

	if filters.Limit > 0 {
		query = query.Limit(filters.Limit, "")
	}

	if filters.Offset > 0 {
		query = query.Range(filters.Offset, filters.Offset+filters.Limit-1, "")
	}

	var votes []models.Vote
	_, err := query.ExecuteTo(&votes)
	if err != nil {
		utils.LogError(err, "Failed to list votes", nil)
		return nil, fmt.Errorf("failed to list votes: %w", err)
	}

	// Convert to slice of pointers
	votePointers := make([]*models.Vote, len(votes))
	for i := range votes {
		votePointers[i] = &votes[i]
	}

	return votePointers, nil
}

// CreateVoteOptions creates multiple vote options
func (r *VoteRepository) CreateVoteOptions(ctx context.Context, options []*models.VoteOption) error {
	if len(options) == 0 {
		return nil
	}

	// Convert to slice of models for insertion
	optionModels := make([]models.VoteOption, len(options))
	for i, option := range options {
		option.ID = uuid.New().String()
		option.CreatedAt = time.Now()
		option.UpdatedAt = time.Now()
		optionModels[i] = *option
	}

	var result []models.VoteOption
	_, err := r.client.From("vote_options").Insert(optionModels, false, "", "", "").ExecuteTo(&result)
	if err != nil {
		utils.LogError(err, "Failed to create vote options", map[string]interface{}{
			"vote_id": options[0].VoteID,
			"count":   len(options),
		})
		return fmt.Errorf("failed to create vote options: %w", err)
	}

	return nil
}

// GetVoteOptions retrieves options for a vote
func (r *VoteRepository) GetVoteOptions(ctx context.Context, voteID string) ([]*models.VoteOption, error) {
	var options []models.VoteOption
	_, err := r.client.From("vote_options").Select("*", "", false).Eq("vote_id", voteID).Order("created_at", &postgrest.OrderOpts{Ascending: true}).ExecuteTo(&options)
	if err != nil {
		utils.LogError(err, "Failed to get vote options", map[string]interface{}{
			"vote_id": voteID,
		})
		return nil, fmt.Errorf("failed to get vote options: %w", err)
	}

	// Convert to slice of pointers
	optionPointers := make([]*models.VoteOption, len(options))
	for i := range options {
		optionPointers[i] = &options[i]
	}

	return optionPointers, nil
}

// UpdateVoteOption updates a vote option
func (r *VoteRepository) UpdateVoteOption(ctx context.Context, option *models.VoteOption) error {
	option.UpdatedAt = time.Now()

	var result []models.VoteOption
	_, err := r.client.From("vote_options").Update(option, "", "").Eq("id", option.ID).ExecuteTo(&result)
	if err != nil {
		utils.LogError(err, "Failed to update vote option", map[string]interface{}{
			"option_id": option.ID,
		})
		return fmt.Errorf("failed to update vote option: %w", err)
	}

	if len(result) == 0 {
		return fmt.Errorf("vote option not found")
	}

	return nil
}

// DeleteVoteOption deletes a vote option
func (r *VoteRepository) DeleteVoteOption(ctx context.Context, id string) error {
	_, err := r.client.From("vote_options").Delete("", "").Eq("id", id).ExecuteTo(nil)
	if err != nil {
		utils.LogError(err, "Failed to delete vote option", map[string]interface{}{
			"option_id": id,
		})
		return fmt.Errorf("failed to delete vote option: %w", err)
	}

	return nil
}

// CreateUserVote creates a user's vote
func (r *VoteRepository) CreateUserVote(ctx context.Context, userVote *models.UserVote) error {
	userVote.ID = uuid.New().String()
	userVote.VotedAt = time.Now()

	var result []models.UserVote
	_, err := r.client.From("user_votes").Insert(userVote, false, "", "", "").ExecuteTo(&result)
	if err != nil {
		utils.LogError(err, "Failed to create user vote", map[string]interface{}{
			"vote_id": userVote.VoteID,
			"user_id": userVote.UserID,
		})
		return fmt.Errorf("failed to create user vote: %w", err)
	}

	if len(result) > 0 {
		*userVote = result[0]
	}

	return nil
}

// GetUserVote retrieves a user's vote for a specific vote
func (r *VoteRepository) GetUserVote(ctx context.Context, voteID, userID string) (*models.UserVote, error) {
	var userVotes []models.UserVote
	_, err := r.client.From("user_votes").Select("*", "", false).Eq("vote_id", voteID).Eq("user_id", userID).ExecuteTo(&userVotes)
	if err != nil {
		utils.LogError(err, "Failed to get user vote", map[string]interface{}{
			"vote_id": voteID,
			"user_id": userID,
		})
		return nil, fmt.Errorf("failed to get user vote: %w", err)
	}

	if len(userVotes) == 0 {
		return nil, fmt.Errorf("user vote not found")
	}

	return &userVotes[0], nil
}

// HasUserVoted checks if a user has voted for a specific vote
func (r *VoteRepository) HasUserVoted(ctx context.Context, voteID, userID string) (bool, error) {
	var userVotes []models.UserVote
	_, err := r.client.From("user_votes").Select("id", "", false).Eq("vote_id", voteID).Eq("user_id", userID).ExecuteTo(&userVotes)
	if err != nil {
		utils.LogError(err, "Failed to check if user voted", map[string]interface{}{
			"vote_id": voteID,
			"user_id": userID,
		})
		return false, fmt.Errorf("failed to check if user voted: %w", err)
	}

	return len(userVotes) > 0, nil
}

// GetVoteResults gets the vote results for a specific vote
func (r *VoteRepository) GetVoteResults(ctx context.Context, voteID string) (map[string]int, error) {
	// Get all vote options for this vote
	options, err := r.GetVoteOptions(ctx, voteID)
	if err != nil {
		return nil, err
	}

	// Get all user votes for this vote
	var userVotes []models.UserVote
	_, err = r.client.From("user_votes").Select("*", "", false).Eq("vote_id", voteID).ExecuteTo(&userVotes)
	if err != nil {
		utils.LogError(err, "Failed to get vote results", map[string]interface{}{
			"vote_id": voteID,
		})
		return nil, fmt.Errorf("failed to get vote results: %w", err)
	}

	// Count votes for each option
	results := make(map[string]int)
	for _, option := range options {
		results[option.ID] = 0
	}

	for _, userVote := range userVotes {
		for _, selectedOption := range userVote.SelectedOptions {
			if _, exists := results[selectedOption]; exists {
				results[selectedOption]++
			}
		}
	}

	return results, nil
}

// GetVotedUsers gets the list of users who voted for a specific vote
func (r *VoteRepository) GetVotedUsers(ctx context.Context, voteID string) ([]string, error) {
	var userVotes []models.UserVote
	_, err := r.client.From("user_votes").Select("user_id", "", false).Eq("vote_id", voteID).Order("voted_at", &postgrest.OrderOpts{Ascending: false}).ExecuteTo(&userVotes)
	if err != nil {
		utils.LogError(err, "Failed to get voted users", map[string]interface{}{
			"vote_id": voteID,
		})
		return nil, fmt.Errorf("failed to get voted users: %w", err)
	}

	userIDs := make([]string, len(userVotes))
	for i, userVote := range userVotes {
		userIDs[i] = userVote.UserID
	}

	return userIDs, nil
}

// GetTotalVoteCount gets the total number of votes for a specific vote
func (r *VoteRepository) GetTotalVoteCount(ctx context.Context, voteID string) (int, error) {
	var userVotes []models.UserVote
	_, err := r.client.From("user_votes").Select("id", "", false).Eq("vote_id", voteID).ExecuteTo(&userVotes)
	if err != nil {
		utils.LogError(err, "Failed to get total vote count", map[string]interface{}{
			"vote_id": voteID,
		})
		return 0, fmt.Errorf("failed to get total vote count: %w", err)
	}

	return len(userVotes), nil
}
