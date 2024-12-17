import type { Badge } from '@/utils/types'
import { ref } from 'vue'
import { defineStore } from 'pinia'

export const useBadgesStore = defineStore('badges', () => {
  const badges = ref<Badge>({
    Bot: [],
    Homepieces: [],
    Wongames: [],
    Elo: [],
    Wontime: [],
    Gameplayed: [],
  })

  const firstTime = ref(true)
  const badgeChanged = ref(false)

  const fetchBadges = async () => {
    try {
      const response = await fetch('/api/badge')
      if (!response.ok) {
        throw new Error('Failed to fetch badges')
      }
      const data = await response.json()
      if (firstTime.value) {
        firstTime.value = false
      } else if (sumBadge(badges.value) !== sumBadge(data)) {
        badgeChanged.value = true
      }
      badges.value = data
    } catch (error) {
      console.error('Error fetching badges:', error)
    }
  }

  const haveBagesChanged = () => {
    if (badgeChanged.value) {
      badgeChanged.value = false
      return true
    }
    return false
  }

  function sumBadge(badge: Badge): number {
    return Object.values(badge).reduce((acc, val) => acc + val, 0)
  }

  return { badges, fetchBadges, haveBagesChanged }
})
