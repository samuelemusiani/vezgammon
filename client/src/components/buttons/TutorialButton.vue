<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import '@fortawesome/fontawesome-free/css/all.min.css'

// Modal visibility state
const isModalOpen = ref(false)

// Stepper control states
const currentStep = ref(1)
const totalSteps = 4

// Function to navigate steps
const goToNextStep = () => {
  if (currentStep.value < totalSteps) currentStep.value++
}

const goToPreviousStep = () => {
  if (currentStep.value > 1) currentStep.value--
}

// Function to close modal
const closeModal = () => {
  isModalOpen.value = false
  currentStep.value = 1 // Reset stepper when closing
}

const handleKeydown = (event: KeyboardEvent) => {
  if (event.key === 'Escape') closeModal()
}

onMounted(() => {
  document.addEventListener('keydown', handleKeydown)
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleKeydown)
})
</script>
<template>
  <div class="fixed bottom-4 left-4 z-50">
    <!-- Floating Button -->
    <button
      class="btn btn-circle btn-primary btn-lg"
      @click="isModalOpen = true"
    >
      <i class="fas fa-info"></i>
    </button>

    <!-- Modal -->
    <dialog v-if="isModalOpen" class="modal modal-open">
      <div
        class="modal-box w-full max-w-[60%] border-2 border-primary bg-base-200"
      >
        <!-- Modal Headers -->
        <div class="mb-4 flex items-center justify-between">
          <h3 class="text-xl font-bold">Tutorial Videos</h3>
          <button class="btn btn-circle btn-ghost btn-sm" @click="closeModal">
            <i class="fas fa-times"></i>
          </button>
        </div>

        <!-- Stepper -->
        <div class="mb-6 flex-grow">
          <div class="mb-4 flex justify-between">
            <span
              v-for="step in totalSteps"
              :key="step"
              class="flex h-8 w-8 items-center justify-center rounded-full"
              :class="{
                'bg-primary text-primary-content': currentStep === step,
                'bg-base-300 text-base-content': currentStep !== step,
              }"
            >
              {{ step }}
            </span>
          </div>

          <!-- Step Content -->
          <div v-if="currentStep === 1">
            <img
              class="aspect-video w-full rounded-lg object-cover"
              src="/tutorial-dice.gif"
              alt="tutorial-dice"
            />
            <p class="mt-2 text-xl">Step 1: Roll the dice!</p>
          </div>
          <div v-else-if="currentStep === 2">
            <img
              class="aspect-video w-full rounded-lg object-cover"
              src="/tutorial-move.gif"
              alt="tutorial-move"
            />
            <p class="mt-2 text-xl">Step 2: move the checkers!</p>
          </div>
          <div v-else-if="currentStep === 3">
            <img
              class="aspect-video w-full rounded-lg object-cover"
              src="/tutorial-exit.gif"
              alt="tutorial-exit"
            />
            <p class="mt-2 text-xl">Step 3: Exiting</p>
          </div>
          <div
            v-else-if="currentStep === 4"
            class="flex min-h-[400px] w-full flex-col items-center justify-center"
          >
            <h1 class="text-center text-2xl font-bold">Tutorial Completed!</h1>
            <h2 class="mb-5 text-center text-lg font-semibold opacity-75">
              Go win your first match!
            </h2>
          </div>
        </div>

        <!-- Step Navigation -->
        <div class="modal-action">
          <button
            :disabled="currentStep === 1"
            class="btn btn-neutral"
            @click="goToPreviousStep"
          >
            I'll go back
          </button>
          <button
            v-if="currentStep < totalSteps"
            class="btn btn-primary"
            @click="goToNextStep"
          >
            Understood
          </button>
          <button
            v-if="currentStep === totalSteps"
            class="btn btn-success"
            @click="closeModal"
          >
            Finish
          </button>
        </div>
      </div>

      <form method="dialog" class="modal-backdrop" @click="closeModal">
        <button>close</button>
      </form>
    </dialog>
  </div>
</template>
