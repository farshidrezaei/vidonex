import { ref } from 'vue'

export interface SystemDiagnostics {
  os: string
  architecture: string
  cpu_count: number
  available_accelerators: string[]
  ffmpeg_version: string
  ffprobe_version: string
  go_version: string
}

export interface VersionCheckResponse {
  current_version: string
  latest_version: string
  has_update: boolean
  release_name?: string
  release_notes?: string
  release_url?: string
  published_at?: string
  diagnostics: SystemDiagnostics
  checked_at: string
}

const isChecking = ref(false)
const hasUpdate = ref(false)
const currentVersion = ref('1.6.0')
const latestVersion = ref('1.6.0')
const releaseName = ref('')
const releaseNotes = ref('')
const releaseUrl = ref('')
const publishedAt = ref('')
const systemDiagnostics = ref<SystemDiagnostics | null>(null)
const lastCheckedAt = ref<Date | null>(null)
const checkError = ref<string | null>(null)

export const isAboutModalOpen = ref(false)
export const isUpdateBannerDismissed = ref(false)

export function useAppUpdate() {
  async function checkForUpdates(force = false) {
    if (isChecking.value) return
    isChecking.value = true
    checkError.value = null

    try {
      const query = force ? '?force=true' : ''
      const res = await $fetch<VersionCheckResponse>(`/api/version/check${query}`)

      if (res) {
        currentVersion.value = res.current_version || '1.6.0'
        latestVersion.value = res.latest_version || currentVersion.value
        hasUpdate.value = !!res.has_update
        releaseName.value = res.release_name || ''
        releaseNotes.value = res.release_notes || ''
        releaseUrl.value = res.release_url || 'https://github.com/farshidrezaei/vidonex/releases'
        publishedAt.value = res.published_at || ''
        if (res.diagnostics) {
          systemDiagnostics.value = res.diagnostics
        }
        lastCheckedAt.value = new Date()
      }
    } catch (err: any) {
      console.warn('[UpdateChecker] Could not fetch version from backend, querying fallback', err)
      checkError.value = err?.message || 'Failed to check updates'

      // Client-side direct GitHub fallback if backend check fails
      try {
        const ghRelease = await $fetch<any>('https://api.github.com/repos/farshidrezaei/vidonex/releases/latest', {
          headers: { Accept: 'application/vnd.github.v3+json' },
        })
        if (ghRelease?.tag_name) {
          const cleanTag = ghRelease.tag_name.replace(/^v/, '')
          latestVersion.value = cleanTag
          hasUpdate.value = isNewerVersion(cleanTag, currentVersion.value)
          releaseName.value = ghRelease.name || ''
          releaseNotes.value = ghRelease.body || ''
          releaseUrl.value = ghRelease.html_url || 'https://github.com/farshidrezaei/vidonex/releases'
          publishedAt.value = ghRelease.published_at || ''
          checkError.value = null
          lastCheckedAt.value = new Date()
        }
      } catch (ghErr: any) {
        console.warn('[UpdateChecker] Fallback check also failed', ghErr)
      }
    } finally {
      isChecking.value = false
    }
  }

  function isNewerVersion(remote: string, current: string): boolean {
    const p1 = remote.split('.').map((x) => parseInt(x, 10) || 0)
    const p2 = current.split('.').map((x) => parseInt(x, 10) || 0)
    const len = Math.max(p1.length, p2.length)
    for (let i = 0; i < len; i++) {
      const v1 = p1[i] || 0
      const v2 = p2[i] || 0
      if (v1 > v2) return true
      if (v1 < v2) return false
    }
    return false
  }

  function openReleaseUrl() {
    if (typeof window !== 'undefined' && releaseUrl.value) {
      window.open(releaseUrl.value, '_blank', 'noopener,noreferrer')
    }
  }

  return {
    isChecking,
    hasUpdate,
    currentVersion,
    latestVersion,
    releaseName,
    releaseNotes,
    releaseUrl,
    publishedAt,
    systemDiagnostics,
    lastCheckedAt,
    checkError,
    isAboutModalOpen,
    isUpdateBannerDismissed,
    checkForUpdates,
    openReleaseUrl,
  }
}
