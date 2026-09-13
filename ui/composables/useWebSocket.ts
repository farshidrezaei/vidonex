import { useRenderStore } from '~/stores/render'

export function useWebSocket() {
  const renderStore = useRenderStore()
  const isConnected = ref(false)
  let socket: WebSocket | null = null
  let reconnectTimeout: ReturnType<typeof setTimeout> | null = null

  function connect() {
    if (import.meta.server) return

    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const host = window.location.host
    const wsUrl = `${protocol}//${host}/ws`

    try {
      socket = new WebSocket(wsUrl)

      socket.onopen = () => {
        isConnected.value = true
        console.log('[Vidonex WS] Connected to server hub')
      }

      socket.onmessage = (event) => {
        try {
          const envelope = JSON.parse(event.data)
          switch (envelope.type) {
            case 'render_progress':
              renderStore.handleProgressEvent(envelope.payload)
              break
            case 'render_complete':
              renderStore.handleCompleteEvent(envelope.payload)
              break
            case 'render_failed':
              renderStore.handleFailedEvent(envelope.payload)
              break
            case 'render_cancelled':
              renderStore.isRendering = false
              break
          }
        } catch (err) {
          console.error('[Vidonex WS] Error parsing message', err)
        }
      }

      socket.onclose = () => {
        isConnected.value = false
        console.log('[Vidonex WS] Disconnected, scheduling reconnect...')
        if (reconnectTimeout) clearTimeout(reconnectTimeout)
        reconnectTimeout = setTimeout(connect, 3000)
      }

      socket.onerror = (err) => {
        console.warn('[Vidonex WS] Socket error', err)
        socket?.close()
      }
    } catch (err) {
      console.error('[Vidonex WS] Connection error', err)
    }
  }

  function disconnect() {
    if (reconnectTimeout) clearTimeout(reconnectTimeout)
    socket?.close()
  }

  return {
    isConnected,
    connect,
    disconnect,
  }
}
