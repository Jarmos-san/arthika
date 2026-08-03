-- name: CreateAssetClass :exec
INSERT INTO asset_classes (id, name, description)
VALUES (?, ?, ?);

-- name: ListAssetClasses :many
SELECT id, name, description
FROM asset_classes
ORDER BY name COLLATE NOCASE ASC;

-- name: FindAssetClassByID :one
SELECT id, name, description
FROM asset_classes
WHERE id = ?
LIMIT 1;

-- name: UpdateAssetClass :one
UPDATE asset_classes
SET name = ?, description = ?
WHERE id = ?
RETURNING id, name, description;

-- name: DeleteAssetClass :one
DELETE FROM asset_classes
WHERE id = ?
RETURNING id;
