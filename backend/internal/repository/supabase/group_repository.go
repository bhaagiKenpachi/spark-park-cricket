package supabase

import (
	"context"
	"fmt"
	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/repository/interfaces"
	"spark-park-cricket-backend/internal/utils"

	"github.com/supabase-community/supabase-go"
)

type groupRepository struct {
	client *supabase.Client
}

// NewGroupRepository creates a new group repository
func NewGroupRepository(client *supabase.Client) interfaces.GroupRepository {
	return &groupRepository{
		client: client,
	}
}

// CreateGroup creates a new group
func (r *groupRepository) CreateGroup(ctx context.Context, group *models.Group) error {
	var result []models.Group
	_, err := r.client.From("groups").Insert(group, false, "", "", "").ExecuteTo(&result)
	if err != nil {
		utils.LogError(err, "Failed to create group", map[string]interface{}{
			"group_id": group.ID,
			"name":     group.Name,
		})
		return fmt.Errorf("failed to create group: %w", err)
	}

	utils.LogInfo("Group created successfully", map[string]interface{}{
		"group_id": group.ID,
		"name":     group.Name,
		"type":     group.Type,
	})

	return nil
}

// GetGroup retrieves a group by ID
func (r *groupRepository) GetGroup(ctx context.Context, groupID string) (*models.Group, error) {
	var groups []models.Group
	_, err := r.client.From("groups").Select("*", "", false).Eq("id", groupID).ExecuteTo(&groups)
	if err != nil {
		utils.LogError(err, "Failed to get group", map[string]interface{}{
			"group_id": groupID,
		})
		return nil, fmt.Errorf("failed to get group: %w", err)
	}

	if len(groups) == 0 {
		return nil, fmt.Errorf("group not found")
	}

	return &groups[0], nil
}

// GetGroupWithCreator retrieves a group with creator information
func (r *groupRepository) GetGroupWithCreator(ctx context.Context, groupID string) (*models.GroupWithCreator, error) {
	// For now, get the group and creator separately
	group, err := r.GetGroup(ctx, groupID)
	if err != nil {
		return nil, err
	}

	// Get creator info
	var users []models.User
	_, err = r.client.From("users").Select("name", "", false).Eq("id", group.CreatedBy).ExecuteTo(&users)
	if err != nil || len(users) == 0 {
		return &models.GroupWithCreator{
			Group:       *group,
			CreatorName: "Unknown",
		}, nil
	}

	return &models.GroupWithCreator{
		Group:       *group,
		CreatorName: users[0].Name,
	}, nil
}

// GetGroupWithMembers retrieves a group with its members
func (r *groupRepository) GetGroupWithMembers(ctx context.Context, groupID string) (*models.GroupWithMembers, error) {
	group, err := r.GetGroupWithCreator(ctx, groupID)
	if err != nil {
		return nil, err
	}

	members, err := r.GetGroupMembers(ctx, groupID)
	if err != nil {
		return nil, err
	}

	return &models.GroupWithMembers{
		Group:       group.Group,
		Members:     members,
		MemberCount: len(members),
		CreatorName: group.CreatorName,
	}, nil
}

// GetGroupWithMembersAndUsers retrieves a group with members and their user details
func (r *groupRepository) GetGroupWithMembersAndUsers(ctx context.Context, groupID string) (*models.GroupWithMembersAndUsers, error) {
	group, err := r.GetGroupWithCreator(ctx, groupID)
	if err != nil {
		return nil, err
	}

	membersWithUsers, err := r.GetGroupMembersWithUsers(ctx, groupID)
	if err != nil {
		return nil, err
	}

	return &models.GroupWithMembersAndUsers{
		Group:       group.Group,
		Members:     membersWithUsers,
		MemberCount: len(membersWithUsers),
		CreatorName: group.CreatorName,
	}, nil
}

