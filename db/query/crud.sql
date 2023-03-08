-- name: CreateReport :exec
INSERT INTO "report" (
    "product_id", "description"
) VALUES (
    $1, $2
);

-- name: HandleReport :one
UPDATE "report"
SET "status" = 'handled'
WHERE "id" = $1
RETURNING *;

-- name: DeleteReport :exec
DELETE FROM "report"
WHERE "id" = $1;

-- name: GetAllReport :many
SELECT * FROM "report" WHERE "status" = 'waiting';