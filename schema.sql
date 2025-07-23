CREATE TABLE IF NOT EXISTS feed (title TEXT PRIMARY KEY, url TEXT NOT NULL);

CREATE TABLE IF NOT EXISTS feed_item (
    feed_title TEXT NOT NULL,
    guid TEXT NOT NULL,
    pub_date INTEGER NOT NULL,
    processed_at INTEGER
);
