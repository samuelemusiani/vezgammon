<template>
  <div class="flex h-full w-full items-center justify-center">
    <div class="flex h-[90%] w-[80%] flex-col items-center justify-center rounded-md border-8 border-primary bg-base-100">
      <div v-if="tournament" class="text-center">
        <div v-if="!showBracket">
          <h1 class="text-5xl font-bold mb-6 retro-title">Tournament Lobby</h1>
          <p class="text-accent font-semibold">Owner: {{ owner ? 'me' : tournament?.owner }}</p>
          <div class="grid grid-cols-2 grid-rows-2 gap-4 mt-16 mb-16">
          <div
              v-for="(player, index) in tournament.users" :key="index"
              class=" p-4 retro-box"
              :class="{
                'text-primary': player === myUsername,
                'text-black': player !== myUsername
              }"
            >
              <span class="font-semibold">{{ player }}</span>
            </div>

            <div
              v-for="n in (4 - tournament.users.length)"
              :key="`empty-${n}`"
              class="flex items-center justify-center p-4 rounded-lg bg-gray-50 border-2 border-dashed border-gray-200"
            >
              <span class="text-gray-400">Waiting for player...</span>
            </div>
        </div>

          <div v-if="owner" class="flex justify-center gap-2 mt-2">
            <button
              class="retro-button"
              @click="deleteTournament"
            >
              Delete Tournament
            </button>
            <button
              class="retro-button"
              @click="startTournament"
              :disabled="!showStartButton"
              :style="{ textShadow: !showStartButton ? 'none' : '2px 2px 0 rgba(0, 0, 0, 0.2)' }"
            >
              Start Tournament
            </button>
          </div>
          <div v-else>
            <button
              class="retro-button"
              @click="exitTournament"
            >
              Exit Tournament
            </button>
            <button
              class="retro-button"
              disabled
              :style="{ textShadow: 'none' }"
            >
              {{ tournament.users.length === 4 ? 'Tournament full' : 'Waiting for players...' }}
            </button>
          </div>
        </div>

        <!-- Tournament Bracket -->
        <div
          v-else
          class="w-full h-full flex flex-col space-y-4"
        >
          <h2 class="text-5xl font-bold retro-title mb-4">Tournament Bracket</h2>
          <p class="text-accent font-semibold">Owner: {{ owner ? 'me' : tournament?.owner }}</p>
          <div class="flex flex-row justify-between gap-8">
          <div class="flex flex-col gap-4 w-1/4">
            <!-- Semi-Final 1 -->
            <div class="flex flex-col items-center space-y-2">
              <div
                class="retro-box w-full p-3 text-center font-semibold"
                :class="{
            'text-green-500': semifinal1Winner === tournament.users[0],
            'text-red-500': semifinal1Winner === tournament.users[1]
          }"
              >
                {{ tournament.users[0] || 'TBD' }}
              </div>
              <div
                class="retro-box w-full p-3 text-center font-semibold"
                :class="{
            'text-green-500': semifinal1Winner === tournament.users[1],
            'text-red-500': semifinal1Winner === tournament.users[0]
          }"
              >
                {{ tournament.users[1] || 'TBD' }}
              </div>
            </div>

            <!-- Semi-Final 2 -->
            <div class="flex flex-col items-center space-y-2">
              <div
                class="retro-box w-full p-3 text-center font-semibold"
                :class="{
            'text-green-500': semifinal2Winner === tournament.users[2],
            'text-red-500': semifinal2Winner === tournament.users[3]
          }"
              >
                {{ tournament.users[2] || 'TBD' }}
              </div>
              <div
                class="retro-box w-full p-3 text-center font-semibold"
                :class="{
            'text-green-500': semifinal2Winner === tournament.users[3],
            'text-red-500': semifinal2Winner === tournament.users[2]
          }"
              >
                {{ tournament.users[3] || 'TBD' }}
              </div>
            </div>
          </div>
            <div class="flex flex-col items-center space-y-2 w-full gap-2">
            <!-- Final 1 place-->
            <div class="flex flex-row items-center space-x-2 w-full h-1/4 mt-10">
              <div
                class="retro-box w-full h-full p-3 text-center text-2xl font-bold"
                :class="{
            'text-green-500': tournamentWinner === semifinal1Winner,
            'text-red-500': tournamentWinner === semifinal2Winner
          }"
              >
                {{ semifinal1Winner || 'TBD' }}
              </div>
              <div
                class="retro-box w-full h-full p-3 text-center text-2xl font-bold"
                :class="{
            'text-green-500': tournamentWinner === semifinal2Winner,
            'text-red-500': tournamentWinner === semifinal1Winner
          }"
              >
                {{ semifinal2Winner || 'TBD' }}
              </div>
            </div>
              <!-- Final 3 place-->
              <div class="flex flex-row items-center space-x-2 w-2/3 h-1/4 mt-10">
                <div
                  class="retro-box w-full h-full p-3 text-center text-2xl font-bold"
                  :class="{
            'text-green-500': thirdplace === semifinal1Looser,
            'text-red-500': thirdplace === semifinal2Looser
          }"
                >
                  {{ semifinal1Looser || 'TBD' }}
                </div>
                <div
                  class="retro-box w-full h-full p-3 text-center text-2xl font-bold"
                  :class="{
            'text-green-500': thirdplace === semifinal2Looser,
            'text-red-500': thirdplace === semifinal1Looser
          }"
                >
                  {{ semifinal2Looser || 'TBD' }}
                </div>
              </div>
            </div>
          </div>

          <!-- Match Progression Buttons (for demonstration) -->
          <div class="flex space-x-4 mt-6">
            <button
              class="px-6 py-3 bg-blue-500 text-white rounded-lg hover:bg-blue-600"
              @click="progressSemifinal1"
            >
              Semi-Final 1 Winner
            </button>
            <button
              class="px-6 py-3 bg-blue-500 text-white rounded-lg hover:bg-blue-600"
              @click="progressSemifinal2"
            >
              Semi-Final 2 Winner
            </button>
            <button
              class="px-6 py-3 bg-green-500 text-white rounded-lg hover:bg-green-600"
              @click="crownChampion"
            >
              Crown Champion
            </button>
          </div>
        </div>
      </div>

      <div v-else class="flex flex-col items-center gap-6">
        <h1 class="text-3xl font-bold text-primary">Not in Tournament</h1>
        <button
          class="px-6 py-3 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors duration-300"
          @click="router.push('/')"
        >
          Return to Home
        </button>
      </div>
    </div>
  </div>
