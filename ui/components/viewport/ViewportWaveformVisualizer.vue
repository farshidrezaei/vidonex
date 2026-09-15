<template>
  <div
    class="absolute select-none overflow-visible group"
    :class="isSelected ? 'pointer-events-auto ring-1.5 ring-indigo-400/90 shadow-[0_0_15px_rgba(99,102,241,0.3)]' : 'pointer-events-auto cursor-pointer hover:ring-1 hover:ring-indigo-400/50'"
    :style="containerStyle"
    @mousedown="handleMouseDown"
  >
    <!-- Live Waveform Canvas -->
    <canvas
      ref="canvasRef"
      :width="canvasResolution.width"
      :height="canvasResolution.height"
      class="w-full h-full block"
    ></canvas>

    <!-- Transform Handles & Bounding Box (Shown when selected) -->
    <template v-if="isSelected">
      <!-- Top-Left -->
      <div
        class="absolute -top-1.5 -left-1.5 w-3 h-3 bg-white border-2 border-indigo-600 rounded-xs cursor-nwse-resize z-30 shadow-md hover:scale-125 transition-transform"
        @mousedown.stop="startResize('tl', $event)"
      ></div>
      <!-- Top-Center -->
      <div
        class="absolute -top-1.5 left-1/2 -translate-x-1/2 w-3 h-3 bg-white border-2 border-indigo-600 rounded-xs cursor-ns-resize z-30 shadow-md hover:scale-125 transition-transform"
        @mousedown.stop="startResize('tc', $event)"
      ></div>
      <!-- Top-Right -->
      <div
        class="absolute -top-1.5 -right-1.5 w-3 h-3 bg-white border-2 border-indigo-600 rounded-xs cursor-nesw-resize z-30 shadow-md hover:scale-125 transition-transform"
        @mousedown.stop="startResize('tr', $event)"
      ></div>
      <!-- Middle-Right -->
      <div
        class="absolute top-1/2 -right-1.5 -translate-y-1/2 w-3 h-3 bg-white border-2 border-indigo-600 rounded-xs cursor-ew-resize z-30 shadow-md hover:scale-125 transition-transform"
        @mousedown.stop="startResize('mr', $event)"
      ></div>
      <!-- Bottom-Right -->
      <div
        class="absolute -bottom-1.5 -right-1.5 w-3 h-3 bg-white border-2 border-indigo-600 rounded-xs cursor-nwse-resize z-30 shadow-md hover:scale-125 transition-transform"
        @mousedown.stop="startResize('br', $event)"
      ></div>
      <!-- Bottom-Center -->
      <div
        class="absolute -bottom-1.5 left-1/2 -translate-x-1/2 w-3 h-3 bg-white border-2 border-indigo-600 rounded-xs cursor-ns-resize z-30 shadow-md hover:scale-125 transition-transform"
        @mousedown.stop="startResize('bc', $event)"
      ></div>
      <!-- Bottom-Left -->
      <div
        class="absolute -bottom-1.5 -left-1.5 w-3 h-3 bg-white border-2 border-indigo-600 rounded-xs cursor-nesw-resize z-30 shadow-md hover:scale-125 transition-transform"
        @mousedown.stop="startResize('bl', $event)"
      ></div>
      <!-- Middle-Left -->
      <div
        class="absolute top-1/2 -left-1.5 -translate-y-1/2 w-3 h-3 bg-white border-2 border-indigo-600 rounded-xs cursor-ew-resize z-30 shadow-md hover:scale-125 transition-transform"
        @mousedown.stop="startResize('ml', $event)"
      ></div>

      <!-- Quick Dimension Tag Overlay -->
      <div
        class="absolute -bottom-6 left-1/2 -translate-x-1/2 px-1.5 py-0.5 rounded bg-gray-900/90 border border-gray-700 text-[9px] font-mono text-gray-300 pointer-events-none whitespace-nowrap shadow"
      >
        {{ Math.round(targetDimensions.width) }} × {{ Math.round(targetDimensions.height) }} px
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { usePlaybackStore } from '~/stores/playback'
import { useTimelineStore } from '~/stores/timeline'
import {
  showGuideCenterX,
  showGuideCenterY,
} from '~/composables/useTransformGizmo'
import type { TrackSpec, ClipSpec, WaveformSpec, PositionSpec } from '~/types/spec'

