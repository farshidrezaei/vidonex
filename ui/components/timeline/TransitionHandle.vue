<template>
  <div
    class="absolute inset-y-0 z-20 pointer-events-none"
    :style="{ left: `${positionX}px` }"
  >
    <!-- Visual Transition Span Ribbon Across Both Clips -->
    <div
      v-if="transition"
      class="absolute top-1.5 bottom-1.5 -translate-x-1/2 rounded-md border border-amber-500/80 bg-gradient-to-r from-amber-500/25 via-amber-500/40 to-amber-500/25 shadow-[0_0_12px_rgba(245,158,11,0.25)] backdrop-blur-xs flex items-center justify-between px-2 cursor-pointer transition-all hover:bg-amber-500/35 hover:border-amber-400 select-none group/trans pointer-events-auto overflow-hidden"
      :style="{ width: `${widthPixels}px` }"
      :title="`Transition: ${transition.type} (${transition.duration}s)\nClick to edit`"
      @click.stop="openModal"
    >
      <!-- Diagonal Pattern Accent -->
      <div class="absolute inset-0 opacity-15 pointer-events-none bg-[linear-gradient(45deg,rgba(255,255,255,0.4)_25%,transparent_25%,transparent_50%,rgba(255,255,255,0.4)_50%,rgba(255,255,255,0.4)_75%,transparent_75%,transparent)] bg-[length:12px_12px]"></div>

      <!-- Left Half: Covers fromClip end -->
      <div class="flex items-center gap-1 z-10 min-w-0 pr-1">
        <UIcon :name="transitionIcon" class="w-3.5 h-3.5 text-amber-300 flex-shrink-0" />
        <span class="text-[10px] font-semibold font-mono text-amber-200 truncate capitalize">
          {{ transition.type }}
        </span>
      </div>

      <!-- Center Cut Divider Line & Badge -->
      <div class="absolute left-1/2 top-0 bottom-0 -translate-x-1/2 w-0 border-r-2 border-dashed border-amber-300/80 pointer-events-none flex flex-col items-center justify-center">
        <div class="w-4 h-4 rounded-full bg-amber-500 text-black shadow-md flex items-center justify-center ring-2 ring-amber-300/60">
          <UIcon :name="transitionIcon" class="w-2.5 h-2.5 pointer-events-none" />
        </div>
      </div>

      <!-- Right Half: Covers toClip start -->
      <div class="flex items-center gap-1 z-10 pl-1">
        <span class="text-[10px] font-mono font-bold text-amber-950 bg-amber-400/90 px-1.5 py-0.5 rounded shadow-sm">
          {{ Number(transition.duration).toFixed(1) }}s
        </span>
      </div>
    </div>

    <!-- Small Plus Handle Button (When No Transition) -->
    <div
      v-else
      class="absolute top-1/2 -translate-y-1/2 -translate-x-1/2 pointer-events-auto"
    >
      <button
        class="w-5 h-5 rounded-full flex items-center justify-center transition shadow-md cursor-pointer bg-gray-800 hover:bg-gray-700 text-gray-400 border border-gray-700 hover:scale-110"
        :title="$t('timeline.add_transition') || 'Add Transition'"
        @click.stop="openModal"
      >
        <UIcon name="i-heroicons-plus" class="w-3 h-3 pointer-events-none" />
      </button>
    </div>

    <!-- Transition Settings Modal -->
    <UModal
      v-model:open="isModalOpen"
      :title="$t('timeline.transitions.modal_title') || 'Transition Settings'"
      :ui="{ content: 'sm:max-w-lg bg-gray-950 border border-gray-800 pointer-events-auto' }"
    >
      <template #body>
        <div class="space-y-5 text-xs">
          <!-- Clip Junction Breadcrumb -->
          <div class="flex items-center justify-between p-2 rounded bg-gray-900/60 border border-gray-800 text-[11px] text-gray-400">
            <div class="flex items-center gap-1.5 truncate max-w-[45%]">
              <UIcon name="i-heroicons-film" class="w-3.5 h-3.5 text-indigo-400 flex-shrink-0" />
              <span class="truncate font-mono text-gray-200">{{ fromClip.id }}</span>
            </div>
            <UIcon name="i-heroicons-arrow-right" class="w-3.5 h-3.5 text-amber-400 flex-shrink-0" />
            <div class="flex items-center gap-1.5 truncate max-w-[45%]">
              <UIcon name="i-heroicons-film" class="w-3.5 h-3.5 text-indigo-400 flex-shrink-0" />
              <span class="truncate font-mono text-gray-200">{{ toClip.id }}</span>
            </div>
          </div>

          <!-- Categorized Transition Gallery -->
          <div class="space-y-3">
            <div class="flex items-center justify-between gap-2">
              <label class="block text-gray-400 font-medium uppercase tracking-wider text-[10px]">
                {{ $t('timeline.transitions.select_type') || 'Transition Type' }}
              </label>
              <!-- Search filter -->
              <div class="relative w-44">
                <UIcon name="i-heroicons-magnifying-glass" class="w-3.5 h-3.5 text-gray-500 absolute left-2 top-1/2 -translate-y-1/2 pointer-events-none" />
                <input
                  v-model="transitionSearchQuery"
                  type="text"
                  placeholder="Filter transitions..."
                  class="w-full bg-gray-900 border border-gray-800 rounded-md pl-7 pr-2 py-1 text-xs text-gray-200 placeholder-gray-500 focus:outline-none focus:border-amber-500"
                />
              </div>
            </div>

            <!-- Category Pills Filter -->
            <div class="flex items-center gap-1 overflow-x-auto pb-1 text-[11px] scrollbar-none">
              <button
                v-for="cat in transitionCategories"
                :key="cat.id"
                type="button"
                class="px-2.5 py-1 rounded-full whitespace-nowrap transition cursor-pointer"
                :class="selectedCategory === cat.id
                  ? 'bg-amber-500 text-black font-semibold shadow-sm'
                  : 'bg-gray-900 hover:bg-gray-800 text-gray-400 hover:text-gray-200 border border-gray-800'"
                @click="selectedCategory = cat.id"
              >
                {{ cat.label }} ({{ cat.count }})
              </button>
            </div>

            <!-- Transition Cards Grid (Scrollable) -->
            <div class="grid grid-cols-2 sm:grid-cols-3 gap-2 max-h-64 overflow-y-auto pr-1">
              <button
                v-for="opt in filteredTransitionOptions"
                :key="opt.value"
                type="button"
                class="flex items-center gap-2.5 p-2 rounded-lg border text-left transition cursor-pointer"
                :class="selectedType === opt.value
                  ? (opt.value === 'none'
                      ? 'border-gray-500 bg-gray-800/80 text-white ring-1 ring-gray-400/50'
                      : 'border-amber-500 bg-amber-500/15 text-amber-300 ring-1 ring-amber-400/50')
                  : 'border-gray-800/80 bg-gray-900/50 hover:bg-gray-850 text-gray-300 hover:border-gray-700'"
                @click="selectedType = opt.value"
              >
                <div
                  class="w-7 h-7 rounded flex items-center justify-center flex-shrink-0"
                  :class="selectedType === opt.value
                    ? (opt.value === 'none' ? 'bg-gray-700 text-gray-200' : 'bg-amber-500 text-black')
                    : 'bg-gray-800 text-gray-400'"
                >
                  <UIcon :name="opt.icon" class="w-4 h-4" />
                </div>
                <div class="min-w-0">
                  <div class="font-medium text-[11px] truncate">{{ opt.label }}</div>
                  <div class="text-[9px] text-gray-500 truncate">{{ opt.desc }}</div>
                </div>
              </button>
            </div>
            <div v-if="filteredTransitionOptions.length === 0" class="text-center py-6 text-gray-500 text-xs">
              No transitions found matching "{{ transitionSearchQuery }}"
            </div>
          </div>

          <!-- Duration Section (Hidden if 'none') -->
          <div v-if="selectedType !== 'none'" class="space-y-3 pt-3 border-t border-gray-800/80">
            <div class="flex items-center justify-between">
              <label class="text-gray-400 font-medium uppercase tracking-wider text-[10px]">
                {{ $t('timeline.transitions.duration_label') || 'Duration (Seconds)' }}
              </label>
              <span class="font-mono text-amber-300 text-xs font-semibold">{{ selectedDuration.toFixed(2) }}s</span>
            </div>

            <!-- Number input and slider -->
            <div class="flex items-center gap-3">
              <div class="relative w-24">
                <input
                  v-model.number="selectedDuration"
                  type="number"
                  step="0.1"
                  min="0.1"
                  :max="maxAllowedDuration"
                  class="w-full bg-gray-900 border border-gray-800 rounded px-2.5 py-1.5 font-mono text-gray-200 focus:outline-none focus:border-amber-500 text-xs pr-6"
                />
                <span class="absolute right-2 top-1/2 -translate-y-1/2 text-gray-500 text-[10px]">s</span>
              </div>
              <input
                v-model.number="selectedDuration"
                type="range"
                step="0.05"
                min="0.1"
                :max="maxAllowedDuration"
                class="flex-1 accent-amber-500 h-1.5 bg-gray-800 rounded cursor-pointer"
              />
            </div>

            <!-- Quick Duration Presets -->
            <div class="flex items-center gap-1.5 pt-1">
              <span class="text-[10px] text-gray-500 mr-1">Presets:</span>
              <button
                v-for="preset in availablePresets"
                :key="preset"
                type="button"
                class="px-2 py-0.5 rounded text-[10px] font-mono border transition cursor-pointer"
                :class="Math.abs(selectedDuration - preset) < 0.01
                  ? 'bg-amber-500 text-black border-amber-400 font-semibold'
                  : 'bg-gray-900 text-gray-400 border-gray-800 hover:border-gray-700 hover:text-gray-200'"
                @click="selectedDuration = preset"
              >
                {{ preset.toFixed(1) }}s
              </button>
            </div>
          </div>

          <!-- Modal Footer Actions -->
          <div class="flex items-center justify-between pt-4 border-t border-gray-800/80">
            <div>
              <UButton
                v-if="transition"
                color="error"
                variant="ghost"
                size="sm"
                icon="i-heroicons-trash"
                @click="deleteTransition"
              >
                {{ $t('timeline.transitions.remove') || 'Remove' }}
              </UButton>
            </div>
            <div class="flex items-center gap-2">
              <UButton
                color="neutral"
                variant="ghost"
                size="sm"
                @click="isModalOpen = false"
              >
                {{ $t('app.cancel') || 'Cancel' }}
              </UButton>
              <UButton
                size="sm"
                icon="i-heroicons-check"
                class="bg-amber-500 hover:bg-amber-400 text-black font-semibold"
                @click="applyTransition"
              >
                {{ $t('app.save') || 'Apply' }}
              </UButton>
            </div>
          </div>
        </div>
      </template>
    </UModal>
  </div>