</template>


<script setup lang="ts">
import {onMounted, onUnmounted, ref} from 'vue'
import router from "@/router";
import { useWebSocketStore } from '@/stores/websocket'
import type {Tournament, WSMessage} from "@/utils/types";

import {useToast} from 'vue-toast-notification';
import 'vue-toast-notification/dist/theme-sugar.css';

const $toast = useToast();

const semifinal1Winner = ref<string | undefined>('TBD')
const semifinal2Winner = ref<string | undefined>('TBD')
const semifinal1Looser = ref<string | undefined>('TBD')
const semifinal2Looser = ref<string | undefined>('TBD')
const tournamentWinner = ref<string | undefined>('')
const thirdplace = ref<string | undefined>('')


const tournament = ref<Tournament | null>(null)
const myUsername = ref('')
const showBracket = ref(false)
const showStartButton = ref(false)

const tournamentId = router.currentRoute.value.params.id
const owner = ref<boolean>(false)
const webSocketStore = useWebSocketStore()

const fetchTournament = async () => {
  try {
    const response = await fetch(`/api/tournament/${tournamentId}`)
    tournament.value = await response.json()
    if(tournament.value?.users.length === 4)
      showStartButton.value = true
  }
  catch (error) {
    console.error('tournament: ' + error)
  }
}

