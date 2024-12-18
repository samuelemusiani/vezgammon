<template>
  <div
    v-if="showRotateMessage"
    class="fixed inset-0 z-50 flex items-center justify-center bg-gray-800/80 backdrop-blur-sm"
  >
    <div
      class="relative mx-4 rounded-xl bg-white p-8 text-center shadow-2xl dark:bg-gray-800"
    >
      <div class="phone-animation mb-6">
        <div class="phone">
          <div class="phone-inner"></div>
        </div>
      </div>
      <h2 class="mb-2 text-xl font-bold text-gray-800 dark:text-white">
        Please Rotate Your Device
      </h2>
      <p class="text-gray-600 dark:text-gray-300">
        For the best experience, please use landscape mode
      </p>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'

const showRotateMessage = ref(false)

const checkOrientation = () => {
  showRotateMessage.value = window.innerHeight > window.innerWidth
}

onMounted(() => {
  checkOrientation()
  window.addEventListener('orientationchange', checkOrientation)
  window.addEventListener('resize', checkOrientation)
})

onUnmounted(() => {
  window.removeEventListener('orientationchange', checkOrientation)
  window.removeEventListener('resize', checkOrientation)
})
</script>

<style scoped>
.phone-animation {
  position: relative;
  width: 100px;
  height: 100px;
  margin: 0 auto;
}

.phone {
  position: absolute;
  left: 50%;
  top: 50%;
  transform: translate(-50%, -50%);
  width: 60px;
  height: 90px;
  border: 3px solid #8b4513;
  border-radius: 8px;
  animation: rotate 3s infinite ease-in-out;
  box-shadow:
    0 0 0 2px #d2691e,
    inset 0 0 10px rgba(0, 0, 0, 0.2);
}

.phone-inner {
  position: absolute;
  left: 50%;
  top: 50%;
  transform: translate(-50%, -50%);
  width: 40px;
  height: 70px;
  background: #ffe5c9;
  border: 2px solid #8b4513;
  border-radius: 4px;
  box-shadow: inset 0 0 10px rgba(139, 69, 19, 0.3);
}

@keyframes rotate {
  0% {
    transform: translate(-50%, -50%) rotate(0deg);
  }
  20% {
    transform: translate(-50%, -50%) rotate(-90deg);
  }
  80% {
    transform: translate(-50%, -50%) rotate(-90deg);
  }
  100% {
    transform: translate(-50%, -50%) rotate(0deg);
  }
}
</style>
