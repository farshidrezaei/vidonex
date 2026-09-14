import { usePlaybackStore } from '~/stores/playback'
import { useTimelineStore } from '~/stores/timeline'
import { useDesktop } from '~/composables/useDesktop'
import type { ClipSpec, TrackSpec } from '~/types/spec'

const IMAGE_EXTENSIONS = new Set(['png', 'jpg', 'jpeg', 'webp', 'svg', 'gif', 'bmp', 'ico'])

function isAudioPlayableSource(source?: string): boolean {
  if (!source) return false
  const ext = source.toLowerCase().split('.').pop() || ''
  return !IMAGE_EXTENSIONS.has(ext)
}

export function useTimelineAudio() {
  const playbackStore = usePlaybackStore()
  const timelineStore = useTimelineStore()
  const { resolveMediaUrl } = useDesktop()

  // Cache of clipId -> HTMLAudioElement
  const audioElementPool = new Map<string, HTMLAudioElement>()

  function getOrCreateAudio(clip: ClipSpec): HTMLAudioElement | null {
    if (!isAudioPlayableSource(clip.source)) {
      return null
    }

    const expectedSrc = resolveMediaUrl(clip.source)
    let audio = audioElementPool.get(clip.id)

    if (!audio) {
      audio = new Audio(expectedSrc)
      audio.preload = 'auto'
      audioElementPool.set(clip.id, audio)
    } else {
      // Check if source changed
      if (audio.src !== expectedSrc) {
        audio.src = expectedSrc
      }
    }

    return audio
  }

  function pauseAll() {
    for (const audio of audioElementPool.values()) {
      if (!audio.paused) {
        try {
          audio.pause()
        } catch {
          // Ignore
        }
      }
    }
  }

  function syncAudio() {
    if (typeof window === 'undefined') return

    const currentTime = playbackStore.currentTime
    const isPlaying = playbackStore.isPlaying
    const masterMuted = playbackStore.isMuted
    const masterVolume = playbackStore.volume

    // Collect all clips that can emit audio
    const activeClipIds = new Set<string>()

    // Retain only audio tracks and unmuted video tracks
    for (const track of timelineStore.tracks) {
      if (track.kind !== 'audio' && track.kind !== 'video') continue
      if (track.muted) continue

      const trackVolume = track.volume ?? 1.0

      for (const clip of track.clips || []) {
        if (!clip.source || !isAudioPlayableSource(clip.source)) continue

        const start = Number(clip.start) || 0
        const duration = Number(clip.duration) || 0
        const trim = Number(clip.trim) || 0
        const speed = Math.max(0.1, Number(clip.speed) || 1.0)
        const end = start + duration

        const isActive = currentTime >= start && currentTime < end

        if (isActive) {
          activeClipIds.add(clip.id)
          const audio = getOrCreateAudio(clip)
          if (!audio) continue

          const elapsedInClip = currentTime - start
          const targetAudioTime = Math.max(0, trim + (elapsedInClip * speed))

          // Calculate effective volume with track, clip, master, and fade in/out
          let clipVol = (clip.volume ?? 1.0) * trackVolume * (masterMuted ? 0 : masterVolume)

          const fadeIn = Number(clip.fade_in) || 0
          if (fadeIn > 0 && elapsedInClip < fadeIn) {
            clipVol *= (elapsedInClip / fadeIn)
          }

          const fadeOut = Number(clip.fade_out) || 0
          const remainingTime = duration - elapsedInClip
          if (fadeOut > 0 && remainingTime < fadeOut) {
            clipVol *= (remainingTime / fadeOut)
          }

          audio.volume = Math.max(0, Math.min(1, clipVol))
          audio.playbackRate = Math.max(0.25, Math.min(4.0, speed * playbackStore.playbackRate))

          if (isPlaying) {
            // Check drift against target time
            const drift = Math.abs(audio.currentTime - targetAudioTime)
            if (audio.paused || drift > 0.15) {
              audio.currentTime = targetAudioTime
            }
            if (audio.paused) {
              audio.play().catch(() => {
                // Autoplay policy or media abort ignored
              })
            }
          } else {
            if (!audio.paused) {
              audio.pause()
            }
            // Sync playhead scrubbing while paused
            if (Math.abs(audio.currentTime - targetAudioTime) > 0.05) {
              audio.currentTime = targetAudioTime
            }
          }
        }
      }
    }

    // Pause any pooled audio elements that are not currently active
    for (const [clipId, audio] of audioElementPool.entries()) {
      if (!activeClipIds.has(clipId) && !audio.paused) {
        try {
          audio.pause()
        } catch {
          // Ignore
        }
      }
    }
  }

  // Watch playback state
  watch(() => playbackStore.isPlaying, (playing) => {
    if (!playing) {
      pauseAll()
    } else {
      syncAudio()
    }
  })

  // Watch current playhead time
  watch(() => playbackStore.currentTime, () => {
    syncAudio()
  })

  // Watch volume / mute changes
  watch([() => playbackStore.volume, () => playbackStore.isMuted], () => {
    syncAudio()
  })

  // Prune unused clip audio instances when tracks change
  watch(() => timelineStore.tracks, (tracks) => {
    const existingClipIds = new Set<string>()
    for (const t of tracks) {
      for (const c of t.clips || []) {
        existingClipIds.add(c.id)
      }
    }
    for (const [id, audio] of audioElementPool.entries()) {
      if (!existingClipIds.has(id)) {
        try {
          audio.pause()
          audio.src = ''
        } catch {
          // Ignore
        }
        audioElementPool.delete(id)
      }
    }
  }, { deep: true })

  onUnmounted(() => {
    pauseAll()
    for (const audio of audioElementPool.values()) {
      try {
        audio.src = ''
      } catch {
        // Ignore
      }
    }
    audioElementPool.clear()
  })

  return {
    syncAudio,
    pauseAll,
  }
}
