const BASE = import.meta.env.VITE_API_URL || 'http://localhost:8080';

function token() { return localStorage.getItem('token'); }

async function req<T>(method: string, path: string, body?: unknown): Promise<T> {
    const headers: Record<string, string> = { 'Content-Type': 'application/json' };
    const t = token();
    if (t) headers['Authorization'] = `Bearer ${t}`;

    const res = await fetch(`${BASE}${path}`, {
        method,
        headers,
        body: body !== undefined ? JSON.stringify(body) : undefined,
    });

    if (res.status === 401) {
        const text = await res.text();
        const data = text ? JSON.parse(text) : null;
        if (t) {
            // Had a token and it's no longer valid — session expired
            localStorage.removeItem('token');
            window.location.href = '/login';
        }
        throw { status: 401, data };
    }

    const text = await res.text();
    const data = text ? JSON.parse(text) : null;
    if (!res.ok) throw { status: res.status, data };
    return data as T;
}

const get   = <T>(path: string)                 => req<T>('GET',    path);
const post  = <T>(path: string, body?: unknown) => req<T>('POST',   path, body);
const put   = <T>(path: string, body?: unknown) => req<T>('PUT',    path, body);
const patch = <T>(path: string, body?: unknown) => req<T>('PATCH',  path, body);
const del   = <T>(path: string)                 => req<T>('DELETE', path);

// Auth
export const login = (email: string, password: string) => 
    post<{token:string}>('/auth/login', { email, password });
export const register = (username: string, email: string, password: string, registration_code: string) =>
    post<{ token: string }>('/auth/register', { username, email, password, registration_code });
export const getMe    = () => get<any>('/auth/me');

// Characters
export const listCharacters      = () => get<any[]>('/characters/');
export const getCharacter        = (id: string) => get<any>(`/characters/${id}`);
export const createCharacter     = (data: any)  => post<any>('/characters/', data);
export const deleteCharacter     = (id: string) => del<any>(`/characters/${id}`);
export const updateCharacterInfo = (id: string, data: any) => patch<any>(`/characters/${id}/info`, data);
export const applyDamage         = (id: string, amount: number) => post<any>(`/characters/${id}/damage`, { amount });
export const applyHeal           = (id: string, amount: number) => post<any>(`/characters/${id}/heal`, { amount });
export const addTempHP           = (id: string, amount: number) => post<any>(`/characters/${id}/temp-hp`, { amount });
export const recordDeathSave     = (id: string, success: boolean) => post<any>(`/characters/${id}/death-save`, { success });
export const longRest            = (id: string) => post<any>(`/characters/${id}/long-rest`);
export const shortRest           = (id: string, hit_dice_remaining: number, current_hp: number) =>
    post<any>(`/characters/${id}/short-rest`, { hit_dice_remaining, current_hp });
export const updateConditions    = (id: string, conditions: string[]) =>
    put<any>(`/characters/${id}/conditions`, { conditions });

// Inventory
export const listInventory       = (id: string) => get<any[]>(`/characters/${id}/inventory/`);
export const createInventoryItem = (id: string, data: any) => post<any>(`/characters/${id}/inventory/`, data);
export const updateInventoryItem = (id: string, itemId: string, data: any) =>
    patch<any>(`/characters/${id}/inventory/${itemId}`, data);
export const deleteInventoryItem = (id: string, itemId: string) =>
    del<any>(`/characters/${id}/inventory/${itemId}`);
export const attuneItem          = (id: string, itemId: string) =>
    post<any>(`/characters/${id}/inventory/${itemId}/attune`);
export const unattuneItem        = (id: string, itemId: string) =>
    post<any>(`/characters/${id}/inventory/${itemId}/unattune`);

// Spells
export const listSpells          = (id: string) => get<any[]>(`/characters/${id}/spells/`);
export const createSpell         = (id: string, data: any) => post<any>(`/characters/${id}/spells/`, data);
export const deleteSpell         = (id: string, spellId: string) =>
    del<any>(`/characters/${id}/spells/${spellId}`);
export const toggleSpellPrepared = (id: string, spellId: string) =>
    post<any>(`/characters/${id}/spells/${spellId}/toggle-prepared`);
export const listSpellSlots      = (id: string) => get<any[]>(`/characters/${id}/spell-slots/`);
export const useSpellSlot        = (id: string, spell_level: number) =>
    post<any>(`/characters/${id}/spell-slots/use`, { spell_level });

// Combat
export const listEncounters    = () => get<any[]>('/combat/');
export const createEncounter   = (name: string) => post<any>('/combat/', { name });
export const getActiveEncounter= () => get<any>('/combat/active');
export const startEncounter    = (id: string) => post<any>(`/combat/${id}/start`);
export const endEncounter      = (id: string) => post<any>(`/combat/${id}/end`);
export const nextRound         = (id: string) => post<any>(`/combat/${id}/next-round`);
export const listParticipants  = (id: string) => get<any[]>(`/combat/${id}/participants`);
export const addParticipant    = (id: string, data: any) => post<any>(`/combat/${id}/participants`, data);
export const removeParticipant = (encId: string, partId: string) =>
    del<any>(`/combat/${encId}/participants/${partId}`);
export const participantDamage = (encId: string, partId: string, amount: number) =>
    post<any>(`/combat/${encId}/participants/${partId}/damage`, { amount });
export const participantHeal   = (encId: string, partId: string, amount: number) =>
    post<any>(`/combat/${encId}/participants/${partId}/heal`, { amount });
