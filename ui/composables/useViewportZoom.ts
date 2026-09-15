export const zoom = ref(1.0)
export const panX = ref(0)
export const panY = ref(0)
export const isSpacePressed = ref(false)
export const isPanning = ref(false)
export const isFullscreen = ref(false)

export const ZOOM_PRESETS = [
  { label: 'Fit', value: 1.0, isFit: true },
  { label: '25%', value: 0.25 },
  { label: '50%', value: 0.5 },
  { label: '75%', value: 0.75 },
  { label: '100%', value: 1.0 },
  { label: '150%', value: 1.5 },
  { label: '200%', value: 2.0 },
]

export function resetViewportZoom() {
  zoom.value = 1.0
  panX.value = 0
  panY.value = 0
}

if (typeof window !== 'undefined') {
  window.addEventListener('fullscreenchange', () => {
    isFullscreen.value = !!document.fullscreenElement
  })
}

export function useViewportZoom() {
  const isFit = computed(() => zoom.value === 1.0 && panX.value === 0 && panY.value === 0)
  const zoomDisplay = computed(() => {
    if (isFit.value) return 'Fit'
    return `${Math.round(zoom.value * 100)}%`
  })

  function setZoom(level: number, resetPan = false) {
    zoom.value = Math.min(5.0, Math.max(0.1, Number(level.toFixed(2))))
    if (resetPan) {
      panX.value = 0
      panY.value = 0
    }
  }

  function zoomIn() {
    setZoom(zoom.value * 1.25)
  }

  function zoomOut() {
    setZoom(zoom.value / 1.25)
  }

  function resetZoom() {
    zoom.value = 1.0
    panX.value = 0
    panY.value = 0
  }

  function handleWheel(event: WheelEvent) {
    event.preventDefault()
    event.stopPropagation()

    // Smooth continuous zooming with mouse wheel
    const zoomFactor = event.deltaY < 0 ? 1.15 : 0.87
    const newZoom = Math.min(5.0, Math.max(0.1, zoom.value * zoomFactor))
    zoom.value = Number(newZoom.toFixed(2))
  }

  function startPan(event: MouseEvent) {
    event.preventDefault()
    event.stopPropagation()

    isPanning.value = true
    const startX = event.clientX
    const startY = event.clientY
    const initialPanX = panX.value
    const initialPanY = panY.value

    const onMouseMove = (moveEvent: MouseEvent) => {
      panX.value = Math.round(initialPanX + (moveEvent.clientX - startX))
      panY.value = Math.round(initialPanY + (moveEvent.clientY - startY))
    }

    const onMouseUp = () => {
      isPanning.value = false
      window.removeEventListener('mousemove', onMouseMove)
      window.removeEventListener('mouseup', onMouseUp)
    }

    window.addEventListener('mousemove', onMouseMove)
    window.addEventListener('mouseup', onMouseUp)
  }

  async function toggleFullscreen(targetElement?: HTMLElement | null) {
    try {
      if (!document.fullscreenElement) {
        const el = targetElement || document.documentElement
        if (el.requestFullscreen) {
          await el.requestFullscreen()
        }
      } else {
        if (document.exitFullscreen) {
          await document.exitFullscreen()
        }
      }
    } catch (err) {
      console.error('Fullscreen toggle failed:', err)
    }
  }

  return {
    zoom,
    panX,
    panY,
    isSpacePressed,
    isPanning,
    isFullscreen,
    isFit,
    zoomDisplay,
    ZOOM_PRESETS,
    setZoom,
    zoomIn,
    zoomOut,
    resetZoom,
    handleWheel,
    startPan,
    toggleFullscreen,
  }
}
