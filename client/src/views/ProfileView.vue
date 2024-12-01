<script setup lang="ts">
import { ref } from 'vue'
import type { User } from '@/utils/types'
import router from '@/router'
import { useTheme } from '@/composables/useTheme'
import Badges from '@/components/Badges.vue'

const badges = ref()

fetch('/api/badge')
  .then(response => response.json())
  .then(data => {
    badges.value = data
  })
  .catch(e => {
    console.error('Error fetching badges:', e)
  })

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

async function goBack() {
  router.push({ name: 'home' })
}

const { currentTheme, themeOptions, changeTheme } = useTheme()
</script>

<template>
  <div class="flex h-full items-center justify-center">
    <div
      class="card w-3/4 rounded-xl border-8 border-primary bg-base-100 shadow-xl"
    >
      <div class="card-body">
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

          <div class="divider divider-neutral">Your Badges</div>

          <Badges :badges="badges" />

          <div class="mt-10 flex items-center justify-center gap-5">
            <button class="btn-seconday btn" @click="goBack">GO BACK</button>
            <button class="btn btn-primary" @click="logout">LOGOUT</button>
            <div class="">
              <div class="dropdown dropdown-end">
                <div tabindex="0" role="button" class="btn m-1">
                  Theme: {{ currentTheme }}
                </div>
                <ul
                  tabindex="0"
                  class="menu dropdown-content w-52 rounded-box border-4 border-primary bg-base-100 p-2 shadow"
                >
                  <li
                    v-for="theme in themeOptions"
                    :key="theme"
                    @click="changeTheme(theme)"
                  >
                    <a>{{ theme }}</a>
                  </li>
                </ul>
              </div>
            </div>
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
