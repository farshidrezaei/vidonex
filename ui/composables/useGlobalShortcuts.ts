import { useTimelineStore } from '~/stores/timeline'
import { usePlaybackStore } from '~/stores/playback'
import { useProjectStore } from '~/stores/project'

export const isShortcutsModalOpen = ref(false)

export function useGlobalShortcuts() {
  const timelineStore = useTimelineStore()
  const playbackStore = usePlaybackStore()
  const projectStore = useProjectStore()

  function handleKeyDown(event: KeyboardEvent) {
    const target = event.target as HTMLElement
    // Ignore when user is typing in form inputs, textareas, selects, or contenteditable
    if (
      target &&
      (['INPUT', 'TEXTAREA', 'SELECT'].includes(target.tagName) ||
        target.isContentEditable ||
        target.closest('input, textarea, [contenteditable="true"]'))
    ) {
      return
    }

    const isCtrlOrCmd = event.ctrlKey || event.metaKey
    const code = event.code
    const key = event.key ? event.key.toLowerCase() : ''

    // 1. Undo: Ctrl+Z / Cmd+Z (without Shift)
    if (isCtrlOrCmd && (code === 'KeyZ' || key === 'z' || key === 'ظ') && !event.shiftKey) {
      event.preventDefault()
      event.stopPropagation()
      timelineStore.undo()
      return
    }

    // 2. Redo: Ctrl+Y / Cmd+Y OR Ctrl+Shift+Z / Cmd+Shift+Z
    if (
      (isCtrlOrCmd && (code === 'KeyY' || key === 'y' || key === 'غ')) ||
      (isCtrlOrCmd && event.shiftKey && (code === 'KeyZ' || key === 'z' || key === 'ظ'))
    ) {
      event.preventDefault()
      event.stopPropagation()
      timelineStore.redo()
      return
    }

    // 3. Duplicate: Ctrl+D / Cmd+D
    if (isCtrlOrCmd && (code === 'KeyD' || key === 'd' || key === 'ی')) {
      event.preventDefault()
      event.stopPropagation()
      if (timelineStore.selectedClipId) {
        timelineStore.duplicateClip(timelineStore.selectedClipId)
      }
      return
    }

    // 4. Play / Pause toggle: Space
    if (code === 'Space' || key === ' ' || key === 'spacebar') {
      event.preventDefault()
      event.stopPropagation()
      playbackStore.togglePlay()
      return
    }

    // 5. Split Clip: S or C
    if (!isCtrlOrCmd && !event.altKey && (code === 'KeyS' || code === 'KeyC' || key === 's' || key === 'c' || key === 'س' || key === 'ژ')) {
      event.preventDefault()
      event.stopPropagation()
      timelineStore.splitClipAtPlayhead(playbackStore.currentTime)
      return
    }

    // 6. Delete Clip: Delete or Backspace
    if (code === 'Delete' || code === 'Backspace' || key === 'delete' || key === 'backspace') {
      event.preventDefault()
      event.stopPropagation()
      if (timelineStore.selectedClipId) {
        timelineStore.removeClip(timelineStore.selectedClipId)
      }
      return
    }

    // 7. Move Selected Clip with Arrow Keys in Preview / Canvas
    if (timelineStore.selectedClip) {
      const step = event.shiftKey ? 10 : 1
      const clip = timelineStore.selectedClip

      if (code === 'ArrowUp' || key === 'arrowup') {
        event.preventDefault()
        event.stopPropagation()
        if (!clip.position) clip.position = { x: 0, y: 0, alignment: 'center' }
        clip.position.y = (clip.position.y || 0) - step
        timelineStore.pushHistoryState('Nudge Clip Up')
        return
      }

      if (code === 'ArrowDown' || key === 'arrowdown') {
        event.preventDefault()
        event.stopPropagation()
        if (!clip.position) clip.position = { x: 0, y: 0, alignment: 'center' }
        clip.position.y = (clip.position.y || 0) + step
        timelineStore.pushHistoryState('Nudge Clip Down')
        return
      }

      if (code === 'ArrowLeft' || key === 'arrowleft') {
        event.preventDefault()
        event.stopPropagation()
        if (!clip.position) clip.position = { x: 0, y: 0, alignment: 'center' }
        clip.position.x = (clip.position.x || 0) - step
        timelineStore.pushHistoryState('Nudge Clip Left')
        return
      }

      if (code === 'ArrowRight' || key === 'arrowright') {
        event.preventDefault()
        event.stopPropagation()
        if (!clip.position) clip.position = { x: 0, y: 0, alignment: 'center' }
        clip.position.x = (clip.position.x || 0) + step
        timelineStore.pushHistoryState('Nudge Clip Right')
        return
      }
    }

    // 8. Step Backward 1 frame: Left Arrow (when no clip selected) or J
    if (!isCtrlOrCmd && (code === 'ArrowLeft' || code === 'KeyJ' || key === 'arrowleft' || key === 'j' || key === 'ت')) {
      event.preventDefault()
      event.stopPropagation()
      playbackStore.stepFrame(-1, projectStore.frameRate)
      return
    }

    // 9. Step Forward 1 frame: Right Arrow (when no clip selected) or L
    if (!isCtrlOrCmd && (code === 'ArrowRight' || code === 'KeyL' || key === 'arrowright' || key === 'l' || key === 'م')) {
      event.preventDefault()
      event.stopPropagation()
      playbackStore.stepFrame(1, projectStore.frameRate)
      return
    }

    // 10. Jump to Beginning: Home or 0
    if (!isCtrlOrCmd && (code === 'Home' || code === 'Numpad0')) {
      event.preventDefault()
      event.stopPropagation()
      playbackStore.seek(0)
      return
    }

    // 11. Toggle Snapping: N
    if (!isCtrlOrCmd && (code === 'KeyN' || key === 'n' || key === 'د')) {
      event.preventDefault()
      event.stopPropagation()
      timelineStore.isSnappingEnabled = !timelineStore.isSnappingEnabled
      return
    }

    // 11. Zoom In: + or =
    if (!isCtrlOrCmd && (code === 'Equal' || code === 'NumpadAdd' || key === '+' || key === '=')) {
      event.preventDefault()
      event.stopPropagation()
      timelineStore.pixelsPerSecond = Math.min(200, timelineStore.pixelsPerSecond + 15)
      return
    }

    // 12. Zoom Out: -
    if (!isCtrlOrCmd && (code === 'Minus' || code === 'NumpadSubtract' || key === '-' || key === '_')) {
      event.preventDefault()
      event.stopPropagation()
      timelineStore.pixelsPerSecond = Math.max(15, timelineStore.pixelsPerSecond - 15)
      return
    }

    // 13. Deselect: Escape
    if (code === 'Escape' || key === 'escape') {
      timelineStore.selectedClipId = null
      return
    }

    // 14. Cheat sheet Help: ? or Ctrl+/
    if (key === '?' || (isCtrlOrCmd && (code === 'Slash' || key === '/'))) {
      event.preventDefault()
      event.stopPropagation()
      isShortcutsModalOpen.value = !isShortcutsModalOpen.value
    }
  }

  function registerGlobalListeners() {
    if (typeof window !== 'undefined') {
      window.addEventListener('keydown', handleKeyDown, { capture: true })
    }
  }

  function unregisterGlobalListeners() {
    if (typeof window !== 'undefined') {
      window.removeEventListener('keydown', handleKeyDown, { capture: true })
    }
  }

  return {
    isShortcutsModalOpen,
    registerGlobalListeners,
    unregisterGlobalListeners,
  }
}
