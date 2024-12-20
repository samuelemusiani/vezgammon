<template>
  <div class="flex h-full w-full items-center justify-center">
    <div class="retro-box p-12 text-center">
      <h2 class="retro-title mb-6 text-4xl">Game Invitation</h2>
      <p class="mb-8 text-xl">
        You have been invited to play a match of VezGammon!
      </p>
      <div class="flex justify-center gap-4">
        <button @click="acceptInvite" class="retro-button">Accept</button>
        <button @click="declineInvite" class="retro-button">Decline</button>
      </div>
    </div>

    <!-- Error Modal -->
    <dialog id="error_modal" class="modal">
      <div class="retro-box modal-box">
        <h3 class="retro-title mb-4 text-center text-3xl font-bold">
          Invalid Link
        </h3>
        <p class="text-center text-lg">{{ errorMessage }}</p>
        <div class="modal-action flex justify-center">
          <form method="dialog">
            <button @click="() => router.push('/')" class="retro-button">
              Go Home
            </button>
          </form>
        </div>
      </div>
    </dialog>
  </div>
</template>

<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router'
import { ref } from 'vue'
import { vfetch } from '@/utils/fetch'

const route = useRoute()
const router = useRouter()
const errorMessage = ref('')

const showError = (message: string) => {
  errorMessage.value = message
  const errorModal = document.getElementById('error_modal') as HTMLDialogElement
  errorModal.showModal()
}

const acceptInvite = async () => {
  try {
    const code = route.params.code
    const response = await vfetch(`/api/play/invite/${code}`)

    if (!response.ok) {
      showError(
        'Invalid or expired invitation link. Ask your friend to send a new one.',
      )
      return
    }

    router.push('/game')
  } catch (error) {
    console.error('Error accepting invite:', error)
    showError('Something went wrong. Please try again.')
  }
}

const declineInvite = () => {
  router.push('/')
}
</script>
