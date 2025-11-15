package logentry

import (
	"errors"
	"text/template"
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

// Template definitions for SQL queries.
var (
	// Use `define` so you can reuse parts if needed later.
	queryTmpl = template.Must(template.New("logentry_queries").Parse(`
		{{ define "insert" }}
			INSERT INTO {{ .Table }} (level, message, timestamp)
			VALUES ($1, $2, $3)
			RETURNING id
		{{ end }}

        {{ define "selectByID" }}
			SELECT id, level, message, timestamp
			FROM {{ .Table }}
			WHERE id = $1
        {{ end }}

        {{ define "selectAll" }}
			SELECT id, level, message, timestamp
			FROM {{ .Table }}
        {{ end }}

        {{ define "search" }}
			SELECT id, level, message, timestamp
			FROM {{ .Table }}
			WHERE {{ .WhereClause }}
			ORDER BY timestamp DESC
			LIMIT ${{ .LimitPos }} OFFSET ${{ .OffsetPos }}
        {{ end }}
    `))
)
