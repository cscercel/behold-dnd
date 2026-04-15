<script setup lang="ts">
import { onMounted, computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useCharactersStore } from '@/stores/characters'
import { useCharacter } from '@/composables/useCharacter'
import { useAuthStore } from '@/stores/auth'

const route  = useRoute()
const router = useRouter()
const store  = useCharactersStore()
const auth   = useAuthStore()
const id     = route.params.id as string

onMounted(() => store.fetchOne(id))

const derived = computed(() =>
  store.current ? useCharacter(store.current) : null
)

// HP action panel
const hpAction = ref<'damage' | 'heal' | 'temp' | null>(null)
const hpAmount = ref(0)
const actionLoading = ref(false)

async function submitHPAction() {
  if (!hpAmount.value || hpAmount.value <= 0) return
  actionLoading.value = true
  try {
    if (hpAction.value === 'damage') await store.applyDamage(id, hpAmount.value)
    if (hpAction.value === 'heal')   await store.heal(id, hpAmount.value)
    if (hpAction.value === 'temp')   await store.addTempHP(id, hpAmount.value)
    hpAction.value = null
    hpAmount.value = 0
  } finally {
    actionLoading.value = false
  }
}

async function handleLongRest() {
  if (confirm('Take a long rest? This will restore HP, hit dice, spell slots, and clear conditions.')) {
    await store.longRest(id)
  }
}

function hpColor(pct: number): string {
  if (pct > 50) return 'var(--hp-high)'
  if (pct > 25) return 'var(--hp-mid)'
  return 'var(--hp-low)'
}

const ABILITIES = [
  { key: 'strength',     label: 'STR' },
  { key: 'dexterity',    label: 'DEX' },
  { key: 'constitution', label: 'CON' },
  { key: 'intelligence', label: 'INT' },
  { key: 'wisdom',       label: 'WIS' },
  { key: 'charisma',     label: 'CHA' },
] as const

const SAVES = [
  { key: 'strength',     label: 'Strength',     prof: 'save_prof_strength'     },
  { key: 'dexterity',    label: 'Dexterity',    prof: 'save_prof_dexterity'    },
  { key: 'constitution', label: 'Constitution', prof: 'save_prof_constitution' },
  { key: 'intelligence', label: 'Intelligence', prof: 'save_prof_intelligence' },
  { key: 'wisdom',       label: 'Wisdom',       prof: 'save_prof_wisdom'       },
  { key: 'charisma',     label: 'Charisma',     prof: 'save_prof_charisma'     },
] as const
</script>

