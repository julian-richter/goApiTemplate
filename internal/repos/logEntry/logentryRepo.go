package logentry

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	dbpkg "github.com/julian-richter/ApiTemplate/internal/db"
	modelpkg "github.com/julian-richter/ApiTemplate/internal/models/logentry"
)

// WithCache enables caching by providing a cache client and key prefix.
func WithCache(client dbpkg.ValkeyClientInterface, prefix string) RepoOption {
	return func(r *Repo) {
		r.cacheClient = client
		r.cachePrefix = prefix
	}
}

// NewRepo creates a new log entry repository. Pass WithCache if caching is desired.
func NewRepo(pgPool *pgxpool.Pool, opts ...RepoOption) *Repo {
	r := &Repo{pgPool: pgPool}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

func (r *Repo) tableName() string {
	return "log_entries"
}

func (r *Repo) cacheKey(id int) string {
	return fmt.Sprintf("%slogentry:%d", r.cachePrefix, id)
}

// Save persists or updates a LogEntry, and caches it if configured.
func (r *Repo) Save(ctx context.Context, entry *modelpkg.LogEntry) error {
	tmplData := struct {
		Table string
	}{
		Table: r.tableName(),
	}

	var query string
	var err error

	if entry.ID <= 0 {
		// New entry - use INSERT with RETURNING id
		var buf bytes.Buffer
		if err = queryTmpl.ExecuteTemplate(&buf, "insert", tmplData); err != nil {
			return fmt.Errorf("Repo.Save: template execution error for insert: %w", err)
		}
		query = buf.String()

		// INSERT (level, message, timestamp) VALUES ($1,$2,$3) RETURNING id
		err = r.pgPool.QueryRow(ctx, query,
			entry.Level,
			entry.Message,
			entry.Timestamp,
		).Scan(&entry.ID)

		if err != nil {
			return fmt.Errorf("Repo.Save: insert failed: %w", err)
		}
	} else {
		// Existing entry - use UPDATE
		var buf bytes.Buffer
		if err = queryTmpl.ExecuteTemplate(&buf, "update", tmplData); err != nil {
			return fmt.Errorf("Repo.Save: template execution error for update: %w", err)
		}
		query = buf.String()

		// UPDATE table SET level=$1, message=$2, timestamp=$3 WHERE id=$4
		tag, err := r.pgPool.Exec(ctx, query,
			entry.Level,
			entry.Message,
			entry.Timestamp,
			entry.ID,
		)

		if err != nil {
			return fmt.Errorf("Repo.Save: update failed: %w", err)
		}

		if tag.RowsAffected() == 0 {
			return fmt.Errorf("Repo.Save: no rows affected, entry with id %d not found", entry.ID)
		}
	}

	// Update cache if configured
	if r.cacheClient != nil {
		key := r.cacheKey(entry.ID)
		if b, err := json.Marshal(entry); err == nil {
			_ = r.cacheClient.Set(ctx, key, string(b), 0)
		}
	}

	return nil
}

// GetByID retrieves a LogEntry by ID, optionally using cache.
func (r *Repo) GetByID(ctx context.Context, id int, useCache bool, ttl time.Duration) (*modelpkg.LogEntry, error) {
	var entry modelpkg.LogEntry

	if useCache && r.cacheClient != nil {
		key := r.cacheKey(id)
		if val, err := r.cacheClient.Get(ctx, key); err == nil {
			if err2 := json.Unmarshal([]byte(val), &entry); err2 == nil {
				return &entry, nil
			}
		}
	}

	tmplData := struct {
		Table string
	}{
		Table: r.tableName(),
	}
	var buf bytes.Buffer
	if err := queryTmpl.ExecuteTemplate(&buf, "selectByID", tmplData); err != nil {
		return nil, fmt.Errorf("Repo.GetByID: template execution error: %w", err)
	}
	query := buf.String()

	row := r.pgPool.QueryRow(ctx, query, id)
	err := row.Scan(&entry.ID, &entry.Level, &entry.Message, &entry.Timestamp)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("Repo.GetByID: row scan error: %w", err)
	}

	if useCache && r.cacheClient != nil {
		key := r.cacheKey(entry.ID)
		if b, err := json.Marshal(&entry); err == nil {
			_ = r.cacheClient.Set(ctx, key, string(b), ttl)
		}
	}

	return &entry, nil
}

// All returns all log entries, no caching by default.
func (r *Repo) All(ctx context.Context) ([]*modelpkg.LogEntry, error) {
	tmplData := struct {
		Table string
	}{
		Table: r.tableName(),
	}
	var buf bytes.Buffer
	if err := queryTmpl.ExecuteTemplate(&buf, "selectAll", tmplData); err != nil {
		return nil, fmt.Errorf("Repo.All: template execution error: %w", err)
	}
	query := buf.String()

	rows, err := r.pgPool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("Repo.All: query error: %w", err)
	}
	defer rows.Close()

	var result []*modelpkg.LogEntry
	for rows.Next() {
		var e modelpkg.LogEntry
		if err := rows.Scan(&e.ID, &e.Level, &e.Message, &e.Timestamp); err != nil {
			return nil, fmt.Errorf("Repo.All: row scan error: %w", err)
		}
		result = append(result, &e)
	}

	// detect mid-stream errors.
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Repo.All: rows error: %w", err)
	}

	return result, nil
}

// Search returns log entries matching filters in SearchParams.
func (r *Repo) Search(ctx context.Context, params SearchParams) ([]*modelpkg.LogEntry, error) {
	const maxLimit = 1000

	whereClauses := []string{"1=1"}
	args := []interface{}{}
	argPos := 1

	if params.Level != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("level = $%d", argPos))
		args = append(args, params.Level)
		argPos++
	}
	if params.MessageContains != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("message ILIKE $%d", argPos))
		args = append(args, "%"+params.MessageContains+"%")
		argPos++
	}
	if params.Since != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("timestamp >= $%d", argPos))
		args = append(args, *params.Since)
		argPos++
	}
	if params.Until != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("timestamp <= $%d", argPos))
		args = append(args, *params.Until)
		argPos++
	}

	limit := params.Limit
	if limit <= 0 {
		limit = 100
	}
	if limit > maxLimit {
		limit = maxLimit
	}

	offset := params.Offset
	if offset < 0 {
		offset = 0
	}

	tmplData := struct {
		Table       string
		WhereClause string
		LimitPos    int
		OffsetPos   int
	}{
		Table:       r.tableName(),
		WhereClause: strings.Join(whereClauses, " AND "),
		LimitPos:    argPos,
		OffsetPos:   argPos + 1,
	}

	var buf bytes.Buffer
	if err := queryTmpl.ExecuteTemplate(&buf, "search", tmplData); err != nil {
		return nil, fmt.Errorf("Repo.Search: template execution error: %w", err)
	}
	query := buf.String()

	args = append(args, limit, offset)

	rows, err := r.pgPool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("Repo.Search: query error: %w", err)
	}
	defer rows.Close()

	var result []*modelpkg.LogEntry
	for rows.Next() {
		var e modelpkg.LogEntry
		if err := rows.Scan(&e.ID, &e.Level, &e.Message, &e.Timestamp); err != nil {
			return nil, fmt.Errorf("Repo.Search: row scan error: %w", err)
		}
		result = append(result, &e)
	}

	// detect mid-stream / final iteration errors
	if err := rows.Err(); err != nil {
		return make([]*modelpkg.LogEntry, 0), fmt.Errorf("Repo.Search: rows iteration error: %w", err)
	}

	return result, nil
}
