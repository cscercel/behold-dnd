<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { combatAPI } from '@/api/combat'
import { characterAPI } from '@/api/characters'
import type { CombatEncounter, CombatParticipant, Character } from '@/types'

const router = useRouter()
const auth   = useAuthStore()

const encounters        = ref<CombatEncounter[]>([])
const encountersLoading = ref(false)
const activeEncounter   = ref<CombatEncounter | null>(null)
const participants      = ref<CombatParticipant[]>([])
const participantsLoading = ref(false)

onMounted(async () => {
  await loadEncounters()
  const active = encounters.value.find(e => e.is_active)
  if (active) await openEncounter(active)
})

async function loadEncounters() {
  encountersLoading.value = true
  try { encounters.value = await combatAPI.list() }
  finally { encountersLoading.value = false }
}

async function openEncounter(enc: CombatEncounter) {
  activeEncounter.value = enc
  participantsLoading.value = true
  try { participants.value = await combatAPI.listParticipants(enc.id) }
  finally { participantsLoading.value = false }
}

function closeEncounter() {
  activeEncounter.value = null
  participants.value = []
}

const sortedParticipants = computed(() =>
  [...participants.value].sort((a, b) => b.initiative - a.initiative)
)
const activeParticipants   = computed(() => sortedParticipants.value.filter(p => p.is_active))
const inactiveParticipants = computed(() => sortedParticipants.value.filter(p => !p.is_active))

const showCreateForm   = ref(false)
const newEncounterName = ref('')

async function createEncounter() {
  if (!newEncounterName.value.trim()) return
  const enc = await combatAPI.create(newEncounterName.value.trim())
  encounters.value.unshift(enc)
  newEncounterName.value = ''
  showCreateForm.value = false
  await openEncounter(enc)
}

async function deleteEncounter(enc: CombatEncounter) {
  if (!confirm(`Delete "${enc.name}"? This cannot be undone.`)) return
  await combatAPI.delete(enc.id)
  encounters.value = encounters.value.filter(e => e.id !== enc.id)
  if (activeEncounter.value?.id === enc.id) closeEncounter()
}

async function startEncounter() {
  if (!activeEncounter.value) return
  activeEncounter.value = await combatAPI.start(activeEncounter.value.id)
  encounters.value = encounters.value.map(e =>
    e.id === activeEncounter.value!.id ? activeEncounter.value! : e
  )
}

async function endEncounter() {
  if (!activeEncounter.value) return
  if (!confirm('End this encounter?')) return
  activeEncounter.value = await combatAPI.end(activeEncounter.value.id)
  encounters.value = encounters.value.map(e =>
    e.id === activeEncounter.value!.id ? activeEncounter.value! : e
  )
}

async function nextRound() {
  if (!activeEncounter.value) return
  activeEncounter.value = await combatAPI.nextRound(activeEncounter.value.id)
}

const showAddParticipant = ref(false)
const allCharacters      = ref<Character[]>([])
const charsLoading       = ref(false)
const addMode            = ref<'character' | 'manual'>('character')
const selectedCharID     = ref('')
const participantCount   = ref(1)
const manualDraft        = ref({ name: '', initiative: 0, current_hp: 0, max_hp: 0, armor_class: 10, speed: 30 })
const addingParticipant  = ref(false)

async function openAddParticipant() {
  showAddParticipant.value = true
  if (!allCharacters.value.length) {
    charsLoading.value = true
    try { allCharacters.value = await charactersAPI.list() }
    finally { charsLoading.value = false }
  }
}

async function addParticipant() {
  if (!activeEncounter.value) return
  addingParticipant.value = true
  try {
    if (addMode.value === 'character' && selectedCharID.value) {
      const char  = allCharacters.value.find(c => c.id === selectedCharID.value)
      const count = participantCount.value || 1
      const added: CombatParticipant[] = []
      for (let i = 0; i < count; i++) {
        const initiative = Math.floor(Math.random() * 20) + 1
        const name = count > 1 ? `${char?.name} ${i + 1}` : char?.name ?? ''
        const p = await combatAPI.addParticipant(activeEncounter.value.id, {
          character_id: selectedCharID.value,
          name,
          initiative,
        })
        added.push(p)
      }
      participants.value.push(...added)
    } else {
      const p = await combatAPI.addParticipant(activeEncounter.value.id, { ...manualDraft.value })
      participants.value.push(p)
    }
    showAddParticipant.value = false
    selectedCharID.value = ''
    participantCount.value = 1
    manualDraft.value = { name: '', initiative: 0, current_hp: 0, max_hp: 0, armor_class: 10, speed: 30 }
  } finally {
    addingParticipant.value = false
  }
}

async function removeParticipant(p: CombatParticipant) {
  if (!activeEncounter.value) return
  if (!confirm(`Remove ${p.name} from the encounter?`)) return
  await combatAPI.removeParticipant(activeEncounter.value.id, p.id)
  participants.value = participants.value.filter(x => x.id !== p.id)
}

