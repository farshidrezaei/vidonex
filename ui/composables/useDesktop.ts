// Composable providing seamless interaction with the native Wails v2 desktop environment.

export interface DesktopServerInfo {
  port: number
  url: string
  ws: string
}

export interface SystemCapabilities {
  os: string
  architecture: string
  cpu_count: number
  available_accelerators: string[]
  ffmpeg_version: string
  ffprobe_version: string
  server_address: string
  server_port: number
}

export interface ImportedAssetResult {
  asset?: any
  error?: string
}

export function useDesktop() {
  const isDesktop = computed(() => {
    if (typeof window === 'undefined') return false
    return !!((window as any).go?.desktop?.App || (window as any).runtime)
  })

  async function getServerInfo(): Promise<DesktopServerInfo | null> {
    if (!isDesktop.value) return null
    try {
      return await (window as any).go.desktop.App.GetServerInfo()
    } catch (err) {
      console.warn('[Desktop] Failed retrieving server info', err)
      return null
    }
  }

  async function getSystemCapabilities(): Promise<SystemCapabilities | null> {
    if (!isDesktop.value) return null
    try {
      return await (window as any).go.desktop.App.GetSystemCapabilities()
    } catch (err) {
      console.warn('[Desktop] Failed retrieving system capabilities', err)
      return null
    }
  }

  async function selectMediaFiles(): Promise<string[]> {
    if (!isDesktop.value) return []
    try {
      const files = await (window as any).go.desktop.App.SelectMediaFiles()
      return files || []
    } catch (err) {
      console.error('[Desktop] Failed selecting media files', err)
      return []
    }
  }

  async function saveProjectFileDialog(defaultName = 'project.vidonyx.json'): Promise<string> {
    if (!isDesktop.value) return ''
    try {
      return await (window as any).go.desktop.App.SaveProjectFileDialog(defaultName)
    } catch (err) {
      console.error('[Desktop] Failed saving project file dialog', err)
      return ''
    }
  }

  async function openProjectFileDialog(): Promise<string> {
    if (!isDesktop.value) return ''
    try {
      return await (window as any).go.desktop.App.OpenProjectFileDialog()
    } catch (err) {
      console.error('[Desktop] Failed opening project file dialog', err)
      return ''
    }
  }

  async function selectExportDirectory(): Promise<string> {
    if (!isDesktop.value) return ''
    try {
      return await (window as any).go.desktop.App.SelectExportDirectory()
    } catch (err) {
      console.error('[Desktop] Failed selecting export directory', err)
      return ''
    }
  }

  async function importLocalMedia(projectID: string, filePaths: string[]): Promise<ImportedAssetResult[]> {
    if (!isDesktop.value || filePaths.length === 0) return []
    try {
      return await (window as any).go.desktop.App.ImportLocalMedia(projectID, filePaths)
    } catch (err) {
      console.error('[Desktop] Failed importing local media', err)
      return []
    }
  }

  return {
    isDesktop,
    getServerInfo,
    getSystemCapabilities,
    selectMediaFiles,
    saveProjectFileDialog,
    openProjectFileDialog,
    selectExportDirectory,
    importLocalMedia,
  }
}
