<template>
  <UModal
    v-model:open="projectStore.isExportOpen"
    :title="$t('render.title')"
    :ui="{ content: 'sm:max-w-xl bg-gray-950 border border-gray-800' }"
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

        <!-- Error Alert -->
        <div v-if="renderStore.renderError" class="p-3 bg-red-950/40 border border-red-800/80 rounded-lg text-xs text-red-300">
          {{ renderStore.renderError }}
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

      <!-- Render Complete Success -->
      <div v-else-if="renderStore.completedDownloadUrl" class="py-6 space-y-4 text-center">
        <div class="w-14 h-14 rounded-full bg-emerald-500/20 border border-emerald-500/30 text-emerald-400 mx-auto flex items-center justify-center">
          <UIcon name="i-heroicons-check-badge" class="w-8 h-8" />
        </div>

        <div class="space-y-1">
          <h4 class="text-base font-semibold text-gray-100">{{ $t('render.completed') }}</h4>
          <p class="text-xs text-gray-400">Your video was encoded and verified with Vidonyx engine.</p>
        </div>

        <div class="flex items-center justify-center gap-3 pt-2">
          <a
            :href="renderStore.completedDownloadUrl"
            download
            class="inline-flex items-center gap-2 px-4 py-2 rounded-lg bg-emerald-600 hover:bg-emerald-500 text-white font-medium text-xs shadow-lg shadow-emerald-600/20 transition"
          >
            <UIcon name="i-heroicons-arrow-down-tray" class="w-4 h-4" />
            {{ $t('render.download_file') }}
          </a>

          <UButton
            size="sm"
            color="neutral"
            variant="ghost"
            @click="resetModal"
          >
            Done
          </UButton>
        </div>
      </div>
    </template>
  </UModal>
</template>

<script setup lang="ts">
import { useProjectStore } from '~/stores/project'
import { useRenderStore } from '~/stores/render'
import { useTimelineStore } from '~/stores/timeline'

const { t } = useI18n()
const projectStore = useProjectStore()
const renderStore = useRenderStore()
const timelineStore = useTimelineStore()

const formatOptions = [
  { label: 'MP4 Video (H.264 / AAC - Standard Web)', value: 'mp4' },
  { label: 'MKV Video (Matroska High Fidelity)', value: 'mkv' },
  { label: 'MOV Video (QuickTime ProRes / H.264)', value: 'mov' },
  { label: 'WebM Video (VP9 / Opus)', value: 'webm' },
  { label: 'GIF Animation (Palettegen / Loop)', value: 'gif' },
]

const gpuOptions = computed(() => [
  { label: t('render.gpus.none'), value: 'none' },
  { label: t('render.gpus.nvenc'), value: 'nvenc' },
  { label: t('render.gpus.videotoolbox'), value: 'videotoolbox' },
  { label: t('render.gpus.qsv'), value: 'qsv' },
  { label: t('render.gpus.vaapi'), value: 'vaapi' },
])

const presetOptions = [
  { label: 'YouTube Standard 1080p 60fps', value: 'youtube_1080p' },
  { label: 'YouTube 4K Ultra HD 60fps', value: 'youtube_4k' },
  { label: 'TikTok / Shorts / Reels Vertical 1080p 60fps', value: 'tiktok' },
]

async function startRender() {
  const spec = timelineStore.toVideoSpec()
  await renderStore.startRender(spec)
}

function resetModal() {
  renderStore.completedDownloadUrl = null
  projectStore.isExportOpen = false
}
</script>
