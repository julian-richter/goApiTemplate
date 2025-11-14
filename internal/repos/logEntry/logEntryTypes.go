package logentry

import (
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	dbpkg "github.com/julian-richter/ApiTemplate/internal/db"
)

// Sentinel not-found error used by handlers.
var ErrNotFound = errors.New("log entry not found")

// RepoOption applies optional settings to Repo.
type RepoOption func(*Repo)

// Repo persists and optionally caches log entries.
type Repo struct {
	pgPool      *pgxpool.Pool
	cacheClient dbpkg.ValkeyClientInterface
	cachePrefix string
}

// SearchParams holds optional filters for searching log entries.
type SearchParams struct {
	Level           string     // exact log level, empty means ignore
	MessageContains string     // substring to search in message, empty means ignore
	Since           *time.Time // if non-nil, only entries after this time
	Until           *time.Time // if non-nil, only entries before this time
	Limit           int        // max results to return (0 means default)
	Offset          int        // number of results to skip
}

// Query constants for log entry operations.
const (
	InsertOrUpdateQuery = `
        INSERT INTO %s (id, level, message, timestamp)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (id) DO UPDATE
          SET level   = EXCLUDED.level,
              message = EXCLUDED.message,
              timestamp = EXCLUDED.timestamp
    `
	SelectByIDQuery = `
        SELECT id, level, message, timestamp
        FROM %s
        WHERE id = $1
    `
	SelectAllQuery = `
        SELECT id, level, message, timestamp
        FROM %s
    `
)
