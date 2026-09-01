import { useTimelineStore } from '~/stores/timeline'
import { usePlaybackStore } from '~/stores/playback'
import type { TrackSpec } from '~/types/spec'

export function useSnapping() {
  const timelineStore = useTimelineStore()
  const playbackStore = usePlaybackStore()

  function getSnapPoints(excludeClipId?: string): number[] {
    const points = new Set<number>()
    points.add(0) // Timeline start
    points.add(playbackStore.currentTime) // Playhead position

    for (const track of timelineStore.tracks) {
      for (const clip of track.clips || []) {
        if (clip.id === excludeClipId) continue
        const start = Number(clip.start) || 0
        const duration = Number(clip.duration) || 0
        points.add(start)
        points.add(start + duration)
      }
    }

    return Array.from(points).sort((a, b) => a - b)
  }

  function snapTime(
    targetStart: number,
    duration = 0,
    excludeClipId?: string,
    thresholdPixels = 6
  ): { time: number; snapped: boolean; snapType?: 'start' | 'end' | 'playhead' } {
    if (timelineStore.isSnappingEnabled === false) {
      return { time: Math.max(0, targetStart), snapped: false }
    }

    const thresholdSeconds = thresholdPixels / timelineStore.pixelsPerSecond
    const points = getSnapPoints(excludeClipId)

    let bestStart = Math.max(0, targetStart)
    let minDiff = Infinity
    let snapped = false
    let snapType: 'start' | 'end' | 'playhead' | undefined = undefined

    // 1. Check left edge snapping (start -> snapPoint)
    for (const point of points) {
      const diff = Math.abs(point - targetStart)
      if (diff <= thresholdSeconds && diff < minDiff) {
        minDiff = diff
        bestStart = Math.max(0, point)
        snapped = true
        snapType = point === playbackStore.currentTime ? 'playhead' : 'start'
      }
    }

    // 2. Check right edge snapping (start + duration -> snapPoint)
    if (duration > 0) {
      const targetEnd = targetStart + duration
      for (const point of points) {
        const diff = Math.abs(point - targetEnd)
        if (diff <= thresholdSeconds && diff < minDiff) {
          const candidateStart = point - duration
          if (candidateStart >= 0) {
            minDiff = diff
            bestStart = candidateStart
            snapped = true
            snapType = point === playbackStore.currentTime ? 'playhead' : 'end'
          }
        }
      }
    }

    return { time: bestStart, snapped, snapType }
  }

  function checkTrackOverlap(
    track: TrackSpec,
    clipIdToExclude: string | null | undefined,
    start: number,
    duration: number
  ): { hasOverlap: boolean; overlappingClipId?: string } {
    const targetStart = Math.max(0, start)
    const targetEnd = targetStart + duration
    const epsilon = 0.005 // Tolerance for floating-point tangent/touching edges

    for (const clip of track.clips || []) {
      if (clip.id === clipIdToExclude) continue
      const cStart = Number(clip.start) || 0
      const cEnd = cStart + (Number(clip.duration) || 0)

      // Strict overlap check: interval intersection length > epsilon
      const overlapStart = Math.max(targetStart, cStart)
      const overlapEnd = Math.min(targetEnd, cEnd)

      if (overlapStart < overlapEnd - epsilon) {
        return { hasOverlap: true, overlappingClipId: clip.id }
      }
    }

    return { hasOverlap: false }
  }

  return {
    snapTime,
    getSnapPoints,
    checkTrackOverlap,
  }
}
