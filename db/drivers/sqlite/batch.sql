-- name: Delete :many
DELETE FROM keys WHERE name IN (sqlc.slice('names')) RETURNING value;

-- name: Get :many
SELECT name, value FROM keys WHERE name IN (sqlc.slice('names'));