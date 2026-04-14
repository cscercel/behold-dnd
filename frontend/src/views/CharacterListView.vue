<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useCharactersStore } from '@/stores/characters'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const store = useCharactersStore()
const auth = useAuthStore()

onMounted(() => store.fetchAll())
</script>

<template>
  <div class="character-list">
    <header>
      <h1>Characters</h1>
      <nav>
        <button v-if="auth.isDM" @click="router.push({ name: 'combat' })">
          Combat Tracker
        </button>
        <button @click="auth.logout()">Logout</button>
      </nav>
    </header>

    <p v-if="store.loading">Loading...</p>
    <p v-else-if="store.error" class="error">{{ store.error }}</p>

    <div v-else class="grid">
      <div
        v-for="char in store.characters"
        :key="char.id"
        class="card"
        @click="router.push({ name: 'character-sheet', params: { id: char.id } })"
      >
        <h2>{{ char.name }}</h2>
        <p>{{ char.race }} {{ char.class }} — Level {{ char.level }}</p>
        <p>HP: {{ char.current_hp }} / {{ char.max_hp }}</p>
        <span v-if="char.is_npc" class="badge">NPC</span>
      </div>
    </div>
  </div>
</template>
