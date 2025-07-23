CREATE TABLE IF NOT EXISTS feed (
    title TEXT PRIMARY KEY,
    url TEXT NOT NULL,
    last_guid TEXT NOT NULL
);
