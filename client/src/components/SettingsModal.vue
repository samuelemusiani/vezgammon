<template>
  <dialog id="settings_modal" class="modal">
    <div class="modal-box w-full border-2 border-primary bg-base-200">
      <h3 class="mb-4 text-center text-2xl font-bold text-primary">Settings</h3>

      <div class="form-control w-full">
        <label class="label cursor-pointer">
          <span class="label-text">Sound Effects</span>
          <input
            type="checkbox"
            class="toggle toggle-primary"
            :checked="audioStore.isAudioEnabled"
            @change="toggleAudio"
          />
        </label>
      </div>

      <div class="form-control">
        <label class="label" for="theme-selector">
          <span class="label-text">Theme</span>
        </label>
        <select
          id="theme-selector"
          class="select select-bordered w-full bg-base-100"
          :value="currentTheme"
          @change="e => changeTheme((e.target as HTMLSelectElement).value)"
        >
          <option v-for="theme in themeOptions" :key="theme" :value="theme">
            {{ theme }}
          </option>
        </select>
      </div>

      <!-- Change Password -->
      <div class="form-control mt-4">
        <button
          @click="togglePasswordSection"
          class="btn btn-primary btn-sm w-full"
        >
          {{
            showPasswordSection ? 'Hide Password Settings' : 'Change Password'
          }}
        </button>

        <div v-if="showPasswordSection" class="mt-4 space-y-2">
          <div>
            <label class="label" for="current-password">
              <span class="label-text">Current Password</span>
            </label>
            <input
              id="current-password"
              type="password"
              v-model="currentPassword"
              class="input input-bordered w-full"
              placeholder="Enter current password"
            />
          </div>

          <div>
            <label class="label" for="new-password">
              <span class="label-text">New Password</span>
            </label>
            <input
              id="new-password"
              type="password"
              v-model="newPassword"
              class="input input-bordered w-full"
              placeholder="Enter new password"
            />
          </div>

          <div>
            <label class="label" for="new-password-confirm">
              <span class="label-text">Confirm New Password</span>
            </label>
            <input
              id="new-password-confirm"
              type="password"
              v-model="confirmPassword"
              class="input input-bordered w-full"
              placeholder="Confirm new password"
            />
          </div>
          <div class="mt-2 flex items-center justify-between">
            <span
              class="text-sm"
              :class="{
                'text-success': passwdMessage.includes('successfully'),
                'text-error': !passwdMessage.includes('successfully'),
              }"
              >{{ passwdMessage }}</span
            >
            <button
              @click="handleChangePassword"
              class="btn btn-primary btn-sm"
              :disabled="!isFormValid"
            >
              Confirm
            </button>
          </div>
        </div>
      </div>

      <!-- Change Avatar -->
      <div class="form-control mt-4">
        <button
          @click="toggleAvatarSection"
          class="btn btn-primary btn-sm w-full"
        >
          {{ showAvatarSection ? 'Hide Avatar Settings' : 'Change Avatar' }}
        </button>

        <div v-if="showAvatarSection" class="mt-4 space-y-2">
          <div>
            <label class="label" for="current-avatar">
              <span class="label-text">Current Avatar</span>
              <input
                id="current-avatar"
                type="text"
                readonly
                class="input input-bordered w-full"
                v-model="currentAvatar"
              />
              <img
                class="ml-2 h-16 w-16 rounded-full border-2 border-primary"
                :src="currentAvatar"
                alt="Current avatar"
              />
            </label>
            <label class="label" for="new-avatar">
              <span class="label-text">New Avatar</span>
              <input
                id="new-avatar"
                type="text"
                class="input input-bordered w-full"
                v-model="newAvatar"
              />
              <img
                v-if="newAvatar && isValid(newAvatar)"
                class="ml-2 h-16 w-16 rounded-full border-2 border-primary"
                :src="newAvatar"
                alt="Not found"
              />
            </label>
          </div>
          <div class="mt-2 flex items-center justify-between">
            <span
              class="text-sm"
              :class="{
                'text-success': avatarMessage.includes('successfully'),
                'text-error': !avatarMessage.includes('successfully'),
              }"
              >{{ avatarMessage }}</span
            >
            <button
              @click="handleChangeAvatar"
              class="btn btn-primary btn-sm"
              :disabled="!isValid(newAvatar)"
            >
              Confirm
            </button>
          </div>
        </div>
      </div>

      <div class="modal-action">
        <form method="dialog">
          <button class="btn btn-primary">Close</button>
        </form>
      </div>
    </div>

    <form method="dialog" class="modal-backdrop">
      <button>close</button>
    </form>
  </dialog>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useAudioStore } from '@/stores/audio'
import { useTheme } from '@/composables/useTheme'
import { vfetch } from '@/utils/fetch'

const audioStore = useAudioStore()
const { currentTheme, themeOptions, changeTheme } = useTheme()

const currentPassword = ref('')
const newPassword = ref('')
const confirmPassword = ref('')
const showPasswordSection = ref(false)
const passwdMessage = ref('')

const showAvatarSection = ref(false)
const avatarMessage = ref('')
const currentAvatar = ref('')
const newAvatar = ref('')

const fetchAvatar = async () => {
  try {
    const res = await vfetch('/api/session')
    const data = await res.json()
    currentAvatar.value = data.avatar
  } catch (e: any) {
    console.error('Error fetching avatar:', e.message)
  }
}

onMounted(() => {
  fetchAvatar()
})

const togglePasswordSection = () => {
  showPasswordSection.value = !showPasswordSection.value
  if (!showPasswordSection.value) {
    currentPassword.value = ''
    newPassword.value = ''
    confirmPassword.value = ''
  }
}

function isValid(newAvatar: string) {
  const dicebearRegex =
    /^https:\/\/api\.dicebear\.com\/\d+\.x\/[^/]+\/svg\?seed=[^&]+$/
  return dicebearRegex.test(newAvatar)
}

const toggleAvatarSection = () => {
  showAvatarSection.value = !showAvatarSection.value
  if (!showAvatarSection.value) {
    newAvatar.value = ''
  }
}

const isFormValid = computed(() => {
  return (
    currentPassword.value.length > 7 &&
    newPassword.value.length > 7 &&
    confirmPassword.value.length > 7
  )
})

const handleChangePassword = async () => {
  if (newPassword.value !== confirmPassword.value) {
    passwdMessage.value = 'New passwords do not match'
    return
  } else if (newPassword.value === currentPassword.value) {
    passwdMessage.value = 'New password must be different from current password'
    return
  }

  try {
    passwdMessage.value = ''

    const res = await vfetch('/api/pass', {
      method: 'PATCH',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        old_pass: currentPassword.value,
        new_pass: newPassword.value,
      }),
    })

    if (!res.ok) {
      const err = await res.json()
      throw new Error(err.error)
    }

    passwdMessage.value = 'Password changed successfully'
    currentPassword.value = ''
    newPassword.value = ''
    confirmPassword.value = ''
  } catch (e: any) {
    passwdMessage.value = e.message
  }
}

const handleChangeAvatar = async () => {
  try {
    avatarMessage.value = ''

    const res = await vfetch('/api/avatar', {
      method: 'PATCH',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        avatar: newAvatar.value,
      }),
    })

    if (!res.ok) {
      const err = await res.json()
      throw new Error(err.error)
    }

    avatarMessage.value = 'Avatar changed successfully'
  } catch (e: any) {
    avatarMessage.value = e.message
  }
}

const toggleAudio = () => {
  audioStore.toggleAudio()
}
</script>
