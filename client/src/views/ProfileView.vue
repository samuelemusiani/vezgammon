<script setup lang="ts">
import { ref } from 'vue'
import type { User } from '@/utils/types'
import router from '@/router'

const session = ref<User | undefined>()
const error = ref<string>('')

fetch('/api/session')
  .then(response => {
    if (!response.ok) {
      throw new Error('During profile fetch: ' + response.statusText)
    }
    return response.json()
  })
  .then(data => {
    session.value = data
  })
  .catch(e => {
    console.error(e)
  })

async function logout() {
  await fetch('/api/logout', { method: 'POST' })
  router.push({ name: 'login' })
}
</script>

<template>
  <div
    class="retro-background flex min-h-screen items-center justify-center bg-base-200"
  >
    <div class="card w-1/2 bg-base-100 shadow-xl">
      <div class="retro-box card-body">
        <h2 class="text-center text-2xl font-bold">Profile</h2>

        <div class="divider divider-neutral"></div>

        <div v-if="session" class="flex flex-col gap-4">
          <div>
            <span> Username: </span>
            <span class="text-lg font-bold"> {{ session.username }} </span>
          </div>
          <div>
            <span> Mail: </span>
            <span class="text-lg font-bold">{{ session.mail }}</span>
          </div>
          <div>
            <span> Fullname: </span>
            <span class="text-lg font-bold">
              {{ session.firstname }} {{ session.lastname }}
            </span>
          </div>

          <div class="mt-10 flex justify-center">
            <button class="btn btn-warning" @click="logout">LOGOUT</button>
          </div>
        </div>

        <div v-else class="text-error">
          {{ error }}
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.retro-background {
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

.retro-box {
  background-color: #ffe5c9;
  border: 5px solid #8b4513;
  box-shadow:
    0 0 0 4px #d2691e,
    inset 0 0 20px rgba(0, 0, 0, 0.2);
}
</style>
