<script setup lang="ts">
import { ref } from 'vue'

type BadgeCategory = 'bot' | 'elo' | 'pieces' | 'played' | 'time' | 'win'

interface Badges {
  bot: number
  elo: number
  pieces: number
  played: number
  time: number
  win: number
}

const props = defineProps<{
  badges: Badges | undefined
}>()

const hoveredBadge = ref<{ category: BadgeCategory; level: number } | null>(
  null,
)

const badgeCategories: Record<BadgeCategory, string> = {
  bot: '🤖 Bot Difficulties',
  elo: '📈 ELO Rating',
  pieces: '♟️ Checker Home',
  played: '🎲 Games Played',
  time: '⏰ Fastest Win',
  win: '🏆 Victories',
}

const getBadgeDescription = (category: string, level: number) => {
  const descriptions: Record<string, string[]> = {
    bot: [
      'Beat your first bot! 🤖',
      'Bot Challenger! 🤖⚡',
      'Supreme Bot Master! 🤖👑',
    ],
    elo: ['Rising Star! ⭐', 'Elite Player! ⭐⚡', 'Legendary Champion! 👑'],
    pieces: ['50! 👣', '100! 🎯', '200! 🎮'],
    played: [
      'Rookie Player (1)! 🎲',
      'Dedicated Player (10)! 🎲🎯',
      'Veteran Player (100)! 🏆',
    ],
    time: [
      'Win in 10 minutes! ⌛',
      'Win in 5 minutes! ⏰',
      'Win in 3 minutes! ⚡',
    ],
    win: ['First Victory! 🎉', '10 Victory! 🔥', '50 Victory! 👑'],
  }
  return descriptions[category]?.[level - 1]
}

</script>
<template>
  <div class="grid grid-cols-2 gap-3 md:grid-cols-3">
    <div
      v-for="(value, category) in badges"
      :key="category"
      class="badge-category"
    >
      <h3 class="mb-2 text-center font-bold text-primary">
        {{ badgeCategories[category] }}
      </h3>

      <div class="flex flex-wrap justify-center gap-4">
        <div
          v-for="level in 3"
          :key="`${category}_${level}`"
          class="cursor-pointer"
          @mouseenter="hoveredBadge = { category, level }"
          @mouseleave="hoveredBadge = null"
          @click="openBadgeModal(category, level)"
        >
          <div class="relative">
            <img
              :src="`/badges/${category}_${level}.png`"
              :alt="`${category} badge level ${level}`"
              class="badge-image"
              :class="{
                unearned: level > value,
                'hover:scale-[1.4]': level <= value,
              }"
            />

            <div
              v-if="
                hoveredBadge?.category === category &&
                hoveredBadge?.level === level &&
                level <= value
              "
              class="badge-tooltip"
            >
              {{ getBadgeDescription(category, level) }}
            </div>
          </div>
        </div>
      </div>

      <div class="mt-2 text-center">
        <span class="badge badge-primary badge-md"> {{ value }}/3 </span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.badge-category {
  @apply rounded-xl border-0 border-primary bg-base-200 p-4 transition-all duration-200;
}
.badge-category:hover {
  @apply border-2 border-primary bg-base-300 shadow-xl;
}

.category-highlighted {
  @apply border-2 border-primary bg-base-300;
  transform: translateY(-4px);
  box-shadow: 0 10px 20px rgba(0, 0, 0, 0.2);
}

.badge-image {
  @apply h-20 w-20 transition-all duration-300;
}

.badge-image.unearned {
  @apply opacity-30 grayscale;
}

.badge-tooltip {
  @apply absolute -top-12 left-1/2 -translate-x-1/2 rounded-lg bg-primary p-2 text-sm text-primary-content;
  min-width: max-content;
  animation: fadeIn 0.2s ease-in-out;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translate(-50%, -20px);
  }
  to {
    opacity: 1;
    transform: translate(-50%, 0);
  }
}
</style>
