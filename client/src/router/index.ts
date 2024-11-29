import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '@/views/HomeView.vue'
import LoginView from '@/views/LoginView.vue'
import RegisterView from '@/views/RegisterView.vue'
import ProfileView from '@/views/ProfileView.vue'
import BoardView from '@/views/BoardView.vue'
import StatsView from '@/views/StatsView.vue'
import WipView from '@/views/WipView.vue'
import PlayerStatsView from '@/views/PlayerStatsView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
    },
    {
      path: '/login',
      name: 'login',
      component: LoginView,
    },
    {
      path: '/register',
      name: 'register',
      component: RegisterView,
    },
    {
      path: '/profile',
      name: 'profile',
      component: ProfileView,
    },
    {
      path: '/stats',
      name: 'stats',
      component: StatsView,
    },
    {
      path: '/game',
      name: 'game',
      component: BoardView,
    },
    {
      path: '/player/:id',
      name: 'player',
      component: PlayerStatsView,
    },
    {
      path: '/wip',
      name: 'wip',
      component: WipView,
    },
  ],
})

router.beforeEach(async to => {
  const publicPages = ['/login', '/register', '/player:id']
  const authRequired = !publicPages.includes(to.path)

  // Probably should use the Pinia authStore
  if (authRequired) {
    try {
      const res = await fetch('/api/session')

      if (!res.ok) {
        return '/login'
      }
    } catch (err) {
      console.error(err)
      return '/login'
    }
  }
})

export default router
