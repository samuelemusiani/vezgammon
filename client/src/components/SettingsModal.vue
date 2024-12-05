<template>
  <dialog id="settings_modal" class="modal">
    <div class="modal-box border-2 border-primary bg-base-200">
      <h3 class="mb-4 text-center text-2xl font-bold text-primary">Settings</h3>

      <div class="form-control">
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
        <label class="label">
          <span class="label-text">Theme</span>
        </label>
        <select
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
            <label class="label">
              <span class="label-text">Current Password</span>
            </label>
            <input
              type="password"
              v-model="currentPassword"
              class="input input-bordered w-full"
              placeholder="Enter current password"
            />
          </div>

          <div>
            <label class="label">
              <span class="label-text">New Password</span>
            </label>
            <input
              type="password"
              v-model="newPassword"
              class="input input-bordered w-full"
              placeholder="Enter new password"
            />
          </div>

          <div>
            <label class="label">
              <span class="label-text">Confirm New Password</span>
            </label>
            <input
              type="password"
              v-model="confirmPassword"
              class="input input-bordered w-full"
              placeholder="Confirm new password"
            />
          </div>
        </div>
      </div>

      <div class="mt-2 flex items-center justify-between">
        <span class="text-sm text-error">{{ errorMessage }}</span>
        <button
          @click="handleChangePassword"
          class="btn btn-primary btn-sm"
          :disabled="!isFormValid"
        >
          Confirm
        </button>
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
import { ref, computed } from 'vue'
import { useAudioStore } from '@/stores/audio'
import { useTheme } from '@/composables/useTheme'

const audioStore = useAudioStore()
const { currentTheme, themeOptions, changeTheme } = useTheme()

// Password state
const currentPassword = ref('')
const newPassword = ref('')
const confirmPassword = ref('')
const showPasswordSection = ref(false)
const errorMessage = ref('')

const togglePasswordSection = () => {
  showPasswordSection.value = !showPasswordSection.value
  if (!showPasswordSection.value) {
    currentPassword.value = ''
    newPassword.value = ''
    confirmPassword.value = ''
  }
}

// Computed property per validare il form
const isFormValid = computed(() => {
  return (
    currentPassword.value.length > 8 &&
    newPassword.value.length > 8 &&
    confirmPassword.value.length > 8 &&
    newPassword.value === confirmPassword.value &&
    newPassword.value !== currentPassword.value
  )
})

// Handler per il cambio password
const handleChangePassword = async () => {
  if (!isFormValid.value) {
    if (newPassword.value !== confirmPassword.value) {
      errorMessage.value = 'New passwords do not match'
    } else if (newPassword.value === currentPassword.value) {
      errorMessage.value =
        'New password must be different from current password'
    }
    return
  }

  try {
    errorMessage.value = ''

    // TODO: when backend is ready
    /*await fetch('/api/change-password', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        currentPassword: currentPassword.value,
        newPassword: newPassword.value,
      }),
      })*/

    // Reset form
    currentPassword.value = ''
    newPassword.value = ''
    confirmPassword.value = ''
  } catch (error: any) {
    errorMessage.value = 'Failed to change password'
  }
}

const toggleAudio = () => {
  audioStore.toggleAudio()
}
</script>
