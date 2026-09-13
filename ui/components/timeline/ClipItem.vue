<template>
  <div
    :draggable="!isTrimming"
    class="absolute top-1.5 bottom-1.5 rounded-lg overflow-hidden select-none group transition-[border-color,box-shadow,opacity] duration-150 cursor-grab active:cursor-grabbing"
    :class="[
      isSelected ? 'border-2 border-indigo-400 shadow-lg shadow-indigo-500/30 z-20 ring-1 ring-indigo-400/40' : 'border border-gray-700/80 hover:border-gray-500 z-10',
      isBeingDragged ? 'opacity-30 scale-95 border-dashed border-indigo-300' : '',
      isTrimming ? 'cursor-ew-resize' : '',
      trackKindColorClass
    ]"
    :style="clipStyle"
    @click.stop="handleSelect"
    @dragstart="handleDragStart"
    @dragend="handleDragEnd"
  >
    <UContextMenu :items="contextMenuItems">
      <div class="w-full h-full relative">
        <!-- Background Filmstrip Thumbnail for Video / Image Clips -->
        <div
          v-if="thumbnailUrl && (trackKind === 'video' || trackKind === 'overlay')"
          class="absolute inset-0 opacity-25 pointer-events-none overflow-hidden flex items-center"
        >
          <img
            :src="thumbnailUrl"
            class="w-full h-full object-cover filter brightness-75 contrast-125"
            alt=""
          />
          <div class="absolute inset-0 bg-gradient-to-r from-gray-950/85 via-transparent to-gray-950/85"></div>
        </div>

        <!-- Audio Waveform Visualization SVG for Audio / Waveform Clips -->
        <div
          v-else-if="trackKind === 'audio' || trackKind === 'waveform'"
          class="absolute inset-0 pointer-events-none flex items-center px-1 overflow-hidden"
        >
          <svg
            class="w-full h-6"
            :class="trackKind === 'audio' ? 'text-emerald-400/50' : 'text-purple-400/50'"
            preserveAspectRatio="none"
          >
            <defs>
              <pattern
                :id="`waveform-pattern-${clip.id}`"
                width="48"
                height="24"
                patternUnits="userSpaceOnUse"
              >
                <!-- Subtle Center Baseline -->
                <line x1="0" y1="12" x2="48" y2="12" stroke="currentColor" stroke-width="0.75" opacity="0.25" />

                <!-- Symmetrical Bounded Waveform Bars -->
                <rect x="1" y="7" width="2" height="10" rx="1" fill="currentColor" />
                <rect x="5" y="4" width="2" height="16" rx="1" fill="currentColor" />
                <rect x="9" y="8" width="2" height="8" rx="1" fill="currentColor" />
                <rect x="13" y="2" width="2" height="20" rx="1" fill="currentColor" />
                <rect x="17" y="5" width="2" height="14" rx="1" fill="currentColor" />
                <rect x="21" y="9" width="2" height="6" rx="1" fill="currentColor" />
                <rect x="25" y="3" width="2" height="18" rx="1" fill="currentColor" />
                <rect x="29" y="6" width="2" height="12" rx="1" fill="currentColor" />
                <rect x="33" y="1" width="2" height="22" rx="1" fill="currentColor" />
                <rect x="37" y="7" width="2" height="10" rx="1" fill="currentColor" />
                <rect x="41" y="4" width="2" height="16" rx="1" fill="currentColor" />
                <rect x="45" y="8" width="2" height="8" rx="1" fill="currentColor" />
              </pattern>
            </defs>
            <rect width="100%" height="100%" :fill="`url(#waveform-pattern-${clip.id})`" />
          </svg>
        </div>

        <!-- Left Trim / Resize Handle (100% Isolated from Drag) -->
        <div
          class="absolute left-0 top-0 bottom-0 w-3 bg-white/10 hover:bg-indigo-400 hover:w-3.5 cursor-ew-resize opacity-0 group-hover:opacity-100 transition z-30 flex items-center justify-center pointer-events-auto"
          @mousedown.stop="startLeftTrim"
          @dragstart.stop.prevent
          @click.stop
        >
          <div class="w-0.5 h-4 bg-white/90 rounded-full"></div>
        </div>

        <!-- Clip Content Body (Fills Track Height) -->
        <div class="w-full h-full px-2.5 flex items-center justify-between text-xs text-white/90 font-medium pointer-events-none relative z-10">
          <div class="flex items-center gap-2 truncate">
            <!-- Leading Thumbnail or Icon -->
            <div
              v-if="thumbnailUrl && (trackKind === 'video' || trackKind === 'overlay')"
              class="w-8 h-8 rounded overflow-hidden flex-shrink-0 border border-white/20 shadow bg-gray-950 flex items-center justify-center"
            >
              <img
                :src="thumbnailUrl"
                class="w-full h-full object-cover"
                alt=""
              />
            </div>
            <div
              v-else
              class="w-7 h-7 rounded bg-black/40 flex items-center justify-center flex-shrink-0 border border-white/10"
            >
              <UIcon :name="clipIcon" class="w-4 h-4 opacity-80" />
            </div>

            <div class="flex flex-col truncate justify-center">
              <span class="truncate text-[11px] font-medium text-white drop-shadow">{{ clipTitle }}</span>
              <span v-if="clip.trim && clip.trim > 0" class="text-[9px] font-mono text-gray-400 leading-tight">
                trim +{{ clip.trim.toFixed(1) }}s
              </span>
            </div>
          </div>

          <div class="flex items-center gap-1 flex-shrink-0">
            <span class="text-[10px] font-mono text-white/90 bg-black/60 px-1.5 py-0.5 rounded border border-white/10 shadow-sm">
              {{ formattedDuration }}
            </span>
          </div>
        </div>

        <!-- Right Trim / Resize Handle (100% Isolated from Drag) -->
        <div
          class="absolute right-0 top-0 bottom-0 w-3 bg-white/10 hover:bg-indigo-400 hover:w-3.5 cursor-ew-resize opacity-0 group-hover:opacity-100 transition z-30 flex items-center justify-center pointer-events-auto"
          @mousedown.stop="startRightTrim"
          @dragstart.stop.prevent
          @click.stop
        >
          <div class="w-0.5 h-4 bg-white/90 rounded-full"></div>
        </div>
      </div>
    </UContextMenu>
  </div>
