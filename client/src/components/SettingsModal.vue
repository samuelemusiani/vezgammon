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
import { useAudioStore } from '@/stores/audio'
import { useTheme } from '@/composables/useTheme'

const audioStore = useAudioStore()
const { currentTheme, themeOptions, changeTheme } = useTheme()

const toggleAudio = () => {
  audioStore.toggleAudio()
}
</script>