interface Props {
  track: TrackSpec
  clip?: ClipSpec | null
  canvasWidth: number
  canvasHeight: number
  displayWidth: number
  displayHeight: number
}

const props = defineProps<Props>()
const playbackStore = usePlaybackStore()
const timelineStore = useTimelineStore()
const canvasRef = ref<HTMLCanvasElement | null>(null)

// Check if this waveform element is currently selected
const isSelected = computed(() => {
  if (props.clip) {
    return timelineStore.selectedClipId === props.clip.id
  }
  return timelineStore.selectedTrackId === props.track.id
})

// Active configuration from clip or track
const waveformConfig = computed<WaveformSpec>(() => {
  return props.clip?.waveform || props.track.waveform || {}
})

// Mode and visual styling
const mode = computed(() => (waveformConfig.value.mode || 'peak_to_peak').toLowerCase())
const primaryColor = computed(() => waveformConfig.value.color || '#00E6B4')
const secondaryColor = computed(() => waveformConfig.value.secondary_color || '#6366F1')
const density = computed(() => Number(waveformConfig.value.bar_density) || 48)
const roundness = computed(() => Number(waveformConfig.value.bar_roundness) ?? 6)
const isGlow = computed(() => waveformConfig.value.glow ?? true)
const glowRadius = computed(() => Number(waveformConfig.value.glow_radius) || 8)
const lineWidth = computed(() => Number(waveformConfig.value.line_width) || 3)
const isFillArea = computed(() => waveformConfig.value.fill_area ?? false)
const innerRadiusRatio = computed(() => Number(waveformConfig.value.inner_radius) || 0.45)
const opacity = computed(() => props.clip?.opacity ?? (waveformConfig.value.opacity ?? 1.0))

// Scale factor from canvas resolution to display stage
const displayScale = computed(() => {
  if (props.canvasWidth <= 0) return 1
  return props.displayWidth / props.canvasWidth
})

// Natural canvas pixel dimensions
const targetDimensions = computed(() => {
  const w = Number(waveformConfig.value.width) || Math.round(props.canvasWidth * 0.75)
  const h = Number(waveformConfig.value.height) || (mode.value === 'circular' ? Math.min(w, 400) : 180)
  return {
    width: Math.max(120, w),
    height: Math.max(50, h),
  }
})

// HiDPI internal canvas resolution for razor-sharp rendering
const canvasResolution = computed(() => {
  const scale = 2
  return {
    width: targetDimensions.value.width * scale,
    height: targetDimensions.value.height * scale,
  }
})

const position = computed<PositionSpec>(() => {
  return props.clip?.position || waveformConfig.value.position || { x: 0, y: 0, alignment: 'center' }
})

const containerStyle = computed(() => {
  const baseScale = displayScale.value
  const visualW = targetDimensions.value.width * baseScale
  const visualH = targetDimensions.value.height * baseScale

  const posX = (position.value.x ?? 0) * baseScale
  const posY = (position.value.y ?? 0) * baseScale

  const left = (props.displayWidth - visualW) / 2 + posX
  const top = (props.displayHeight - visualH) / 2 + posY

  return {
    left: `${Math.round(left)}px`,
    top: `${Math.round(top)}px`,
    width: `${Math.round(visualW)}px`,
    height: `${Math.round(visualH)}px`,
    opacity: opacity.value,
    zIndex: (props.track.z_index ?? 10) + (isSelected.value ? 10 : 0),
    cursor: isSelected.value ? 'move' : 'pointer',
  }
})

// Helper to mutate configuration on clip or track
function updateWaveformSpec(patch: Partial<WaveformSpec>) {
  if (props.clip) {
    if (!props.clip.waveform) {
      props.clip.waveform = {
        mode: 'peak_to_peak',
        color: '#00E6B4',
        secondary_color: '#6366F1',
        width: targetDimensions.value.width,
        height: targetDimensions.value.height,
        position: { x: position.value.x, y: position.value.y, alignment: 'center' },
      }
    }
    Object.assign(props.clip.waveform, patch)
  } else {
    if (!props.track.waveform) {
      props.track.waveform = {
        mode: 'peak_to_peak',
        color: '#00E6B4',
        secondary_color: '#6366F1',
        width: targetDimensions.value.width,
        height: targetDimensions.value.height,
        position: { x: position.value.x, y: position.value.y, alignment: 'center' },
      }
    }
    Object.assign(props.track.waveform, patch)
  }
}