</template>

<script setup lang="ts">
import { useTimelineStore } from '~/stores/timeline'
import { usePlaybackStore } from '~/stores/playback'
import { useMediaStore } from '~/stores/media'
import { useSnapping } from '~/composables/useSnapping'
import type { ClipSpec, TrackKind } from '~/types/spec'

const props = defineProps<{
  clip: ClipSpec
  trackId: string
  trackKind: TrackKind
}>()

const timelineStore = useTimelineStore()
const playbackStore = usePlaybackStore()
const mediaStore = useMediaStore()
const { snapTime } = useSnapping()

const isTrimming = ref(false)

const isSelected = computed(() => timelineStore.selectedClipId === props.clip.id)
const isBeingDragged = computed(() => timelineStore.draggedTimelineClip?.clip.id === props.clip.id)

const matchedAsset = computed(() =>
  mediaStore.assets.find(
    (a) => a.file_path === props.clip.source || a.file_name === props.clip.source || a.id === props.clip.source
  )
)

const isStaticImage = computed(() => {
  if (matchedAsset.value?.file_type === 'image') return true
  const src = props.clip.source?.toLowerCase() || ''
  return ['png', 'jpg', 'jpeg', 'webp', 'svg', 'bmp', 'gif'].some((ext) => src.endsWith('.' + ext))
})

const nativeDuration = computed(() => {
  if (isStaticImage.value) return Infinity
  if (matchedAsset.value?.duration_seconds && matchedAsset.value.duration_seconds > 0) {
    return matchedAsset.value.duration_seconds
  }
  return Infinity
})

const thumbnailUrl = computed(() => {
  if (matchedAsset.value?.thumbnail_path) {
    return `/api/media/files/${matchedAsset.value.thumbnail_path}`
  }
  if (props.clip.source) {
    return `/api/media/files/${props.clip.source}`
  }
  return null
})

const clipStyle = computed(() => {
  const start = Number(props.clip.start) || 0
  const duration = Number(props.clip.duration) || 0
  const left = start * timelineStore.pixelsPerSecond
  const width = Math.max(10, duration * timelineStore.pixelsPerSecond)

  return {
    left: `${left}px`,
    width: `${width}px`,
  }
})

