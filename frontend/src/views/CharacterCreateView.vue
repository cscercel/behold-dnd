<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useCharactersStore } from '@/stores/characters'
import { featuresAPI } from '@/api/features'

const router = useRouter()
const auth = useAuthStore()
const store = useCharactersStore()

const step = ref(1)
const TOTAL_STEPS = 6
const saving = ref(false)
const error = ref<string | null>(null)

const form = ref({
  // Step 1 — Identity
  name: '',
  race: '',
  raceCustom: '',
  classSelect: '',
  classCustom: '',
  class: '',
  level: 1,
  background: '',
  alignment: '',
  xp: 0,
  is_npc: false,

  // Step 2 — Ability scores
  strength: 10,
  dexterity: 10,
  constitution: 10,
  intelligence: 10,
  wisdom: 10,
  charisma: 10,

  // Step 3 — Saves & Skills
  save_prof_strength: false,
  save_prof_dexterity: false,
  save_prof_constitution: false,
  save_prof_intelligence: false,
  save_prof_wisdom: false,
  save_prof_charisma: false,
  skill_acrobatics: 0,
  skill_animal_handling: 0,
  skill_arcana: 0,
  skill_athletics: 0,
  skill_deception: 0,
  skill_history: 0,
  skill_insight: 0,
  skill_intimidation: 0,
  skill_investigation: 0,
  skill_medicine: 0,
  skill_nature: 0,
  skill_perception: 0,
  skill_performance: 0,
  skill_persuasion: 0,
  skill_religion: 0,
  skill_sleight_of_hand: 0,
  skill_stealth: 0,
  skill_survival: 0,

  // Step 4 — Combat
  max_hp: 8,
  current_hp: 8,
  armor_class: 10,
  speed: 30,
  hit_dice_type: 8,
  hit_dice_remaining: 1,

  // Step 5 — Spellcasting & Training
  spellcasting_ability: '',
  training_armor: [] as string[],
  training_weapons: [] as string[],
  training_tools: [] as string[],
  training_languages: [] as string[],
  copper: 0,
  silver: 0,
  electrum: 0,
  gold: 0,
  platinum: 0,

  // Step 6 — Personality
  personality_traits: '',
  ideals: '',
  bonds: '',
  flaws: '',
  notes: '',
})

// ── Step validation ──────────────────────────────────────────
const canProceed = computed(() => {
  if (step.value === 1) {
    const hasName = form.value.name.trim() !== ''
    const hasClass = form.value.classSelect === 'custom'
      ? form.value.classCustom.trim() !== ''
      : form.value.classSelect !== ''
    return hasName && hasClass
  }
  return true
})

// ── Ability modifier display ─────────────────────────────────
function mod(score: number): string {
  const m = Math.floor((score - 10) / 2)
  return m >= 0 ? `+${m}` : `${m}`
}

// ── Class change handler ─────────────────────────────────────
function onClassChange() {
  if (form.value.classSelect !== 'custom') {
    form.value.class = form.value.classSelect
    syncHitDice()
  }
}

// ── Skill cycle helper ───────────────────────────────────────
function cycleSkill(field: string) {
  const cur = (form.value as any)[field] as number
  ;(form.value as any)[field] = (cur + 1) % 3
}
function skillLabel(level: number) {
  if (level === 2) return 'E'
  if (level === 1) return '●'
  return '○'
}

// ── Training array helpers ───────────────────────────────────
function toggleTraining(arr: string[], value: string) {
  const idx = arr.indexOf(value)
  if (idx === -1) arr.push(value)
  else arr.splice(idx, 1)
}

// ── HP sync ──────────────────────────────────────────────────
function syncHP() {
  form.value.current_hp = form.value.max_hp
}

// ── Hit dice sync ────────────────────────────────────────────
function syncHitDice() {
  form.value.hit_dice_remaining = form.value.level
}

// ── Submit ───────────────────────────────────────────────────
async function submit() {
  // Resolve custom race/class
  if (form.value.race === 'custom') form.value.race = form.value.raceCustom
  if (form.value.classSelect === 'custom') form.value.class = form.value.classCustom
  else form.value.class = form.value.classSelect

  saving.value = true
  error.value = null
  try {
    const payload = {
      ...form.value,
      owner_id: form.value.is_npc ? null : auth.user?.id,
      current_hp: form.value.max_hp,
      temp_hp: 0,
      conditions: [],
      resistances: [],
      vulnerabilities: [],
      immunities: [],
      inspiration: false,
      attunement_slots: 3,
      death_save_successes: 0,
      death_save_failures: 0,
    }

    const created = await store.create(payload)

    // Seed base actions every character has
    await Promise.all(BASE_ACTIONS.map(a => featuresAPI.create(created.id, a)))

    router.push({ name: 'character-sheet', params: { id: created.id } })
  } catch (e: any) {
    error.value = e.message ?? 'Failed to create character'
  } finally {
    saving.value = false
  }
}