<template>
  <div class="sheet-page">
    <!-- Header -->
    <header class="page-header">
      <div class="header-left">
        <button class="btn btn-ghost back-btn" @click="router.back()">← Back</button>
        <div class="header-brand">
          <span class="header-icon">⚔</span>
          <span class="header-title">Behold DnD</span>
        </div>
      </div>
      <div class="header-actions">
        <button class="btn btn-ghost" @click="auth.logout()">Logout</button>
      </div>
    </header>

    <div v-if="store.loading" class="loading-state">
      <div class="skeleton" style="height: 40px; width: 30%; margin-bottom: 8px;" />
      <div class="skeleton" style="height: 20px; width: 20%;" />
    </div>

    <p v-else-if="store.error" class="page-error">{{ store.error }}</p>

    <template v-else-if="store.current && derived">
      <!-- Character identity bar -->
      <div class="identity-bar">
        <div class="identity-portrait">{{ store.current.name[0] }}</div>
        <div class="identity-info">
          <h1 class="identity-name">{{ store.current.name }}</h1>
          <p class="identity-sub">
            {{ store.current.race }} · {{ store.current.class }} · {{ store.current.background }}
            <span v-if="store.current.alignment"> · {{ store.current.alignment }}</span>
          </p>
        </div>
        <div class="identity-level">
          <span class="level-num">{{ store.current.level }}</span>
          <span class="level-label">Level</span>
        </div>
        <div class="identity-xp">
          <span class="xp-val">{{ store.current.xp }}</span>
          <span class="xp-label">XP</span>
        </div>
      </div>

      <!-- Quick stats bar -->
      <div class="quickstats-bar">
        <div class="quickstat">
          <span class="quickstat-val">{{ store.current.armor_class }}</span>
          <span class="quickstat-label">Armor Class</span>
        </div>
        <div class="quickstat-divider" />
        <div class="quickstat">
          <span class="quickstat-val">{{ derived.initiative.value }}</span>
          <span class="quickstat-label">Initiative</span>
        </div>
        <div class="quickstat-divider" />
        <div class="quickstat">
          <span class="quickstat-val">{{ store.current.speed }} ft</span>
          <span class="quickstat-label">Speed</span>
        </div>
        <div class="quickstat-divider" />
        <div class="quickstat">
          <span class="quickstat-val">{{ derived.signedModifier(derived.proficiencyBonus.value) }}</span>
          <span class="quickstat-label">Proficiency</span>
        </div>
        <div class="quickstat-divider" />
        <div class="quickstat">
          <span class="quickstat-val">{{ derived.passivePerception.value }}</span>
          <span class="quickstat-label">Passive Perc.</span>
        </div>
        <div class="quickstat-divider" />
        <div class="quickstat">
          <span class="quickstat-val">{{ store.current.hit_dice_remaining }}d{{ store.current.hit_dice_type }}</span>
          <span class="quickstat-label">Hit Dice</span>
        </div>
      </div>

      <!-- Main sheet content -->
      <div class="sheet-body">

        <!-- LEFT: Ability scores + saves -->
        <div class="sheet-col sheet-col-left">

          <!-- Ability Scores -->
          <div class="sheet-panel">
            <div class="section-header">
              <span class="section-title">Ability Scores</span>
            </div>
            <div class="ability-grid">
              <div
                v-for="ab in ABILITIES"
                :key="ab.key"
                class="ability-block"
              >
                <span class="ability-label">{{ ab.label }}</span>
                <span class="ability-score">{{ store.current[ab.key as keyof typeof store.current] }}</span>
                <span class="ability-mod">{{ derived.signedModifier(derived.modifiers.value[ab.key]) }}</span>
              </div>
            </div>
          </div>

          <!-- Saving Throws -->
          <div class="sheet-panel">
            <div class="section-header">
              <span class="section-title">Saving Throws</span>
            </div>
            <div class="save-list">
              <div
                v-for="save in SAVES"
                :key="save.key"
                class="save-row"
              >
                <span
                  class="prof-dot"
                  :class="{ 'prof-dot-active': store.current[save.prof as keyof typeof store.current] }"
                />
                <span class="save-val">
                  {{ derived.signedModifier(derived.savingThrows.value[save.key]) }}
                </span>
                <span class="save-label">{{ save.label }}</span>
              </div>
            </div>
          </div>

          <!-- Conditions -->
          <div v-if="store.current.conditions?.length" class="sheet-panel">
            <div class="section-header">
              <span class="section-title">Conditions</span>
            </div>
            <div class="condition-chips">
              <span
                v-for="cond in store.current.conditions"
                :key="cond"
                class="condition-tag"
              >{{ cond }}</span>
            </div>
          </div>
        </div>

        <!-- MIDDLE: HP + Skills -->
        <div class="sheet-col sheet-col-mid">

          <!-- HP Panel -->
          <div class="sheet-panel hp-panel">
            <div class="section-header">
              <span class="section-title">Hit Points</span>
              <div class="hp-actions-row">
                <button
                  class="hp-action-btn"
                  :class="{ active: hpAction === 'damage' }"
                  @click="hpAction = hpAction === 'damage' ? null : 'damage'"
                >Damage</button>
                <button
                  class="hp-action-btn hp-action-heal"
                  :class="{ active: hpAction === 'heal' }"
                  @click="hpAction = hpAction === 'heal' ? null : 'heal'"
                >Heal</button>
                <button
                  class="hp-action-btn hp-action-temp"
                  :class="{ active: hpAction === 'temp' }"
                  @click="hpAction = hpAction === 'temp' ? null : 'temp'"
                >Temp</button>
              </div>
            </div>

            <!-- HP Display -->
            <div class="hp-display">
              <div class="hp-numbers">
                <span class="hp-current">{{ store.current.current_hp }}</span>
                <span class="hp-sep">/</span>
                <span class="hp-max">{{ store.current.max_hp }}</span>
                <span v-if="store.current.temp_hp > 0" class="hp-temp">
                  +{{ store.current.temp_hp }} temp
                </span>
              </div>
              <div class="hp-bar" style="margin-top: 8px;">
                <div
                  class="hp-bar-fill"
                  :style="{
                    width: derived.hpPercentage.value + '%',
                    background: hpColor(derived.hpPercentage.value)
                  }"
                />
              </div>
            </div>

            <!-- HP action input -->
            <div v-if="hpAction" class="hp-input-row">
              <input
                v-model.number="hpAmount"
                class="input hp-input"
                type="number"
                min="1"
                placeholder="Amount..."
                @keyup.enter="submitHPAction"
              />
              <button
                class="btn btn-primary"
                :disabled="actionLoading || hpAmount <= 0"
                @click="submitHPAction"
              >Apply</button>
              <button class="btn btn-ghost" @click="hpAction = null; hpAmount = 0">✕</button>
            </div>
          </div>

          <!-- Death saves -->
          <div
            v-if="store.current.current_hp === 0"
            class="sheet-panel death-saves-panel"
          >
            <div class="section-header">
              <span class="section-title">Death Saving Throws</span>
            </div>
            <div class="death-saves">
              <div class="death-row">
                <span class="death-label">Successes</span>
                <div class="death-dots">
                  <span
                    v-for="i in 3"
                    :key="i"
                    class="death-dot death-dot-success"
                    :class="{ filled: i <= store.current.death_save_successes }"
                  />
                </div>
              </div>
              <div class="death-row">
                <span class="death-label">Failures</span>
                <div class="death-dots">
                  <span
                    v-for="i in 3"
                    :key="i"
                    class="death-dot death-dot-failure"
                    :class="{ filled: i <= store.current.death_save_failures }"
                  />
                </div>
              </div>
            </div>
          </div>

          <!-- Skills -->
          <div class="sheet-panel">
            <div class="section-header">
              <span class="section-title">Skills</span>
            </div>
            <div class="skill-list">
              <div
                v-for="skill in derived.skills.value"
                :key="skill.name"
                class="skill-row"
              >
                <span
                  class="prof-dot"
                  :class="{
                    'prof-dot-active': skill.profLevel >= 1,
                    'prof-dot-expert': skill.profLevel === 2
                  }"
                />
                <span class="skill-val">{{ skill.display }}</span>
                <span class="skill-name">{{ skill.name }}</span>
                <span class="skill-ability">{{ skill.ability.slice(0, 3).toUpperCase() }}</span>
              </div>
            </div>
          </div>

          <!-- Rest buttons -->
          <div class="rest-row">
            <button class="btn btn-ghost rest-btn" @click="handleLongRest">🌙 Long Rest</button>
          </div>
        </div>

        <!-- RIGHT: Combat + Traits -->
        <div class="sheet-col sheet-col-right">

          <!-- Currency -->
          <div class="sheet-panel">
            <div class="section-header">
              <span class="section-title">Currency</span>
            </div>
            <div class="currency-grid">
              <div class="currency-item">
                <span class="currency-val">{{ store.current.copper }}</span>
                <span class="currency-label copper">CP</span>
              </div>
              <div class="currency-item">
                <span class="currency-val">{{ store.current.silver }}</span>
                <span class="currency-label silver">SP</span>
              </div>
              <div class="currency-item">
                <span class="currency-val">{{ store.current.electrum }}</span>
                <span class="currency-label electrum">EP</span>
              </div>
              <div class="currency-item">
                <span class="currency-val">{{ store.current.gold }}</span>
                <span class="currency-label gold">GP</span>
              </div>
              <div class="currency-item">
                <span class="currency-val">{{ store.current.platinum }}</span>
                <span class="currency-label platinum">PP</span>
              </div>
            </div>
          </div>

          <!-- Defenses -->
          <div
            v-if="store.current.resistances?.length || store.current.immunities?.length || store.current.vulnerabilities?.length"
            class="sheet-panel"
          >
            <div class="section-header">
              <span class="section-title">Defenses</span>
            </div>
            <div v-if="store.current.resistances?.length" class="defense-group">
              <span class="defense-type">Resistances</span>
              <div class="defense-tags">
                <span
                  v-for="r in store.current.resistances"
                  :key="r"
                  class="defense-tag resist"
                >{{ r }}</span>
              </div>
            </div>
            <div v-if="store.current.immunities?.length" class="defense-group">
              <span class="defense-type">Immunities</span>
              <div class="defense-tags">
                <span
                  v-for="im in store.current.immunities"
                  :key="im"
                  class="defense-tag immune"
                >{{ im }}</span>
              </div>
            </div>
            <div v-if="store.current.vulnerabilities?.length" class="defense-group">
              <span class="defense-type">Vulnerabilities</span>
              <div class="defense-tags">
                <span
                  v-for="v in store.current.vulnerabilities"
                  :key="v"
                  class="defense-tag vuln"
                >{{ v }}</span>
              </div>
            </div>
          </div>

          <!-- Proficiencies & Languages -->
          <div class="sheet-panel">
            <div class="section-header">
              <span class="section-title">Proficiencies</span>
            </div>
            <div v-if="store.current.training_armor?.length" class="proficiency-group">
              <span class="prof-type">Armor</span>
              <p class="prof-vals">{{ store.current.training_armor.join(', ') }}</p>
            </div>
            <div v-if="store.current.training_weapons?.length" class="proficiency-group">
              <span class="prof-type">Weapons</span>
              <p class="prof-vals">{{ store.current.training_weapons.join(', ') }}</p>
            </div>
            <div v-if="store.current.training_tools?.length" class="proficiency-group">
              <span class="prof-type">Tools</span>
              <p class="prof-vals">{{ store.current.training_tools.join(', ') }}</p>
            </div>
            <div v-if="store.current.training_languages?.length" class="proficiency-group">
              <span class="prof-type">Languages</span>
              <p class="prof-vals">{{ store.current.training_languages.join(', ') }}</p>
            </div>
          </div>

          <!-- Personality -->
          <div class="sheet-panel">
            <div class="section-header">
              <span class="section-title">Personality</span>
            </div>
            <div v-if="store.current.personality_traits" class="trait-group">
              <span class="trait-type">Traits</span>
              <p class="trait-text">{{ store.current.personality_traits }}</p>
            </div>
            <div v-if="store.current.ideals" class="trait-group">
              <span class="trait-type">Ideals</span>
              <p class="trait-text">{{ store.current.ideals }}</p>
            </div>
            <div v-if="store.current.bonds" class="trait-group">
              <span class="trait-type">Bonds</span>
              <p class="trait-text">{{ store.current.bonds }}</p>
            </div>
            <div v-if="store.current.flaws" class="trait-group">
              <span class="trait-type">Flaws</span>
              <p class="trait-text">{{ store.current.flaws }}</p>
            </div>
          </div>

          <!-- Notes -->
          <div v-if="store.current.notes" class="sheet-panel">
            <div class="section-header">
              <span class="section-title">Notes</span>
            </div>
            <p class="notes-text">{{ store.current.notes }}</p>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<style scoped>
