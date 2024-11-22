<script setup lang="ts">
import { ref } from 'vue'
import router from '@/router'

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
  <div
    class="retro-background flex min-h-screen items-center justify-center bg-base-200"
  >
    <div class="card w-96 bg-base-100 shadow-xl">
      <div class="retro-box card-body">
        <h2 class="card-title">Login</h2>
        <form @click.prevent="">
          <div class="form-control">
            <label for="input-mail" class="label">
              <span class="label-text"> Username or mail </span>
            </label>
            <div
              class="input input-bordered flex items-center gap-2 has-[:invalid]:border-error"
            >
              <input
                id="input-text"
                type="text"
                class="X-required grow invalid:text-error"
                placeholder="Enter Username"
                v-model="username"
              />
            </div>
          </div>

          <div class="form-control">
            <label for="input-password" class="label">
              <span class="label-text"> Password </span>
            </label>
            <div
              class="input input-bordered flex items-center gap-2 has-[:invalid]:border-error"
            >
              <input
                id="input-password"
                type="password"
                class="X-required grow invalid:text-error"
                placeholder="Enter password"
                v-model="passwd"
              />
            </div>
            <div class="label">
              <a href="#" class="link-hover link label-text-alt"
                >Forgot password?</a
              >
            </div>
          </div>

          <div class="mt-5 text-error">
            {{ err }}
          </div>

          <div class="form-control mt-6">
            <button class="retro-button btn btn-primary" @click="login">
              Login
            </button>
          </div>
        </form>
        <div class="divider">OR</div>
        <div class="text-center">
          <p>Don't have an account?</p>
          <RouterLink to="/register" class="link link-primary"
            >Sign up now</RouterLink
          >
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

.retro-button {
  @apply btn;
  background: #d2691e;
  color: white;
  border: 3px solid #8b4513;
  font-family: 'Arial Black', serif;
  text-transform: uppercase;
  text-shadow: 2px 2px 0 rgba(0, 0, 0, 0.2);
  box-shadow: 0 2px 0 #8b4513;
  font-size: 1.1rem;

  &:hover {
    transform: translateY(2px);
    box-shadow:
      inset 0 0 10px rgba(0, 0, 0, 0.2),
      0 0px 0 #8b4513;
    cursor: url('/tortellino.png'), auto;
  }
}
</style>
