import { defineStore } from 'pinia'
import type { Project, AspectRatioPreset } from '~/types/project'
import type { VideoSpec } from '~/types/spec'

export const ASPECT_RATIO_PRESETS: AspectRatioPreset[] = [
  { id: 'youtube_16_9', name: '16:9 Landscape (1080p)', width: 1920, height: 1080, ratio: '16:9', icon: 'i-heroicons-tv' },
  { id: 'youtube_4k', name: '16:9 4K UHD (3840x2160)', width: 3840, height: 2160, ratio: '16:9', icon: 'i-heroicons-sparkles' },
  { id: 'tiktok_9_16', name: '9:16 Vertical (Reels/TikTok)', width: 1080, height: 1920, ratio: '9:16', icon: 'i-heroicons-device-phone-mobile' },
  { id: 'instagram_square', name: '1:1 Square (Instagram)', width: 1080, height: 1080, ratio: '1:1', icon: 'i-heroicons-squares-2x2' },
  { id: 'cinema_wide', name: '21:9 Ultra-Wide Cinema', width: 2560, height: 1080, ratio: '21:9', icon: 'i-heroicons-film' },
]

export const useProjectStore = defineStore('project', () => {
  const currentProject = ref<Project | null>(null)
  const projectsList = ref<Project[]>([])
  const isSaving = ref(false)
  const isSettingsOpen = ref(false)
  const isGraphOpen = ref(false)
  const isExportOpen = ref(false)
  const lastSavedAt = ref<Date | null>(null)

  const canvasWidth = computed(() => currentProject.value?.width || 1920)
  const canvasHeight = computed(() => currentProject.value?.height || 1080)
  const aspectRatio = computed(() => `${canvasWidth.value}:${canvasHeight.value}`)
  const frameRate = computed(() => currentProject.value?.frame_rate || 30.0)
  const backgroundColor = computed(() => currentProject.value?.background_color || '#000000')

  async function fetchProjects() {
    try {
      const response = await $fetch<{ success: boolean; data: Project[] }>('/api/projects')
      if (response.success && response.data) {
        projectsList.value = response.data
      }
    } catch (err) {
      console.error('Failed fetching projects', err)
    }
  }

  async function createProject(name = 'Untitled Project', width = 1920, height = 1080, frameRate = 30.0) {
    try {
      const response = await $fetch<{ success: boolean; data: Project }>('/api/projects', {
        method: 'POST',
        body: {
          name,
          width,
          height,
          frame_rate: frameRate,
          background_color: '#000000',
        },
      })
      if (response.success && response.data) {
        currentProject.value = response.data
        await fetchProjects()
        return response.data
      }
    } catch (err) {
      console.error('Failed creating project', err)
    }
    return null
  }

  async function loadProject(id: string) {
    try {
      const response = await $fetch<{ success: boolean; data: Project }>(`/api/projects/${id}`)
      if (response.success && response.data) {
        currentProject.value = response.data
        return response.data
      }
    } catch (err) {
      console.error('Failed loading project', err)
    }
    return null
  }

  async function saveCurrentProject(spec: VideoSpec) {
    if (!currentProject.value) return

    isSaving.value = true
    try {
      await $fetch(`/api/projects/${currentProject.value.id}`, {
        method: 'PUT',
        body: {
          name: currentProject.value.name,
          description: currentProject.value.description || '',
          width: currentProject.value.width,
          height: currentProject.value.height,
          frame_rate: currentProject.value.frame_rate,
          background_color: currentProject.value.background_color,
          specification: spec,
        },
      })
      lastSavedAt.value = new Date()
    } catch (err) {
      console.error('Failed saving project', err)
    } finally {
      isSaving.value = false
    }
  }

  function setPreset(preset: AspectRatioPreset) {
    if (!currentProject.value) return
    currentProject.value.width = preset.width
    currentProject.value.height = preset.height
    if (currentProject.value.specification) {
      if (!currentProject.value.specification.canvas) {
        currentProject.value.specification.canvas = {
          width: preset.width,
          height: preset.height,
          frame_rate: currentProject.value.frame_rate || 30.0,
          background_color: currentProject.value.background_color || '#000000',
        }
      } else {
        currentProject.value.specification.canvas.width = preset.width
        currentProject.value.specification.canvas.height = preset.height
      }
      saveCurrentProject(currentProject.value.specification)
    }
  }

  function setAspectRatio(ratio: string) {
    const found = ASPECT_RATIO_PRESETS.find((p) => p.ratio === ratio || p.id === ratio)
    if (found) {
      setPreset(found)
    }
  }

  return {
    currentProject,
    projectsList,
    isSaving,
    isSettingsOpen,
    isGraphOpen,
    isExportOpen,
    lastSavedAt,
    canvasWidth,
    canvasHeight,
    aspectRatio,
    frameRate,
    backgroundColor,
    fetchProjects,
    createProject,
    loadProject,
    saveCurrentProject,
    setPreset,
    setAspectRatio,
  }
})
