import { useAuthStore } from '@/stores/auth'

const BASE_URL = 'http://localhost:8080'

async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
  const auth = useAuthStore()

  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(options.headers as Record<string, string>),
  }

  if (auth.token) {
    headers['Authorization'] = `Bearer ${auth.token}`
  }

  const response = await fetch(`${BASE_URL}/${path}`, {
    ...options,
    headers,
  })

  // Log out on token expiration
  if (response.status === 401) {
    auth.logout()
    window.location.href = '/login'
    throw new Error('Unauthorized')
  }

  // If no content
  if (response.status === 204) {
    return null as T
  }

  const data = await response.json()

  // Surface error message
  if (!response.ok) {
    throw new Error(data.error ?? 'Request failed')
  }

  return data as T
}

export const client = {
  get: <T>(path: string) =>
    request<T>(path),

  post: <T>(path: string, body?: unknown) =>
    request<T>(path, {
      method: 'POST',
      body: body !== undefined ? JSON.stringify(body) : undefined,
    }),

  patch: <T>(path: string, body: unknown) =>
    request<T>(path, {
      method: 'PATCH',
      body: JSON.stringify(body),
    }),

  put: <T>(path: string, body: unknown) =>
    request<T>(path, {
      method: 'PUT',
      body: JSON.stringify(body),
    }),

  delete: <T>(path: string) =>
    request<T>(path, { method: 'DELETE' }),
}