</template>

<script setup lang="ts">
import { useTimelineStore } from '~/stores/timeline'
import type { TransitionSpec, ClipSpec } from '~/types/spec'

const props = defineProps<{
  trackId: string
  fromClip: ClipSpec
  toClip: ClipSpec
  transition?: TransitionSpec
}>()

const { t } = useI18n()
const timelineStore = useTimelineStore()

const positionX = computed(() => {
  const fromEnd = (Number(props.fromClip.start) || 0) + (Number(props.fromClip.duration) || 0)
  return fromEnd * timelineStore.pixelsPerSecond
})

const transitionDuration = computed(() => Math.max(0.1, parseFloat(String(props.transition?.duration)) || 1.0))
const widthPixels = computed(() => transitionDuration.value * timelineStore.pixelsPerSecond)

const isModalOpen = ref(false)
const selectedType = ref('dissolve')
const selectedDuration = ref(1.0)

const maxClipDuration = computed(() => {
  const durA = parseFloat(String(props.fromClip?.duration)) || 5.0
  const durB = parseFloat(String(props.toClip?.duration)) || 5.0
  return Math.min(durA, durB)
})

const maxAllowedDuration = computed(() => {
  // Centered transition covers D/2 from each clip, so D <= 2 * min(durA, durB)
  const maxD = maxClipDuration.value * 2
  return Math.max(0.2, Math.min(10.0, Number(maxD.toFixed(2))))
})

