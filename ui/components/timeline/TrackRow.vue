<template>
  <div
    ref="trackRowRef"
    class="h-14 border-b border-gray-800/80 bg-gray-950/40 relative timeline-subgrid transition-colors min-w-full"
    :class="[
      isDragTarget && isCompatibleDrag && !isOverlapping ? 'bg-indigo-950/20' : '',
      isDragTarget && (!isCompatibleDrag || isOverlapping) ? 'bg-red-950/20 cursor-not-allowed' : ''
    ]"
    :style="{ width: `${totalTimelineWidth}px` }"
    @click="handleRowClick"
    @dragover.prevent="handleDragOver"
    @dragenter.prevent="handleDragEnter"
    @dragleave="handleDragLeave"
    @drop.prevent.stop="handleDrop"
  >
    <!-- Render Existing Clips -->
    <TimelineClipItem
      v-for="clip in track.clips || []"
      :key="clip.id"
      :clip="clip"
      :track-id="track.id"
      :track-kind="track.kind"
    />

    <!-- Render Transition Handles between Adjacent Clips (Video tracks) -->
    <template v-if="track.kind === 'video' || track.kind === 'overlay'">
      <TimelineTransitionHandle
        v-for="(pair, index) in adjacentClipPairs"
        :key="`trans_${index}`"
        :track-id="track.id"
        :from-clip="pair.from"
        :to-clip="pair.to"
        :transition="getTransition(pair.from.id, pair.to.id)"
      />
    </template>

    <!-- Animated Ghost Drop Preview (Compatible & No-Overlap Case) -->
    <div
      v-if="isDragTarget && isCompatibleDrag && !isOverlapping && ghostBounds && previewTitle"
      class="absolute top-1 bottom-1 rounded-lg border-2 border-indigo-400 timeline-drag-shimmer pointer-events-none z-30 flex items-center justify-between px-3 text-xs text-white shadow-[0_0_20px_rgba(99,102,241,0.4)] backdrop-blur-sm transition-all duration-75"
      :style="{ left: `${ghostBounds.left}px`, width: `${ghostBounds.width}px` }"
    >
      <div class="flex items-center gap-2 truncate">
        <!-- Floating Thumbnail / Icon -->
        <div class="w-6 h-6 rounded overflow-hidden flex-shrink-0 border border-indigo-300/60 bg-gray-900 shadow-md flex items-center justify-center">
          <img
            v-if="previewThumbnail"
            :src="resolveMediaUrl(previewThumbnail)"
            class="w-full h-full object-cover"
            alt=""
          />
          <UIcon
            v-else
            :name="previewIsAudio ? 'i-heroicons-speaker-wave' : 'i-heroicons-film'"
            class="w-4 h-4 text-indigo-300"
          />
        </div>

        <span class="text-[11px] font-sans font-semibold truncate text-white drop-shadow">
          {{ previewTitle }}
        </span>
      </div>

      <!-- Time & Snap Type Badge -->
      <div class="flex items-center gap-1 flex-shrink-0">
        <span
          v-if="ghostBounds.snapped"
          class="text-[9px] font-mono text-emerald-200 bg-emerald-950/80 border border-emerald-500/50 px-1.5 py-0.5 rounded shadow flex items-center gap-0.5"
        >
          <UIcon name="i-heroicons-bolt" class="w-3 h-3 text-emerald-400" />
          {{ ghostBounds.snapType === 'playhead' ? 'Playhead' : 'Magnet' }}
        </span>
        <span class="text-[10px] font-mono text-indigo-100 font-bold bg-indigo-600/90 px-2 py-0.5 rounded-full border border-indigo-300/50 shadow-md">
          {{ ghostBounds.duration.toFixed(1) }}s
        </span>
      </div>
    </div>

    <!-- Overlap Error Warning Overlay -->
    <div
      v-else-if="isDragTarget && isCompatibleDrag && isOverlapping && ghostBounds"
      class="absolute top-1 bottom-1 rounded-lg border-2 border-red-500 timeline-drag-incompatible pointer-events-none z-30 flex items-center justify-center gap-2 px-3 text-[11px] font-sans text-red-100 font-semibold shadow-[0_0_15px_rgba(239,68,68,0.5)] backdrop-blur-sm transition-all duration-75"
      :style="{ left: `${ghostBounds.left}px`, width: `${ghostBounds.width}px` }"
    >
      <UIcon name="i-heroicons-exclamation-triangle" class="w-4 h-4 text-red-300 flex-shrink-0 animate-bounce" />
      <span class="truncate drop-shadow">
        Cannot overlap with existing clip (Release cancels)
      </span>
    </div>

    <!-- Incompatible Track Type Warning Overlay -->
    <div
      v-else-if="isDragTarget && !isCompatibleDrag && ghostBounds"
      class="absolute top-1 bottom-1 rounded-lg border-2 border-red-500 timeline-drag-incompatible pointer-events-none z-30 flex items-center justify-center gap-2 px-3 text-[11px] font-sans text-red-100 font-semibold shadow-[0_0_15px_rgba(239,68,68,0.5)] backdrop-blur-sm transition-all duration-75"
      :style="{ left: `${ghostBounds.left}px`, width: `${ghostBounds.width}px` }"
    >
      <UIcon name="i-heroicons-no-symbol" class="w-4 h-4 text-red-300 flex-shrink-0 animate-bounce" />
      <span class="truncate drop-shadow">
        {{ incompatibleErrorMessage }}
      </span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useTimelineStore, isSourceCompatibleWithTrack } from '~/stores/timeline'
