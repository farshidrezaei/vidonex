<template>
  <div
    ref="viewportContainerRef"
    class="h-full flex flex-col bg-gray-950/90 relative overflow-hidden select-none"
  >
    <!-- Viewport Canvas Stage -->
    <div
      ref="stageRef"
      class="flex-1 flex items-center justify-center p-3 relative overflow-hidden"
      :class="isSpacePressed && isPanning ? 'cursor-grabbing' : (isSpacePressed ? 'cursor-grab' : 'cursor-default')"
      @wheel.prevent="handleWheel"
      @mousedown="handleStageMouseDown"
      @dblclick="resetZoom"
    >
      <!-- Floating Pan/Zoom Reset Overlay Pill (visible when panned or zoomed) -->
      <div
        v-if="!isFit"
        class="absolute top-3 right-3 z-30 flex items-center gap-1.5 px-2.5 py-1 rounded-full bg-gray-900/80 backdrop-blur border border-gray-700/60 shadow-lg text-xs font-mono text-gray-300 pointer-events-auto transition-all"
      >
        <span>{{ zoomDisplay }}</span>
        <button
          class="px-1.5 py-0.5 rounded bg-indigo-600/70 hover:bg-indigo-500 text-white text-[10px] font-sans font-medium transition-colors"
          title="Reset View to Fit"
          @click.stop="resetZoom"
        >
          Fit
        </button>
      </div>

      <!-- Scaled Aspect Ratio Canvas Screen -->
      <div
        class="canvas-stage relative shadow-2xl rounded-sm flex items-center justify-center"
        :class="isPanning || isDragging || isResizing || isRotating ? '!transition-none' : 'transition-transform duration-75'"
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
              :src="getMediaSourceUrl(clip.source)"
              class="w-full h-full object-cover pointer-events-none select-none rounded-sm"
              alt=""
            />
            <video
              v-else-if="clip.source"
              :ref="(el) => setVideoRef(clip.id, el as HTMLVideoElement | null)"
              :src="getMediaSourceUrl(clip.source)"
              class="w-full h-full object-cover pointer-events-none select-none rounded-sm"
              muted
              playsinline
              preload="auto"
              @loadedmetadata="handleVideoLoaded(clip, $event)"
            ></video>
            <div
              v-else
              class="w-full h-full bg-indigo-500/20 border border-indigo-500/40 flex items-center justify-center text-xs font-mono text-indigo-300"
            >
              {{ clip.id }}
            </div>
          </div>

          <!-- Render Active Waveform Visualizers -->
          <ViewportWaveformVisualizer
            v-for="(layer, idx) in activeWaveformLayers"
            :key="`${layer.track.id}_${layer.clip?.id || idx}`"
            :track="layer.track"
            :clip="layer.clip"
            :canvas-width="projectStore.canvasWidth"
            :canvas-height="projectStore.canvasHeight"
            :display-width="displayDimensions.width"
            :display-height="displayDimensions.height"
          />
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

        <!-- Transform & Crop Gizmo Overlays -->
        <ViewportCropGizmo
          v-if="isCropping"
          :canvas-width="projectStore.canvasWidth"
          :canvas-height="projectStore.canvasHeight"
          :display-width="displayDimensions.width"
          :display-height="displayDimensions.height"
        />
        <ViewportTransformGizmo
          v-else
          :canvas-width="projectStore.canvasWidth"
          :canvas-height="projectStore.canvasHeight"
          :display-width="displayDimensions.width"
          :display-height="displayDimensions.height"
        />
      </div>
    </div>

    <!-- Playback Transport Controls Toolbar -->
    <ViewportPlaybackToolbar @toggle-fullscreen="onToggleFullscreen" />
  </div>
</template>

<script setup lang="ts">
import { useProjectStore } from '~/stores/project'
import { usePlaybackStore } from '~/stores/playback'
import { useTimelineStore } from '~/stores/timeline'
import { useMediaStore } from '~/stores/media'
import { useViewportZoom } from '~/composables/useViewportZoom'
import {
  showGuideCenterX,
  showGuideCenterY,
  showGuideLeft,
  showGuideRight,
  showGuideTop,
  showGuideBottom,
  calculateClipBounds,
  isDragging,
  isResizing,
  isRotating,
  isCropping,
  useTransformGizmo,
} from '~/composables/useTransformGizmo'
import {
  computeClipTransitionState,
  getTransitionStyleModifiers,
} from '~/composables/useTransitionPreview'
import { useTimelineAudio } from '~/composables/useTimelineAudio'
import { useDesktop } from '~/composables/useDesktop'
import type { ClipSpec, TrackSpec } from '~/types/spec'

const projectStore = useProjectStore()
const playbackStore = usePlaybackStore()
const timelineStore = useTimelineStore()
const mediaStore = useMediaStore()
const { resolveMediaUrl } = useDesktop()
const { startDrag } = useTransformGizmo()
const {
  zoom,
  panX,
  panY,
  isSpacePressed,
  isPanning,
  isFit,
  zoomDisplay,
  resetZoom,
  handleWheel,
  startPan,
  toggleFullscreen,
} = useViewportZoom()

