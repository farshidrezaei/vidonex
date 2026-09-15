<template>
  <div class="flex items-center gap-1.5 px-2 py-0.5 bg-gray-900/90 rounded border border-gray-800 text-[10px] font-mono select-none">
    <!-- Channel Labels & Dual Horizontal Meters -->
    <div class="flex flex-col gap-0.5 justify-center">
      <!-- Left Channel -->
      <div class="flex items-center gap-1">
        <span class="text-[8px] text-gray-500 font-bold w-2">L</span>
        <div class="relative w-16 h-1.5 bg-gray-950 rounded-sm overflow-hidden flex items-center">
          <!-- Meter Fill Gradient -->
          <div
            class="h-full transition-all duration-75 ease-out rounded-sm"
            :style="{
              width: `${leftPercent}%`,
              background: getMeterGradient(leftDb)
            }"
          ></div>
          <!-- Peak Hold Tick -->
          <div
            v-if="leftPeakPercent > 0"
            class="absolute top-0 bottom-0 w-0.5 bg-white shadow-sm transition-all duration-300"
            :style="{ left: `${leftPeakPercent}%` }"
          ></div>
        </div>
      </div>

      <!-- Right Channel -->
      <div class="flex items-center gap-1">
        <span class="text-[8px] text-gray-500 font-bold w-2">R</span>
        <div class="relative w-16 h-1.5 bg-gray-950 rounded-sm overflow-hidden flex items-center">
          <!-- Meter Fill Gradient -->
          <div
            class="h-full transition-all duration-75 ease-out rounded-sm"
            :style="{
              width: `${rightPercent}%`,
              background: getMeterGradient(rightDb)
            }"
          ></div>
          <!-- Peak Hold Tick -->
          <div
            v-if="rightPeakPercent > 0"
            class="absolute top-0 bottom-0 w-0.5 bg-white shadow-sm transition-all duration-300"
            :style="{ left: `${rightPeakPercent}%` }"
          ></div>
        </div>
      </div>
    </div>

    <!-- dBFS Numeric Indicator -->
    <span
      class="w-11 text-right text-[9px] font-semibold tabular-nums"
      :class="currentPeakDb > -3 ? 'text-red-400' : currentPeakDb > -12 ? 'text-amber-400' : 'text-emerald-400'"
    >
      {{ displayDbText }}
    </span>
  </div>
</template>

<script setup lang="ts">
import { usePlaybackStore } from '~/stores/playback'
import { useTimelineStore } from '~/stores/timeline'

const playbackStore = usePlaybackStore()
const timelineStore = useTimelineStore()

const leftDb = ref(-60)
const rightDb = ref(-60)
const leftPeakDb = ref(-60)
const rightPeakDb = ref(-60)

let peakDecayTimer: number | null = null
let animationFrameId: number | null = null

// Converts dBFS (-60 to 0) to percentage (0 to 100)
function dbToPercent(db: number): number {
  if (db <= -60) return 0
  if (db >= 0) return 100
  return Math.min(100, Math.max(0, ((db + 60) / 60) * 100))
}

const leftPercent = computed(() => dbToPercent(leftDb.value))
const rightPercent = computed(() => dbToPercent(rightDb.value))
const leftPeakPercent = computed(() => dbToPercent(leftPeakDb.value))
const rightPeakPercent = computed(() => dbToPercent(rightPeakDb.value))

const currentPeakDb = computed(() => Math.max(leftDb.value, rightDb.value))

const displayDbText = computed(() => {
  if (currentPeakDb.value <= -59) return '-∞ dB'
  return `${currentPeakDb.value.toFixed(1)} dB`
})

function getMeterGradient(db: number): string {
  if (db > -3) {
    return 'linear-gradient(90deg, #10b981 0%, #f59e0b 70%, #ef4444 100%)'
  }
  if (db > -12) {
    return 'linear-gradient(90deg, #10b981 0%, #f59e0b 100%)'
  }
  return 'linear-gradient(90deg, #059669 0%, #10b981 100%)'
}

// Check if any audio track has an active clip at current playhead time
function hasActiveAudioAtPlayhead(): boolean {
  const currentTime = playbackStore.currentTime
  for (const track of timelineStore.tracks) {
    if (track.kind === 'audio' || track.kind === 'waveform' || track.kind === 'video') {
      for (const clip of track.clips || []) {
        const start = clip.start || 0
        const end = start + (clip.duration || 0)
        if (currentTime >= start && currentTime <= end) {
          return true
        }
      }
    }
  }
  return false
}

// Simulated high-fidelity Audio Peak Meter tick
function updateMeter() {
  if (playbackStore.isPlaying && !playbackStore.isMuted && playbackStore.volume > 0 && hasActiveAudioAtPlayhead()) {
    const vol = playbackStore.volume
    // Dynamic noise envelope for organic movement
    const t = performance.now() / 100
    const rawLeft = Math.sin(t * 1.7) * 0.3 + Math.cos(t * 3.1) * 0.2 + 0.5
    const rawRight = Math.cos(t * 1.9) * 0.3 + Math.sin(t * 2.8) * 0.2 + 0.5

    // Map to realistic master dB (-36 to -2 dB scaled by master volume)
    const targetLeft = -36 + rawLeft * 32 + (vol - 1) * 20
    const targetRight = -36 + rawRight * 32 + (vol - 1) * 20

    leftDb.value = Math.max(-60, Math.min(0, targetLeft))
    rightDb.value = Math.max(-60, Math.min(0, targetRight))

    // Update peak hold
    if (leftDb.value > leftPeakDb.value) leftPeakDb.value = leftDb.value
    if (rightDb.value > rightPeakDb.value) rightPeakDb.value = rightDb.value
  } else {
    // Smooth decay to zero
    leftDb.value = Math.max(-60, leftDb.value - 2.5)
    rightDb.value = Math.max(-60, rightDb.value - 2.5)
    leftPeakDb.value = Math.max(-60, leftPeakDb.value - 1.2)
    rightPeakDb.value = Math.max(-60, rightPeakDb.value - 1.2)
  }

  animationFrameId = requestAnimationFrame(updateMeter)
}

onMounted(() => {
  animationFrameId = requestAnimationFrame(updateMeter)
})

onUnmounted(() => {
  if (animationFrameId !== null) {
    cancelAnimationFrame(animationFrameId)
  }
  if (peakDecayTimer) {
    clearInterval(peakDecayTimer)
  }
})
</script>
