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

          <!-- Transition Types Grid -->
          <div class="space-y-2">
            <label class="block text-gray-400 font-medium uppercase tracking-wider text-[10px]">
              {{ $t('timeline.transitions.select_type') || 'Transition Type' }}
            </label>
            <div class="grid grid-cols-2 sm:grid-cols-3 gap-2">
              <button
                v-for="opt in transitionOptions"
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

const transitionIcon = computed(() => {
  const type = (props.transition?.type || selectedType.value || 'dissolve').toLowerCase()
  switch (type) {
    case 'fade':
      return 'i-heroicons-eye-slash'
    case 'wipeleft':
      return 'i-heroicons-arrow-left'
    case 'wiperight':
      return 'i-heroicons-arrow-right'
    case 'wipeup':
      return 'i-heroicons-arrow-up'
    case 'wipedown':
      return 'i-heroicons-arrow-down'
    case 'slideleft':
      return 'i-heroicons-arrow-left-start-on-rectangle'
    case 'slideright':
      return 'i-heroicons-arrow-right-start-on-rectangle'
    case 'circleopen':
    case 'circlecrop':
      return 'i-heroicons-lifebuoy'
    case 'dissolve':
    default:
      return 'i-heroicons-sparkles'
  }
})

const transitionOptions = computed(() => [
  {
    value: 'none',
    label: t('timeline.transitions.none') || 'None',
    desc: t('timeline.transitions.none_desc') || 'No transition (Cut)',
    icon: 'i-heroicons-no-symbol',
  },
  {
    value: 'dissolve',
    label: t('timeline.transitions.dissolve') || 'Dissolve',
    desc: 'Cross dissolve',
    icon: 'i-heroicons-sparkles',
  },
  {
    value: 'fade',
    label: t('timeline.transitions.fade') || 'Fade',
    desc: 'Dip to black',
    icon: 'i-heroicons-eye-slash',
  },
  {
    value: 'wipeleft',
    label: t('timeline.transitions.wipeleft') || 'Wipe Left',
    desc: 'Right to left',
    icon: 'i-heroicons-arrow-left',
  },
  {
    value: 'wiperight',
    label: t('timeline.transitions.wiperight') || 'Wipe Right',
    desc: 'Left to right',
    icon: 'i-heroicons-arrow-right',
  },
  {
    value: 'wipeup',
    label: t('timeline.transitions.wipeup') || 'Wipe Up',
    desc: 'Bottom to top',
    icon: 'i-heroicons-arrow-up',
  },
  {
    value: 'wipedown',
    label: t('timeline.transitions.wipedown') || 'Wipe Down',
    desc: 'Top to bottom',
    icon: 'i-heroicons-arrow-down',
  },
  {
    value: 'slideleft',
    label: t('timeline.transitions.slideleft') || 'Slide Left',
    desc: 'Slide left',
    icon: 'i-heroicons-arrow-left-start-on-rectangle',
  },
  {
    value: 'slideright',
    label: t('timeline.transitions.slideright') || 'Slide Right',
    desc: 'Slide right',
    icon: 'i-heroicons-arrow-right-start-on-rectangle',
  },
  {
    value: 'circleopen',
    label: t('timeline.transitions.circleopen') || 'Circle Open',
    desc: 'Iris circle open',
    icon: 'i-heroicons-lifebuoy',
  },
])

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
