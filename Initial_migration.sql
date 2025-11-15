CREATE TABLE IF NOT EXISTS log_entries (
                                           id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
                                           level TEXT NOT NULL,
                                           message TEXT NOT NULL,
                                           timestamp TIMESTAMPTZ NOT NULL
);
