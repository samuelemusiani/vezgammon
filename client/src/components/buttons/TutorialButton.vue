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
      class="retro-button circle rounded-full bg-primary"
      @click="isModalOpen = true"
    >
      <i class="fas fa-info"></i>
    </button>

    <!-- Modal -->
    <div
      v-if="isModalOpen"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50"
    >
      <div
        class="retro-box relative flex min-h-[50%] w-full max-w-[60%] flex-col rounded-lg bg-white p-6"
      >
        <!-- Modal Headers -->
        <div class="mb-4 flex items-center justify-between">
          <h3 class="text-xl font-bold text-gray-700">Tutorial Videos</h3>
          <!-- Close Button -->
          <button class="text-gray-400 hover:text-gray-600" @click="closeModal">
            <i class="fas fa-times fa-lg"></i>
          </button>
        </div>

        <!-- Stepper -->
        <div class="mb-6 flex-grow">
          <div class="mb-4 flex justify-between">
            <span
              v-for="step in totalSteps"
              :key="step"
              class="flex h-8 w-8 items-center justify-center rounded-full border-2"
              :class="{
                'border-primary bg-primary text-white': currentStep === step,
                'border-gray-300 bg-gray-200 text-gray-500':
                  currentStep !== step,
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
            <p class="text-xl">Step 1: Roll the dice!</p>
          </div>
          <div v-else-if="currentStep === 2">
            <img
              class="aspect-video w-full rounded-lg object-cover"
              src="/tutorial-move.gif"
              alt="tutorial-move"
            />
            <p class="text-xl">Step 2: move the checkers!</p>
          </div>
          <div v-else-if="currentStep === 3">
            <img
              class="aspect-video w-full rounded-lg object-cover"
              src="/tutorial-exit.gif"
              alt="tutorial-exit"
            />
            <p class="text-xl">Step 3: Exiting</p>
          </div>
          <div
            v-else-if="currentStep === 4"
            class="flex min-h-[400px] w-full flex-col items-center justify-center"
          >
            <h1 class="text-center text-2xl font-bold text-gray-700">
              Tutorial Completed!
            </h1>
            <h2 class="mb-5 text-center text-lg font-semibold text-gray-600">
              Go win your first match!
            </h2>
          </div>
        </div>

        <!-- Step Navigation -->
        <div class="mt-auto flex w-full justify-between">
          <button
            :disabled="currentStep === 1"
            class="retro-button bg-gray-300"
            @click="goToPreviousStep"
          >
            I'll go back
          </button>
          <button
            v-if="currentStep < totalSteps"
            class="retro-button bg-primary"
            @click="goToNextStep"
          >
            Understood
          </button>
          <button
            v-if="currentStep === totalSteps"
            class="retro-button bg-green-500"
            @click="closeModal"
          >
            Finish
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
<style scoped>
.retro-button {
  @apply btn btn-primary btn-lg border-4 border-accent text-white;
  text-transform: uppercase;
  text-shadow: 2px 2px 0 rgba(0, 0, 0, 0.2);
  box-shadow: 0 2px 0 #8b4513;
  font-size: 1.1rem;
  height: 6vh;

  &.circle {
    width: 60px;
    height: 60px;
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
</style>