.sheet-page {
  min-height: 100vh;
  background: var(--bg-dark);
}

/* Header */
.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 24px;
  background: var(--bg-surface);
  border-bottom: 1px solid var(--border);
  position: sticky;
  top: 0;
  z-index: 10;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.header-brand {
  display: flex;
  align-items: center;
  gap: 8px;
}

.header-icon {
  font-size: 18px;
  color: var(--accent-gold);
}

.header-title {
  font-family: var(--font-deco);
  font-size: 15px;
  letter-spacing: 0.08em;
}

.back-btn {
  padding: 6px 12px;
  font-size: 12px;
}

.header-actions { display: flex; gap: 12px; }

/* Identity bar */
.identity-bar {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px 24px;
  background: var(--bg-surface);
  border-bottom: 1px solid var(--border);
}

.identity-portrait {
  width: 56px;
  height: 56px;
  border-radius: 50%;
  background: var(--bg-card);
  border: 2px solid var(--accent-gold-dim);
  display: flex;
  align-items: center;
  justify-content: center;
  font-family: var(--font-display);
  font-size: 24px;
  font-weight: 600;
  color: var(--accent-gold);
  flex-shrink: 0;
}

.identity-info { flex: 1; }

.identity-name {
  font-family: var(--font-display);
  font-size: 24px;
  font-weight: 600;
  letter-spacing: 0.04em;
}

