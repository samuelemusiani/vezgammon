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
          class="retro-box p-6 relative"
          @mouseenter="play()"
        >
          <div class="flex justify-between items-center">
            <div>
              <h2 class="text-2xl font-bold text-[#8b4513]">{{ tournament.name }}</h2>
              <p class="text-[#d2691e]">
                {{ tournament.startDate }} | Owner: {{ tournament.owner }}
              </p>
            </div>
            <div class="flex flex-row">
              <div v-for="i in tournament.participants" :key="i">
                <svg class="w-6 h-6 text-gray-800" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="currentColor" viewBox="0 0 24 24">
                  <path fill-rule="evenodd" d="M12 20a7.966 7.966 0 0 1-5.002-1.756l.002.001v-.683c0-1.794 1.492-3.25 3.333-3.25h3.334c1.84 0 3.333 1.456 3.333 3.25v.683A7.966 7.966 0 0 1 12 20ZM2 12C2 6.477 6.477 2 12 2s10 4.477 10 10c0 5.5-4.44 9.963-9.932 10h-.138C6.438 21.962 2 17.5 2 12Zm10-5c-1.84 0-3.333 1.455-3.333 3.25S10.159 13.5 12 13.5c1.84 0 3.333-1.455 3.333-3.25S13.841 7 12 7Z" clip-rule="evenodd"/>
                </svg>

              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
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

.retro-box:hover {
  transform: scale(1.02);
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
