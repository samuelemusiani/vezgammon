<template>
  <div class="flex h-full w-full items-center justify-center">
    <div
      class="flex h-[90%] w-[80%] flex-col items-center justify-center rounded-md border-8 border-primary bg-base-100"
    >
      <!-- crea un bottone per tornare indeitro dentro la card -->
      <button
        @click="router.push('/')"
        class="retro-button absolute left-[12%] top-[10%] p-2"
      >
        Back
      </button>
      <div class="text-center">
        <h1 class="retro-title mb-2 p-4 text-5xl font-bold">Leaderboard</h1>
        <div class="text-xl font-bold text-accent">Compete for the glory</div>
      </div>
      <div class="flex w-full flex-col items-center justify-center">
        <div
          class="no-scrollbar max-h-[calc(100vh-300px)] w-full max-w-4xl space-y-3 overflow-y-auto p-8"
        >
          <div
            v-for="(el, index) in leaderboard?.slice(0, 10)"
            :key="index"
            class="retro-box relative w-full cursor-pointer hover:scale-[1.02]"
            @mouseenter="play()"
            @click="router.push('/player/' + el.username)"
          >
            <div class="flex items-center justify-between">
              <div class="flex items-center space-x-4">
                <div
                  class="flex h-12 w-12 items-center justify-center text-2xl font-bold"
                  :style="{ backgroundColor: getColorForRank(index) }"
                >
                  {{ index + 1 }}
                </div>
                <h2
                  class="text-2xl font-bold"
                  :class="myUsername === el.username ? 'text-primary' : ''"
                >
                  {{ el.username }}
                </h2>
              </div>
              <div class="mr-4 text-xl font-semibold text-accent">
                {{ el.elo }}
              </div>
            </div>
          </div>
          <!-- divisorio -->
          <div v-if="Me && myRank" class="h-1 w-full bg-primary"></div>
          <div
            v-if="Me && myRank"
            class="retro-box relative w-full cursor-pointer hover:scale-[1.02]"
            @mouseenter="play()"
            @click="router.push('/player/' + Me?.username)"
          >
            <div class="flex items-center justify-between">
              <div class="flex items-center space-x-4">
                <div
                  class="flex h-12 w-12 items-center justify-center text-2xl font-bold"
                >
                  {{ myRank + 1 }}
                </div>
                <h2 class="text-2xl font-bold text-primary">
                  {{ Me.username }}
                </h2>
              </div>
              <div class="mr-4 text-xl font-semibold text-accent">
                {{ Me.elo }}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import type { GameStats, LeaderBoardUser } from '@/utils/types'
import { useSound } from '@vueuse/sound'
import buttonSfx from '@/utils/sounds/button.mp3'
import router from '@/router'

const { play } = useSound(buttonSfx, { volume: 0.3 })

const leaderboard = ref<LeaderBoardUser[] | null>(null)
const myUsername = ref<string | null>(null)
const myRank = ref<number | null>(null)
const Me = ref<LeaderBoardUser | null>(null)

onMounted(async () => {
  await fetchMe()
  await fetchLeaderboard()
})

async function fetchLeaderboard() {
  const response = await fetch('/api/stats')
  const stats: GameStats = await response.json()
  leaderboard.value = stats.leaderboard.sort(
    (a: LeaderBoardUser, b: LeaderBoardUser) => {
      if (b.elo !== a.elo) {
        return b.elo - a.elo
      }
      return a.username.localeCompare(b.username)
    },
  )
  console.log(leaderboard.value, myUsername.value)
  if (
    !leaderboard.value
      .slice(0, 10)
      .map(el => el.username)
      .includes(myUsername.value || '')
  ) {
    console.log('not found')
    Me.value =
      leaderboard.value.find(el => el.username === myUsername.value) || null
    myRank.value = leaderboard.value.findIndex(
      el => el.username === myUsername.value,
    )
  }
}

const fetchMe = async () => {
  try {
    console.log('fetchMe')
    const response = await fetch('/api/session')
    console.log(response)
    const user = await response.json()
    console.log(user)
    myUsername.value = user.username
  } catch (error) {
    console.error('me: ' + error)
  }
}

// Function to generate dynamic colors for rank indicators
function getColorForRank(index: number) {
  const colors = [
    '#FFD700', // Gold
    '#C0C0C0', // Silver
    '#CD7F32', // Bronze
  ]
  return colors[index] || ''
}
</script>

<style scoped>
/* Hide scrollbar for Chrome, Safari and Opera */
.no-scrollbar::-webkit-scrollbar {
  display: none;
}
/* Hide scrollbar for IE, Edge and Firefox */
.no-scrollbar {
  -ms-overflow-style: none; /* IE and Edge */
  scrollbar-width: none; /* Firefox */
}
</style>
