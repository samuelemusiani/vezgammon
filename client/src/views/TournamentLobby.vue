<template>
  <div class="flex h-full w-full items-center justify-center">
    <div
      class="relative flex h-[94%] w-[95%] flex-col rounded-md border-4 border-primary bg-base-100 sm:w-[90%] md:w-[85%] lg:w-[80%] lg:border-8"
    >
      <div
        class="flex h-full flex-col items-center justify-center overflow-y-auto p-4"
      >
        <button
          @click="router.push('/tournaments')"
          class="retro-button absolute left-4 top-4 z-50"
        >
          Back
        </button>
        <div v-if="tournament" class="w-full">
          <div class="flex flex-1 flex-col items-center">
            <h1
              class="retro-title mb-8 w-60 text-center text-2xl font-bold md:w-full md:p-4 md:text-3xl lg:text-4xl xl:text-7xl"
            >
              Tournament {{ showBracket ? 'Brackets' : 'Lobby' }}
            </h1>
            <div class="mb-2 font-bold text-accent md:text-lg lg:text-xl">
              Owner:
              {{ tournament?.owner == myUsername ? 'me' : tournament?.owner }}
            </div>
            <div class="w-[100px]"></div>
          </div>
          <div v-if="!showBracket">
            <div
              class="mb-4 grid grid-cols-2 grid-rows-2 gap-2 sm:mb-6 sm:gap-4 md:mb-8"
            >
              <div
                v-for="(player, index) in tournament?.users"
                :key="index"
                class="retro-box relative flex min-h-[40px] items-center justify-center p-2 sm:min-h-[50px] sm:p-4 md:min-h-[60px]"
                :class="{
                  'text-primary': player === myUsername,
                  italic: ['Enzo', 'Caterina', 'Giovanni'].includes(player),
                  'text-black': player !== myUsername,
                }"
              >
                <span class="text-sm font-semibold sm:text-base md:text-lg">
                  {{ index + 1 }}. {{ player }}
                </span>
                <button
                  v-if="
                    ['Enzo', 'Caterina', 'Giovanni'].includes(player) &&
                    myUsername == tournament.owner
                  "
                  class="absolute right-1 top-1 cursor-pointer text-red-500 transition-colors hover:scale-[1.10] hover:text-red-600"
                  @click="deleteBot(player)"
                >
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    width="25"
                    height="25"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2.5"
                  >
                    <rect x="3" y="3" width="18" height="18" rx="2" />
                    <path d="M15.5 8.5l-7 7" />
                    <path d="M8.5 8.5l7 7" />
                  </svg>
                </button>
              </div>

              <div
                v-for="n in 4 - tournament.users.length"
                :key="`empty-${n}`"
                class="flex items-center justify-center rounded-lg border-2 border-dashed border-gray-200 bg-gray-50 p-4"
              >
                <span class="text-gray-400">Waiting for player...</span>
              </div>
            </div>

            <div
              v-if="tournament?.owner == myUsername"
              class="mt-2 flex justify-center gap-2"
            >
              <button class="retro-button" @click="deleteTournament">
                Delete Tournament
              </button>
              <div class="relative">
                <button
                  class="retro-button"
                  :disabled="showStartButton"
                  :style="{
                    textShadow: showStartButton
                      ? 'none'
                      : '2px 2px 0 rgba(0, 0, 0, 0.2)',
                  }"
                  @click="botDropDown = !botDropDown"
                >
                  Add Bot
                </button>

                <div
                  v-if="botDropDown"
                  class="absolute top-full mb-2 w-full rounded-lg"
                >
                  <ul class="flex flex-col gap-1 py-1.5">
                    <li
                      v-for="bot in ['easy', 'medium', 'hard']"
                      :key="bot"
                      class="retro-button w-full"
                      @click="addBot(bot)"
                    >
                      {{ bot }}
                    </li>
                  </ul>
                </div>
              </div>
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
          <div v-else class="flex flex-col gap-8">
            <div v-for="i in [0, 1]" :key="i" class="flex flex-col gap-4">
              <p class="text-lg font-bold text-[#8b4513]">
                {{ i === 0 ? 'Semi-Finals' : 'Finals' }}
              </p>
              <div v-for="j in [0, 1]" :key="j" class="flex flex-col gap-2">
                <div
                  class="retro-box flex flex-row items-center justify-between p-4"
                >
                  <div
                    class="flex w-[42%] flex-row items-center justify-start gap-1"
                    :style="
                      tournament?.games?.[i * 2 + j]?.status === 'winp2'
                        ? 'filter: grayscale(100%); opacity: 0.5'
                        : ''
                    "
                  >
                    <img
                      :src="avatar[i * 4 + j * 2]"
                      alt="avatar"
                      class="h-12 w-12 rounded-full object-cover"
                    />
                    <span>{{ boxes[i * 4 + j * 2] }}</span>
                  </div>
                  <div class="flex w-1/3 items-center justify-center">
                    <span v-if="i === 0">vs</span>
                    <span v-else-if="j === 1">
                      <svg
                        height="64"
                        width="64"
                        version="1.1"
                        id="Layer_1"
                        xmlns="http://www.w3.org/2000/svg"
                        xmlns:xlink="http://www.w3.org/1999/xlink"
                        viewBox="0 0 300.439 300.439"
                        xml:space="preserve"
                        fill="#000000"
                        class="mx-auto"
                      >
                        <g id="SVGRepo_bgCarrier" stroke-width="0"></g>
                        <g
                          id="SVGRepo_tracerCarrier"
                          stroke-linecap="round"
                          stroke-linejoin="round"
                        ></g>
                        <g id="SVGRepo_iconCarrier">
                          <g>
                            <path
                              style="fill: #bf392c"
                              d="M276.967,0h-84.498L70.415,178.385h84.498L276.967,0z"
                            ></path>
                            <path
                              style="fill: #e2574c"
                              d="M23.472,0h84.498l122.053,178.385h-84.498L23.472,0z"
                            ></path>
                            <path
                              style="fill: #efc75e"
                              d="M154.914,93.887c57.271,0,103.276,46.005,103.276,103.276s-46.005,103.276-103.276,103.276 S51.638,254.434,51.638,197.163S97.643,93.887,154.914,93.887z"
                            ></path>
                            <path
                              style="fill: #d7b354"
                              d="M154.914,122.053c-41.31,0-75.11,33.799-75.11,75.11s33.799,75.11,75.11,75.11 s75.11-33.799,75.11-75.11S196.224,122.053,154.914,122.053z M154.914,253.495c-30.983,0-56.332-25.35-56.332-56.332 s25.35-56.332,56.332-56.332s56.332,25.35,56.332,56.332S185.896,253.495,154.914,253.495z"
                            ></path>
                          </g>
                        </g>
                      </svg>
                    </span>
                    <span v-else>
                      <svg
                        height="64"
                        width="64"
                        version="1.1"
                        id="Layer_1"
                        xmlns="http://www.w3.org/2000/svg"
                        xmlns:xlink="http://www.w3.org/1999/xlink"
                        viewBox="0 0 300.439 300.439"
                        xml:space="preserve"
                        fill="#000000"
                        class="mx-auto"
                      >
                        <g id="SVGRepo_bgCarrier" stroke-width="0"></g>
                        <g
                          id="SVGRepo_tracerCarrier"
                          stroke-linecap="round"
                          stroke-linejoin="round"
                        ></g>
                        <g id="SVGRepo_iconCarrier">
                          <g>
                            <path
                              style="fill: #bf392c"
                              d="M276.967,0h-84.498L70.415,178.385h84.498L276.967,0z"
                            ></path>
                            <path
                              style="fill: #e2574c"
                              d="M23.472,0h84.498l122.053,178.385h-84.498L23.472,0z"
                            ></path>
                            <path
                              style="fill: #ed9d5d"
                              d="M154.914,93.887c57.271,0,103.276,46.005,103.276,103.276s-46.005,103.276-103.276,103.276 S51.638,254.434,51.638,197.163S97.643,93.887,154.914,93.887z"
                            ></path>
                            <path
                              style="fill: #d58d54"
                              d="M154.914,122.053c-41.31,0-75.11,33.799-75.11,75.11s33.799,75.11,75.11,75.11 s75.11-33.799,75.11-75.11S196.224,122.053,154.914,122.053z M154.914,253.495c-30.983,0-56.332-25.35-56.332-56.332 s25.35-56.332,56.332-56.332s56.332,25.35,56.332,56.332S185.896,253.495,154.914,253.495z"
                            ></path>
                          </g>
                        </g>
                      </svg>
                    </span>
                  </div>
                  <div
                    class="flex w-[42%] flex-row items-center justify-end gap-1"
                    :style="
                      tournament?.games?.[i * 2 + j]?.status === 'winp1'
                        ? 'filter: grayscale(100%); opacity: 0.5'
                        : ''
                    "
                  >
                    <span>{{ boxes[i * 4 + j * 2 + 1] }}</span>
                    <img
                      :src="avatar[i * 4 + j * 2 + 1]"
                      alt="avatar"
                      class="h-12 w-12 rounded-full object-cover"
                    />
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div v-else class="flex flex-col items-center gap-6">
          <h1 class="text-3xl font-bold text-primary">
            Tournament does not exists
          </h1>
          <button class="retro-button" @click="router.push('/')">
            Return to Home
          </button>
        </div>
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

