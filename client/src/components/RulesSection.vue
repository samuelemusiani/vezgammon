<template>
  <div class="text-primary">
    <h3 class="mb-6 text-center text-3xl font-bold text-primary">
      Backgammon Rules
    </h3>
    <div
      class="rule-section"
      v-for="(section, index) in ruleSections"
      :key="index"
      :class="{ active: activeSection === index }"
      @click="activeSection = index"
    >
      <h4 class="mb-2 flex cursor-pointer justify-between text-xl font-bold">
        <span>{{ section.title }}</span>
        <span class="text-sm">Click to expand</span>
      </h4>

      <div class="animate-fadeIn" v-show="activeSection === index">
        <p class="mb-2">{{ section.content }}</p>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { ref } from 'vue'

const activeSection = ref<number>()

const ruleSections = [
  {
    title: '🎯 Objective',
    content:
      'The objective of backgammon is to move all your checkers into your home board and then bear them off. The first player to remove all their checkers wins the game.',
  },
  {
    title: '⚫ Moving Checkers',
    content: `Players take turns rolling two dice and moving their checkers counterclockwise around the board. Each die roll allows you to move one or two different checkers:
     The numbers on each die represent separate moves. For example, if you roll 5-3, you can move one checker 5 spaces and another (or the same) checker 3 spaces. If you roll doubles (same number on both dice), you have to play each number twice. For example, if you roll 6-6, you can make four moves of 6 spaces each.
`,
  },
  {
    title: '⚔️ Hitting and Entering',
    content:
      "When a point is occupied by a single checker, it is called a 'blot' and can be hit by an opponent's checker. If your checker is hit, it must be placed on the bar and re-entered before any other move can be made.",
  },
  {
    title: '🏆 Bearing Off',
    content:
      "Once all of a player's checkers are in their home board, they can start bearing off. This is done by moving checkers off the board according to the dice rolls.",
  },
]
</script>

<style scoped>
.rule-section {
  @apply rounded-lg border-2 border-primary/35 bg-primary/5 p-4 transition-all duration-300 ease-in-out;
}

.rule-section:hover {
  @apply border-primary;
  transform: translateX(7px);
}

.rule-section.active {
  @apply border-primary bg-primary/10;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(-10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.animate-fadeIn {
  animation: fadeIn 0.4s ease-in-out;
}
</style>
