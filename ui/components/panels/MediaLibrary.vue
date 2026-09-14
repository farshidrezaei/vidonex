<template>
  <div
    class="h-full flex flex-col bg-gray-900/60 select-none relative transition-colors"
    :class="isDraggingOver ? 'bg-indigo-950/20 ring-2 ring-inset ring-indigo-500/40' : ''"
    @dragover.prevent="isDraggingOver = true"
    @dragleave.prevent="handlePanelDragLeave"
    @drop.prevent="handleUploadDrop"
  >
    <!-- Hidden Custom Sleek Drag Ghost Badge -->
    <div
      ref="dragGhostRef"
      class="fixed -top-[9999px] left-0 pointer-events-none flex items-center gap-2 px-3 py-1.5 bg-gray-900 text-white rounded-full shadow-2xl border border-indigo-400 ring-2 ring-indigo-500/50 backdrop-blur text-xs font-medium z-[9999]"
    >
      <div class="w-5 h-5 rounded-full overflow-hidden flex items-center justify-center bg-indigo-950 border border-indigo-400/50 flex-shrink-0">
        <img
          v-if="ghostAsset?.thumbnail_path"
          :src="getThumbnailUrl(ghostAsset.thumbnail_path)"
          class="w-full h-full object-cover"
          alt=""
        />
        <UIcon
          v-else-if="ghostAsset?.file_type === 'audio'"
          name="i-heroicons-speaker-wave"
          class="w-3.5 h-3.5 text-emerald-400"
        />
        <UIcon
          v-else-if="getAssetCategory(ghostAsset) === 'subtitle'"
          name="i-heroicons-document-text"
          class="w-3.5 h-3.5 text-violet-400"
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

    <!-- Panel Header: Title, View Switcher & Add Button -->
    <div class="p-2.5 border-b border-gray-800 flex items-center justify-between gap-2 bg-gray-950/40 flex-shrink-0">
      <div class="flex items-center gap-1.5 min-w-0">
        <UIcon name="i-heroicons-folder" class="w-4 h-4 text-indigo-400 flex-shrink-0" />
        <span class="font-semibold text-xs tracking-wider uppercase text-gray-300 truncate">
          {{ $t('media.title') }}
        </span>
        <span class="text-[10px] font-mono px-1.5 py-0.2 rounded-full bg-gray-800 text-gray-400 border border-gray-700/60">
          {{ mediaStore.assets.length }}
        </span>
      </div>

      <div class="flex items-center gap-1 flex-shrink-0">
        <!-- View Switcher (Tree View vs Grid View) -->
        <div class="flex items-center bg-gray-900 border border-gray-800 rounded p-0.5">
          <button
            type="button"
            class="p-1 rounded text-xs transition-colors"
            :class="viewMode === 'tree' ? 'bg-indigo-600 text-white shadow-sm' : 'text-gray-400 hover:text-gray-200'"
            :title="$t('media.tree_view')"
            @click="viewMode = 'tree'"
          >
            <UIcon name="i-heroicons-bars-3-bottom-left" class="w-3.5 h-3.5" />
          </button>
          <button
            type="button"
            class="p-1 rounded text-xs transition-colors"
            :class="viewMode === 'grid' ? 'bg-indigo-600 text-white shadow-sm' : 'text-gray-400 hover:text-gray-200'"
            :title="$t('media.grid_view')"
            @click="viewMode = 'grid'"
          >
            <UIcon name="i-heroicons-squares-2x2" class="w-3.5 h-3.5" />
          </button>
        </div>

        <!-- Add / Upload Media Button -->
        <UButton
          size="xs"
          color="primary"
          variant="soft"
          icon="i-heroicons-plus"
          :loading="mediaStore.isUploading"
          :title="$t('media.upload')"
          @click="triggerFileSelect"
        >
          {{ $t('media.upload') }}
        </UButton>
      </div>
    </div>

    <!-- Fast Search Filter Bar -->
    <div class="p-2 border-b border-gray-800/60 flex items-center gap-1.5 bg-gray-950/20 flex-shrink-0">
      <div class="relative flex-1">
        <UIcon name="i-heroicons-magnifying-glass" class="w-3.5 h-3.5 text-gray-500 absolute left-2 top-1/2 -translate-y-1/2 pointer-events-none" />
        <input
          v-model="searchQuery"
          type="text"
          :placeholder="$t('media.search_placeholder')"
          class="w-full bg-gray-900/90 border border-gray-800 rounded-md pl-7 pr-7 py-1 text-xs text-gray-200 placeholder-gray-500 focus:outline-none focus:border-indigo-500 transition-colors"
        />
        <button
          v-if="searchQuery"
          type="button"
          class="absolute right-2 top-1/2 -translate-y-1/2 text-gray-500 hover:text-gray-300"
          @click="searchQuery = ''"
        >
          <UIcon name="i-heroicons-x-mark" class="w-3.5 h-3.5" />
        </button>
      </div>
      <input ref="fileInputRef" type="file" multiple accept="video/*,audio/*,image/*,.srt,.vtt" class="hidden" @change="handleFileInput" />
    </div>

    <!-- Media Library Content Area -->
    <div class="flex-1 overflow-y-auto px-2 py-2 space-y-1 min-h-0">
      <!-- Empty Library State -->
      <div
        v-if="mediaStore.assets.length === 0 && !mediaStore.isUploading"
        class="text-center py-10 px-4 space-y-3 border border-dashed border-gray-800 rounded-xl m-1 hover:border-gray-700 bg-gray-950/30 transition cursor-pointer"
        @click="triggerFileSelect"
      >
        <div class="w-10 h-10 rounded-full bg-gray-900 border border-gray-800 flex items-center justify-center mx-auto text-gray-400">
          <UIcon name="i-heroicons-arrow-up-tray" class="w-5 h-5 text-indigo-400" />
        </div>
        <div class="space-y-1">
          <p class="text-xs font-medium text-gray-300">{{ $t('media.empty') }}</p>
          <p class="text-[11px] text-gray-500 max-w-[200px] mx-auto leading-relaxed">
            {{ $t('media.dropzone') }}
          </p>
        </div>
      </div>

      <!-- Search No Results State -->
      <div v-else-if="filteredAssets.length === 0" class="text-center py-8 text-gray-500 text-xs">
        {{ $t('media.no_results') }}
      </div>

      <!-- MODE 1: COMPACT TREE VIEW (Default) -->
      <div v-else-if="viewMode === 'tree'" class="space-y-1.5">
        <div
          v-for="category in activeCategories"
          :key="category.id"
          class="rounded-lg border border-gray-800/60 bg-gray-950/40 overflow-hidden"
        >
          <!-- Collapsible Category Branch Header -->
          <div
            class="flex items-center justify-between px-2.5 py-1.5 hover:bg-gray-800/50 cursor-pointer select-none transition-colors group"
            @click="toggleCategory(category.id)"
          >
            <div class="flex items-center gap-1.5 min-w-0">
              <UIcon
                name="i-heroicons-chevron-right"
                class="w-3.5 h-3.5 text-gray-500 transition-transform duration-200"
                :class="expandedCategories[category.id] ? 'rotate-90 text-gray-300' : ''"
              />
              <UIcon :name="category.icon" class="w-3.5 h-3.5 flex-shrink-0" :class="category.textColor" />
              <span class="text-xs font-semibold text-gray-300 tracking-wide truncate">
                {{ $t(category.labelKey) }}
              </span>
            </div>

            <span
              class="text-[10px] font-mono px-1.5 py-0.5 rounded-full"
              :class="category.badgeColor"
            >
              {{ categorizedAssets[category.id].length }}
            </span>
          </div>

          <!-- Category Children Items (Leaf Rows) -->
          <div
            v-if="expandedCategories[category.id]"
            class="px-1 pb-1 pt-0.5 space-y-0.5"
          >
            <div
              v-for="asset in categorizedAssets[category.id]"
              :key="asset.id"
              draggable="true"
              class="group/item flex items-center justify-between px-2 py-1 rounded-md text-xs hover:bg-gray-800/80 border border-transparent hover:border-gray-700/60 cursor-grab active:cursor-grabbing transition-all duration-150"
              :class="mediaStore.draggedAsset?.id === asset.id ? 'opacity-30 scale-95 border-indigo-500 bg-indigo-500/10' : ''"
              :title="$t('media.drag_hint')"
              @dragstart="handleAssetDragStart($event, asset)"
              @dragend="handleAssetDragEnd"
            >
              <!-- Left: Inline Thumbnail / Icon & File Name -->
              <div class="flex items-center gap-2 min-w-0 flex-1 pointer-events-none">
                <div class="w-5 h-5 rounded overflow-hidden flex-shrink-0 bg-gray-900 border border-gray-800 flex items-center justify-center relative">
                  <img
                    v-if="asset.thumbnail_path"
                    :src="getThumbnailUrl(asset.thumbnail_path)"
                    class="w-full h-full object-cover"
                    alt=""
                  />
                  <UIcon
                    v-else-if="asset.file_type === 'audio'"
                    name="i-heroicons-speaker-wave"
                    class="w-3 h-3 text-emerald-400"
                  />
                  <UIcon
                    v-else-if="getAssetCategory(asset) === 'subtitle'"
                    name="i-heroicons-document-text"
                    class="w-3 h-3 text-violet-400"
                  />
                  <UIcon
                    v-else
                    name="i-heroicons-film"
                    class="w-3 h-3 text-indigo-400"
                  />
                </div>

                <!-- File Name -->
                <span class="truncate text-[11px] font-medium text-gray-200 group-hover/item:text-white">
                  {{ asset.file_name }}
                </span>
              </div>

              <!-- Right: Metadata Badge & Hover Actions -->
              <div class="flex items-center gap-1.5 flex-shrink-0">
                <!-- Meta Info Tag (Duration or Size) -->
                <span
                  v-if="asset.duration_seconds > 0"
                  class="text-[10px] font-mono text-gray-500 group-hover/item:text-gray-400"
                >
                  {{ formatSeconds(asset.duration_seconds) }}
                </span>
                <span
                  v-else-if="asset.file_size_bytes > 0"
                  class="text-[10px] font-mono text-gray-500 group-hover/item:text-gray-400"
                >
                  {{ formatFileSize(asset.file_size_bytes) }}
                </span>

                <!-- Action Buttons: Visible on row hover (NO Add-to-timeline button, pure drag & drop!) -->
                <div class="flex items-center gap-0.5 opacity-0 group-hover/item:opacity-100 transition-opacity" @mousedown.stop>
                  <!-- Inspect Metadata -->
                  <button
                    type="button"
                    class="p-1 rounded hover:bg-gray-700 text-gray-400 hover:text-gray-200 transition-colors"
                    :title="$t('media.probe_details')"
                    @click="mediaStore.openProbeModal(asset)"
                  >
                    <UIcon name="i-heroicons-information-circle" class="w-3.5 h-3.5" />
                  </button>
                  <!-- Delete Asset -->
                  <button
                    type="button"
                    class="p-1 rounded hover:bg-red-950/80 text-gray-400 hover:text-red-400 transition-colors"
                    :title="$t('media.delete')"
                    @click="mediaStore.deleteAsset(asset.id)"
                  >
                    <UIcon name="i-heroicons-trash" class="w-3.5 h-3.5" />
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- MODE 2: GRID / CARD VIEW -->
      <div v-else-if="viewMode === 'grid'" class="grid grid-cols-2 gap-2">
        <div
          v-for="asset in filteredAssets"
          :key="asset.id"
          draggable="true"
          class="group bg-gray-950/60 hover:bg-gray-800/60 border border-gray-800/80 hover:border-indigo-500/60 rounded-lg p-1.5 transition-all duration-150 flex flex-col gap-1.5 cursor-grab active:cursor-grabbing relative"
          :class="mediaStore.draggedAsset?.id === asset.id ? 'opacity-30 scale-95 border-indigo-500' : ''"
          :title="$t('media.drag_hint')"
          @dragstart="handleAssetDragStart($event, asset)"
          @dragend="handleAssetDragEnd"
        >
          <!-- Thumbnail Container -->
          <div class="w-full aspect-video rounded bg-gray-900 overflow-hidden relative border border-gray-800 flex items-center justify-center">
            <img
              v-if="asset.thumbnail_path"
              :src="getThumbnailUrl(asset.thumbnail_path)"
              class="w-full h-full object-cover pointer-events-none"
              alt=""
            />
            <UIcon
              v-else-if="asset.file_type === 'audio'"
              name="i-heroicons-speaker-wave"
              class="w-6 h-6 text-emerald-400"
            />
            <UIcon
              v-else-if="getAssetCategory(asset) === 'subtitle'"
              name="i-heroicons-document-text"
              class="w-6 h-6 text-violet-400"
            />
            <UIcon
              v-else
              name="i-heroicons-film"
              class="w-6 h-6 text-indigo-400"
            />

            <!-- Duration Badge -->
            <span
              v-if="asset.duration_seconds > 0"
              class="absolute bottom-1 right-1 bg-black/80 text-[9px] font-mono text-gray-300 px-1 rounded pointer-events-none"
            >
              {{ formatSeconds(asset.duration_seconds) }}
            </span>
          </div>

          <!-- Title & Actions -->
          <div class="flex items-center justify-between gap-1 min-w-0">
            <p class="text-[11px] font-medium text-gray-200 truncate flex-1" :title="asset.file_name">
              {{ asset.file_name }}
            </p>

            <div class="flex items-center gap-0.5 opacity-0 group-hover:opacity-100 transition-opacity" @mousedown.stop>
              <button
                type="button"
                class="p-0.5 rounded text-gray-400 hover:text-gray-200"
                :title="$t('media.probe_details')"
                @click="mediaStore.openProbeModal(asset)"
              >
                <UIcon name="i-heroicons-information-circle" class="w-3.5 h-3.5" />
              </button>
              <button
                type="button"
                class="p-0.5 rounded text-gray-400 hover:text-red-400"
                :title="$t('media.delete')"
                @click="mediaStore.deleteAsset(asset.id)"
              >
                <UIcon name="i-heroicons-trash" class="w-3.5 h-3.5" />
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Drag Hint Footer Bar -->
    <div
      v-if="mediaStore.assets.length > 0"
      class="px-2.5 py-1.5 border-t border-gray-800/80 bg-gray-950/60 flex items-center justify-between text-[10px] text-gray-500 flex-shrink-0"
    >
      <span class="flex items-center gap-1 truncate">
        <UIcon name="i-heroicons-cursor-arrow-rays" class="w-3.5 h-3.5 text-indigo-400" />
        <span>{{ $t('media.drag_hint') }}</span>
      </span>

      <button
        type="button"
        class="text-indigo-400 hover:text-indigo-300 font-medium transition-colors cursor-pointer"
        @click="triggerFileSelect"
      >
        + {{ $t('media.upload') }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useMediaStore } from '~/stores/media'
