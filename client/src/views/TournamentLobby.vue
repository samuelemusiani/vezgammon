<template>
  <div class="flex h-full w-full items-center justify-center">
    <div
      class="flex h-[94%] w-[80%] flex-col items-center justify-center overflow-y-auto rounded-md border-8 border-primary bg-base-100"
    >
      <button
        @click="router.push('/tournaments')"
        class="retro-button absolute left-[12%] top-[10%] p-2"
      >
        Back
      </button>
      <div v-if="tournament" class="text-center">
        <div class="mb-8 flex flex-col items-center text-center xl:mb-16">
          <h1
            class="retro-title mb-8 w-60 text-2xl font-bold md:w-full md:p-4 md:text-3xl lg:text-4xl xl:text-7xl"
          >
            Tournament {{ showBracket? 'Brackets' : 'Lobby' }}
          </h1>
          <div class="font-bold text-accent md:text-lg lg:text-xl">
            Owner: {{ tournament?.owner == myUsername ? 'me' : tournament?.owner }}
          </div>
        </div>
        <div v-if="!showBracket">
          <div class="mb-16 mt-16 grid grid-cols-2 grid-rows-2 gap-4">
            <div
              v-for="(player, index) in tournament?.users"
              :key="index"
              class="retro-box p-4"
              :class="{
                'text-primary': player === myUsername,
                'italic' : ['Enzo', 'Caterina', 'Giovanni'].includes(player),
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

          <div v-if="tournament?.owner == myUsername" class="mt-2 flex justify-center gap-2">
            <button class="retro-button" @click="deleteTournament">
              Delete Tournament
            </button>
            <div class="relative">
              <button
                class="retro-button"
                :disabled="showStartButton"
                :style="{
                  textShadow: showStartButton ? 'none' : '2px 2px 0 rgba(0, 0, 0, 0.2)',
                }"
                @click="botDropDown = !botDropDown"
              >
                Add Bot
              </button>

              <div
                v-if="botDropDown"
                class="absolute top-full mb-2 w-full rounded-lg"
              >
                <ul class="py-1.5 flex flex-col gap-1">
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
        <div v-else>
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
            </div>
          </div>
        </div>
      </div>

      <div v-else class="flex flex-col items-center gap-6">
        <h1 class="text-3xl font-bold text-primary">Tournament does not exists</h1>
        <button
          class="retro-button"
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
])
const finals = ref<Array<string>>([
  'Finalist 1',
  'Finalist 2',
  'Consolation 1',
  'Consolation 2',
])

const tournamentId = router.currentRoute.value.params.id
const webSocketStore = useWebSocketStore()

const fetchTournament = async () => {
  try {
    const response = await fetch(`/api/tournament/${tournamentId}`)
    tournament.value = await response.json()
    showStartButton.value = tournament.value?.users.length === 4
    if (tournament.value?.status === 'in_progress' || tournament.value?.status === 'ended') {
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
    const response = await fetch('/api/session')
    const user = await response.json()
    myUsername.value = user.username
  } catch (error) {
    console.error('me: ' + error)
  }
}

onMounted(async () => {
  await fetchTournament()
  await fetchMe()
  showBracket.value = tournament.value?.status === 'in_progress' || tournament.value?.status === 'ended'
  try {
    webSocketStore.connect()
    webSocketStore.addMessageHandler(handleMessage)
  } catch (error) {
    console.error('websocket: ' + error)
  }
  const response = await fetch('/api/play')
  if (response.ok) {
    $toast.success('Get ready for the next round!')
    setTimeout(() => {
      router.push('/game')
    }, 3000)
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
    }, 3000)
  }
}

function startTournament() {
  try {
    fetch(`/api/tournament/${tournamentId}/start`, {
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
    fetch(`/api/tournament/${tournamentId}`, {
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
    fetch(`/api/tournament/${tournamentId}/cancel`, {
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
  const username = difficulty.toLowerCase() == 'easy' ? 'Enzo' : difficulty == 'medium' ? 'Caterina' : 'Giovanni'
  try {
    fetch(`/api/tournament/${tournamentId}/invite`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify([{username}]),
    })
    botDropDown.value = false
  } catch (error) {
    console.error(error)
  }
}
</script>
