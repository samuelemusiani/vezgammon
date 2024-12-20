<script setup lang="ts">
import { onMounted, ref } from 'vue'
import type { User, Badge } from '@/utils/types'
import router from '@/router'
import Badges from '@/components/Badges.vue'
import { vfetch } from '@/utils/fetch'

const badges = ref<Badge | null>()

const fetchBadges = async () => {
  try {
    const res = await vfetch('/api/badge')
    const data = await res.json()
    badges.value = data
  } catch (e: any) {
    console.error('Error fetching badges:', e.message)
  }
}

const session = ref<User | undefined>()
const error = ref<string>('')

const fetchSession = async () => {
  vfetch('/api/session')
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
}

onMounted(() => {
  fetchSession()
  fetchBadges()
})

async function logout() {
  await vfetch('/api/logout', { method: 'POST' })
  router.push({ name: 'login' })
}

async function goBack() {
  router.push({ name: 'home' })
}
</script>

<template>
  <div class="flex h-full items-center justify-center">
    <div
      class="card h-[94%] w-3/4 overflow-y-auto rounded-xl border-8 border-primary bg-base-100 shadow-xl"
    >
      <div class="card-body">
        <h2 class="text-center text-xl font-bold md:text-3xl xl:text-4xl">
          Profile
        </h2>

        <div class="divider divider-neutral"></div>

        <div v-if="session" class="flex flex-col gap-4">
          <div
            class="m-4 flex flex-col items-center justify-between gap-4 md:flex-row"
          >
            <div
              class="flex w-full flex-col gap-4 md:text-xl xl:flex-row xl:justify-evenly"
            >
              <div>
                <span> Username: </span>
                <span class="font-bold"> {{ session.username }} </span>
              </div>
              <div>
                <span> Mail: </span>
                <span class="font-bold">{{ session.mail }}</span>
              </div>
              <div>
                <span> Fullname: </span>
                <span class="font-bold">
                  {{ session.firstname }} {{ session.lastname }}
                </span>
              </div>
            </div>
            <div>
              <img
                class="h-32 w-32 rounded-full border-4 border-primary"
                :src="session.avatar"
                alt="User avatar"
              />
            </div>
          </div>

          <div class="divider divider-neutral">Your Badges</div>

          <Badges v-if="badges" :badges="badges" />

          <div class="mt-8 flex items-center justify-center gap-5">
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