const selectedParticipant = ref<CombatParticipant | null>(null)
const pAction             = ref<'damage' | 'heal' | 'conditions' | 'initiative' | null>(null)
const pAmount             = ref(0)
const pInitiative         = ref(0)
const pConditions         = ref<string[]>([])
const pActionLoading      = ref(false)

const CONDITION_OPTIONS = [
  'Blinded', 'Charmed', 'Deafened', 'Exhaustion', 'Frightened',
  'Grappled', 'Incapacitated', 'Invisible', 'Paralyzed', 'Petrified',
  'Poisoned', 'Prone', 'Restrained', 'Stunned', 'Unconscious',
]

function openParticipantModal(p: CombatParticipant) {
  selectedParticipant.value = p
  pAction.value = null
  pAmount.value = 0
  pInitiative.value = p.initiative
  pConditions.value = [...(p.conditions ?? [])]
}

function closeParticipantModal() {
  selectedParticipant.value = null
  pAction.value = null
}

function toggleCondition(cond: string) {
  const idx = pConditions.value.indexOf(cond)
  if (idx === -1) pConditions.value.push(cond)
  else pConditions.value.splice(idx, 1)
}

async function applyParticipantAction() {
  if (!activeEncounter.value || !selectedParticipant.value) return
  pActionLoading.value = true
  const encId = activeEncounter.value.id
  const pId   = selectedParticipant.value.id
  try {
    let updated: CombatParticipant | null = null
    if (pAction.value === 'damage' && pAmount.value > 0)
      updated = await combatAPI.damageParticipant(encId, pId, pAmount.value)
    else if (pAction.value === 'heal' && pAmount.value > 0)
      updated = await combatAPI.healParticipant(encId, pId, pAmount.value)
    else if (pAction.value === 'initiative')
      updated = await combatAPI.updateInitiative(encId, pId, pInitiative.value)
    else if (pAction.value === 'conditions')
      updated = await combatAPI.updateConditions(encId, pId, pConditions.value)

    if (updated) {
      const idx = participants.value.findIndex(p => p.id === pId)
      if (idx !== -1) participants.value[idx] = updated
      selectedParticipant.value = updated
    }
    pAction.value = null
    pAmount.value = 0
  } finally {
    pActionLoading.value = false
  }
}

async function toggleConcentration(p: CombatParticipant) {
  if (!activeEncounter.value) return
  const updated = await combatAPI.toggleConcentration(activeEncounter.value.id, p.id)
  const idx = participants.value.findIndex(x => x.id === p.id)
  if (idx !== -1) participants.value[idx] = updated
  if (selectedParticipant.value?.id === p.id) selectedParticipant.value = updated
}

async function deactivateParticipant(p: CombatParticipant) {
  if (!activeEncounter.value) return
  const updated = await combatAPI.deactivate(activeEncounter.value.id, p.id)
  const idx = participants.value.findIndex(x => x.id === p.id)
  if (idx !== -1) participants.value[idx] = updated
  if (selectedParticipant.value?.id === p.id) selectedParticipant.value = updated
}

function hpColor(p: CombatParticipant): string {
  const pct = p.max_hp > 0 ? (p.current_hp / p.max_hp) * 100 : 0
  if (pct > 50) return 'var(--hp-high)'
  if (pct > 25) return 'var(--hp-mid)'
  return 'var(--hp-low)'
}

function hpPct(p: CombatParticipant): number {
  return p.max_hp > 0 ? Math.round((p.current_hp / p.max_hp) * 100) : 0
}

