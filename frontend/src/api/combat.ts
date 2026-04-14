import { client } from './client'
import type { CombatEncounter, CombatParticipant } from '@/types'


export const combatAPI = {
  list: () =>
    client.get<CombatEncounter[]>('/combat'),

  get: (id: string) =>
    client.get<CombatEncounter>(`/combat/${id}`),

  getActive: () =>
    client.get<CombatEncounter>('/combat/active'),

  create: (name: string) =>
    client.post<CombatEncounter>('/combat', { name }),

  delete: (id: string) =>
    client.delete<null>(`/combat/${id}`),

  start: (id: string) =>
    client.post<CombatEncounter>(`/combat/${id}/start`),

  end: (id: string) =>
    client.post<CombatEncounter>(`/combat/${id}/end`),

  nextRound: (id: string) =>
    client.post<CombatEncounter>(`/combat/${id}/next-round`),

  listParticipants: (encounterId: string) =>
    client.get<CombatParticipant[]>(`/combat/${encounterId}/participants`),

  addParticipant: (encounterId: string, data: {
    character_id?: string
    name?: string,
    initiative: number,
    current_hp?: number,
    max_hp?: number,
    temp_hp?: number,
    armor_class?: number,
    speed?: number
  }) =>
    client.post<CombatParticipant>(`/combat/${encounterId}/participants`, data),

  removeParticipant: (encounterId: string, participantId: string) =>
    client.delete<null>(`/combat/${encounterId}/participants/${participantId}`),

  damageParticipant: (encounterId: string, participantId: string, amount: number) =>
    client.post<CombatParticipant>(`/combat/${encounterId}/participants/${participantId}/damage`, { amount }),

  healParticipant: (encounterId: string, participantId: string, amount: number) =>
    client.post<CombatParticipant>(`/combat/${encounterId}/participants/${participantId}/heal`, { amount }),

  UpdateTempHp: (encounterId: string, participantId: string, amount: number) =>
    client.put<CombatParticipant>(`/combat/${encounterId}/participants/${participantId}/temp-hp`, { amount }),

  updateInitiative: (encounterId: string, participantId: string, initiative: number) =>
    client.put<CombatParticipant>(`/combat/${encounterId}/participants/${participantId}/initiative`, { initiative }),

  updateConditions: (encounterId: string, participantId: string, conditions: string[]) =>
    client.put<CombatParticipant>(`/combat/${encounterId}/participants/${participantId}/conditions`, { conditions }),

  toggleConcentration: (encounterId: string, participantId: string) =>
    client.post<CombatParticipant>(`/combat/${encounterId}/participants/${participantId}/toggle-concentration`),

  deactivate: (encounterId: string, participantId: string) =>
    client.post<CombatParticipant>(`/combat/${encounterId}/participants/${participantId}/deactivate`),
}
