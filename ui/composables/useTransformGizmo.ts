import { useTimelineStore } from '~/stores/timeline'
import { useMediaStore } from '~/stores/media'
import { useProjectStore } from '~/stores/project'
import type { ClipSpec } from '~/types/spec'

export type HandleType = 'tl' | 'tc' | 'tr' | 'mr' | 'br' | 'bc' | 'bl' | 'ml' | 'rot'

export interface ClipBounds {
  left: number
  top: number
  width: number
  height: number
  rotation: number
  opacity: number
}

// Global reactive smart guide state for viewport overlay
export const showGuideCenterX = ref(false)
export const showGuideCenterY = ref(false)
export const showGuideLeft = ref(false)
export const showGuideRight = ref(false)
export const showGuideTop = ref(false)
export const showGuideBottom = ref(false)

export function calculateClipBounds(
  clip: ClipSpec,
  canvasWidth: number,
  canvasHeight: number,
  displayWidth: number,
  displayHeight: number,
  mediaAssets: Array<{ id: string; file_name: string; file_path: string; width?: number; height?: number }>
): ClipBounds {
  // Find asset metadata for native aspect ratio
  const asset = mediaAssets.find(
    (a) => a.file_path === clip.source || a.file_name === clip.source || a.id === clip.source
  )

  let nativeW = asset?.width || 0
  let nativeH = asset?.height || 0

  // Default to canvas dimensions if no probe data yet
  if (!nativeW || !nativeH) {
    nativeW = canvasWidth
    nativeH = canvasHeight
  }

  const assetRatio = nativeW / nativeH

  // Contain fit within canvas to determine base dimensions on master canvas
  let baseCanvasW = canvasWidth
  let baseCanvasH = canvasWidth / assetRatio
  if (baseCanvasH > canvasHeight) {
    baseCanvasH = canvasHeight
    baseCanvasW = baseCanvasH * assetRatio
  }

  // Display viewport scaling factor
  const displayScale = canvasWidth > 0 ? displayWidth / canvasWidth : 1
  const scale = clip.scale || 1.0
  const rotation = clip.rotation || 0
  const posX = (clip.position?.x || 0) * displayScale
  const posY = (clip.position?.y || 0) * displayScale

  const displayW = baseCanvasW * scale * displayScale
  const displayH = baseCanvasH * scale * displayScale

  const left = (displayWidth - displayW) / 2 + posX
  const top = (displayHeight - displayH) / 2 + posY

  return {
    left,
    top,
    width: displayW,
    height: displayH,
    rotation,
    opacity: clip.opacity ?? 1.0,
  }
}

