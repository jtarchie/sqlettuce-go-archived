-- name: Delete :many
DELETE FROM keys WHERE name IN (sqlc.slice('names')) RETURNING value;