function updatePosition(x: number, y: number) {
  if (props.clip) {
    if (!props.clip.position) {
      props.clip.position = { x: 0, y: 0, alignment: 'center' }
    }
    props.clip.position.x = x
    props.clip.position.y = y
  }
  updateWaveformSpec({
    position: { x, y, alignment: 'center' },
  })
}

// -------------------------------------------------------------
// Interactive Direct Manipulation: Drag to Move & Handle Resize
// -------------------------------------------------------------
function handleMouseDown(event: MouseEvent) {
  event.stopPropagation()
  event.preventDefault()

  if (props.clip) {
    timelineStore.selectedClipId = props.clip.id
  } else {
    timelineStore.selectedTrackId = props.track.id
  }

  // Start dragging position
  const startX = event.clientX
  const startY = event.clientY
  const initialPosX = position.value.x ?? 0
  const initialPosY = position.value.y ?? 0
  const scale = displayScale.value

  function onMouseMove(moveEvent: MouseEvent) {
    const deltaX = (moveEvent.clientX - startX) / scale
    const deltaY = (moveEvent.clientY - startY) / scale

    let nextX = Math.round(initialPosX + deltaX)
    let nextY = Math.round(initialPosY + deltaY)

    // Snap to center
    if (Math.abs(nextX) < 8) {
      nextX = 0
      showGuideCenterX.value = true
    } else {
      showGuideCenterX.value = false
    }
    if (Math.abs(nextY) < 8) {
      nextY = 0
      showGuideCenterY.value = true
    } else {
      showGuideCenterY.value = false
    }

    updatePosition(nextX, nextY)
  }

  function onMouseUp() {
    window.removeEventListener('mousemove', onMouseMove)
    window.removeEventListener('mouseup', onMouseUp)
    showGuideCenterX.value = false
    showGuideCenterY.value = false
    timelineStore.pushHistoryState('Move Waveform')
  }

  window.addEventListener('mousemove', onMouseMove)
  window.addEventListener('mouseup', onMouseUp)
}

function startResize(handle: string, event: MouseEvent) {
  event.stopPropagation()
  event.preventDefault()

  const startX = event.clientX
  const startY = event.clientY
  const startW = targetDimensions.value.width
  const startH = targetDimensions.value.height
  const initialPosX = position.value.x ?? 0
  const initialPosY = position.value.y ?? 0
  const scale = displayScale.value

  function onResizeMove(moveEvent: MouseEvent) {
    const deltaX = (moveEvent.clientX - startX) / scale
    const deltaY = (moveEvent.clientY - startY) / scale

    let newW = startW
    let newH = startH
    let newPosX = initialPosX
    let newPosY = initialPosY

    if (handle.includes('r')) {
      newW = Math.max(120, startW + deltaX)
      newPosX = initialPosX + deltaX / 2
    }
    if (handle.includes('l')) {
      newW = Math.max(120, startW - deltaX)
      newPosX = initialPosX + deltaX / 2
    }
    if (handle.includes('b')) {
      newH = Math.max(50, startH + deltaY)
      newPosY = initialPosY + deltaY / 2
    }
    if (handle.includes('t')) {
      newH = Math.max(50, startH - deltaY)
      newPosY = initialPosY + deltaY / 2
    }

    updatePosition(Math.round(newPosX), Math.round(newPosY))
    updateWaveformSpec({
      width: Math.round(newW),
      height: Math.round(newH),
    })
  }

  function onResizeUp() {
    window.removeEventListener('mousemove', onResizeMove)
    window.removeEventListener('mouseup', onResizeUp)
    timelineStore.pushHistoryState('Resize Waveform')
  }

  window.addEventListener('mousemove', onResizeMove)
  window.addEventListener('mouseup', onResizeUp)
}

// -------------------------------------------------------------
// Live Waveform Canvas Render Engine (5 Modes)
// -------------------------------------------------------------
let animationFrameId: number | null = null

