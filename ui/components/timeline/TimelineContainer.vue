<template>
  <div class="h-full flex flex-col bg-gray-950 select-none overflow-hidden" dir="ltr">
    <!-- Timeline Actions Toolbar -->
    <div class="h-10 border-b border-gray-800 px-3 bg-gray-950/80 flex items-center justify-between z-20">
      <!-- Left: Add Track, History, Edit Tools -->
      <div class="flex items-center gap-1.5">
        <!-- Add Track -->
        <UDropdownMenu :items="addTrackMenuItems">
          <UButton
            size="xs"
            color="neutral"
            variant="soft"
            icon="i-heroicons-plus"
          >
            {{ $t('timeline.add_track') }}
          </UButton>
        </UDropdownMenu>

        <div class="h-4 w-px bg-gray-800 mx-1"></div>

        <!-- Undo (Ctrl+Z) -->
        <UButton
          size="xs"
          color="neutral"
          variant="ghost"
          icon="i-heroicons-arrow-uturn-left"
          :disabled="timelineStore.undoStack.length === 0"
          title="Undo (Ctrl+Z)"
          @click="timelineStore.undo"
        />

        <!-- Redo (Ctrl+Y / Ctrl+Shift+Z) -->
        <UButton
          size="xs"
          color="neutral"
          variant="ghost"
          icon="i-heroicons-arrow-uturn-right"
          :disabled="timelineStore.redoStack.length === 0"
          title="Redo (Ctrl+Y / Ctrl+Shift+Z)"
          @click="timelineStore.redo"
        />

        <div class="h-4 w-px bg-gray-800 mx-1"></div>

        <!-- Split at Playhead (S) -->
        <UButton
          size="xs"
          color="neutral"
          variant="ghost"
          icon="i-heroicons-scissors"
          :disabled="!timelineStore.selectedClipId"
          :title="$t('timeline.split_at_playhead') + ' (S)'"
          @click="splitClip"
        >
          <span class="text-[11px] font-medium hidden sm:inline">{{ $t('timeline.split_at_playhead') }}</span>
        </UButton>

        <!-- Duplicate Clip (Ctrl+D) -->
        <UButton
          size="xs"
          color="neutral"
          variant="ghost"
          icon="i-heroicons-document-duplicate"
          :disabled="!timelineStore.selectedClipId"
          title="Duplicate Clip (Ctrl+D)"
          @click="duplicateClip"
        />

        <!-- Delete Selected (Del) -->
        <UButton
          size="xs"
          color="error"
          variant="ghost"
          icon="i-heroicons-trash"
          :disabled="!timelineStore.selectedClipId"
          :title="$t('timeline.delete_selected') + ' (Del)'"
          @click="deleteSelectedClip"
        />

        <div class="h-4 w-px bg-gray-800 mx-1"></div>

        <!-- Magnetic Snapping Toggle (N) -->
        <UButton
          size="xs"
          :color="timelineStore.isSnappingEnabled ? 'primary' : 'neutral'"
          :variant="timelineStore.isSnappingEnabled ? 'soft' : 'ghost'"
          icon="i-heroicons-bolt"
          :title="`Snapping / Magnet: ${timelineStore.isSnappingEnabled ? 'ON' : 'OFF'} (N)`"
          @click="timelineStore.isSnappingEnabled = !timelineStore.isSnappingEnabled"
        >
          <span class="text-[10px] font-mono font-semibold hidden md:inline">
            {{ timelineStore.isSnappingEnabled ? 'Snap ON' : 'Snap OFF' }}
          </span>
        </UButton>
      </div>

      <!-- Right: Zoom & Keyboard Help -->
      <div class="flex items-center gap-1">
        <!-- Keyboard Shortcuts Help (?) -->
        <UButton
          icon="i-heroicons-question-mark-circle"
          size="xs"
          color="neutral"
          variant="ghost"
          title="Keyboard Shortcuts (?)"
          @click="isShortcutsModalOpen = true"
        />

        <div class="h-4 w-px bg-gray-800 mx-1"></div>

        <!-- Zoom Controls -->
        <UButton
          icon="i-heroicons-minus"
          size="xs"
          color="neutral"
          variant="ghost"
          title="Zoom Out (-)"
          @click="timelineStore.pixelsPerSecond = Math.max(15, timelineStore.pixelsPerSecond - 15)"
        />
        <span class="font-mono text-[11px] text-gray-500 w-10 text-center">{{ timelineStore.pixelsPerSecond }}px</span>
        <UButton
          icon="i-heroicons-plus"
          size="xs"
          color="neutral"
          variant="ghost"
          title="Zoom In (+)"
          @click="timelineStore.pixelsPerSecond = Math.min(200, timelineStore.pixelsPerSecond + 15)"
        />
      </div>
    </div>

    <!-- Timeline Workspace (Split Layout) -->
    <div class="flex-1 flex overflow-hidden min-h-0">
      <!-- Left Column: Track Headers -->
      <div class="w-56 flex-shrink-0 border-r border-gray-800 flex flex-col bg-gray-950 z-10 min-h-0">
        <!-- Ruler spacer -->
        <div class="h-7 border-b border-gray-800 bg-gray-950 px-3 flex items-center flex-shrink-0">
          <span class="text-[10px] font-mono text-gray-500 uppercase tracking-wider">{{ $t('timeline.tracks') }}</span>
        </div>

        <!-- Track Headers List -->
        <div ref="trackHeadersScrollRef" class="flex-1 overflow-y-hidden min-h-0">
          <TimelineTrackHeader
            v-for="track in timelineStore.tracks"
            :key="track.id"
            :track="track"
          />
        </div>
      </div>

      <!-- Right Column: Scrollable Tracks & Ruler Area -->
      <div
        ref="timelineScrollRef"
        class="flex-1 flex flex-col overflow-x-auto overflow-y-auto relative timeline-grid min-h-0"
        @scroll="handleTimelineScroll"
      >
        <!-- Time Ruler -->
        <TimelineRuler />

        <!-- Track Rows Lane Container -->
        <div
          class="relative min-w-full flex-1"
          :style="{ width: `${timelineStore.timelineDuration * timelineStore.pixelsPerSecond}px` }"
        >
          <!-- Playhead Needle -->
          <TimelinePlayhead />

          <!-- Track Rows -->
          <TimelineTrackRow
            v-for="track in timelineStore.tracks"
            :key="track.id"
            :track="track"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useTimelineStore } from '~/stores/timeline'
