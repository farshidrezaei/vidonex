<template>
  <div class="h-full flex flex-col bg-gray-950/90 relative overflow-hidden select-none">
    <!-- Viewport Canvas Stage -->
    <div ref="stageRef" class="flex-1 flex items-center justify-center p-6 relative overflow-hidden">
      <!-- Scaled Aspect Ratio Canvas Screen -->
      <div
        class="canvas-stage relative shadow-2xl transition-all rounded-sm overflow-hidden flex items-center justify-center"
        :style="screenStyle"
      >
        <!-- Canvas Background (Color or Checkerboard) -->
        <div
          class="absolute inset-0"
          :class="projectStore.backgroundColor === 'transparent' ? 'canvas-checkerboard' : ''"
          :style="{ backgroundColor: projectStore.backgroundColor }"
        ></div>

        <!-- Render Active Visual Clips at Playhead -->
        <div
          v-for="clip in activeVisualClips"
          :key="clip.id"
          class="absolute pointer-events-auto select-none flex items-center justify-center"
          :style="getClipRenderStyle(clip)"
          @click.stop="timelineStore.selectedClipId = clip.id"
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
} from '~/composables/useTransformGizmo'
import type { ClipSpec } from '~/types/spec'

const projectStore = useProjectStore()
const playbackStore = usePlaybackStore()
const timelineStore = useTimelineStore()
const mediaStore = useMediaStore()

const stageRef = ref<HTMLDivElement | null>(null)
const stageDimensions = useElementSize(stageRef)

const displayDimensions = computed(() => {
  const availableW = Math.max(100, stageDimensions.width.value - 48)
  const availableH = Math.max(100, stageDimensions.height.value - 48)

  const canvasW = projectStore.canvasWidth
  const canvasH = projectStore.canvasHeight
  const targetRatio = canvasW / canvasH

  let width = availableW
  let height = width / targetRatio

  if (height > availableH) {
    height = availableH
    width = height * targetRatio
  }

  return { width, height }
})

const screenStyle = computed(() => ({
  width: `${displayDimensions.value.width}px`,
  height: `${displayDimensions.value.height}px`,
}))

const activeVisualClips = computed(() => {
  const time = playbackStore.currentTime
  const visualTracks = timelineStore.tracks.filter((t) => t.kind === 'video' || t.kind === 'overlay')
  const clips: ClipSpec[] = []

  for (const track of visualTracks) {
    for (const clip of track.clips || []) {
      const start = Number(clip.start) || 0
      const duration = Number(clip.duration) || 0
      if (time >= start && time <= start + duration) {
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

  return {
    left: `${bounds.left}px`,
    top: `${bounds.top}px`,
    width: `${bounds.width}px`,
    height: `${bounds.height}px`,
    transform: `rotate(${bounds.rotation}deg)`,
    opacity: bounds.opacity,
    transformOrigin: 'center center',
  }
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
