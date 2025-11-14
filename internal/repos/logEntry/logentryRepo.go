package logentry

import (
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
	query := strings.TrimSpace(fmt.Sprintf(InsertOrUpdateQuery, r.tableName()))

	_, err := r.pgPool.Exec(ctx, query,
		entry.ID,
		entry.Level,
		entry.Message,
		entry.Timestamp,
	)
	if err != nil {
		return fmt.Errorf("Repo.Save: db exec error: %w", err)
	}

	if r.cacheClient != nil {
		key := r.cacheKey(entry.ID)
		b, err := json.Marshal(entry)
		if err == nil {
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

	query := strings.TrimSpace(fmt.Sprintf(SelectByIDQuery, r.tableName()))
	row := r.pgPool.QueryRow(ctx, query, id)
	err := row.Scan(&entry.ID, &entry.Level, &entry.Message, &entry.Timestamp)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("Repo.GetByID: row scan error: %w", err)
	}

	if useCache && r.cacheClient != nil {
		key := r.cacheKey(id)
		if b, err := json.Marshal(&entry); err == nil {
			_ = r.cacheClient.Set(ctx, key, string(b), ttl)
		}
	}

	return &entry, nil
}

// All returns all log entries, no caching by default.
func (r *Repo) All(ctx context.Context) ([]*modelpkg.LogEntry, error) {
	query := strings.TrimSpace(fmt.Sprintf(SelectAllQuery, r.tableName()))
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

	query := strings.TrimSpace(fmt.Sprintf(`
        SELECT id, level, message, timestamp
        FROM %s
        WHERE %s
        ORDER BY timestamp DESC
        LIMIT $%d OFFSET $%d
    `, r.tableName(), strings.Join(whereClauses, " AND "), argPos, argPos+1))
	args = append(args, limit, offset)

	rows, err := r.pgPool.Query(ctx, query, args...)
	if err != nil {
		return make([]*modelpkg.LogEntry, 0), fmt.Errorf("Repo.Search: query error: %w", err)
	}
	defer rows.Close()

	result := make([]*modelpkg.LogEntry, 0)
	for rows.Next() {
		var e modelpkg.LogEntry
		if err := rows.Scan(&e.ID, &e.Level, &e.Message, &e.Timestamp); err != nil {
			return make([]*modelpkg.LogEntry, 0), fmt.Errorf("Repo.Search: row scan error: %w", err)
		}
		result = append(result, &e)
	}

	return result, nil
}