import { usePlaybackStore } from '~/stores/playback'
import { isShortcutsModalOpen } from '~/composables/useGlobalShortcuts'

const { t } = useI18n()
const timelineStore = useTimelineStore()
const playbackStore = usePlaybackStore()

const trackHeadersScrollRef = ref<HTMLDivElement | null>(null)
const timelineScrollRef = ref<HTMLDivElement | null>(null)

const addTrackMenuItems = computed(() => [
  [
    {
      label: t('timeline.track_types.video'),
      icon: 'i-heroicons-film',
      onSelect: () => timelineStore.addTrack('video'),
    },
    {
      label: t('timeline.track_types.audio'),
      icon: 'i-heroicons-speaker-wave',
      onSelect: () => timelineStore.addTrack('audio'),
    },
    {
      label: t('timeline.track_types.overlay'),
      icon: 'i-heroicons-square-2-stack',
      onSelect: () => timelineStore.addTrack('overlay'),
    },
    {
      label: t('timeline.track_types.subtitle'),
      icon: 'i-heroicons-document-text',
      onSelect: () => timelineStore.addTrack('subtitle'),
    },
    {
      label: t('timeline.track_types.waveform'),
      icon: 'i-heroicons-chart-bar',
      onSelect: () => timelineStore.addTrack('waveform'),
    },
  ],
])

function splitClip() {
  timelineStore.splitClipAtPlayhead(playbackStore.currentTime)
}

function duplicateClip() {
  if (timelineStore.selectedClipId) {
    timelineStore.duplicateClip(timelineStore.selectedClipId)
  }
}

function deleteSelectedClip() {
  if (timelineStore.selectedClipId) {
    timelineStore.removeClip(timelineStore.selectedClipId)
  }
}

function handleTimelineScroll(event: Event) {
  const target = event.target as HTMLElement
  if (trackHeadersScrollRef.value) {
    trackHeadersScrollRef.value.scrollTop = target.scrollTop
  }
}
</script>
