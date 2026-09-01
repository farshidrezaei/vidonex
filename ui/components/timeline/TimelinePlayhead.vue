<template>
  <div
    class="absolute top-0 bottom-0 z-30 pointer-events-none transition-none"
    :style="{ left: `${positionX}px` }"
  >
    <!-- Scrubber Head -->
    <div
      class="w-3.5 h-4 bg-indigo-500 rounded-b -translate-x-1/2 flex items-center justify-center cursor-ew-resize pointer-events-auto shadow-md shadow-indigo-500/40"
      @mousedown="startScrubbing"
    >
      <div class="w-1 h-2 bg-white/80 rounded-full"></div>
    </div>

    <!-- Playhead Needle Line -->
    <div class="w-0.5 h-full bg-indigo-500 -translate-x-1/2 shadow-[0_0_8px_rgba(99,102,241,0.6)]"></div>
  </div>
</template>

<script setup lang="ts">
import { usePlaybackStore } from '~/stores/playback'
import { useTimelineStore } from '~/stores/timeline'
import { useSnapping } from '~/composables/useSnapping'

const playbackStore = usePlaybackStore()
const timelineStore = useTimelineStore()
const { snapTime } = useSnapping()

const positionX = computed(() => playbackStore.currentTime * timelineStore.pixelsPerSecond)

function startScrubbing(event: MouseEvent) {
  event.stopPropagation()
  const startX = event.clientX
  const initialTime = playbackStore.currentTime

  const onMouseMove = (moveEvent: MouseEvent) => {
    const deltaX = moveEvent.clientX - startX
    const deltaTime = deltaX / timelineStore.pixelsPerSecond
    const rawTime = Math.max(0, initialTime + deltaTime)
    const { time: snappedTime } = snapTime(rawTime)
    playbackStore.seek(snappedTime)
  }

  const onMouseUp = () => {
    window.removeEventListener('mousemove', onMouseMove)
    window.removeEventListener('mouseup', onMouseUp)
  }

  window.addEventListener('mousemove', onMouseMove)
  window.addEventListener('mouseup', onMouseUp)
}
</script>
