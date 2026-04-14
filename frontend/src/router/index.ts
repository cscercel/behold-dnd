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

    if (to.meta.requiresAuth && !auth.isAuthenticated) {
        return { name: 'login' }
    }

    if (to.meta.requiresDM && !auth.isDM) {
        return { name: 'characters' }
    }

    // Rehydrate user info after a page refresh
    if (auth.isAuthenticated && !auth.user) {
        await auth.fetchMe()
    }
})

export default router
