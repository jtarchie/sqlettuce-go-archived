-- name: Get :one
SELECT value
FROM strings
WHERE name = @name;