// ── Step labels ──────────────────────────────────────────────
const STEPS = [
  'Identity',
  'Ability Scores',
  'Skills & Saves',
  'Combat',
  'Training',
  'Personality',
]

// ── Static data ──────────────────────────────────────────────
const ALIGNMENTS = [
  'Lawful Good', 'Neutral Good', 'Chaotic Good',
  'Lawful Neutral', 'True Neutral', 'Chaotic Neutral',
  'Lawful Evil', 'Neutral Evil', 'Chaotic Evil',
]

const RACES = [
  'Dragonborn', 'Dwarf', 'Elf', 'Gnome', 'Half-Elf',
  'Half-Orc', 'Halfling', 'Human', 'Tiefling',
]

const CLASSES = [
  'Artificer', 'Barbarian', 'Bard', 'Cleric', 'Druid',
  'Fighter', 'Monk', 'Paladin', 'Ranger', 'Rogue',
  'Sorcerer', 'Warlock', 'Wizard',
]

const HIT_DICE = [6, 8, 10, 12]

const ABILITIES = [
  { key: 'strength',     label: 'Strength',     abbr: 'STR' },
  { key: 'dexterity',    label: 'Dexterity',    abbr: 'DEX' },
  { key: 'constitution', label: 'Constitution', abbr: 'CON' },
  { key: 'intelligence', label: 'Intelligence', abbr: 'INT' },
  { key: 'wisdom',       label: 'Wisdom',       abbr: 'WIS' },
  { key: 'charisma',     label: 'Charisma',     abbr: 'CHA' },
] as const

const SAVES = [
  { key: 'strength',     label: 'Strength',     field: 'save_prof_strength'     },
  { key: 'dexterity',    label: 'Dexterity',    field: 'save_prof_dexterity'    },
  { key: 'constitution', label: 'Constitution', field: 'save_prof_constitution' },
  { key: 'intelligence', label: 'Intelligence', field: 'save_prof_intelligence' },
  { key: 'wisdom',       label: 'Wisdom',       field: 'save_prof_wisdom'       },
  { key: 'charisma',     label: 'Charisma',     field: 'save_prof_charisma'     },
] as const

const SKILLS = [
  { label: 'Acrobatics',      field: 'skill_acrobatics',      ability: 'DEX' },
  { label: 'Animal Handling', field: 'skill_animal_handling',  ability: 'WIS' },
  { label: 'Arcana',          field: 'skill_arcana',           ability: 'INT' },
  { label: 'Athletics',       field: 'skill_athletics',        ability: 'STR' },
  { label: 'Deception',       field: 'skill_deception',        ability: 'CHA' },
  { label: 'History',         field: 'skill_history',          ability: 'INT' },
  { label: 'Insight',         field: 'skill_insight',          ability: 'WIS' },
  { label: 'Intimidation',    field: 'skill_intimidation',     ability: 'CHA' },
  { label: 'Investigation',   field: 'skill_investigation',    ability: 'INT' },
  { label: 'Medicine',        field: 'skill_medicine',         ability: 'WIS' },
  { label: 'Nature',          field: 'skill_nature',           ability: 'INT' },
  { label: 'Perception',      field: 'skill_perception',       ability: 'WIS' },
  { label: 'Performance',     field: 'skill_performance',      ability: 'CHA' },
  { label: 'Persuasion',      field: 'skill_persuasion',       ability: 'CHA' },
  { label: 'Religion',        field: 'skill_religion',         ability: 'INT' },
  { label: 'Sleight of Hand', field: 'skill_sleight_of_hand',  ability: 'DEX' },
  { label: 'Stealth',         field: 'skill_stealth',          ability: 'DEX' },
  { label: 'Survival',        field: 'skill_survival',         ability: 'WIS' },
] as const