import { useProjectStore } from '~/stores/project'
import { useDesktop } from '~/composables/useDesktop'
import type { MediaAsset } from '~/types/project'

interface CategoryDefinition {
  id: 'video' | 'audio' | 'image' | 'subtitle' | 'other'
  labelKey: string
  icon: string
  badgeColor: string
  textColor: string
}

const CATEGORIES: CategoryDefinition[] = [
  {
    id: 'video',
    labelKey: 'media.categories.video',
    icon: 'i-heroicons-film',
    badgeColor: 'bg-indigo-500/20 text-indigo-300 border border-indigo-500/30',
    textColor: 'text-indigo-400',
  },
  {
    id: 'audio',
    labelKey: 'media.categories.audio',
    icon: 'i-heroicons-speaker-wave',
    badgeColor: 'bg-emerald-500/20 text-emerald-300 border border-emerald-500/30',
    textColor: 'text-emerald-400',
  },
  {
    id: 'image',
    labelKey: 'media.categories.image',
    icon: 'i-heroicons-photo',
    badgeColor: 'bg-amber-500/20 text-amber-300 border border-amber-500/30',
    textColor: 'text-amber-400',
  },
  {
    id: 'subtitle',
    labelKey: 'media.categories.subtitle',
    icon: 'i-heroicons-document-text',
    badgeColor: 'bg-violet-500/20 text-violet-300 border border-violet-500/30',
    textColor: 'text-violet-400',
  },
]