// ListGroups retrieves a list of groups with pagination
func (r *groupRepository) ListGroups(ctx context.Context, limit, offset int) ([]*models.GroupWithCreator, error) {
	var groups []models.Group
	query := r.client.From("groups").Select("*", "", false)

	if limit > 0 {
		query = query.Limit(limit, "")
	}

	_, err := query.ExecuteTo(&groups)
	if err != nil {
		utils.LogError(err, "Failed to list groups", map[string]interface{}{
			"limit":  limit,
			"offset": offset,
		})
		return nil, fmt.Errorf("failed to list groups: %w", err)
	}

	// Convert to GroupWithCreator
	result := make([]*models.GroupWithCreator, len(groups))
	for i, group := range groups {
		// Get creator name
		var users []models.User
		_, err = r.client.From("users").Select("name", "", false).Eq("id", group.CreatedBy).ExecuteTo(&users)
		creatorName := "Unknown"
		if err == nil && len(users) > 0 {
			creatorName = users[0].Name
		}

		result[i] = &models.GroupWithCreator{
			Group:       group,
			CreatorName: creatorName,
		}
	}

	return result, nil
}

// ListGroupsByType retrieves groups filtered by type
func (r *groupRepository) ListGroupsByType(ctx context.Context, groupType models.GroupType, limit, offset int) ([]*models.GroupWithCreator, error) {
	var groups []models.Group
	query := r.client.From("groups").Select("*", "", false).Eq("type", string(groupType))

	if limit > 0 {
		query = query.Limit(limit, "")
	}

	_, err := query.ExecuteTo(&groups)
	if err != nil {
		utils.LogError(err, "Failed to list groups by type", map[string]interface{}{
			"type":   groupType,
			"limit":  limit,
			"offset": offset,
		})
		return nil, fmt.Errorf("failed to list groups by type: %w", err)
	}

	// Convert to GroupWithCreator
	result := make([]*models.GroupWithCreator, len(groups))
	for i, group := range groups {
		// Get creator name
		var users []models.User
		_, err = r.client.From("users").Select("name", "", false).Eq("id", group.CreatedBy).ExecuteTo(&users)
		creatorName := "Unknown"
		if err == nil && len(users) > 0 {
			creatorName = users[0].Name
		}

		result[i] = &models.GroupWithCreator{
			Group:       group,
			CreatorName: creatorName,
		}
	}

	return result, nil
}

// ListGroupsByCreator retrieves groups created by a specific user
func (r *groupRepository) ListGroupsByCreator(ctx context.Context, creatorID string, limit, offset int) ([]*models.GroupWithCreator, error) {
	var groups []models.Group
	query := r.client.From("groups").Select("*", "", false).Eq("created_by", creatorID)

	if limit > 0 {
		query = query.Limit(limit, "")
	}

	_, err := query.ExecuteTo(&groups)
	if err != nil {
		utils.LogError(err, "Failed to list groups by creator", map[string]interface{}{
			"creator_id": creatorID,
			"limit":      limit,
			"offset":     offset,
		})
		return nil, fmt.Errorf("failed to list groups by creator: %w", err)
	}

	// Convert to GroupWithCreator
	result := make([]*models.GroupWithCreator, len(groups))
	for i, group := range groups {
		// Get creator name
		var users []models.User
		_, err = r.client.From("users").Select("name", "", false).Eq("id", group.CreatedBy).ExecuteTo(&users)
		creatorName := "Unknown"
		if err == nil && len(users) > 0 {
			creatorName = users[0].Name
		}

		result[i] = &models.GroupWithCreator{
			Group:       group,
			CreatorName: creatorName,
		}
	}

	return result, nil
}

// ListGroupsByUser retrieves groups that a user is a member of
func (r *groupRepository) ListGroupsByUser(ctx context.Context, userID string, limit, offset int) ([]*models.GroupWithCreator, error) {
	// Get group IDs from group_members table
	var members []models.GroupMember
	query := r.client.From("group_members").Select("group_id", "", false).Eq("user_id", userID)

	if limit > 0 {
		query = query.Limit(limit, "")
	}

	_, err := query.ExecuteTo(&members)
	if err != nil {
		utils.LogError(err, "Failed to list groups by user", map[string]interface{}{
			"user_id": userID,
			"limit":   limit,
			"offset":  offset,
		})
		return nil, fmt.Errorf("failed to list groups by user: %w", err)
	}

	if len(members) == 0 {
		return []*models.GroupWithCreator{}, nil
	}

	// Get groups for these IDs
	groupIDs := make([]string, len(members))
	for i, member := range members {
		groupIDs[i] = member.GroupID
	}

	var groups []models.Group
	_, err = r.client.From("groups").Select("*", "", false).In("id", groupIDs).ExecuteTo(&groups)
	if err != nil {
		return nil, fmt.Errorf("failed to get groups: %w", err)
	}

	// Convert to GroupWithCreator
	result := make([]*models.GroupWithCreator, len(groups))
	for i, group := range groups {
		// Get creator name
		var users []models.User
		_, err = r.client.From("users").Select("name", "", false).Eq("id", group.CreatedBy).ExecuteTo(&users)
		creatorName := "Unknown"
		if err == nil && len(users) > 0 {
			creatorName = users[0].Name
		}

		result[i] = &models.GroupWithCreator{
			Group:       group,
			CreatorName: creatorName,
		}
	}

	return result, nil
}