const ARMOR_OPTIONS    = ['Light Armor', 'Medium Armor', 'Heavy Armor', 'Shields']
const WEAPON_OPTIONS   = ['Simple Weapons', 'Martial Weapons', 'Hand Crossbows', 'Longswords', 'Rapiers', 'Shortswords']
const TOOL_OPTIONS     = ["Thieves' Tools", "Artisan's Tools", 'Herbalism Kit', 'Musical Instrument', "Navigator's Tools", "Poisoner's Kit"]
const LANGUAGE_OPTIONS = ['Common', 'Dwarvish', 'Elvish', 'Giant', 'Gnomish', 'Goblin', 'Halfling', 'Orc', 'Abyssal', 'Celestial', 'Draconic', 'Deep Speech', 'Infernal', 'Primordial', 'Sylvan', 'Undercommon']

const BASE_ACTIONS = [
  // Actions
  { name: 'Attack',          action_type: 'action',       source: 'Basic Rules', description: 'Make one melee or ranged attack, or grapple.' },
  { name: 'Cast a Spell',    action_type: 'action',       source: 'Basic Rules', description: 'Cast a spell with a casting time of 1 action.' },
  { name: 'Dash',            action_type: 'action',       source: 'Basic Rules', description: 'Gain extra movement equal to your speed for the turn.' },
  { name: 'Disengage',       action_type: 'action',       source: 'Basic Rules', description: 'Your movement does not provoke opportunity attacks for the rest of the turn.' },
  { name: 'Dodge',           action_type: 'action',       source: 'Basic Rules', description: 'Until the start of your next turn, attacks against you have disadvantage and you have advantage on Dexterity saving throws.' },
  { name: 'Help',            action_type: 'action',       source: 'Basic Rules', description: 'Give an ally advantage on their next ability check or attack roll against a creature within 5ft of you.' },
  { name: 'Hide',            action_type: 'action',       source: 'Basic Rules', description: 'Make a Stealth check to hide. Requires obscurement from enemies.' },
  { name: 'Ready',           action_type: 'action',       source: 'Basic Rules', description: 'Prepare a trigger and a reaction to use when the trigger occurs before the start of your next turn.' },
  { name: 'Search',          action_type: 'action',       source: 'Basic Rules', description: 'Devote your attention to finding something. Make a Perception or Investigation check.' },
  { name: 'Use an Object',   action_type: 'action',       source: 'Basic Rules', description: 'Interact with a second object or use a special object feature.' },
  // Bonus actions
  { name: 'Off-Hand Attack', action_type: 'bonus_action', source: 'Basic Rules', description: 'When you take the Attack action with a light melee weapon, you can use your bonus action to attack with a different light melee weapon in your other hand. No ability modifier is added to damage unless it is negative.' },
  // Reactions
  { name: 'Opportunity Attack', action_type: 'reaction', source: 'Basic Rules', description: 'When a hostile creature you can see moves out of your reach, make one melee attack against it.' },
  { name: 'Ready Trigger',      action_type: 'reaction', source: 'Basic Rules', description: 'Execute the action you prepared with the Ready action when your chosen trigger occurs.' },
]
</script>

