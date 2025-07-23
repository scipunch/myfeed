-- name: SaveFeed :exec
INSERT INTO
    feed (title, url, last_guid)
VALUES
    (?, ?, ?);
