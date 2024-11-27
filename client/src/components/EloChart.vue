<script setup lang="ts">
import { ref, onMounted } from 'vue'
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
  PointElement
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
      className: 'border-primary bg-secondary',
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
  <div class="p-2 borer-primary mt-2 border-8 bg-">
    <Line :data="chartData" :options="chartOptions" />
  </div>
</template>
