<template>
  <dialog id="settings_modal" class="modal">
    <div class="retro-box modal-box">
      <h3 class="retro-title mb-4 text-center text-2xl font-bold">Settings</h3>

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

      <div class="modal-action">
        <form method="dialog">
          <button class="retro-button">Close</button>
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

const audioStore = useAudioStore()

const toggleAudio = () => {
  audioStore.toggleAudio()
}
</script>