export function useTransformGizmo() {
  const timelineStore = useTimelineStore()
  const mediaStore = useMediaStore()
  const projectStore = useProjectStore()

  const isDragging = ref(false)
  const isResizing = ref(false)
  const isRotating = ref(false)
  const activeHandle = ref<HandleType | null>(null)

  const selectedClip = computed(() => timelineStore.selectedClip)

  function clearGuides() {
    showGuideCenterX.value = false
    showGuideCenterY.value = false
    showGuideLeft.value = false
    showGuideRight.value = false
    showGuideTop.value = false
    showGuideBottom.value = false
  }

  function startDrag(event: MouseEvent, displayScale: number) {
    if (!selectedClip.value) return
    isDragging.value = true

    const startX = event.clientX
    const startY = event.clientY
    const initialPos = {
      x: selectedClip.value.position?.x || 0,
      y: selectedClip.value.position?.y || 0,
    }

    const scaleFactor = displayScale > 0 ? 1 / displayScale : 1
    const canvasW = projectStore.canvasWidth
    const canvasH = projectStore.canvasHeight

    // Calculate native clip size in canvas units
    const bounds = calculateClipBounds(
      selectedClip.value,
      canvasW,
      canvasH,
      canvasW,
      canvasH,
      mediaStore.assets
    )
    const halfClipW = bounds.width / 2
    const halfClipH = bounds.height / 2

    const onMouseMove = (moveEvent: MouseEvent) => {
      if (!isDragging.value || !selectedClip.value) return

      const deltaX = (moveEvent.clientX - startX) * scaleFactor
      const deltaY = (moveEvent.clientY - startY) * scaleFactor

      let newX = initialPos.x + deltaX
      let newY = initialPos.y + deltaY

      const snapThreshold = 6 * scaleFactor

      // 1. Vertical Center Snapping (X = 0)
      if (Math.abs(newX) < snapThreshold) {
        newX = 0
        showGuideCenterX.value = true
        showGuideLeft.value = false
        showGuideRight.value = false
      }
      // 2. Left Edge Snapping (Left of clip touches Left of Canvas: newX - halfClipW = -canvasW / 2)
      else if (Math.abs(newX - halfClipW - (-canvasW / 2)) < snapThreshold) {
        newX = -canvasW / 2 + halfClipW
        showGuideLeft.value = true
        showGuideCenterX.value = false
        showGuideRight.value = false
      }
      // 3. Right Edge Snapping (Right of clip touches Right of Canvas: newX + halfClipW = canvasW / 2)
      else if (Math.abs(newX + halfClipW - canvasW / 2) < snapThreshold) {
        newX = canvasW / 2 - halfClipW
        showGuideRight.value = true
        showGuideCenterX.value = false
        showGuideLeft.value = false
      } else {
        showGuideCenterX.value = false
        showGuideLeft.value = false
        showGuideRight.value = false
      }

      // 4. Horizontal Center Snapping (Y = 0)
      if (Math.abs(newY) < snapThreshold) {
        newY = 0
        showGuideCenterY.value = true
        showGuideTop.value = false
        showGuideBottom.value = false
      }
      // 5. Top Edge Snapping (Top of clip touches Top of Canvas: newY - halfClipH = -canvasH / 2)
      else if (Math.abs(newY - halfClipH - (-canvasH / 2)) < snapThreshold) {
        newY = -canvasH / 2 + halfClipH
        showGuideTop.value = true
        showGuideCenterY.value = false
        showGuideBottom.value = false
      }
      // 6. Bottom Edge Snapping (Bottom of clip touches Bottom of Canvas: newY + halfClipH = canvasH / 2)
      else if (Math.abs(newY + halfClipH - canvasH / 2) < snapThreshold) {
        newY = canvasH / 2 - halfClipH
        showGuideBottom.value = true
        showGuideCenterY.value = false
        showGuideTop.value = false
      } else {
        showGuideCenterY.value = false
        showGuideTop.value = false
        showGuideBottom.value = false
      }

      if (!selectedClip.value.position) {
        selectedClip.value.position = { x: Math.round(newX), y: Math.round(newY), alignment: 'center' }
      } else {
        selectedClip.value.position.x = Math.round(newX)
        selectedClip.value.position.y = Math.round(newY)
      }
    }

    const onMouseUp = () => {
      isDragging.value = false
      clearGuides()
      timelineStore.pushHistoryState('Move Clip Position')
      window.removeEventListener('mousemove', onMouseMove)
      window.removeEventListener('mouseup', onMouseUp)
    }

    window.addEventListener('mousemove', onMouseMove)
    window.addEventListener('mouseup', onMouseUp)
  }

  function startResize(handle: HandleType, event: MouseEvent, bounds: ClipBounds) {
    if (!selectedClip.value) return
    event.stopPropagation()

    isResizing.value = true
    activeHandle.value = handle

    const startX = event.clientX
    const startY = event.clientY
    const initialScale = selectedClip.value.scale || 1.0

    // Center point in screen coordinates
    const stage = document.querySelector('.canvas-stage')
    const rect = stage?.getBoundingClientRect()
    const stageLeft = rect ? rect.left : 0
    const stageTop = rect ? rect.top : 0

    const centerX = stageLeft + bounds.left + bounds.width / 2
    const centerY = stageTop + bounds.top + bounds.height / 2
    const initialDist = Math.hypot(startX - centerX, startY - centerY)

    const onMouseMove = (moveEvent: MouseEvent) => {
      if (!isResizing.value || !selectedClip.value) return

      const currentDist = Math.hypot(moveEvent.clientX - centerX, moveEvent.clientY - centerY)
      if (initialDist > 0) {
        const ratio = currentDist / initialDist
        let newScale = Math.max(0.05, Math.min(10.0, initialScale * ratio))

        // Magnetic snap to 1.0 (100% native scale) - Subtle
        if (Math.abs(newScale - 1.0) < 0.015) {
          newScale = 1.0
          showGuideCenterX.value = true
          showGuideCenterY.value = true
        } else {
          showGuideCenterX.value = false
          showGuideCenterY.value = false
        }

        selectedClip.value.scale = Math.round(newScale * 100) / 100
      }
    }

    const onMouseUp = () => {
      isResizing.value = false
      activeHandle.value = null
      clearGuides()
      timelineStore.pushHistoryState('Resize Clip Scale')
      window.removeEventListener('mousemove', onMouseMove)
      window.removeEventListener('mouseup', onMouseUp)
    }

    window.addEventListener('mousemove', onMouseMove)
    window.addEventListener('mouseup', onMouseUp)
  }

  function startRotate(event: MouseEvent, bounds: ClipBounds) {
    if (!selectedClip.value) return
    event.stopPropagation()

    isRotating.value = true
    const stage = document.querySelector('.canvas-stage')
    const rect = stage?.getBoundingClientRect()
    const stageLeft = rect ? rect.left : 0
    const stageTop = rect ? rect.top : 0

    const centerX = stageLeft + bounds.left + bounds.width / 2
    const centerY = stageTop + bounds.top + bounds.height / 2

    const initialAngle = Math.atan2(event.clientY - centerY, event.clientX - centerX) * (180 / Math.PI)
    const initialClipRotation = selectedClip.value.rotation || 0

    const onMouseMove = (moveEvent: MouseEvent) => {
      if (!isRotating.value || !selectedClip.value) return

      const currentAngle = Math.atan2(moveEvent.clientY - centerY, moveEvent.clientX - centerX) * (180 / Math.PI)
      let deltaAngle = currentAngle - initialAngle
      let newRotation = (initialClipRotation + deltaAngle) % 360
      if (newRotation > 180) newRotation -= 360
      if (newRotation < -180) newRotation += 360

      // Snap to 0, 45, 90, 135, 180, -45, -90, -135, -180 degrees
      const snapAngles = [0, 45, 90, 135, 180, -45, -90, -135, -180]
      for (const snap of snapAngles) {
        if (Math.abs(newRotation - snap) < 4) {
          newRotation = snap
          break
        }
      }

      selectedClip.value.rotation = Math.round(newRotation)
    }

    const onMouseUp = () => {
      isRotating.value = false
      timelineStore.pushHistoryState('Rotate Clip')
      window.removeEventListener('mousemove', onMouseMove)
      window.removeEventListener('mouseup', onMouseUp)
    }

    window.addEventListener('mousemove', onMouseMove)
    window.addEventListener('mouseup', onMouseUp)
  }

  return {
    isDragging,
    isResizing,
    isRotating,
    activeHandle,
    showGuideCenterX,
    showGuideCenterY,
    showGuideLeft,
    showGuideRight,
    showGuideTop,
    showGuideBottom,
    startDrag,
    startResize,
    startRotate,
  }
}