const availablePresets = computed(() => {
  const allPresets = [0.5, 1.0, 1.5, 2.0, 3.0, 4.0]
  const filtered = allPresets.filter((p) => p <= maxAllowedDuration.value)
  if (filtered.length === 0) {
    return [Number(maxAllowedDuration.value.toFixed(1))]
  }
  return filtered
})

const transitionSearchQuery = ref('')
const selectedCategory = ref('all')

const transitionCategories = computed(() => {
  const counts: Record<string, number> = {
    all: allTransitionsCatalog.length,
    fades: allTransitionsCatalog.filter((t) => t.category === 'fades').length,
    wipes: allTransitionsCatalog.filter((t) => t.category === 'wipes').length,
    slides: allTransitionsCatalog.filter((t) => t.category === 'slides').length,
    shapes: allTransitionsCatalog.filter((t) => t.category === 'shapes').length,
    zooms: allTransitionsCatalog.filter((t) => t.category === 'zooms').length,
  }
  return [
    { id: 'all', label: 'All', count: counts.all },
    { id: 'fades', label: 'Fades & Dissolve', count: counts.fades },
    { id: 'wipes', label: 'Wipes', count: counts.wipes },
    { id: 'slides', label: 'Slides & Pushes', count: counts.slides },
    { id: 'shapes', label: 'Shapes & Iris', count: counts.shapes },
    { id: 'zooms', label: 'Zooms & Warps', count: counts.zooms },
  ]
})

