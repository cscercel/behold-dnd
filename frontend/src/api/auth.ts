import { client } from './client'
import router from '@/router'
import type { User } from '@/types'

export const authAPI = {
  login: (email: string, password: string) =>
    client.post<{ token: string; user: User }>('/auth/login', { email, password }),

  me: () =>
    client.get<User>('/auth/me'),

  register: (username: string, email: string, password: string, role: string, registrationCode: string) =>
    client.post<User>('/auth/register', {
      username,
      email,
      password,
      role,
      registration_code: registrationCode,
    }),
}
