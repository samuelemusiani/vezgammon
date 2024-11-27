<script setup lang="ts">
import { ref } from 'vue'
import router from '@/router'

const mail = ref('')
const passwd1 = ref('')
const passwd2 = ref('')
const username = ref('')
const firstname = ref('')
const lastname = ref('')

const err = ref('')

async function register() {
  if (!validate()) {
    return
  }

  try {
    const response = await fetch('/api/register', {
      method: 'POST',
      body: JSON.stringify({
        mail: mail.value,
        password: passwd1.value,
        username: username.value,
        firstname: firstname.value,
        lastname: lastname.value,
      }),
    })

    if (!response.ok) {
      throw new Error('During registration: ' + (await response.text()))
    }

    router.push({ name: 'login' })
  } catch (e) {
    console.error(e)
    err.value = 'Error during registration'
  }
}

function validate() {
  const el = document.getElementsByClassName('X-required')
  for (const e of el) {
    e.setAttribute('required', '')
  }

  let ok = true

  let e = ''
  if (
    mail.value.length == 0 ||
    !mail.value.includes('@') ||
    !mail.value.includes('.', mail.value.indexOf('@'))
  ) {
    e = 'Insert a valid email'
    ok = false
  } else if (passwd1.value.length < 8) {
    e = 'Insert a password of at least 8 character'
    ok = false
  } else if (passwd1.value != passwd2.value) {
    e = 'Passwords do not match'
    ok = false
  } else if (username.value.length < 3) {
    e = 'Username too short'
    ok = false
  }
  err.value = e

  return ok
}
</script>

<template>
  <div class="flex items-center justify-center h-full w-full">
    <div class="card w-[40%] min-w-96 h-[90%] bg-base-100 shadow-md border border-primary border-8 overflow-auto">
      <div class="card-body">
        <h2 class="card-title">Register</h2>
        <form @submit.prevent="register" class="grid grid-cols-2 gap-4">
          <div class="form-control">
            <label for="input-mail" class="label">
              <span class="label-text">Email</span>
            </label>
            <input
              id="input-mail"
              type="email"
              placeholder="email@example.com"
              v-model="mail"
              class="input input-bordered w-full focus:ring-primary bg-base-200 invalid:border-error invalid:text-error"
            />
          </div>

          <div class="form-control">
            <label for="input-username" class="label">
              <span class="label-text">Username</span>
            </label>
            <input
              id="input-username"
              type="text"
              placeholder="Master725"
              v-model="username"
              class="input input-bordered w-full focus:ring-primary bg-base-200"
            />
          </div>

          <div class="form-control">
            <label for="input-password1" class="label">
              <span class="label-text">Password</span>
            </label>
            <input
              id="input-password1"
              type="password"
              placeholder="Choose a password"
              v-model="passwd1"
              minlength="8"
              class="input input-bordered w-full focus:ring-primary bg-base-200 invalid:border-error invalid:text-error"
            />
          </div>

          <div class="form-control">
            <label for="input-password2" class="label">
              <span class="label-text">Confirm Password</span>
            </label>
            <input
              id="input-password2"
              type="password"
              placeholder="Confirm password"
              v-model="passwd2"
              minlength="8"
              class="input input-bordered w-full focus:ring-primary bg-base-200 invalid:border-error invalid:text-error"
            />
          </div>

          <div class="form-control">
            <label for="input-name" class="label">
              <span class="label-text">First Name</span>
            </label>
            <input
              id="input-name"
              type="text"
              placeholder="Alex"
              v-model="firstname"
              class="input input-bordered w-full focus:ring-primary bg-base-200"
            />
          </div>

          <div class="form-control">
            <label for="input-surname" class="label">
              <span class="label-text">Last Name</span>
            </label>
            <input
              id="input-surname"
              type="text"
              placeholder="Harrison"
              v-model="lastname"
              class="input input-bordered w-full focus:ring-primary bg-base-200"
            />
          </div>

          <div class="col-span-2 text-error mt-6">
            {{ err }}
          </div>

          <div class="form-control col-span-2 mt-2">
            <button type="submit" class="btn btn-primary border-secondary border-4">
              Register
            </button>
          </div>

          <div class="col-span-2 divider">OR</div>

          <div class="col-span-2 text-center">
            <p>Already have an account?</p>
            <RouterLink to="/login" class="link link-primary">
              Go back to Login
            </RouterLink>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>
