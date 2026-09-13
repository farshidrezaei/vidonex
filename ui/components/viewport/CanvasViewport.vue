<template>
  <div class="h-full flex flex-col bg-gray-950/90 relative overflow-hidden select-none">
    <!-- Viewport Canvas Stage -->
    <div ref="stageRef" class="flex-1 flex items-center justify-center p-3 relative overflow-hidden">
      <!-- Scaled Aspect Ratio Canvas Screen -->
      <div
        class="canvas-stage relative shadow-2xl transition-all rounded-sm flex items-center justify-center"
        :style="screenStyle"
      >
        <!-- Canvas Screen Frame (Clipped Media & Background) -->
        <div class="absolute inset-0 overflow-hidden rounded-sm pointer-events-none">
          <!-- Canvas Background (Color or Checkerboard) -->
          <div
            class="absolute inset-0 pointer-events-auto"
            :class="projectStore.backgroundColor === 'transparent' ? 'canvas-checkerboard' : ''"
            :style="{ backgroundColor: projectStore.backgroundColor }"
            @mousedown="timelineStore.selectedClipId = null"
          ></div>

          <!-- Render Active Visual Clips at Playhead -->
          <div
            v-for="clip in activeVisualClips"
            :key="clip.id"
            class="absolute pointer-events-auto select-none flex items-center justify-center cursor-move"
            :style="getClipRenderStyle(clip)"
            @mousedown.stop="handleClipMouseDown(clip, $event)"
          >
            <!-- Video / Image Asset Element -->
            <img
              v-if="isImage(clip.source)"
              :src="`/api/media/files/${clip.source}`"
              class="w-full h-full object-cover pointer-events-none select-none rounded-sm"
              alt=""
            />
            <video
              v-else-if="clip.source"
              :src="`/api/media/files/${clip.source}`"
              class="w-full h-full object-cover pointer-events-none select-none rounded-sm"
              :currentTime="getClipCurrentTime(clip)"
              muted
              playsinline
            ></video>
            <div
              v-else
              class="w-full h-full bg-indigo-500/20 border border-indigo-500/40 flex items-center justify-center text-xs font-mono text-indigo-300"
            >
              {{ clip.id }}
            </div>
          </div>
        </div>

        <!-- Smart Alignment Guidelines Overlay -->
        <!-- Center Vertical Line (X = 0) -->
        <div
          v-if="showGuideCenterX"
          class="absolute top-0 bottom-0 left-1/2 -translate-x-1/2 w-0.5 bg-cyan-400 z-20 pointer-events-none shadow-[0_0_10px_rgba(34,211,238,0.9)] animate-pulse"
        ></div>

        <!-- Center Horizontal Line (Y = 0) -->
        <div
          v-if="showGuideCenterY"
          class="absolute left-0 right-0 top-1/2 -translate-y-1/2 h-0.5 bg-cyan-400 z-20 pointer-events-none shadow-[0_0_10px_rgba(34,211,238,0.9)] animate-pulse"
        ></div>

        <!-- Center Pivot Magnet Crosshair Dot -->
        <div
          v-if="showGuideCenterX && showGuideCenterY"
          class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-3 h-3 rounded-full bg-white ring-4 ring-cyan-400/60 z-30 pointer-events-none shadow-[0_0_15px_rgba(34,211,238,1)] animate-ping"
        ></div>

        <!-- Left Canvas Border -->
        <div
          v-if="showGuideLeft"
          class="absolute top-0 bottom-0 left-0 w-1 bg-indigo-400 z-20 pointer-events-none shadow-[0_0_10px_rgba(129,140,248,0.9)]"
        ></div>

        <!-- Right Canvas Border -->
        <div
          v-if="showGuideRight"
          class="absolute top-0 bottom-0 right-0 w-1 bg-indigo-400 z-20 pointer-events-none shadow-[0_0_10px_rgba(129,140,248,0.9)]"
        ></div>

        <!-- Top Canvas Border -->
        <div
          v-if="showGuideTop"
          class="absolute left-0 right-0 top-0 h-1 bg-indigo-400 z-20 pointer-events-none shadow-[0_0_10px_rgba(129,140,248,0.9)]"
        ></div>

        <!-- Bottom Canvas Border -->
        <div
          v-if="showGuideBottom"
          class="absolute left-0 right-0 bottom-0 h-1 bg-indigo-400 z-20 pointer-events-none shadow-[0_0_10px_rgba(129,140,248,0.9)]"
        ></div>

        <!-- Transform Gizmo Overlay -->
        <ViewportTransformGizmo
          :canvas-width="projectStore.canvasWidth"
          :canvas-height="projectStore.canvasHeight"
          :display-width="displayDimensions.width"
          :display-height="displayDimensions.height"
        />
      </div>
    </div>

    <!-- Playback Transport Controls Toolbar -->
    <ViewportPlaybackToolbar />
  </div>
