export interface User {
  id: string
  username: string
  email: string
  role: 'player' | 'dm'
}

export interface Character {
  id: string
  owner_id: string | null
  is_npc: boolean

  // Identity
  name: string
  race: string
  class: string
  level: number
  background: string
  alignment: string
  xp: number

  // Ability Scores
  strength: number
  dexterity: number
  constitution: number
  intelligence: number
  wisdom: number
  charisma: number

  // Saving Throw proficiencies
  save_prof_strength: boolean
  save_prof_dexterity: boolean
  save_prof_constitution: boolean
  save_prof_intelligence: boolean
  save_prof_wisdom: boolean
  save_prof_charisma: boolean

  // Skill proficiencies (0 = none, 1 = proficient, 2 = expertise)
  skill_acrobatics: number
  skill_animal_handling: number
  skill_arcana: number
  skill_athletics: number
  skill_deception: number
  skill_history: number
  skill_insight: number
  skill_intimidation: number
  skill_investigation: number
  skill_medicine: number
  skill_nature: number
  skill_perception: number
  skill_performance: number
  skill_persuasion: number
  skill_religion: number
  skill_sleight_of_hand: number
  skill_stealth: number
  skill_survival: number

  // HP
  max_hp: number
  current_hp: number
  temp_hp: number

  // Death Saves
  death_saves_successes: number
  death_saves_failures: number

  // Misc
  inspiration: boolean
  attunement_slots: number
  training_armor: string[]
  training_weapons: string[]
  training_tools: string[]
  training_languages: string[]

  // Currency
  copper: number
  silver: number
  electrum: number
  gold: number
  platinum: number

  // Conditions and defenses
  conditions: string[]
  resistances: string[]
  vulnerabilities: string[]
  immunities: string[]

  // Flavor
  personality_traits: string
  ideals: string
  bonds: string
  flaws: string
  notes: string

  created_at: string
  updated_at: string

}

export interface InventoryItem {
  id: string
  character_id: string
  name: string
  quantity: number
  weight: number
  description: string
  is_equipped: boolean
  requires_attunement: boolean
  is_attuned: boolean
  created_at: string
}

export interface Spell {
  id: string
  character_id: string
  name: string
  level: number
  school: string
  casting_time: string
  range: string
  components: string
  duration: string
  description: string
  is_prepared: boolean
  created_at: string
}

export interface SpellSlot {
  id: string
  character_id: string
  spell_level: number
  total: number
  used: number
}

export interface CombatEncounter {
  id: string
  name: string
  is_active: boolean
  round: number
  created_at: string
}

export interface CombatParticipant {
  id: string
  encounter_id: string
  character_id: string | null
  name: string
  initiative: number
  current_hp: number
  max_hp: number
  temp_hp: number
  armor_class: number
  speed: number
  conditions: string[]
  concentration: boolean
  is_active: boolean
  notes: string
}

// Partial update payloads — all fields optional
export type UpdateCharacterPayload = Partial<Omit<Character,
  'id' | 'owner_id' | 'created_at' | 'updated_at'
>>

export type UpdateInventoryItemPayload = Partial<Omit<InventoryItem,
  'id' | 'character_id' | 'created_at'
>>

export type UpdateSpellPayload = Partial<Omit<Spell,
  'id' | 'character_id' | 'created_at'
>>
