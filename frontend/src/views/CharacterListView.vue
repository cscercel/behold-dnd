<script setup lang="ts">
import { onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useCharactersStore } from '@/stores/characters'
import { useAuthStore } from '@/stores/auth'
import type { Character } from '@/types'

const router = useRouter()
const store = useCharactersStore()
const auth = useAuthStore()

onMounted(() => store.fetchAll())

// Players only see their own characters
const myCharacters = computed(() =>
  store.characters.filter(
    (c) => !c.is_npc && c.owner_id === auth.user?.id
  )
)

function hpColor(char: Character): string {
  const pct = char.max_hp > 0 ? char.current_hp / char.max_hp : 0
  if (pct > 0.5) return 'var(--hp-high)'
  if (pct > 0.25) return 'var(--hp-mid)'
  return 'var(--hp-low)'
}

function hpPercent(char: Character): number {
  return char.max_hp > 0 ? Math.round((char.current_hp / char.max_hp) * 100) : 0
}
</script>

<template>
  <div class="characters-page">
    <!-- Header -->
    <header class="page-header">
      <div class="header-brand">
        <span class="header-icon">⚔</span>
        <span class="header-title">Behold DnD</span>
      </div>
      <div class="header-actions">
        <span class="header-user">{{ auth.user?.username }}</span>
        <button class="btn btn-ghost" @click="auth.logout()">Logout</button>
      </div>
    </header>

    <main class="page-content">
      <!-- Hero section -->
      <div class="hero fade-in">
        <h1 class="hero-title">Your Characters</h1>
        <p class="hero-sub">Choose your hero or forge a new legend</p>
      </div>

      <div class="ornament" />

      <!-- Loading -->
      <div v-if="store.loading" class="character-grid">
        <div v-for="i in 3" :key="i" class="character-card-skeleton">
          <div class="skeleton" style="height: 24px; width: 60%; margin-bottom: 8px;" />
          <div class="skeleton" style="height: 16px; width: 40%; margin-bottom: 16px;" />
          <div class="skeleton" style="height: 8px;" />
        </div>
      </div>

      <!-- Error -->
      <p v-else-if="store.error" class="page-error">{{ store.error }}</p>

      <!-- Empty -->
      <div v-else-if="myCharacters.length === 0" class="empty-state fade-in">
        <div class="empty-icon">📜</div>
        <p class="empty-title">No characters yet</p>
        <p class="empty-sub">Your story is waiting to be written</p>
      </div>

      <!-- Character grid -->
      <div v-else class="character-grid fade-in">
        <div
          v-for="char in myCharacters"
          :key="char.id"
          class="character-card"
          @click="router.push({ name: 'character-sheet', params: { id: char.id } })"
        >
          <!-- Card header -->
          <div class="char-header">
            <div class="char-portrait">{{ char.name[0] }}</div>
            <div class="char-identity">
              <h2 class="char-name">{{ char.name }}</h2>
              <p class="char-sub">{{ char.race }} {{ char.class }}</p>
            </div>
            <div class="char-level">
              <span class="level-num">{{ char.level }}</span>
              <span class="level-label">Level</span>
            </div>
          </div>

          <div class="ornament" style="margin: 12px 0;" />

          <!-- Stats row -->
          <div class="char-stats">
            <div class="char-stat">
              <span class="char-stat-val">{{ char.armor_class }}</span>
              <span class="char-stat-label">AC</span>
            </div>
            <div class="char-stat">
              <span class="char-stat-val">{{ char.speed }}</span>
              <span class="char-stat-label">Speed</span>
            </div>
            <div class="char-stat">
              <span class="char-stat-val">{{ char.current_hp }}/{{ char.max_hp }}</span>
              <span class="char-stat-label">HP</span>
            </div>
          </div>

          <!-- HP bar -->
          <div class="hp-bar" style="margin-top: 10px;">
            <div
              class="hp-bar-fill"
              :style="{ width: hpPercent(char) + '%', background: hpColor(char) }"
            />
          </div>

          <!-- Conditions -->
          <div v-if="char.conditions?.length" class="char-conditions">
            <span
              v-for="cond in char.conditions"
              :key="cond"
              class="condition-tag"
            >{{ cond }}</span>
          </div>

          <div class="char-footer">
            <span class="char-xp">{{ char.xp }} XP</span>
            <span class="char-arrow">→</span>
          </div>
        </div>
      </div>

      <!-- Create character -->
      <div class="create-section fade-in">
        <button
          class="btn btn-ghost create-btn"
          @click="router.push({ name: 'character-create' })"
        >
          + Create New Character
        </button>
      </div>
    </main>
  </div>
</template>

<style scoped>
.characters-page {
  min-height: 100vh;
  background: var(--bg-dark);
}

