import { client } from './client'
import type { Feature, UpdateFeaturePayload } from '@/types'


export const featuresAPI = {
    list: (characterId: string) =>
        client.get<Feature[]>(`/characters/${characterId}/features`),

    create: (characterId: string, data: Partial<Feature>) =>
        client.post<Feature>(`/characters/${characterId}/features`, data),

    update: (characterId: string, featureId: string, data: UpdateFeaturePayload) =>
        client.patch<Feature>(`/characters/${characterId}/features/${featureId}`, data),

    delete: (characterId: string, featureId: string) =>
        client.delete<null>(`/characters/${characterId}/features/${featureId}`),
}
