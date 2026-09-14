<template>
  <UModal
    v-model:open="projectStore.isExportOpen"
    :title="$t('render.title')"
    :ui="{ content: renderStore.completedDownloadUrl ? 'sm:max-w-2xl bg-gray-950 border border-gray-800' : 'sm:max-w-xl bg-gray-950 border border-gray-800' }"
  >
    <template #body>
      <!-- Configuration Form (When not rendering) -->
      <div v-if="!renderStore.isRendering && !renderStore.completedDownloadUrl" class="space-y-4 text-xs">
        <!-- Output Container -->
        <div>
          <label class="block text-gray-400 mb-1">{{ $t('render.output_format') }}</label>
          <USelect
            v-model="renderStore.selectedFormat"
            :items="formatOptions"
            size="sm"
          />
        </div>

        <!-- Output Destination (Merged Folder + File Name) -->
        <div>
          <div class="flex items-center justify-between mb-1">
            <label class="block text-gray-400">
              {{ $t('render.output_destination') || 'Output Destination & Filename' }}
            </label>
            <span v-if="!renderStore.customOutputPath" class="text-[10px] text-gray-500 font-mono">
              {{ $t('render.default_location') || 'Default system folder' }}
            </span>
            <button
              v-else
              type="button"
              class="text-[10px] text-indigo-400 hover:text-indigo-300 transition cursor-pointer"
              @click="renderStore.customOutputPath = ''"
            >
              {{ $t('render.reset_default') || 'Reset to default' }}
            </button>
          </div>

          <div class="flex items-center gap-2">
            <input
              v-model="renderStore.customOutputPath"
              type="text"
              :placeholder="defaultOutputPathPlaceholder"
              class="flex-1 min-w-0 bg-gray-900 border border-gray-800 rounded px-2.5 py-1.5 font-mono text-xs text-gray-200 placeholder-gray-600 focus:outline-none focus:border-indigo-500"
            />
            <UButton
              v-if="isDesktop"
              size="xs"
              color="primary"
              variant="soft"
              icon="i-heroicons-folder-open"
              @click="browseOutputPath"
            >
              {{ $t('render.browse') || 'Browse...' }}
            </UButton>
          </div>
        </div>

        <!-- GPU Acceleration -->
        <div>
          <label class="block text-gray-400 mb-1">{{ $t('render.gpu_acceleration') }}</label>
          <USelect
            v-model="renderStore.selectedGpu"
            :items="gpuOptions"
            size="sm"
          />
        </div>

        <!-- Platform Preset -->
        <div>
          <label class="block text-gray-400 mb-1">{{ $t('render.preset') }}</label>
          <USelect
            v-model="renderStore.selectedPreset"
            :items="presetOptions"
            size="sm"
          />
        </div>

        <!-- CRF & Bitrate -->
        <div class="grid grid-cols-2 gap-3 pt-2 border-t border-gray-800">
          <div>
            <div class="flex justify-between text-gray-400 mb-1">
              <span>{{ $t('render.crf') }}</span>
              <span class="font-mono text-gray-300">{{ renderStore.selectedCrf }}</span>
            </div>
            <input
              v-model.number="renderStore.selectedCrf"
              type="range"
              min="14"
              max="35"
              step="1"
              class="w-full h-1 bg-gray-800 rounded appearance-none accent-indigo-500"
            />
          </div>

          <div>
            <label class="block text-gray-400 mb-1">{{ $t('render.bitrate') }}</label>
            <input
              v-model="renderStore.selectedBitrate"
              type="text"
              class="w-full bg-gray-900 border border-gray-800 rounded px-2 py-1 font-mono text-gray-300 focus:outline-none focus:border-indigo-500"
            />
          </div>
        </div>

        <!-- Collapsible Error Alert with Copy Button -->
        <div
          v-if="renderStore.renderError"
          class="rounded-lg border border-red-800/80 bg-red-950/50 p-3 space-y-2 text-xs text-red-200 shadow-sm"
        >
          <div class="flex items-start justify-between gap-2">
            <div class="flex items-start gap-2 min-w-0">
              <UIcon name="i-heroicons-exclamation-circle" class="w-4 h-4 text-red-400 flex-shrink-0 mt-0.5" />
              <div class="min-w-0">
                <div class="font-semibold text-red-300">
                  {{ $t('render.error_title') || 'Render Execution Failed' }}
                </div>
                <div class="text-[11px] text-red-300/80 font-mono truncate mt-0.5" :title="renderStore.renderError">
                  {{ errorSummary }}
                </div>
              </div>
            </div>

            <div class="flex items-center gap-1.5 flex-shrink-0">
              <!-- Copy Button -->
              <button
                type="button"
                class="flex items-center gap-1 px-2 py-1 rounded bg-red-900/40 hover:bg-red-900/70 border border-red-700/50 text-[11px] text-red-200 transition cursor-pointer select-none"
                :title="$t('render.copy_error') || 'Copy full error log'"
                @click="copyErrorText"
              >
                <UIcon :name="isCopied ? 'i-heroicons-check' : 'i-heroicons-clipboard-document'" class="w-3.5 h-3.5 text-red-300" />
                <span>{{ isCopied ? ($t('render.copied') || 'Copied!') : ($t('render.copy_error') || 'Copy') }}</span>
              </button>

              <!-- Collapse / Expand Toggle Button -->
              <button
                type="button"
                class="flex items-center gap-1 px-2 py-1 rounded bg-red-900/40 hover:bg-red-900/70 border border-red-700/50 text-[11px] text-red-200 transition cursor-pointer select-none"
                @click="isErrorDetailsOpen = !isErrorDetailsOpen"
              >
                <span>{{ isErrorDetailsOpen ? ($t('render.hide_details') || 'Hide') : ($t('render.show_details') || 'Details') }}</span>
                <UIcon :name="isErrorDetailsOpen ? 'i-heroicons-chevron-up' : 'i-heroicons-chevron-down'" class="w-3.5 h-3.5 text-red-300" />
              </button>
            </div>
          </div>

          <!-- Expanded Full Error Log (Collapsible) -->
          <div v-if="isErrorDetailsOpen" class="pt-2 border-t border-red-800/40">
            <pre class="max-h-48 overflow-y-auto rounded bg-black/80 border border-red-900/50 p-2.5 text-[10px] font-mono text-red-300/90 whitespace-pre-wrap break-all select-all leading-relaxed">{{ renderStore.renderError }}</pre>
          </div>
        </div>

        <UButton
          size="md"
          color="primary"
          block
          icon="i-heroicons-bolt"
          class="font-semibold mt-4 bg-gradient-to-r from-indigo-500 to-purple-600 hover:from-indigo-600 hover:to-purple-700 shadow-lg shadow-indigo-500/20"
          @click="startRender"
        >
          {{ $t('render.start_render') }}
        </UButton>
      </div>

      <!-- Live Rendering Progress -->
      <div v-else-if="renderStore.isRendering" class="py-6 space-y-6 text-center">
        <div class="relative w-24 h-24 mx-auto flex items-center justify-center">
          <div class="absolute inset-0 rounded-full border-4 border-gray-800 animate-pulse"></div>
          <div class="text-2xl font-bold font-mono text-indigo-400">
            {{ renderStore.progressPercentage }}%
          </div>
        </div>

        <div class="space-y-1">
          <h4 class="text-sm font-semibold text-gray-200">{{ $t('render.rendering_status') }}</h4>
          <p class="text-xs text-gray-400 font-mono">
            {{ renderStore.currentFPS }} FPS | {{ renderStore.renderSpeed }}x Speed
          </p>
        </div>

        <!-- Progress Bar -->
        <div class="w-full bg-gray-900 rounded-full h-2 overflow-hidden border border-gray-800">
          <div
            class="bg-gradient-to-r from-indigo-500 to-purple-600 h-full transition-all duration-300"
            :style="{ width: `${renderStore.progressPercentage}%` }"
          ></div>
        </div>

        <UButton
          size="sm"
          color="error"
          variant="soft"
          icon="i-heroicons-x-circle"
          @click="renderStore.cancelRender"
        >
          {{ $t('render.cancel_render') }}
        </UButton>
      </div>

      <!-- Render Complete Success with Video Player Preview -->
      <div v-else-if="renderStore.completedDownloadUrl" class="py-3 space-y-4">
        <!-- Status Header Banner -->
        <div class="flex items-center justify-between p-3 rounded-lg bg-emerald-950/40 border border-emerald-800/50 text-emerald-300 shadow-sm">
          <div class="flex items-center gap-2.5 min-w-0">
            <div class="w-8 h-8 rounded-full bg-emerald-500/20 border border-emerald-500/30 flex items-center justify-center flex-shrink-0">
              <UIcon name="i-heroicons-check-badge" class="w-5 h-5 text-emerald-400" />
            </div>
            <div class="min-w-0">
              <h4 class="text-xs font-semibold text-emerald-200 leading-tight">
                {{ $t('render.completed') }}
              </h4>
              <p class="text-[11px] text-emerald-400/80 truncate mt-0.5">
                {{ $t('render.completed_desc') }}
              </p>
            </div>
          </div>

          <div v-if="formattedFileSize" class="px-2 py-1 rounded bg-emerald-900/50 border border-emerald-700/40 text-[10px] font-mono text-emerald-200 flex-shrink-0">
            {{ formattedFileSize }}
          </div>
        </div>

        <!-- Video Player Preview Viewport -->
        <div class="relative w-full rounded-xl overflow-hidden bg-black/95 border border-gray-800 shadow-2xl flex items-center justify-center min-h-[220px] max-h-[360px]">
          <img
            v-if="isGifFormat"
            :src="renderStore.completedDownloadUrl"
            alt="Rendered animation output"
            class="max-h-[340px] w-auto max-w-full object-contain mx-auto"
          />
          <video
            v-else
            ref="previewVideoRef"
            :src="renderStore.completedDownloadUrl"
            controls
            playsinline
            preload="metadata"
            class="max-h-[340px] w-auto max-w-full object-contain mx-auto focus:outline-none"
          ></video>
        </div>

        <!-- Modal Actions Footer -->
        <div class="flex items-center justify-between gap-2 pt-1 border-t border-gray-900">
          <UButton
            size="sm"
            color="neutral"
            variant="ghost"
            icon="i-heroicons-arrow-path"
            @click="newRender"
          >
            {{ $t('render.new_render') }}
          </UButton>

          <div class="flex items-center gap-2">
            <a
              :href="renderStore.completedDownloadUrl"
              download
              class="inline-flex items-center gap-1.5 px-3.5 py-1.5 rounded-lg bg-emerald-600 hover:bg-emerald-500 text-white font-medium text-xs shadow-lg shadow-emerald-600/20 transition cursor-pointer select-none"
            >
              <UIcon name="i-heroicons-arrow-down-tray" class="w-4 h-4" />
              <span>{{ $t('render.download_file') }}</span>
            </a>

            <UButton
              size="sm"
              color="neutral"
              variant="soft"
              @click="resetModal"
            >
              {{ $t('render.done') }}
            </UButton>
          </div>
        </div>
      </div>
    </template>
  </UModal>
