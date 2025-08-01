// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: query.sql

package db

import (
	"context"
)

const saveFeed = `-- name: SaveFeed :exec
INSERT INTO
    feed (title, url, last_guid)
VALUES
    (?, ?, ?)
`

type SaveFeedParams struct {
	Title    string
	Url      string
	LastGuid string
}

func (q *Queries) SaveFeed(ctx context.Context, arg SaveFeedParams) error {
	_, err := q.db.ExecContext(ctx, saveFeed, arg.Title, arg.Url, arg.LastGuid)
	return err
}
