<template>
  <div class="flex h-full w-full items-center justify-center">
    <div
      class="flex h-[90%] w-[80%] flex-col items-center justify-center rounded-md border-8 border-primary bg-base-100"
    >
      <!-- Game Title -->
      <div class="mb-16 text-center">
        <h1 class="retro-title mb-8 p-4 text-7xl font-bold">VezGammon</h1>
        <div class="text-xl font-bold text-accent">
          The Ultimate Backgammon Experience
        </div>
      </div>

      <!-- Button Container -->
      <div class="relative flex w-full max-w-4xl items-center justify-center">
        <!-- Left Button (Stats) -->
        <div class="absolute left-8">
          <button
            @click="navigateTo('/stats')"
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
            class="retro-button"
            @mouseenter="(e: MouseEvent) => play()"
            @click="router.push('/leaderboard')"
          >
            LEADERBOARD
          </button>
          <button
            @mouseenter="(e: MouseEvent) => play()"
            @click="openRulesModal"
            class="retro-button font-bold"
          >
            RULES
          </button>
          <button
            @mouseenter="(e: MouseEvent) => play()"
            @click="openSettingsModal"
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
          <template v-if="modals === 0">
            <button
              @mouseenter="(e: MouseEvent) => play()"
              @click="startLocalGame"
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
              @click="showOnlineMenu"
              class="retro-button"
            >
              Play Online
            </button>
            <button
              @mouseenter="(e: MouseEvent) => play()"
              @click="playTutorial"
              class="retro-button"
            >
              Play Tutorial
            </button>
          </template>

          <template v-else-if="modals === 1">
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

          <template v-else-if="modals === 2">
            <button
              @mouseenter="(e: MouseEvent) => play()"
              @click="startRandomGame"
              class="retro-button"
            >
              Random Opponent
            </button>
            <button
              @mouseenter="(e: MouseEvent) => play()"
              @click="createInviteLink"
              class="retro-button"
            >
              Invite Friend
            </button>
            <div v-if="inviteLink" class="mt-4">
              <div class="flex items-center gap-2 rounded bg-base-200 p-2">
                <input
                  type="text"
                  :value="inviteLink"
                  class="w-full bg-transparent p-2"
                  readonly
                />

                <button
                  @click="copyInviteLink"
                  class="retro-button px-4"
                  :class="{ 'bg-success': linkCopied }"
                >
                  {{ linkCopied ? 'Copied!' : 'Copy' }}
                </button>
              </div>
            </div>
            <div v-if="inviteLink" class="flex justify-center gap-2">
              <TelegramShareButton
                :url="inviteLink"
                title="Do you want to play with me? Join me on VezGammon!"
              />

              <WhatsappShareButton
                :url="inviteLink"
                title="Do you want to play with me? Join me on VezGammon!"
              />
            </div>
          </template>

          <template v-else-if="modals === 3">
            <div class="flex flex-row justify-between gap-2">
              <input
                v-model="tourn_name"
                type="text"
                class="input flex-grow border-2 border-primary"
                placeholder="Tournament name"
              />
              <button
                @mouseenter="(e: MouseEvent) => play()"
                @click="create_tourn"
                class="retro-button"
              >
                New
              </button>
            </div>
            <button
              @mouseenter="(e: MouseEvent) => play()"
              @click="router.push('/tournaments')"
              class="retro-button"
            >
              Join
            </button>
          </template>
        </div>

        <!-- Close button -->
        <div class="modal-action w-full">
          <form method="dialog" class="flex w-full justify-between">
            <button
              v-if="modals !== 0"
              @click="backToGameMode"
              class="retro-button"
            >
              Back
            </button>
            <button class="retro-button ml-auto" @click="backToGameMode">
              Close
            </button>
          </form>
        </div>
      </div>

      <form method="dialog" class="modal-backdrop">
        <button>close</button>
      </form>
    </dialog>

    <dialog id="waiting_modal" class="modal">
      <div class="retro-box modal-box text-center">
        <h3 class="retro-title mb-4 text-2xl font-bold">
          Waiting for Opponent
        </h3>
        <div class="flex flex-col items-center gap-4">
          <div class="loading loading-spinner loading-lg"></div>
          <p class="text-lg">Searching for an opponent...</p>
        </div>
        <div class="modal-action">
          <form method="dialog">
            <button @click="handleCancelMatchmaking" class="retro-button">
              Cancel
            </button>
          </form>
        </div>
      </div>
    </dialog>

    <!-- Resume Modal -->
    <dialog id="resume_game_modal" class="modal">
      <div class="retro-box modal-box">
        <h3 class="retro-title mb-4 text-center text-2xl font-bold">
          Game in Progress
        </h3>
        <p class="mb-4 text-center">You have an ongoing game.</p>
        <div class="flex flex-col gap-4">
          <button
            @mouseenter="(e: MouseEvent) => play()"
            @click="resumeGame"
            class="retro-button"
          >
            Resume Game
          </button>
          <button
            @mouseenter="(e: MouseEvent) => play()"
            @click="leaveGame"
            class="retro-button"
          >
            Leave Game
          </button>
        </div>
      </div>

      <form method="dialog" class="modal-backdrop">
        <button>close</button>
      </form>
    </dialog>

    <SettingsModal />

    <dialog id="rules_modal" class="modal">
      <div
        class="modal-box max-h-[85vh] max-w-3xl overflow-y-auto border-4 border-primary"
      >
        <RulesSection />

        <div class="modal-action">
          <form method="dialog">
            <button class="retro-button">Close</button>
          </form>
        </div>
      </div>
    </dialog>
  </div>
