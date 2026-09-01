<template>
  <div
    class="h-14 border-b border-gray-800/80 bg-gray-950 px-3 flex items-center justify-between select-none cursor-pointer transition"
    :class="isSelected ? 'border-l-2 border-l-indigo-500 bg-gray-900/60' : 'hover:bg-gray-900/30'"
    @click="timelineStore.selectedTrackId = track.id"
  >
    <!-- Left: Track Kind & ID -->
    <div class="flex items-center gap-2 truncate">
      <div class="w-6 h-6 rounded flex items-center justify-center text-xs" :class="kindBadgeClass">
        <UIcon :name="kindIcon" class="w-3.5 h-3.5" />
      </div>
      <span class="text-xs font-medium text-gray-300 truncate">{{ track.id }}</span>
    </div>

    <!-- Right: Track Actions (Mute, Delete) -->
    <div class="flex items-center gap-1">
      <!-- Mute Button -->
      <button
        class="p-1 rounded text-gray-500 hover:text-gray-300 transition"
        :class="track.muted ? 'text-red-400' : ''"
        :title="track.muted ? $t('timeline.unmute_track') : $t('timeline.mute_track')"
        @click.stop="toggleMute"
      >
        <UIcon :name="track.muted ? 'i-heroicons-speaker-x-mark' : 'i-heroicons-speaker-wave'" class="w-3.5 h-3.5" />
      </button>

      <!-- Delete Track -->
      <button
        class="p-1 rounded text-gray-600 hover:text-red-400 transition"
        :title="$t('timeline.tracks')"
        @click.stop="timelineStore.removeTrack(track.id)"
      >
        <UIcon name="i-heroicons-trash" class="w-3.5 h-3.5" />
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useTimelineStore } from '~/stores/timeline'
import type { TrackSpec } from '~/types/spec'

const props = defineProps<{
  track: TrackSpec
}>()

const timelineStore = useTimelineStore()
const isSelected = computed(() => timelineStore.selectedTrackId === props.track.id)

const kindIcon = computed(() => {
  switch (props.track.kind) {
    case 'audio':
      return 'i-heroicons-speaker-wave'
    case 'waveform':
      return 'i-heroicons-chart-bar'
    case 'subtitle':
      return 'i-heroicons-document-text'
    case 'overlay':
      return 'i-heroicons-square-2-stack'
    default:
      return 'i-heroicons-film'
  }
})

const kindBadgeClass = computed(() => {
  switch (props.track.kind) {
    case 'audio':
      return 'bg-emerald-500/20 text-emerald-400'
    case 'waveform':
      return 'bg-purple-500/20 text-purple-400'
    case 'subtitle':
      return 'bg-amber-500/20 text-amber-400'
    default:
      return 'bg-indigo-500/20 text-indigo-400'
  }
})

function toggleMute() {
  props.track.muted = !props.track.muted
  timelineStore.pushHistoryState('Toggle Track Mute')
}
</script>
