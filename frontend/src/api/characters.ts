import { client } from './client'
import type { Character, UpdateCharacterPayload } from '@/types'


export const characterAPI = {
  list: () =>
    client.get<Character[]>('/list-characters'),

  get: (id: string) =>
    client.get<Character>(`/characters/${id}`),

  create: (data: Partial<Character>) =>
    client.post<Character>('/characters', data),

  update: (id: string, data: UpdateCharacterPayload) =>
    client.patch<Character>(`/characters/${id}`, data),

  delete: (id: string) =>
    client.delete<null>(`/characters/${id}`),

  // Game mechanics
  damage: (id: string, amount: number) =>
    client.post<Character>(`/characters/${id}/damage`, { amount }),

  heal: (id: string, amount: number) =>
    client.post<Character>(`/characters/${id}/heal`, { amount }),

  tempHP: (id: string, amount: number) =>
    client.post<Character>(`/characters/${id}/temp-hp`, { amount }),

  deathSave: (id: string, success: boolean) =>
    client.post<Character>(`/characters/${id}/death-save`, { success }),

  longRest: (id: string) =>
    client.post<Character>(`/characters/${id}/long-rest`),

  shortRest: (id: string, hitDiceUsed: number, hpRegained: number) =>
    client.post<Character>(`/characters/${id}/short-rest`, {
      hit_dice_used: hitDiceUsed,
      hp_regained: hpRegained,
    }),

  updateConditions: (id: string, conditions: string[]) =>
    client.put<Character>(`/characters/${id}/conditions`, { conditions }),
}
