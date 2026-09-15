import { usePlaybackStore } from '~/stores/playback'
import { useTimelineStore } from '~/stores/timeline'
import { useDesktop } from '~/composables/useDesktop'
import type { ClipSpec, TrackSpec } from '~/types/spec'

const IMAGE_EXTENSIONS = new Set(['png', 'jpg', 'jpeg', 'webp', 'svg', 'gif', 'bmp', 'ico', 'tiff', 'tif'])

function isAudioPlayableSource(source?: string): boolean {
  if (!source) return false
  const cleanSource = source.split('?')[0].split('#')[0]
  const ext = cleanSource.toLowerCase().split('.').pop() || ''
  return !IMAGE_EXTENSIONS.has(ext)
}

export function useTimelineAudio() {
  const playbackStore = usePlaybackStore()
  const timelineStore = useTimelineStore()
  const { resolveMediaUrl } = useDesktop()

  // Cache of clipId -> HTMLAudioElement
  const audioElementPool = new Map<string, HTMLAudioElement>()

  // Web Audio Context for unlocking autoplay policy on first user interaction
  let audioContext: AudioContext | null = null
  let isUnlocked = false

  function unlockAudio() {
    if (isUnlocked && audioContext && audioContext.state === 'running') return
    try {
      const AudioCtx = window.AudioContext || (window as any).webkitAudioContext
      if (AudioCtx) {
        if (!audioContext) {
          audioContext = new AudioCtx()
        }
        if (audioContext.state === 'suspended') {
          audioContext.resume().catch(() => {})
        }
      }
      // Play and immediately pause any dormant elements to acquire user-gesture blessing
      for (const audio of audioElementPool.values()) {
        const promise = audio.play()
        if (promise !== undefined) {
          promise.then(() => {
            if (!playbackStore.isPlaying) {
              audio.pause()
            }
          }).catch(() => {})
        }
      }
      isUnlocked = true
    } catch {
      // Audio context unlock failure ignored
    }
  }

  // Bind global user gesture unlockers
  if (typeof window !== 'undefined') {
    const handleGesture = () => {
      unlockAudio()
      window.removeEventListener('click', handleGesture)
      window.removeEventListener('keydown', handleGesture)
      window.removeEventListener('touchstart', handleGesture)
    }
    window.addEventListener('click', handleGesture, { passive: true, once: true })
    window.addEventListener('keydown', handleGesture, { passive: true, once: true })
    window.addEventListener('touchstart', handleGesture, { passive: true, once: true })
  }

  function getOrCreateAudio(clip: ClipSpec): HTMLAudioElement | null {
    if (!isAudioPlayableSource(clip.source)) {
      return null
    }

    const expectedSrc = resolveMediaUrl(clip.source)
    if (!expectedSrc) return null

    let audio = audioElementPool.get(clip.id)

    if (!audio) {
      audio = new Audio()
      audio.preload = 'auto'
      audio.src = expectedSrc
      audio.addEventListener('loadedmetadata', () => {
        if (playbackStore.isPlaying) {
          syncAudio()
        }
      })
      audioElementPool.set(clip.id, audio)
    } else {
      if (audio.src !== expectedSrc && !audio.src.endsWith(expectedSrc)) {
        try {
          const curUrl = new URL(audio.src, window.location.origin)
          const expUrl = new URL(expectedSrc, window.location.origin)
          if (curUrl.pathname !== expUrl.pathname) {
            audio.src = expectedSrc
          }
        } catch {
          audio.src = expectedSrc
        }
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

    // Retain audio tracks, waveform tracks, and unmuted video/overlay tracks
    for (const track of timelineStore.tracks) {
      if (track.kind !== 'audio' && track.kind !== 'video' && track.kind !== 'overlay' && track.kind !== 'waveform') continue
      if (track.muted) continue

      const trackVolume = Number.isFinite(track.volume) ? Number(track.volume) : 1.0

      for (const clip of track.clips || []) {
        if (!clip.source || !isAudioPlayableSource(clip.source)) continue

        const start = parseFloat(String(clip.start)) || 0
        const duration = parseFloat(String(clip.duration)) || 0
        const trim = parseFloat(String(clip.trim || 0)) || 0
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
          const clipRawVol = Number.isFinite(clip.volume) ? Number(clip.volume) : 1.0
          const masterRawVol = masterMuted ? 0 : (Number.isFinite(masterVolume) ? Number(masterVolume) : 1.0)
          let clipVol = clipRawVol * trackVolume * masterRawVol

          const fadeIn = parseFloat(String(clip.fade_in || 0)) || 0
          if (fadeIn > 0 && elapsedInClip < fadeIn) {
            clipVol *= (elapsedInClip / fadeIn)
          }

          const fadeOut = parseFloat(String(clip.fade_out || 0)) || 0
          const remainingTime = duration - elapsedInClip
          if (fadeOut > 0 && remainingTime < fadeOut) {
            clipVol *= (remainingTime / fadeOut)
          }

          audio.volume = Math.max(0, Math.min(1, isNaN(clipVol) ? 1 : clipVol))
          const effectiveRate = Math.max(0.25, Math.min(4.0, speed * (playbackStore.playbackRate || 1.0)))
          if (Number.isFinite(effectiveRate)) {
            audio.playbackRate = effectiveRate
          }

          if (isPlaying) {
            // Only seek if audio is paused or drifting by more than 0.35s to prevent continuous stutter
            const drift = Math.abs(audio.currentTime - targetAudioTime)
            if (audio.paused || drift > 0.35) {
              try {
                audio.currentTime = targetAudioTime
              } catch {
                // Ignore seek errors if metadata is not ready yet
              }
            }
            if (audio.paused) {
              const playPromise = audio.play()
              if (playPromise !== undefined) {
                playPromise.catch((err) => {
                  console.warn('[useTimelineAudio] Playback interrupted or blocked:', err)
                })
              }
            }
          } else {
            if (!audio.paused) {
              audio.pause()
            }
            if (Math.abs(audio.currentTime - targetAudioTime) > 0.05) {
              try {
                audio.currentTime = targetAudioTime
              } catch {
                // Ignore
              }
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

  // Pre-load audio when composable mounts
  onMounted(() => {
    syncAudio()
  })

  // Watch playback state
  watch(() => playbackStore.isPlaying, (playing) => {
    if (!playing) {
      pauseAll()
    } else {
      unlockAudio()
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
    if (audioContext && audioContext.state !== 'closed') {
      audioContext.close().catch(() => {})
    }
  })

  return {
    syncAudio,
    pauseAll,
  }
}
