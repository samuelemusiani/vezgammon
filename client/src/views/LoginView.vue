<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useTheme } from '@/composables/useTheme' // We'll create this
import router from '@/router'

// Theme management
const { currentTheme, themeOptions, changeTheme } = useTheme()

const username = ref('')
const passwd = ref('')
const err = ref('')

async function login() {
  if (!validate()) {
    return
  }

  try {
    const response = await fetch('/api/login', {
      method: 'POST',
      body: JSON.stringify({
        username: username.value,
        password: passwd.value,
      }),
    })

    if (!response.ok) {
      throw new Error('During login: ' + (await response.text()))
    }

    router.push({ name: 'home' })
  } catch (e) {
    console.error(e)
    err.value = 'Error during login'
  }
}

function validate() {
  const el = document.getElementsByClassName('X-required')
  for (const e of el) {
    e.setAttribute('required', '')
  }

  let ok = true
  let e = ''
  if (username.value.length == 0) {
    e = "Username can't be empty "
    ok = false
  } else if (passwd.value.length == 0) {
    e = "Password can't be empty "
    ok = false
  }

  err.value = e

  return ok
}
</script>

<template>
  <div class="flex min-h-screen items-center justify-center bg-neutral">
    <div class="card w-96 bg-base-100 shadow-xl">
      <div class="card-body">
        <div class="mb-4 flex justify-end">
          <div class="dropdown dropdown-end">
            <div tabindex="0" role="button" class="btn m-1">
              Theme: {{ currentTheme }}
            </div>
            <ul
              tabindex="0"
              class="menu dropdown-content z-[1] w-52 rounded-box bg-base-100 p-2 shadow"
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

        <h2 class="card-title">Login</h2>
        <form @submit.prevent="login">
          <div class="form-control">
            <label class="label">
              <span class="label-text">Username or mail</span>
            </label>
            <input
              type="text"
              placeholder="Enter Username"
              v-model="username"
              class="input input-bordered w-full"
            />
          </div>

          <div class="form-control mt-4">
            <label class="label">
              <span class="label-text">Password</span>
            </label>
            <input
              type="password"
              placeholder="Enter password"
              v-model="passwd"
              class="input input-bordered w-full"
            />
            <label class="label">
              <a href="#" class="link-hover link label-text-alt"
                >Forgot password?</a
              >
            </label>
          </div>

          <div class="mt-4 text-error">
            {{ err }}
          </div>

          <div class="form-control mt-6">
            <button type="submit" class="btn btn-primary">Login</button>
          </div>
        </form>

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
</template>