import { usePlaybackStore } from '~/stores/playback'
import { useSnapping } from '~/composables/useSnapping'
import { useDesktop } from '~/composables/useDesktop'
import type { TrackSpec, ClipSpec, TransitionSpec } from '~/types/spec'
import type { MediaAsset } from '~/types/project'

const { resolveMediaUrl } = useDesktop()

const props = defineProps<{
  track: TrackSpec
}>()

const timelineStore = useTimelineStore()
const playbackStore = usePlaybackStore()
const mediaStore = useMediaStore()
const { snapTime, checkTrackOverlap } = useSnapping()

const trackRowRef = ref<HTMLDivElement | null>(null)
const isDragTarget = ref(false)
const isOverlapping = ref(false)
const ghostBounds = ref<{
  left: number
  width: number
  duration: number
  snapped: boolean
  snapType?: string
} | null>(null)

// Check active drag source (Asset or existing Timeline Clip)
const activeAsset = computed(() => mediaStore.draggedAsset)
const activeTimelineClip = computed(() => timelineStore.draggedTimelineClip)

const previewSource = computed(() => {
  if (activeAsset.value) return activeAsset.value.file_path
  if (activeTimelineClip.value) return activeTimelineClip.value.clip.source
  return undefined
})

const previewTitle = computed(() => {
  if (activeAsset.value) return activeAsset.value.file_name
  if (activeTimelineClip.value) {
    const src = activeTimelineClip.value.clip.source
    return src ? src.split('/').pop() : activeTimelineClip.value.clip.id
  }
  return null
})

const previewThumbnail = computed(() => {
  if (activeAsset.value?.thumbnail_path) return activeAsset.value.thumbnail_path
  return null
})

const previewIsAudio = computed(() => {
  const src = previewSource.value?.toLowerCase() || ''
  return src.endsWith('.mp3') || src.endsWith('.wav') || src.endsWith('.aac') || src.endsWith('.flac') || props.track.kind === 'audio'
})

const isCompatibleDrag = computed(() => {
  const src = previewSource.value
  if (!src) return true
  return isSourceCompatibleWithTrack(src, props.track.kind)
})

const incompatibleErrorMessage = computed(() => {
  const src = previewSource.value?.toLowerCase() || ''
  const ext = src.split('.').pop() || ''
  const isImage = ['png', 'jpg', 'jpeg', 'webp', 'svg', 'bmp', 'gif'].includes(ext)
  const isAudio = ['mp3', 'wav', 'aac', 'flac', 'ogg', 'm4a'].includes(ext)

  if (isImage && (props.track.kind === 'audio' || props.track.kind === 'waveform')) {
    return 'Image cannot be placed in audio track'
  }
  if (isAudio && (props.track.kind === 'video' || props.track.kind === 'overlay')) {
    return 'Audio cannot be placed in video track'
  }
  return 'Incompatible track type'
})