<template>
  <div class="create-page">
    <!-- Header -->
    <header class="page-header">
      <div class="header-left">
        <button class="btn btn-ghost back-btn" @click="router.back()">← Back</button>
        <span class="header-icon">⚔</span>
        <span class="header-title">Beyond DnD</span>
      </div>
      <button class="btn btn-ghost" @click="auth.logout()">Logout</button>
    </header>

    <main class="create-main">
      <!-- Page title -->
      <div class="create-hero fade-in">
        <h1 class="create-title">{{ auth.isDM && form.is_npc ? 'Create NPC' : 'Create Character' }}</h1>
        <p class="create-sub">Step {{ step }} of {{ TOTAL_STEPS }} — {{ STEPS[step - 1] }}</p>
      </div>

      <!-- Progress bar -->
      <div class="progress-bar">
        <div class="progress-fill" :style="{ width: ((step - 1) / (TOTAL_STEPS - 1)) * 100 + '%' }" />
      </div>

      <!-- Step indicators -->
      <div class="step-indicators">
        <div
          v-for="(label, i) in STEPS"
          :key="i"
          class="step-indicator"
          :class="{ active: step === i + 1, done: step > i + 1 }"
          @click="step > i + 1 ? step = i + 1 : null"
        >
          <div class="step-dot">{{ step > i + 1 ? '✓' : i + 1 }}</div>
          <span class="step-label">{{ label }}</span>
        </div>
      </div>

      <div class="ornament" />

      <!-- ── Step 1: Identity ── -->
      <div v-if="step === 1" class="step-content fade-in">
        <!-- NPC toggle for DM only -->
        <div v-if="auth.isDM" class="npc-toggle-row">
          <label class="toggle-label">
            <div class="toggle-track" :class="{ active: form.is_npc }" @click="form.is_npc = !form.is_npc">
              <div class="toggle-thumb" />
            </div>
            <span>Creating an NPC / Monster</span>
          </label>
        </div>

        <div class="form-grid">
          <div class="field field-wide">
            <label class="field-label">Character Name *</label>
            <input v-model="form.name" class="input" type="text" placeholder="Aldric Stonehammer" />
          </div>

          <div class="field">
            <label class="field-label">Race</label>
            <select v-model="form.race" class="input">
              <option value="">— Select Race —</option>
              <option v-for="r in RACES" :key="r">{{ r }}</option>
              <option value="custom">Other / Custom</option>
            </select>
            <input
              v-if="form.race === 'custom'"
              v-model="form.raceCustom"
              class="input"
              style="margin-top:6px"
              type="text"
              placeholder="Enter race..."
            />
          </div>

          <div class="field">
            <label class="field-label">Class *</label>
            <select v-model="form.classSelect" class="input" @change="onClassChange">
              <option value="">— Select Class —</option>
              <option v-for="c in CLASSES" :key="c">{{ c }}</option>
              <option value="custom">Other / Custom</option>
            </select>
            <input
              v-if="form.classSelect === 'custom'"
              v-model="form.classCustom"
              class="input"
              style="margin-top:6px"
              type="text"
              placeholder="Enter class..."
            />
          </div>

          <div class="field field-sm">
            <label class="field-label">Level</label>
            <input v-model.number="form.level" class="input" type="number" min="1" max="20" @change="syncHitDice" />
          </div>

          <div class="field">
            <label class="field-label">Background</label>
            <input v-model="form.background" class="input" type="text" placeholder="Soldier, Acolyte..." />
          </div>

          <div class="field">
            <label class="field-label">Alignment</label>
            <select v-model="form.alignment" class="input">
              <option value="">— Select —</option>
              <option v-for="a in ALIGNMENTS" :key="a">{{ a }}</option>
            </select>
          </div>

          <div class="field field-sm">
            <label class="field-label">Starting XP</label>
            <input v-model.number="form.xp" class="input" type="number" min="0" />
          </div>
        </div>
      </div>

      <!-- ── Step 2: Ability Scores ── -->
      <div v-if="step === 2" class="step-content fade-in">
        <p class="step-hint">Enter your ability scores. The modifier is calculated automatically.</p>
        <div class="ability-grid">
          <div v-for="ab in ABILITIES" :key="ab.key" class="ability-card">
            <span class="ability-abbr">{{ ab.abbr }}</span>
            <span class="ability-full">{{ ab.label }}</span>
            <input
              v-model.number="(form as any)[ab.key]"
              class="ability-input"
              type="number"
              min="1"
              max="30"
            />
            <span class="ability-mod">{{ mod((form as any)[ab.key]) }}</span>
          </div>
        </div>
      </div>

      <!-- ── Step 3: Skills & Saves ── -->
      <div v-if="step === 3" class="step-content fade-in">
        <div class="two-col">
          <div class="sub-panel">
            <h3 class="sub-panel-title">Saving Throws</h3>
            <p class="step-hint">Check the abilities you are proficient in.</p>
            <div class="save-list">
              <label v-for="save in SAVES" :key="save.key" class="save-check-row">
                <input
                  type="checkbox"
                  :checked="(form as any)[save.field]"
                  @change="(form as any)[save.field] = !((form as any)[save.field])"
                  class="prof-checkbox"
                />
                <span class="save-check-label">{{ save.label }}</span>
              </label>
            </div>
          </div>

          <div class="sub-panel">
            <h3 class="sub-panel-title">Skills</h3>
            <p class="step-hint">Click to cycle: None → Proficient (●) → Expertise (E)</p>
            <div class="skill-list">
              <div v-for="skill in SKILLS" :key="skill.field" class="skill-create-row">
                <span
                  class="skill-prof-btn"
                  :class="`level-${(form as any)[skill.field]}`"
                  @click="cycleSkill(skill.field)"
                >{{ skillLabel((form as any)[skill.field]) }}</span>
                <span class="skill-create-name">{{ skill.label }}</span>
                <span class="skill-create-ability">{{ skill.ability }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- ── Step 4: Combat ── -->
      <div v-if="step === 4" class="step-content fade-in">
        <div class="form-grid">
          <div class="field">
            <label class="field-label">Max HP</label>
            <input v-model.number="form.max_hp" class="input" type="number" min="1" @input="syncHP" />
          </div>

          <div class="field">
            <label class="field-label">Armor Class</label>
            <input v-model.number="form.armor_class" class="input" type="number" min="1" />
          </div>

          <div class="field">
            <label class="field-label">Speed (ft)</label>
            <input v-model.number="form.speed" class="input" type="number" min="0" step="5" />
          </div>

          <div class="field">
            <label class="field-label">Hit Dice Type</label>
            <select v-model.number="form.hit_dice_type" class="input">
              <option v-for="d in HIT_DICE" :key="d" :value="d">d{{ d }}</option>
            </select>
          </div>

          <div class="field">
            <label class="field-label">Hit Dice Remaining</label>
            <input v-model.number="form.hit_dice_remaining" class="input" type="number" min="0" :max="form.level" />
          </div>
        </div>

        <div class="combat-preview">
          <div class="preview-stat"><span class="preview-val">{{ form.max_hp }}</span><span class="preview-label">HP</span></div>
          <div class="preview-stat"><span class="preview-val">{{ form.armor_class }}</span><span class="preview-label">AC</span></div>
          <div class="preview-stat"><span class="preview-val">{{ form.speed }}ft</span><span class="preview-label">Speed</span></div>
          <div class="preview-stat"><span class="preview-val">{{ form.hit_dice_remaining }}d{{ form.hit_dice_type }}</span><span class="preview-label">Hit Dice</span></div>
        </div>
      </div>

      <!-- ── Step 5: Spellcasting & Training ── -->
      <div v-if="step === 5" class="step-content fade-in">
        <div class="form-grid">
          <div class="field">
            <label class="field-label">Spellcasting Ability</label>
            <select v-model="form.spellcasting_ability" class="input">
              <option value="">None (non-spellcaster)</option>
              <option v-for="ab in ABILITIES" :key="ab.key" :value="ab.key">{{ ab.label }}</option>
            </select>
          </div>
        </div>

        <div class="training-grid">
          <div class="training-group">
            <h3 class="sub-panel-title">Armor</h3>
            <label v-for="opt in ARMOR_OPTIONS" :key="opt" class="check-option">
              <input type="checkbox" :checked="form.training_armor.includes(opt)" @change="toggleTraining(form.training_armor, opt)" class="prof-checkbox" />
              {{ opt }}
            </label>
          </div>

          <div class="training-group">
            <h3 class="sub-panel-title">Weapons</h3>
            <label v-for="opt in WEAPON_OPTIONS" :key="opt" class="check-option">
              <input type="checkbox" :checked="form.training_weapons.includes(opt)" @change="toggleTraining(form.training_weapons, opt)" class="prof-checkbox" />
              {{ opt }}
            </label>
          </div>

          <div class="training-group">
            <h3 class="sub-panel-title">Tools</h3>
            <label v-for="opt in TOOL_OPTIONS" :key="opt" class="check-option">
              <input type="checkbox" :checked="form.training_tools.includes(opt)" @change="toggleTraining(form.training_tools, opt)" class="prof-checkbox" />
              {{ opt }}
            </label>
          </div>

          <div class="training-group">
            <h3 class="sub-panel-title">Languages</h3>
            <label v-for="opt in LANGUAGE_OPTIONS" :key="opt" class="check-option">
              <input type="checkbox" :checked="form.training_languages.includes(opt)" @change="toggleTraining(form.training_languages, opt)" class="prof-checkbox" />
              {{ opt }}
            </label>
          </div>
        </div>

        <div class="ornament" style="margin: 20px 0" />

        <h3 class="sub-panel-title" style="margin-bottom:12px">Starting Currency</h3>
        <div class="currency-form">
          <div class="currency-field"><label class="currency-label copper">CP</label><input v-model.number="form.copper" class="input currency-input" type="number" min="0" /></div>
          <div class="currency-field"><label class="currency-label silver">SP</label><input v-model.number="form.silver" class="input currency-input" type="number" min="0" /></div>
          <div class="currency-field"><label class="currency-label electrum">EP</label><input v-model.number="form.electrum" class="input currency-input" type="number" min="0" /></div>
          <div class="currency-field"><label class="currency-label gold">GP</label><input v-model.number="form.gold" class="input currency-input" type="number" min="0" /></div>
          <div class="currency-field"><label class="currency-label platinum">PP</label><input v-model.number="form.platinum" class="input currency-input" type="number" min="0" /></div>
        </div>
      </div>

      <!-- ── Step 6: Personality ── -->
      <div v-if="step === 6" class="step-content fade-in">
        <p class="step-hint">All optional — you can fill these in later on the character sheet.</p>
        <div class="personality-form">
          <div class="field">
            <label class="field-label">Personality Traits</label>
            <textarea v-model="form.personality_traits" class="input textarea" rows="2" placeholder="I always have a plan for when things go wrong..." />
          </div>
          <div class="field">
            <label class="field-label">Ideals</label>
            <textarea v-model="form.ideals" class="input textarea" rows="2" placeholder="Honor. I don't steal from others in the trade..." />
          </div>
          <div class="field">
            <label class="field-label">Bonds</label>
            <textarea v-model="form.bonds" class="input textarea" rows="2" placeholder="I'm trying to pay off an old debt..." />
          </div>
          <div class="field">
            <label class="field-label">Flaws</label>
            <textarea v-model="form.flaws" class="input textarea" rows="2" placeholder="I have a tell when I'm lying..." />
          </div>
          <div class="field">
            <label class="field-label">Notes</label>
            <textarea v-model="form.notes" class="input textarea" rows="3" placeholder="Additional notes about your character..." />
          </div>
        </div>

        <div class="ornament" style="margin: 20px 0" />

        <!-- Summary -->
        <div class="summary-card">
          <h3 class="summary-title">Summary</h3>
          <div class="summary-grid">
            <div class="summary-item"><span class="summary-label">Name</span><span class="summary-val">{{ form.name || '—' }}</span></div>
            <div class="summary-item"><span class="summary-label">Class</span><span class="summary-val">{{ (form.classSelect === 'custom' ? form.classCustom : form.classSelect) || '—' }} {{ form.level }}</span></div>
            <div class="summary-item"><span class="summary-label">Race</span><span class="summary-val">{{ (form.race === 'custom' ? form.raceCustom : form.race) || '—' }}</span></div>
            <div class="summary-item"><span class="summary-label">HP</span><span class="summary-val">{{ form.max_hp }}</span></div>
            <div class="summary-item"><span class="summary-label">AC</span><span class="summary-val">{{ form.armor_class }}</span></div>
            <div class="summary-item"><span class="summary-label">Speed</span><span class="summary-val">{{ form.speed }}ft</span></div>
            <div class="summary-item"><span class="summary-label">STR</span><span class="summary-val">{{ form.strength }} ({{ mod(form.strength) }})</span></div>
            <div class="summary-item"><span class="summary-label">DEX</span><span class="summary-val">{{ form.dexterity }} ({{ mod(form.dexterity) }})</span></div>
            <div class="summary-item"><span class="summary-label">CON</span><span class="summary-val">{{ form.constitution }} ({{ mod(form.constitution) }})</span></div>
            <div class="summary-item"><span class="summary-label">INT</span><span class="summary-val">{{ form.intelligence }} ({{ mod(form.intelligence) }})</span></div>
            <div class="summary-item"><span class="summary-label">WIS</span><span class="summary-val">{{ form.wisdom }} ({{ mod(form.wisdom) }})</span></div>
            <div class="summary-item"><span class="summary-label">CHA</span><span class="summary-val">{{ form.charisma }} ({{ mod(form.charisma) }})</span></div>
          </div>
          <p style="margin-top:10px;font-family:var(--font-display);font-size:10px;letter-spacing:.08em;text-transform:uppercase;color:var(--text-muted)">
            Base actions (Attack, Dash, Dodge, etc.) will be added automatically.
          </p>
        </div>

        <p v-if="error" class="create-error">{{ error }}</p>
      </div>

      <!-- Navigation -->
      <div class="step-nav">
        <button v-if="step > 1" class="btn btn-ghost" @click="step--">← Previous</button>
        <div style="flex:1" />
        <button
          v-if="step < TOTAL_STEPS"
          class="btn btn-primary"
          :disabled="!canProceed"
          @click="step++"
        >Next →</button>
        <button
          v-else
          class="btn btn-primary"
          :disabled="saving || !form.name || (!form.classSelect && !form.classCustom)"
          @click="submit"
        >{{ saving ? 'Creating...' : '⚔ Create Character' }}</button>
      </div>
    </main>
  </div>
</template>

<style scoped>
.create-page { min-height: 100vh; background: var(--bg-dark); }

.page-header { display: flex; align-items: center; justify-content: space-between; padding: 10px 20px; background: var(--bg-surface); border-bottom: 1px solid var(--border); position: sticky; top: 0; z-index: 20; }
.header-left { display: flex; align-items: center; gap: 12px; }
.header-icon { font-size: 18px; color: var(--accent-gold); }
.header-title { font-family: var(--font-deco); font-size: 15px; letter-spacing: .08em; }
.back-btn { padding: 5px 10px; font-size: 11px; }

.create-main { max-width: 760px; margin: 0 auto; padding: 40px 24px 80px; }

.create-hero { text-align: center; margin-bottom: 24px; }
.create-title {
  font-family: var(--font-display); font-size: 30px; font-weight: 600; letter-spacing: .06em;
  background: linear-gradient(180deg, #e8dcc8 0%, var(--accent-gold) 100%);
  -webkit-background-clip: text; -webkit-text-fill-color: transparent; background-clip: text;
  margin-bottom: 6px;
}
.create-sub { font-family: var(--font-body); font-style: italic; color: var(--text-secondary); font-size: 16px; }

.progress-bar { height: 3px; background: var(--bg-card); border-radius: 2px; margin-bottom: 20px; overflow: hidden; }
.progress-fill { height: 100%; background: var(--accent-gold); border-radius: 2px; transition: width .4s ease; }

.step-indicators { display: flex; justify-content: space-between; margin-bottom: 28px; }
.step-indicator { display: flex; flex-direction: column; align-items: center; gap: 5px; cursor: default; flex: 1; }
.step-indicator.done { cursor: pointer; }
.step-dot { width: 28px; height: 28px; border-radius: 50%; display: flex; align-items: center; justify-content: center; font-family: var(--font-display); font-size: 11px; font-weight: 600; border: 1px solid var(--border-light); background: var(--bg-card); color: var(--text-muted); transition: all .2s; }
.step-indicator.active .step-dot { background: var(--accent-gold); border-color: var(--accent-gold); color: #0a0c10; }
.step-indicator.done .step-dot { background: var(--accent-gold-dim); border-color: var(--accent-gold-dim); color: var(--accent-gold); }
.step-label { font-family: var(--font-display); font-size: 9px; letter-spacing: .08em; text-transform: uppercase; color: var(--text-muted); text-align: center; }
.step-indicator.active .step-label { color: var(--accent-gold); }

.step-content { min-height: 320px; }
.step-hint { font-family: var(--font-body); font-style: italic; font-size: 14px; color: var(--text-muted); margin-bottom: 16px; }

.npc-toggle-row { margin-bottom: 20px; }
.toggle-label { display: flex; align-items: center; gap: 10px; cursor: pointer; font-family: var(--font-body); font-size: 15px; color: var(--text-secondary); }
.toggle-track { width: 36px; height: 20px; border-radius: 10px; background: var(--border-light); position: relative; transition: background .2s; flex-shrink: 0; }
.toggle-track.active { background: var(--accent-gold-dim); }
.toggle-thumb { width: 14px; height: 14px; border-radius: 50%; background: var(--text-muted); position: absolute; top: 3px; left: 3px; transition: transform .2s, background .2s; }
.toggle-track.active .toggle-thumb { transform: translateX(16px); background: var(--accent-gold); }

.form-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 14px; }
.field { display: flex; flex-direction: column; gap: 5px; }
.field-wide { grid-column: span 2; }
.field-sm { max-width: 140px; }
.field-label { font-family: var(--font-display); font-size: 10px; font-weight: 600; letter-spacing: .1em; text-transform: uppercase; color: var(--text-muted); }

.ability-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 14px; }
.ability-card { background: var(--bg-card); border: 1px solid var(--border); border-radius: var(--radius-md); display: flex; flex-direction: column; align-items: center; gap: 4px; padding: 16px 12px; transition: border-color .2s; }
.ability-card:focus-within { border-color: var(--accent-gold-dim); }
.ability-abbr { font-family: var(--font-display); font-size: 14px; font-weight: 600; letter-spacing: .1em; color: var(--text-muted); }
.ability-full { font-family: var(--font-body); font-size: 12px; color: var(--text-muted); }
.ability-input { width: 70px; padding: 8px; text-align: center; background: var(--bg-input); border: 1px solid var(--border-light); border-radius: var(--radius-sm); color: var(--text-primary); font-family: var(--font-display); font-size: 24px; font-weight: 600; outline: none; transition: border-color .2s; }
.ability-input:focus { border-color: var(--accent-gold-dim); }
.ability-mod { font-family: var(--font-display); font-size: 18px; font-weight: 700; color: var(--accent-gold); }

.two-col { display: grid; grid-template-columns: 1fr 1fr; gap: 20px; }
.sub-panel { background: var(--bg-card); border: 1px solid var(--border); border-radius: var(--radius-sm); padding: 14px; }
.sub-panel-title { font-family: var(--font-display); font-size: 11px; font-weight: 600; letter-spacing: .12em; text-transform: uppercase; color: var(--text-muted); margin-bottom: 10px; }

.save-list { display: flex; flex-direction: column; gap: 6px; }
.save-check-row { display: flex; align-items: center; gap: 8px; cursor: pointer; padding: 4px; border-radius: 2px; transition: background .1s; }
.save-check-row:hover { background: var(--bg-surface); }
.prof-checkbox { accent-color: var(--accent-gold); width: 13px; height: 13px; cursor: pointer; }
.save-check-label { font-family: var(--font-body); font-size: 14px; color: var(--text-secondary); }

.skill-list { display: flex; flex-direction: column; gap: 2px; }
.skill-create-row { display: flex; align-items: center; gap: 6px; padding: 3px 2px; border-radius: 2px; transition: background .1s; }
.skill-create-row:hover { background: var(--bg-surface); }
.skill-prof-btn { width: 16px; height: 16px; border-radius: 50%; flex-shrink: 0; display: flex; align-items: center; justify-content: center; font-size: 9px; cursor: pointer; border: 1px solid var(--border-light); color: var(--text-muted); transition: all .15s; user-select: none; }
.skill-prof-btn.level-1 { background: var(--accent-gold-dim); border-color: var(--accent-gold); color: var(--accent-gold); }
.skill-prof-btn.level-2 { background: var(--accent-gold); border-color: var(--accent-gold); color: #000; }
.skill-create-name { font-family: var(--font-body); font-size: 13px; color: var(--text-secondary); flex: 1; }
.skill-create-ability { font-family: var(--font-display); font-size: 9px; color: var(--text-muted); letter-spacing: .06em; }

.combat-preview { display: flex; gap: 0; margin-top: 20px; background: var(--bg-card); border: 1px solid var(--border); border-radius: var(--radius-sm); overflow: hidden; }
.preview-stat { flex: 1; text-align: center; padding: 12px; border-right: 1px solid var(--border); }
.preview-stat:last-child { border-right: none; }
.preview-val { display: block; font-family: var(--font-display); font-size: 22px; font-weight: 700; color: var(--accent-gold); line-height: 1; }
.preview-label { display: block; font-family: var(--font-display); font-size: 9px; letter-spacing: .1em; text-transform: uppercase; color: var(--text-muted); margin-top: 3px; }

.training-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 14px; margin-top: 14px; }
.training-group { background: var(--bg-card); border: 1px solid var(--border); border-radius: var(--radius-sm); padding: 12px; }
.check-option { display: flex; align-items: center; gap: 8px; cursor: pointer; font-family: var(--font-body); font-size: 13px; color: var(--text-secondary); padding: 3px 0; }
.check-option:hover { color: var(--text-primary); }

.currency-form { display: flex; gap: 10px; }
.currency-field { display: flex; flex-direction: column; align-items: center; gap: 4px; }
.currency-label { font-family: var(--font-display); font-size: 11px; font-weight: 600; letter-spacing: .08em; }
.currency-input { width: 70px; text-align: center; }
.copper   { color: #b87333; }
.silver   { color: #aaa9ad; }
.electrum { color: #7ecece; }
.gold     { color: var(--accent-gold); }
.platinum { color: #e5e4e2; }

.personality-form { display: flex; flex-direction: column; gap: 12px; }
.textarea { resize: vertical; font-family: var(--font-body); line-height: 1.5; }

.summary-card { background: var(--bg-card); border: 1px solid var(--border-gold); border-radius: var(--radius-md); padding: 16px; }
.summary-title { font-family: var(--font-display); font-size: 11px; font-weight: 600; letter-spacing: .12em; text-transform: uppercase; color: var(--text-muted); margin-bottom: 12px; }
.summary-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 8px; }
.summary-item { display: flex; flex-direction: column; gap: 2px; }
.summary-label { font-family: var(--font-display); font-size: 9px; letter-spacing: .1em; text-transform: uppercase; color: var(--text-muted); }
.summary-val { font-family: var(--font-display); font-size: 14px; font-weight: 600; color: var(--text-primary); }

.step-nav { display: flex; align-items: center; margin-top: 32px; padding-top: 20px; border-top: 1px solid var(--border); }

.create-error { color: var(--accent-red-bright); font-family: var(--font-body); font-size: 14px; margin-top: 12px; text-align: center; padding: 8px; background: rgba(139,26,26,.1); border: 1px solid var(--accent-red); border-radius: var(--radius-sm); }
</style>
