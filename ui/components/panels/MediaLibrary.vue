<template>
  <div class="h-full flex flex-col bg-gray-900/60 select-none relative">
    <!-- Hidden Custom Sleek Drag Ghost Badge -->
    <div
      ref="dragGhostRef"
      class="fixed -top-[9999px] left-0 pointer-events-none flex items-center gap-2 px-3 py-1.5 bg-gray-900 text-white rounded-full shadow-2xl border border-indigo-400 ring-2 ring-indigo-500/50 backdrop-blur text-xs font-medium z-[9999]"
    >
      <div class="w-5 h-5 rounded-full overflow-hidden flex items-center justify-center bg-indigo-950 border border-indigo-400/50 flex-shrink-0">
        <img
          v-if="ghostAsset?.thumbnail_path"
          :src="`/api/media/files/${ghostAsset.thumbnail_path}`"
          class="w-full h-full object-cover"
          alt=""
        />
        <UIcon
          v-else-if="ghostAsset?.file_type === 'audio'"
          name="i-heroicons-speaker-wave"
          class="w-3.5 h-3.5 text-emerald-400"
        />
        <UIcon
          v-else
          name="i-heroicons-film"
          class="w-3.5 h-3.5 text-indigo-400"
        />
      </div>

      <span class="max-w-[130px] truncate text-[11px] font-medium text-gray-100">
        {{ ghostAsset?.file_name || 'Media Asset' }}
      </span>

      <span class="text-[9px] bg-indigo-500/30 text-indigo-200 px-1.5 py-0.5 rounded-full font-mono uppercase tracking-wider">
        {{ ghostAsset?.file_type }}
      </span>
    </div>

    <!-- Header -->
    <div class="p-3 border-b border-gray-800 flex items-center justify-between">
      <div class="flex items-center gap-2">
        <UIcon name="i-heroicons-folder" class="w-4 h-4 text-indigo-400" />
        <span class="font-semibold text-xs tracking-wider uppercase text-gray-300">{{ $t('media.title') }}</span>
      </div>

      <div>
        <UButton
          size="xs"
          color="primary"
          variant="soft"
          icon="i-heroicons-plus"
          :loading="mediaStore.isUploading"
          @click="triggerFileSelect"
        >
          {{ $t('media.upload') }}
        </UButton>
      </div>
    </div>

    <!-- Dropzone Area -->
    <div
      class="p-3 m-3 border-2 border-dashed rounded-xl transition flex flex-col items-center justify-center text-center gap-2 cursor-pointer"
      :class="isDraggingOver ? 'border-indigo-500 bg-indigo-500/10' : 'border-gray-800 hover:border-gray-700 bg-gray-950/40'"
      @dragover.prevent="isDraggingOver = true"
      @dragleave.prevent="isDraggingOver = false"
      @drop.prevent="handleUploadDrop"
      @click="triggerFileSelect"
    >
      <input ref="fileInputRef" type="file" multiple accept="video/*,audio/*,image/*" class="hidden" @change="handleFileInput" />
      <div class="w-10 h-10 rounded-full bg-gray-800/80 flex items-center justify-center text-gray-400 group-hover:text-indigo-400">
        <UIcon name="i-heroicons-arrow-up-tray" class="w-5 h-5" />
      </div>
      <p class="text-xs text-gray-400 max-w-[200px] leading-relaxed">
        {{ $t('media.dropzone') }}
      </p>
    </div>

    <!-- Assets List -->
    <div class="flex-1 overflow-y-auto px-3 pb-3 space-y-2">
      <div v-if="mediaStore.assets.length === 0 && !mediaStore.isUploading" class="text-center py-8 text-gray-500 text-xs">
        {{ $t('media.empty') }}
      </div>

      <!-- Draggable Media Asset Card -->
      <div
        v-for="asset in mediaStore.assets"
        :key="asset.id"
        draggable="true"
        class="group bg-gray-950/60 hover:bg-gray-800/60 border border-gray-800/80 hover:border-indigo-500/60 rounded-lg p-2 transition-all duration-200 flex items-center gap-3 relative cursor-grab active:cursor-grabbing"
        :class="mediaStore.draggedAsset?.id === asset.id ? 'opacity-30 scale-95 border-indigo-500 ring-2 ring-indigo-500/30' : 'hover:shadow-lg hover:shadow-indigo-500/10'"
        @dragstart="handleAssetDragStart($event, asset)"
        @dragend="handleAssetDragEnd"
      >
        <!-- Thumbnail / Icon -->
        <div class="w-12 h-12 rounded bg-gray-900 flex-shrink-0 flex items-center justify-center overflow-hidden relative border border-gray-800">
          <img
            v-if="asset.thumbnail_path"
            :src="`/api/media/files/${asset.thumbnail_path}`"
            class="w-full h-full object-cover pointer-events-none"
            alt=""
          />
          <UIcon
            v-else-if="asset.file_type === 'audio'"
            name="i-heroicons-speaker-wave"
            class="w-6 h-6 text-emerald-400 pointer-events-none"
          />
          <UIcon
            v-else
            name="i-heroicons-film"
            class="w-6 h-6 text-indigo-400 pointer-events-none"
          />

          <!-- Duration Badge -->
          <span
            v-if="asset.duration_seconds > 0"
            class="absolute bottom-0.5 right-0.5 bg-black/80 text-[9px] font-mono text-gray-300 px-1 rounded pointer-events-none"
          >
            {{ formatSeconds(asset.duration_seconds) }}
          </span>
        </div>

        <!-- Info -->
        <div class="flex-1 min-w-0 pointer-events-none">
          <p class="text-xs font-medium text-gray-200 truncate" :title="asset.file_name">
            {{ asset.file_name }}
          </p>
          <div class="flex items-center gap-2 mt-1 text-[10px] text-gray-500 font-mono">
            <span v-if="asset.width && asset.height">{{ asset.width }}x{{ asset.height }}</span>
            <span v-if="asset.frame_rate">{{ Math.round(asset.frame_rate) }}fps</span>
            <span>{{ formatFileSize(asset.file_size_bytes) }}</span>
          </div>
        </div>

        <!-- Action Buttons -->
        <div class="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition z-10" @mousedown.stop>
          <!-- Add to Timeline -->
          <UButton
            icon="i-heroicons-plus-circle"
            size="2xs"
            color="primary"
            variant="ghost"
            :title="$t('media.add_to_timeline')"
            @click="addAssetToTimeline(asset)"
          />
          <!-- Probe -->
          <UButton
            icon="i-heroicons-information-circle"
            size="2xs"
            color="neutral"
            variant="ghost"
            :title="$t('media.probe_details')"
            @click="mediaStore.openProbeModal(asset)"
          />
          <!-- Delete -->
          <UButton
            icon="i-heroicons-trash"
            size="2xs"
            color="error"
            variant="ghost"
            :title="$t('media.delete')"
            @click="mediaStore.deleteAsset(asset.id)"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useMediaStore } from '~/stores/media'