const trackKindColorClass = computed(() => {
  switch (props.trackKind) {
    case 'audio':
      return 'bg-gradient-to-r from-emerald-950/95 to-teal-900/95'
    case 'waveform':
      return 'bg-gradient-to-r from-purple-950/95 to-indigo-900/95'
    case 'subtitle':
      return 'bg-gradient-to-r from-amber-950/95 to-yellow-900/95'
    default:
      return 'bg-gradient-to-r from-blue-950/95 to-indigo-900/95'
  }
})

const clipIcon = computed(() => {
  switch (props.trackKind) {
    case 'audio':
      return 'i-heroicons-speaker-wave'
    case 'waveform':
      return 'i-heroicons-chart-bar'
    case 'subtitle':
      return 'i-heroicons-document-text'
    default:
      return 'i-heroicons-film'
  }
})

const clipTitle = computed(() => {
  if (props.clip.source) {
    return props.clip.source.split('/').pop() || props.clip.id
  }
  return props.clip.id
})

const formattedDuration = computed(() => {
  const d = Number(props.clip.duration) || 0
  return `${d.toFixed(1)}s`
})

// Right-Click Context Menu Items Definition
const contextMenuItems = computed(() => {
  const start = Number(props.clip.start) || 0
  const duration = Number(props.clip.duration) || 0

  return [
    // Header / Quick Info
    [
      {
        label: clipTitle.value,
        icon: clipIcon.value,
        disabled: true,
      },
      {
        label: `Start: ${start.toFixed(2)}s | Duration: ${duration.toFixed(2)}s`,
        icon: 'i-heroicons-clock',
        disabled: true,
      },
    ],
    // Primary Actions
    [
      {
        label: 'Duplicate Clip',
        icon: 'i-heroicons-document-duplicate',
        kbds: ['Ctrl', 'D'],
        onSelect: () => timelineStore.duplicateClip(props.clip.id),
      },
      {
        label: 'Split at Playhead',
        icon: 'i-heroicons-scissors',
        kbds: ['S'],
        disabled: playbackStore.currentTime <= start || playbackStore.currentTime >= start + duration,
        onSelect: () => timelineStore.splitClipAtPlayhead(playbackStore.currentTime),
      },
    ],
    // Info & Properties
    [
      {
        label: 'Inspect Properties',
        icon: 'i-heroicons-adjustments-horizontal',
        onSelect: () => {
          timelineStore.selectedClipId = props.clip.id
          timelineStore.selectedTrackId = props.trackId
        },
      },
      ...(matchedAsset.value
        ? [
            {
              label: 'View Probe Details',
              icon: 'i-heroicons-information-circle',
              onSelect: () => {
                if (matchedAsset.value) {
                  mediaStore.openProbeModal(matchedAsset.value)
                }
              },
            },
          ]
        : []),
    ],
    // Destructive
    [
      {
        label: 'Delete Clip',
        icon: 'i-heroicons-trash',
        color: 'error' as const,
        kbds: ['Del'],
        onSelect: () => timelineStore.removeClip(props.clip.id),
      },
    ],
  ]
})

function handleSelect() {
  timelineStore.selectedClipId = props.clip.id
  timelineStore.selectedTrackId = props.trackId
}

function handleDragStart(event: DragEvent) {
  if (isTrimming.value) {
    event.preventDefault()
    return
  }

  const targetEl = event.currentTarget as HTMLElement
  const rect = targetEl.getBoundingClientRect()
  const grabOffsetPixels = Math.max(0, event.clientX - rect.left)
  const grabOffsetTime = grabOffsetPixels / timelineStore.pixelsPerSecond

  timelineStore.selectedClipId = props.clip.id
  timelineStore.selectedTrackId = props.trackId
  timelineStore.draggedTimelineClip = {
    clip: props.clip,
    sourceTrackId: props.trackId,
    grabOffsetTime,
  }

  if (event.dataTransfer) {
    event.dataTransfer.setData(
      'application/vidonex-clip',
      JSON.stringify({
        clipId: props.clip.id,
        sourceTrackId: props.trackId,
        grabOffsetTime,
      })
    )
    event.dataTransfer.setData('text/plain', `clip:${props.trackId}:${props.clip.id}`)
    event.dataTransfer.effectAllowed = 'all'
    event.dataTransfer.setDragImage(targetEl, grabOffsetPixels, Math.max(0, event.clientY - rect.top))
  }
}

