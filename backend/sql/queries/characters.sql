-- name: GetCharacter :one
SELECT * FROM characters
WHERE id = $1;

-- name: ListCharacters :many
SELECT * FROM characters
ORDER BY name;

-- name: ListPlayerCharacters :many
SELECT * FROM characters
WHERE is_npc = FALSE
ORDER BY name;

-- name: ListNPCs :many
SELECT * FROM characters
WHERE is_npc = TRUE
ORDER BY name;

-- name: CreateCharacter :one
INSERT INTO characters (
    owner_id,
    is_npc,
    name,
    race,
    class,
    level,
    max_hp,
    current_hp,
    armor_class,
    speed
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
RETURNING *;

-- name: UpdateCharacterHP :one
UPDATE characters
SET
    current_hp = $2,
    temp_hp    = $3,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteCharacter :exec
DELETE FROM characters
WHERE id = $1;
