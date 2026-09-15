<template>
  <div class="space-y-2 select-none">
    <div class="flex items-center justify-between text-[11px] text-gray-400">
      <span>{{ $t('inspector.motion.curve_preview') || 'Easing Curve Graph' }}</span>
      <span class="font-mono text-indigo-400 text-[10px]">{{ easing }}</span>
    </div>

    <!-- Interactive Canvas Graph Container -->
    <div class="relative w-full h-28 bg-gray-950 rounded-lg border border-gray-800/80 p-2 overflow-hidden group">
      <canvas
        ref="canvasRef"
        class="w-full h-full block"
      ></canvas>

      <!-- Overlay Play Button to Preview Easing in Motion -->
      <button
        type="button"
        class="absolute bottom-2 right-2 p-1.5 rounded bg-gray-900/90 hover:bg-indigo-600 border border-gray-700/80 text-white transition-all shadow-md"
        :title="$t('inspector.motion.test_curve') || 'Preview Curve Motion'"
        @click="triggerAnimationPreview"
      >
        <UIcon name="i-heroicons-play" class="w-3.5 h-3.5" />
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{
  easing: string
}>()

const canvasRef = ref<HTMLCanvasElement | null>(null)
let previewProgress = 0
let isPreviewing = false
let animFrameId: number | null = null

// Mathematical Easing Function Evaluator for t in [0, 1]
function evaluateEasing(t: number, easingName: string): number {
  const clampT = Math.max(0, Math.min(1, t))
  switch (easingName) {
    case 'linear':
      return clampT
    case 'easeInQuad':
      return clampT * clampT
    case 'easeOutQuad':
      return 1 - (1 - clampT) * (1 - clampT)
    case 'easeInOutQuad':
      return clampT < 0.5 ? 2 * clampT * clampT : 1 - Math.pow(-2 * clampT + 2, 2) / 2
    case 'easeInOutCubic':
      return clampT < 0.5 ? 4 * clampT * clampT * clampT : 1 - Math.pow(-2 * clampT + 2, 3) / 2
    default:
      return clampT < 0.5 ? 2 * clampT * clampT : 1 - Math.pow(-2 * clampT + 2, 2) / 2
  }
}

function drawCurve() {
  const canvas = canvasRef.value
  if (!canvas) return

  const ctx = canvas.getContext('2d')
  if (!ctx) return

  const rect = canvas.getBoundingClientRect()
  const width = Math.max(10, Math.floor(rect.width || 200))
  const height = Math.max(10, Math.floor(rect.height || 96))

  const dpr = window.devicePixelRatio || 1
  canvas.width = width * dpr
  canvas.height = height * dpr
  ctx.scale(dpr, dpr)

  ctx.clearRect(0, 0, width, height)

  const padLeft = 14
  const padRight = 14
  const padTop = 12
  const padBottom = 12

  const graphWidth = width - padLeft - padRight
  const graphHeight = height - padTop - padBottom

  // Subtle Background Grid
  ctx.strokeStyle = 'rgba(255, 255, 255, 0.05)'
  ctx.lineWidth = 1
  for (let i = 0; i <= 4; i++) {
    const gx = padLeft + (graphWidth / 4) * i
    const gy = padTop + (graphHeight / 4) * i
    ctx.beginPath()
    ctx.moveTo(gx, padTop)
    ctx.lineTo(gx, padTop + graphHeight)
    ctx.stroke()

    ctx.beginPath()
    ctx.moveTo(padLeft, gy)
    ctx.lineTo(padLeft + graphWidth, gy)
    ctx.stroke()
  }

  // Draw Gradient Area Under Curve
  const samples = 60
  const areaGradient = ctx.createLinearGradient(0, padTop, 0, padTop + graphHeight)
  areaGradient.addColorStop(0, 'rgba(99, 102, 241, 0.25)')
  areaGradient.addColorStop(1, 'rgba(99, 102, 241, 0.0)')

  ctx.fillStyle = areaGradient
  ctx.beginPath()
  ctx.moveTo(padLeft, padTop + graphHeight)

  for (let i = 0; i <= samples; i++) {
    const t = i / samples
    const val = evaluateEasing(t, props.easing)
    const x = padLeft + t * graphWidth
    const y = padTop + graphHeight - val * graphHeight
    ctx.lineTo(x, y)
  }
  ctx.lineTo(padLeft + graphWidth, padTop + graphHeight)
  ctx.closePath()
  ctx.fill()

  // Draw Smooth Curve Line
  ctx.strokeStyle = '#818cf8'
  ctx.lineWidth = 2.5
  ctx.beginPath()
  for (let i = 0; i <= samples; i++) {
    const t = i / samples
    const val = evaluateEasing(t, props.easing)
    const x = padLeft + t * graphWidth
    const y = padTop + graphHeight - val * graphHeight
    if (i === 0) {
      ctx.moveTo(x, y)
    } else {
      ctx.lineTo(x, y)
    }
  }
  ctx.stroke()

  // Draw Endpoint Handles
  ctx.fillStyle = '#6366f1'
  ctx.strokeStyle = '#ffffff'
  ctx.lineWidth = 1.5

  // Start point (0, 0)
  ctx.beginPath()
  ctx.arc(padLeft, padTop + graphHeight, 3.5, 0, Math.PI * 2)
  ctx.fill()
  ctx.stroke()

  // End point (1, 1)
  ctx.beginPath()
  ctx.arc(padLeft + graphWidth, padTop, 3.5, 0, Math.PI * 2)
  ctx.fill()
  ctx.stroke()

  // Draw Animated Preview Tracer Dot
  if (isPreviewing) {
    const easeVal = evaluateEasing(previewProgress, props.easing)
    const dotX = padLeft + previewProgress * graphWidth
    const dotY = padTop + graphHeight - easeVal * graphHeight

    ctx.fillStyle = '#38bdf8'
    ctx.shadowColor = '#38bdf8'
    ctx.shadowBlur = 8
    ctx.beginPath()
    ctx.arc(dotX, dotY, 4.5, 0, Math.PI * 2)
    ctx.fill()
    ctx.shadowBlur = 0
  }
}

function triggerAnimationPreview() {
  if (isPreviewing) return
  isPreviewing = true
  previewProgress = 0
  const startTime = performance.now()
  const duration = 1200 // 1.2s animation

  function step(now: number) {
    const elapsed = now - startTime
    previewProgress = Math.min(1, elapsed / duration)
    drawCurve()

    if (previewProgress < 1) {
      animFrameId = requestAnimationFrame(step)
    } else {
      setTimeout(() => {
        isPreviewing = false
        drawCurve()
      }, 200)
    }
  }

  animFrameId = requestAnimationFrame(step)
}

watch(
  () => props.easing,
  () => {
    nextTick(drawCurve)
  },
  { immediate: true }
)

onMounted(() => {
  nextTick(drawCurve)
})

onUnmounted(() => {
  if (animFrameId !== null) {
    cancelAnimationFrame(animFrameId)
  }
})
</script>