// UpdateGroup updates an existing group
func (r *groupRepository) UpdateGroup(ctx context.Context, group *models.Group) error {
	_, err := r.client.From("groups").Update(group, "", "").Eq("id", group.ID).ExecuteTo(&group)
	if err != nil {
		utils.LogError(err, "Failed to update group", map[string]interface{}{
			"group_id": group.ID,
		})
		return fmt.Errorf("failed to update group: %w", err)
	}

	utils.LogInfo("Group updated successfully", map[string]interface{}{
		"group_id": group.ID,
	})

	return nil
}

// DeleteGroup deletes a group
func (r *groupRepository) DeleteGroup(ctx context.Context, groupID string) error {
	_, err := r.client.From("groups").Delete("", "").Eq("id", groupID).ExecuteTo(&[]models.Group{})
	if err != nil {
		utils.LogError(err, "Failed to delete group", map[string]interface{}{
			"group_id": groupID,
		})
		return fmt.Errorf("failed to delete group: %w", err)
	}

	utils.LogInfo("Group deleted successfully", map[string]interface{}{
		"group_id": groupID,
	})

	return nil
}

// AddGroupMember adds a user to a group
func (r *groupRepository) AddGroupMember(ctx context.Context, member *models.GroupMember) error {
	var result []models.GroupMember
	_, err := r.client.From("group_members").Insert(member, false, "", "", "").ExecuteTo(&result)
	if err != nil {
		utils.LogError(err, "Failed to add group member", map[string]interface{}{
			"group_id": member.GroupID,
			"user_id":  member.UserID,
		})
		return fmt.Errorf("failed to add group member: %w", err)
	}

	utils.LogInfo("Group member added successfully", map[string]interface{}{
		"group_id": member.GroupID,
		"user_id":  member.UserID,
		"role":     member.Role,
	})

	return nil
}

// RemoveGroupMember removes a user from a group
func (r *groupRepository) RemoveGroupMember(ctx context.Context, groupID, userID string) error {
	_, err := r.client.From("group_members").Delete("", "").Eq("group_id", groupID).Eq("user_id", userID).ExecuteTo(&[]models.GroupMember{})
	if err != nil {
		utils.LogError(err, "Failed to remove group member", map[string]interface{}{
			"group_id": groupID,
			"user_id":  userID,
		})
		return fmt.Errorf("failed to remove group member: %w", err)
	}

	utils.LogInfo("Group member removed successfully", map[string]interface{}{
		"group_id": groupID,
		"user_id":  userID,
	})

	return nil
}

// GetGroupMember retrieves a specific group member
func (r *groupRepository) GetGroupMember(ctx context.Context, groupID, userID string) (*models.GroupMember, error) {
	var members []models.GroupMember
	_, err := r.client.From("group_members").Select("*", "", false).Eq("group_id", groupID).Eq("user_id", userID).ExecuteTo(&members)
	if err != nil {
		utils.LogError(err, "Failed to get group member", map[string]interface{}{
			"group_id": groupID,
			"user_id":  userID,
		})
		return nil, fmt.Errorf("failed to get group member: %w", err)
	}

	if len(members) == 0 {
		return nil, fmt.Errorf("group member not found")
	}

	return &members[0], nil
}

// GetGroupMembers retrieves all members of a group
func (r *groupRepository) GetGroupMembers(ctx context.Context, groupID string) ([]*models.GroupMember, error) {
	var members []models.GroupMember
	_, err := r.client.From("group_members").Select("*", "", false).Eq("group_id", groupID).ExecuteTo(&members)
	if err != nil {
		utils.LogError(err, "Failed to get group members", map[string]interface{}{
			"group_id": groupID,
		})
		return nil, fmt.Errorf("failed to get group members: %w", err)
	}

	// Convert to pointers
	result := make([]*models.GroupMember, len(members))
	for i := range members {
		result[i] = &members[i]
	}

	return result, nil
}

