<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import '@fortawesome/fontawesome-free/css/all.min.css'

// Modal visibility state
const isModalOpen = ref(false)

// Stepper control states
const currentStep = ref(1)
const totalSteps = 3

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
      class="rounded-full retro-button circle bg-primary"
      @click="isModalOpen = true"
    >
      <i class="fas fa-info"></i>
    </button>

    <!-- Modal -->
    <div
      v-if="isModalOpen"
      class="flex fixed inset-0 z-50 justify-center items-center bg-black bg-opacity-50"
    >
      <div class="p-6 w-full max-w-lg bg-white rounded-lg retro-box">

        <!-- Modal Headers -->
        <div class="flex justify-between items-center mb-4">
          <h3 class="text-xl font-bold text-gray-700">Tutorial Videos</h3>
          <!-- Close Button -->
          <button class="text-gray-400 hover:text-gray-600" @click="closeModal">
            <i class="fas fa-times fa-lg"></i>
          </button>
        </div>
        <!-- Stepper -->
        <div class="mb-6">
          <div class="flex justify-between mb-4">
            <span
              v-for="step in totalSteps"
              :key="step"
              class="flex justify-center items-center w-8 h-8 rounded-full border-2"
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
            <video
              class="w-full h-48 object-cover rounded-lg"
              src="/tutorial-1.mp4"
              controls
              ></video>
            <p class="text-gray-600">Step 1: Upload your video files (TODO).</p>
            <!-- Add your input components here -->
          </div>
          <div v-else-if="currentStep === 2">
            <p class="text-gray-600">
              Step 2: Review and organize videos (TODO).
            </p>
          </div>
          <div v-else-if="currentStep === 3">
            <p class="text-gray-600">Step 3: Confirm and save (TODO).</p>
          </div>
        </div>

        <!-- Step Navigation -->
        <div class="flex justify-between">
          <button
            v-if="currentStep > 1"
            class="bg-gray-300 retro-button"
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
            class="bg-green-500 retro-button"
            @click="closeModal"
          >
            Finish
          </button>
        </div>
        <button
          class="absolute top-3 right-3 text-gray-400 hover:text-gray-600"
          @click="closeModal"
        >
          <i class="fas fa-times"></i>
        </button>
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