interface TransitionItem {
  value: string
  label: string
  desc: string
  category: 'fades' | 'wipes' | 'slides' | 'shapes' | 'zooms' | 'none'
  icon: string
}

const allTransitionsCatalog: TransitionItem[] = [
  { value: 'none', label: 'None (Cut)', desc: 'Cut immediately', category: 'none', icon: 'i-heroicons-no-symbol' },
  // Fades
  { value: 'dissolve', label: 'Cross Dissolve', desc: 'Smooth cross fade', category: 'fades', icon: 'i-heroicons-sparkles' },
  { value: 'fade', label: 'Dip to Black', desc: 'Fade through black', category: 'fades', icon: 'i-heroicons-eye-slash' },
  { value: 'fadeblack', label: 'Fade Black', desc: 'Full dip black', category: 'fades', icon: 'i-heroicons-moon' },
  { value: 'fadewhite', label: 'Flash to White', desc: 'Bright flash fade', category: 'fades', icon: 'i-heroicons-sun' },
  { value: 'fadegrays', label: 'Fade Grayscale', desc: 'Desaturate cross fade', category: 'fades', icon: 'i-heroicons-adjustments-horizontal' },
  // Wipes
  { value: 'wipeleft', label: 'Wipe Left', desc: 'Linear wipe to left', category: 'wipes', icon: 'i-heroicons-arrow-left' },
  { value: 'wiperight', label: 'Wipe Right', desc: 'Linear wipe to right', category: 'wipes', icon: 'i-heroicons-arrow-right' },
  { value: 'wipeup', label: 'Wipe Up', desc: 'Linear wipe upward', category: 'wipes', icon: 'i-heroicons-arrow-up' },
  { value: 'wipedown', label: 'Wipe Down', desc: 'Linear wipe downward', category: 'wipes', icon: 'i-heroicons-arrow-down' },
  { value: 'wipetl', label: 'Wipe Top-Left', desc: 'Diagonal wipe top-left', category: 'wipes', icon: 'i-heroicons-arrow-up-left' },
  { value: 'wipetr', label: 'Wipe Top-Right', desc: 'Diagonal wipe top-right', category: 'wipes', icon: 'i-heroicons-arrow-up-right' },
  { value: 'wipebl', label: 'Wipe Bottom-Left', desc: 'Diagonal wipe bottom-left', category: 'wipes', icon: 'i-heroicons-arrow-down-left' },
  { value: 'wipebr', label: 'Wipe Bottom-Right', desc: 'Diagonal wipe bottom-right', category: 'wipes', icon: 'i-heroicons-arrow-down-right' },
  // Slides & Pushes
  { value: 'slideleft', label: 'Slide Left', desc: 'Slide over from right', category: 'slides', icon: 'i-heroicons-arrow-left-start-on-rectangle' },
  { value: 'slideright', label: 'Slide Right', desc: 'Slide over from left', category: 'slides', icon: 'i-heroicons-arrow-right-start-on-rectangle' },
  { value: 'slideup', label: 'Slide Up', desc: 'Slide up from bottom', category: 'slides', icon: 'i-heroicons-arrow-up' },
  { value: 'slidedown', label: 'Slide Down', desc: 'Slide down from top', category: 'slides', icon: 'i-heroicons-arrow-down' },
  { value: 'smoothleft', label: 'Smooth Left', desc: 'Soft directional push left', category: 'slides', icon: 'i-heroicons-chevron-double-left' },
  { value: 'smoothright', label: 'Smooth Right', desc: 'Soft directional push right', category: 'slides', icon: 'i-heroicons-chevron-double-right' },
  { value: 'smoothup', label: 'Smooth Up', desc: 'Soft directional push up', category: 'slides', icon: 'i-heroicons-chevron-double-up' },
  { value: 'smoothdown', label: 'Smooth Down', desc: 'Soft directional push down', category: 'slides', icon: 'i-heroicons-chevron-double-down' },
  // Shapes & Iris
  { value: 'circleopen', label: 'Iris Circle Open', desc: 'Expanding circle', category: 'shapes', icon: 'i-heroicons-lifebuoy' },
  { value: 'circleclose', label: 'Iris Circle Close', desc: 'Closing circle mask', category: 'shapes', icon: 'i-heroicons-stop-circle' },
  { value: 'circlecrop', label: 'Circle Crop Iris', desc: 'Circular reveal', category: 'shapes', icon: 'i-heroicons-arrow-path' },
  { value: 'rectcrop', label: 'Rectangle Crop', desc: 'Expanding rectangle aperture', category: 'shapes', icon: 'i-heroicons-rectangle-stack' },
  { value: 'horzopen', label: 'Horizontal Open', desc: 'Doors open horizontally', category: 'shapes', icon: 'i-heroicons-arrows-pointing-out' },
  { value: 'horzclose', label: 'Horizontal Close', desc: 'Doors close horizontally', category: 'shapes', icon: 'i-heroicons-arrows-pointing-in' },
  { value: 'vertopen', label: 'Vertical Open', desc: 'Doors open vertically', category: 'shapes', icon: 'i-heroicons-arrows-up-down' },
  { value: 'vertclose', label: 'Vertical Close', desc: 'Doors close vertically', category: 'shapes', icon: 'i-heroicons-bars-2' },
  // Zooms & Warps
  { value: 'zoomin', label: 'Zoom In', desc: 'Dynamic scale burst zoom', category: 'zooms', icon: 'i-heroicons-magnifying-glass-plus' },
  { value: 'radial', label: 'Radial Clock Wipe', desc: 'Circular radar sweep', category: 'zooms', icon: 'i-heroicons-clock' },
  { value: 'pixelize', label: 'Pixelize / Glitch', desc: 'Mosaic block pixelation', category: 'zooms', icon: 'i-heroicons-squares-plus' },
  { value: 'squeezev', label: 'Squeeze Vertical', desc: 'Vertical compression', category: 'zooms', icon: 'i-heroicons-bars-3-bottom-left' },
  { value: 'squeezeh', label: 'Squeeze Horizontal', desc: 'Horizontal compression', category: 'zooms', icon: 'i-heroicons-bars-3' },
  { value: 'hlslice', label: 'Horizontal Slice', desc: 'Multi-slice blinds', category: 'zooms', icon: 'i-heroicons-table-cells' },
  { value: 'vuslice', label: 'Vertical Slice', desc: 'Vertical louvers blinds', category: 'zooms', icon: 'i-heroicons-queue-list' },
  { value: 'distance', label: 'Distance Fade', desc: '3D distance perspective', category: 'zooms', icon: 'i-heroicons-cube' },
]