// GetGroupMembersWithUsers retrieves group members with their user details
func (r *groupRepository) GetGroupMembersWithUsers(ctx context.Context, groupID string) ([]*models.GroupMemberWithUser, error) {
	var members []models.GroupMember
	_, err := r.client.From("group_members").Select("*", "", false).Eq("group_id", groupID).ExecuteTo(&members)
	if err != nil {
		utils.LogError(err, "Failed to get group members with users", map[string]interface{}{
			"group_id": groupID,
		})
		return nil, fmt.Errorf("failed to get group members with users: %w", err)
	}

	// Get user details for each member
	result := make([]*models.GroupMemberWithUser, len(members))
	for i, member := range members {
		var users []models.User
		_, err = r.client.From("users").Select("name,email", "", false).Eq("id", member.UserID).ExecuteTo(&users)

		userName := "Unknown"
		userEmail := ""
		if err == nil && len(users) > 0 {
			userName = users[0].Name
			userEmail = users[0].Email
		}

		result[i] = &models.GroupMemberWithUser{
			GroupMember: member,
			UserName:    userName,
			UserEmail:   userEmail,
		}
	}

	return result, nil
}

// UpdateGroupMemberRole updates a group member's role
func (r *groupRepository) UpdateGroupMemberRole(ctx context.Context, groupID, userID, role string) error {
	updateData := map[string]interface{}{
		"role": role,
	}

	_, err := r.client.From("group_members").Update(updateData, "", "").Eq("group_id", groupID).Eq("user_id", userID).ExecuteTo(&[]models.GroupMember{})
	if err != nil {
		utils.LogError(err, "Failed to update group member role", map[string]interface{}{
			"group_id": groupID,
			"user_id":  userID,
			"role":     role,
		})
		return fmt.Errorf("failed to update group member role: %w", err)
	}

	utils.LogInfo("Group member role updated successfully", map[string]interface{}{
		"group_id": groupID,
		"user_id":  userID,
		"role":     role,
	})

	return nil
}

// IsUserInGroup checks if a user is a member of a group
func (r *groupRepository) IsUserInGroup(ctx context.Context, groupID, userID string) (bool, error) {
	var members []models.GroupMember
	_, err := r.client.From("group_members").Select("id", "", false).Eq("group_id", groupID).Eq("user_id", userID).ExecuteTo(&members)
	if err != nil {
		utils.LogError(err, "Failed to check if user is in group", map[string]interface{}{
			"group_id": groupID,
			"user_id":  userID,
		})
		return false, fmt.Errorf("failed to check if user is in group: %w", err)
	}

	return len(members) > 0, nil
}

// GetUserGroups retrieves all groups that a user is a member of
func (r *groupRepository) GetUserGroups(ctx context.Context, userID string) ([]*models.Group, error) {
	// Get group IDs from group_members table
	var members []models.GroupMember
	_, err := r.client.From("group_members").Select("group_id", "", false).Eq("user_id", userID).ExecuteTo(&members)
	if err != nil {
		utils.LogError(err, "Failed to get user groups", map[string]interface{}{
			"user_id": userID,
		})
		return nil, fmt.Errorf("failed to get user groups: %w", err)
	}

	if len(members) == 0 {
		return []*models.Group{}, nil
	}

	// Get groups for these IDs
	groupIDs := make([]string, len(members))
	for i, member := range members {
		groupIDs[i] = member.GroupID
	}

	var groups []models.Group
	_, err = r.client.From("groups").Select("*", "", false).In("id", groupIDs).ExecuteTo(&groups)
	if err != nil {
		return nil, fmt.Errorf("failed to get groups: %w", err)
	}

	// Convert to pointers
	result := make([]*models.Group, len(groups))
	for i := range groups {
		result[i] = &groups[i]
	}

	return result, nil
}

