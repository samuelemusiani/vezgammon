<script setup lang="ts">
import { ref } from 'vue'
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

const chartData = ref({
  labels: props.elo.map((_, index) => `Game ${index + 1}`),
  datasets: [
    {
      label: 'Elo Rating',
      data: props.elo,
      fill: true,
      borderColor: 'rgb(75, 192, 192)',
      backgroundColor: 'rgba(75, 192, 192, 0.2)',
    },
  ],
})

const chartOptions = ref({
  responsive: true,
  maintainAspectRatio: false,
  scales: {
    x: {
      title: {
        display: true,
        text: 'Games',
      },
    },
    y: {
      title: {
        display: true,
        text: 'Elo Rating',
      },
      beginAtZero: false,
    },
  },
})
</script>

<template>
  <div class="mt-2 h-[300px] border-8 border-primary bg-base-200 p-2">
    <Line :data="chartData" :options="chartOptions" />
  </div>
</template>
