import { defineStore } from 'pinia'
import { usePlaybackStore } from '~/stores/playback'
import { useProjectStore } from '~/stores/project'
import type { VideoSpec, TrackSpec, ClipSpec, TransitionSpec, TrackKind } from '~/types/spec'
import type { MediaAsset } from '~/types/project'

export interface DraggedClipInfo {
  clip: ClipSpec
  sourceTrackId: string
  grabOffsetTime?: number
}

export function isSourceCompatibleWithTrack(source: string | undefined, targetKind: TrackKind): boolean {
  if (!source) return true
  const ext = source.toLowerCase().split('.').pop() || ''
  const isImage = ['png', 'jpg', 'jpeg', 'webp', 'svg', 'bmp', 'gif'].includes(ext)
  const isAudio = ['mp3', 'wav', 'aac', 'flac', 'ogg', 'm4a'].includes(ext)

  if (isImage) {
    return targetKind === 'video' || targetKind === 'overlay'
  }
  if (isAudio) {
    return targetKind === 'audio' || targetKind === 'waveform'
  }
  return true // Video files can go to video, overlay, audio, or waveform
}

export function findNonOverlappingStart(
  track: TrackSpec,
  clipIdToExclude: string | null,
  targetStart: number,
  duration: number
): number {
  const safeStart = Math.max(0, targetStart)
  const otherClips = (track.clips || [])
    .filter((c) => c.id !== clipIdToExclude)
    .map((c) => ({
      id: c.id,
      start: Number(c.start) || 0,
      duration: Number(c.duration) || 0,
      end: (Number(c.start) || 0) + (Number(c.duration) || 0),
    }))
    .sort((a, b) => a.start - b.start)

  if (otherClips.length === 0) {
    return safeStart
  }

  // Check if safeStart has an overlap with any other clip
  const hasOverlap = (start: number) => {
    const end = start + duration
    return otherClips.some((c) => Math.max(start, c.start) < Math.min(end, c.end) - 0.001)
  }

  if (!hasOverlap(safeStart)) {
    return safeStart
  }

  // Find candidate non-overlapping slot closest to safeStart
  const candidateStarts: number[] = [0]
  for (const c of otherClips) {
    candidateStarts.push(c.end) // Right edge of clip c
    if (c.start >= duration) {
      candidateStarts.push(c.start - duration) // Left edge of clip c
    }
  }

  const validCandidates = candidateStarts.filter((pos) => pos >= 0 && !hasOverlap(pos))

  if (validCandidates.length > 0) {
    validCandidates.sort((a, b) => Math.abs(a - safeStart) - Math.abs(b - safeStart))
    return validCandidates[0]
  }

  // Fallback: Place at the very end of the last clip
  const lastClip = otherClips[otherClips.length - 1]
  return lastClip ? lastClip.end : safeStart
}

