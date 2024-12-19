<template>
  <div class="flex h-full w-full items-center justify-center">
    <div
      class="flex h-[90%] w-[80%] flex-col items-center justify-center rounded-md border-8 border-primary bg-base-100"
    >
      <div v-if="tournament" class="text-center">
        <div v-if="!showBracket">
          <h1 class="retro-title mb-6 text-5xl font-bold">Tournament Lobby</h1>
          <p class="font-semibold text-accent">
            Owner: {{ owner ? 'me' : tournament?.owner }}
          </p>
          <div class="mb-16 mt-16 grid grid-cols-2 grid-rows-2 gap-4">
            <div
              v-for="(player, index) in tournament?.users"
              :key="index"
              class="retro-box p-4"
              :class="{
                'text-primary': player === myUsername,
                'text-black': player !== myUsername,
              }"
            >
              <span class="font-semibold">{{ player }}</span>
            </div>

            <div
              v-for="n in 4 - tournament.users.length"
              :key="`empty-${n}`"
              class="flex items-center justify-center rounded-lg border-2 border-dashed border-gray-200 bg-gray-50 p-4"
            >
              <span class="text-gray-400">Waiting for player...</span>
            </div>
          </div>

          <div v-if="owner" class="mt-2 flex justify-center gap-2">
            <button class="retro-button" @click="deleteTournament">
              Delete Tournament
            </button>
            <button
              class="retro-button"
              @click="startTournament"
              :disabled="!showStartButton"
              :style="{
                textShadow: !showStartButton
                  ? 'none'
                  : '2px 2px 0 rgba(0, 0, 0, 0.2)',
              }"
            >
              Start Tournament
            </button>
          </div>
          <div v-else class="mt-2 flex justify-center gap-2">
            <button class="retro-button" @click="exitTournament">
              Exit Tournament
            </button>
            <button
              class="retro-button"
              disabled
              :style="{ textShadow: 'none' }"
            >
              {{
                tournament.users.length === 4
                  ? 'Tournament full'
                  : 'Waiting for players...'
              }}
            </button>
          </div>
        </div>

        <!-- Tournament Bracket -->
        <div v-else class="flex h-full w-full flex-col space-y-4">
          <h2 class="retro-title mb-4 text-5xl font-bold">
            Tournament Bracket
          </h2>
          <p class="font-semibold text-accent">
            Owner: {{ owner ? 'me' : tournament?.owner }}
          </p>
          <div class="flex flex-row justify-between gap-8">
            <div class="flex w-1/4 flex-col gap-4">
              <!-- Semi-Final 1 -->
              <div class="flex flex-col items-center space-y-2">
                <div
                  v-for="(box, index) in boxes.slice(0, 2)"
                  :key="index"
                  class="retro-box w-full p-3 text-center font-semibold"
                  :style="{
                    color:
                      tournament?.games[0]?.status === 'winp1'
                        ? index === 0
                          ? 'green'
                          : 'red'
                        : tournament?.games[0]?.status === 'winp2'
                          ? index === 1
                            ? 'green'
                            : 'red'
                          : '',
                  }"
                >
                  {{ box }}
                </div>
              </div>

              <!-- Semi-Final 2 -->
              <div class="flex flex-col items-center space-y-2">
                <div
                  v-for="(box, index) in boxes.slice(2, 4)"
                  :key="index"
                  class="retro-box w-full p-3 text-center font-semibold"
                  :style="{
                    color:
                      tournament?.games[1]?.status === 'winp1'
                        ? index === 0
                          ? 'green'
                          : 'red'
                        : tournament?.games[1]?.status === 'winp2'
                          ? index === 1
                            ? 'green'
                            : 'red'
                          : '',
                  }"
                >
                  {{ box }}
                </div>
              </div>
            </div>
            <div class="flex w-full flex-col items-center gap-2 space-y-2">
              <!-- Final 1 place-->
              <div
                class="mt-10 flex h-1/4 w-full flex-row items-center space-x-2"
              >
                <div
                  v-for="(box, index) in finals.slice(0, 2)"
                  :key="index"
                  class="retro-box h-full w-full p-3 text-center text-2xl font-bold"
                  :style="{
                    color:
                      tournament?.games[2]?.status === 'winp1'
                        ? index === 0
                          ? 'green'
                          : 'red'
                        : tournament?.games[2]?.status === 'winp2'
                          ? index === 1
                            ? 'green'
                            : 'red'
                          : '',
                  }"
                >
                  {{ box }}
                </div>
              </div>
              <!-- Final 3 place-->
              <div
                class="mt-10 flex h-1/4 w-3/4 flex-row items-center space-x-2"
              >
                <div
                  v-for="(box, index) in finals.slice(2, 4)"
                  :key="index"
                  class="retro-box h-full w-full p-3.5 text-center font-bold"
                  :style="{
                    color:
                      tournament?.games[3]?.status === 'winp1'
                        ? index === 0
                          ? 'green'
                          : 'red'
                        : tournament?.games[3]?.status === 'winp2'
                          ? index === 1
                            ? 'green'
                            : 'red'
                          : '',
                  }"
                >
                  {{ box }}
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div v-else class="flex flex-col items-center gap-6">
        <h1 class="text-3xl font-bold text-primary">Not in Tournament</h1>
        <button
          class="rounded-lg bg-blue-500 px-6 py-3 text-white transition-colors duration-300 hover:bg-blue-600"
          @click="router.push('/')"
        >
          Return to Home
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
import router from '@/router'
import { useWebSocketStore } from '@/stores/websocket'
import type { Tournament, WSMessage } from '@/utils/types'
import { vfetch } from '@/utils/fetch'