</template>

<script setup lang="ts">
import { useProjectStore } from '~/stores/project'
import { usePlaybackStore } from '~/stores/playback'
import { useTimelineStore } from '~/stores/timeline'
import { useMediaStore } from '~/stores/media'
import {
  showGuideCenterX,
  showGuideCenterY,
  showGuideLeft,
  showGuideRight,
  showGuideTop,
  showGuideBottom,
  calculateClipBounds,
  useTransformGizmo,
} from '~/composables/useTransformGizmo'
import {
  computeClipTransitionState,
  getTransitionStyleModifiers,
} from '~/composables/useTransitionPreview'
import { useTimelineAudio } from '~/composables/useTimelineAudio'
import type { ClipSpec, TrackSpec } from '~/types/spec'

const projectStore = useProjectStore()
const playbackStore = usePlaybackStore()
const timelineStore = useTimelineStore()
const mediaStore = useMediaStore()
const { startDrag } = useTransformGizmo()

// Initialize multi-track synchronized audio playback engine for live preview
useTimelineAudio()

function handleClipMouseDown(clip: ClipSpec, event: MouseEvent) {
  event.stopPropagation()
  event.preventDefault()
  timelineStore.selectedClipId = clip.id
  const displayScale = projectStore.canvasWidth > 0 ? displayDimensions.value.width / projectStore.canvasWidth : 1
  startDrag(event, displayScale)
}

const stageRef = ref<HTMLDivElement | null>(null)
const stageDimensions = useElementSize(stageRef)

// Padding offset corresponding to p-3 padding (12px each side = 24px total)
const VIEWPORT_PADDING_OFFSET = 24

const displayDimensions = computed(() => {
  const stageW = stageDimensions.width.value
  const stageH = stageDimensions.height.value

  const availableW = Math.max(50, stageW - VIEWPORT_PADDING_OFFSET)
  const availableH = Math.max(50, stageH - VIEWPORT_PADDING_OFFSET)

  const canvasW = projectStore.canvasWidth || 1920
  const canvasH = projectStore.canvasHeight || 1080
  const targetRatio = canvasW / canvasH
  const availableRatio = availableW / availableH

  let width: number
  let height: number

  // Dynamic Fit:
  // If the available area is proportionally narrower than the aspect ratio (availableRatio < targetRatio),
  // the video is constrained by available width -> Fit on the X axis.
  // If the available area is proportionally wider than the aspect ratio (availableRatio >= targetRatio),
  // the video is constrained by available height -> Fit on the Y axis.
  if (availableRatio < targetRatio) {
    width = availableW
    height = width / targetRatio
  } else {
    height = availableH
    width = height * targetRatio
  }

  return {
    width: Math.round(width),
    height: Math.round(height),
  }
})

const screenStyle = computed(() => ({
  width: `${displayDimensions.value.width}px`,
  height: `${displayDimensions.value.height}px`,
}))

function findTrackForClip(clipId: string): TrackSpec | undefined {
  return timelineStore.tracks.find((t) => t.clips?.some((c) => c.id === clipId))
}

