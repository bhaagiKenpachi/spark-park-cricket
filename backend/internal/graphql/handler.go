package graphql

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"spark-park-cricket-backend/internal/interfaces"
	"spark-park-cricket-backend/pkg/websocket"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

const (
	resolverContextKey contextKey = "resolver_context"
)

// GraphQLHandler handles GraphQL requests
type GraphQLHandler struct {
	schema *graphql.Schema
	hub    *websocket.Hub
}

// GraphQLRequest represents a GraphQL request
type GraphQLRequest struct {
	Query         string                 `json:"query"`
	Variables     map[string]interface{} `json:"variables,omitempty"`
	OperationName string                 `json:"operationName,omitempty"`
}

// GraphQLResponse represents a GraphQL response
type GraphQLResponse struct {
	Data   interface{}            `json:"data,omitempty"`
	Errors []GraphQLError         `json:"errors,omitempty"`
	Extras map[string]interface{} `json:"extras,omitempty"`
}

// GraphQLError represents a GraphQL error
type GraphQLError struct {
	Message    string                 `json:"message"`
	Locations  []GraphQLErrorLocation `json:"locations,omitempty"`
	Path       []interface{}          `json:"path,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
}

// GraphQLErrorLocation represents the location of a GraphQL error
type GraphQLErrorLocation struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

// NewGraphQLHandler creates a new GraphQL handler
func NewGraphQLHandler(scorecardService interfaces.ScorecardServiceInterface, hub *websocket.Hub) *GraphQLHandler {
	log.Printf("DEBUG: Creating new GraphQL handler")

	// Create resolver context
	resolverCtx := &ResolverContext{
		ScorecardService: scorecardService,
		Hub:              hub,
	}
	log.Printf("DEBUG: Resolver context created")

	// Create schema with resolver context
	schema, err := createSchemaWithContext(resolverCtx)
	if err != nil {
		log.Fatalf("Failed to create GraphQL schema: %v", err)
	}
	log.Printf("DEBUG: GraphQL schema created successfully")

	return &GraphQLHandler{
		schema: schema,
		hub:    hub,
	}
}

// createSchemaWithContext creates a GraphQL schema with resolver context
func createSchemaWithContext(resolverCtx *ResolverContext) (*graphql.Schema, error) {
	log.Printf("DEBUG: Creating GraphQL schema with context")

	// Create query type with resolver context
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"liveScorecard": &graphql.Field{
				Type: liveScorecardType,
				Args: graphql.FieldConfigArgument{
					"match_id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					// Add resolver context to the context
					ctx := context.WithValue(p.Context, resolverContextKey, resolverCtx)
					p.Context = ctx
					return resolveLiveScorecard(p)
				},
			},
			"matchDetails": &graphql.Field{
				Type: matchDetailsType,
				Args: graphql.FieldConfigArgument{
					"match_id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					log.Printf("DEBUG: matchDetails resolver called")
					// Add resolver context to the context
					ctx := context.WithValue(p.Context, resolverContextKey, resolverCtx)
					p.Context = ctx
					return resolveMatchDetails(p)
				},
			},
			"matchStatistics": &graphql.Field{
				Type: matchStatisticsType,
				Args: graphql.FieldConfigArgument{
					"match_id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					// Add resolver context to the context
					ctx := context.WithValue(p.Context, resolverContextKey, resolverCtx)
					p.Context = ctx
					return resolveMatchStatistics(p)
				},
			},
			"inningsScore": &graphql.Field{
				Type: graphql.NewObject(graphql.ObjectConfig{
					Name: "InningsScore",
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
					},
				}),
				Args: graphql.FieldConfigArgument{
					"match_id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"innings_number": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					// Add resolver context to the context
					ctx := context.WithValue(p.Context, resolverContextKey, resolverCtx)
					p.Context = ctx
					return resolveInningsScore(p)
				},
			},
			"inningsDetails": &graphql.Field{
				Type: inningsSummaryType,
				Args: graphql.FieldConfigArgument{
					"match_id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"innings_number": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					// Add resolver context to the context
					ctx := context.WithValue(p.Context, resolverContextKey, resolverCtx)
					p.Context = ctx
					return resolveInningsDetails(p)
				},
			},
			"overDetails": &graphql.Field{
				Type: overSummaryType,
				Args: graphql.FieldConfigArgument{
					"match_id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"innings_number": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
					"over_number": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					// Add resolver context to the context
					ctx := context.WithValue(p.Context, resolverContextKey, resolverCtx)
					p.Context = ctx
					return resolveOverDetails(p)
				},
			},
			"latestOver": &graphql.Field{
				Type: overSummaryType,
				Args: graphql.FieldConfigArgument{
					"match_id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"innings_number": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					// Add resolver context to the context
					ctx := context.WithValue(p.Context, resolverContextKey, resolverCtx)
					p.Context = ctx
					return resolveLatestOver(p)
				},
			},
			"allOvers": &graphql.Field{
				Type: graphql.NewList(overSummaryType),
				Args: graphql.FieldConfigArgument{
					"match_id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"innings_number": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					// Add resolver context to the context
					ctx := context.WithValue(p.Context, resolverContextKey, resolverCtx)
					p.Context = ctx
					return resolveAllOvers(p)
				},
			},
			"ballDetails": &graphql.Field{
				Type: graphql.NewList(ballSummaryType),
				Args: graphql.FieldConfigArgument{
					"match_id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"innings_number": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
					"over_number": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					// Add resolver context to the context
					ctx := context.WithValue(p.Context, resolverContextKey, resolverCtx)
					p.Context = ctx
					return resolveBallDetails(p)
				},
			},
			"matchTeams": &graphql.Field{
				Type: graphql.NewList(teamType),
				Args: graphql.FieldConfigArgument{
					"match_id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					// Add resolver context to the context
					ctx := context.WithValue(p.Context, resolverContextKey, resolverCtx)
					p.Context = ctx
					return resolveMatchTeams(p)
				},
			},
			"matchPlayers": &graphql.Field{
				Type: graphql.NewList(playerType),
				Args: graphql.FieldConfigArgument{
					"match_id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					// Add resolver context to the context
					ctx := context.WithValue(p.Context, resolverContextKey, resolverCtx)
					p.Context = ctx
					return resolveMatchPlayers(p)
				},
			},
			"playerStatistics": &graphql.Field{
				Type: graphql.NewList(playerStatisticsType),
				Args: graphql.FieldConfigArgument{
					"match_id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					// Add resolver context to the context
					ctx := context.WithValue(p.Context, resolverContextKey, resolverCtx)
					p.Context = ctx
					return resolvePlayerStatistics(p)
				},
			},
		},
	})

	// Create subscription type with resolver context
	subscriptionType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Subscription",
		Fields: graphql.Fields{
			"scorecardUpdated": &graphql.Field{
				Type: liveScorecardType,
				Args: graphql.FieldConfigArgument{
					"match_id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					// Add resolver context to the context
					ctx := context.WithValue(p.Context, resolverContextKey, resolverCtx)
					p.Context = ctx
					return resolveScorecardSubscription(p)
				},
			},
		},
	})

	log.Printf("DEBUG: Creating GraphQL schema with %d query fields", len(queryType.Fields()))

	// Log all query field names
	for fieldName := range queryType.Fields() {
		log.Printf("DEBUG: Query field: %s", fieldName)
	}

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:        queryType,
		Subscription: subscriptionType,
	})
	if err != nil {
		log.Printf("DEBUG: Error creating GraphQL schema: %v", err)
		return nil, err
	}

	log.Printf("DEBUG: GraphQL schema created successfully")
	return &schema, nil
}

// ServeHTTP handles HTTP requests for GraphQL
func (h *GraphQLHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Only allow POST requests for GraphQL
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var req GraphQLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding GraphQL request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("DEBUG: GraphQL request received - Query: %s", req.Query)
	log.Printf("DEBUG: GraphQL request variables: %+v", req.Variables)

	// Execute GraphQL query
	result := graphql.Do(graphql.Params{
		Schema:         *h.schema,
		RequestString:  req.Query,
		VariableValues: req.Variables,
		OperationName:  req.OperationName,
		Context:        r.Context(),
	})

	// Handle errors
	if len(result.Errors) > 0 {
		log.Printf("GraphQL errors: %v", result.Errors)
	}

	// Create response
	var graphqlErrors []GraphQLError
	for _, err := range result.Errors {
		graphqlErrors = append(graphqlErrors, GraphQLError{
			Message: err.Message,
		})
	}

	response := GraphQLResponse{
		Data:   result.Data,
		Errors: graphqlErrors,
	}

	// Set content type
	w.Header().Set("Content-Type", "application/json")

	// Write response
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding GraphQL response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// GetPlaygroundHandler returns a GraphQL playground handler for development
func (h *GraphQLHandler) GetPlaygroundHandler() http.Handler {
	playgroundHandler := handler.New(&handler.Config{
		Schema:   h.schema,
		Pretty:   true,
		GraphiQL: true,
	})

	return playgroundHandler
}

// BroadcastScorecardUpdate broadcasts a scorecard update to WebSocket clients
func (h *GraphQLHandler) BroadcastScorecardUpdate(matchID string, scorecard interface{}) {
	if h.hub != nil {
		// Create update message
		updateMessage := map[string]interface{}{
			"type": "scorecard_update",
			"data": scorecard,
		}

		// Broadcast to the match room
		h.hub.BroadcastToRoom(matchID, updateMessage)
	}
}
