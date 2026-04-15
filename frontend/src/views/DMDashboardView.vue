<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const auth = useAuthStore()
</script>

<template>
  <div class="dm-page">
    <!-- Header -->
    <header class="page-header">
      <div class="header-brand">
        <span class="header-icon">⚔</span>
        <span class="header-title">Behold DnD</span>
      </div>
      <div class="header-actions">
        <span class="header-badge badge badge-dm">Dungeon Master</span>
        <span class="header-user">{{ auth.user?.username }}</span>
        <button class="btn btn-ghost" @click="auth.logout()">Logout</button>
      </div>
    </header>

    <main class="page-content">
      <!-- Hero -->
      <div class="hero fade-in">
        <h1 class="hero-title">Dungeon Master's Sanctum</h1>
        <p class="hero-sub">Command your world, shape your story</p>
      </div>

      <div class="ornament" />

      <!-- Main options -->
      <div class="dm-grid fade-in">

        <!-- Characters -->
        <div class="dm-card" @click="router.push({ name: 'characters' })">
          <div class="dm-card-icon">🧙</div>
          <div class="dm-card-body">
            <h2 class="dm-card-title">Characters</h2>
            <p class="dm-card-desc">View and manage all player characters, NPCs, and creatures in your campaign</p>
          </div>
          <span class="dm-card-arrow">→</span>
        </div>

        <!-- Combat -->
        <div class="dm-card dm-card-combat" @click="router.push({ name: 'combat' })">
          <div class="dm-card-icon">⚔️</div>
          <div class="dm-card-body">
            <h2 class="dm-card-title">Combat Tracker</h2>
            <p class="dm-card-desc">Run encounters, manage initiative order, track HP and conditions in real time</p>
          </div>
          <span class="dm-card-arrow">→</span>
        </div>

        <!-- Create Character -->
        <div class="dm-card" @click="router.push({ name: 'character-create' })">
          <div class="dm-card-icon">📜</div>
          <div class="dm-card-body">
            <h2 class="dm-card-title">Create Character</h2>
            <p class="dm-card-desc">Forge a new player character, NPC, or monster for your campaign</p>
          </div>
          <span class="dm-card-arrow">→</span>
        </div>

      </div>
    </main>
  </div>
</template>

<style scoped>
.dm-page {
  min-height: 100vh;
  background: var(--bg-dark);
}

/* Reuse same header style as CharacterListView */
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
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 16px;
}

.header-badge { font-size: 10px; }

.header-user {
  font-family: var(--font-display);
  font-size: 13px;
  color: var(--text-secondary);
}

/* Content */
.page-content {
  max-width: 860px;
  margin: 0 auto;
  padding: 48px 32px;
}

.hero { text-align: center; margin-bottom: 32px; }

.hero-title {
  font-family: var(--font-display);
  font-size: 34px;
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

/* DM cards */
.dm-grid {
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin-top: 24px;
}

.dm-card {
  display: flex;
  align-items: center;
  gap: 20px;
  padding: 24px 28px;
  background: var(--bg-card);
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: border-color 0.2s, box-shadow 0.2s, transform 0.15s;
}

.dm-card:hover {
  border-color: var(--accent-gold-dim);
  box-shadow: 0 0 24px rgba(201,168,71,0.1);
  transform: translateX(4px);
}

.dm-card-combat:hover {
  border-color: var(--accent-red);
  box-shadow: 0 0 24px rgba(139,26,26,0.2);
}

.dm-card-icon {
  font-size: 32px;
  flex-shrink: 0;
  width: 56px;
  height: 56px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg-surface);
  border: 1px solid var(--border-light);
  border-radius: var(--radius-sm);
}

.dm-card-body { flex: 1; }

.dm-card-title {
  font-family: var(--font-display);
  font-size: 18px;
  font-weight: 600;
  letter-spacing: 0.04em;
  margin-bottom: 4px;
}

.dm-card-desc {
  font-family: var(--font-body);
  font-size: 15px;
  color: var(--text-secondary);
  line-height: 1.4;
}

.dm-card-arrow {
  font-size: 20px;
  color: var(--accent-gold-dim);
  flex-shrink: 0;
  transition: color 0.2s, transform 0.2s;
}

.dm-card:hover .dm-card-arrow {
  color: var(--accent-gold);
  transform: translateX(4px);
}
</style>