const totalTimelineWidth = computed(() => timelineStore.timelineDuration * timelineStore.pixelsPerSecond)

const adjacentClipPairs = computed(() => {
  const clips = props.track.clips || []
  const sorted = [...clips].sort((a, b) => (Number(a.start) || 0) - (Number(b.start) || 0))
  const pairs: { from: ClipSpec; to: ClipSpec }[] = []

  for (let i = 0; i < sorted.length - 1; i++) {
    const current = sorted[i]
    const next = sorted[i + 1]
    if (current && next) {
      const currentEnd = (Number(current.start) || 0) + (Number(current.duration) || 0)
      const nextStart = Number(next.start) || 0

      if (Math.abs(currentEnd - nextStart) < 0.5) {
        pairs.push({ from: current, to: next })
      }
    }
  }

  return pairs
})

function getTransition(fromId: string, toId: string): TransitionSpec | undefined {
  return props.track.transitions?.find((t) => t.from === fromId && t.to === toId)
}

function handleRowClick() {
  timelineStore.selectedTrackId = props.track.id
}

function handleDragEnter() {
  isDragTarget.value = true
}

function handleDragOver(event: DragEvent) {
  event.preventDefault()
  if (!trackRowRef.value) return
  isDragTarget.value = true

  if (event.dataTransfer) {
    event.dataTransfer.dropEffect = 'copy'
  }

  const rect = trackRowRef.value.getBoundingClientRect()
  
  let duration = 5.0
  let clipIdToExclude: string | null = null
  let grabOffset = 0

  if (activeTimelineClip.value) {
    duration = Number(activeTimelineClip.value.clip.duration) || 5.0
    clipIdToExclude = activeTimelineClip.value.clip.id
    grabOffset = activeTimelineClip.value.grabOffsetTime || 0
  } else if (activeAsset.value) {
    if (activeAsset.value.file_type === 'image') {
      duration = 5.0
    } else if (activeAsset.value.duration_seconds > 0) {
      duration = activeAsset.value.duration_seconds
    }
  }

  const mouseX = event.clientX - rect.left - (grabOffset * timelineStore.pixelsPerSecond)
  const rawTime = Math.max(0, mouseX / timelineStore.pixelsPerSecond)

  // 1. Magnetic Snapping to playhead, 0s, and all other clip edges
  const { time: candidateStart, snapped, snapType } = snapTime(rawTime, duration, clipIdToExclude || undefined)

  // 2. Strict Overlap Validation for ghost indicator
  const { hasOverlap } = checkTrackOverlap(props.track, clipIdToExclude, candidateStart, duration)
  isOverlapping.value = hasOverlap

  ghostBounds.value = {
    left: candidateStart * timelineStore.pixelsPerSecond,
    width: Math.max(10, duration * timelineStore.pixelsPerSecond),
    duration,
    snapped,
    snapType,
  }
}

function handleDragLeave(event: DragEvent) {
  if (trackRowRef.value && !trackRowRef.value.contains(event.relatedTarget as Node)) {
    isDragTarget.value = false
    isOverlapping.value = false
    ghostBounds.value = null
  }
}

