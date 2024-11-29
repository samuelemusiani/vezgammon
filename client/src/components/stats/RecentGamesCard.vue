<template>
  <div class="card bg-base-200">
    <div class="card-body">
      <h3 class="card-title">Recent Games</h3>
      <div class="overflow-x-auto">
        <table class="table">
          <tbody>
            <tr v-for="game in games" :key="game.id" class="text-xl">
              <td>{{ game.player1 }} vs {{ game.player2 }}</td>
              <td>
                <span
                  :class="
                    game.status === game.current_player
                      ? 'text-success'
                      : 'text-error'
                  "
                >
                  {{ game.status }}
                </span>
              </td>
              <td class="text-right">
                {{ new Date(game.start).toDateString() }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { GameState } from '@/utils/game/types'

defineProps<{
  games: GameState[]
}>()
</script>