</template>

<script setup lang="ts">
import MedalIcon from '@/utils/icons/MedalIcon.vue'
import ProfileIcon from '@/utils/icons/ProfileIcon.vue'
import RulesSection from '@/components/RulesSection.vue'
import WhatsappShareButton from '@/components/buttons/WhatsappShare.vue'
import TelegramShareButton from '@/components/buttons/TelegramShare.vue'
import SettingsModal from '@/components/SettingsModal.vue'

import router from '@/router'
import { useSound } from '@vueuse/sound'
import buttonSfx from '@/utils/sounds/button.mp3'
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useWebSocketStore } from '@/stores/websocket'
import { useAudioStore } from '@/stores/audio'
import type { WSMessage } from '@/utils/types'

const { play: playSound } = useSound(buttonSfx, { volume: 0.3 })
const webSocketStore = useWebSocketStore()
// 0 for base, 1 for bot difficulty, 2 for online options, 3 for tournaments options,
const modals = ref(0)
const inviteLink = ref('')
const linkCopied = ref(false)
const audioStore = useAudioStore()

onMounted(() => {
  webSocketStore.connect()
  webSocketStore.addMessageHandler(handleMatchmaking)
  checkIfInGame()
})

onUnmounted(() => {
  webSocketStore.removeMessageHandler(handleMatchmaking)
})

const handleMatchmaking = (message: WSMessage) => {
  if (message.type === 'game_found') {
    const waitingModal = document.getElementById(
      'waiting_modal',
    ) as HTMLDialogElement
    waitingModal.close()
    router.push('/game')
  }
}

const play = () => {
  if (audioStore.isAudioEnabled) {
    playSound()
  }
}

const showOnlineOptions = () => {
  const playModal = document.getElementById('play_modal') as HTMLDialogElement
  playModal.close()
  const onlineModal = document.getElementById(
    'online_options_modal',
  ) as HTMLDialogElement
  onlineModal.showModal()
}

const startRandomGame = () => {
  const modal = document.getElementById('play_modal') as HTMLDialogElement
  modal.close()
  startOnlineGame()
}

const createInviteLink = async () => {
  try {
    const response = await fetch('/api/play/invite')
    const data = await response.json()
    inviteLink.value = `${window.location.origin}/invite/${data.Link}`
    linkCopied.value = false
  } catch (error) {
    console.error('Error creating invite link:', error)
  }
}

