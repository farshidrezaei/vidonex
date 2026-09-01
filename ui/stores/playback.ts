import { defineStore } from 'pinia'

export const usePlaybackStore = defineStore('playback', () => {
  const currentTime = ref(0)
  const duration = ref(30)
  const isPlaying = ref(false)
  const playbackRate = ref(1.0)
  const volume = ref(1.0)
  const isMuted = ref(false)
  const isLooping = ref(false)
  
  let animationFrameId: number | null = null
  let lastTimestamp = 0

  function play() {
    if (isPlaying.value) return
    if (currentTime.value >= duration.value) {
      currentTime.value = 0
    }
    isPlaying.value = true
    lastTimestamp = performance.now()
    animationFrameId = requestAnimationFrame(loop)
  }

  function pause() {
    isPlaying.value = false
    if (animationFrameId !== null) {
      cancelAnimationFrame(animationFrameId)
      animationFrameId = null
    }
  }

  function togglePlay() {
    if (isPlaying.value) {
      pause()
    } else {
      play()
    }
  }

  function seek(time: number) {
    currentTime.value = Math.max(0, Math.min(time, duration.value))
  }

  function stepFrame(frameDelta: number, fps = 30) {
    pause()
    const frameDuration = 1 / fps
    seek(currentTime.value + frameDelta * frameDuration)
  }

  function jumpToStart() {
    pause()
    currentTime.value = 0
  }

  function jumpToEnd() {
    pause()
    currentTime.value = duration.value
  }

  function loop(timestamp: number) {
    if (!isPlaying.value) return

    const deltaTime = (timestamp - lastTimestamp) / 1000
    lastTimestamp = timestamp

    currentTime.value += deltaTime * playbackRate.value

    if (currentTime.value >= duration.value) {
      if (isLooping.value) {
        currentTime.value = 0
      } else {
        currentTime.value = duration.value
        pause()
        return
      }
    }

    animationFrameId = requestAnimationFrame(loop)
  }

  function formatTimecode(seconds: number, fps = 30): string {
    const totalFrames = Math.floor(seconds * fps)
    const frames = totalFrames % fps
    const totalSecs = Math.floor(seconds)
    const secs = totalSecs % 60
    const mins = Math.floor(totalSecs / 60) % 60
    const hours = Math.floor(totalSecs / 3600)

    const pad = (n: number, z = 2) => String(n).padStart(z, '0')
    return `${pad(hours)}:${pad(mins)}:${pad(secs)}:${pad(frames)}`
  }

  return {
    currentTime,
    duration,
    isPlaying,
    playbackRate,
    volume,
    isMuted,
    isLooping,
    play,
    pause,
    togglePlay,
    seek,
    stepFrame,
    jumpToStart,
    jumpToEnd,
    formatTimecode,
  }
})