export const useTimelineStore = defineStore('timeline', () => {
  const playbackStore = usePlaybackStore()
  const projectStore = useProjectStore()

  const tracks = ref<TrackSpec[]>([
    {
      id: 'track_video_1',
      kind: 'video',
      clips: [],
      transitions: [],
    },
    {
      id: 'track_audio_1',
      kind: 'audio',
      clips: [],
    },
  ])

  const selectedClipId = ref<string | null>(null)
  const selectedTrackId = ref<string | null>(null)
  const pixelsPerSecond = ref<number>(50) // Zoom level
  const isSnappingEnabled = ref<boolean>(true)

  // Dragged clip for cross-track movement
  const draggedTimelineClip = ref<DraggedClipInfo | null>(null)

  // Undo / Redo history stacks
  const undoStack = ref<{ tracks: string; description: string }[]>([])
  const redoStack = ref<{ tracks: string; description: string }[]>([])
  const maxHistoryLength = 50

  const selectedClip = computed(() => {
    if (!selectedClipId.value) return null
    for (const track of tracks.value) {
      const clip = track.clips?.find((c) => c.id === selectedClipId.value)
      if (clip) return clip
    }
    return null
  })

  const selectedTrack = computed(() => {
    if (!selectedTrackId.value) return null
    return tracks.value.find((track) => track.id === selectedTrackId.value) || null
  })

  // Total duration of all tracks
  const calculatedDuration = computed(() => {
    let maxTime = 10 // Minimum 10 seconds empty timeline
    for (const track of tracks.value) {
      for (const clip of track.clips || []) {
        const start = Number(clip.start) || 0
        const duration = Number(clip.duration) || 0
        if (start + duration > maxTime) {
          maxTime = start + duration
        }
      }
    }
    return maxTime
  })

  watch(calculatedDuration, (newDuration) => {
    playbackStore.duration = Math.max(10, Math.ceil(newDuration) + 2)
  }, { immediate: true })

  let isInternalLoading = false
  let autoSaveTimeout: ReturnType<typeof setTimeout> | null = null

  function saveCurrentTimeline() {
    if (isInternalLoading) return
    projectStore.saveCurrentProject(toVideoSpec())
  }

  function triggerAutoSave() {
    if (isInternalLoading) return
    if (autoSaveTimeout) clearTimeout(autoSaveTimeout)
    autoSaveTimeout = setTimeout(() => {
      saveCurrentTimeline()
    }, 300)
  }

  watch(
    tracks,
    () => {
      triggerAutoSave()
    },
    { deep: true }
  )

  function pushHistoryState(description = 'Edit Timeline') {
    undoStack.value.push({
      tracks: JSON.stringify(tracks.value),
      description,
    })
    if (undoStack.value.length > maxHistoryLength) {
      undoStack.value.shift()
    }
    redoStack.value = [] // Clear redo stack on new mutation

    // Trigger immediate auto-save
    saveCurrentTimeline()
  }

  function undo() {
    if (undoStack.value.length === 0) return
    const currentState = { tracks: JSON.stringify(tracks.value), description: 'Current' }
    redoStack.value.push(currentState)

    const previousState = undoStack.value.pop()
    if (previousState) {
      tracks.value = JSON.parse(previousState.tracks)
      projectStore.saveCurrentProject(toVideoSpec())
    }
  }

  function redo() {
    if (redoStack.value.length === 0) return
    const currentState = { tracks: JSON.stringify(tracks.value), description: 'Current' }
    undoStack.value.push(currentState)

    const nextState = redoStack.value.pop()
    if (nextState) {
      tracks.value = JSON.parse(nextState.tracks)
      projectStore.saveCurrentProject(toVideoSpec())
    }
  }

  function addTrack(kind: TrackKind) {
    pushHistoryState(`Add ${kind} track`)
    const id = `track_${kind}_${tracks.value.length + 1}`
    tracks.value.push({
      id,
      kind,
      clips: [],
      transitions: kind === 'video' || kind === 'overlay' ? [] : undefined,
    })
    selectedTrackId.value = id
  }

  function removeTrack(trackId: string) {
    pushHistoryState('Remove track')
    tracks.value = tracks.value.filter((track) => track.id !== trackId)
    if (selectedTrackId.value === trackId) {
      selectedTrackId.value = null
    }
  }

  function addClipToTrack(trackId: string, asset: MediaAsset, startTime?: number) {
    const track = tracks.value.find((t) => t.id === trackId)
    if (!track) return

    // Type compatibility check
    if (!isSourceCompatibleWithTrack(asset.file_path, track.kind)) {
      console.warn(`Asset ${asset.file_name} is incompatible with ${track.kind} track`)
      return
    }

    pushHistoryState(`Add clip ${asset.file_name}`)

    let duration = 5.0
    if (asset.file_type === 'image') {
      duration = 5.0
    } else if (asset.duration_seconds > 0) {
      duration = asset.duration_seconds
    }

    const requestedStart = startTime !== undefined ? startTime : getNextAvailableStartTime(track)
    // Non-overlapping collision prevention
    const start = findNonOverlappingStart(track, null, requestedStart, duration)

    const newClip: ClipSpec = {
      id: `clip_${Date.now().toString(36)}`,
      source: asset.file_path,
      start,
      duration,
      speed: 1.0,
      scale: 1.0,
      rotation: 0,
      opacity: 1.0,
      volume: 1.0,
      blend_mode: 'normal',
      position: {
        x: 0,
        y: 0,
        alignment: 'center',
      },
    }

    if (!track.clips) {
      track.clips = []
    }
    track.clips.push(newClip)
    selectedClipId.value = newClip.id

    // Expand total timeline duration if needed
    const clipEnd = start + duration
    if (clipEnd + 10 > playbackStore.duration) {
      playbackStore.duration = Math.ceil(clipEnd + 10)
    }
  }

  function moveClipToTrack(sourceTrackId: string, targetTrackId: string, clipId: string, targetStart: number): boolean {
    const sourceTrack = tracks.value.find((t) => t.id === sourceTrackId)
    const targetTrack = tracks.value.find((t) => t.id === targetTrackId)
    if (!sourceTrack || !targetTrack || !sourceTrack.clips) return false

    const clipIndex = sourceTrack.clips.findIndex((c) => c.id === clipId)
    if (clipIndex === -1) return false

    const clip = sourceTrack.clips[clipIndex]

    // Type compatibility check
    if (!isSourceCompatibleWithTrack(clip.source, targetTrack.kind)) {
      console.warn(`Clip ${clip.id} is incompatible with ${targetTrack.kind} track`)
      return false
    }

    const duration = Number(clip.duration) || 5.0
    // Prevent collision and overlapping in target track
    const validStart = findNonOverlappingStart(targetTrack, clipId, targetStart, duration)

    pushHistoryState(`Move clip to ${targetTrack.id}`)

    // Remove from source track
    sourceTrack.clips.splice(clipIndex, 1)

    // Update start time
    clip.start = validStart

    // Add to target track
    if (!targetTrack.clips) {
      targetTrack.clips = []
    }
    targetTrack.clips.push(clip)
    selectedClipId.value = clip.id
    selectedTrackId.value = targetTrack.id

    // Expand total timeline duration if needed
    const clipEnd = validStart + duration
    if (clipEnd + 10 > playbackStore.duration) {
      playbackStore.duration = Math.ceil(clipEnd + 10)
    }

    return true
  }

  function getNextAvailableStartTime(track: TrackSpec): number {
    if (!track.clips || track.clips.length === 0) return 0
    let lastEnd = 0
    for (const clip of track.clips) {
      const end = (Number(clip.start) || 0) + (Number(clip.duration) || 0)
      if (end > lastEnd) {
        lastEnd = end
      }
    }
    return lastEnd
  }

  function duplicateClip(clipId: string) {
    for (const track of tracks.value) {
      if (!track.clips) continue
      const clip = track.clips.find((c) => c.id === clipId)
      if (!clip) continue

      pushHistoryState(`Duplicate clip`)
      const duration = Number(clip.duration) || 5.0
      const preferredStart = (Number(clip.start) || 0) + duration
      const validStart = findNonOverlappingStart(track, null, preferredStart, duration)

      const duplicatedClip: ClipSpec = {
        ...JSON.parse(JSON.stringify(clip)),
        id: `clip_${Date.now().toString(36)}`,
        start: validStart,
      }

      track.clips.push(duplicatedClip)
      selectedClipId.value = duplicatedClip.id
      selectedTrackId.value = track.id
      return
    }
  }

  function removeClip(clipId: string) {
    pushHistoryState('Delete clip')
    for (const track of tracks.value) {
      if (track.clips) {
        track.clips = track.clips.filter((clip) => clip.id !== clipId)
      }
      if (track.transitions) {
        track.transitions = track.transitions.filter((t) => t.from !== clipId && t.to !== clipId)
      }
    }
    if (selectedClipId.value === clipId) {
      selectedClipId.value = null
    }
  }

  function splitClipAtPlayhead(playheadTime: number) {
    if (!selectedClipId.value) return

    for (const track of tracks.value) {
      if (!track.clips) continue
      const clipIndex = track.clips.findIndex((c) => c.id === selectedClipId.value)
      if (clipIndex === -1) continue

      const clip = track.clips[clipIndex]
      const start = Number(clip.start) || 0
      const duration = Number(clip.duration) || 0

      if (playheadTime > start && playheadTime < start + duration) {
        pushHistoryState('Split clip')

        const firstDuration = playheadTime - start
        const secondDuration = duration - firstDuration
        const originalTrim = Number(clip.trim) || 0

        clip.duration = firstDuration

        const secondClip: ClipSpec = {
          ...JSON.parse(JSON.stringify(clip)),
          id: `clip_${Date.now().toString(36)}`,
          start: playheadTime,
          duration: secondDuration,
          trim: originalTrim + firstDuration,
        }

        track.clips.splice(clipIndex + 1, 0, secondClip)
        selectedClipId.value = secondClip.id
        return
      }
    }
  }

  function addTransition(trackId: string, fromClipId: string, toClipId: string, type = 'dissolve', duration = 1.0) {
    const track = tracks.value.find((t) => t.id === trackId)
    if (!track) return

    const dur = Math.max(0.1, parseFloat(String(duration)) || 1.0)

    if (!track.transitions) {
      track.transitions = []
    }

    const existingIndex = track.transitions.findIndex((t) => t.from === fromClipId && t.to === toClipId)
    const newTransition: TransitionSpec = {
      from: fromClipId,
      to: toClipId,
      type,
      duration: dur,
    }

    if (existingIndex !== -1) {
      track.transitions[existingIndex] = newTransition
    } else {
      track.transitions.push(newTransition)
    }

    pushHistoryState(`Add transition ${type} (${dur}s)`)
    saveCurrentTimeline()
  }

  function loadFromSpec(spec: VideoSpec) {
    if (spec.tracks && spec.tracks.length > 0) {
      isInternalLoading = true
      const parsedTracks: TrackSpec[] = JSON.parse(JSON.stringify(spec.tracks))
      for (const track of parsedTracks) {
        if (track.transitions) {
          for (const trans of track.transitions) {
            if (trans.duration !== undefined) {
              trans.duration = parseFloat(String(trans.duration)) || 1.0
            }
          }
        }
      }
      tracks.value = parsedTracks
      selectedClipId.value = null
      selectedTrackId.value = tracks.value[0]?.id || null
      undoStack.value = []
      redoStack.value = []
      setTimeout(() => {
        isInternalLoading = false
      }, 100)
    }
  }

  function toVideoSpec(): VideoSpec {
    const exportedTracks = tracks.value.map((track, index) => ({
      ...track,
      z_index: track.z_index !== undefined ? track.z_index : index,
    }))
    return {
      version: '1.0',
      canvas: {
        width: projectStore.canvasWidth,
        height: projectStore.canvasHeight,
        frame_rate: projectStore.frameRate,
        background_color: projectStore.backgroundColor,
      },
      tracks: JSON.parse(JSON.stringify(exportedTracks)),
    }
  }

  return {
    tracks,
    selectedClipId,
    selectedTrackId,
    selectedClip,
    selectedTrack,
    pixelsPerSecond,
    isSnappingEnabled,
    draggedTimelineClip,
    undoStack,
    redoStack,
    pushHistoryState,
    undo,
    redo,
    addTrack,
    removeTrack,
    addClipToTrack,
    moveClipToTrack,
    removeClip,
    duplicateClip,
    splitClipAtPlayhead,
    addTransition,
    loadFromSpec,
    toVideoSpec,
  }
})
