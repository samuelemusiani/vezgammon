<template>
  <div class="flex h-full w-full items-center justify-center">
    <div
      class="flex h-[94%] w-[80%] flex-col items-center overflow-y-auto rounded-md border-8 border-primary bg-base-100 md:justify-center"
    >
      <div class="mb-8 flex flex-col items-center text-center xl:mb-16">
        <h1
          class="retro-title mb-8 w-60 text-2xl font-bold md:w-full md:p-4 md:text-3xl lg:text-4xl xl:text-7xl"
        >
          404 Not Found
        </h1>
        <div class="font-bold text-accent md:text-lg lg:text-xl">
          The page you are looking for does not exist.
        </div>
        <button
          @click="router.push('/')"
          class="retro-button mt-8 p-2 md:mt-12 md:p-4"
        >
          Return Home
        </button>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import router from '@/router'
</script>
