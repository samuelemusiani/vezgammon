import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '@/views/HomeView.vue'
import LoginView from '@/views/LoginView.vue'
import RegisterView from '@/views/RegisterView.vue'
import ProfileView from '@/views/ProfileView.vue'
import TournamentsView from '@/views/TournamentsView.vue'
import TournamentLobbyView from '@/views/TournamentLobby.vue'
import BoardView from '@/views/BoardView.vue'
import StatsView from '@/views/StatsView.vue'
import WipView from '@/views/WipView.vue'
import PlayerStatsView from '@/views/PlayerStatsView.vue'
import InviteView from '@/views/InviteView.vue'
import LeaderBoardView from '@/views/LeaderBoardView.vue'

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
      path: '/player/:username',
      name: 'player',
      component: PlayerStatsView,
    },
    {
      path: '/tournaments',
      name: 'tournaments',
      component: TournamentsView,
    },
    {
      path: '/tournaments/:id',
      name: 'tournament_lobby',
      component: TournamentLobbyView,
    },
    {
      path: '/leaderboard',
      name: 'leaderboard',
      component: LeaderBoardView,
    },
    {
      path: '/wip',
      name: 'wip',
      component: WipView,
    },
    {
      path: '/invite/:code',
      name: 'invite',
      component: InviteView,
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
  const authRequired =
    !publicPages.includes(to.path) && !to.path.startsWith('/player')

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