function handleDragEnd() {
  setTimeout(() => {
    timelineStore.draggedTimelineClip = null
  }, 150)
}

function startLeftTrim(event: MouseEvent) {
  isTrimming.value = true
  const startX = event.clientX
  const initialStart = Number(props.clip.start) || 0
  const initialDuration = Number(props.clip.duration) || 0
  const initialTrim = Number(props.clip.trim) || 0
  const originalEnd = initialStart + initialDuration

  const track = timelineStore.tracks.find((t) => t.id === props.trackId)
  const prevClips = (track?.clips || [])
    .filter((c) => c.id !== props.clip.id && Number(c.start) < initialStart)
    .sort((a, b) => Number(b.start) - Number(a.start))
  const prevClipEnd = prevClips.length > 0 ? (Number(prevClips[0].start) || 0) + (Number(prevClips[0].duration) || 0) : 0

  // Left trim: images can expand to 0; video/audio cannot expand before start of media (trim >= 0)
  const minStartFromTrim = isStaticImage.value ? 0 : initialStart - initialTrim
  const leftBoundary = Math.max(prevClipEnd, minStartFromTrim, 0)

  const onMouseMove = (moveEvent: MouseEvent) => {
    const deltaX = moveEvent.clientX - startX
    const deltaTime = deltaX / timelineStore.pixelsPerSecond

    const rawNewStart = initialStart + deltaTime
    const { time: snappedStart } = snapTime(rawNewStart, 0, props.clip.id)

    const clampedStart = Math.max(leftBoundary, Math.min(originalEnd - 0.2, snappedStart))
    const actualDelta = clampedStart - initialStart

    props.clip.start = clampedStart
    props.clip.duration = Math.max(0.2, initialDuration - actualDelta)
    if (!isStaticImage.value) {
      props.clip.trim = Math.max(0, initialTrim + actualDelta)
    }
  }

  const onMouseUp = () => {
    isTrimming.value = false
    timelineStore.pushHistoryState('Trim Clip Start')
    window.removeEventListener('mousemove', onMouseMove)
    window.removeEventListener('mouseup', onMouseUp)
  }

  window.addEventListener('mousemove', onMouseMove)
  window.addEventListener('mouseup', onMouseUp)
}

function startRightTrim(event: MouseEvent) {
  isTrimming.value = true
  const startX = event.clientX
  const initialDuration = Number(props.clip.duration) || 0
  const clipStart = Number(props.clip.start) || 0
  const currentTrim = Number(props.clip.trim) || 0

  const track = timelineStore.tracks.find((t) => t.id === props.trackId)
  const nextClips = (track?.clips || [])
    .filter((c) => c.id !== props.clip.id && Number(c.start) > clipStart)
    .sort((a, b) => Number(a.start) - Number(b.start))
  const nextClipStart = nextClips.length > 0 ? Number(nextClips[0].start) || Infinity : Infinity

  // Strict Max Duration enforcement: Video/Audio cannot exceed original native length
  let maxAllowedDuration = Infinity
  if (nativeDuration.value < Infinity) {
    maxAllowedDuration = Math.max(0.2, nativeDuration.value - currentTrim)
  }

  const maxEnd = Math.min(nextClipStart, clipStart + maxAllowedDuration)

  const onMouseMove = (moveEvent: MouseEvent) => {
    const deltaX = moveEvent.clientX - startX
    const deltaTime = deltaX / timelineStore.pixelsPerSecond

    const rawNewEnd = clipStart + initialDuration + deltaTime
    const { time: snappedEnd } = snapTime(rawNewEnd, 0, props.clip.id)

    const clampedEnd = Math.min(maxEnd, Math.max(clipStart + 0.2, snappedEnd))
    props.clip.duration = Math.max(0.2, clampedEnd - clipStart)
  }

  const onMouseUp = () => {
    isTrimming.value = false
    timelineStore.pushHistoryState('Trim Clip End')
    window.removeEventListener('mousemove', onMouseMove)
    window.removeEventListener('mouseup', onMouseUp)
  }

  window.addEventListener('mousemove', onMouseMove)
  window.addEventListener('mouseup', onMouseUp)
}
</script>