</template>

<script setup lang="ts">
import { useProjectStore } from '~/stores/project'
import { useRenderStore } from '~/stores/render'
import { useTimelineStore } from '~/stores/timeline'
import { useDesktop } from '~/composables/useDesktop'

const { t } = useI18n()
const projectStore = useProjectStore()
const renderStore = useRenderStore()
const timelineStore = useTimelineStore()
const { isDesktop, saveVideoFileDialog } = useDesktop()

const defaultOutputPathPlaceholder = computed(() => {
  const ext = renderStore.selectedFormat || 'mp4'
  const projName = projectStore.currentProject?.name?.toLowerCase().replace(/[^a-z0-9]+/g, '_') || 'video'
  return `${projName}.${ext}`
})

async function browseOutputPath() {
  const projName = projectStore.currentProject?.name?.toLowerCase().replace(/[^a-z0-9]+/g, '_') || 'composition'
  const ext = renderStore.selectedFormat || 'mp4'
  const defaultName = `${projName}.${ext}`
  const chosenPath = await saveVideoFileDialog(defaultName, ext)
  if (chosenPath) {
    renderStore.customOutputPath = chosenPath
  }
}

const formatOptions = computed(() => [
  { label: t('render.formats.mp4'), value: 'mp4' },
  { label: t('render.formats.mkv'), value: 'mkv' },
  { label: t('render.formats.mov'), value: 'mov' },
  { label: t('render.formats.webm'), value: 'webm' },
  { label: t('render.formats.gif'), value: 'gif' },
])