function drawWaveform() {
  const canvas = canvasRef.value
  if (!canvas) return
  const ctx = canvas.getContext('2d')
  if (!ctx) return

  const width = canvas.width
  const height = canvas.height
  const halfHeight = height / 2

  ctx.clearRect(0, 0, width, height)
  ctx.save()

  const currentMode = mode.value
  const pColor = primaryColor.value
  const sColor = secondaryColor.value
  const isPlaying = playbackStore.isPlaying
  const currentTime = playbackStore.currentTime

  // 1. Dual-Color Gradient
  const gradient = ctx.createLinearGradient(0, 0, width, height)
  gradient.addColorStop(0, pColor)
  gradient.addColorStop(1, sColor)

  if (isGlow.value) {
    ctx.shadowColor = pColor
    ctx.shadowBlur = glowRadius.value * 2
  }

  if (currentMode === 'bars' || currentMode === 'p2p' || currentMode === 'peak_to_peak') {
    // Mode 1: Mirrored Peak-to-Peak Bipolar Bars
    const barCount = Math.max(16, Math.min(128, density.value))
    const barSpacing = width / barCount
    const barWidth = Math.max(3, barSpacing * 0.65)
    const radius = Math.min(barWidth / 2, roundness.value * 2)

    ctx.fillStyle = gradient

    for (let i = 0; i < barCount; i++) {
      const seed = i * 0.35 + currentTime * 8
      const harmonic = Math.abs(Math.sin(seed) * 0.5 + Math.cos(seed * 2.1) * 0.3 + Math.sin(seed * 0.7) * 0.2)
      const amplitude = isPlaying ? Math.max(0.08, harmonic) : Math.max(0.06, Math.sin(i * 0.2 + currentTime) * 0.3 + 0.4)

      const barHeight = amplitude * (height * 0.85)
      const x = i * barSpacing + (barSpacing - barWidth) / 2
      const y = halfHeight - barHeight / 2

      ctx.beginPath()
      ctx.roundRect(x, y, barWidth, Math.max(4, barHeight), radius)
      ctx.fill()
    }
  } else if (currentMode === 'spectrum') {
    // Mode 2: Bottom-aligned Equalizer Spectrum with Floating Peak Caps
    const barCount = Math.max(16, Math.min(96, density.value))
    const barSpacing = width / barCount
    const barWidth = Math.max(3, barSpacing * 0.7)
    const radius = Math.min(barWidth / 2, roundness.value * 2)

    ctx.fillStyle = gradient

    for (let i = 0; i < barCount; i++) {
      const seed = i * 0.4 + currentTime * 9
      const harmonic = Math.abs(Math.sin(seed) * 0.6 + Math.cos(seed * 1.8) * 0.4)
      const amplitude = isPlaying ? Math.max(0.05, harmonic) : Math.max(0.05, Math.sin(i * 0.25) * 0.4 + 0.4)

      const barHeight = amplitude * (height * 0.8)
      const x = i * barSpacing + (barSpacing - barWidth) / 2
      const y = height - barHeight

      // Vertical Bar
      ctx.beginPath()
      ctx.roundRect(x, y, barWidth, Math.max(4, barHeight), [radius, radius, 0, 0])
      ctx.fill()

      // Floating Glowing Peak Cap Dot
      const capY = Math.max(2, y - 6)
      ctx.fillStyle = '#FFFFFF'
      ctx.beginPath()
      ctx.arc(x + barWidth / 2, capY, Math.max(1.5, barWidth / 4), 0, Math.PI * 2)
      ctx.fill()
      ctx.fillStyle = gradient
    }
  } else if (currentMode === 'circular') {
    // Mode 3: Radial Circular Audiogram with Inner Hole
    const centerX = width / 2
    const centerY = halfHeight
    const baseRadius = Math.min(centerX, centerY) * innerRadiusRatio.value
    const rayCount = Math.max(32, Math.min(96, density.value))

    ctx.strokeStyle = gradient
    ctx.lineWidth = Math.max(2, (width / rayCount) * 0.5)
    ctx.lineCap = 'round'

    // Inner Ring Avatar Circle
    ctx.beginPath()
    ctx.arc(centerX, centerY, baseRadius, 0, Math.PI * 2)
    ctx.strokeStyle = `${pColor}40`
    ctx.lineWidth = 2
    ctx.stroke()

    ctx.strokeStyle = gradient
    ctx.lineWidth = Math.max(2, (width / rayCount) * 0.5)

    for (let i = 0; i < rayCount; i++) {
      const angle = (i / rayCount) * Math.PI * 2
      const seed = i * 0.45 + currentTime * 7
      const amp = Math.abs(Math.sin(seed) * 0.65 + Math.cos(seed * 2.2) * 0.35)
      const rayLen = isPlaying ? amp * (baseRadius * 0.9) + 4 : 8

      const x1 = centerX + Math.cos(angle) * baseRadius
      const y1 = centerY + Math.sin(angle) * baseRadius
      const x2 = centerX + Math.cos(angle) * (baseRadius + rayLen)
      const y2 = centerY + Math.sin(angle) * (baseRadius + rayLen)

      ctx.beginPath()
      ctx.moveTo(x1, y1)
      ctx.lineTo(x2, y2)
      ctx.stroke()
    }
  } else if (currentMode === 'dots') {
    // Mode 4: Frequency Matrix Dots
    const colCount = Math.max(20, Math.min(80, density.value))
    const rowCount = 8
    const spacingX = width / colCount
    const spacingY = height / rowCount

    for (let col = 0; col < colCount; col++) {
      const seed = col * 0.35 + currentTime * 8
      const colAmp = isPlaying
        ? Math.abs(Math.sin(seed) * 0.6 + Math.cos(seed * 1.7) * 0.4)
        : Math.abs(Math.sin(col * 0.2) * 0.5 + 0.3)

      const activeRows = Math.round(colAmp * rowCount)

      for (let row = 0; row < rowCount; row++) {
        const fromBottom = rowCount - 1 - row
        const isActive = fromBottom <= activeRows

        const x = col * spacingX + spacingX / 2
        const y = row * spacingY + spacingY / 2
        const dotRadius = isActive ? Math.max(2, spacingX * 0.35) : Math.max(1, spacingX * 0.15)

        ctx.fillStyle = isActive ? gradient : `${pColor}20`
        ctx.beginPath()
        ctx.arc(x, y, dotRadius, 0, Math.PI * 2)
        ctx.fill()
      }
    }
  } else {
    // Mode 5: Smooth Neon Wave / Oscilloscope
    ctx.strokeStyle = gradient
    ctx.lineWidth = lineWidth.value * 2
    ctx.lineJoin = 'round'
    ctx.lineCap = 'round'

    const points = 120
    const step = width / points
    const wavePoints: Array<{ x: number; y: number }> = []

    for (let i = 0; i <= points; i++) {
      const x = i * step
      const seed = i * 0.18 + currentTime * 9
      const amp = isPlaying
        ? Math.sin(seed) * 0.6 + Math.cos(seed * 2.3) * 0.3
        : Math.sin(i * 0.15 + currentTime * 2) * 0.35

      const y = halfHeight + amp * (halfHeight * 0.75)
      wavePoints.push({ x, y })
    }

    ctx.beginPath()
    for (let i = 0; i < wavePoints.length; i++) {
      const pt = wavePoints[i]
      if (i === 0) {
        ctx.moveTo(pt.x, pt.y)
      } else {
        const prev = wavePoints[i - 1]
        const midX = (prev.x + pt.x) / 2
        const midY = (prev.y + pt.y) / 2
        ctx.quadraticCurveTo(prev.x, prev.y, midX, midY)
      }
    }
    ctx.stroke()

    // Optional Semi-Transparent Area Fill
    if (isFillArea.value) {
      const fillGradient = ctx.createLinearGradient(0, halfHeight, 0, height)
      fillGradient.addColorStop(0, `${pColor}40`)
      fillGradient.addColorStop(1, 'transparent')

      ctx.fillStyle = fillGradient
      ctx.lineTo(width, height)
      ctx.lineTo(0, height)
      ctx.closePath()
      ctx.fill()
    }
  }

  ctx.restore()
}

function startRenderLoop() {
  function loop() {
    drawWaveform()
    if (playbackStore.isPlaying) {
      animationFrameId = requestAnimationFrame(loop)
    }
  }
  stopRenderLoop()
  animationFrameId = requestAnimationFrame(loop)
}

function stopRenderLoop() {
  if (animationFrameId !== null) {
    cancelAnimationFrame(animationFrameId)
    animationFrameId = null
  }
}

watch(
  () => playbackStore.isPlaying,
  (playing) => {
    if (playing) {
      startRenderLoop()
    } else {
      stopRenderLoop()
      drawWaveform()
    }
  },
  { immediate: true }
)

watch(
  () => [
    playbackStore.currentTime,
    waveformConfig.value,
    targetDimensions.value,
    props.canvasWidth,
    props.canvasHeight,
  ],
  () => {
    drawWaveform()
  },
  { deep: true }
)

onMounted(() => {
  drawWaveform()
})

onUnmounted(() => {
  stopRenderLoop()
})
</script>