const filteredTransitionOptions = computed(() => {
  return allTransitionsCatalog.filter((opt) => {
    // 1. Category Filter
    if (selectedCategory.value !== 'all') {
      if (opt.value !== 'none' && opt.category !== selectedCategory.value) {
        return false
      }
    }
    // 2. Search Query Filter
    if (transitionSearchQuery.value.trim()) {
      const q = transitionSearchQuery.value.toLowerCase().trim()
      return opt.label.toLowerCase().includes(q) || opt.value.toLowerCase().includes(q) || opt.desc.toLowerCase().includes(q)
    }
    return true
  })
})

const transitionIcon = computed(() => {
  const type = (props.transition?.type || selectedType.value || 'dissolve').toLowerCase()
  const found = allTransitionsCatalog.find((t) => t.value === type)
  return found?.icon || 'i-heroicons-sparkles'
})

function openModal() {
  if (props.transition) {
    selectedType.value = props.transition.type || 'dissolve'
    selectedDuration.value = Math.min(
      maxAllowedDuration.value,
      Math.max(0.1, Number(props.transition.duration) || 1.0)
    )
  } else {
    selectedType.value = 'dissolve'
    selectedDuration.value = Math.min(1.0, maxAllowedDuration.value)
  }
  isModalOpen.value = true
}

function applyTransition() {
  if (selectedType.value === 'none') {
    deleteTransition()
  } else {
    const finalDuration = Math.min(
      maxAllowedDuration.value,
      Math.max(0.1, Number(selectedDuration.value) || 1.0)
    )
    timelineStore.addTransition(
      props.trackId,
      props.fromClip.id,
      props.toClip.id,
      selectedType.value,
      Number(finalDuration.toFixed(2))
    )
  }
  isModalOpen.value = false
}

function deleteTransition() {
  const track = timelineStore.tracks.find((t) => t.id === props.trackId)
  if (track && track.transitions) {
    track.transitions = track.transitions.filter((t) => !(t.from === props.fromClip.id && t.to === props.toClip.id))
    timelineStore.pushHistoryState('Remove transition')
  }
  isModalOpen.value = false
}
</script>
