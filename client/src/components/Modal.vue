<script setup lang="ts">
defineProps<{
  show: boolean
  title: string
  confirmText?: string
  cancelText?: string
  confirmVariant?: 'success' | 'danger' | 'primary'
}>()

const emit = defineEmits<{
  (e: 'confirm'): void
  (e: 'cancel'): void
}>()
</script>

<template>
  <div
    v-if="show"
    class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50"
  >
    <div class="retro-box max-w-md rounded-lg p-6 text-center">
      <h3 class="mb-4 text-xl font-bold text-amber-900">
        {{ title }}
      </h3>

      <div class="mb-6 text-gray-700">
        <slot></slot>
      </div>

      <div class="flex justify-center gap-4">
        <button
          v-if="confirmText"
          @click="emit('confirm')"
          class="retro-button"
          :class="{
            'bg-green-700 hover:bg-green-800': confirmVariant === 'success',
            'bg-red-700 hover:bg-red-800': confirmVariant === 'danger',
          }"
        >
          {{ confirmText }}
        </button>
        <button v-if="cancelText" @click="emit('cancel')" class="retro-button">
          {{ cancelText }}
        </button>
      </div>
    </div>
  </div>
</template>