const fetchMe = async () => {
  try {
    const response = await fetch('/api/session')
    const user = await response.json()
    myUsername.value = user.username
  }
  catch (error) {
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
  }
  catch (error) {
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
  }
  else if(message.type === 'tournament_new_user_enrolled') {
    await fetchTournament()
    $toast.info('Someone joined the tournament :)')
  }
  else if(message.type === 'tournament_user_left') {
    await fetchTournament()
    $toast.warning('Someone left the tournament :(')
  }
  else if(message.type === 'game_tournament_ready') {
    $toast.success('Tournament is starting!')
    showBracket.value = true
    // wait for 3 seconds, then router push to /game
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
        'Content-Type': 'application/json'
      }
    })
  }
  catch (error) {
    console.log(error)
  }
}

function exitTournament() {
  try {
    fetch(`/api/tournament/${tournamentId}`, {
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json'
      }
    })
    $toast.warning('You left this tournament')
    router.push('/')
  }
  catch (error) {
    console.error(error)
  }
}

function deleteTournament() {
  try {
    fetch(`/api/tournament/${tournamentId}/cancel`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      }
    })
  }
  catch (error) {
    console.error(error)
  }
}

function progressSemifinal1() {
  semifinal1Winner.value = tournament.value?.users[0] === semifinal1Winner.value
    ? tournament.value?.users[1]
    : tournament.value?.users[0]
  semifinal1Looser.value = tournament.value?.users[0] === semifinal1Winner.value
    ? tournament.value?.users[1]
    : tournament.value?.users[0]
}

function progressSemifinal2() {
  semifinal2Winner.value = tournament.value?.users[2] === semifinal2Winner.value
    ? tournament.value?.users[3]
    : tournament.value?.users[2]
  semifinal2Looser.value = tournament.value?.users[2] === semifinal2Winner.value
    ? tournament.value?.users[3]
    : tournament.value?.users[2]
}

function crownChampion() {
  if (semifinal1Winner.value && semifinal2Winner.value) {
    tournamentWinner.value = Math.random() > 0.5
      ? semifinal1Winner.value
      : semifinal2Winner.value
    thirdplace.value = Math.random() > 0.5
      ? semifinal1Looser.value
      : semifinal2Looser.value
  }
}
</script>


<style scoped>
/* Inherits retro styles from parent component */
.retro-box {
  background-color: #ffe5c9;
  border: 5px solid #8b4513;
  box-shadow:
    0 0 0 4px #d2691e,
    inset 0 0 20px rgba(0, 0, 0, 0.2);
  transition: transform 0.2s;
}

.retro-title {
  color: #ffd700;
  text-shadow:
    4px 4px 0 #8b4513,
    -1px -1px 0 #000,
    1px -1px 0 #000,
    -1px 1px 0 #000,
    1px 1px 0 #000;
  letter-spacing: 3px;
  animation: move-title 8s ease-in-out infinite alternate;
  border-bottom: 2px solid #8b4513;
}

.retro-button {
  @apply btn bg-primary text-white font-bold;
  border: 3px solid #8b4513;
  text-transform: uppercase;
  text-shadow: 2px 2px 0 rgba(0, 0, 0, 0.2);
  box-shadow: 0 2px 0 #8b4513;
  font-size: 1.1rem;
  height: 6vh;

  &.circle {
    width: 70px;
    height: 70px;
    border-radius: 50%;
  }

  &:hover {
    transform: translateY(2px);
    box-shadow:
      inset 0 0 10px rgba(0, 0, 0, 0.2),
      0 0px 0 #8b4513;
    cursor: url('/tortellino.png'), auto;
  }
}

@keyframes move-title {
  from {
    transform: rotate(-4deg);
  }
  to {
    transform: rotate(4deg);
  }
}


</style>

