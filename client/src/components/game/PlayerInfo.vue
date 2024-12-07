<script setup lang="ts">
interface PlayerInfoProps {
  username: string
  elo: number
  isCurrentTurn: boolean
  isOpponent?: boolean
  avatar: string
}

defineProps<PlayerInfoProps>()
</script>

<template>
  <div class="flex flex-col items-center">
    <div class="relative mb-2">
      <div class="h-16 w-16 overflow-hidden rounded-full bg-gray-200">
        <img
          :src="avatar"
          :alt="`${isOpponent ? 'Opponent' : 'Player'} avatar`"
          class="h-full w-full object-cover"
        />
      </div>
      <div
        :class="[
          'absolute -bottom-1 right-0 h-4 w-4 rounded-full border-2 border-white',
          isCurrentTurn ? 'bg-green-500' : 'bg-gray-300',
        ]"
      ></div>
    </div>
    <h3 class="text-lg font-bold">{{ username }}</h3>
    <p class="text-gray-600">ELO: {{ elo }}</p>
  </div>
</template>
