-- name: SaveFeed :exec
INSERT INTO
    feed (title, url)
VALUES
    (?, ?);

-- name: SaveFeedItem :exec
INSERT INTO
    feed_item (feed_title, guid, pub_date, processed_at)
VALUES
    (?, ?, ?, ?);
