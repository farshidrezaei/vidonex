<template>
  <div
    v-if="selectedClip"
    class="absolute pointer-events-auto border-2 border-indigo-500 cursor-move transition-shadow shadow-lg shadow-indigo-500/30"
    :style="gizmoStyle"
    @mousedown="handleDrag"
  >
    <!-- Center Pivot Dot -->
    <div class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-2 h-2 rounded-full bg-indigo-400 pointer-events-none"></div>

    <!-- Rotation Handle Stem & Knob -->
    <div class="absolute -top-6 left-1/2 -translate-x-1/2 w-0.5 h-6 bg-indigo-500 pointer-events-none"></div>
    <div
      class="absolute -top-8 left-1/2 -translate-x-1/2 w-4 h-4 rounded-full bg-white border-2 border-indigo-500 cursor-grab active:cursor-grabbing hover:scale-125 transition-transform"
      @mousedown="handleRotate"
    ></div>

    <!-- 8 Resize Handles -->
    <!-- Top-Left -->
    <div
      class="absolute -top-1.5 -left-1.5 w-3 h-3 bg-white border-2 border-indigo-500 cursor-nwse-resize hover:scale-125 transition-transform"
      @mousedown="handleResize('tl', $event)"
    ></div>

    <!-- Top-Center -->
    <div
      class="absolute -top-1.5 left-1/2 -translate-x-1/2 w-3 h-3 bg-white border-2 border-indigo-500 cursor-ns-resize hover:scale-125 transition-transform"
      @mousedown="handleResize('tc', $event)"
    ></div>

    <!-- Top-Right -->
    <div
      class="absolute -top-1.5 -right-1.5 w-3 h-3 bg-white border-2 border-indigo-500 cursor-nesw-resize hover:scale-125 transition-transform"
      @mousedown="handleResize('tr', $event)"
    ></div>

    <!-- Middle-Right -->
    <div
      class="absolute top-1/2 -right-1.5 -translate-y-1/2 w-3 h-3 bg-white border-2 border-indigo-500 cursor-ew-resize hover:scale-125 transition-transform"
      @mousedown="handleResize('mr', $event)"
    ></div>

    <!-- Bottom-Right -->
    <div
      class="absolute -bottom-1.5 -right-1.5 w-3 h-3 bg-white border-2 border-indigo-500 cursor-nwse-resize hover:scale-125 transition-transform"
      @mousedown="handleResize('br', $event)"
    ></div>

    <!-- Bottom-Center -->
    <div
      class="absolute -bottom-1.5 left-1/2 -translate-x-1/2 w-3 h-3 bg-white border-2 border-indigo-500 cursor-ns-resize hover:scale-125 transition-transform"
      @mousedown="handleResize('bc', $event)"
    ></div>

    <!-- Bottom-Left -->
    <div
      class="absolute -bottom-1.5 -left-1.5 w-3 h-3 bg-white border-2 border-indigo-500 cursor-nesw-resize hover:scale-125 transition-transform"
      @mousedown="handleResize('bl', $event)"
    ></div>

    <!-- Middle-Left -->
    <div
      class="absolute top-1/2 -left-1.5 -translate-y-1/2 w-3 h-3 bg-white border-2 border-indigo-500 cursor-ew-resize hover:scale-125 transition-transform"
      @mousedown="handleResize('ml', $event)"
    ></div>
  </div>
</template>

<script setup lang="ts">
import { useTimelineStore } from '~/stores/timeline'
import { useMediaStore } from '~/stores/media'
import { useTransformGizmo, calculateClipBounds, type HandleType, type ClipBounds } from '~/composables/useTransformGizmo'

const props = defineProps<{
  canvasWidth: number
  canvasHeight: number
  displayWidth: number
  displayHeight: number
}>()

const timelineStore = useTimelineStore()
const mediaStore = useMediaStore()
const { startDrag, startResize, startRotate } = useTransformGizmo()

const selectedClip = computed(() => timelineStore.selectedClip)

const currentBounds = computed<ClipBounds>(() => {
  if (!selectedClip.value) {
    return { left: 0, top: 0, width: 0, height: 0, rotation: 0, opacity: 1 }
  }
  return calculateClipBounds(
    selectedClip.value,
    props.canvasWidth,
    props.canvasHeight,
    props.displayWidth,
    props.displayHeight,
    mediaStore.assets
  )
})

const gizmoStyle = computed(() => {
  if (!selectedClip.value) return {}
  const b = currentBounds.value

  return {
    left: `${b.left}px`,
    top: `${b.top}px`,
    width: `${b.width}px`,
    height: `${b.height}px`,
    transform: `rotate(${b.rotation}deg)`,
    transformOrigin: 'center center',
  }
})

const displayScale = computed(() => (props.canvasWidth > 0 ? props.displayWidth / props.canvasWidth : 1))

function handleDrag(event: MouseEvent) {
  startDrag(event, displayScale.value)
}

function handleResize(handle: HandleType, event: MouseEvent) {
  startResize(handle, event, currentBounds.value)
}

function handleRotate(event: MouseEvent) {
  startRotate(event, currentBounds.value)
}
</script>
