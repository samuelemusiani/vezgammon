import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '@/views/HomeView.vue'
import LoginView from '@/views/LoginView.vue'
import RegisterView from '@/views/RegisterView.vue'
import ProfileView from '@/views/ProfileView.vue'
import BoardView from '@/views/BoardView.vue'
import WipView from '@/views/WipView.vue'
import TournamentsView from '@/views/TournamentsView.vue'

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
      path: '/game',
      name: 'game',
      component: BoardView,
    },
    {
      path: '/tournaments',
      name: 'tournaments',
      component: TournamentsView,
    },
    {
      path: '/wip',
      name: 'wip',
      component: WipView,
    },
    {
      path: '/about',
      name: 'about',
      // route level code-splitting
      // this generates a separate chunk (About.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import('../views/AboutView.vue'),
    },
  ],
})

router.beforeEach(async to => {
  const publicPages = ['/login', '/register']
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
