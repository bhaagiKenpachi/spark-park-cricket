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

	// Query type
	queryType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"liveScorecard": &graphql.Field{
				Type: liveScorecardType,
				Args: graphql.FieldConfigArgument{
					"match_id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: resolveLiveScorecard,
			},
		},
	})

	// Subscription type for real-time updates
	subscriptionType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Subscription",
		Fields: graphql.Fields{
			"scorecardUpdated": &graphql.Field{
				Type: liveScorecardType,
				Args: graphql.FieldConfigArgument{
					"match_id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: resolveScorecardSubscription,
			},
		},
	})
)

// Schema represents the GraphQL schema
var Schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:        queryType,
	Subscription: subscriptionType,
})
