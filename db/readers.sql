-- name: Get :one
SELECT value
FROM keys
WHERE name = @name;