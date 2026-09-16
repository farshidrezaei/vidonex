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

// In-App Self-Update Reactive State
const isUpgrading = ref(false)
const upgradeProgress = ref(0)
const upgradeStage = ref<string>('')
const upgradeMessage = ref<string>('')
const upgradeError = ref<string | null>(null)
const isUpgradeComplete = ref(false)

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

  async function startInAppUpdate() {
    if (isUpgrading.value) return
    isUpgrading.value = true
    upgradeProgress.value = 5
    upgradeStage.value = 'checking'
    upgradeMessage.value = 'Initiating update...'
    upgradeError.value = null
    isUpgradeComplete.value = false

    try {
      const response = await fetch('/api/version/upgrade', {
        method: 'POST',
      })

      if (!response.ok) {
        throw new Error(`Server returned HTTP ${response.status}: ${response.statusText}`)
      }

      if (!response.body) {
        throw new Error('ReadableStream not supported on this browser')
      }

      const reader = response.body.getReader()
      const decoder = new TextDecoder('utf-8')
      let buffer = ''

      while (true) {
        const { done, value } = await reader.read()
        if (done) break

        buffer += decoder.decode(value, { stream: true })
        const lines = buffer.split('\n')
        buffer = lines.pop() || ''

        for (const line of lines) {
          const trimmed = line.trim()
          if (trimmed.startsWith('data:')) {
            try {
              const eventData = JSON.parse(trimmed.slice(5).trim())
              if (eventData.percentage !== undefined) {
                upgradeProgress.value = Math.round(eventData.percentage)
              }
              if (eventData.stage) {
                upgradeStage.value = eventData.stage
              }
              if (eventData.message) {
                upgradeMessage.value = eventData.message
              }
              if (eventData.error) {
                upgradeError.value = eventData.error
              }
              if (eventData.stage === 'completed') {
                isUpgradeComplete.value = true
                currentVersion.value = latestVersion.value
                hasUpdate.value = false
              }
              if (eventData.stage === 'failed') {
                throw new Error(eventData.error || eventData.message || 'Update failed')
              }
            } catch (jsonErr: any) {
              // Ignore malformed ping lines
            }
          }
        }
      }
    } catch (err: any) {
      console.error('[UpdateChecker] Self-update failed:', err)
      upgradeError.value = err?.message || 'Failed to update Vidonex'
      upgradeStage.value = 'failed'
    } finally {
      isUpgrading.value = false
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
    isUpgrading,
    upgradeProgress,
    upgradeStage,
    upgradeMessage,
    upgradeError,
    isUpgradeComplete,
    checkForUpdates,
    startInAppUpdate,
    openReleaseUrl,
  }
}

