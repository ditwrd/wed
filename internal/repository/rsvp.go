package repository

import (
	"context"

	"github.com/ditwrd/wed/internal/model"
)

// RSVPRepository defines the interface for RSVP data access
type RSVP interface {
	Create(ctx context.Context, rsvp *model.RSVP) error
	GetByID(ctx context.Context, id string) (*model.RSVP, error)
	GetAll(ctx context.Context) ([]model.RSVP, error)
	GetPaginated(ctx context.Context, limit, offset int) ([]model.RSVP, error)
	GetCount(ctx context.Context) (int, error)
	GetLatestMessages(ctx context.Context, limit int) ([]model.RSVP, error)
	GetStats(ctx context.Context) (map[string]interface{}, error)
}