const copyInviteLink = async () => {
  try {
    await navigator.clipboard.writeText(inviteLink.value)
    linkCopied.value = true
    setTimeout(() => {
      linkCopied.value = false
    }, 2000)
  } catch (error) {
    console.error('Error copying to clipboard:', error)
  }
}

const checkIfInGame = async () => {
  const response = await fetch('/api/play')
  if (response.ok) {
    const resumeModal = document.getElementById(
      'resume_game_modal',
    ) as HTMLDialogElement
    resumeModal.showModal()
  }
}

const resumeGame = () => {
  const modal = document.getElementById(
    'resume_game_modal',
  ) as HTMLDialogElement
  modal.close()
  router.push('/game')
}

const leaveGame = async () => {
  try {
    await fetch('/api/play', { method: 'DELETE' })
    const modal = document.getElementById(
      'resume_game_modal',
    ) as HTMLDialogElement
    modal.close()
  } catch (error) {
    console.error('Error leaving game:', error)
  }
}

const handleCancelMatchmaking = async () => {
  await fetch('/api/play/search', { method: 'DELETE' })
  const waitingModal = document.getElementById(
    'waiting_modal',
  ) as HTMLDialogElement
  waitingModal.close()
}

const openSettingsModal = () => {
  const modal = document.getElementById('settings_modal') as HTMLDialogElement
  modal.showModal()
}

const modalTitle = computed(() => {
  switch (modals.value) {
    case 0:
      return 'Select Game Mode'
    case 1:
      return 'Select AI Difficulty'
    case 2:
      return 'Play Online'
    case 3:
      return 'Tournaments'
  }
})

const showAIDifficulty = () => {
  modals.value = 1
}

const showOnlineMenu = () => {
  modals.value = 2
}

const showTournamentMenu = () => {
  modals.value = 3
}

const backToGameMode = () => {
  modals.value = 0
  inviteLink.value = ''
}

const startGameWithAI = async (
  difficulty: 'easy' | 'medium' | 'hard',
  variant: null | string = null,
) => {
  const modal = document.getElementById('play_modal') as HTMLDialogElement
  modal.close()
  modals.value = 0

  try {
    await fetch(`/api/play/bot/${difficulty}`)
    const destination = variant ? `/game?variant=${variant}` : '/game'
    router.push(destination)
  } catch (error) {
    console.error('Error starting game with AI:', error)
  }
}

const startOnlineGame = async () => {
  try {
    const modal = document.getElementById('play_modal') as HTMLDialogElement
    modal.close()

    const waitingModal = document.getElementById(
      'waiting_modal',
    ) as HTMLDialogElement
    waitingModal.showModal()

    await fetch('/api/play/search')
  } catch (error) {
    console.error('Error starting online game:', error)
    // In caso di errore, chiudi il modale di attesa
    const waitingModal = document.getElementById(
      'waiting_modal',
    ) as HTMLDialogElement
    waitingModal.close()
  }
}

const navigateTo = (path: string) => {
  router.push(path)
}

const openPlayModal = () => {
  const modal = document.getElementById('play_modal') as HTMLDialogElement
  modal.showModal()
}

const playTutorial = () => {
  const modal = document.getElementById('play_modal') as HTMLDialogElement
  modal.close()
  // launch a easy game with bot and show tutorial variant
  startGameWithAI('easy', 'tutorial')
}

const startLocalGame = async () => {
  const modal = document.getElementById('play_modal') as HTMLDialogElement
  modal.close()

  await fetch('/api/play/local')
  router.push('/game')
}

const openRulesModal = () => {
  const modal = document.getElementById('rules_modal') as HTMLDialogElement
  modal.showModal()
}

const tourn_name = ref('')

function create_tourn() {
  fetch('/api/tournament/create', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ name: tourn_name.value }),
  })
    .then(res => res.json())
    .then(data => {
      const id = data.id
      console.log(data)
      router.push('/tournaments/' + id)
    })
    .catch(err => console.error(err))
}
</script>