.identity-sub {
  font-family: var(--font-body);
  font-size: 15px;
  color: var(--text-secondary);
  margin-top: 2px;
}

.identity-level,
.identity-xp {
  text-align: center;
  flex-shrink: 0;
  padding: 0 16px;
}

.level-num {
  display: block;
  font-family: var(--font-display);
  font-size: 28px;
  font-weight: 700;
  color: var(--accent-gold);
  line-height: 1;
}

.level-label,
.xp-label {
  display: block;
  font-family: var(--font-display);
  font-size: 9px;
  letter-spacing: 0.12em;
  text-transform: uppercase;
  color: var(--text-muted);
  margin-top: 2px;
}

.xp-val {
  display: block;
  font-family: var(--font-display);
  font-size: 18px;
  font-weight: 600;
  color: var(--text-secondary);
  line-height: 1;
}

/* Quickstats bar */
.quickstats-bar {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 12px 24px;
  background: var(--bg-card);
  border-bottom: 1px solid var(--border);
  gap: 0;
  overflow-x: auto;
}

.quickstat {
  text-align: center;
  padding: 0 24px;
}

.quickstat-val {
  display: block;
  font-family: var(--font-display);
  font-size: 20px;
  font-weight: 600;
  color: var(--text-primary);
  line-height: 1;
}