// AssignGroupsToVote assigns groups to a vote
func (r *groupRepository) AssignGroupsToVote(ctx context.Context, voteID string, groupIDs []string) error {
	// Create vote-group associations
	voteGroups := make([]models.VoteGroup, len(groupIDs))
	for i, groupID := range groupIDs {
		voteGroups[i] = models.VoteGroup{
			ID:      fmt.Sprintf("%s-%s", voteID, groupID), // Simple ID generation
			VoteID:  voteID,
			GroupID: groupID,
		}
	}

	var result []models.VoteGroup
	_, err := r.client.From("vote_groups").Insert(voteGroups, false, "", "", "").ExecuteTo(&result)
	if err != nil {
		utils.LogError(err, "Failed to assign groups to vote", map[string]interface{}{
			"vote_id":   voteID,
			"group_ids": groupIDs,
		})
		return fmt.Errorf("failed to assign groups to vote: %w", err)
	}

	utils.LogInfo("Groups assigned to vote successfully", map[string]interface{}{
		"vote_id":   voteID,
		"group_ids": groupIDs,
	})

	return nil
}

// RemoveGroupsFromVote removes groups from a vote
func (r *groupRepository) RemoveGroupsFromVote(ctx context.Context, voteID string, groupIDs []string) error {
	_, err := r.client.From("vote_groups").Delete("", "").Eq("vote_id", voteID).In("group_id", groupIDs).ExecuteTo(&[]models.VoteGroup{})
	if err != nil {
		utils.LogError(err, "Failed to remove groups from vote", map[string]interface{}{
			"vote_id":   voteID,
			"group_ids": groupIDs,
		})
		return fmt.Errorf("failed to remove groups from vote: %w", err)
	}

	utils.LogInfo("Groups removed from vote successfully", map[string]interface{}{
		"vote_id":   voteID,
		"group_ids": groupIDs,
	})

	return nil
}

// GetVoteGroups retrieves all groups assigned to a vote
func (r *groupRepository) GetVoteGroups(ctx context.Context, voteID string) ([]*models.Group, error) {
	// Get group IDs from vote_groups table
	var voteGroups []models.VoteGroup
	_, err := r.client.From("vote_groups").Select("group_id", "", false).Eq("vote_id", voteID).ExecuteTo(&voteGroups)
	if err != nil {
		utils.LogError(err, "Failed to get vote groups", map[string]interface{}{
			"vote_id": voteID,
		})
		return nil, fmt.Errorf("failed to get vote groups: %w", err)
	}

	if len(voteGroups) == 0 {
		return []*models.Group{}, nil
	}

	// Get groups for these IDs
	groupIDs := make([]string, len(voteGroups))
	for i, vg := range voteGroups {
		groupIDs[i] = vg.GroupID
	}

	var groups []models.Group
	_, err = r.client.From("groups").Select("*", "", false).In("id", groupIDs).ExecuteTo(&groups)
	if err != nil {
		return nil, fmt.Errorf("failed to get groups: %w", err)
	}

	// Convert to pointers
	result := make([]*models.Group, len(groups))
	for i := range groups {
		result[i] = &groups[i]
	}

	return result, nil
}

// GetGroupVotes retrieves all votes assigned to a group
func (r *groupRepository) GetGroupVotes(ctx context.Context, groupID string) ([]*models.Vote, error) {
	// Get vote IDs from vote_groups table
	var voteGroups []models.VoteGroup
	_, err := r.client.From("vote_groups").Select("vote_id", "", false).Eq("group_id", groupID).ExecuteTo(&voteGroups)
	if err != nil {
		utils.LogError(err, "Failed to get group votes", map[string]interface{}{
			"group_id": groupID,
		})
		return nil, fmt.Errorf("failed to get group votes: %w", err)
	}

	if len(voteGroups) == 0 {
		return []*models.Vote{}, nil
	}

	// Get votes for these IDs
	voteIDs := make([]string, len(voteGroups))
	for i, vg := range voteGroups {
		voteIDs[i] = vg.VoteID
	}

	var votes []models.Vote
	_, err = r.client.From("votes").Select("*", "", false).In("id", voteIDs).ExecuteTo(&votes)
	if err != nil {
		return nil, fmt.Errorf("failed to get votes: %w", err)
	}

	// Convert to pointers
	result := make([]*models.Vote, len(votes))
	for i := range votes {
		result[i] = &votes[i]
	}

	return result, nil
}

