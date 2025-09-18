package graphql

import (
	"spark-park-cricket-backend/internal/models"

	"github.com/graphql-go/graphql"
)

// Define GraphQL types for cricket scorecard
var (
	// BallType enum
	ballTypeEnum = graphql.NewEnum(graphql.EnumConfig{
		Name: "BallType",
		Values: graphql.EnumValueConfigMap{
			"GOOD": &graphql.EnumValueConfig{
				Value: models.BallTypeGood,
			},
			"WIDE": &graphql.EnumValueConfig{
				Value: models.BallTypeWide,
			},
			"NO_BALL": &graphql.EnumValueConfig{
				Value: models.BallTypeNoBall,
			},
			"DEAD_BALL": &graphql.EnumValueConfig{
				Value: models.BallTypeDeadBall,
			},
		},
	})

	// RunType enum
	runTypeEnum = graphql.NewEnum(graphql.EnumConfig{
		Name: "RunType",
		Values: graphql.EnumValueConfigMap{
			"ZERO": &graphql.EnumValueConfig{
				Value: models.RunTypeZero,
			},
			"ONE": &graphql.EnumValueConfig{
				Value: models.RunTypeOne,
			},
			"TWO": &graphql.EnumValueConfig{
				Value: models.RunTypeTwo,
			},
			"THREE": &graphql.EnumValueConfig{
				Value: models.RunTypeThree,
			},
			"FOUR": &graphql.EnumValueConfig{
				Value: models.RunTypeFour,
			},
			"FIVE": &graphql.EnumValueConfig{
				Value: models.RunTypeFive,
			},
			"SIX": &graphql.EnumValueConfig{
				Value: models.RunTypeSix,
			},
			"SEVEN": &graphql.EnumValueConfig{
				Value: models.RunTypeSeven,
			},
			"EIGHT": &graphql.EnumValueConfig{
				Value: models.RunTypeEight,
			},
			"NINE": &graphql.EnumValueConfig{
				Value: models.RunTypeNine,
			},
			"NO_BALL": &graphql.EnumValueConfig{
				Value: models.RunTypeNB,
			},
			"WIDE": &graphql.EnumValueConfig{
				Value: models.RunTypeWD,
			},
			"LEG_BYES": &graphql.EnumValueConfig{
				Value: models.RunTypeLB,
			},
			"WICKET": &graphql.EnumValueConfig{
				Value: models.RunTypeWC,
			},
		},
	})

	// TeamType enum
	teamTypeEnum = graphql.NewEnum(graphql.EnumConfig{
		Name: "TeamType",
		Values: graphql.EnumValueConfigMap{
			"A": &graphql.EnumValueConfig{
				Value: models.TeamTypeA,
			},
			"B": &graphql.EnumValueConfig{
				Value: models.TeamTypeB,
			},
		},
	})

	// BallSummary type
	ballSummaryType = graphql.NewObject(graphql.ObjectConfig{
		Name: "BallSummary",
		Fields: graphql.Fields{
			"ball_number": &graphql.Field{
				Type: graphql.Int,
			},
			"ball_type": &graphql.Field{
				Type: ballTypeEnum,
			},
			"run_type": &graphql.Field{
				Type: runTypeEnum,
			},
			"runs": &graphql.Field{
				Type: graphql.Int,
			},
			"byes": &graphql.Field{
				Type: graphql.Int,
			},
			"is_wicket": &graphql.Field{
				Type: graphql.Boolean,
			},
			"wicket_type": &graphql.Field{
				Type: graphql.String,
			},
		},
	})

	// OverSummary type
	overSummaryType = graphql.NewObject(graphql.ObjectConfig{
		Name: "OverSummary",
		Fields: graphql.Fields{
			"over_number": &graphql.Field{
				Type: graphql.Int,
			},
			"total_runs": &graphql.Field{
				Type: graphql.Int,
			},
			"total_balls": &graphql.Field{
				Type: graphql.Int,
			},
			"total_wickets": &graphql.Field{
				Type: graphql.Int,
			},
			"status": &graphql.Field{
				Type: graphql.String,
			},
			"balls": &graphql.Field{
				Type: graphql.NewList(ballSummaryType),
			},
		},
	})

	// ExtrasSummary type
	extrasSummaryType = graphql.NewObject(graphql.ObjectConfig{
		Name: "ExtrasSummary",
		Fields: graphql.Fields{
			"byes": &graphql.Field{
				Type: graphql.Int,
			},
			"leg_byes": &graphql.Field{
				Type: graphql.Int,
			},
			"wides": &graphql.Field{
				Type: graphql.Int,
			},
			"no_balls": &graphql.Field{
				Type: graphql.Int,
			},
			"total": &graphql.Field{
				Type: graphql.Int,
			},
		},
	})

	// InningsSummary type
	inningsSummaryType = graphql.NewObject(graphql.ObjectConfig{
		Name: "InningsSummary",
		Fields: graphql.Fields{
			"innings_number": &graphql.Field{
				Type: graphql.Int,
			},
			"batting_team": &graphql.Field{
				Type: teamTypeEnum,
			},
			"total_runs": &graphql.Field{
				Type: graphql.Int,
			},
			"total_wickets": &graphql.Field{
				Type: graphql.Int,
			},
			"total_overs": &graphql.Field{
				Type: graphql.Float,
			},
			"total_balls": &graphql.Field{
				Type: graphql.Int,
			},
			"status": &graphql.Field{
				Type: graphql.String,
			},
			"extras": &graphql.Field{
				Type: extrasSummaryType,
			},
			"overs": &graphql.Field{
				Type: graphql.NewList(overSummaryType),
			},
		},
	})

	// LiveScorecard type - optimized for live updates
	liveScorecardType = graphql.NewObject(graphql.ObjectConfig{
		Name: "LiveScorecard",
		Fields: graphql.Fields{
			"match_id": &graphql.Field{
				Type: graphql.String,
			},
			"match_number": &graphql.Field{
				Type: graphql.Int,
			},
			"series_name": &graphql.Field{
				Type: graphql.String,
			},
			"team_a": &graphql.Field{
				Type: graphql.String,
			},
			"team_b": &graphql.Field{
				Type: graphql.String,
			},
			"total_overs": &graphql.Field{
				Type: graphql.Int,
			},
			"toss_winner": &graphql.Field{
				Type: teamTypeEnum,
			},
			"toss_type": &graphql.Field{
				Type: graphql.String,
			},
			"current_innings": &graphql.Field{
				Type: graphql.Int,
			},
			"match_status": &graphql.Field{
				Type: graphql.String,
			},
			"current_score": &graphql.Field{
				Type: graphql.NewObject(graphql.ObjectConfig{
					Name: "CurrentScore",
					Fields: graphql.Fields{
						"runs": &graphql.Field{
							Type: graphql.Int,
						},
						"wickets": &graphql.Field{
							Type: graphql.Int,
						},
						"overs": &graphql.Field{
							Type: graphql.Float,
						},
						"balls": &graphql.Field{
							Type: graphql.Int,
						},
						"run_rate": &graphql.Field{
							Type: graphql.Float,
						},
					},
				}),
			},
			"current_over": &graphql.Field{
				Type: overSummaryType,
			},
			"innings": &graphql.Field{
				Type: graphql.NewList(inningsSummaryType),
			},
		},
	})

	// MatchDetails type for comprehensive match information
	matchDetailsType = graphql.NewObject(graphql.ObjectConfig{
		Name: "MatchDetails",
		Fields: graphql.Fields{
			"match_id": &graphql.Field{
				Type: graphql.String,
			},
			"match_number": &graphql.Field{
				Type: graphql.Int,
			},
			"series_name": &graphql.Field{
				Type: graphql.String,
			},
			"team_a": &graphql.Field{
				Type: graphql.String,
			},
			"team_b": &graphql.Field{
				Type: graphql.String,
			},
			"total_overs": &graphql.Field{
				Type: graphql.Int,
			},
			"toss_winner": &graphql.Field{
				Type: teamTypeEnum,
			},
			"toss_type": &graphql.Field{
				Type: graphql.String,
			},
			"current_innings": &graphql.Field{
				Type: graphql.Int,
			},
			"match_status": &graphql.Field{
				Type: graphql.String,
			},
			"batting_team": &graphql.Field{
				Type: teamTypeEnum,
			},
			"team_a_player_count": &graphql.Field{
				Type: graphql.Int,
			},
			"team_b_player_count": &graphql.Field{
				Type: graphql.Int,
			},
		},
	})

	// MatchStatistics type for match-level statistics
	matchStatisticsType = graphql.NewObject(graphql.ObjectConfig{
		Name: "MatchStatistics",
		Fields: graphql.Fields{
			"total_runs": &graphql.Field{
				Type: graphql.Int,
			},
			"total_wickets": &graphql.Field{
				Type: graphql.Int,
			},
			"total_overs": &graphql.Field{
				Type: graphql.Float,
			},
			"total_balls": &graphql.Field{
				Type: graphql.Int,
			},
			"run_rate": &graphql.Field{
				Type: graphql.Float,
			},
			"extras": &graphql.Field{
				Type: extrasSummaryType,
			},
			"innings_count": &graphql.Field{
				Type: graphql.Int,
			},
		},
	})

	// Team type for team information
	teamType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Team",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"players_count": &graphql.Field{
				Type: graphql.Int,
			},
			"created_at": &graphql.Field{
				Type: graphql.String,
			},
			"updated_at": &graphql.Field{
				Type: graphql.String,
			},
		},
	})

	// Player type for player information
	playerType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Player",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"team_id": &graphql.Field{
				Type: graphql.String,
			},
			"created_at": &graphql.Field{
				Type: graphql.String,
			},
			"updated_at": &graphql.Field{
				Type: graphql.String,
			},
		},
	})

	// PlayerStatistics type for player performance statistics
	playerStatisticsType = graphql.NewObject(graphql.ObjectConfig{
		Name: "PlayerStatistics",
		Fields: graphql.Fields{
			"player_id": &graphql.Field{
				Type: graphql.String,
			},
			"player_name": &graphql.Field{
				Type: graphql.String,
			},
			"team_id": &graphql.Field{
				Type: graphql.String,
			},
			"runs_scored": &graphql.Field{
				Type: graphql.Int,
			},
			"balls_faced": &graphql.Field{
				Type: graphql.Int,
			},
			"wickets_taken": &graphql.Field{
				Type: graphql.Int,
			},
			"overs_bowled": &graphql.Field{
				Type: graphql.Float,
			},
			"runs_conceded": &graphql.Field{
				Type: graphql.Int,
			},
			"strike_rate": &graphql.Field{
				Type: graphql.Float,
			},
			"economy_rate": &graphql.Field{
				Type: graphql.Float,
			},
		},
	})

	// Query type - removed as it's not used (schema is created dynamically in handler.go)
)

// Note: Schema is now created dynamically in the handler with resolver context
// The static schema above is kept for reference but not exported
