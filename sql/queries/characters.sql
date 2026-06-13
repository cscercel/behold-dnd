-- name: CreateCharacter :one
INSERT INTO characters (
    owner_id, is_npc,
    name, race, class, level, background, alignment, xp,
    strength, dexterity, constitution, intelligence, wisdom, charisma,
    save_prof_strength, save_prof_dexterity, save_prof_constitution,
    save_prof_intelligence, save_prof_wisdom, save_prof_charisma,
    skill_acrobatics, skill_animal_handling, skill_arcana, skill_athletics,
    skill_deception, skill_history, skill_insight, skill_intimidation,
    skill_investigation, skill_medicine, skill_nature, skill_perception,
    skill_performance, skill_persuasion, skill_religion, skill_sleight_of_hand,
    skill_stealth, skill_survival,
    max_hp, current_hp, temp_hp, armor_class, speed,
    hit_dice_type, hit_dice_remaining,
    inspiration, attunement_slots,
    training_armor, training_weapons, training_tools, training_languages,
    copper, silver, electrum, gold, platinum,
    conditions, resistances, vulnerabilities, immunities,
    personality_traits, ideals, bonds, flaws, notes,
    spellcasting_ability
) VALUES (
    $1,  $2,  $3,  $4,  $5,  $6,  $7,  $8,  $9,
    $10, $11, $12, $13, $14, $15,
    $16, $17, $18, $19, $20, $21,
    $22, $23, $24, $25, $26, $27, $28, $29, $30,
    $31, $32, $33, $34, $35, $36, $37, $38, $39,
    $40, $41, $42, $43, $44, $45, $46,
    $47, $48,
    $49, $50, $51, $52,
    $53, $54, $55, $56, $57,
    $58, $59, $60, $61,
    $62, $63, $64, $65, $66,
    $67
)
RETURNING *;

-- name: GetCharacter :one
SELECT 
    *,
    CEILING(1 + level / 4.0) AS proficiency_bonus,

    FLOOR((strength - 10) / 2.0) AS strength_modifier,
    FLOOR((dexterity - 10) / 2.0) AS dexterity_modifier,
    FLOOR((constitution - 10) / 2.0) AS constitution_modifier,
    FLOOR((intelligence - 10) / 2.0) AS intelligence_modifier,
    FLOOR((wisdom - 10) / 2.0) AS wisdom_modifier,
    FLOOR((charisma - 10) / 2.0) AS charisma_modifier,

    FLOOR((strength - 10) / 2.0) + save_prof_strength::int * CEILING(1 + level / 4.0) AS strength_saving_throw_modifier,
    FLOOR((dexterity - 10) / 2.0) + save_prof_dexterity::int * CEILING(1 + level / 4.0) AS dexterity_saving_throw_modifier,
    FLOOR((constitution - 10) / 2.0) + save_prof_constitution::int * CEILING(1 + level / 4.0) AS constitution_saving_throw_modifier,
    FLOOR((intelligence - 10) / 2.0) + save_prof_intelligence::int * CEILING(1 + level / 4.0) AS intelligence_saving_throw_modifier,
    FLOOR((wisdom - 10) / 2.0) + save_prof_wisdom::int * CEILING(1 + level / 4.0) AS wisdom_saving_throw_modifier,
    FLOOR((charisma - 10) / 2.0) + save_prof_charisma::int * CEILING(1 + level / 4.0) AS charisma_saving_throw_modifier,

    FLOOR((dexterity - 10) / 2.0) + skill_acrobatics * CEILING(1 + level / 4.0) AS skill_acrobatics_modifier,
    FLOOR((wisdom - 10) / 2.0) + skill_animal_handling * CEILING(1 + level / 4.0) AS skill_animal_handling_modifier,
    FLOOR((intelligence - 10) / 2.0) + skill_arcana * CEILING(1 + level / 4.0) AS skill_arcana_modifier,
    FLOOR((strength - 10) / 2.0) + skill_athletics * CEILING(1 + level / 4.0) AS skill_athletics_modifier,
    FLOOR((charisma - 10) / 2.0) + skill_deception * CEILING(1 + level / 4.0) AS skill_deception_modifier,
    FLOOR((intelligence - 10) / 2.0) + skill_history * CEILING(1 + level / 4.0) AS skill_history_modifier,
    FLOOR((wisdom - 10) / 2.0) + skill_insight * CEILING(1 + level / 4.0) AS skill_insight_modifier,
    FLOOR((charisma - 10) / 2.0) + skill_intimidation * CEILING(1 + level / 4.0) AS skill_intimidation_modifier,
    FLOOR((intelligence - 10) / 2.0) + skill_investigation * CEILING(1 + level / 4.0) AS skill_investigation_modifier,
    FLOOR((wisdom - 10) / 2.0) + skill_medicine * CEILING(1 + level / 4.0) AS skill_medicine_modifier,
    FLOOR((intelligence - 10) / 2.0) + skill_nature * CEILING(1 + level / 4.0) AS skill_nature_modifier,
    FLOOR((wisdom - 10) / 2.0) + skill_perception * CEILING(1 + level / 4.0) AS skill_perception_modifier,
    FLOOR((charisma - 10) / 2.0) + skill_performance * CEILING(1 + level / 4.0) AS skill_performance_modifier,
    FLOOR((charisma - 10) / 2.0) + skill_persuasion * CEILING(1 + level / 4.0) AS skill_persuasion_modifier,
    FLOOR((intelligence - 10) / 2.0) + skill_religion * CEILING(1 + level / 4.0) AS skill_religion_modifier,
    FLOOR((dexterity - 10) / 2.0) + skill_sleight_of_hand * CEILING(1 + level / 4.0) AS skill_sleight_of_hand_modifier,
    FLOOR((dexterity - 10) / 2.0) + skill_stealth * CEILING(1 + level / 4.0) AS skill_stealth_modifier,
    FLOOR((wisdom - 10) / 2.0) + skill_survival * CEILING(1 + level / 4.0) AS skill_survival_modifier,

    10 + FLOOR((wisdom - 10) / 2.0) + skill_perception * CEILING(1 + level / 4.0) AS passive_perception,
    10 + FLOOR((intelligence - 10) / 2.0) + skill_investigation * CEILING(1 + level / 4.0) AS passive_investigation,
    10 + FLOOR((wisdom - 10) / 2.0) + skill_insight * CEILING(1 + level / 4.0) AS passive_insight,

    FLOOR((dexterity - 10) / 2.0) AS initiative
