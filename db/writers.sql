-- name: Set :exec
INSERT INTO strings (name, value)
VALUES (@name, @value) ON CONFLICT(name) DO
UPDATE
SET value = excluded.value;