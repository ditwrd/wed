package repository

import "github.com/ditwrd/wed/internal/model"

// RSVPRepository defines the interface for RSVP data access
type RSVP interface {
	// Create inserts a new RSVP entry
	Create(rsvp *model.RSVP) error

	// GetByID retrieves an RSVP entry by its ID
	GetByID(id string) (*model.RSVP, error)

	// GetAll retrieves all RSVP entries
	GetAll() ([]model.RSVP, error)

	// GetPaginated retrieves RSVP entries with pagination
	GetPaginated(limit, offset int) ([]model.RSVP, error)

	// GetCount returns the total number of RSVP entries
	GetCount() (int, error)

	// GetLatestMessages retrieves the latest non-empty messages (up to limit)
	GetLatestMessages(limit int) ([]model.RSVP, error)

	// GetStats returns statistics about RSVP entries
	GetStats() (map[string]interface{}, error)
}
