<script setup lang="ts">
import { ref } from 'vue'
import type { User } from '@/utils/types'
import router from '@/router'
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
          </div>
        </div>

        <div v-else class="text-error">
          {{ error }}
        </div>
      </div>
    </div>
  </div>
</template>