const gpuOptions = computed(() => [
  { label: t('render.gpus.none'), value: 'none' },
  { label: t('render.gpus.nvenc'), value: 'nvenc' },
  { label: t('render.gpus.videotoolbox'), value: 'videotoolbox' },
  { label: t('render.gpus.qsv'), value: 'qsv' },
  { label: t('render.gpus.vaapi'), value: 'vaapi' },
])

const presetOptions = computed(() => [
  { label: t('render.presets.youtube_1080p'), value: 'youtube_1080p' },
  { label: t('render.presets.youtube_4k'), value: 'youtube_4k' },
  { label: t('render.presets.tiktok'), value: 'tiktok' },
])

const isErrorDetailsOpen = ref(false)
const isCopied = ref(false)

const errorSummary = computed(() => {
  if (!renderStore.renderError) return ''
  const trimmed = renderStore.renderError.trim()
  const firstLine = trimmed.split('\n')[0] || ''
  return firstLine.length > 110 ? firstLine.slice(0, 107) + '...' : firstLine
})

const previewVideoRef = ref<HTMLVideoElement | null>(null)

const isGifFormat = computed(() => {
  if (renderStore.selectedFormat === 'gif') return true
  if (renderStore.completedDownloadUrl?.toLowerCase().endsWith('.gif')) return true
  return false
})