FROM characters
WHERE id = $1;

-- name: ListCharacters :many
SELECT * FROM characters
ORDER BY name;

-- name: ListUserCharacters :many
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

-- name: UpdateCharacterInfo :one
UPDATE characters
SET
    name                   = COALESCE(sqlc.narg('name'), name),
    race                   = COALESCE(sqlc.narg('race'), race),
    background             = COALESCE(sqlc.narg('background'), background),
    alignment              = COALESCE(sqlc.narg('alignment'), alignment),
    inspiration            = COALESCE(sqlc.narg('inspiration'), inspiration),
    speed                  = COALESCE(sqlc.narg('speed'), speed),
    ideals                 = COALESCE(sqlc.narg('ideals'), ideals),
    bonds                  = COALESCE(sqlc.narg('bonds'), bonds),
    flaws                  = COALESCE(sqlc.narg('flaws'), flaws),
    notes                  = COALESCE(sqlc.narg('notes'), notes),
    updated_at             = NOW()
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: UpdateCharacterAbilityScores :one
UPDATE characters
SET
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
    updated_at             = NOW()
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: UpdateCharacterSkills :one
UPDATE characters
SET
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
    updated_at             = NOW()
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: UpdateCharacterLevel :one
UPDATE characters
SET
    class                  = COALESCE(sqlc.narg('class'), class),
    level                  = COALESCE(sqlc.narg('level'), level),
    xp                     = COALESCE(sqlc.narg('xp'), xp),
    max_hp                 = COALESCE(sqlc.narg('max_hp'), max_hp),
    hit_dice_type          = COALESCE(sqlc.narg('hit_dice_type'), hit_dice_type),
    hit_dice_remaining     = COALESCE(sqlc.narg('hit_dice_remaining'), hit_dice_remaining),
    updated_at             = NOW()
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: UpdateCharacterTraining :one
UPDATE characters
SET
    armor_class            = COALESCE(sqlc.narg('armor_class'), armor_class),
    attunement_slots       = COALESCE(sqlc.narg('attunement_slots'), attunement_slots),
    training_armor         = COALESCE(sqlc.narg('training_armor'), training_armor),
    training_weapons       = COALESCE(sqlc.narg('training_weapons'), training_weapons),
    training_tools         = COALESCE(sqlc.narg('training_tools'), training_tools),
    training_languages     = COALESCE(sqlc.narg('training_languages'), training_languages),
    spellcasting_ability   = COALESCE(sqlc.narg('spellcasting_ability'), spellcasting_ability),
    resistances            = COALESCE(sqlc.narg('resistances'), resistances),
    vulnerabilities        = COALESCE(sqlc.narg('vulnerabilities'), vulnerabilities),
    immunities             = COALESCE(sqlc.narg('immunities'), immunities),
    updated_at             = NOW()
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: UpdateCharacterCurrency :one
UPDATE characters
SET
    copper                 = COALESCE(sqlc.narg('copper'), copper),
    silver                 = COALESCE(sqlc.narg('silver'), silver),
    electrum               = COALESCE(sqlc.narg('electrum'), electrum),
    gold                   = COALESCE(sqlc.narg('gold'), gold),
    platinum               = COALESCE(sqlc.narg('platinum'), platinum),
    conditions             = COALESCE(sqlc.narg('conditions'), conditions),
    updated_at             = NOW()
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: UpdateCharacterHP :one
UPDATE characters
SET
    current_hp = $2,
    temp_hp    = $3,
    updated_at = NOW()
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
