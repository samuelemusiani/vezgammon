<template>
  <div class="flex min-h-screen items-center justify-center">
    <div class="flex max-h-full w-full flex-col items-center justify-center">
      <!-- Game Title -->
      <div class="mb-32 text-center">
        <h1 class="text-7xl">VezGammon</h1>
        <h4>The Ultimate Backgammon Experience</h4>
      </div>

      <!-- Button Container -->
      <div class="relative flex w-full max-w-4xl items-center justify-center">
        <!-- Left Button (Stats) -->
        <div class="mr-20">
          <button class="retro-button circle" title="Statistics">
            <MedalIcon />
          </button>
        </div>

        <!-- Central Buttons -->
        <div class="flex w-full max-w-sm flex-col gap-6">
          <button class="retro-button">PLAY</button>
          <button class="retro-button">RULES</button>
          <button class="retro-button">SETTINGS</button>
        </div>

        <!-- Right Button (Profile) -->
        <div class="ml-20">
          <button class="retro-button circle" title="Profile">
            <ProfileIcon />
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import MedalIcon from '@/utils/icons/MedalIcon.vue'
import ProfileIcon from '@/utils/icons/ProfileIcon.vue'
</script>

<style scoped>
.retro-button {
  @apply btn;
  background: #d2691e;
  color: white;
  border: 3px solid #8b4513;
  box-shadow: 0 2px 0 #8b4513;

  &.circle {
    width: 70px;
    height: 70px;
    border-radius: 50%;
  }
}
</style>
