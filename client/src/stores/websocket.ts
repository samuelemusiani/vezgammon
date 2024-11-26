import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useWebSocketStore = defineStore('websocket', () => {
  const socket = ref<WebSocket | null>(null)
  const isConnected = ref(false)
  const lastMessage = ref<string | null>(null)

  const connect = () => {
    try {
      if (isConnected.value) return // Already connected
      socket.value = new WebSocket('ws://130.136.201.194:8080/api/ws')
      isConnected.value = true

      socket.value.onopen = () => {
        isConnected.value = true
        console.log('WebSocket connected')
      }

      socket.value.onclose = () => {
        isConnected.value = false
        console.log('WebSocket disconnected')
      }

      socket.value.onerror = error => {
        console.error('WebSocket error:', error)
      }

      socket.value.onmessage = event => {
        lastMessage.value = event.data
        console.log('Received message:', event.data)
      }
    } catch (error) {
      console.error('Error connecting to WebSocket:', error)
    }
  }

  const disconnect = () => {
    if (socket.value) {
      socket.value.close()
      socket.value = null
    }
  }

  const sendMessage = (message: any) => {
    if (socket.value && isConnected.value) {
      socket.value.send(JSON.stringify(message))
    }
  }

  return {
    socket,
    isConnected,
    lastMessage,
    connect,
    disconnect,
    sendMessage,
  }
})
