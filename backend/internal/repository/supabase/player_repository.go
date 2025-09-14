package supabase

import (
	"context"
	"fmt"
	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/repository/interfaces"

	"github.com/supabase-community/supabase-go"
)

type playerRepository struct {
	client *supabase.Client
	schema string
}

// NewPlayerRepository creates a new player repository
func NewPlayerRepository(client *supabase.Client, schema string) interfaces.PlayerRepository {
	return &playerRepository{
		client: client,
		schema: schema,
	}
}

func (r *playerRepository) Create(ctx context.Context, player *models.Player) error {
	_, err := r.client.From("players").Insert(player, false, "", "", "").ExecuteTo(&player)
	return err
}

func (r *playerRepository) GetByID(ctx context.Context, id string) (*models.Player, error) {
	var result []models.Player
	_, err := r.client.From("players").Select("*", "", false).Eq("id", id).ExecuteTo(&result)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("player not found")
	}
	return &result[0], nil
}

func (r *playerRepository) GetAll(ctx context.Context, filters *models.PlayerFilters) ([]*models.Player, error) {
	var result []models.Player
	tableName := fmt.Sprintf("%s.players", r.schema)
	query := r.client.From(tableName).Select("*", "", false)

	if filters != nil {
		if filters.Limit > 0 {
			query = query.Limit(filters.Limit, "")
		}
		// Note: Offset is not supported by this Supabase client version
		// Use Range method for pagination if needed
		if filters.TeamID != nil && *filters.TeamID != "" {
			query = query.Eq("team_id", *filters.TeamID)
		}
	}

	_, err := query.ExecuteTo(&result)
	if err != nil {
		return nil, err
	}

	// Convert to slice of pointers
	players := make([]*models.Player, len(result))
	for i := range result {
		players[i] = &result[i]
	}
	return players, nil
}

func (r *playerRepository) Update(ctx context.Context, id string, player *models.Player) error {
	_, err := r.client.From("players").Update(player, "", "").Eq("id", id).ExecuteTo(&player)
	return err
}

func (r *playerRepository) Delete(ctx context.Context, id string) error {
	_, err := r.client.From("players").Delete("", "").Eq("id", id).ExecuteTo(nil)
	return err
}

func (r *playerRepository) GetByTeamID(ctx context.Context, teamID string) ([]*models.Player, error) {
	var result []models.Player
	_, err := r.client.From("players").Select("*", "", false).Eq("team_id", teamID).ExecuteTo(&result)
	if err != nil {
		return nil, err
	}

	// Convert to slice of pointers
	players := make([]*models.Player, len(result))
	for i := range result {
		players[i] = &result[i]
	}
	return players, nil
}

func (r *playerRepository) Count(ctx context.Context) (int64, error) {
	var result []models.Player
	_, err := r.client.From("players").Select("*", "", false).ExecuteTo(&result)
	if err != nil {
		return 0, err
	}
	return int64(len(result)), nil
}
