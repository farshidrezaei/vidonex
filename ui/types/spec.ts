export interface VideoSpec {
  version: string
  canvas: CanvasSpec
  fps?: number | string
  output?: string
  preset?: string
  hardware_acceleration?: string
  encoding?: EncodingSpec
  tracks: TrackSpec[]
  metadata?: Record<string, string>
}

export interface CanvasSpec {
  width?: number
  height?: number
  preset?: string
  background_color?: string
}

export interface EncodingSpec {
  video_codec?: string
  audio_codec?: string
  pixel_format?: string
  audio_bitrate?: string
  crf?: number
  preset?: string
}

export type TrackKind = 'video' | 'audio' | 'overlay' | 'subtitle' | 'waveform'

export interface TrackSpec {
  id: string
  kind: TrackKind
  z_index?: number
  volume?: number
  muted?: boolean
  duck_under?: string
  ducking?: DuckingSpec
  subtitles?: SubtitlesSpec
  waveform?: WaveformSpec
  clips?: ClipSpec[]
  transitions?: TransitionSpec[]
}

export interface ClipSpec {
  id: string
  source?: string
  start: number | string // e.g. 2.5 or "2.5s"
  duration: number | string
  trim?: number | string
  speed?: number
  volume?: number
  opacity?: number
  scale?: number
  rotation?: number
  blend_mode?: string
  fade_in?: number | string
  fade_out?: number | string
  position?: PositionSpec
  crop?: CropSpec
  color_grading?: ColorGradingSpec
  chroma_key?: ChromaKeySpec
  keyframes?: KeyframesSpec
}

export interface CropSpec {
  top?: number // insets in percentage (0 to 100) or pixels
  bottom?: number
  left?: number
  right?: number
}

export interface ColorGradingSpec {
  brightness?: number // -1.0 to 1.0 (default 0.0)
  contrast?: number // 0.0 to 3.0 (default 1.0)
  saturation?: number // 0.0 to 3.0 (default 1.0)
  gamma?: number // 0.1 to 3.0 (default 1.0)
  temperature?: number // -1.0 to 1.0 (default 0.0)
  tint?: number // -1.0 to 1.0 (default 0.0)
  lut_file?: string
}

export interface PositionSpec {
  x?: number
  y?: number
  alignment?: 'center' | 'top_left' | 'top_center' | 'top_right' | 'center_left' | 'center_right' | 'bottom_left' | 'bottom_center' | 'bottom_right'
}

export interface ChromaKeySpec {
  color: string // hex e.g. "#00FF00"
  similarity?: number
  blend?: number
  despill?: boolean
}

export interface DuckingSpec {
  threshold_db?: number
  duck_db?: number
  attack_ms?: number
  release_ms?: number
}

export interface SubtitlesSpec {
  file?: string
  srt_content?: string
  vtt_content?: string
  style?: SubtitleStyleSpec
}

export interface SubtitleStyleSpec {
  font_name?: string
  font_size?: number
  color?: string
  outline_color?: string
  outline_width?: number
  box?: boolean
  box_color?: string
  alignment?: 'bottom_center' | 'top_center' | 'center'
}

export interface WaveformSpec {
  mode?: 'peak_to_peak' | 'bars' | 'circular' | 'wave'
  color?: string
  secondary_color?: string
  height?: number
  line_width?: number
}

export interface TransitionSpec {
  id?: string
  type: string // "fade", "dissolve", "wipeleft", "wiperight", "circleopen"
  duration: number | string
  from?: string
  to?: string
}

export interface KeyframesSpec {
  ken_burns?: KenBurnsSpec
}

export interface KenBurnsSpec {
  start_rect: { x: number; y: number; width: number; height: number }
  end_rect: { x: number; y: number; width: number; height: number }
  easing?: 'linear' | 'ease_in' | 'ease_out' | 'ease_in_out' | 'elastic'
}