const viewportContainerRef = ref<HTMLDivElement | null>(null)

function onToggleFullscreen() {
  toggleFullscreen(viewportContainerRef.value)
}

function handleStageMouseDown(event: MouseEvent) {
  // If space is held down or middle click (button === 1), start panning
  if (isSpacePressed.value || event.button === 1) {
    startPan(event)
    return
  }

  // If clicked directly on the empty background (outside canvas screen), deselect active clip
  if (event.target === stageRef.value) {
    timelineStore.selectedClipId = null
  }
}

// Initialize multi-track synchronized audio playback engine for live preview
useTimelineAudio()

function handleClipMouseDown(clip: ClipSpec, event: MouseEvent) {
  // If space is pressed or middle click, pan tool takes precedence over element dragging
  if (isSpacePressed.value || event.button === 1) {
    startPan(event)
    return
  }

  event.stopPropagation()
  event.preventDefault()
  timelineStore.selectedClipId = clip.id
  const baseScale = projectStore.canvasWidth > 0 ? displayDimensions.value.width / projectStore.canvasWidth : 1
  const displayScale = baseScale * zoom.value
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
  transform: `translate3d(${panX.value}px, ${panY.value}px, 0) scale(${zoom.value})`,
  transformOrigin: 'center center',
}))

function handleKeyDown(event: KeyboardEvent) {
  const target = event.target as HTMLElement
  if (target && (['INPUT', 'TEXTAREA', 'SELECT'].includes(target.tagName) || target.isContentEditable)) {
    return
  }
  if (event.code === 'Space' || event.key === ' ') {
    isSpacePressed.value = true
  }
}

function handleKeyUp(event: KeyboardEvent) {
  if (event.code === 'Space' || event.key === ' ') {
    isSpacePressed.value = false
  }
}

onMounted(() => {
  window.addEventListener('keydown', handleKeyDown)
  window.addEventListener('keyup', handleKeyUp)
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeyDown)
  window.removeEventListener('keyup', handleKeyUp)
})

function findTrackForClip(clipId: string): TrackSpec | undefined {
  return timelineStore.tracks.find((t) => t.clips?.some((c) => c.id === clipId))
}

function findTrackIndexForClip(clipId: string): number {
  return timelineStore.tracks.findIndex((t) => t.clips?.some((c) => c.id === clipId))
}

function getTrackBaseZIndex(trackId: string): number {
  const trackIndex = timelineStore.tracks.findIndex((t) => t.id === trackId)
  if (trackIndex === -1) return 10
  const track = timelineStore.tracks[trackIndex]
  // Base z-index derived from user-configured z_index (if provided) or its natural timeline index.
  // We allocate 100 integer slots per track layer (e.g. Track 0: 100, Track 1: 200, Track 2: 300...)
  // This completely prevents transition zIndexExtra (+5) from ever bleeding into higher tracks.
  const layerIndex = track.z_index !== undefined ? track.z_index : trackIndex
  return (layerIndex + 1) * 100
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

  // Stable sort clips by their effective track layer to ensure consistent DOM paint order
  return clips.sort((a, b) => {
    const trackIndexA = findTrackIndexForClip(a.id)
    const trackIndexB = findTrackIndexForClip(b.id)
    const trackA = trackIndexA !== -1 ? timelineStore.tracks[trackIndexA] : undefined
    const trackB = trackIndexB !== -1 ? timelineStore.tracks[trackIndexB] : undefined
    const zA = trackA?.z_index !== undefined ? trackA.z_index : trackIndexA
    const zB = trackB?.z_index !== undefined ? trackB.z_index : trackIndexB
    if (zA !== zB) {
      return zA - zB
    }
    return (Number(a.start) || 0) - (Number(b.start) || 0)
  })
})

interface ActiveWaveformLayer {
  track: TrackSpec
  clip?: ClipSpec | null
}

