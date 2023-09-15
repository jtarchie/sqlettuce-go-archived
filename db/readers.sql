-- name: Get :one
SELECT value
FROM keys
WHERE name = @name;

-- name: Substr :one
SELECT SUBSTR(
  value,
  IIF(@start < 0,
    @start,
    @start + 1
  ),
  IIF(@end < 0,
    LENGTH(value) - @end,
    @start + @end + 1
  )
) FROM keys WHERE name = @name;