const activeVisualClips = computed(() => {
  const time = playbackStore.currentTime
  const visualTracks = timelineStore.tracks.filter((t) => t.kind === 'video' || t.kind === 'overlay')
  const clips: ClipSpec[] = []

  for (const track of visualTracks) {
    for (const clip of track.clips || []) {
      const start = Number(clip.start) || 0
      const duration = Number(clip.duration) || 0

      const end = start + duration
      let isVisible = time >= start && time <= end

      // Check if clip is active as an outgoing or incoming clip in a centered transition
      if (!isVisible && track.transitions && track.transitions.length > 0) {
        // 1. Outgoing clip remains visible during the transition's second half [cutTime, cutTime + halfDuration]
        const outgoingTransition = track.transitions.find((t) => t.from === clip.id)
        if (outgoingTransition) {
          const transDuration = Math.max(0.01, Number(outgoingTransition.duration) || 1.0)
          const halfDuration = transDuration / 2
          const cutTime = end
          if (time >= cutTime && time <= cutTime + halfDuration) {
            isVisible = true
          }
        }

        // 2. Incoming clip becomes visible during the transition's first half [cutTime - halfDuration, cutTime]
        const incomingTransition = track.transitions.find((t) => t.to === clip.id)
        if (incomingTransition) {
          const partnerClip = track.clips?.find((c) => c.id === incomingTransition.from)
          if (partnerClip) {
            const cutTime = (Number(partnerClip.start) || 0) + (Number(partnerClip.duration) || 0)
            const transDuration = Math.max(0.01, Number(incomingTransition.duration) || 1.0)
            const halfDuration = transDuration / 2
            if (time >= cutTime - halfDuration && time <= cutTime) {
              isVisible = true
            }
          }
        }
      }

      if (isVisible) {
        clips.push(clip)
      }
    }
  }

  return clips
})

function getClipRenderStyle(clip: ClipSpec) {
  const bounds = calculateClipBounds(
    clip,
    projectStore.canvasWidth,
    projectStore.canvasHeight,
    displayDimensions.value.width,
    displayDimensions.value.height,
    mediaStore.assets
  )

  const track = findTrackForClip(clip.id)
  const transitionInfo = track
    ? computeClipTransitionState(clip, track, playbackStore.currentTime)
    : { isActive: false, role: 'none' as const, progress: 0 }
  const modifiers = getTransitionStyleModifiers(transitionInfo)

  let opacity = bounds.opacity
  if (modifiers.opacity !== undefined) {
    opacity *= modifiers.opacity
  }

  let transform = `rotate(${bounds.rotation}deg)`
  if (modifiers.transformExtra) {
    transform = `${transform} ${modifiers.transformExtra}`
  }

  const style: Record<string, any> = {
    left: `${bounds.left}px`,
    top: `${bounds.top}px`,
    width: `${bounds.width}px`,
    height: `${bounds.height}px`,
    transform,
    opacity,
    mixBlendMode: (clip.blend_mode && clip.blend_mode !== 'normal') ? clip.blend_mode : 'normal',
    transformOrigin: 'center center',
    zIndex: (track?.z_index || 0) + (modifiers.zIndexExtra || 0),
  }

  if (modifiers.clipPath) {
    style.clipPath = modifiers.clipPath
  }

  return style
}

function isImage(source?: string): boolean {
  if (!source) return false
  const ext = source.toLowerCase().split('.').pop() || ''
  return ['png', 'jpg', 'jpeg', 'webp', 'svg', 'gif', 'bmp'].includes(ext)
}

function getClipCurrentTime(clip: ClipSpec): number {
  const start = Number(clip.start) || 0
  const trim = Number(clip.trim) || 0
  const speed = Number(clip.speed) || 1.0

  const elapsed = (playbackStore.currentTime - start) * speed
  return Math.max(0, trim + elapsed)
}
</script>