const $toast = useToast()

const tournament = ref<Tournament | null>(null)
const myUsername = ref('')
const botDropDown = ref(false)
const showBracket = ref(false)
const showStartButton = ref(false)
const boxes = ref<Array<string>>([
  'Player 1',
  'Player 2',
  'Player 3',
  'Player 4',
  'Consolation 1',
  'Consolation 2',
  'Finalist 1',
  'Finalist 2',
])
const avatar = ref<Array<string>>(
  new Array(8).fill('https://api.dicebear.com/9.x/icons/svg?seed=Christian'),
)

const tournamentId = router.currentRoute.value.params.id
const webSocketStore = useWebSocketStore()

const fetchTournament = async () => {
  try {
    const response = await vfetch(`/api/tournament/${tournamentId}`)
    tournament.value = await response.json()
    console.log(tournament.value)
    showStartButton.value = tournament.value?.users.length === 4
    if (
      tournament.value?.status === 'in_progress' ||
      tournament.value?.status === 'ended'
    ) {
      if (tournament.value?.games[0]) {
        boxes.value[0] = tournament.value?.games[0].player1
        avatar.value[0] = await vfetch(
          `/api/player/${tournament.value?.games[0].player1}/avatar`,
        ).then(res => res.json())
        boxes.value[1] = tournament.value?.games[0].player2
        avatar.value[1] = await vfetch(
          `/api/player/${tournament.value?.games[0].player2}/avatar`,
        ).then(res => res.json())
      }
      if (tournament.value?.games[1]) {
        boxes.value[2] = tournament.value?.games[1].player1
        avatar.value[2] = await vfetch(
          `/api/player/${tournament.value?.games[1].player1}/avatar`,
        ).then(res => res.json())
        boxes.value[3] = tournament.value?.games[1].player2
        avatar.value[3] = await vfetch(
          `/api/player/${tournament.value?.games[1].player2}/avatar`,
        ).then(res => res.json())
      }
      if (tournament.value?.games[2]) {
        boxes.value[4] = tournament.value?.games[2].player1
        avatar.value[4] = await vfetch(
          `/api/player/${tournament.value?.games[2].player1}/avatar`,
        ).then(res => res.json())
        boxes.value[5] = tournament.value?.games[2].player2
        avatar.value[5] = await vfetch(
          `/api/player/${tournament.value?.games[2].player2}/avatar`,
        ).then(res => res.json())
      }
      if (tournament.value?.games[3]) {
        boxes.value[6] = tournament.value?.games[3].player1
        avatar.value[6] = await vfetch(
          `/api/player/${tournament.value?.games[3].player1}/avatar`,
        ).then(res => res.json())
        boxes.value[7] = tournament.value?.games[3].player2
        avatar.value[7] = await vfetch(
          `/api/player/${tournament.value?.games[3].player2}/avatar`,
        ).then(res => res.json())
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
  showBracket.value =
    tournament.value?.status === 'in_progress' ||
    tournament.value?.status === 'ended'
  try {
    webSocketStore.connect()
    webSocketStore.addMessageHandler(handleMessage)
  } catch (error) {
    console.error('websocket: ' + error)
  }
  const response = await vfetch('/api/play')
  if (response.ok) {
    $toast.success('Get ready for the next round!')
    setTimeout(() => {
      router.push('/game')
    }, 5000)
  }
})

onUnmounted(() => {
  webSocketStore.removeMessageHandler(handleMessage)
})

const handleMessage = async (message: WSMessage) => {
  if (message.type === 'tournament_cancelled') {
    $toast.error('Tournament has been cancelled')
    await router.push('/')
  } else if (message.type === 'tournament_new_user_enrolled') {
    await fetchTournament()
    $toast.info('Someone joined the tournament :)')
  } else if (message.type === 'tournament_new_bot_enrolled') {
    await fetchTournament()
    $toast.info('A bot has been added to the tournament')
  } else if (message.type === 'tournament_user_left') {
    await fetchTournament()
    $toast.warning('Someone left the tournament :(')
  } else if (message.type == 'tournament_bot_left') {
    await fetchTournament()
    $toast.warning('A bot has been removed from the tournament')
  } else if (message.type === 'game_tournament_ready') {
    await fetchTournament()
    if (!showBracket.value) $toast.success('Tournament is starting!')
    else {
      $toast.success('Get ready for the next round!')
    }
    showBracket.value = true
    setTimeout(() => {
      router.push('/game')
    }, 5000)
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

async function addBot(difficulty: string) {
  let username
  switch (difficulty.toLowerCase()) {
    case 'easy':
      username = 'Enzo'
      break
    case 'medium':
      username = 'Caterina'
      break
    case 'hard':
      username = 'Giovanni'
      break
    default:
      username = 'Invalid'
      break
  }
  try {
    vfetch(`/api/tournament/${tournamentId}/invite`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify([{ username }]),
    })
    botDropDown.value = false
  } catch (error) {
    console.error(error)
  }
}

function deleteBot(username: string) {
  try {
    vfetch(`/api/tournament/${tournamentId}/deletebot`, {
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify([{ username }]),
    })
  } catch (error) {
    console.error(error)
  }
}
</script>
