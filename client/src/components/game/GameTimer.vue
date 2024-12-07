<script setup lang="ts">
defineProps<{
  timeLeft: number
  isMyTurn: boolean
}>()
</script>

<template>
  <div v-show="isMyTurn" class="w-full">
    <div class="mb-2 text-center text-xl font-bold">
      Time Left: {{ timeLeft }}s
    </div>
    <div class="h-2 w-full rounded-full bg-gray-200">
      <div
        class="h-full rounded-full bg-primary transition-all duration-1000"
        :style="{ width: `${(timeLeft / 60) * 100}%` }"
        :class="{
          'bg-red-500': timeLeft <= 10,
          'bg-yellow-500': timeLeft <= 30 && timeLeft > 10,
        }"
      ></div>
    </div>
  </div>
</template>