// IsGroupAssignedToVote checks if a group is assigned to a vote
func (r *groupRepository) IsGroupAssignedToVote(ctx context.Context, voteID, groupID string) (bool, error) {
	var voteGroups []models.VoteGroup
	_, err := r.client.From("vote_groups").Select("id", "", false).Eq("vote_id", voteID).Eq("group_id", groupID).ExecuteTo(&voteGroups)
	if err != nil {
		utils.LogError(err, "Failed to check if group is assigned to vote", map[string]interface{}{
			"vote_id":  voteID,
			"group_id": groupID,
		})
		return false, fmt.Errorf("failed to check if group is assigned to vote: %w", err)
	}

	return len(voteGroups) > 0, nil
}

// GetGroupVoteResults retrieves voting results for a specific group
func (r *groupRepository) GetGroupVoteResults(ctx context.Context, voteID, groupID string) (*models.GroupVoteResult, error) {
	// Get group name
	group, err := r.GetGroup(ctx, groupID)
	if err != nil {
		return nil, err
	}

	// Get vote options
	var options []models.VoteOption
	_, err = r.client.From("vote_options").Select("*", "", false).Eq("vote_id", voteID).ExecuteTo(&options)
	if err != nil {
		return nil, fmt.Errorf("failed to get vote options: %w", err)
	}

	// Get user votes for this group
	var userVotes []models.UserVote
	_, err = r.client.From("user_votes").Select("*", "", false).Eq("vote_id", voteID).ExecuteTo(&userVotes)
	if err != nil {
		return nil, fmt.Errorf("failed to get user votes: %w", err)
	}

	// Filter votes by group members
	var groupMembers []models.GroupMember
	_, err = r.client.From("group_members").Select("user_id", "", false).Eq("group_id", groupID).ExecuteTo(&groupMembers)
	if err != nil {
		return nil, fmt.Errorf("failed to get group members: %w", err)
	}

	memberUserIDs := make(map[string]bool)
	for _, member := range groupMembers {
		memberUserIDs[member.UserID] = true
	}

	// Calculate results
	results := make(map[string]int)
	resultsWithNames := make(map[string][]models.VoterInfo)
	votedMembers := make([]string, 0)

	for _, vote := range userVotes {
		if memberUserIDs[vote.UserID] {
			votedMembers = append(votedMembers, vote.UserID)
			for _, optionID := range vote.SelectedOptions {
				results[optionID]++

				// Get voter info
				var users []models.User
				_, err = r.client.From("users").Select("name", "", false).Eq("id", vote.UserID).ExecuteTo(&users)
				if err == nil && len(users) > 0 {
					voterInfo := models.VoterInfo{
						UserID:   vote.UserID,
						UserName: users[0].Name,
						VotedAt:  vote.VotedAt.Format("2006-01-02 15:04:05"),
					}
					resultsWithNames[optionID] = append(resultsWithNames[optionID], voterInfo)
				}
			}
		}
	}

	return &models.GroupVoteResult{
		GroupID:          groupID,
		GroupName:        group.Name,
		TotalVotes:       len(votedMembers),
		Results:          results,
		ResultsWithNames: resultsWithNames,
		VotedMembers:     votedMembers,
	}, nil
}

// GetVoteResultsByGroups retrieves voting results broken down by groups
func (r *groupRepository) GetVoteResultsByGroups(ctx context.Context, voteID string) ([]*models.GroupVoteResult, error) {
	// Get all groups assigned to this vote
	groups, err := r.GetVoteGroups(ctx, voteID)
	if err != nil {
		return nil, err
	}

	results := make([]*models.GroupVoteResult, len(groups))
	for i, group := range groups {
		groupResult, err := r.GetGroupVoteResults(ctx, voteID, group.ID)
		if err != nil {
			return nil, err
		}
		results[i] = groupResult
	}

	return results, nil
}

