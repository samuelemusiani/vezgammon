<template>
  <div class="fixed bottom-4 right-4 z-50">
    <!-- Chat Button -->
    <button
      @click="toggleChat"
      class="retro-button circle bg-amber-700 hover:bg-amber-800"
    >
      <i class="fas fa-comments"></i>

      <div
        v-if="unreadMessages"
        class="absolute -right-2 -top-2 h-5 w-5 rounded-full bg-red-500 text-xs text-white"
      >
        {{ unreadMessages }}
      </div>
    </button>

    <!-- Chat Window -->
    <div
      v-if="isOpen"
      class="retro-box absolute bottom-16 right-0 h-96 w-80 overflow-hidden bg-white"
    >
      <!-- Chat Header -->
      <div
        class="flex h-12 items-center justify-between border-b-2 border-amber-900 bg-amber-800 px-4"
      >
        <h3 class="font-bold text-white">Game Chat</h3>
        <button @click="toggleChat" class="text-white hover:text-gray-300">
          <i class="fas fa-times"></i>
        </button>
      </div>

      <!-- Messages Container -->
      <div
        ref="messagesContainer"
        class="h-[calc(100%-6rem)] overflow-y-auto p-4"
      >
        <div v-for="(msg, index) in messages" :key="index" class="mb-2">
          <div
            :class="[
              'flex',
              msg.sender === myUsername ? 'justify-end' : 'justify-start',
            ]"
          >
            <div
              :class="[
                'max-w-[80%] rounded-lg px-4 py-2',
                msg.sender === myUsername
                  ? 'bg-amber-700 text-white'
                  : 'bg-amber-100 text-amber-900',
              ]"
            >
              <div class="mb-1 text-xs opacity-50">
                {{ msg.sender }}
              </div>
              {{ msg.payload }}
            </div>
          </div>
        </div>
      </div>

      <!-- Input Area -->
      <div
        class="absolute bottom-0 left-0 flex h-12 w-full border-t-2 border-amber-900"
      >
        <input
          v-model="newMessage"
          type="text"
          placeholder="Type a message..."
          class="flex-1 bg-amber-50 px-4"
          @keyup.enter="sendMessage"
        />
        <button
          @click="sendMessage"
          class="w-12 bg-amber-700 text-white hover:bg-amber-800"
        >
          <i class="fas fa-paper-plane"></i>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick, watch } from 'vue'
import '@fortawesome/fontawesome-free/css/all.min.css'

import { useWebSocketStore } from '@/stores/websocket'
import tinSfx from '@/utils/sounds/tintin.mp3'
import type { WSMessage } from '@/utils/types'
import { useSound } from '@vueuse/sound'

const props = defineProps<{
  myUsername: string
  opponentUsername: string
}>()

const { play: playTin } = useSound(tinSfx, { volume: 0.5 })

const webSocketStore = useWebSocketStore()
const isOpen = ref(false)
const messages = ref<{ sender: string; payload: string }[]>([])
const newMessage = ref('')
const unreadMessages = ref(0)
const messagesContainer = ref<HTMLElement | null>(null)

const toggleChat = () => {
  isOpen.value = !isOpen.value
  if (isOpen.value) {
    unreadMessages.value = 0
    nextTick(() => {
      scrollToBottom()
    })
  }
}

const scrollToBottom = () => {
  if (messagesContainer.value) {
    messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
  }
}

const sendMessage = () => {
  if (!newMessage.value.trim()) return

  // Invia il messaggio tramite WebSocket
  webSocketStore.sendMessage({
    type: 'chat_message',
    payload: newMessage.value,
  })
  // Aggiungi il messaggio localmente
  messages.value.push({
    sender: props.myUsername,
    payload: newMessage.value,
  })

  newMessage.value = ''
  nextTick(() => {
    scrollToBottom()
  })
}

// Gestione dei messaggi in arrivo
const handleIncomingMessage = (message: WSMessage) => {
  if (message.type === 'chat_message') {
    messages.value.push({
      sender: props.opponentUsername,
      payload: message.payload,
    })
    if (!isOpen.value) {
      playTin()
      unreadMessages.value++
    }
    nextTick(() => {
      scrollToBottom()
    })
  }
}

onMounted(() => {
  webSocketStore.addMessageHandler(handleIncomingMessage)
})

watch(messages, () => {
  nextTick(() => {
    scrollToBottom()
  })
})
</script>

<style scoped>
.retro-button {
  @apply btn;
  background: #d2691e;
  color: white;
  border: 3px solid #8b4513;
  font-family: 'Arial Black', serif;
  text-transform: uppercase;
  text-shadow: 2px 2px 0 rgba(0, 0, 0, 0.2);
  box-shadow: 0 2px 0 #8b4513;
  font-size: 1.1rem;
  height: 6vh;

  &.circle {
    width: 60px;
    height: 60px;
    border-radius: 50%;
  }

  &:hover {
    transform: translateY(2px);
    box-shadow:
      inset 0 0 10px rgba(0, 0, 0, 0.2),
      0 0px 0 #8b4513;
    cursor: url('/tortellino.png'), auto;
  }
}

.retro-box {
  background-color: #ffe5c9;
  border: 5px solid #8b4513;
  box-shadow:
    0 0 0 4px #d2691e,
    inset 0 0 20px rgba(0, 0, 0, 0.2);
}
</style>
