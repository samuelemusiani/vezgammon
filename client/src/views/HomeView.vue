<template>
  <div class="retro-background">
    <div class="flex max-h-full w-full flex-col items-center justify-center">
      <!-- Game Title -->
      <div class="mb-32 text-center">
        <h1 class="retro-title text-7xl">VezGammon</h1>
        <div class="retro-subtitle">The Ultimate Backgammon Experience</div>
      </div>

      <!-- Button Container -->
      <div class="relative flex w-full max-w-4xl items-center justify-center">
        <!-- Left Button (Stats) -->
        <div class="absolute left-8">
          <button class="retro-button circle" title="Statistics">
            <MedalIcon />
          </button>
        </div>

        <!-- Central Buttons -->
        <div class="flex w-full max-w-sm flex-col gap-6">
          <button @click="openPlayModal" class="retro-button">PLAY</button>
          <button class="retro-button">RULES</button>
          <button class="retro-button">SETTINGS</button>
        </div>

        <!-- Right Button (Profile) -->
        <div class="absolute right-8">
          <button
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
          Select Game Mode
        </h3>
        <!-- Options -->
        <div class="flex flex-col gap-4">
          <button @click="startGame('local')" class="retro-button">
            Local Game (2 Players)
          </button>
          <button @click="startGame('ai')" class="retro-button">
            Play vs AI
          </button>
          <button @click="startGame('online')" class="retro-button">
            Online Match
          </button>
        </div>

        <!-- Close button -->
        <div class="modal-action">
          <form method="dialog">
            <button class="retro-button">Close</button>
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

const navigateTo = (path: string) => {
  router.push(path)
}

const openPlayModal = () => {
  const modal = document.getElementById('play_modal') as HTMLDialogElement
  modal.showModal()
}

const startGame = (mode: 'local' | 'ai' | 'online') => {
  const modal = document.getElementById('play_modal') as HTMLDialogElement
  modal.close()

  switch (mode) {
    case 'local':
      router.push('/game')
      break
    case 'ai':
      router.push('/game')
      break
    case 'online':
      router.push('/game')
      break
  }
}
</script>

<style scoped>
.retro-background {
  @apply flex min-h-screen items-center justify-center;
  background: #2c1810;
  background-image: repeating-linear-gradient(
      45deg,
      rgba(139, 69, 19, 0.1) 0px,
      rgba(139, 69, 19, 0.1) 2px,
      transparent 2px,
      transparent 10px
    ),
    repeating-linear-gradient(
      -45deg,
      rgba(139, 69, 19, 0.1) 0px,
      rgba(139, 69, 19, 0.1) 2px,
      transparent 2px,
      transparent 10px
    );
  cursor: url('/tortellino.png'), auto;
  border: 6px solid #d2691e;
}

.retro-title {
  font-family: 'Arial Black', serif;
  color: #ffd700;
  text-shadow:
    4px 4px 0 #8b4513,
    -1px -1px 0 #000,
    1px -1px 0 #000,
    -1px 1px 0 #000,
    1px 1px 0 #000;
  letter-spacing: 3px;
  animation: move-title 5s ease-in-out infinite alternate;
  padding-bottom: 10px;
  margin-bottom: 20px;
  border-bottom: 2px solid #8b4513;
}

.retro-subtitle {
  font-family: 'Arial Black', serif;
  color: #d2691e;
  font-size: 1.2rem;
  text-shadow: 2px 2px 0 rgba(0, 0, 0, 0.5);
  letter-spacing: 1px;
}

.retro-box {
  background-color: #ffe5c9;
  border: 5px solid #8b4513;
  box-shadow:
    0 0 0 4px #d2691e,
    inset 0 0 20px rgba(0, 0, 0, 0.2);
}

.retro-button {
  @apply btn;
  background: #d2691e;
  color: white;
  border: 3px solid #8b4513;
  font-family: 'Arial Black', serif;
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
    transform: rotate(-2deg);
  }
  to {
    transform: rotate(0deg);
  }
}
</style>
