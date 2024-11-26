<script setup lang="ts">
import { ref } from 'vue'
import router from '@/router'
import { useTheme } from '@/composables/useTheme'

const username = ref('')
const passwd = ref('')
const err = ref('')

// Gestione dei temi
const { currentTheme, themeOptions, changeTheme } = useTheme()

async function login() {
  if (!validate()) {
    return
  }

  try {
    const response = await fetch('/api/login', {
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
      const message = await response.text()
      throw new Error(message || 'Error during login')
    }

    router.push({ name: 'home' })
  } catch (e) {
    console.error(e)
    err.value = e instanceof Error ? e.message : 'Unexpected error'
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
  <div class="flex items-center justify-center h-full w-full">
    <div class="card w-96 bg-base-100 shadow-md border border-primary border-8">
      <div class="card-body">
        <h2 class="card-title">Login</h2>
        <div class="card-body">
          <!-- Gestione Tema -->
          <div class="mb-4 flex justify-end">
            <div class="dropdown dropdown-end">
              <div tabindex="0" role="button" class="btn m-1">
                Theme: {{ currentTheme }}
              </div>
              <ul
                tabindex="0"
                class="p-2 shadow menu dropdown-content bg-base-100 rounded-box w-52"
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

          <!-- Form di Login -->
          <form @submit.prevent="login">
            <div class="form-control">
              <label class="label">
                <span class="label-text">Username or Email</span>
              </label>
              <input
                type="text"
                placeholder="Enter Username"
                v-model="username"
                class="input input-bordered w-full focus:ring-primary bg-base-200"
              />
            </div>

            <div class="form-control mt-4">
              <label class="label">
                <span class="label-text">Password</span>
              </label>
              <input
                type="password"
                placeholder="Enter Password"
                v-model="passwd"
                class="input input-bordered w-full focus:ring-primary bg-base-200"
              />
              <label class="label">
                <a href="#" class="link-hover link label-text-alt">Forgot password?</a>
              </label>
            </div>

            <!-- Errore -->
            <div class="mt-4 text-error">
              {{ err }}
            </div>

            <!-- Pulsante di invio -->
            <div class="form-control mt-6">
              <button type="submit" class="btn btn-primary border-secondary border-4">Login</button>
            </div>
          </form>

          <!-- Divider e Link di Registrazione -->
          <div class="divider">OR</div>
          <div class="text-center">
            <p>Don't have an account?</p>
            <RouterLink to="/register" class="link link-primary">
              Sign up now
            </RouterLink>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
