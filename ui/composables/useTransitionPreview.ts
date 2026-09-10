import type { ClipSpec, TrackSpec, TransitionSpec } from '~/types/spec'

export interface TransitionActiveInfo {
  isActive: boolean
  role: 'from' | 'to' | 'none'
  progress: number // 0.0 to 1.0 (0.5 is exactly at the cut line)
  transition?: TransitionSpec
  partnerClip?: ClipSpec
  transitionStart?: number
  transitionEnd?: number
  cutTime?: number
}

/**
 * Computes transition state centered on the cut boundary between fromClip and toClip.
 * Total transition window: [cutTime - duration/2, cutTime + duration/2]
 */
export function computeClipTransitionState(
  clip: ClipSpec,
  track: TrackSpec,
  currentTime: number
): TransitionActiveInfo {
  if (!track.transitions || track.transitions.length === 0 || !track.clips) {
    return { isActive: false, role: 'none', progress: 0 }
  }

  for (const transition of track.transitions) {
    const fromClip = track.clips.find((c) => c.id === transition.from)
    const toClip = track.clips.find((c) => c.id === transition.to)
    if (!fromClip || !toClip) continue

    const transitionDuration = Math.max(0.01, Number(transition.duration) || 1.0)
    const halfDuration = transitionDuration / 2
    const cutTime = (Number(fromClip.start) || 0) + (Number(fromClip.duration) || 0)
    const transitionStart = cutTime - halfDuration
    const transitionEnd = cutTime + halfDuration

    if (currentTime >= transitionStart && currentTime <= transitionEnd) {
      const rawProgress = (currentTime - transitionStart) / transitionDuration
      const progress = Math.min(1.0, Math.max(0.0, rawProgress))

      if (clip.id === fromClip.id) {
        return {
          isActive: true,
          role: 'from',
          progress,
          transition,
          partnerClip: toClip,
          transitionStart,
          transitionEnd,
          cutTime,
        }
      }

      if (clip.id === toClip.id) {
        return {
          isActive: true,
          role: 'to',
          progress,
          transition,
          partnerClip: fromClip,
          transitionStart,
          transitionEnd,
          cutTime,
        }
      }
    }
  }

  return { isActive: false, role: 'none', progress: 0 }
}

export function getTransitionStyleModifiers(transitionInfo: TransitionActiveInfo): {
  opacity?: number
  clipPath?: string
  transformExtra?: string
  zIndexExtra?: number
} {
  if (!transitionInfo.isActive || !transitionInfo.transition) {
    return {}
  }

  const { role, progress, transition } = transitionInfo
  const type = (transition.type || 'dissolve').toLowerCase()

  switch (type) {
    case 'dissolve':
      if (role === 'from') {
        return { opacity: 1.0 - progress }
      }
      return { opacity: progress, zIndexExtra: 5 }

    case 'fade':
      // Dip through black / transparency centered at cut (progress = 0.5)
      if (role === 'from') {
        const fadeOut = progress < 0.5 ? 1.0 - progress * 2 : 0
        return { opacity: fadeOut }
      }
      const fadeIn = progress >= 0.5 ? (progress - 0.5) * 2 : 0
      return { opacity: fadeIn, zIndexExtra: 5 }

    case 'wipeleft':
      // Incoming clip wipes in from right to left
      if (role === 'to') {
        const percentage = (1.0 - progress) * 100
        return {
          clipPath: `inset(0 0 0 ${percentage.toFixed(2)}%)`,
          zIndexExtra: 5,
        }
      }
      return {}

    case 'wiperight':
      // Incoming clip wipes in from left to right
      if (role === 'to') {
        const percentage = (1.0 - progress) * 100
        return {
          clipPath: `inset(0 ${percentage.toFixed(2)}% 0 0)`,
          zIndexExtra: 5,
        }
      }
      return {}

    case 'wipeup':
      // Incoming clip wipes in from bottom to top
      if (role === 'to') {
        const percentage = (1.0 - progress) * 100
        return {
          clipPath: `inset(${percentage.toFixed(2)}% 0 0 0)`,
          zIndexExtra: 5,
        }
      }
      return {}

    case 'wipedown':
      // Incoming clip wipes in from top to bottom
      if (role === 'to') {
        const percentage = (1.0 - progress) * 100
        return {
          clipPath: `inset(0 0 ${percentage.toFixed(2)}% 0)`,
          zIndexExtra: 5,
        }
      }
      return {}

    case 'circleopen':
    case 'circlecrop':
      // Expanding circle revealing incoming clip
      if (role === 'to') {
        const radius = progress * 75
        return {
          clipPath: `circle(${radius.toFixed(2)}% at center)`,
          zIndexExtra: 5,
        }
      }
      return {}

    case 'slideleft':
      // Outgoing slides left, incoming slides in from right
      if (role === 'from') {
        const offset = -progress * 100
        return { transformExtra: `translateX(${offset.toFixed(2)}%)` }
      }
      const incomingOffsetLeft = (1.0 - progress) * 100
      return {
        transformExtra: `translateX(${incomingOffsetLeft.toFixed(2)}%)`,
        zIndexExtra: 5,
      }

    case 'slideright':
      // Outgoing slides right, incoming slides in from left
      if (role === 'from') {
        const offset = progress * 100
        return { transformExtra: `translateX(${offset.toFixed(2)}%)` }
      }
      const incomingOffsetRight = (progress - 1.0) * 100
      return {
        transformExtra: `translateX(${incomingOffsetRight.toFixed(2)}%)`,
        zIndexExtra: 5,
      }

    default:
      if (role === 'from') {
        return { opacity: 1.0 - progress }
      }
      return { opacity: progress, zIndexExtra: 5 }
  }
}
