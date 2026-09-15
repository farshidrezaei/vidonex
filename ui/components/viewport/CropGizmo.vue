<template>
  <div
    v-if="selectedClip"
    class="absolute pointer-events-auto border-2 border-dashed border-amber-400 bg-amber-500/10 shadow-[0_0_20px_rgba(245,158,11,0.3)] select-none z-30"
    :style="cropGizmoStyle"
  >
    <!-- Crop Dimmed Inset Mask Visualizer -->
    <div class="absolute inset-0 pointer-events-none">
      <!-- Top Rule-of-thirds Line -->
      <div class="absolute left-0 right-0 top-1/3 border-b border-amber-300/30"></div>
      <!-- Bottom Rule-of-thirds Line -->
      <div class="absolute left-0 right-0 top-2/3 border-b border-amber-300/30"></div>
      <!-- Left Rule-of-thirds Line -->
      <div class="absolute top-0 bottom-0 left-1/3 border-r border-amber-300/30"></div>
      <!-- Right Rule-of-thirds Line -->
      <div class="absolute top-0 bottom-0 left-2/3 border-r border-amber-300/30"></div>
    </div>

    <!-- Thick Crop Corner Brackets -->
    <!-- Top-Left -->
    <div class="absolute -top-1 -left-1 w-4 h-4 border-t-4 border-l-4 border-amber-400 cursor-nwse-resize" @mousedown.stop="startCropDrag('tl', $event)"></div>
    <!-- Top-Right -->
    <div class="absolute -top-1 -right-1 w-4 h-4 border-t-4 border-r-4 border-amber-400 cursor-nesw-resize" @mousedown.stop="startCropDrag('tr', $event)"></div>
    <!-- Bottom-Left -->
    <div class="absolute -bottom-1 -left-1 w-4 h-4 border-b-4 border-l-4 border-amber-400 cursor-nesw-resize" @mousedown.stop="startCropDrag('bl', $event)"></div>
    <!-- Bottom-Right -->
    <div class="absolute -bottom-1 -right-1 w-4 h-4 border-b-4 border-r-4 border-amber-400 cursor-nwse-resize" @mousedown.stop="startCropDrag('br', $event)"></div>

    <!-- 4 Edge Midpoint Handles -->
    <!-- Top -->
    <div class="absolute -top-1.5 left-1/2 -translate-x-1/2 w-8 h-2.5 bg-amber-400 rounded-sm cursor-ns-resize shadow-md" @mousedown.stop="startCropDrag('t', $event)"></div>
    <!-- Bottom -->
    <div class="absolute -bottom-1.5 left-1/2 -translate-x-1/2 w-8 h-2.5 bg-amber-400 rounded-sm cursor-ns-resize shadow-md" @mousedown.stop="startCropDrag('b', $event)"></div>
    <!-- Left -->
    <div class="absolute top-1/2 -left-1.5 -translate-y-1/2 w-2.5 h-8 bg-amber-400 rounded-sm cursor-ew-resize shadow-md" @mousedown.stop="startCropDrag('l', $event)"></div>
    <!-- Right -->
    <div class="absolute top-1/2 -right-1.5 -translate-y-1/2 w-2.5 h-8 bg-amber-400 rounded-sm cursor-ew-resize shadow-md" @mousedown.stop="startCropDrag('r', $event)"></div>

    <!-- Floating Crop Action Controls Pill (Top Center) -->
    <div class="absolute -top-10 left-1/2 -translate-x-1/2 flex items-center gap-1.5 px-3 py-1 bg-gray-950/95 border border-amber-500/60 rounded-full shadow-2xl text-[11px] text-gray-200 z-40 whitespace-nowrap">
      <UIcon name="i-heroicons-scissors" class="w-3.5 h-3.5 text-amber-400" />
      <span class="font-semibold text-amber-300">{{ $t('crop.title') || 'Crop Active' }}</span>
      <span class="font-mono text-[10px] text-gray-400">
        T:{{ Math.round(cropInsets.top) }}% B:{{ Math.round(cropInsets.bottom) }}% L:{{ Math.round(cropInsets.left) }}% R:{{ Math.round(cropInsets.right) }}%
      </span>
      <button
        type="button"
        class="px-2 py-0.5 rounded bg-amber-500 hover:bg-amber-400 text-black font-semibold text-[10px] transition cursor-pointer"
        @click.stop="isCropping = false"
      >
        {{ $t('crop.done') || 'Done' }}
      </button>
      <button
        type="button"
        class="px-1.5 py-0.5 rounded bg-gray-800 hover:bg-gray-700 text-gray-300 text-[10px] transition cursor-pointer"
        @click.stop="resetCrop"
      >
        {{ $t('crop.reset') || 'Reset' }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useTimelineStore } from '~/stores/timeline'