.quickstat-label {
  display: block;
  font-family: var(--font-display);
  font-size: 9px;
  letter-spacing: 0.12em;
  text-transform: uppercase;
  color: var(--text-muted);
  margin-top: 3px;
}

.quickstat-divider {
  width: 1px;
  height: 32px;
  background: var(--border);
  flex-shrink: 0;
}

/* Sheet body */
.sheet-body {
  display: grid;
  grid-template-columns: 200px 1fr 240px;
  gap: 0;
  min-height: calc(100vh - 200px);
}

.sheet-col {
  padding: 16px;
  border-right: 1px solid var(--border);
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.sheet-col-right { border-right: none; }

/* Panel */
.sheet-panel {
  background: var(--bg-card);
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  padding: 14px;
}

/* Ability scores */
.ability-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 8px;
}

.ability-block {
  background: var(--bg-surface);
  border: 1px solid var(--border-light);
  border-radius: var(--radius-sm);
  text-align: center;
  padding: 10px 6px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
}

.ability-label {
  font-family: var(--font-display);
  font-size: 9px;
  letter-spacing: 0.1em;
  text-transform: uppercase;
  color: var(--text-muted);
}

.ability-score {
  font-family: var(--font-display);
  font-size: 20px;
  font-weight: 600;
  color: var(--text-primary);
  line-height: 1;
}

.ability-mod {
  font-family: var(--font-display);
  font-size: 13px;
  color: var(--accent-gold);
}

/* Saving throws */
.save-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.save-row {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 3px 0;
}

.prof-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  border: 1px solid var(--border-light);
  background: transparent;
  flex-shrink: 0;
}

.prof-dot-active {
  background: var(--accent-gold);
  border-color: var(--accent-gold);
}

.prof-dot-expert {
  background: var(--accent-gold);
  border-color: var(--accent-gold);
  box-shadow: 0 0 0 2px var(--accent-gold-dim);
}

.save-val {
  font-family: var(--font-display);
  font-size: 13px;
  font-weight: 600;
  color: var(--text-primary);
  width: 28px;
  text-align: right;
}

.save-label {
  font-family: var(--font-body);
  font-size: 14px;
  color: var(--text-secondary);
}

/* Conditions */
.condition-chips {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}

.condition-tag {
  padding: 2px 8px;
  background: rgba(139,26,26,0.2);
  border: 1px solid var(--accent-red);
  border-radius: 2px;
  font-family: var(--font-display);
  font-size: 10px;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: #f5a0a0;
}

/* HP */
.hp-display { margin-top: 10px; }

.hp-numbers {
  display: flex;
  align-items: baseline;
  gap: 4px;
}

.hp-current {
  font-family: var(--font-display);
  font-size: 36px;
  font-weight: 700;
  color: var(--text-primary);
  line-height: 1;
}

.hp-sep {
  font-family: var(--font-display);
  font-size: 20px;
  color: var(--text-muted);
}

.hp-max {
  font-family: var(--font-display);
  font-size: 20px;
  color: var(--text-secondary);
}

.hp-temp {
  font-family: var(--font-display);
  font-size: 13px;
  color: var(--accent-gold);
  margin-left: 8px;
}

.hp-actions-row {
  display: flex;
  gap: 4px;
}

.hp-action-btn {
  padding: 4px 10px;
  font-family: var(--font-display);
  font-size: 10px;
  font-weight: 600;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  border: 1px solid var(--border-light);
  border-radius: var(--radius-sm);
  background: transparent;
  color: var(--text-muted);
  cursor: pointer;
  transition: all 0.15s;
}

.hp-action-btn:hover,
.hp-action-btn.active {
  border-color: var(--accent-red);
  color: var(--accent-red-bright);
}

.hp-action-heal:hover,
.hp-action-heal.active {
  border-color: var(--hp-high);
  color: #4ade80;
}

.hp-action-temp:hover,
.hp-action-temp.active {
  border-color: var(--accent-gold-dim);
  color: var(--accent-gold);
}

.hp-input-row {
  display: flex;
  gap: 8px;
  margin-top: 10px;
}

.hp-input {
  flex: 1;
  padding: 8px 10px;
}

