import { defineStore } from 'pinia'
import type { MediaAsset } from '~/types/project'

export const useMediaStore = defineStore('media', () => {
  const assets = ref<MediaAsset[]>([])
  const isUploading = ref(false)
  const uploadProgress = ref(0)
  const activeProbeAsset = ref<MediaAsset | null>(null)
  const isProbeModalOpen = ref(false)
  
  // Currently dragged asset for live timeline previews
  const draggedAsset = ref<MediaAsset | null>(null)

  async function fetchAssets(projectId: string) {
    if (!projectId) return
    try {
      const response = await $fetch<{ success: boolean; data: MediaAsset[] }>(`/api/projects/${projectId}/assets`)
      if (response.success && response.data) {
        assets.value = response.data
      }
    } catch (err) {
      console.error('Failed fetching project assets', err)
    }
  }

  async function uploadFile(projectId: string, file: File): Promise<MediaAsset | null> {
    if (!projectId || !file) return null

    isUploading.value = true
    uploadProgress.value = 0

    const formData = new FormData()
    formData.append('project_id', projectId)
    formData.append('file', file)

    try {
      const response = await $fetch<{ success: boolean; data: MediaAsset }>('/api/media/upload', {
        method: 'POST',
        body: formData,
      })

      if (response.success && response.data) {
        assets.value.unshift(response.data)
        return response.data
      }
    } catch (err) {
      console.error('Failed uploading media asset', err)
    } finally {
      isUploading.value = false
      uploadProgress.value = 0
    }
    return null
  }

  async function deleteAsset(assetId: string) {
    try {
      await $fetch(`/api/media/${assetId}`, { method: 'DELETE' })
      assets.value = assets.value.filter((asset) => asset.id !== assetId)
    } catch (err) {
      console.error('Failed deleting asset', err)
    }
  }

  function openProbeModal(asset: MediaAsset) {
    activeProbeAsset.value = asset
    isProbeModalOpen.value = true
  }

  function setDraggedAsset(asset: MediaAsset | null) {
    draggedAsset.value = asset
  }

  return {
    assets,
    isUploading,
    uploadProgress,
    activeProbeAsset,
    isProbeModalOpen,
    draggedAsset,
    fetchAssets,
    uploadFile,
    deleteAsset,
    openProbeModal,
    setDraggedAsset,
  }
})
