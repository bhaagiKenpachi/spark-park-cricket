package models

import (
	"time"
)

// Series represents a cricket tournament or competition
type Series struct {
	ID        string    `json:"id,omitempty" db:"id,omitempty"`
	Name      string    `json:"name" db:"name"`
	StartDate time.Time `json:"start_date" db:"start_date"`
	EndDate   time.Time `json:"end_date" db:"end_date"`
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" db:"updated_at,omitempty"`
}

// CreateSeriesRequest represents the request to create a new series
type CreateSeriesRequest struct {
	Name      string    `json:"name" validate:"required,min=3,max=255"`
	StartDate time.Time `json:"start_date" validate:"required"`
	EndDate   time.Time `json:"end_date" validate:"required,gtfield=StartDate"`
}

// UpdateSeriesRequest represents the request to update a series
type UpdateSeriesRequest struct {
	Name      *string    `json:"name,omitempty" validate:"omitempty,min=3,max=255"`
	StartDate *time.Time `json:"start_date,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`
}

// SeriesFilters represents filters for listing series
type SeriesFilters struct {
	Limit  int `json:"limit" validate:"min=1,max=100"`
	Offset int `json:"offset" validate:"min=0"`
}
