CREATE TABLE IF NOT EXISTS log_entries (
                                           id INTEGER PRIMARY KEY,
                                           level TEXT NOT NULL,
                                           message TEXT NOT NULL,
                                           timestamp TIMESTAMPTZ NOT NULL
);
