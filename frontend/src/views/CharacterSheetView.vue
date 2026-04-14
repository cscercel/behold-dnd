<script setup lang="ts">
import { onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { useCharactersStore } from '@/stores/characters'
import { useCharacter } from '@/composables/useCharacter'
import AbilityScores from '@/components/character/AbilityScores.vue'
import HitPoints from '@/components/character/HitPoints.vue'
import SkillList from '@/components/character/SkillList.vue'

const route = useRoute()
const store = useCharactersStore()
const id = route.params.id as string

onMounted(() => store.fetchOne(id))

const derived = computed(() =>
  store.current ? useCharacter(store.current) : null
)
</script>

<template>
  <div class="character-sheet">
    <p v-if="store.loading">Loading...</p>
    <p v-else-if="store.error" class="error">{{ store.error }}</p>

    <template v-else-if="store.current && derived">
      <header>
        <h1>{{ store.current.name }}</h1>
        <p>
          {{ store.current.race }} {{ store.current.class }}
          · Level {{ store.current.level }}
          · {{ store.current.background }}
        </p>
        <p>
          Initiative: {{ derived.initiative.value }}
          · Passive Perception: {{ derived.passivePerception.value }}
          · Proficiency Bonus: {{ derived.signedModifier(derived.proficiencyBonus.value) }}
          · AC: {{ store.current.armor_class }}
          · Speed: {{ store.current.speed }}ft
        </p>
      </header>

      <div class="sheet-grid">
        <AbilityScores
          :character="store.current"
          :modifiers="derived.modifiers.value"
          :saving-throws="derived.savingThrows.value"
          :proficiency-bonus="derived.proficiencyBonus.value"
        />

        <HitPoints
          :character="store.current"
          :hp-percentage="derived.hpPercentage.value"
          @damage="(amount) => store.applyDamage(id, amount)"
          @heal="(amount) => store.heal(id, amount)"
          @temp-hp="(amount) => store.addTempHP(id, amount)"
          @long-rest="store.longRest(id)"
        />

        <SkillList :skills="derived.skills.value" />
      </div>
    </template>
  </div>
</template>
