<template>
  <footer class="h-9 border-t border-gray-800 bg-gray-950 px-4 flex items-center justify-between text-xs text-gray-400 select-none z-30">
    <!-- Left: Timecode & Duration -->
    <div class="flex items-center gap-4 font-mono">
      <div class="flex items-center gap-1.5 text-gray-200">
        <span class="text-indigo-400 font-semibold">{{ formattedTime }}</span>
        <span class="text-gray-600">/</span>
        <span class="text-gray-400">{{ formattedDuration }}</span>
      </div>

      <div class="h-3.5 w-px bg-gray-800"></div>

      <div class="flex items-center gap-1 text-gray-400">
        <span class="text-gray-500">FPS:</span>
        <span class="text-gray-300">{{ projectStore.frameRate }}</span>
      </div>
    </div>

    <!-- Right: Snapping & Zoom Controls -->
    <div class="flex items-center gap-4">
      <!-- Magnetic Snapping -->
      <button
        class="flex items-center gap-1.5 px-2 py-0.5 rounded text-xs transition"
        :class="timelineStore.snappingEnabled ? 'bg-indigo-500/20 text-indigo-400 border border-indigo-500/30' : 'text-gray-500 hover:text-gray-300'"
        @click="timelineStore.snappingEnabled = !timelineStore.snappingEnabled"
      >
        <UIcon name="i-heroicons-bolt" class="w-3.5 h-3.5" />
        <span>{{ $t('timeline.snap_toggle') }}</span>
      </button>

      <div class="h-3.5 w-px bg-gray-800"></div>

      <!-- Zoom Slider -->
      <div class="flex items-center gap-2">
        <button
          class="hover:text-gray-200 transition"
          :title="$t('timeline.zoom_out')"
          @click="zoomOut"
        >
          <UIcon name="i-heroicons-minus" class="w-3.5 h-3.5" />
        </button>

        <input
          v-model.number="timelineStore.pixelsPerSecond"
          type="range"
          min="15"
          max="200"
          step="5"
          class="w-24 h-1 bg-gray-800 rounded-lg appearance-none cursor-pointer accent-indigo-500"
        />

        <button
          class="hover:text-gray-200 transition"
          :title="$t('timeline.zoom_in')"
          @click="zoomIn"
        >
          <UIcon name="i-heroicons-plus" class="w-3.5 h-3.5" />
        </button>

        <span class="font-mono text-gray-500 w-10 text-right">{{ timelineStore.pixelsPerSecond }}px/s</span>
      </div>
    </div>
  </footer>
</template>

<script setup lang="ts">
import { usePlaybackStore } from '~/stores/playback'
import { useProjectStore } from '~/stores/project'
import { useTimelineStore } from '~/stores/timeline'

const playbackStore = usePlaybackStore()
const projectStore = useProjectStore()
const timelineStore = useTimelineStore()

const formattedTime = computed(() => playbackStore.formatTimecode(playbackStore.currentTime, projectStore.frameRate))
const formattedDuration = computed(() => playbackStore.formatTimecode(playbackStore.duration, projectStore.frameRate))

function zoomIn() {
  timelineStore.pixelsPerSecond = Math.min(200, timelineStore.pixelsPerSecond + 15)
}

function zoomOut() {
  timelineStore.pixelsPerSecond = Math.max(15, timelineStore.pixelsPerSecond - 15)
}
</script>
