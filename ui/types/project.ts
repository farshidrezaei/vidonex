import type { VideoSpec } from './spec'

export interface Project {
  id: string
  name: string
  description?: string
  width: number
  height: number
  frame_rate: number
  background_color: string
  specification: string | VideoSpec
  created_at: string
  updated_at: string
}

export interface MediaAsset {
  id: string
  project_id: string
  file_name: string
  file_path: string
  file_type: 'video' | 'audio' | 'image'
  file_size_bytes: number
  duration_seconds: number
  width: number
  height: number
  frame_rate: number
  sample_rate: number
  channels: number
  thumbnail_path?: string
  waveform_data?: string
  created_at: string
}

export interface RenderJob {
  id: string
  project_id: string
  status: 'pending' | 'rendering' | 'completed' | 'failed' | 'cancelled'
  progress_percentage: number
  current_frame: number
  current_fps: number
  current_time_seconds: number
  render_speed: number
  output_path: string
  error_message?: string
  created_at: string
  finished_at?: string
}

export interface AspectRatioPreset {
  id: string
  name: string
  width: number
  height: number
  ratio: string
  icon: string
}