import { useMediaStore } from '~/stores/media'
import { calculateClipBounds, isCropping } from '~/composables/useTransformGizmo'

const props = defineProps<{
  canvasWidth: number
  canvasHeight: number
  displayWidth: number
  displayHeight: number
}>()

const timelineStore = useTimelineStore()
const mediaStore = useMediaStore()
const selectedClip = computed(() => timelineStore.selectedClip)

const cropInsets = computed(() => {
  return {
    top: selectedClip.value?.crop?.top || 0,
    bottom: selectedClip.value?.crop?.bottom || 0,
    left: selectedClip.value?.crop?.left || 0,
    right: selectedClip.value?.crop?.right || 0,
  }
})

const currentBounds = computed(() => {
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

const cropGizmoStyle = computed(() => {
  const b = currentBounds.value
  const insets = cropInsets.value

  const insetLeftPx = (b.width * insets.left) / 100
  const insetRightPx = (b.width * insets.right) / 100
  const insetTopPx = (b.height * insets.top) / 100
  const insetBottomPx = (b.height * insets.bottom) / 100

  const actualLeft = b.left + insetLeftPx
  const actualTop = b.top + insetTopPx
  const actualWidth = Math.max(10, b.width - insetLeftPx - insetRightPx)
  const actualHeight = Math.max(10, b.height - insetTopPx - insetBottomPx)

  return {
    left: `${actualLeft}px`,
    top: `${actualTop}px`,
    width: `${actualWidth}px`,
    height: `${actualHeight}px`,
    transform: `rotate(${b.rotation}deg)`,
    transformOrigin: 'center center',
  }
})

function resetCrop() {
  if (!selectedClip.value) return
  selectedClip.value.crop = { top: 0, bottom: 0, left: 0, right: 0 }
  timelineStore.pushHistoryState('Reset Crop')
}

function startCropDrag(handle: string, startEvent: MouseEvent) {
  startEvent.stopPropagation()
  startEvent.preventDefault()

  if (!selectedClip.value) return
  if (!selectedClip.value.crop) {
    selectedClip.value.crop = { top: 0, bottom: 0, left: 0, right: 0 }
  }

  const startX = startEvent.clientX
  const startY = startEvent.clientY
  const originalCrop = { ...selectedClip.value.crop }
  const totalW = currentBounds.value.width || 1
  const totalH = currentBounds.value.height || 1

  function onMouseMove(moveEvent: MouseEvent) {
    const deltaX = moveEvent.clientX - startX
    const deltaY = moveEvent.clientY - startY

    const deltaLeftPct = (deltaX / totalW) * 100
    const deltaRightPct = (-deltaX / totalW) * 100
    const deltaTopPct = (deltaY / totalH) * 100
    const deltaBottomPct = (-deltaY / totalH) * 100

    if (!selectedClip.value || !selectedClip.value.crop) return

    if (handle.includes('l')) {
      const maxLeft = 100 - (originalCrop.right || 0) - 5
      selectedClip.value.crop.left = Math.max(0, Math.min(maxLeft, (originalCrop.left || 0) + deltaLeftPct))
    }
    if (handle.includes('r')) {
      const maxRight = 100 - (originalCrop.left || 0) - 5
      selectedClip.value.crop.right = Math.max(0, Math.min(maxRight, (originalCrop.right || 0) + deltaRightPct))
    }
    if (handle.includes('t')) {
      const maxTop = 100 - (originalCrop.bottom || 0) - 5
      selectedClip.value.crop.top = Math.max(0, Math.min(maxTop, (originalCrop.top || 0) + deltaTopPct))
    }
    if (handle.includes('b')) {
      const maxBottom = 100 - (originalCrop.top || 0) - 5
      selectedClip.value.crop.bottom = Math.max(0, Math.min(maxBottom, (originalCrop.bottom || 0) + deltaBottomPct))
    }
  }

  function onMouseUp() {
    window.removeEventListener('mousemove', onMouseMove)
    window.removeEventListener('mouseup', onMouseUp)
    timelineStore.pushHistoryState('Update Crop')
  }

  window.addEventListener('mousemove', onMouseMove)
  window.addEventListener('mouseup', onMouseUp)
}
</script>
