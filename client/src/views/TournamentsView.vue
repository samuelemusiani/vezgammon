<template>
  <div class="flex h-full w-full items-center justify-center">
      <div class="flex h-[90%] w-[80%] flex-col items-center justify-center rounded-md border-8 border-primary bg-base-100">
      <!-- Page Title -->
      <div class="flex flex-col gap-4 text-center mb-4">
        <h1 class="retro-title text-primary text-5xl">Tournaments</h1>
      </div>

      <!-- Tournaments List -->
      <div class="w-full max-w-4xl overflow-y-auto max-h-[calc(100vh-300px)] space-y-6 p-8 no-scrollbar">
        <div
          v-for="tournament in tournaments"
          :key="tournament.id"
          class="retro-box p-6 relative hover:scale-[1.02]"
          @mouseenter="play()"
          @click="openTournamentModal(tournament)"
        >
          <div class="flex justify-between items-center">
            <div>
              <h2 class="text-2xl font-bold text-[#8b4513]">{{ tournament.name }}</h2>
              <p class="text-[#d2691e]">
                {{ tournament.startDate }} | Owner: {{ tournament.owner }}
              </p>
            </div>
            <div class="text-primary text-xl">
              {{ tournament.participants }}/4
            </div>
          </div>
        </div>
      </div>
    </div>
    <dialog id="select_tournament" class="modal">
      <div class="retro-box modal-box">
        <h3 class="retro-title mb-4 text-center text-2xl font-bold">
          {{ selectedTournament?.name }}
        </h3>
        <!-- Options -->
        <div class="flex flex-col gap-4">
          <div class="text-center">
            {{ selectedTournament?.startDate }} | Owner: {{ selectedTournament?.owner }}
          </div>
          <div class="flex flex-row justify-center gap-8 items-center">
            <div v-for="i in selectedTournament?.participants" :key="i">
              <div class="h-16 w-16 overflow-hidden rounded-full bg-gray-200 border-primary border-2 hover:scale-[1.02]">
                <img
                  :src="`https://api.dicebear.com/6.x/avataaars/svg?seed=${users[i]}`"
                  alt="Opponent avatar"
                  class="h-full w-full object-cover"
                />
              </div>
              <div class="text-center text-[#8b4513] font-bold">{{ users[i-1] }}</div>
            </div>
          </div>
          <div class="flex w-full justify-between">
            <button @mouseenter="play" @click="closeTournamentModal" class="retro-button">
              CLOSE
            </button>
            <button @mouseenter="play" @click="joinTournament(selectedTournament.id)" class="retro-button">
              JOIN TOURNAMENT
            </button>
          </div>
        </div>
      </div>
          </dialog>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useSound } from '@vueuse/sound'
import buttonSfx from '@/utils/sounds/button.mp3'

const { play } = useSound(buttonSfx, { volume: 0.3 })

const tournaments = ref([
  {
    id: 1,
    name: 'Beginners Challenge',
    startDate: 'December 1st, 2024 - 3:00 PM',
    owner: 'Omar',
    participants: 3,
  },
  {
    id: 2,
    name: 'Pro Masters Tournament',
    startDate: 'January 15th, 2025 - 6:00 PM',
    owner: 'Lele',
    participants: 2,
  },
  {
    id: 3,
    name: 'Global Backgammon Championship',
    startDate: 'March 20th, 2025 - 8:00 PM',
    owner: 'Samu',
    participants: 1,
  },
  {
    id: 4,
    name: 'Retro Gaming Tournament',
    startDate: 'April 5th, 2025 - 10:00 AM',
    owner: 'Lollo',
    participants: 4,
  },
  {
    id: 5,
    name: 'Retro Gaming Tournament',
    startDate: 'April 5th, 2025 - 10:00 AM',
    owner: 'Fabio',
    participants: 2,
  }
])

tournaments.value = tournaments.value.filter(tournament => tournament.participants < 4)

const users = [ 'Omar', 'Lele', 'Samu' ]
// Selected tournament for modal
const selectedTournament = ref(null)

// Open modal with selected tournament
const openTournamentModal = (tournament) => {
  selectedTournament.value = tournament
  document.getElementById('select_tournament').showModal()
}

// Close modal
const closeTournamentModal = () => {
  document.getElementById('select_tournament').close()
  selectedTournament.value = null
}

const joinTournament = async (tournamentId: number) => {
  try {
    await fetch(`/api/tournaments/${tournamentId}/join`)
    // Additional logic for tournament registration
  } catch (error) {
    console.error('Error joining tournament:', error)
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

  /* Hide scrollbar for Chrome, Safari and Opera */
  .no-scrollbar::-webkit-scrollbar {
    display: none;
  }
  /* Hide scrollbar for IE, Edge and Firefox */
  .no-scrollbar {
    -ms-overflow-style: none;  /* IE and Edge */
    scrollbar-width: none;  /* Firefox */
  }

</style>