import { useToast } from 'vue-toast-notification'
import 'vue-toast-notification/dist/theme-sugar.css'

const $toast = useToast()

const tournament = ref<Tournament | null>(null)
const myUsername = ref('')
const showBracket = ref(false)
const showStartButton = ref(false)
const boxes = ref<Array<string>>([
  'Player 1',
  'Player 2',
  'Player 3',
  'Player 4',
])
const finals = ref<Array<string>>([
  'Finalist 1',
  'Finalist 2',
  'Consolation 1',
  'Consolation 2',
])

const tournamentId = router.currentRoute.value.params.id
const owner = ref<boolean>(false)
const webSocketStore = useWebSocketStore()

const fetchTournament = async () => {
  try {
    const response = await vfetch(`/api/tournament/${tournamentId}`)
    tournament.value = await response.json()
    if (tournament.value?.users.length === 4) showStartButton.value = true
    if (tournament.value?.status === 'in_progress') {
      showBracket.value = true
      if (tournament.value?.games[0]) {
        boxes.value[0] = tournament.value?.games[0].player1
        boxes.value[1] = tournament.value?.games[0].player2
      }
      if (tournament.value?.games[1]) {
        boxes.value[2] = tournament.value?.games[1].player1
        boxes.value[3] = tournament.value?.games[1].player2
      }
      if (tournament.value?.games[2]) {
        finals.value[0] = tournament.value?.games[3].player1
        finals.value[1] = tournament.value?.games[3].player2
      }
      if (tournament.value?.games[3]) {
        finals.value[2] = tournament.value?.games[2].player1
        finals.value[3] = tournament.value?.games[2].player2
      }
    }
  } catch (error) {
    console.error('tournament: ' + error)
  }
}

const fetchMe = async () => {
  try {
    const response = await vfetch('/api/session')
    const user = await response.json()
    myUsername.value = user.username
  } catch (error) {
    console.error('me: ' + error)
  }
}

onMounted(async () => {
  await fetchTournament()
  await fetchMe()
  owner.value = tournament.value?.owner === myUsername.value
  showBracket.value = tournament.value?.status === 'in_progress'
  try {
    webSocketStore.connect()
    webSocketStore.addMessageHandler(handleMessage)
  } catch (error) {
    console.error('websocket: ' + error)
  }
})

onUnmounted(() => {
  webSocketStore.removeMessageHandler(handleMessage)
})

const handleMessage = async (message: WSMessage) => {
  console.log('TOURNAMENTS: Received message:', message)
  if (message.type === 'tournament_cancelled') {
    $toast.error('Tournament has been cancelled')
    await router.push('/')
  } else if (message.type === 'tournament_new_user_enrolled') {
    await fetchTournament()
    $toast.info('Someone joined the tournament :)')
  } else if (message.type === 'tournament_user_left') {
    await fetchTournament()
    $toast.warning('Someone left the tournament :(')
  } else if (message.type === 'game_tournament_ready') {
    await fetchTournament()
    if (!showBracket.value) $toast.success('Tournament is starting!')
    else {
      $toast.success('Get ready for the next round!')
    }
    showBracket.value = true
    setTimeout(() => {
      router.push('/game')
    }, 3000)
  }
}

function startTournament() {
  try {
    vfetch(`/api/tournament/${tournamentId}/start`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
    })
  } catch (error) {
    console.log(error)
  }
}

function exitTournament() {
  try {
    vfetch(`/api/tournament/${tournamentId}`, {
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json',
      },
    })
    $toast.warning('You left this tournament')
    router.push('/')
  } catch (error) {
    console.error(error)
  }
}

function deleteTournament() {
  try {
    vfetch(`/api/tournament/${tournamentId}/cancel`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
    })
  } catch (error) {
    console.error(error)
  }
}
</script>