/* Header */
.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 32px;
  background: var(--bg-surface);
  border-bottom: 1px solid var(--border);
  position: sticky;
  top: 0;
  z-index: 10;
}

.header-brand {
  display: flex;
  align-items: center;
  gap: 10px;
}

.header-icon {
  font-size: 20px;
  color: var(--accent-gold);
  filter: drop-shadow(0 0 6px rgba(201,168,71,0.4));
}

.header-title {
  font-family: var(--font-deco);
  font-size: 16px;
  letter-spacing: 0.08em;
  color: var(--text-primary);
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 16px;
}

.header-user {
  font-family: var(--font-display);
  font-size: 13px;
  color: var(--text-secondary);
}

/* Hero */
.page-content {
  max-width: 960px;
  margin: 0 auto;
  padding: 48px 32px;
}

.hero { text-align: center; margin-bottom: 32px; }

.hero-title {
  font-family: var(--font-display);
  font-size: 36px;
  font-weight: 600;
  letter-spacing: 0.06em;
  background: linear-gradient(180deg, #e8dcc8 0%, var(--accent-gold) 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin-bottom: 8px;
}

.hero-sub {
  font-family: var(--font-body);
  font-style: italic;
  color: var(--text-secondary);
  font-size: 17px;
}

/* Grid */
.character-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 20px;
  margin-top: 24px;
}

.character-card {
  background: var(--bg-card);
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  padding: 20px;
  cursor: pointer;
  transition: border-color 0.2s, box-shadow 0.2s, transform 0.15s;
}

.character-card:hover {
  border-color: var(--accent-gold-dim);
  box-shadow: 0 0 24px rgba(201,168,71,0.12);
  transform: translateY(-2px);
}

.character-card-skeleton {
  background: var(--bg-card);
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  padding: 20px;
}

.char-header {
  display: flex;
  align-items: center;
  gap: 14px;
}

.char-portrait {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  background: var(--bg-surface);
  border: 2px solid var(--accent-gold-dim);
  display: flex;
  align-items: center;
  justify-content: center;
  font-family: var(--font-display);
  font-size: 20px;
  font-weight: 600;
  color: var(--accent-gold);
  flex-shrink: 0;
}

.char-identity { flex: 1; min-width: 0; }

.char-name {
  font-family: var(--font-display);
  font-size: 17px;
  font-weight: 600;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  letter-spacing: 0.03em;
}

.char-sub {
  font-family: var(--font-body);
  font-size: 14px;
  color: var(--text-secondary);
  margin-top: 2px;
}

.char-level {
  display: flex;
  flex-direction: column;
  align-items: center;
  flex-shrink: 0;
}

.level-num {
  font-family: var(--font-display);
  font-size: 22px;
  font-weight: 700;
  color: var(--accent-gold);
  line-height: 1;
}

.level-label {
  font-family: var(--font-display);
  font-size: 9px;
  letter-spacing: 0.1em;
  text-transform: uppercase;
  color: var(--text-muted);
}

.char-stats {
  display: flex;
  gap: 0;
}

.char-stat {
  flex: 1;
  text-align: center;
  padding: 6px 0;
}

.char-stat-val {
  display: block;
  font-family: var(--font-display);
  font-size: 15px;
  font-weight: 600;
  color: var(--text-primary);
}

.char-stat-label {
  display: block;
  font-family: var(--font-display);
  font-size: 9px;
  letter-spacing: 0.1em;
  text-transform: uppercase;
  color: var(--text-muted);
  margin-top: 2px;
}

.char-conditions {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  margin-top: 10px;
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

.char-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 12px;
}

.char-xp {
  font-family: var(--font-display);
  font-size: 11px;
  color: var(--text-muted);
  letter-spacing: 0.08em;
}

.char-arrow {
  color: var(--accent-gold-dim);
  font-size: 16px;
  transition: transform 0.2s;
}

.character-card:hover .char-arrow {
  transform: translateX(4px);
  color: var(--accent-gold);
}

/* Empty state */
.empty-state {
  text-align: center;
  padding: 64px 32px;
}

.empty-icon { font-size: 48px; margin-bottom: 16px; }

.empty-title {
  font-family: var(--font-display);
  font-size: 20px;
  color: var(--text-secondary);
  margin-bottom: 8px;
}

.empty-sub {
  font-family: var(--font-body);
  font-style: italic;
  color: var(--text-muted);
}

/* Create section */
.create-section {
  display: flex;
  justify-content: center;
  margin-top: 32px;
}

.create-btn {
  padding: 12px 32px;
  font-size: 13px;
}

.page-error {
  color: var(--accent-red-bright);
  text-align: center;
  padding: 16px;
}
</style>
