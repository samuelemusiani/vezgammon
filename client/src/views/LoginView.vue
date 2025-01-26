<script setup lang="ts">
import { ref } from 'vue'
import router from '@/router'
import { vfetch } from '@/utils/fetch'

const username = ref('')
const passwd = ref('')
const err = ref('')

async function login() {
  if (!validate()) {
    return
  }

  try {
    const response = await vfetch('/api/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        username: username.value,
        password: passwd.value,
      }),
    })

    if (!response.ok) {
      const message = await response.json()
      throw new Error('Username/Email or Password are invalid')
    }

    router.push({ name: 'home' })
  } catch (e) {
    console.error(e)
    err.value = 'Username or Password are invalid'
  }
}

function validate() {
  let ok = true
  let message = ''

  if (!username.value.trim()) {
    message = "Username can't be empty."
    ok = false
  } else if (!passwd.value.trim()) {
    message = "Password can't be empty."
    ok = false
  }

  err.value = message
  return ok
}
</script>

<template>
  <div class="flex h-full w-full items-center justify-center">
    <div
      class="card max-h-[90%] w-96 overflow-y-auto border-8 border-primary bg-base-100 shadow-md"
    >
      <div class="card-body">
        <h2 class="card-title">Login</h2>
        <div class="card-body">
          <!-- Form di Login -->
          <form @submit.prevent="login">
            <div class="form-control">
              <label class="label" for="username">
                <span class="label-text">Username or Email</span>
              </label>
              <input
                id="username"
                type="text"
                placeholder="Enter Username"
                v-model="username"
                class="input input-bordered w-full bg-base-200 focus:ring-primary"
              />
            </div>

            <div class="form-control mt-4">
              <label class="label" for="password">
                <span class="label-text">Password</span>
              </label>
              <input
                id="password"
                type="password"
                placeholder="Enter Password"
                v-model="passwd"
                class="input input-bordered w-full bg-base-200 focus:ring-primary"
              />
            </div>

            <!-- Errore -->
            <div class="mt-4 text-error">
              {{ err }}
            </div>

            <!-- Pulsante di invio -->
            <div class="form-control mt-6">
              <button
                type="submit"
                class="btn btn-primary border-4 border-secondary"
              >
                Login
              </button>
            </div>
          </form>

          <!-- Divider e Link di Registrazione -->
          <div class="divider">OR</div>
          <div class="text-center">
            <p>Don't have an account?</p>
            <router-link to="/register" class="link link-primary">
              Sign up now
            </router-link>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
