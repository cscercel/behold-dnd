-- name: GetInventoryItem :one
SELECT * FROM inventory_items
WHERE id = $1;

-- name: ListInventoryItems :many
SELECT * FROM inventory_items
WHERE character_id = $1
ORDER BY name;

-- name: CreateInventoryItem :one
INSERT INTO inventory_items (
    character_id,
    name,
    quantity,
    weight,
    value,
    description,
    is_equipped,
    requires_attunement,
    is_attuned
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: UpdateInventoryItem :one
UPDATE inventory_items
SET
    name                = COALESCE(sqlc.narg('name'), name),
    quantity            = COALESCE(sqlc.narg('quantity'), quantity),
    weight              = COALESCE(sqlc.narg('weight'), weight),
    value               = COALESCE(sqlc.narg('value'), value),
    description         = COALESCE(sqlc.narg('description'), description),
    is_equipped         = COALESCE(sqlc.narg('is_equipped'), is_equipped),
    requires_attunement = COALESCE(sqlc.narg('requires_attunement'), requires_attunement),
    is_attuned          = COALESCE(sqlc.narg('is_attuned'), is_attuned)
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeleteInventoryItem :exec
DELETE FROM inventory_items
WHERE id = $1;

-- name: CountAttunedItems :one
SELECT COUNT(*) FROM inventory_items
WHERE character_id = $1
AND is_attuned = TRUE;
