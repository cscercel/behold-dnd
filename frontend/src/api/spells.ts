import { client } from './client'
import type { Spell, SpellSlot, UpdateSpellPayload } from '@/types'


export const spellsAPI = {
  list: (characterId: string) =>
    client.get<Spell[]>(`/characters/${characterId}/spells`),

  create: (characterId: string, data: Partial<Spell>) =>
    client.post<Spell>(`/characters/${characterId}/spells`, data),

  update: (characterId: string, spellId: string, data: UpdateSpellPayload) =>
    client.patch<Spell>(`/characters/${characterId}/spells/${spellId}`, data),

  delete: (characterId: string, spellId: string) =>
    client.delete<null>(`/characters/${characterId}/spells/${spellId}`),

  togglePrepared: (characterId: string, spellId: string) =>
    client.post<Spell>(`/characters/${characterId}/spells/${spellId}/toggle-prepared`),

  listSlots: (characterId: string) =>
    client.get<SpellSlot[]>(`/characters/${characterId}/spell-slots`),

  upsertSlot: (characterId: string, spellLevel: number, total: number, used: number) =>
    client.put<SpellSlot>(`/characters/${characterId}/spell-slots`, {
      spell_level: spellLevel,
      total: total,
      used: used,
    }),

  useSlot: (characterId: string, spellLevel: number) =>
    client.post<SpellSlot>(`/characters/${characterId}/spell-slots/use`, {
      spell_level: spellLevel,
    }),
}
