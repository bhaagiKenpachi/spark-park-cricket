package supabase

import (
	"context"
	"fmt"
	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/repository/interfaces"

	"github.com/supabase-community/supabase-go"
)

type teamRepository struct {
	client *supabase.Client
}

// NewTeamRepository creates a new team repository
func NewTeamRepository(client *supabase.Client) interfaces.TeamRepository {
	return &teamRepository{
		client: client,
	}
}

func (r *teamRepository) Create(ctx context.Context, team *models.Team) error {
	_, err := r.client.From("teams").Insert(team, false, "", "", "").ExecuteTo(&team)
	return err
}

func (r *teamRepository) GetByID(ctx context.Context, id string) (*models.Team, error) {
	var result []models.Team
	_, err := r.client.From("teams").Select("*", "", false).Eq("id", id).ExecuteTo(&result)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("team not found")
	}
	return &result[0], nil
}

func (r *teamRepository) GetAll(ctx context.Context, filters *models.TeamFilters) ([]*models.Team, error) {
	var result []models.Team
	query := r.client.From("teams").Select("*", "", false)

	if filters != nil {
		if filters.Limit > 0 {
			query = query.Limit(filters.Limit, "")
		}
		// Note: Offset is not supported by this Supabase client version
		// Use Range method for pagination if needed
	}

	_, err := query.ExecuteTo(&result)
	if err != nil {
		return nil, err
	}

	// Convert to slice of pointers
	teams := make([]*models.Team, len(result))
	for i := range result {
		teams[i] = &result[i]
	}
	return teams, nil
}

func (r *teamRepository) Update(ctx context.Context, id string, team *models.Team) error {
	_, err := r.client.From("teams").Update(team, "", "").Eq("id", id).ExecuteTo(&team)
	return err
}

func (r *teamRepository) Delete(ctx context.Context, id string) error {
	_, err := r.client.From("teams").Delete("", "").Eq("id", id).ExecuteTo(nil)
	return err
}

func (r *teamRepository) Count(ctx context.Context) (int64, error) {
	var result []models.Team
	_, err := r.client.From("teams").Select("*", "", false).ExecuteTo(&result)
	if err != nil {
		return 0, err
	}
	return int64(len(result)), nil
}
