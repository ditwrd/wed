package sqlite

import (
	"context"
	"errors"
	"time"

	"github.com/ditwrd/wed/internal/model"
	"github.com/ditwrd/wed/internal/repository"
	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
	"github.com/jmoiron/sqlx"
)

var (
	ErrNotFound = errors.New("rsvp not found")
)

type TxFunc func(ctx context.Context, tx *sqlx.Tx) error

func withTx(ctx context.Context, db *sqlx.DB, fn TxFunc) error {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()
	if err = fn(ctx, tx); err != nil {
		return err
	}
	return tx.Commit()
}

// SQLiteRSVPRepository implements RSVPRepository interface using SQLite
type SQLiteRSVPRepository struct {
	db *sqlx.DB
}

// NewSQLiteRSVPRepository creates a new SQLiteRSVPRepository
func NewSQLiteRSVPRepository(db *sqlx.DB) repository.RSVP {
	return &SQLiteRSVPRepository{
		db: db,
	}
}

// rsvpColumn is a struct with database tags for SQLite operations
type rsvpColumn struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	Attending bool      `db:"attending"`
	Message   string    `db:"message"`
	GroupName string    `db:"group_name"`
	CreatedAt time.Time `db:"created_at"`
}

// Create inserts a new RSVP entry
func (r *SQLiteRSVPRepository) Create(ctx context.Context, rsvp *model.RSVP) error {
	// Generate a new UUID if not provided
	if rsvp.ID == "" {
		rsvp.ID = uuid.New().String()
	}

	// Set created time if not provided
	if rsvp.CreatedAt.IsZero() {
		rsvp.CreatedAt = time.Now()
	}

	// Convert model.RSVP to rsvpDB for database operations
	rsvpDB := rsvpColumn{
		ID:        rsvp.ID,
		Name:      rsvp.Name,
		Attending: rsvp.Attending,
		Message:   rsvp.Message,
		GroupName: rsvp.GroupName,
		CreatedAt: rsvp.CreatedAt,
	}

	// Use go-sqlbuilder for INSERT
	sb := sqlbuilder.SQLite.NewInsertBuilder()
	sb.InsertInto("rsvps")
	sb.Cols("id", "name", "attending", "message", "group_name", "created_at")
	sb.Values(
		rsvpDB.ID,
		rsvpDB.Name,
		rsvpDB.Attending,
		rsvpDB.Message,
		rsvpDB.GroupName,
		rsvpDB.CreatedAt,
	)

	query, args := sb.Build()
	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

// GetByID retrieves an RSVP entry by its ID
func (r *SQLiteRSVPRepository) GetByID(ctx context.Context, id string) (*model.RSVP, error) {
	var rsvp rsvpColumn

	// Use go-sqlbuilder for SELECT
	sb := sqlbuilder.SQLite.NewSelectBuilder()
	sb.Select("id", "name", "attending", "message", "group_name", "created_at")
	sb.From("rsvps")
	sb.Where(sb.Equal("id", id))

	query, args := sb.Build()
	err := r.db.GetContext(ctx, &rsvp, query, args...)
	if err != nil {
		return nil, ErrNotFound
	}

	// Convert rsvpDB to model.RSVP
	modelRSVP := model.RSVP{
		ID:        rsvp.ID,
		Name:      rsvp.Name,
		Attending: rsvp.Attending,
		Message:   rsvp.Message,
		GroupName: rsvp.GroupName,
		CreatedAt: rsvp.CreatedAt,
	}

	return &modelRSVP, nil
}

// GetAll retrieves all RSVP entries
func (r *SQLiteRSVPRepository) GetAll(ctx context.Context) ([]model.RSVP, error) {
	var rsvps []rsvpColumn

	// Use go-sqlbuilder for SELECT
	sb := sqlbuilder.SQLite.NewSelectBuilder()
	sb.Select("id", "name", "attending", "message", "group_name", "created_at")
	sb.From("rsvps")
	sb.OrderBy("created_at DESC")

	query, args := sb.Build()
	err := r.db.SelectContext(ctx, &rsvps, query, args...)
	if err != nil {
		return nil, err
	}

	// Convert []rsvpDB to []model.RSVP
	modelRSVPs := make([]model.RSVP, len(rsvps))
	for i, rsvp := range rsvps {
		modelRSVPs[i] = model.RSVP{
			ID:        rsvp.ID,
			Name:      rsvp.Name,
			Attending: rsvp.Attending,
			Message:   rsvp.Message,
			GroupName: rsvp.GroupName,
			CreatedAt: rsvp.CreatedAt,
		}
	}

	return modelRSVPs, nil
}

// GetPaginated retrieves RSVP entries with pagination
func (r *SQLiteRSVPRepository) GetPaginated(ctx context.Context, limit, offset int) ([]model.RSVP, error) {
	var rsvps []rsvpColumn

	// Use go-sqlbuilder for SELECT with pagination
	sb := sqlbuilder.SQLite.NewSelectBuilder()
	sb.Select("id", "name", "attending", "message", "group_name", "created_at")
	sb.From("rsvps")
	sb.OrderBy("created_at DESC")
	sb.Limit(limit)
	sb.Offset(offset)

	query, args := sb.Build()
	err := r.db.SelectContext(ctx, &rsvps, query, args...)
	if err != nil {
		return nil, err
	}

	// Convert []rsvpDB to []model.RSVP
	modelRSVPs := make([]model.RSVP, len(rsvps))
	for i, rsvp := range rsvps {
		modelRSVPs[i] = model.RSVP{
			ID:        rsvp.ID,
			Name:      rsvp.Name,
			Attending: rsvp.Attending,
			Message:   rsvp.Message,
			GroupName: rsvp.GroupName,
			CreatedAt: rsvp.CreatedAt,
		}
	}

	return modelRSVPs, nil
}

// GetCount returns the total number of RSVP entries
func (r *SQLiteRSVPRepository) GetCount(ctx context.Context) (int, error) {
	var count int

	// Use go-sqlbuilder for COUNT
	sb := sqlbuilder.SQLite.NewSelectBuilder()
	sb.Select("COUNT(*)")
	sb.From("rsvps")

	query, args := sb.Build()
	err := r.db.GetContext(ctx, &count, query, args...)
	return count, err
}

// GetLatestMessages retrieves the latest non-empty messages (up to limit)
func (r *SQLiteRSVPRepository) GetLatestMessages(ctx context.Context, limit int) ([]model.RSVP, error) {
	var messages []rsvpColumn

	sb := sqlbuilder.SQLite.NewSelectBuilder()
	sb.Select("id", "name", "attending", "message", "group_name", "created_at")
	sb.From("rsvps")
	sb.Where("message IS NOT NULL AND message != ''")
	sb.OrderBy("created_at DESC")
	sb.Limit(limit)

	query, args := sb.Build()
	err := r.db.SelectContext(ctx, &messages, query, args...)
	if err != nil {
		return nil, err
	}

	// Convert []rsvpDB to []model.RSVP
	modelRSVPs := make([]model.RSVP, len(messages))
	for i, message := range messages {
		modelRSVPs[i] = model.RSVP{
			ID:        message.ID,
			Name:      message.Name,
			Attending: message.Attending,
			Message:   message.Message,
			GroupName: message.GroupName,
			CreatedAt: message.CreatedAt,
		}
	}

	return modelRSVPs, nil
}

// GetStats returns statistics about RSVP entries
func (r *SQLiteRSVPRepository) GetStats(ctx context.Context) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total RSVP count
	var totalCount int
	sb := sqlbuilder.SQLite.NewSelectBuilder()
	sb.Select("COUNT(*)")
	sb.From("rsvps")
	query, args := sb.Build()
	err := r.db.GetContext(ctx, &totalCount, query, args...)
	if err != nil {
		return nil, err
	}
	stats["total"] = totalCount

	// Attending count
	var attendingCount int
	sb = sqlbuilder.SQLite.NewSelectBuilder()
	sb.Select("COUNT(*)")
	sb.From("rsvps")
	sb.Where("attending = 1")
	query, args = sb.Build()
	err = r.db.GetContext(ctx, &attendingCount, query, args...)
	if err != nil {
		return nil, err
	}
	stats["attending"] = attendingCount

	// Not attending count
	var notAttendingCount int
	sb = sqlbuilder.SQLite.NewSelectBuilder()
	sb.Select("COUNT(*)")
	sb.From("rsvps")
	sb.Where("attending = 0")
	query, args = sb.Build()
	err = r.db.GetContext(ctx, &notAttendingCount, query, args...)
	if err != nil {
		return nil, err
	}
	stats["not_attending"] = notAttendingCount

	// Group count
	var groupCount int
	sb = sqlbuilder.SQLite.NewSelectBuilder()
	sb.Select("COUNT(DISTINCT group_name)")
	sb.From("rsvps")
	sb.Where("group_name IS NOT NULL AND group_name != ''")
	query, args = sb.Build()
	err = r.db.GetContext(ctx, &groupCount, query, args...)
	if err != nil {
		return nil, err
	}
	stats["groups"] = groupCount

	return stats, nil
}
