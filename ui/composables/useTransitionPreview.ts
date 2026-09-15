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

    case 'circleclose':
      // Circle closing down on outgoing clip revealing incoming
      if (role === 'from') {
        const radius = (1.0 - progress) * 75
        return {
          clipPath: `circle(${radius.toFixed(2)}% at center)`,
          zIndexExtra: 5,
        }
      }
      return {}

    case 'zoomin':
      // Incoming clip zooms in from center (scale 0 to 1)
      if (role === 'to') {
        const scale = Math.max(0.01, progress)
        return {
          transformExtra: `scale(${scale.toFixed(3)})`,
          opacity: Math.min(1.0, progress * 1.5),
          zIndexExtra: 5,
        }
      }
      if (role === 'from') {
        const scale = 1.0 + progress * 0.2
        return {
          transformExtra: `scale(${scale.toFixed(3)})`,
          opacity: Math.max(0, 1.0 - progress * 1.2),
        }
      }
      return {}

    case 'zoomout':
      // Outgoing clip zooms out into distance while incoming clip takes over
      if (role === 'from') {
        const scale = Math.max(0.01, 1.0 - progress)
        return {
          transformExtra: `scale(${scale.toFixed(3)})`,
          opacity: 1.0 - progress,
          zIndexExtra: 5,
        }
      }
      if (role === 'to') {
        return {
          opacity: progress,
        }
      }
      return {}

    case 'slideleft':
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
      if (role === 'from') {
        const offset = progress * 100
        return { transformExtra: `translateX(${offset.toFixed(2)}%)` }
      }
      const incomingOffsetRight = (progress - 1.0) * 100
      return {
        transformExtra: `translateX(${incomingOffsetRight.toFixed(2)}%)`,
        zIndexExtra: 5,
      }

    case 'slideup':
      if (role === 'from') {
        const offset = -progress * 100
        return { transformExtra: `translateY(${offset.toFixed(2)}%)` }
      }
      const incomingOffsetUp = (1.0 - progress) * 100
      return {
        transformExtra: `translateY(${incomingOffsetUp.toFixed(2)}%)`,
        zIndexExtra: 5,
      }

    case 'slidedown':
      if (role === 'from') {
        const offset = progress * 100
        return { transformExtra: `translateY(${offset.toFixed(2)}%)` }
      }
      const incomingOffsetDown = (progress - 1.0) * 100
      return {
        transformExtra: `translateY(${incomingOffsetDown.toFixed(2)}%)`,
        zIndexExtra: 5,
      }

    case 'smoothleft':
      // Smooth eased slide left
      if (role === 'from') {
        const eased = Math.sin((progress * Math.PI) / 2)
        return { transformExtra: `translateX(${(-eased * 100).toFixed(2)}%)` }
      }
      const easedInLeft = 1.0 - Math.sin((progress * Math.PI) / 2)
      return {
        transformExtra: `translateX(${(easedInLeft * 100).toFixed(2)}%)`,
        zIndexExtra: 5,
      }

    case 'smoothright':
      if (role === 'from') {
        const eased = Math.sin((progress * Math.PI) / 2)
        return { transformExtra: `translateX(${(eased * 100).toFixed(2)}%)` }
      }
      const easedInRight = (1.0 - Math.sin((progress * Math.PI) / 2)) * -100
      return {
        transformExtra: `translateX(${easedInRight.toFixed(2)}%)`,
        zIndexExtra: 5,
      }

    case 'smoothup':
      if (role === 'from') {
        const eased = Math.sin((progress * Math.PI) / 2)
        return { transformExtra: `translateY(${(-eased * 100).toFixed(2)}%)` }
      }
      const easedInUp = 1.0 - Math.sin((progress * Math.PI) / 2)
      return {
        transformExtra: `translateY(${(easedInUp * 100).toFixed(2)}%)`,
        zIndexExtra: 5,
      }

    case 'smoothdown':
      if (role === 'from') {
        const eased = Math.sin((progress * Math.PI) / 2)
        return { transformExtra: `translateY(${(eased * 100).toFixed(2)}%)` }
      }
      const easedInDown = (1.0 - Math.sin((progress * Math.PI) / 2)) * -100
      return {
        transformExtra: `translateY(${easedInDown.toFixed(2)}%)`,
        zIndexExtra: 5,
      }

    case 'rectcrop':
      // Expanding rectangle from center
      if (role === 'to') {
        const insetX = (1.0 - progress) * 50
        const insetY = (1.0 - progress) * 50
        return {
          clipPath: `inset(${insetY.toFixed(2)}% ${insetX.toFixed(2)}% ${insetY.toFixed(2)}% ${insetX.toFixed(2)}%)`,
          zIndexExtra: 5,
        }
      }
      return {}

    case 'horzopen':
      // Horizontal curtain opening from center outwards revealing incoming clip
      if (role === 'to') {
        const insetX = (1.0 - progress) * 50
        return {
          clipPath: `inset(0 ${insetX.toFixed(2)}% 0 ${insetX.toFixed(2)}%)`,
          zIndexExtra: 5,
        }
      }
      return {}

    case 'horzclose':
      // Horizontal curtain closing from edges inwards
      if (role === 'from') {
        const insetX = progress * 50
        return {
          clipPath: `inset(0 ${insetX.toFixed(2)}% 0 ${insetX.toFixed(2)}%)`,
          zIndexExtra: 5,
        }
      }
      return {}

    case 'vertopen':
      // Vertical curtain opening from center
      if (role === 'to') {
        const insetY = (1.0 - progress) * 50
        return {
          clipPath: `inset(${insetY.toFixed(2)}% 0 ${insetY.toFixed(2)}% 0)`,
          zIndexExtra: 5,
        }
      }
      return {}

    case 'vertclose':
      // Vertical curtain closing towards center
      if (role === 'from') {
        const insetY = progress * 50
        return {
          clipPath: `inset(${insetY.toFixed(2)}% 0 ${insetY.toFixed(2)}% 0)`,
          zIndexExtra: 5,
        }
      }
      return {}

    case 'wipetl':
      // Wipes towards top-left: outgoing clip shrinks to top-left corner
      if (role === 'from') {
        const p = progress * 100
        return {
          clipPath: `polygon(0 0, ${100 - p}% 0, 0 ${100 - p}%)`,
          zIndexExtra: 5,
        }
      }
      return {}

    case 'wipetr':
      // Wipes towards top-right: outgoing clip shrinks to top-right corner
      if (role === 'from') {
        const p = progress * 100
        return {
          clipPath: `polygon(100% 0, 100% ${100 - p}%, ${p}% 0)`,
          zIndexExtra: 5,
        }
      }
      return {}

    case 'wipebl':
      // Wipes towards bottom-left: outgoing clip shrinks to bottom-left corner
      if (role === 'from') {
        const p = progress * 100
        return {
          clipPath: `polygon(0 100%, 0 ${p}%, ${100 - p}% 100%)`,
          zIndexExtra: 5,
        }
      }
      return {}

    case 'wipebr':
      // Wipes towards bottom-right: outgoing clip shrinks to bottom-right corner
      if (role === 'from') {
        const p = progress * 100
        return {
          clipPath: `polygon(100% 100%, ${p}% 100%, 100% ${p}%)`,
          zIndexExtra: 5,
        }
      }
      return {}

    case 'squeezeh':
      // Incoming clip squeezes in horizontally
      if (role === 'to') {
        return {
          transformExtra: `scaleX(${progress.toFixed(3)})`,
          zIndexExtra: 5,
        }
      }
      return {}

    case 'squeezev':
      // Incoming clip squeezes in vertically
      if (role === 'to') {
        return {
          transformExtra: `scaleY(${progress.toFixed(3)})`,
          zIndexExtra: 5,
        }
      }
      return {}

    case 'fadeblack':
      if (role === 'from') {
        const fadeOut = progress < 0.5 ? 1.0 - progress * 2 : 0
        return { opacity: Math.max(0, fadeOut) }
      }
      const fadeInBlack = progress >= 0.5 ? (progress - 0.5) * 2 : 0
      return { opacity: Math.min(1.0, fadeInBlack), zIndexExtra: 5 }

    case 'fadewhite':
      if (role === 'from') {
        const fadeOut = progress < 0.5 ? 1.0 - progress * 2 : 0
        return { opacity: Math.max(0, fadeOut) }
      }
      const fadeInWhite = progress >= 0.5 ? (progress - 0.5) * 2 : 0
      return { opacity: Math.min(1.0, fadeInWhite), zIndexExtra: 5 }

    case 'fadegrays':
      if (role === 'from') {
        return { opacity: 1.0 - progress }
      }
      return { opacity: progress, zIndexExtra: 5 }

    case 'radial':
      // Clock radar sweep / circular reveal
      if (role === 'to') {
        return {
          clipPath: `circle(${(progress * 85).toFixed(2)}% at center)`,
          transformExtra: `rotate(${((1.0 - progress) * 60).toFixed(1)}deg)`,
          zIndexExtra: 5,
        }
      }
      return {}

    case 'pixelize':
      // Mosaic pixelation and cross-dissolve
      if (role === 'from') {
        return {
          opacity: Math.max(0, 1.0 - progress * 1.4),
          transformExtra: `scale(${1 + progress * 0.05})`,
        }
      }
      return {
        opacity: Math.min(1.0, progress * 1.6),
        zIndexExtra: 5,
      }

    case 'hlslice':
      // Horizontal multi-slice blinds reveal
      if (role === 'to') {
        const percentage = (1.0 - progress) * 100
        return {
          clipPath: `inset(0 0 0 ${percentage.toFixed(2)}%)`,
          zIndexExtra: 5,
        }
      }
      return {}

    case 'vuslice':
      // Vertical louvers blinds reveal
      if (role === 'to') {
        const percentage = (1.0 - progress) * 100
        return {
          clipPath: `inset(0 0 ${percentage.toFixed(2)}% 0)`,
          zIndexExtra: 5,
        }
      }
      return {}

    case 'distance':
      // 3D perspective distance fly-in
      if (role === 'to') {
        const scale = Math.max(0.01, progress)
        return {
          transformExtra: `scale(${scale.toFixed(3)})`,
          opacity: Math.min(1.0, progress * 1.5),
          zIndexExtra: 5,
        }
      }
      if (role === 'from') {
        return {
          opacity: Math.max(0, 1.0 - progress * 1.2),
        }
      }
      return {}

    default:
      if (role === 'from') {
        return { opacity: 1.0 - progress }
      }
      return { opacity: progress, zIndexExtra: 5 }
  }
}
