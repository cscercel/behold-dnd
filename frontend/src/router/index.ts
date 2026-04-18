import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/LoginView.vue'),
      meta: { requiresAuth: false },
    },
    {
      path: '/',
      name: 'characters',
      component: () => import('@/views/CharacterListView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/characters/create',
      name: 'character-create',
      component: () => import('@/views/CharacterCreateView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/dm',
      name: 'dm-dashboard',
      component: () => import('@/views/DMDashboardView.vue'),
      meta: { requiresAuth: true, requiresDM: true },
    },
    {
      path: '/characters/:id',
      name: 'character-sheet',
      component: () => import('@/views/CharacterSheetView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/combat',
      name: 'combat',
      component: () => import('@/views/CombatView.vue'),
      meta: { requiresAuth: true, requiresDM: true },
    },
  ],
})

router.beforeEach(async (to) => {
  const auth = useAuthStore()

  // Rehydrate user after page refresh
  if (auth.isAuthenticated && !auth.user) {
    await auth.fetchMe()
  }

  // Not logged in — redirect to login
  if (to.meta.requiresAuth && !auth.isAuthenticated) {
    return { name: 'login' }
  }

  // Already logged in — redirect away from login page
  if (to.name === 'login' && auth.isAuthenticated) {
    return auth.isDM ? { name: 'dm-dashboard' } : { name: 'characters' }
  }

  // Non-DM trying to access DM-only routes
  if (to.meta.requiresDM && !auth.isDM) {
    return { name: 'characters' }
  }
})

export default router