import { useProjectStore } from '~/stores/project'
import { useTimelineStore } from '~/stores/timeline'
import type { MediaAsset } from '~/types/project'

const mediaStore = useMediaStore()
const projectStore = useProjectStore()
const timelineStore = useTimelineStore()

const isDraggingOver = ref(false)
const dragGhostRef = ref<HTMLDivElement | null>(null)
const ghostAsset = ref<MediaAsset | null>(null)
const fileInputRef = ref<HTMLInputElement | null>(null)

function triggerFileSelect() {
  fileInputRef.value?.click()
}

function handleAssetDragStart(event: DragEvent, asset: MediaAsset) {
  ghostAsset.value = asset
  mediaStore.setDraggedAsset(asset)

  if (event.dataTransfer) {
    event.dataTransfer.setData('vidonyx-asset-id', asset.id)
    event.dataTransfer.setData('application/json', JSON.stringify(asset))
    event.dataTransfer.setData('text/plain', asset.id)
    event.dataTransfer.effectAllowed = 'all'

    // Set custom sleek drag ghost element
    if (dragGhostRef.value) {
      event.dataTransfer.setDragImage(dragGhostRef.value, 30, 15)
    }
  }
}

function handleAssetDragEnd() {
  setTimeout(() => {
    mediaStore.setDraggedAsset(null)
    ghostAsset.value = null
  }, 150)
}

async function handleFileInput(event: Event) {
  const target = event.target as HTMLInputElement
  if (!target.files || !projectStore.currentProject) return

  for (const file of Array.from(target.files)) {
    await mediaStore.uploadFile(projectStore.currentProject.id, file)
  }
}

async function handleUploadDrop(event: DragEvent) {
  isDraggingOver.value = false
  if (!event.dataTransfer?.files || !projectStore.currentProject) return

  for (const file of Array.from(event.dataTransfer.files)) {
    await mediaStore.uploadFile(projectStore.currentProject.id, file)
  }
}

function addAssetToTimeline(asset: MediaAsset) {
  const targetKind = asset.file_type === 'audio' ? 'audio' : 'video'
  let targetTrack = timelineStore.tracks.find((t) => t.kind === targetKind)

  if (!targetTrack) {
    timelineStore.addTrack(targetKind)
    targetTrack = timelineStore.tracks[timelineStore.tracks.length - 1]
  }

  timelineStore.addClipToTrack(targetTrack.id, asset)
}

function formatSeconds(sec: number): string {
  const m = Math.floor(sec / 60)
  const s = Math.floor(sec % 60)
  return `${m}:${String(s).padStart(2, '0')}`
}

function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i]
}
</script>