const mediaStore = useMediaStore()
const projectStore = useProjectStore()
const { isDesktop, selectMediaFiles, importLocalMedia, resolveMediaUrl } = useDesktop()

const searchQuery = ref('')
const viewMode = ref<'tree' | 'grid'>('tree')
const isDraggingOver = ref(false)
const dragGhostRef = ref<HTMLDivElement | null>(null)
const ghostAsset = ref<MediaAsset | null>(null)
const fileInputRef = ref<HTMLInputElement | null>(null)

const expandedCategories = ref<Record<string, boolean>>({
  video: true,
  audio: true,
  image: true,
  subtitle: true,
  other: true,
})

function toggleCategory(categoryId: string) {
  expandedCategories.value[categoryId] = !expandedCategories.value[categoryId]
}

function getAssetCategory(asset?: MediaAsset | null): 'video' | 'audio' | 'image' | 'subtitle' | 'other' {
  if (!asset) return 'other'
  const ext = (asset.file_name || asset.file_path || '').toLowerCase().split('.').pop() || ''
  if (['srt', 'vtt', 'ass', 'sub'].includes(ext)) {
    return 'subtitle'
  }
  if (asset.file_type === 'video' || ['mp4', 'mov', 'mkv', 'webm', 'avi', 'flv'].includes(ext)) {
    return 'video'
  }
  if (asset.file_type === 'audio' || ['mp3', 'wav', 'aac', 'flac', 'ogg', 'm4a'].includes(ext)) {
    return 'audio'
  }
  if (asset.file_type === 'image' || ['png', 'jpg', 'jpeg', 'webp', 'svg', 'gif', 'bmp'].includes(ext)) {
    return 'image'
  }
  return 'other'
}

