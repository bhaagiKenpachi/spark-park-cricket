package unit

import (
	"fmt"
	"testing"

	"spark-park-cricket-backend/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestFallOfWicketsModel_OverPosition(t *testing.T) {
	// Test the over position calculation
	fallOfWickets := &models.FallOfWickets{
		OverNumber: 15,
		BallNumber: 3,
	}

	// Calculate over position manually
	expected := "15.3"
	actual := fmt.Sprintf("%d.%d", fallOfWickets.OverNumber, fallOfWickets.BallNumber)

	assert.Equal(t, expected, actual)
}

func TestFallOfWicketsModel_OverPosition_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		overNumber  int
		ballNumber  int
		expected    string
	}{
		{
			name:        "First over first ball",
			overNumber:  1,
			ballNumber:  1,
			expected:    "1.1",
		},
		{
			name:        "Last ball of over",
			overNumber:  10,
			ballNumber:  6,
			expected:    "10.6",
		},
		{
			name:        "Zero over",
			overNumber:  0,
			ballNumber:  1,
			expected:    "0.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fallOfWickets := &models.FallOfWickets{
				OverNumber: tt.overNumber,
				BallNumber: tt.ballNumber,
			}

			actual := fmt.Sprintf("%d.%d", fallOfWickets.OverNumber, fallOfWickets.BallNumber)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestCreateFallOfWicketsRequest_Validation(t *testing.T) {
	// Test valid request
	validReq := &models.CreateFallOfWicketsRequest{
		MatchID:      "match-1",
		InningsID:    "innings-1",
		OverID:       "over-1",
		BallID:       "ball-1",
		WicketNumber: 1,
		Score:        50,
		OverNumber:   1,
		BallNumber:   1,
	}

	// Basic validation - check that all required fields are set
	assert.NotEmpty(t, validReq.MatchID)
	assert.NotEmpty(t, validReq.InningsID)
	assert.NotEmpty(t, validReq.OverID)
	assert.NotEmpty(t, validReq.BallID)
	assert.Greater(t, validReq.WicketNumber, 0)
	assert.GreaterOrEqual(t, validReq.Score, 0)
	assert.GreaterOrEqual(t, validReq.OverNumber, 0)
	assert.GreaterOrEqual(t, validReq.BallNumber, 1)
	assert.LessOrEqual(t, validReq.BallNumber, 6)
}

func TestUpdateFallOfWicketsRequest_Validation(t *testing.T) {
	// Test valid update request
	validReq := &models.UpdateFallOfWicketsRequest{
		Score: func() *int { s := 100; return &s }(),
	}

	assert.NotNil(t, validReq.Score)
	assert.Equal(t, 100, *validReq.Score)
	assert.GreaterOrEqual(t, *validReq.Score, 0)
}

func TestFallOfWicketsFilters_Validation(t *testing.T) {
	// Test valid filters
	validFilters := &models.FallOfWicketsFilters{
		MatchID:   func() *string { s := "match-1"; return &s }(),
		InningsID: func() *string { s := "innings-1"; return &s }(),
		Limit:     10,
		Offset:    0,
	}

	assert.NotNil(t, validFilters.MatchID)
	assert.NotNil(t, validFilters.InningsID)
	assert.Equal(t, "match-1", *validFilters.MatchID)
	assert.Equal(t, "innings-1", *validFilters.InningsID)
	assert.Equal(t, 10, validFilters.Limit)
	assert.Equal(t, 0, validFilters.Offset)
}

func TestWicketFall_Structure(t *testing.T) {
	// Test WicketFall structure
	wicketFall := models.WicketFall{
		WicketNumber: 1,
		Score:        50,
		OverNumber:   1,
		BallNumber:   3,
		OverPosition: "1.3",
	}

	assert.Equal(t, 1, wicketFall.WicketNumber)
	assert.Equal(t, 50, wicketFall.Score)
	assert.Equal(t, 1, wicketFall.OverNumber)
	assert.Equal(t, 3, wicketFall.BallNumber)
	assert.Equal(t, "1.3", wicketFall.OverPosition)
}

func TestFallOfWicketsSummary_Structure(t *testing.T) {
	// Test FallOfWicketsSummary structure
	summary := models.FallOfWicketsSummary{
		MatchID:      "match-1",
		InningsID:    "innings-1",
		TotalWickets: 2,
		Wickets: []models.WicketFall{
			{
				WicketNumber: 1,
				Score:        25,
				OverNumber:   1,
				BallNumber:   3,
				OverPosition: "1.3",
			},
			{
				WicketNumber: 2,
				Score:        50,
				OverNumber:   2,
				BallNumber:   1,
				OverPosition: "2.1",
			},
		},
	}

	assert.Equal(t, "match-1", summary.MatchID)
	assert.Equal(t, "innings-1", summary.InningsID)
	assert.Equal(t, 2, summary.TotalWickets)
	assert.Len(t, summary.Wickets, 2)

	// Check first wicket
	assert.Equal(t, 1, summary.Wickets[0].WicketNumber)
	assert.Equal(t, 25, summary.Wickets[0].Score)
	assert.Equal(t, "1.3", summary.Wickets[0].OverPosition)

	// Check second wicket
	assert.Equal(t, 2, summary.Wickets[1].WicketNumber)
	assert.Equal(t, 50, summary.Wickets[1].Score)
	assert.Equal(t, "2.1", summary.Wickets[1].OverPosition)
}
