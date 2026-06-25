import { getMe } from './api';

interface User { id: string; username: string; email: string; role: 'player' | 'dm'; }

let user    = $state<User | null>(null);
let loading = $state(true);

export function getUser()    { return user; }
export function isLoading()  { return loading; }

export async function initAuth() {
    const t = localStorage.getItem('token');
    if (t) {
        try { user = await getMe(); } catch { localStorage.removeItem('token'); }
    }
    loading = false;
}

export async function authLogin(token: string) {
    localStorage.setItem('token', token);
    user = await getMe();
}

export function authLogout() {
    localStorage.removeItem('token');
    user = null;
}
