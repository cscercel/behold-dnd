-- +goose Up
CREATE TABLE users (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    username        TEXT        NOT NULL UNIQUE,
    email           TEXT        NOT NULL UNIQUE,
    hashed_password TEXT        NOT NULL,
    role            TEXT        NOT NULL CHECK (role IN ('player', 'dm')),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE characters (
    id         UUID    PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id   UUID    NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    is_npc     BOOLEAN NOT NULL DEFAULT FALSE,

    -- Identity
    name       TEXT    NOT NULL,
    race       TEXT    NOT NULL DEFAULT '',
    class      TEXT    NOT NULL DEFAULT '',
    level      INTEGER NOT NULL DEFAULT 1 CHECK (level BETWEEN 1 AND 20),
    background TEXT    NOT NULL DEFAULT '',
    alignment  TEXT    NOT NULL DEFAULT '',
    xp         INTEGER NOT NULL DEFAULT 0,

    -- Ability scores
    strength     INTEGER NOT NULL DEFAULT 10,
    dexterity    INTEGER NOT NULL DEFAULT 10,
    constitution INTEGER NOT NULL DEFAULT 10,
    intelligence INTEGER NOT NULL DEFAULT 10,
    wisdom       INTEGER NOT NULL DEFAULT 10,
    charisma     INTEGER NOT NULL DEFAULT 10,

    -- Saving throw proficiencies
    save_prof_strength     BOOLEAN NOT NULL DEFAULT FALSE,
    save_prof_dexterity    BOOLEAN NOT NULL DEFAULT FALSE,
    save_prof_constitution BOOLEAN NOT NULL DEFAULT FALSE,
    save_prof_intelligence BOOLEAN NOT NULL DEFAULT FALSE,
    save_prof_wisdom       BOOLEAN NOT NULL DEFAULT FALSE,
    save_prof_charisma     BOOLEAN NOT NULL DEFAULT FALSE,

    -- Skill proficiencies (0 = none, 1 = proficient, 2 = expertise)
    skill_acrobatics      INTEGER NOT NULL DEFAULT 0,
    skill_animal_handling INTEGER NOT NULL DEFAULT 0,
    skill_arcana          INTEGER NOT NULL DEFAULT 0,
    skill_athletics       INTEGER NOT NULL DEFAULT 0,
    skill_deception       INTEGER NOT NULL DEFAULT 0,
    skill_history         INTEGER NOT NULL DEFAULT 0,
    skill_insight         INTEGER NOT NULL DEFAULT 0,
    skill_intimidation    INTEGER NOT NULL DEFAULT 0,
    skill_investigation   INTEGER NOT NULL DEFAULT 0,
    skill_medicine        INTEGER NOT NULL DEFAULT 0,
    skill_nature          INTEGER NOT NULL DEFAULT 0,
    skill_perception      INTEGER NOT NULL DEFAULT 0,
    skill_performance     INTEGER NOT NULL DEFAULT 0,
    skill_persuasion      INTEGER NOT NULL DEFAULT 0,
    skill_religion        INTEGER NOT NULL DEFAULT 0,
    skill_sleight_of_hand INTEGER NOT NULL DEFAULT 0,
    skill_stealth         INTEGER NOT NULL DEFAULT 0,
    skill_survival        INTEGER NOT NULL DEFAULT 0,

    -- Hit points
    max_hp     INTEGER NOT NULL DEFAULT 0,
    current_hp INTEGER NOT NULL DEFAULT 0,
    temp_hp    INTEGER NOT NULL DEFAULT 0,

    -- Combat stats
    armor_class INTEGER NOT NULL DEFAULT 10,
    speed       INTEGER NOT NULL DEFAULT 30,

    -- Hit dice
    hit_dice_type      INTEGER NOT NULL DEFAULT 8,
    hit_dice_remaining INTEGER NOT NULL DEFAULT 1,

    -- Death saves
    death_save_successes INTEGER NOT NULL DEFAULT 0 CHECK (death_save_successes BETWEEN 0 AND 3),
    death_save_failures  INTEGER NOT NULL DEFAULT 0 CHECK (death_save_failures  BETWEEN 0 AND 3),

    -- Spellcasting ability
    spellcasting_ability TEXT   NOT NULL DEFAULT '',

    -- Misc
    inspiration         BOOLEAN NOT NULL DEFAULT FALSE,
    training_armor      TEXT[]  NOT NULL DEFAULT '{}',
    training_weapons    TEXT[]  NOT NULL DEFAULT '{}',
    training_tools      TEXT[]  NOT NULL DEFAULT '{}',
    training_languages  TEXT[]  NOT NULL DEFAULT '{}',
    attunement_slots    INTEGER NOT NULL DEFAULT 3,

    -- Currency
    copper   INTEGER NOT NULL DEFAULT 0,
    silver   INTEGER NOT NULL DEFAULT 0,
    electrum INTEGER NOT NULL DEFAULT 0,
    gold     INTEGER NOT NULL DEFAULT 0,
    platinum INTEGER NOT NULL DEFAULT 0,

    -- Conditions and Defenses
    conditions      TEXT[] NOT NULL DEFAULT '{}',
    resistances     TEXT[] NOT NULL DEFAULT '{}',
    vulnerabilities TEXT[] NOT NULL DEFAULT '{}',
    immunities      TEXT[] NOT NULL DEFAULT '{}',

    -- Roleplay
    personality_traits TEXT NOT NULL DEFAULT '',
    ideals             TEXT NOT NULL DEFAULT '',
    bonds              TEXT NOT NULL DEFAULT '',
    flaws              TEXT NOT NULL DEFAULT '',
    notes              TEXT NOT NULL DEFAULT '',

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE inventory_items (
    id                  UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    character_id        UUID          NOT NULL REFERENCES characters(id) ON DELETE CASCADE,
    name                TEXT          NOT NULL,
    quantity            INTEGER       NOT NULL DEFAULT 1,
    weight              INTEGER       NOT NULL DEFAULT 0,
    description         TEXT          NOT NULL DEFAULT '',
    value               INTEGER       NOT NULL DEFAULT 0,
    is_equipped         BOOLEAN       NOT NULL DEFAULT FALSE,
    requires_attunement BOOLEAN       NOT NULL DEFAULT FALSE,
    is_attuned          BOOLEAN       NOT NULL DEFAULT FALSE,
    created_at          TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE TABLE spells (
    id           UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    character_id UUID        NOT NULL REFERENCES characters(id) ON DELETE CASCADE,
    name         TEXT        NOT NULL,
    level        INTEGER     NOT NULL DEFAULT 0 CHECK (level BETWEEN 0 AND 9),
    school       TEXT        NOT NULL DEFAULT '',
    casting_time TEXT        NOT NULL DEFAULT '',
    range        TEXT        NOT NULL DEFAULT '',
    components   TEXT        NOT NULL DEFAULT '',
    duration     TEXT        NOT NULL DEFAULT '',
    description  TEXT        NOT NULL DEFAULT '',
    is_prepared  BOOLEAN     NOT NULL DEFAULT FALSE,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE spell_slots (
    id           UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    character_id UUID        NOT NULL REFERENCES characters(id) ON DELETE CASCADE,
    spell_level  INTEGER     NOT NULL CHECK (spell_level BETWEEN 1 AND 9),
    total        INTEGER     NOT NULL DEFAULT 0,
    used         INTEGER     NOT NULL DEFAULT 0,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (character_id, spell_level)
);

CREATE TABLE features (
    id           UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    character_id UUID        NOT NULL REFERENCES characters(id) ON DELETE CASCADE,
    name         TEXT        NOT NULL,
    action_type  TEXT        NOT NULL CHECK (action_type IN ('none', 'action', 'bonus_action', 'reaction', 'free')),
    source       TEXT        NOT NULL DEFAULT '',
    description  TEXT        NOT NULL DEFAULT '',
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE combat_encounters (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    name       TEXT        NOT NULL DEFAULT '',
    is_active  BOOLEAN     NOT NULL DEFAULT FALSE,
    round      INTEGER     NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE combat_participants (
    id              UUID    PRIMARY KEY DEFAULT gen_random_uuid(),
    encounter_id    UUID    NOT NULL REFERENCES combat_encounters(id) ON DELETE CASCADE,
    character_id    UUID    NOT NULL REFERENCES characters(id),
    name            TEXT    NOT NULL,
    initiative      INTEGER NOT NULL DEFAULT 0,
    current_hp      INTEGER NOT NULL DEFAULT 0,
    max_hp          INTEGER NOT NULL DEFAULT 0,
    temp_hp         INTEGER NOT NULL DEFAULT 0,
    armor_class     INTEGER NOT NULL DEFAULT 10,
    speed           INTEGER NOT NULL DEFAULT 30,
    conditions      TEXT[]  NOT NULL DEFAULT '{}',
    concentration   BOOLEAN NOT NULL DEFAULT FALSE,
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    notes           TEXT    NOT NULL DEFAULT '',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE combat_participants;
DROP TABLE combat_encounters;
DROP TABLE features;
DROP TABLE spell_slots;
DROP TABLE spells;
DROP TABLE inventory_items;
DROP TABLE characters;
DROP TABLE users;