const activeWaveformLayers = computed<ActiveWaveformLayer[]>(() => {
  const time = playbackStore.currentTime
  const waveformTracks = timelineStore.tracks.filter((t) => t.kind === 'waveform')
  const layers: ActiveWaveformLayer[] = []

  for (const track of waveformTracks) {
    if (track.muted) continue
    if (track.clips && track.clips.length > 0) {
      for (const clip of track.clips) {
        const start = Number(clip.start) || 0
        const duration = Number(clip.duration) || 0
        if (time >= start && time <= start + duration) {
          layers.push({ track, clip })
        }
      }
    } else {
      // Track-level waveform with no specific clips
      layers.push({ track, clip: null })
    }
  }

  return layers
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

  const baseZIndex = track ? getTrackBaseZIndex(track.id) : 100

  const style: Record<string, any> = {
    left: `${bounds.left}px`,
    top: `${bounds.top}px`,
    width: `${bounds.width}px`,
    height: `${bounds.height}px`,
    transform,
    opacity,
    mixBlendMode: (clip.blend_mode && clip.blend_mode !== 'normal') ? clip.blend_mode : 'normal',
    transformOrigin: 'center center',
    zIndex: baseZIndex + (modifiers.zIndexExtra || 0),
  }

  // 1. Apply Live Color Grading Filter
  const filterParts: string[] = []
  if (clip.color_grading) {
    const cg = clip.color_grading
    if (cg.brightness && cg.brightness !== 0) {
      filterParts.push(`brightness(${1 + cg.brightness})`)
    }
    if (cg.contrast && cg.contrast !== 1) {
      filterParts.push(`contrast(${cg.contrast})`)
    }
    if (cg.saturation && cg.saturation !== 1) {
      filterParts.push(`saturate(${cg.saturation})`)
    }
    if (cg.temperature && cg.temperature !== 0) {
      // Warmth shift: positive adds slight sepia/warm hue, negative cool hue
      const hueShift = cg.temperature * 15
      filterParts.push(`hue-rotate(${hueShift}deg)`)
    }
    if (cg.blur && cg.blur > 0) {
      const canvasScale = displayDimensions.value.width / (projectStore.canvasWidth || 1920)
      const displayBlur = Math.max(0.5, cg.blur * 0.2 * canvasScale)
      filterParts.push(`blur(${displayBlur.toFixed(1)}px)`)
    }
    // Highlights & Whites adjustment
    const highlightEffect = (cg.highlights || 0) * 0.25 + (cg.whites || 0) * 0.15
    if (highlightEffect !== 0) {
      filterParts.push(`brightness(${Math.max(0.1, 1 + highlightEffect).toFixed(2)})`)
    }
    // Shadows & Blacks adjustment
    const shadowEffect = (cg.shadows || 0) * 0.2 - (cg.blacks || 0) * 0.15
    if (shadowEffect !== 0) {
      filterParts.push(`contrast(${Math.max(0.2, 1 - shadowEffect).toFixed(2)})`)
    }
    // Sharpen edge contrast enhancement
    if (cg.sharpen && cg.sharpen > 0) {
      const sharpenContrast = 1 + (cg.sharpen / 100) * 0.35
      filterParts.push(`contrast(${sharpenContrast.toFixed(2)})`)
    }
  }
  if (filterParts.length > 0) {
    style.filter = filterParts.join(' ')
  }

  // 2. Apply Inset Crop
  if (clip.crop && (clip.crop.top || clip.crop.bottom || clip.crop.left || clip.crop.right)) {
    const top = clip.crop.top || 0
    const right = clip.crop.right || 0
    const bottom = clip.crop.bottom || 0
    const left = clip.crop.left || 0
    style.clipPath = `inset(${top}% ${right}% ${bottom}% ${left}%)`
  } else if (modifiers.clipPath) {
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

function getMediaSourceUrl(source?: string): string {
  return resolveMediaUrl(source)
}

const videoElements = new Map<string, HTMLVideoElement>()

function setVideoRef(clipId: string, el: HTMLVideoElement | null) {
  if (el) {
    videoElements.set(clipId, el)
    syncSingleVideo(clipId, el)
  } else {
    videoElements.delete(clipId)
  }
}

function handleVideoLoaded(clip: ClipSpec, event: Event) {
  const video = event.target as HTMLVideoElement
  if (!video) return
  syncSingleVideo(clip.id, video)
}

function syncSingleVideo(clipId: string, video: HTMLVideoElement) {
  const clip = activeVisualClips.value.find((c) => c.id === clipId)
  if (!clip) return

  const targetTime = getClipCurrentTime(clip)
  const isPlaying = playbackStore.isPlaying
  const speed = Math.max(0.1, Number(clip.speed) || 1.0)
  video.playbackRate = Math.max(0.25, Math.min(4.0, speed * playbackStore.playbackRate))

  if (isPlaying) {
    const drift = Math.abs(video.currentTime - targetTime)
    if (drift > 0.2 || video.paused) {
      video.currentTime = targetTime
    }
    if (video.paused) {
      video.play().catch(() => {
        // Suppress autoplay policy errors on background tab
      })
    }
  } else {
    if (!video.paused) {
      video.pause()
    }
    if (Math.abs(video.currentTime - targetTime) > 0.04) {
      video.currentTime = targetTime
    }
  }
}

function syncAllVideos() {
  for (const [clipId, video] of videoElements.entries()) {
    syncSingleVideo(clipId, video)
  }
}

watch(() => playbackStore.isPlaying, (playing) => {
  if (!playing) {
    for (const video of videoElements.values()) {
      if (!video.paused) {
        try {
          video.pause()
        } catch {
          // Ignore
        }
      }
    }
  }
  syncAllVideos()
})

watch(() => playbackStore.currentTime, () => {
  syncAllVideos()
})

watch(activeVisualClips, (clips) => {
  const activeIds = new Set(clips.map((c) => c.id))
  for (const [id, video] of videoElements.entries()) {
    if (!activeIds.has(id)) {
      if (!video.paused) {
        try {
          video.pause()
        } catch {
          // Ignore
        }
      }
      videoElements.delete(id)
    }
  }
  nextTick(() => syncAllVideos())
}, { deep: true })

onUnmounted(() => {
  for (const video of videoElements.values()) {
    try {
      video.pause()
      video.src = ''
    } catch {
      // Ignore
    }
  }
  videoElements.clear()
})
</script>