// GetVoteWithGroupResults retrieves a vote with results broken down by groups
func (r *groupRepository) GetVoteWithGroupResults(ctx context.Context, voteID string) (*models.VoteWithGroupResults, error) {
	// Get vote
	var votes []models.Vote
	_, err := r.client.From("votes").Select("*", "", false).Eq("id", voteID).ExecuteTo(&votes)
	if err != nil || len(votes) == 0 {
		return nil, fmt.Errorf("failed to get vote: %w", err)
	}
	vote := votes[0]

	// Get creator name
	var users []models.User
	_, err = r.client.From("users").Select("name", "", false).Eq("id", vote.CreatedBy).ExecuteTo(&users)
	creatorName := "Unknown"
	if err == nil && len(users) > 0 {
		creatorName = users[0].Name
	}

	// Get vote options
	var options []models.VoteOption
	_, err = r.client.From("vote_options").Select("*", "", false).Eq("vote_id", voteID).ExecuteTo(&options)
	if err != nil {
		return nil, fmt.Errorf("failed to get vote options: %w", err)
	}

	// Get group results
	groupResults, err := r.GetVoteResultsByGroups(ctx, voteID)
	if err != nil {
		return nil, err
	}

	// Calculate overall results
	overallResults := make(map[string]int)
	totalVotes := 0
	for _, groupResult := range groupResults {
		totalVotes += groupResult.TotalVotes
		for optionID, count := range groupResult.Results {
			overallResults[optionID] += count
		}
	}

	// Convert options to pointers
	optionPointers := make([]*models.VoteOption, len(options))
	for i := range options {
		optionPointers[i] = &options[i]
	}

	return &models.VoteWithGroupResults{
		Vote:           vote,
		Options:        optionPointers,
		GroupResults:   groupResults,
		OverallResults: overallResults,
		TotalVotes:     totalVotes,
		CreatorName:    creatorName,
	}, nil
}

// GetGroupMemberCount retrieves the number of members in a group
func (r *groupRepository) GetGroupMemberCount(ctx context.Context, groupID string) (int, error) {
	var members []models.GroupMember
	_, err := r.client.From("group_members").Select("id", "", false).Eq("group_id", groupID).ExecuteTo(&members)
	if err != nil {
		utils.LogError(err, "Failed to get group member count", map[string]interface{}{
			"group_id": groupID,
		})
		return 0, fmt.Errorf("failed to get group member count: %w", err)
	}

	return len(members), nil
}

// GetGroupVoteCount retrieves the number of votes assigned to a group
func (r *groupRepository) GetGroupVoteCount(ctx context.Context, groupID string) (int, error) {
	var voteGroups []models.VoteGroup
	_, err := r.client.From("vote_groups").Select("id", "", false).Eq("group_id", groupID).ExecuteTo(&voteGroups)
	if err != nil {
		utils.LogError(err, "Failed to get group vote count", map[string]interface{}{
			"group_id": groupID,
		})
		return 0, fmt.Errorf("failed to get group vote count: %w", err)
	}

	return len(voteGroups), nil
}

// SearchGroups searches for groups by name or description
func (r *groupRepository) SearchGroups(ctx context.Context, query string, limit, offset int) ([]*models.GroupWithCreator, error) {
	var groups []models.Group
	searchQuery := r.client.From("groups").Select("*", "", false)

	// Note: Supabase doesn't support ILIKE in the Go client, so we'll use a simple approach
	// In a real implementation, you might want to use a full-text search or implement this differently

	if limit > 0 {
		searchQuery = searchQuery.Limit(limit, "")
	}

	_, err := searchQuery.ExecuteTo(&groups)
	if err != nil {
		utils.LogError(err, "Failed to search groups", map[string]interface{}{
			"query":  query,
			"limit":  limit,
			"offset": offset,
		})
		return nil, fmt.Errorf("failed to search groups: %w", err)
	}

	// Filter results by query (client-side filtering for now)
	filteredGroups := make([]models.Group, 0)
	for _, group := range groups {
		if contains(group.Name, query) || contains(group.Description, query) {
			filteredGroups = append(filteredGroups, group)
		}
	}

	// Convert to GroupWithCreator
	result := make([]*models.GroupWithCreator, len(filteredGroups))
	for i, group := range filteredGroups {
		// Get creator name
		var users []models.User
		_, err = r.client.From("users").Select("name", "", false).Eq("id", group.CreatedBy).ExecuteTo(&users)
		creatorName := "Unknown"
		if err == nil && len(users) > 0 {
			creatorName = users[0].Name
		}

		result[i] = &models.GroupWithCreator{
			Group:       group,
			CreatorName: creatorName,
		}
	}

	return result, nil
}

// Helper function for case-insensitive string contains
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) && (s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr ||
			containsSubstring(s, substr))))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
