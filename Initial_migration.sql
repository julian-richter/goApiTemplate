CREATE TABLE IF NOT EXISTS log_entries (
                                           id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
                                           level TEXT NOT NULL,
                                           message TEXT NOT NULL,
                                           timestamp TIMESTAMPTZ NOT NULL
);

-- Indexes to support common queries
CREATE INDEX IF NOT EXISTS idx_log_entries_level ON log_entries(level);
CREATE INDEX IF NOT EXISTS idx_log_entries_timestamp ON log_entries(timestamp);