const filteredAssets = computed(() => {
  const query = searchQuery.value.trim().toLowerCase()
  if (!query) return mediaStore.assets
  return mediaStore.assets.filter((a) => a.file_name.toLowerCase().includes(query))
})

const categorizedAssets = computed(() => {
  const groups: Record<string, MediaAsset[]> = {
    video: [],
    audio: [],
    image: [],
    subtitle: [],
    other: [],
  }
  for (const asset of filteredAssets.value) {
    const cat = getAssetCategory(asset)
    if (!groups[cat]) {
      groups[cat] = []
    }
    groups[cat].push(asset)
  }
  return groups
})

// Only display categories that contain assets (or all if library has items and query is empty)
const activeCategories = computed(() => {
  return CATEGORIES.filter((cat) => (categorizedAssets.value[cat.id]?.length || 0) > 0)
})

function getThumbnailUrl(path?: string): string {
  return resolveMediaUrl(path)
}

async function triggerFileSelect() {
  if (isDesktop.value && projectStore.currentProject) {
    const selectedPaths = await selectMediaFiles()
    if (selectedPaths && selectedPaths.length > 0) {
      await importLocalMedia(projectStore.currentProject.id, selectedPaths)
      await mediaStore.fetchAssets(projectStore.currentProject.id)
      return
    }
  }
  fileInputRef.value?.click()
}

function handleAssetDragStart(event: DragEvent, asset: MediaAsset) {
  ghostAsset.value = asset
  mediaStore.setDraggedAsset(asset)

  if (event.dataTransfer) {
    event.dataTransfer.setData('vidonex-asset-id', asset.id)
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

function handlePanelDragLeave(event: DragEvent) {
  const target = event.currentTarget as HTMLElement
  if (target && !target.contains(event.relatedTarget as Node)) {
    isDraggingOver.value = false
  }
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
