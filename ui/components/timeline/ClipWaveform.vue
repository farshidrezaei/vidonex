<template>
  <div class="w-full h-full flex items-center px-1 overflow-hidden pointer-events-none select-none">
    <canvas
      ref="canvasRef"
      class="w-full h-full opacity-60"
      :class="trackKind === 'audio' ? 'text-emerald-400' : 'text-purple-400'"
    ></canvas>
  </div>
</template>

<script setup lang="ts">
import type { TrackKind } from '~/types/spec'

const props = defineProps<{
  clipId: string
  source?: string
  duration: number
  inTrim?: number
  trackKind: TrackKind
  width: number
}>()

const canvasRef = ref<HTMLCanvasElement | null>(null)

// Deterministic Pseudo-Random Generator based on string seed
function createSeededRandom(seedStr: string) {
  let h = 0xdeadbeef
  for (let i = 0; i < seedStr.length; i++) {
    h = Math.imul(h ^ seedStr.charCodeAt(i), 2654435761)
  }
  return function () {
    h = Math.imul(h ^ (h >>> 16), 2246822507)
    h = Math.imul(h ^ (h >>> 13), 3266489909)
    return ((h >>> 0) / 4294967296)
  }
}

function drawWaveform() {
  const canvas = canvasRef.value
  if (!canvas) return

  const ctx = canvas.getContext('2d')
  if (!ctx) return

  const rect = canvas.getBoundingClientRect()
  const displayWidth = Math.max(10, Math.floor(rect.width || props.width || 100))
  const displayHeight = Math.max(10, Math.floor(rect.height || 40))

  const dpr = window.devicePixelRatio || 1
  canvas.width = displayWidth * dpr
  canvas.height = displayHeight * dpr
  ctx.scale(dpr, dpr)

  ctx.clearRect(0, 0, displayWidth, displayHeight)

  const barWidth = 3
  const barGap = 1.5
  const step = barWidth + barGap
  const numBars = Math.floor(displayWidth / step)

  const seed = `${props.source || ''}-${props.clipId}`
  const rand = createSeededRandom(seed)

  // Pre-generate base peaks
  const midY = displayHeight / 2
  const maxAmp = (displayHeight / 2) * 0.85

  const primaryColor = props.trackKind === 'audio' ? '#34d399' : '#c084fc'
  const secondaryColor = props.trackKind === 'audio' ? '#059669' : '#7e22ce'

  // Center baseline
  ctx.strokeStyle = 'rgba(255, 255, 255, 0.15)'
  ctx.lineWidth = 1
  ctx.beginPath()
  ctx.moveTo(0, midY)
  ctx.lineTo(displayWidth, midY)
  ctx.stroke()

  // Draw vertical audio bars
  for (let i = 0; i < numBars; i++) {
    const x = i * step
    // Normalized position along the clip
    const t = (i / numBars) * (props.duration || 10) + (props.inTrim || 0)

    // Layered harmonic envelope for realistic music/speech rhythm
    const r1 = rand()
    const r2 = rand()
    const modulation =
      Math.sin(t * 1.5) * 0.3 +
      Math.sin(t * 4.2) * 0.25 +
      Math.cos(t * 8.7) * 0.15 +
      0.5

    const rawAmp = (r1 * 0.5 + r2 * 0.5) * modulation
    const barHeight = Math.max(2, rawAmp * maxAmp)

    // Gradient fill for depth
    const gradient = ctx.createLinearGradient(0, midY - barHeight, 0, midY + barHeight)
    gradient.addColorStop(0, primaryColor)
    gradient.addColorStop(0.5, secondaryColor)
    gradient.addColorStop(1, primaryColor)

    ctx.fillStyle = gradient

    // Draw symmetrical rounded bar
    ctx.beginPath()
    const radius = 1.2
    const topY = midY - barHeight
    const height = barHeight * 2

    ctx.roundRect(x, topY, barWidth, height, radius)
    ctx.fill()
  }
}

watch(
  () => [props.width, props.duration, props.inTrim, props.source],
  () => {
    nextTick(drawWaveform)
  },
  { immediate: true }
)

onMounted(() => {
  nextTick(drawWaveform)
})
</script>
