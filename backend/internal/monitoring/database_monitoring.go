package monitoring

import (
	"context"
	"time"
)

// DatabaseOperation represents a database operation with monitoring
type DatabaseOperation struct {
	Operation string
	Table     string
	MatchID   string
	StartTime time.Time
	Metrics   *Metrics
}

// NewDatabaseOperation creates a new database operation monitor
func NewDatabaseOperation(metrics *Metrics, operation, table, matchID string) *DatabaseOperation {
	return &DatabaseOperation{
		Operation: operation,
		Table:     table,
		MatchID:   matchID,
		StartTime: time.Now(),
		Metrics:   metrics,
	}
}

// Finish records the operation as successful
func (d *DatabaseOperation) Finish() {
	duration := time.Since(d.StartTime)
	d.Metrics.RecordBallAPIDatabaseOperation(d.Operation, d.Table, d.MatchID, "success", duration)
}

// FinishWithError records the operation as failed
func (d *DatabaseOperation) FinishWithError(errorType string) {
	duration := time.Since(d.StartTime)
	d.Metrics.RecordBallAPIDatabaseError(d.Operation, d.Table, d.MatchID, errorType)
	d.Metrics.RecordBallAPIDatabaseOperation(d.Operation, d.Table, d.MatchID, "error", duration)
}

// WithDatabaseMonitoring wraps a database operation with monitoring
func WithDatabaseMonitoring(
	metrics *Metrics,
	operation, table, matchID string,
	dbFunc func() error,
) error {
	dbOp := NewDatabaseOperation(metrics, operation, table, matchID)

	err := dbFunc()
	if err != nil {
		dbOp.FinishWithError("database_error")
		return err
	}

	dbOp.Finish()
	return nil
}

// WithDatabaseMonitoringContext wraps a database operation with monitoring and context
func WithDatabaseMonitoringContext(
	ctx context.Context,
	metrics *Metrics,
	operation, table, matchID string,
	dbFunc func(ctx context.Context) error,
) error {
	dbOp := NewDatabaseOperation(metrics, operation, table, matchID)

	err := dbFunc(ctx)
	if err != nil {
		dbOp.FinishWithError("database_error")
		return err
	}

	dbOp.Finish()
	return nil
}
