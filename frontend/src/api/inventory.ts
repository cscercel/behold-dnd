import { client } from './client'
import type { InventoryItem, UpdateInventoryItemPayload } from '@/types'


export const inventoryAPI = {
  list: (characterId: string) =>
    client.get<InventoryItem[]>(`/characters/${characterId}/inventory`),

  create: (characterId: string, data: Partial<InventoryItem>) =>
    client.post<InventoryItem>(`/characters/${characterId}/inventory`, data),

  update: (characterId: string, itemId: string, data: UpdateInventoryItemPayload) =>
    client.patch<InventoryItem>(`/characters/${characterId}/inventory/${itemId}`, data),

  delete: (characterId: string, itemId: string) =>
    client.delete<null>(`/characters/${characterId}/inventory/${itemId}`),

  attune: (characterId: string, itemId: string) =>
    client.post<InventoryItem>(`/characters/${characterId}/inventory/${itemId}/attune`),

  unattune: (characterId: string, itemId: string) =>
    client.post<InventoryItem>(`/characters/${characterId}/inventory/${itemId}/unattune`),
}
