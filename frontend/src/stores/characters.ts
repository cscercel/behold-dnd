import { defineStore } from "pinia"
import { ref } from 'vue'
import type { Character, UpdateCharacterPayload } from '@/types'
import { characterAPI } from '@/api/characters'


export const useCharacterStore = defineStore('characters', () => {
  const characters = ref<Character[]>([])
  const current = ref<Character | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function fetchAll() {
    loading.value = true
    error.value = null
    try {
      characters.value = await characterAPI.list()
    } catch (e: any) {
      error.value = e.message ?? 'Failed to load characters'
    } finally {
      loading.value = false
    }
  }

  async function fetchOne(id: string) {
    loading.value = true
    error.value = null
    try {
      current.value = await characterAPI.get(id)
    } catch (e: any) {
      error.value = e.message ?? 'Failed to load character'
    } finally {
      loading.value = false
    }
  }

  async function update(id: string, payload: UpdateCharacterPayload) {
    const updated = await characterAPI.update(id, payload)
    current.value = updated
    const index = characters.value.findIndex((c) => c.id === id)
    if (index !== -1) characters.value[index] = updated
    return updated
  }

  async function applyDamage(id: string, amount: number) {
    current.value = await characterAPI.damage(id, amount)
    return current.value
  }

  async function heal(id: string, amount: number) {
    current.value = await characterAPI.heal(id, amount)
    return current.value
  }

  async function addTempHP(id: string, amount: number) {
    current.value = await characterAPI.tempHP(id, amount)
    return current.value
  }

  async function recordDeathSave(id: string, success: boolean) {
    current.value = await characterAPI.deathSave(id, success)
    return current.value
  }

  async function longRest(id: string) {
    current.value = await characterAPI.longRest(id)
    return current.value
  }

  async function shortRest(id: string, hitDiceUsed: number, hpRegained: number) {
    current.value = await characterAPI.shortRest(id, hitDiceUsed, hpRegained)
    return current.value
  }

  return {
    characters,
    current,
    loading,
    error,
    fetchAll,
    fetchOne,
    update,
    applyDamage,
    heal,
    addTempHP,
    recordDeathSave,
    longRest,
    shortRest,
  }
})
