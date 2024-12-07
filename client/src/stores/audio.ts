import { defineStore } from 'pinia'

export const useAudioStore = defineStore('audio', {
  state: () => ({
    isAudioEnabled: true, // default true
  }),
  actions: {
    toggleAudio() {
      this.isAudioEnabled = !this.isAudioEnabled
    },
  },
})
