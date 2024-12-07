<template>
  <div class="flex h-full w-full items-center justify-center">
    <div
      class="flex h-[90%] w-[80%] flex-col items-center justify-center rounded-md border-8 border-primary bg-base-100"
    >
      <!-- Page Title -->
      <div class="mb-4 flex flex-col gap-4 text-center">
        <h1 class="retro-title text-5xl">Tournaments</h1>
      </div>

      <!-- Tournaments List -->
      <div
        class="no-scrollbar max-h-[calc(100vh-300px)] w-full max-w-4xl space-y-6 overflow-y-auto p-8"
      >
        <div
          v-for="tournament in tournaments"
          :key="tournament.id"
          class="retro-box relative p-6 hover:scale-[1.02]"
          @mouseenter="play()"
          @click="openTournamentModal(tournament)"
        >
          <div class="flex items-center justify-between">
            <div>
              <h2 class="text-2xl font-bold text-[#8b4513]">
                {{ tournament.name }}
              </h2>
              <p class="text-accent">
                {{ dateformatter.format(new Date(tournament.creation_date)) }} |
                Owner: {{ tournament.owner }}
              </p>
            </div>
            <div class="text-xl font-semibold text-accent text-primary">
              {{ tournament.user_number }}/4
            </div>
          </div>
        </div>
      </div>
    </div>
    <dialog id="select_tournament" class="modal">
      <div class="retro-box modal-box" v-if="selectedTournament">
        <h3 class="retro-title mb-4 text-center text-2xl font-bold">
          {{ selectedTournament?.name }}
        </h3>
        <!-- Options -->
        <div class="flex flex-col gap-4">
          <div class="text-center font-semibold text-accent">
            {{
              selectedTournament
                ? dateformatter.format(
                    new Date(selectedTournament.creation_date),
                  )
                : ''
            }}
            | Owner: {{ selectedTournament?.owner }}
          </div>
          <div class="flex flex-row items-center justify-center gap-8">
            <!-- href to user profile in the div below -->
            <div
              v-for="(user, index) in selectedTournament?.users"
              :key="index"
            >
              <div
                class="h-16 w-16 overflow-hidden rounded-full border-2 border-primary bg-gray-200 hover:scale-[1.02]"
              >
                <img
                  :src="avatars[index]"
                  alt="Opponent avatar"
                  class="h-full w-full object-cover"
                />
              </div>
              <div class="text-center font-bold text-[#8b4513]">{{ user }}</div>
            </div>
          </div>
          <div class="flex w-full justify-between">
            <button
              @mouseenter="(e: MouseEvent) => play()"
              @click="closeTournamentModal"
              class="retro-button"
            >
              CLOSE
            </button>
            <button
              @mouseenter="(e: MouseEvent) => play()"
              v-if="selectedTournament?.users.includes(myUsername)"
              @click="joinTournament(selectedTournament.id)"
              class="retro-button"
            >
              SEE TOURNAMENT
            </button>
            <button
              @mouseenter="(e: MouseEvent) => play()"
              v-else
              @click="joinTournament(selectedTournament.id)"
              class="retro-button"
            >
              JOIN TOURNAMENT
            </button>
          </div>
        </div>
      </div>
    </dialog>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useSound } from '@vueuse/sound'
import buttonSfx from '@/utils/sounds/button.mp3'
import router from '@/router'
import type { Tournament } from '@/utils/types'

import { useToast } from 'vue-toast-notification'
import 'vue-toast-notification/dist/theme-sugar.css'

const $toast = useToast()

const { play } = useSound(buttonSfx, { volume: 0.3 })

interface SimpleTournament {
  id: number
  name: string
  creation_date: string
  owner: string
  user_number: number
}

const tournaments = ref<SimpleTournament[] | null>(null)
const myUsername = ref('')

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
  try {
    const response = await fetch('/api/tournament/list')
    tournaments.value = await response.json()
    console.log('Tournaments:', tournaments.value)
  } catch (error) {
    console.error('Error fetching tournaments:', error)
  }
  await fetchMe()
})

// Selected tournament for modal
const selectedTournament = ref<Tournament | null>(null)
const avatars = ref<string[]>([])

// Open modal with selected tournament
const openTournamentModal = async (tournament: SimpleTournament) => {
  const data = await fetch(`/api/tournament/${tournament.id}`)
  selectedTournament.value = await data.json()
  if (selectedTournament.value?.users) {
    for (const user of selectedTournament.value.users) {
      const res = await fetch(`/api/player/${user}/avatar`)
      const avatar = await res.json()
      avatars.value.push(avatar)
    }
  }
  const el = document.getElementById('select_tournament') as HTMLDialogElement
  el.showModal()
}

// Close modal
const closeTournamentModal = async () => {
  const response = await fetch('/api/tournament/list')
  tournaments.value = await response.json()
  const el = document.getElementById('select_tournament') as HTMLDialogElement
  el.close()
  selectedTournament.value = null
}

const joinTournament = async (tournamentId: number) => {
  if (selectedTournament.value?.users.includes(myUsername.value)) {
    await router.push('/tournaments/' + tournamentId)
    return
  }
  try {
    const response = await fetch(`/api/tournament/${tournamentId}`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
    })
    if (response.ok) {
      await router.push('/tournaments/' + tournamentId)
      $toast.info('You joined a tournament!')
    } else {
      console.error('Error joining tournament:', response)
    }
  } catch (error) {
    console.error('Error joining tournament:', error)
  }
}

const dateformatter = new Intl.DateTimeFormat('it-IT', {
  hour12: false,
  dateStyle: 'medium',
  timeStyle: 'short',
})
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