function handleDrop(event: DragEvent) {
  event.preventDefault()
  event.stopPropagation()
  isDragTarget.value = false
  isOverlapping.value = false
  ghostBounds.value = null

  if (!trackRowRef.value) return

  const rect = trackRowRef.value.getBoundingClientRect()

  // Case 1: Moving an existing timeline clip (Cross-track or same-track)
  let clipInfo = activeTimelineClip.value
  if (!clipInfo) {
    const rawClipData = event.dataTransfer?.getData('application/vidonex-clip')
    if (rawClipData) {
      try {
        const parsed = JSON.parse(rawClipData)
        const sourceTrack = timelineStore.tracks.find((t) => t.id === parsed.sourceTrackId)
        const foundClip = sourceTrack?.clips?.find((c) => c.id === parsed.clipId)
        if (foundClip && sourceTrack) {
          clipInfo = {
            clip: foundClip,
            sourceTrackId: sourceTrack.id,
            grabOffsetTime: parsed.grabOffsetTime || 0,
          }
        }
      } catch {}
    }
  }

  if (clipInfo) {
    const { clip, sourceTrackId, grabOffsetTime } = clipInfo
    if (!isSourceCompatibleWithTrack(clip.source, props.track.kind)) {
      timelineStore.draggedTimelineClip = null
      return
    }

    const duration = Number(clip.duration) || 5.0
    const grabOffset = grabOffsetTime || 0
    const mouseX = event.clientX - rect.left - (grabOffset * timelineStore.pixelsPerSecond)
    const rawTime = Math.max(0, mouseX / timelineStore.pixelsPerSecond)
    const { time: candidateStart } = snapTime(rawTime, duration, clip.id)

    // Strict zero-overlap enforcement: Cancel move if overlapping with another clip
    const { hasOverlap } = checkTrackOverlap(props.track, clip.id, candidateStart, duration)
    if (hasOverlap) {
      timelineStore.draggedTimelineClip = null
      return
    }

    timelineStore.moveClipToTrack(sourceTrackId, props.track.id, clip.id, candidateStart)
    timelineStore.draggedTimelineClip = null
    return
  }

  // Case 2: Adding an asset from Media Library
  let asset = activeAsset.value
  if (!asset) {
    const assetId = event.dataTransfer?.getData('vidonex-asset-id') || event.dataTransfer?.getData('text/plain')
    if (assetId) {
      asset = mediaStore.assets.find((a) => a.id === assetId || a.file_path === assetId || a.file_name === assetId) || null
    }
  }
  if (!asset) {
    const rawJson = event.dataTransfer?.getData('application/json')
    if (rawJson) {
      try {
        asset = JSON.parse(rawJson) as MediaAsset
      } catch {}
    }
  }

  // Case 3: Dropping native files directly from OS file manager onto track
  if (!asset && event.dataTransfer?.files && event.dataTransfer.files.length > 0) {
    const projectStore = useProjectStore()
    if (projectStore.currentProject) {
      const mouseX = event.clientX - rect.left
      const rawTime = Math.max(0, mouseX / timelineStore.pixelsPerSecond)
      let currentDropTime = rawTime

      for (const file of Array.from(event.dataTransfer.files)) {
        mediaStore.uploadFile(projectStore.currentProject.id, file).then((uploadedAsset) => {
          if (uploadedAsset && isSourceCompatibleWithTrack(uploadedAsset.file_path, props.track.kind)) {
            const dur = uploadedAsset.file_type === 'image' ? 5.0 : (uploadedAsset.duration_seconds || 5.0)
            const { hasOverlap } = checkTrackOverlap(props.track, null, currentDropTime, dur)
            if (!hasOverlap) {
              timelineStore.addClipToTrack(props.track.id, uploadedAsset, currentDropTime)
              currentDropTime += dur
            }
          }
        })
      }
    }
    return
  }

  if (!asset) {
    mediaStore.setDraggedAsset(null)
    return
  }

  // Check type compatibility
  if (!isSourceCompatibleWithTrack(asset.file_path || asset.file_name, props.track.kind)) {
    mediaStore.setDraggedAsset(null)
    return
  }

  let duration = 5.0
  if (asset.file_type === 'image') {
    duration = 5.0
  } else if (asset.duration_seconds > 0) {
    duration = asset.duration_seconds
  }

  const mouseX = event.clientX - rect.left
  const rawTime = Math.max(0, mouseX / timelineStore.pixelsPerSecond)
  const { time: candidateStart } = snapTime(rawTime, duration)

  // Strict zero-overlap enforcement: Cancel add if dropped over an existing clip
  const { hasOverlap } = checkTrackOverlap(props.track, null, candidateStart, duration)
  if (hasOverlap) {
    mediaStore.setDraggedAsset(null)
    return
  }

  timelineStore.addClipToTrack(props.track.id, asset, candidateStart)
  timelineStore.selectedTrackId = props.track.id
  mediaStore.setDraggedAsset(null)
}
</script>