const formattedFileSize = computed(() => {
  if (!renderStore.totalSize || renderStore.totalSize <= 0) return null
  const bytes = renderStore.totalSize
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
  return `${(bytes / (1024 * 1024)).toFixed(2)} MB`
})

function pausePreview() {
  if (previewVideoRef.value) {
    try {
      previewVideoRef.value.pause()
    } catch {
      // Ignore pause failure
    }
  }
}

watch(() => projectStore.isExportOpen, (isOpen) => {
  if (!isOpen) {
    pausePreview()
  }
})

async function copyErrorText() {
  if (!renderStore.renderError) return
  try {
    await navigator.clipboard.writeText(renderStore.renderError)
    isCopied.value = true
    setTimeout(() => {
      isCopied.value = false
    }, 2000)
  } catch (err) {
    console.error('Failed copying error log to clipboard:', err)
  }
}

async function startRender() {
  isErrorDetailsOpen.value = false
  isCopied.value = false
  pausePreview()
  const spec = timelineStore.toVideoSpec()
  await renderStore.startRender(spec)
}

function resetModal() {
  pausePreview()
  renderStore.completedDownloadUrl = null
  renderStore.renderError = null
  isErrorDetailsOpen.value = false
  isCopied.value = false
  projectStore.isExportOpen = false
}

function newRender() {
  pausePreview()
  renderStore.completedDownloadUrl = null
  renderStore.renderError = null
  isErrorDetailsOpen.value = false
  isCopied.value = false
}
</script>
