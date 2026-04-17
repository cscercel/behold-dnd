-- name: GetCharacter :one
SELECT * FROM characters
WHERE id = $1;

-- name: ListCharacters :many
SELECT * FROM characters
ORDER BY name;

-- name: ListMyCharacters :many
SELECT * FROM characters
WHERE owner_id = $1
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
    name                   = COALESCE(sqlc.narg('name'), name),
    race                   = COALESCE(sqlc.narg('race'), race),
    class                  = COALESCE(sqlc.narg('class'), class),
    level                  = COALESCE(sqlc.narg('level'), level),
    background             = COALESCE(sqlc.narg('background'), background),
    alignment              = COALESCE(sqlc.narg('alignment'), alignment),
    xp                     = COALESCE(sqlc.narg('xp'), xp),
    strength               = COALESCE(sqlc.narg('strength'), strength),
    dexterity              = COALESCE(sqlc.narg('dexterity'), dexterity),
    constitution           = COALESCE(sqlc.narg('constitution'), constitution),
    intelligence           = COALESCE(sqlc.narg('intelligence'), intelligence),
    wisdom                 = COALESCE(sqlc.narg('wisdom'), wisdom),
    charisma               = COALESCE(sqlc.narg('charisma'), charisma),
    save_prof_strength     = COALESCE(sqlc.narg('save_prof_strength'), save_prof_strength),
    save_prof_dexterity    = COALESCE(sqlc.narg('save_prof_dexterity'), save_prof_dexterity),
    save_prof_constitution = COALESCE(sqlc.narg('save_prof_constitution'), save_prof_constitution),
    save_prof_intelligence = COALESCE(sqlc.narg('save_prof_intelligence'), save_prof_intelligence),
    save_prof_wisdom       = COALESCE(sqlc.narg('save_prof_wisdom'), save_prof_wisdom),
    save_prof_charisma     = COALESCE(sqlc.narg('save_prof_charisma'), save_prof_charisma),
    skill_acrobatics       = COALESCE(sqlc.narg('skill_acrobatics'), skill_acrobatics),
    skill_animal_handling  = COALESCE(sqlc.narg('skill_animal_handling'), skill_animal_handling),
    skill_arcana           = COALESCE(sqlc.narg('skill_arcana'), skill_arcana),
    skill_athletics        = COALESCE(sqlc.narg('skill_athletics'), skill_athletics),
    skill_deception        = COALESCE(sqlc.narg('skill_deception'), skill_deception),
    skill_history          = COALESCE(sqlc.narg('skill_history'), skill_history),
    skill_insight          = COALESCE(sqlc.narg('skill_insight'), skill_insight),
    skill_intimidation     = COALESCE(sqlc.narg('skill_intimidation'), skill_intimidation),
    skill_investigation    = COALESCE(sqlc.narg('skill_investigation'), skill_investigation),
    skill_medicine         = COALESCE(sqlc.narg('skill_medicine'), skill_medicine),
    skill_nature           = COALESCE(sqlc.narg('skill_nature'), skill_nature),
    skill_perception       = COALESCE(sqlc.narg('skill_perception'), skill_perception),
    skill_performance      = COALESCE(sqlc.narg('skill_performance'), skill_performance),
    skill_persuasion       = COALESCE(sqlc.narg('skill_persuasion'), skill_persuasion),
    skill_religion         = COALESCE(sqlc.narg('skill_religion'), skill_religion),
    skill_sleight_of_hand  = COALESCE(sqlc.narg('skill_sleight_of_hand'), skill_sleight_of_hand),
    skill_stealth          = COALESCE(sqlc.narg('skill_stealth'), skill_stealth),
    skill_survival         = COALESCE(sqlc.narg('skill_survival'), skill_survival),
    armor_class            = COALESCE(sqlc.narg('armor_class'), armor_class),
    speed                  = COALESCE(sqlc.narg('speed'), speed),
    max_hp                 = COALESCE(sqlc.narg('max_hp'), max_hp),
    hit_dice_type          = COALESCE(sqlc.narg('hit_dice_type'), hit_dice_type),
    hit_dice_remaining     = COALESCE(sqlc.narg('hit_dice_remaining'), hit_dice_remaining),
    inspiration            = COALESCE(sqlc.narg('inspiration'), inspiration),
    attunement_slots       = COALESCE(sqlc.narg('attunement_slots'), attunement_slots),
    training_armor         = COALESCE(sqlc.narg('training_armor'), training_armor),
    training_weapons       = COALESCE(sqlc.narg('training_weapons'), training_weapons),
    training_tools         = COALESCE(sqlc.narg('training_tools'), training_tools),
    training_languages     = COALESCE(sqlc.narg('training_languages'), training_languages),
    copper                 = COALESCE(sqlc.narg('copper'), copper),
    silver                 = COALESCE(sqlc.narg('silver'), silver),
    electrum               = COALESCE(sqlc.narg('electrum'), electrum),
    gold                   = COALESCE(sqlc.narg('gold'), gold),
    platinum               = COALESCE(sqlc.narg('platinum'), platinum),
    conditions             = COALESCE(sqlc.narg('conditions'), conditions),
    resistances            = COALESCE(sqlc.narg('resistances'), resistances),
    vulnerabilities        = COALESCE(sqlc.narg('vulnerabilities'), vulnerabilities),
    immunities             = COALESCE(sqlc.narg('immunities'), immunities),
    personality_traits     = COALESCE(sqlc.narg('personality_traits'), personality_traits),
    ideals                 = COALESCE(sqlc.narg('ideals'), ideals),
    bonds                  = COALESCE(sqlc.narg('bonds'), bonds),
    flaws                  = COALESCE(sqlc.narg('flaws'), flaws),
    notes                  = COALESCE(sqlc.narg('notes'), notes),
    updated_at             = NOW()
WHERE id = sqlc.arg('id')
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
    temp_hp                 = 0,
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
    temp_hp                 = 0,
    updated_at              = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteCharacter :exec
DELETE FROM characters
WHERE id = $1;