/* Death saves */
.death-saves-panel { background: rgba(139,26,26,0.08); border-color: var(--accent-red); }

.death-saves { display: flex; flex-direction: column; gap: 8px; margin-top: 10px; }

.death-row { display: flex; align-items: center; gap: 10px; }

.death-label {
  font-family: var(--font-display);
  font-size: 11px;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: var(--text-muted);
  width: 72px;
}

.death-dots { display: flex; gap: 6px; }

.death-dot {
  width: 14px;
  height: 14px;
  border-radius: 50%;
  border: 1px solid var(--border-light);
  background: transparent;
  transition: all 0.2s;
}

.death-dot-success.filled { background: var(--hp-high); border-color: var(--hp-high); }
.death-dot-failure.filled { background: var(--accent-red-bright); border-color: var(--accent-red-bright); }

/* Skills */
.skill-list { display: flex; flex-direction: column; gap: 2px; }

.skill-row {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 3px 4px;
  border-radius: 2px;
  transition: background 0.1s;
}

.skill-row:hover { background: var(--bg-surface); }

.skill-val {
  font-family: var(--font-display);
  font-size: 12px;
  font-weight: 600;
  color: var(--text-primary);
  width: 24px;
  text-align: right;
}

.skill-name {
  font-family: var(--font-body);
  font-size: 14px;
  color: var(--text-secondary);
  flex: 1;
}

.skill-ability {
  font-family: var(--font-display);
  font-size: 9px;
  letter-spacing: 0.08em;
  color: var(--text-muted);
}

/* Rest */
.rest-row {
  display: flex;
  gap: 8px;
  justify-content: center;
}

.rest-btn { font-size: 12px; padding: 8px 16px; }

/* Currency */
.currency-grid {
  display: flex;
  gap: 8px;
  margin-top: 8px;
}

.currency-item {
  flex: 1;
  text-align: center;
  background: var(--bg-surface);
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  padding: 8px 4px;
}

.currency-val {
  display: block;
  font-family: var(--font-display);
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
}

.currency-label {
  display: block;
  font-family: var(--font-display);
  font-size: 9px;
  letter-spacing: 0.1em;
  margin-top: 3px;
}

.copper  { color: #b87333; }
.silver  { color: #aaa9ad; }
.electrum{ color: #7ecece; }
.gold    { color: var(--accent-gold); }
.platinum{ color: #e5e4e2; }

/* Defenses */
.defense-group { margin-bottom: 10px; }

.defense-type {
  display: block;
  font-family: var(--font-display);
  font-size: 10px;
  letter-spacing: 0.1em;
  text-transform: uppercase;
  color: var(--text-muted);
  margin-bottom: 4px;
}

.defense-tags { display: flex; flex-wrap: wrap; gap: 4px; }

.defense-tag {
  padding: 2px 8px;
  border-radius: 2px;
  font-family: var(--font-display);
  font-size: 10px;
  letter-spacing: 0.06em;
  text-transform: uppercase;
}

.resist  { background: rgba(29,70,92,0.3); border: 1px solid #1a3a5c; color: #7ab3e0; }
.immune  { background: rgba(45,122,79,0.2); border: 1px solid var(--hp-high); color: #4ade80; }
.vuln    { background: rgba(139,26,26,0.2); border: 1px solid var(--accent-red); color: #f5a0a0; }

/* Proficiencies */
.proficiency-group { margin-bottom: 10px; }
.prof-type {
  display: block;
  font-family: var(--font-display);
  font-size: 10px;
  letter-spacing: 0.1em;
  text-transform: uppercase;
  color: var(--text-muted);
  margin-bottom: 3px;
}
.prof-vals {
  font-family: var(--font-body);
  font-size: 14px;
  color: var(--text-secondary);
}

/* Traits */
.trait-group { margin-bottom: 10px; }
.trait-type {
  display: block;
  font-family: var(--font-display);
  font-size: 10px;
  letter-spacing: 0.1em;
  text-transform: uppercase;
  color: var(--text-muted);
  margin-bottom: 3px;
}
.trait-text {
  font-family: var(--font-body);
  font-style: italic;
  font-size: 14px;
  color: var(--text-secondary);
  line-height: 1.5;
}

.notes-text {
  font-family: var(--font-body);
  font-size: 14px;
  color: var(--text-secondary);
  line-height: 1.6;
}

/* Loading */
.loading-state {
  padding: 48px 32px;
}

.page-error {
  color: var(--accent-red-bright);
  text-align: center;
  padding: 32px;
}
</style>
