<script setup lang="ts">
import { watch, computed } from 'vue'
import { Line } from 'vue-chartjs'
import {
  Chart as ChartJS,
  Title,
  Tooltip,
  Legend,
  LineElement,
  CategoryScale,
  LinearScale,
  PointElement,
} from 'chart.js'

import type { TooltipItem } from 'chart.js'

ChartJS.register(
  Title,
  Tooltip,
  Legend,
  LineElement,
  CategoryScale,
  LinearScale,
  PointElement,
)

const props = defineProps<{
  elo: number[]
}>()

const chartData = computed(() => ({
  labels: props.elo.map((_, index) => `${index + 1}`),
  datasets: [
    {
      label: 'Elo Rating',
      data: props.elo,
      fill: true,
      borderColor: 'rgb(75, 192, 192)',
      backgroundColor: 'rgba(75, 192, 192, 0.2)',
      tension: 0.4,
      pointBackgroundColor: 'rgb(75, 192, 192)',
      pointBorderColor: '#fff',
      pointBorderWidth: 0,
      pointRadius: 2,
      pointHoverRadius: 6,
    },
  ],
}))

const chartOptions = computed(() => ({
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      display: true,
      position: 'top' as const,
      labels: {
        color: 'rgb(75, 192, 192)',
        font: {
          size: 14,
        },
      },
    },
    tooltip: {
      backgroundColor: 'rgba(0, 0, 0, 0.8)',
      titleColor: 'rgb(75, 192, 192)',
      bodyColor: '#fff',
      displayColors: false,
      callbacks: {
        label: (context: TooltipItem<'line'>) => `Elo: ${context.parsed.y}`,
      },
    },
  },
  scales: {
    x: {
      title: {
        display: true,
        text: 'Games',
        color: 'rgb(75, 192, 192)',
        font: {
          size: 16,
        },
      },
      grid: {
        display: false,
      },
      ticks: {
        color: 'rgb(75, 192, 192)',
      },
    },
    y: {
      title: {
        display: true,
        text: 'Elo Rating',
        color: 'rgb(75, 192, 192)',
        font: {
          size: 16,
        },
      },
      grid: {
        color: 'rgba(75, 192, 192, 0.1)',
      },
      ticks: {
        color: 'rgb(75, 192, 192)',
      },
      beginAtZero: false,
    },
  },
}))

watch(
  () => props.elo,
  newElo => {
    console.log('Elo data updated:', newElo)
  },
  { deep: true },
)
</script>

<template>
  <div class="h-[300px] rounded-lg border-2 border-primary bg-base-200 p-3">
    <Line :data="chartData" :options="chartOptions" />
  </div>
</template>
