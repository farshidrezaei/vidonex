<template>
  <div
    class="h-7 border-b border-gray-800 bg-gray-950 sticky top-0 z-40 flex-shrink-0 select-none cursor-pointer min-w-full relative"
    :style="{ width: `${totalWidth}px` }"
    @mousedown="handleRulerClick"
  >
    <!-- Active Content Boundary Region Highlight on Ruler -->
    <div
      v-if="timelineStore.contentDuration > 0"
      class="absolute top-0 bottom-0 left-0 bg-indigo-500/10 border-r-2 border-indigo-400/80 pointer-events-none z-10"
      :style="{ width: `${timelineStore.contentDuration * timelineStore.pixelsPerSecond}px` }"
      title="Active Video Content Duration"
    >
      <div class="absolute right-1 bottom-0.5 text-[8px] font-mono text-indigo-300 font-bold bg-indigo-950/90 px-1 rounded shadow-sm border border-indigo-500/30">
        END
      </div>
    </div>

    <!-- Major Time Marks -->
    <div
      v-for="mark in timeMarks"
      :key="mark.time"
      class="absolute top-0 bottom-0 border-l border-gray-700/80 flex flex-col justify-between pl-1 pointer-events-none"
      :style="{ left: `${mark.position}px` }"
    >
      <span class="text-[10px] font-mono text-gray-400 font-medium tracking-tight">
        {{ mark.label }}
      </span>
      <div class="w-px h-1.5 bg-gray-700"></div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useTimelineStore } from '~/stores/timeline'
import { usePlaybackStore } from '~/stores/playback'

const timelineStore = useTimelineStore()
const playbackStore = usePlaybackStore()

const totalWidth = computed(() => timelineStore.timelineDuration * timelineStore.pixelsPerSecond)

const timeMarks = computed(() => {
  const marks = []
  const duration = timelineStore.timelineDuration
  const pps = timelineStore.pixelsPerSecond

  // Dynamic interval: 1s, 2s, 5s, 10s depending on zoom
  let interval = 1
  if (pps < 30) interval = 5
  else if (pps < 60) interval = 2
  else interval = 1

  for (let t = 0; t <= duration; t += interval) {
    const m = Math.floor(t / 60)
    const s = Math.floor(t % 60)
    const label = `${String(m).padStart(2, '0')}:${String(s).padStart(2, '0')}`
    marks.push({
      time: t,
      position: t * pps,
      label,
    })
  }

  return marks
})

function handleRulerClick(event: MouseEvent) {
  const rect = (event.currentTarget as HTMLElement).getBoundingClientRect()
  const clickX = event.clientX - rect.left
  const targetTime = Math.max(0, clickX / timelineStore.pixelsPerSecond)
  playbackStore.seek(targetTime)

  const onMouseMove = (moveEvent: MouseEvent) => {
    const moveX = moveEvent.clientX - rect.left
    const newTime = Math.max(0, moveX / timelineStore.pixelsPerSecond)
    playbackStore.seek(newTime)
  }

  const onMouseUp = () => {
    window.removeEventListener('mousemove', onMouseMove)
    window.removeEventListener('mouseup', onMouseUp)
  }

  window.addEventListener('mousemove', onMouseMove)
  window.addEventListener('mouseup', onMouseUp)
}
</script>
