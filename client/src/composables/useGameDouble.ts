import { ref } from 'vue'
import { vfetch } from '@/utils/fetch'

export function useGameDouble(onDecline?: () => void) {
  const showDoubleModal = ref(false)

  const handleDouble = async () => {
    try {
      const res = await vfetch('/api/play/double', {
        method: 'POST',
      })
      return res.ok
    } catch (err) {
      console.error('Error sending double:', err)
      return false
    }
  }

  const acceptDouble = async () => {
    try {
      const res = await vfetch('/api/play/double', {
        method: 'PUT',
      })
      return res.ok
    } catch (err) {
      console.error('Error accepting double:', err)
      return false
    }
  }

  const declineDouble = async () => {
    try {
      const res = await vfetch('/api/play/double', {
        method: 'DELETE',
      })
      if (res.ok) {
        showDoubleModal.value = false
        if (onDecline) onDecline()
      }

      return res.ok
    } catch (err) {
      console.error('Error declining double:', err)
      return false
    }
  }

  return {
    showDoubleModal,
    handleDouble,
    acceptDouble,
    declineDouble,
  }
}
