<template>
  <div
    class="retro-background flex max-h-full w-full flex-col items-center justify-center"
  >
    <h1 class="text-4xl font-bold">Work in Progress</h1>
    <p class="mt-4 text-lg">This page is still under construction.</p>
  </div>
</template>
<style scoped>
.retro-background {
  @apply flex min-h-screen items-center justify-center;
  background: #2c1810;
  background-image: repeating-linear-gradient(
      45deg,
      rgba(139, 69, 19, 0.1) 0px,
      rgba(139, 69, 19, 0.1) 2px,
      transparent 2px,
      transparent 10px
    ),
    repeating-linear-gradient(
      -45deg,
      rgba(139, 69, 19, 0.1) 0px,
      rgba(139, 69, 19, 0.1) 2px,
      transparent 2px,
      transparent 10px
    );
  cursor: url('/tortellino.png'), auto;
  border: 6px solid #d2691e;
}
</style>