function formatDate(iso: string): string {
  return new Date(iso).toLocaleDateString(undefined, { month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' })
}
</script>

<template>
  <div class="combat-page">

    <!-- Header -->
    <header class="page-header">
      <div class="header-left">
        <button class="btn btn-ghost back-btn" @click="activeEncounter ? closeEncounter() : router.push({ name: 'dm-dashboard' })">
          {{ activeEncounter ? '← Encounters' : '← Dashboard' }}
        </button>
        <span class="header-icon">⚔</span>
        <span class="header-title">Beyond DnD</span>
        <span class="header-badge badge badge-dm">DM</span>
      </div>
      <button class="btn btn-ghost" @click="auth.logout()">Logout</button>
    </header>

    <!-- ── Encounters List ── -->
    <div v-if="!activeEncounter" class="encounters-view">
      <main class="page-content">
        <div class="hero fade-in">
          <h1 class="hero-title">Combat Encounters</h1>
          <p class="hero-sub">Manage your battles, track your foes</p>
        </div>

        <div class="ornament" />

        <div class="create-bar">
          <template v-if="showCreateForm">
            <input
              v-model="newEncounterName"
              class="input encounter-name-input"
              type="text"
              placeholder="Encounter name (e.g. Goblin Ambush)..."
              autofocus
              @keyup.enter="createEncounter"
              @keyup.escape="showCreateForm = false"
            />
            <button class="btn btn-primary" :disabled="!newEncounterName.trim()" @click="createEncounter">Create</button>
            <button class="btn btn-ghost" @click="showCreateForm = false">Cancel</button>
          </template>
          <button v-else class="btn btn-primary" @click="showCreateForm = true">+ New Encounter</button>
        </div>

        <div v-if="encountersLoading" class="encounters-grid">
          <div v-for="i in 3" :key="i" class="encounter-card-skeleton">
            <div class="skeleton" style="height:22px;width:60%;margin-bottom:8px" />
            <div class="skeleton" style="height:14px;width:40%" />
          </div>
        </div>

        <div v-else-if="encounters.length" class="encounters-grid fade-in">
          <div
            v-for="enc in encounters"
            :key="enc.id"
            class="encounter-card"
            :class="{ active: enc.is_active }"
            @click="openEncounter(enc)"
          >
            <div class="encounter-card-header">
              <span v-if="enc.is_active" class="active-badge">● Active</span>
              <button class="icon-btn danger enc-delete" @click.stop="deleteEncounter(enc)" title="Delete">✕</button>
            </div>
            <h2 class="encounter-name">{{ enc.name }}</h2>
            <div class="encounter-meta">
              <span>Round {{ enc.round }}</span>
              <span>{{ formatDate(enc.created_at) }}</span>
            </div>
            <span class="encounter-arrow">→</span>
          </div>
        </div>

        <div v-else class="empty-state fade-in">
          <div class="empty-icon">⚔️</div>
          <p class="empty-title">No encounters yet</p>
          <p class="empty-sub">Create your first encounter to begin tracking combat</p>
        </div>
      </main>
    </div>

    <!-- ── Active Encounter View ── -->
    <div v-else class="tracker-view">

      <div class="encounter-bar">
        <div class="encounter-bar-left">
          <h2 class="encounter-title">{{ activeEncounter.name }}</h2>
          <div class="encounter-info">
            <span class="round-badge">Round {{ activeEncounter.round }}</span>
            <span class="status-badge" :class="activeEncounter.is_active ? 'status-active' : 'status-idle'">
              {{ activeEncounter.is_active ? 'In Progress' : 'Not Started' }}
            </span>
          </div>
        </div>
        <div class="encounter-bar-right">
          <button v-if="!activeEncounter.is_active" class="btn btn-primary enc-btn" @click="startEncounter">▶ Start Encounter</button>
          <template v-else>
            <button class="btn btn-ghost enc-btn" @click="nextRound">Next Round →</button>
            <button class="btn btn-danger enc-btn" @click="endEncounter">■ End Encounter</button>
          </template>
          <button class="btn btn-ghost enc-btn" @click="openAddParticipant">+ Add Combatant</button>
        </div>
      </div>

      <div class="tracker-body">

        <div class="participants-section">
          <div class="section-label">Initiative Order</div>

          <div v-if="participantsLoading" class="participants-list">
            <div v-for="i in 4" :key="i" class="participant-skeleton">
              <div class="skeleton" style="height:20px;width:30%;margin-bottom:6px" />
              <div class="skeleton" style="height:8px" />
            </div>
          </div>

          <div v-else-if="activeParticipants.length" class="participants-list">
            <div
              v-for="p in activeParticipants"
              :key="p.id"
              class="participant-card"
              @click="openParticipantModal(p)"
            >
              <div class="initiative-badge">{{ p.initiative }}</div>
              <div class="participant-body">
                <div class="participant-top">
                  <span class="participant-name">{{ p.name }}</span>
                  <div class="participant-tags">
                    <span v-if="p.concentration" class="tag tag-conc">Conc.</span>
                    <span v-for="cond in p.conditions" :key="cond" class="tag tag-cond">{{ cond }}</span>
                  </div>
                  <div class="participant-hp">
                    <span class="hp-nums">{{ p.current_hp }}<span class="hp-sep">/</span>{{ p.max_hp }}</span>
                    <span class="hp-ac">AC {{ p.armor_class }}</span>
                  </div>
                </div>
                <div class="p-hp-bar">
                  <div class="p-hp-bar-fill" :style="{ width: hpPct(p) + '%', background: hpColor(p) }" />
                </div>
              </div>
              <button class="icon-btn danger p-remove" @click.stop="removeParticipant(p)" title="Remove">✕</button>
            </div>
          </div>

          <div v-else class="empty-participants">
            <p>No combatants yet.</p>
            <button class="btn btn-ghost" style="margin-top:8px" @click="openAddParticipant">+ Add Combatant</button>
          </div>
        </div>

        <div v-if="inactiveParticipants.length" class="participants-section" style="margin-top:20px">
          <div class="section-label section-label-dead">Downed / Removed</div>
          <div class="participants-list participants-list-dead">
            <div
              v-for="p in inactiveParticipants"
              :key="p.id"
              class="participant-card participant-card-dead"
              @click="openParticipantModal(p)"
            >
              <div class="initiative-badge initiative-badge-dead">{{ p.initiative }}</div>
              <div class="participant-body">
                <div class="participant-top">
                  <span class="participant-name" style="opacity:.5">{{ p.name }}</span>
                  <span class="tag tag-dead">Down</span>
                  <div class="participant-hp" style="opacity:.5">
                    <span class="hp-nums">{{ p.current_hp }}<span class="hp-sep">/</span>{{ p.max_hp }}</span>
                  </div>
                </div>
                <div class="p-hp-bar">
                  <div class="p-hp-bar-fill" style="width:0%;background:var(--hp-low)" />
                </div>
              </div>
              <button class="icon-btn danger p-remove" @click.stop="removeParticipant(p)" title="Remove">✕</button>
            </div>
          </div>
        </div>

      </div>
    </div>

    <!-- ── Add Participant Modal ── -->
    <div v-if="showAddParticipant" class="modal-overlay" @click.self="showAddParticipant = false">
      <div class="modal">
        <h3 class="modal-title">Add Combatant</h3>

        <div class="mode-toggle">
          <button class="mode-btn" :class="{ active: addMode === 'character' }" @click="addMode = 'character'">From Character Sheet</button>
          <button class="mode-btn" :class="{ active: addMode === 'manual' }" @click="addMode = 'manual'">Manual Entry</button>
        </div>

        <!-- Character picker -->
        <div v-if="addMode === 'character'" class="modal-form">
          <div v-if="charsLoading" class="skeleton" style="height:40px" />
          <template v-else>
            <label class="field-label">Character</label>

            <div class="char-picker">
              <div v-if="allCharacters.filter(c => !c.is_npc).length" class="char-group-label">Player Characters</div>
              <div
                v-for="c in allCharacters.filter(c => !c.is_npc)"
                :key="c.id"
                class="char-option"
                :class="{ selected: selectedCharID === c.id }"
                @click="selectedCharID = c.id"
              >
                <div class="char-option-portrait">{{ c.name[0] }}</div>
                <div class="char-option-info">
                  <span class="char-option-name">{{ c.name }}</span>
                  <span class="char-option-sub">{{ c.race }} {{ c.class }} · Level {{ c.level }}</span>
                </div>
                <span v-if="selectedCharID === c.id" class="char-option-check">✓</span>
              </div>

              <div v-if="allCharacters.filter(c => c.is_npc).length" class="char-group-label" style="margin-top:4px">NPCs & Monsters</div>
              <div
                v-for="c in allCharacters.filter(c => c.is_npc)"
                :key="c.id"
                class="char-option"
                :class="{ selected: selectedCharID === c.id }"
                @click="selectedCharID = c.id"
              >
                <div class="char-option-portrait npc">{{ c.name[0] }}</div>
                <div class="char-option-info">
                  <span class="char-option-name">{{ c.name }}</span>
                  <span class="char-option-sub">{{ c.race }} {{ c.class }}</span>
                </div>
                <span v-if="selectedCharID === c.id" class="char-option-check">✓</span>
              </div>

              <p v-if="!allCharacters.length" class="empty-chars">No characters found. Create one first.</p>
            </div>

            <label class="field-label" style="margin-top:12px">
              How many?
              <span class="field-hint">(multiple adds numbered copies, e.g. "Goblin 1", "Goblin 2")</span>
            </label>
            <input v-model.number="participantCount" class="input" type="number" min="1" max="20" style="width:80px" />
          </template>
        </div>

        <!-- Manual entry -->
        <div v-else class="modal-form">
          <label class="field-label">Name *</label>
          <input v-model="manualDraft.name" class="input" type="text" placeholder="Goblin Scout" />
          <div class="modal-row-4">
            <div><label class="field-label">Initiative</label><input v-model.number="manualDraft.initiative" class="input" type="number" /></div>
            <div><label class="field-label">Current HP</label><input v-model.number="manualDraft.current_hp" class="input" type="number" min="0" /></div>
            <div><label class="field-label">Max HP</label><input v-model.number="manualDraft.max_hp" class="input" type="number" min="0" /></div>
            <div><label class="field-label">AC</label><input v-model.number="manualDraft.armor_class" class="input" type="number" min="0" /></div>
          </div>
        </div>

        <div class="modal-footer">
          <button class="btn btn-ghost" @click="showAddParticipant = false">Cancel</button>
          <button
            class="btn btn-primary"
            :disabled="addingParticipant || (addMode === 'character' && !selectedCharID) || (addMode === 'manual' && !manualDraft.name)"
            @click="addParticipant"
          >{{ addingParticipant ? 'Adding...' : 'Add to Encounter' }}</button>
        </div>
      </div>
    </div>

    <!-- ── Participant Action Modal ── -->
    <div v-if="selectedParticipant" class="modal-overlay" @click.self="closeParticipantModal">
      <div class="modal modal-participant">
        <div class="p-modal-header">
          <div>
            <h3 class="modal-title">{{ selectedParticipant.name }}</h3>
            <div class="p-modal-stats">
              <span>Initiative: <strong>{{ selectedParticipant.initiative }}</strong></span>
              <span>AC: <strong>{{ selectedParticipant.armor_class }}</strong></span>
              <span>Speed: <strong>{{ selectedParticipant.speed }}ft</strong></span>
            </div>
          </div>
          <button class="btn btn-ghost btn-sm" @click="closeParticipantModal">✕</button>
        </div>

        <div class="p-hp-display">
          <div class="p-hp-numbers">
            <span class="p-hp-current">{{ selectedParticipant.current_hp }}</span>
            <span class="p-hp-sep">/</span>
            <span class="p-hp-max">{{ selectedParticipant.max_hp }}</span>
          </div>
          <div class="p-hp-bar" style="margin-top:6px">
            <div class="p-hp-bar-fill" :style="{ width: hpPct(selectedParticipant) + '%', background: hpColor(selectedParticipant) }" />
          </div>
        </div>

        <div v-if="selectedParticipant.conditions?.length || selectedParticipant.concentration" class="p-conditions">
          <span v-if="selectedParticipant.concentration" class="tag tag-conc">Concentrating</span>
          <span v-for="cond in selectedParticipant.conditions" :key="cond" class="tag tag-cond">{{ cond }}</span>
        </div>

        <div class="p-action-buttons">
          <button class="p-action-btn" :class="{ active: pAction === 'damage' }" @click="pAction = pAction === 'damage' ? null : 'damage'">⚔ Damage</button>
          <button class="p-action-btn p-action-heal" :class="{ active: pAction === 'heal' }" @click="pAction = pAction === 'heal' ? null : 'heal'">+ Heal</button>
          <button class="p-action-btn p-action-init" :class="{ active: pAction === 'initiative' }" @click="pAction = pAction === 'initiative' ? null : 'initiative'">🎲 Initiative</button>
          <button class="p-action-btn p-action-cond" :class="{ active: pAction === 'conditions' }" @click="pAction = pAction === 'conditions' ? null : 'conditions'">⚡ Conditions</button>
        </div>

        <div v-if="pAction === 'damage' || pAction === 'heal'" class="p-input-row">
          <input
            v-model.number="pAmount"
            class="input"
            type="number"
            min="1"
            :placeholder="pAction === 'damage' ? 'Damage amount...' : 'Heal amount...'"
            @keyup.enter="applyParticipantAction"
          />
          <button class="btn btn-primary" :disabled="pActionLoading || pAmount <= 0" @click="applyParticipantAction">Apply</button>
        </div>

        <div v-if="pAction === 'initiative'" class="p-input-row">
          <input v-model.number="pInitiative" class="input" type="number" placeholder="New initiative..." @keyup.enter="applyParticipantAction" />
          <button class="btn btn-primary" :disabled="pActionLoading" @click="applyParticipantAction">Set</button>
        </div>

        <div v-if="pAction === 'conditions'" class="conditions-picker">
          <div class="conditions-grid">
            <label
              v-for="cond in CONDITION_OPTIONS"
              :key="cond"
              class="condition-check"
              :class="{ selected: pConditions.includes(cond) }"
              @click="toggleCondition(cond)"
            >{{ cond }}</label>
          </div>
          <button class="btn btn-primary btn-sm" style="margin-top:10px" :disabled="pActionLoading" @click="applyParticipantAction">Apply Conditions</button>
        </div>

        <div class="p-quick-actions">
          <button
            class="btn btn-ghost btn-sm"
            :class="{ 'conc-active': selectedParticipant.concentration }"
            @click="toggleConcentration(selectedParticipant)"
          >{{ selectedParticipant.concentration ? '◆ Drop Concentration' : '◇ Concentrating' }}</button>
          <button v-if="selectedParticipant.is_active" class="btn btn-danger btn-sm" @click="deactivateParticipant(selectedParticipant)">☠ Down</button>
          <button class="btn btn-ghost btn-sm" @click="removeParticipant(selectedParticipant); closeParticipantModal()">Remove</button>
        </div>
      </div>
    </div>

  </div>
</template>

<style scoped>
.combat-page { min-height: 100vh; background: var(--bg-dark); }

.page-header { display: flex; align-items: center; justify-content: space-between; padding: 10px 20px; background: var(--bg-surface); border-bottom: 1px solid var(--border); position: sticky; top: 0; z-index: 20; }
.header-left { display: flex; align-items: center; gap: 12px; }
.header-icon { font-size: 18px; color: var(--accent-gold); }
.header-title { font-family: var(--font-deco); font-size: 15px; letter-spacing: .08em; }
.header-badge { font-size: 10px; }
.back-btn { padding: 5px 10px; font-size: 11px; }

.encounters-view { min-height: calc(100vh - 53px); }
.page-content { max-width: 900px; margin: 0 auto; padding: 40px 24px; }

.hero { text-align: center; margin-bottom: 28px; }
.hero-title {
  font-family: var(--font-display); font-size: 32px; font-weight: 600; letter-spacing: .06em;
  background: linear-gradient(180deg, #e8dcc8 0%, var(--accent-gold) 100%);
  -webkit-background-clip: text; -webkit-text-fill-color: transparent; background-clip: text;
  margin-bottom: 6px;
}
.hero-sub { font-family: var(--font-body); font-style: italic; color: var(--text-secondary); font-size: 16px; }

.create-bar { display: flex; align-items: center; gap: 10px; margin-bottom: 20px; }
.encounter-name-input { flex: 1; max-width: 360px; }

.encounters-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(260px, 1fr)); gap: 16px; }
.encounter-card { background: var(--bg-card); border: 1px solid var(--border); border-radius: var(--radius-md); padding: 16px; cursor: pointer; position: relative; transition: border-color .2s, box-shadow .2s, transform .15s; display: flex; flex-direction: column; gap: 6px; }
.encounter-card:hover { border-color: var(--accent-gold-dim); box-shadow: 0 0 20px rgba(201,168,71,.1); transform: translateY(-2px); }
.encounter-card.active { border-color: var(--accent-red); box-shadow: 0 0 20px rgba(139,26,26,.15); }
.encounter-card-skeleton { background: var(--bg-card); border: 1px solid var(--border); border-radius: var(--radius-md); padding: 16px; }
.encounter-card-header { display: flex; justify-content: space-between; align-items: center; }
.active-badge { font-family: var(--font-display); font-size: 10px; letter-spacing: .08em; text-transform: uppercase; color: var(--accent-red-bright); }
.enc-delete { opacity: 0; transition: opacity .15s; }
.encounter-card:hover .enc-delete { opacity: 1; }
.encounter-name { font-family: var(--font-display); font-size: 17px; font-weight: 600; letter-spacing: .03em; }
.encounter-meta { display: flex; gap: 12px; font-family: var(--font-display); font-size: 10px; letter-spacing: .06em; text-transform: uppercase; color: var(--text-muted); }
.encounter-arrow { font-size: 16px; color: var(--accent-gold-dim); align-self: flex-end; transition: transform .2s, color .2s; }
.encounter-card:hover .encounter-arrow { transform: translateX(4px); color: var(--accent-gold); }

.empty-state { text-align: center; padding: 64px 32px; }
.empty-icon { font-size: 48px; margin-bottom: 16px; }
.empty-title { font-family: var(--font-display); font-size: 20px; color: var(--text-secondary); margin-bottom: 8px; }
.empty-sub { font-family: var(--font-body); font-style: italic; color: var(--text-muted); }

.tracker-view { display: flex; flex-direction: column; height: calc(100vh - 53px); }
.encounter-bar { display: flex; align-items: center; justify-content: space-between; padding: 12px 20px; background: var(--bg-surface); border-bottom: 1px solid var(--border); flex-wrap: wrap; gap: 10px; }
.encounter-bar-left { display: flex; align-items: center; gap: 14px; }
.encounter-title { font-family: var(--font-display); font-size: 20px; font-weight: 600; letter-spacing: .04em; }
.encounter-info { display: flex; align-items: center; gap: 8px; }
.round-badge { font-family: var(--font-display); font-size: 12px; letter-spacing: .08em; text-transform: uppercase; color: var(--accent-gold); background: rgba(201,168,71,.1); border: 1px solid var(--accent-gold-dim); border-radius: 3px; padding: 2px 8px; }
.status-badge { font-family: var(--font-display); font-size: 10px; letter-spacing: .08em; text-transform: uppercase; border-radius: 3px; padding: 2px 8px; }
.status-active { color: #4ade80; background: rgba(45,122,79,.15); border: 1px solid var(--hp-high); }
.status-idle   { color: var(--text-muted); background: var(--bg-card); border: 1px solid var(--border); }
.encounter-bar-right { display: flex; align-items: center; gap: 8px; flex-wrap: wrap; }
.enc-btn { padding: 7px 14px; font-size: 11px; }

.tracker-body { flex: 1; overflow-y: auto; padding: 16px 20px; max-width: 860px; width: 100%; margin: 0 auto; }

.section-label { font-family: var(--font-display); font-size: 10px; font-weight: 600; letter-spacing: .14em; text-transform: uppercase; color: var(--text-muted); margin-bottom: 10px; }
.section-label-dead { color: var(--accent-red); opacity: .6; }

.participants-list { display: flex; flex-direction: column; gap: 8px; }
.participants-list-dead { opacity: .7; }

.participant-card { display: flex; align-items: center; gap: 12px; background: var(--bg-card); border: 1px solid var(--border); border-radius: var(--radius-sm); padding: 10px 12px; cursor: pointer; position: relative; transition: border-color .15s, background .15s; }
.participant-card:hover { border-color: var(--accent-gold-dim); background: var(--bg-card-hover); }
.participant-skeleton { background: var(--bg-card); border: 1px solid var(--border); border-radius: var(--radius-sm); padding: 12px; }

.initiative-badge { width: 36px; height: 36px; border-radius: 50%; flex-shrink: 0; background: var(--bg-surface); border: 2px solid var(--accent-gold-dim); display: flex; align-items: center; justify-content: center; font-family: var(--font-display); font-size: 14px; font-weight: 700; color: var(--accent-gold); }
.initiative-badge-dead { border-color: var(--border); color: var(--text-muted); }

.participant-body { flex: 1; min-width: 0; }
.participant-top { display: flex; align-items: center; gap: 8px; margin-bottom: 6px; flex-wrap: wrap; }
.participant-name { font-family: var(--font-display); font-size: 15px; font-weight: 600; flex: 1; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }

.participant-tags { display: flex; gap: 4px; flex-wrap: wrap; }
.tag { padding: 1px 6px; border-radius: 2px; font-family: var(--font-display); font-size: 9px; letter-spacing: .06em; text-transform: uppercase; }
.tag-conc { background: rgba(26,58,92,.4); border: 1px solid #1a3a5c; color: #7ab3e0; }
.tag-cond { background: rgba(139,26,26,.2); border: 1px solid var(--accent-red); color: #f5a0a0; }
.tag-dead { background: rgba(60,60,60,.3); border: 1px solid var(--border-light); color: var(--text-muted); }

.participant-hp { display: flex; align-items: center; gap: 8px; flex-shrink: 0; }
.hp-nums { font-family: var(--font-display); font-size: 14px; font-weight: 600; }
.hp-sep  { color: var(--text-muted); }
.hp-ac   { font-family: var(--font-display); font-size: 11px; color: var(--text-muted); letter-spacing: .06em; }

.p-hp-bar { height: 4px; background: var(--bg-surface); border-radius: 2px; overflow: hidden; }
.p-hp-bar-fill { height: 100%; border-radius: 2px; transition: width .3s ease, background .3s ease; }

.p-remove { opacity: 0; transition: opacity .15s; flex-shrink: 0; }
.participant-card:hover .p-remove { opacity: 1; }

.empty-participants { text-align: center; padding: 40px; font-family: var(--font-body); font-style: italic; color: var(--text-muted); }

/* Add participant modal */
.mode-toggle { display: flex; gap: 0; background: var(--bg-surface); border: 1px solid var(--border); border-radius: var(--radius-sm); overflow: hidden; margin-bottom: 16px; }
.mode-btn { flex: 1; padding: 8px; font-family: var(--font-display); font-size: 11px; font-weight: 600; letter-spacing: .08em; text-transform: uppercase; border: none; background: transparent; color: var(--text-muted); cursor: pointer; transition: all .15s; }
.mode-btn.active { background: var(--accent-gold-dim); color: var(--accent-gold); }

.char-picker { max-height: 280px; overflow-y: auto; border: 1px solid var(--border); border-radius: var(--radius-sm); background: var(--bg-surface); }
.char-group-label { font-family: var(--font-display); font-size: 9px; font-weight: 600; letter-spacing: .12em; text-transform: uppercase; color: var(--text-muted); padding: 8px 10px 4px; border-bottom: 1px solid var(--border); }
.char-option { display: flex; align-items: center; gap: 10px; padding: 8px 10px; cursor: pointer; border-bottom: 1px solid var(--border); transition: background .1s; }
.char-option:last-child { border-bottom: none; }
.char-option:hover { background: var(--bg-card-hover); }
.char-option.selected { background: var(--accent-gold-glow); border-left: 2px solid var(--accent-gold); }
.char-option-portrait { width: 32px; height: 32px; border-radius: 50%; background: var(--bg-card); border: 1px solid var(--accent-gold-dim); display: flex; align-items: center; justify-content: center; font-family: var(--font-display); font-size: 14px; font-weight: 600; color: var(--accent-gold); flex-shrink: 0; }
.char-option-portrait.npc { border-color: #1a3a5c; color: #7ab3e0; }
.char-option-info { flex: 1; min-width: 0; }
.char-option-name { display: block; font-family: var(--font-display); font-size: 13px; font-weight: 600; }
.char-option-sub  { display: block; font-family: var(--font-body); font-size: 12px; color: var(--text-muted); margin-top: 1px; }
.char-option-check { font-size: 14px; color: var(--accent-gold); flex-shrink: 0; }
.empty-chars { font-family: var(--font-body); font-style: italic; color: var(--text-muted); font-size: 13px; padding: 16px; text-align: center; }

.modal-row-4 { display: grid; grid-template-columns: repeat(4, 1fr); gap: 8px; margin-top: 10px; }
.field-hint { font-family: var(--font-body); font-style: italic; font-size: 11px; color: var(--text-muted); font-weight: 400; letter-spacing: 0; text-transform: none; }

/* Participant action modal */
.modal-participant { max-width: 420px; }
.p-modal-header { display: flex; align-items: flex-start; justify-content: space-between; margin-bottom: 12px; }
.p-modal-stats { display: flex; gap: 12px; font-family: var(--font-display); font-size: 11px; color: var(--text-muted); letter-spacing: .06em; margin-top: 4px; }
.p-modal-stats strong { color: var(--text-primary); }

.p-hp-display { text-align: center; margin-bottom: 10px; }
.p-hp-numbers { display: flex; align-items: baseline; justify-content: center; gap: 4px; }
.p-hp-current { font-family: var(--font-display); font-size: 40px; font-weight: 700; line-height: 1; }
.p-hp-sep     { font-family: var(--font-display); font-size: 22px; color: var(--text-muted); }
.p-hp-max     { font-family: var(--font-display); font-size: 22px; color: var(--text-secondary); }

.p-conditions { display: flex; flex-wrap: wrap; gap: 4px; margin-bottom: 10px; }

.p-action-buttons { display: grid; grid-template-columns: 1fr 1fr 1fr 1fr; gap: 6px; margin-bottom: 10px; }
.p-action-btn { padding: 7px 4px; font-family: var(--font-display); font-size: 9px; font-weight: 600; letter-spacing: .06em; text-transform: uppercase; text-align: center; border: 1px solid var(--border-light); border-radius: var(--radius-sm); background: transparent; color: var(--text-muted); cursor: pointer; transition: all .15s; }
.p-action-btn:hover, .p-action-btn.active { border-color: var(--accent-red-bright); color: var(--accent-red-bright); background: rgba(139,26,26,.1); }
.p-action-heal:hover, .p-action-heal.active { border-color: var(--hp-high); color: #4ade80; background: rgba(45,122,79,.1); }
.p-action-init:hover, .p-action-init.active { border-color: var(--accent-gold-dim); color: var(--accent-gold); background: var(--accent-gold-glow); }
.p-action-cond:hover, .p-action-cond.active { border-color: #7ab3e0; color: #7ab3e0; background: rgba(26,58,92,.2); }

.p-input-row { display: flex; gap: 8px; margin-bottom: 10px; }
.p-input-row .input { flex: 1; }

.conditions-picker { margin-bottom: 10px; }
.conditions-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 5px; }
.condition-check { padding: 5px 8px; border-radius: var(--radius-sm); cursor: pointer; font-family: var(--font-display); font-size: 10px; letter-spacing: .06em; text-transform: uppercase; border: 1px solid var(--border); color: var(--text-muted); text-align: center; transition: all .15s; user-select: none; }
.condition-check:hover { border-color: var(--border-light); color: var(--text-secondary); }
.condition-check.selected { background: rgba(139,26,26,.2); border-color: var(--accent-red); color: #f5a0a0; }

.p-quick-actions { display: flex; gap: 6px; flex-wrap: wrap; padding-top: 10px; border-top: 1px solid var(--border); margin-top: 4px; }
.conc-active { border-color: #1a3a5c; color: #7ab3e0; }

/* Shared */
.modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,.78); display: flex; align-items: center; justify-content: center; z-index: 100; padding: 20px; }
.modal { background: var(--bg-card); border: 1px solid var(--border-light); border-radius: var(--radius-md); padding: 22px; width: 100%; max-width: 520px; max-height: 90vh; overflow-y: auto; }
.modal-title { font-family: var(--font-display); font-size: 16px; font-weight: 600; letter-spacing: .05em; margin-bottom: 14px; }
.modal-form { display: flex; flex-direction: column; gap: 8px; margin-bottom: 16px; }
.modal-footer { display: flex; justify-content: flex-end; gap: 7px; margin-top: 16px; }
.field-label { font-family: var(--font-display); font-size: 10px; font-weight: 600; letter-spacing: .1em; text-transform: uppercase; color: var(--text-muted); display: block; margin-bottom: 3px; }
.icon-btn { width: 20px; height: 20px; display: flex; align-items: center; justify-content: center; background: none; border: 1px solid var(--border); border-radius: 2px; color: var(--text-muted); font-size: 10px; cursor: pointer; transition: all .15s; }
.icon-btn.danger:hover { border-color: var(--accent-red); color: var(--accent-red-bright); }
.btn-sm { padding: 5px 10px; font-size: 11px; }
</style>
