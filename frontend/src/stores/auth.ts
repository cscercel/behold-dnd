import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User } from '@/types'
import { authAPI } from '@/api/auth'
import router from '@/router'


export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(localStorage.getItem('token'))
  const user = ref<User | null>(null)

  const isAuthenticated = computed(() => !!token.value)
  const isDM = computed(() => user.value?.role === 'dm')

  async function login(email: string, password: string) {
    const response = await authAPI.login(email, password)
    token.value = response.token
    user.value = response.user
    localStorage.setItem('token', response.token)
  }

  async function fetchMe() {
    if (!token.value) return
    user.value = await authAPI.me()
  }

  function logout() {
    token.value = null
    user.value = null
    localStorage.removeItem('token')
    router.push({ name: 'login' }) // Go back to login page
  }

  return { token, user, isAuthenticated, isDM, login, fetchMe, logout }
})
