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

-- name: UpdateCharacter :one
UPDATE characters
SET
    name                    = $2,
    race                    = $3,
    class                   = $4,
    level                   = $5,
    background              = $6,
    alignment              = $7,
    xp                      = $8,
    strength                = $9,
    dexterity               = $10,
    constitution            = $11,
    intelligence            = $12,
    wisdom                  = $13,
    charisma                = $14,
    save_prof_strength      = $15,
    save_prof_dexterity     = $16,
    save_prof_constitution  = $17,
    save_prof_intelligence  = $18,
    save_prof_wisdom        = $19,
    save_prof_charisma      = $20,
    armor_class             = $21,
    speed                   = $22,
    hit_dice_type           = $23,
    hit_dice_remaining      = $24,
    inspiration             = $25,
    updated_at              = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateDeathSaves :one
UPDATE characters
SET
    death_save_successes    = $2,
    death_save_failures     = $3,
    updated_at              = NOW()
WHERE id = $1
RETURNING *;

-- name: ResetDeathSaves :one
UPDATE characters
SET
    death_save_successes    = 0,
    death_save_failures     = 0,
    updated_at              = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateConditions :one
UPDATE characters
SET
    conditions = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: LongRest :one
UPDATE characters
SET
    current_hp              = max_hp,
    hit_dice_remaining      = GREATEST(hit_dice_remaining + (level / 2), level),
    death_save_successes    = 0,
    death_save_failures     = 0,
    conditions              = '{}',
    updated_at              = NOW()
WHERE id = $1
RETURNING *;

-- name: ShortRest :one
UPDATE characters
SET
    hit_dice_remaining      = GREATEST(hit_dice_remaining - $2, 0),
    current_hp              = LEAST(current_hp + $3, max_hp),
    updated_at              = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteCharacter :exec
DELETE FROM characters
WHERE id = $1;
