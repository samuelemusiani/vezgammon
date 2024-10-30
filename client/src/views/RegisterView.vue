<script setup lang="ts">
import { ref } from 'vue'
import router from '@/router'

const email = ref('')
const passwd1 = ref('')
const passwd2 = ref('')
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
        email: email.value,
        password: passwd1.value,
        firstname: firstname.value,
        lastname: lastname.value,
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
  for (let i = 0; i < el.length; i++) {
    el[i].setAttribute('required', '')
  }

  let ok = true

  let e = ''
  if (
    email.value.length == 0 ||
    !email.value.includes('@') ||
    !email.value.includes('.', email.value.indexOf('@'))
  ) {
    e = 'Insert a valid email'
    ok = false
  } else if (passwd1.value.length < 8) {
    e = 'Insert a password of at least 8 character'
    ok = false
  } else if (passwd1.value != passwd2.value) {
    e = 'Passwords do not match'
    ok = false
  } else if (firstname.value.length < 3) {
    e = 'Name too short'
    ok = false
  }
  err.value = e

  return ok
}
</script>

<template>
  <div class="bg-base-200 min-h-screen flex justify-center items-center">
    <div class="bg-base-100 card w-96 shadow-xl">
      <div class="card-body">
        <h2 class="card-title">Register</h2>
        <form @click.prevent="">
          <div class="form-control">
            <label for="input-mail" class="label">
              <span class="label-text"> Email </span>
            </label>
            <div
              class="input input-bordered flex items-center gap-2 has-[:invalid]:border-error"
            >
              <input
                id="input-mail"
                type="email"
                class="grow invalid:text-error X-required"
                v-model="email"
                placeholder="email@example.com"
              />
            </div>
          </div>

          <div class="form-control mt-5">
            <label for="input-password1" class="label">
              <span class="label-text"> Password </span>
            </label>
            <div
              class="input input-bordered flex items-center gap-2 has-[:invalid]:border-error"
            >
              <input
                id="input-password1"
                type="password"
                class="grow invalid:text-error X-required"
                minlength="8"
                v-model="passwd1"
                placeholder="Choose a password"
              />
            </div>
            <label for="input-password2" class="label">
              <span class="label-text"> Confirm Password </span>
            </label>
            <div
              class="input input-bordered flex items-center gap-2 has-[:invalid]:border-error"
            >
              <input
                id="input-password2"
                type="password"
                class="grow invalid:text-error X-required"
                minlength="8"
                v-model="passwd2"
                placeholder="Confirm password"
              />
            </div>

            <div class="form-control mt-5">
              <label for="input-name" class="label">
                <span class="label-text"> First Name </span>
              </label>
              <div
                class="input input-bordered flex items-center gap-2 has-[:invalid]:border-error"
              >
                <input
                  id="input-name"
                  type="text"
                  class="grow invalid:text-error X-required"
                  placeholder="Alex"
                  v-model="firstname"
                />
              </div>

              <label for="input-surname" class="label">
                <span class="label-surname"> Last Name </span>
              </label>
              <div class="input input-bordered flex items-center gap-2">
                <input
                  id="input-surname"
                  type="text"
                  class="grow"
                  placeholder="Harrison"
                  v-model="lastname"
                />
              </div>
            </div>
          </div>

          <div class="text-error mt-2">
            {{ err }}
          </div>

          <div class="form-control mt-6">
            <button class="btn btn-secondary" @click="register">
              Register
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>
