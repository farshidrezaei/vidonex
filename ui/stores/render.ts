import { defineStore } from 'pinia'
import type { RenderJob } from '~/types/project'
import type { VideoSpec } from '~/types/spec'
import { useProjectStore } from './project'

export interface RenderProgressData {
  job_id: string
  percentage: number
  current_frame: number
  current_fps: number
  current_time: number
  render_speed: number
  bitrate?: string
  total_size?: number
}

export const useRenderStore = defineStore('render', () => {
  const projectStore = useProjectStore()

  const isRendering = ref(false)
  const activeJobId = ref<string | null>(null)
  const progressPercentage = ref(0)
  const currentFPS = ref(0)
  const currentRenderTime = ref(0)
  const renderSpeed = ref(0)
  const bitrate = ref('')
  const totalSize = ref(0)
  const completedDownloadUrl = ref<string | null>(null)
  const renderError = ref<string | null>(null)

  // Options
  const selectedFormat = ref('mp4')
  const selectedPreset = ref('youtube_1080p')
  const selectedGpu = ref('none')
  const selectedCrf = ref(23)
  const selectedBitrate = ref('6000k')

  async function startRender(spec: VideoSpec) {
    if (!projectStore.currentProject) return null

    isRendering.value = true
    progressPercentage.value = 0
    currentFPS.value = 0
    currentRenderTime.value = 0
    renderSpeed.value = 0
    completedDownloadUrl.value = null
    renderError.value = null

    try {
      const response = await $fetch<{ job_id: string; status: string; output_path: string }>('/api/render/start', {
        method: 'POST',
        body: {
          project_id: projectStore.currentProject.id,
          specification: spec,
          output_format: selectedFormat.value,
          preset: selectedPreset.value,
          hardware_acceleration: selectedGpu.value === 'none' ? undefined : selectedGpu.value,
        },
      })

      activeJobId.value = response.job_id
      return response.job_id
    } catch (err: any) {
      console.error('Failed initiating render job', err)
      isRendering.value = false
      renderError.value = err?.data?.message || err?.message || 'Failed starting render'
      return null
    }
  }

  async function cancelRender() {
    if (!activeJobId.value) return
    try {
      await $fetch(`/api/render/${activeJobId.value}/cancel`, { method: 'POST' })
      isRendering.value = false
      activeJobId.value = null
    } catch (err) {
      console.error('Failed cancelling render', err)
    }
  }

  function handleProgressEvent(data: RenderProgressData) {
    if (activeJobId.value && data.job_id !== activeJobId.value) return
    isRendering.value = true
    progressPercentage.value = Math.round(data.percentage * 10) / 10
    currentFPS.value = Math.round(data.current_fps * 10) / 10
    currentRenderTime.value = Math.round(data.current_time * 10) / 10
    renderSpeed.value = Math.round(data.render_speed * 100) / 100
    if (data.bitrate) bitrate.value = data.bitrate
    if (data.total_size) totalSize.value = data.total_size
  }

  function handleCompleteEvent(data: { job_id: string; output_path: string; download_url: string }) {
    isRendering.value = false
    progressPercentage.value = 100
    completedDownloadUrl.value = data.download_url
  }

  function handleFailedEvent(data: { job_id: string; error: string }) {
    isRendering.value = false
    renderError.value = data.error
  }

  return {
    isRendering,
    activeJobId,
    progressPercentage,
    currentFPS,
    currentRenderTime,
    renderSpeed,
    bitrate,
    totalSize,
    completedDownloadUrl,
    renderError,
    selectedFormat,
    selectedPreset,
    selectedGpu,
    selectedCrf,
    selectedBitrate,
    startRender,
    cancelRender,
    handleProgressEvent,
    handleCompleteEvent,
    handleFailedEvent,
  }
})
