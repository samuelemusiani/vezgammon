<template>
  <div class="flex h-full items-center justify-center w-full">
    <div class="flex h-[90%] w-[80%] flex-col items-center justify-center bg-base-100 border-8 border-primary rounded-md">
      <!-- Game Title -->
      <div class="mb-32 text-center">
        <h1 class="text-primary retro-title text-7xl mb-8 p-4">VezGammon</h1>
        <div class="text-accent">The Ultimate Backgammon Experience</div>
      </div>

      <!-- Button Container -->
      <div class="relative flex w-full max-w-4xl items-center justify-center">
        <!-- Left Button (Stats) -->
        <div class="absolute left-8">
          <button
            @click="(e: MouseEvent) => router.push('/wip')"
            @mouseenter="(e: MouseEvent) => play()"
            class="retro-button circle"
            title="Statistics"
          >
            <MedalIcon />
          </button>
        </div>

        <!-- Central Buttons -->
        <div class="flex w-full max-w-sm flex-col gap-6">
          <button
            @click="(e: MouseEvent) => openPlayModal()"
            @mouseenter="(e: MouseEvent) => play()"
            class="retro-button"
          >
            PLAY
          </button>
          <button
            @mouseenter="(e: MouseEvent) => play()"
            @click="router.push('/wip')"
            class="retro-button"
          >
            RULES
          </button>
          <button
            @mouseenter="(e: MouseEvent) => play()"
            @click="router.push('/wip')"
            class="retro-button"
          >
            SETTINGS
          </button>
        </div>

        <!-- Right Button (Profile) -->
        <div class="absolute right-8">
          <button
            @mouseenter="(e: MouseEvent) => play()"
            @click="navigateTo('/profile')"
            class="retro-button circle"
            title="Profile"
          >
            <ProfileIcon />
          </button>
        </div>
      </div>
    </div>

    <!-- Play Modal -->
    <dialog id="play_modal" class="modal">
      <div class="retro-box modal-box">
        <h3 class="retro-title mb-4 text-center text-2xl font-bold">
          {{ modalTitle }}
        </h3>
        <!-- Options -->
        <div class="flex flex-col gap-4">
          <template v-if="!showDifficulty">
            <button
              @mouseenter="(e: MouseEvent) => play()"
              @click="startGame('local')"
              class="retro-button"
            >
              Local Game (2 Players)
            </button>
            <button
              @mouseenter="(e: MouseEvent) => play()"
              @click="showAIDifficulty"
              class="retro-button"
            >
              Play vs AI
            </button>
            <button
              @mouseenter="(e: MouseEvent) => play()"
              @click="router.push('/wip')"
              class="retro-button"
            >
              Play Online
            </button>
          </template>

          <template v-else>
            <button
              @mouseenter="(e: MouseEvent) => play()"
              @click="startGameWithAI('easy')"
              class="retro-button"
            >
              Easy
            </button>
            <button
              @mouseenter="(e: MouseEvent) => play()"
              @click="startGameWithAI('medium')"
              class="retro-button"
            >
              Medium
            </button>
            <button
              @mouseenter="(e: MouseEvent) => play()"
              @click="startGameWithAI('hard')"
              class="retro-button"
            >
              Hard
            </button>
          </template>
        </div>

        <!-- Close button -->
        <div class="modal-action w-full">
          <form method="dialog" class="flex w-full justify-between">
            <button
              v-if="showDifficulty"
              @click="backToGameMode"
              class="retro-button"
            >
              Back
            </button>
            <button class="retro-button ml-auto">Close</button>
          </form>
        </div>
      </div>

      <form method="dialog" class="modal-backdrop">
        <button>close</button>
      </form>
    </dialog>
  </div>
</template>

<script setup lang="ts">
import MedalIcon from '@/utils/icons/MedalIcon.vue'
import ProfileIcon from '@/utils/icons/ProfileIcon.vue'
import router from '@/router'
import { useSound } from '@vueuse/sound'
import buttonSfx from '@/utils/sounds/button.mp3'
import { ref, computed } from 'vue'

const { play } = useSound(buttonSfx, { volume: 0.3 })
const showDifficulty = ref(false)

const modalTitle = computed(() => {
  return showDifficulty.value ? 'Choose Difficulty' : 'Select Game Mode'
})

const showAIDifficulty = () => {
  showDifficulty.value = true
}

const backToGameMode = () => {
  showDifficulty.value = false
}

const startGameWithAI = async (difficulty: 'easy' | 'medium' | 'hard') => {
  const modal = document.getElementById('play_modal') as HTMLDialogElement
  modal.close()
  showDifficulty.value = false

  try {
    await fetch(`/api/play/bot/${difficulty}`)
    router.push('/game')
  } catch (error) {
    console.error('Error starting game with AI:', error)
  }
}

const navigateTo = (path: string) => {
  router.push(path)
}

const openPlayModal = () => {
  const modal = document.getElementById('play_modal') as HTMLDialogElement
  modal.showModal()
}

const startGame = async (mode: 'local' | 'ai' | 'online') => {
  const modal = document.getElementById('play_modal') as HTMLDialogElement
  modal.close()

  switch (mode) {
    case 'local':
      await fetch('/api/play/local')
      router.push('/game')
      break
    case 'online':
      router.push('/game')
      break
  }
}
</script>

<style scoped>
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



.retro-box {
  background-color: #ffe5c9;
  border: 5px solid #8b4513;
  box-shadow:
    0 0 0 4px #d2691e,
    inset 0 0 20px rgba(0, 0, 0, 0.2);
}

.retro-button {
  @apply btn bg-primary text-neutral;
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
