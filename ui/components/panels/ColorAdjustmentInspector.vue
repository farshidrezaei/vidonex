<template>
  <div v-if="selectedClip" class="space-y-4 text-xs">
    <!-- Header with Reset Button -->
    <div class="flex items-center justify-between pb-1 border-b border-gray-800/60">
      <span class="text-[11px] font-medium text-gray-400 uppercase tracking-wider flex items-center gap-1.5">
        <UIcon name="i-heroicons-adjustments-horizontal" class="w-3.5 h-3.5 text-indigo-400" />
        {{ $t('adjustments.title') || 'Color & Light' }}
      </span>
      <button
        type="button"
        class="text-[10px] text-gray-500 hover:text-indigo-400 transition cursor-pointer"
        @click="resetAdjustments"
      >
        {{ $t('adjustments.reset') || 'Reset' }}
      </button>
    </div>

    <!-- Section 1: Basic Tone & Exposure (All Centered at 0) -->
    <div class="space-y-3">
      <!-- Brightness -->
      <div class="space-y-1">
        <div class="flex justify-between text-gray-400 text-[11px]">
          <span class="cursor-pointer select-none" @dblclick="brightness = 0">{{ $t('adjustments.brightness') || 'Brightness' }}</span>
          <span class="font-mono" :class="brightness !== 0 ? 'text-indigo-400' : 'text-gray-400'">{{ formatCentered(brightness) }}</span>
        </div>
        <div class="relative flex items-center">
          <div class="w-full h-1.5 bg-gray-800 rounded-full relative overflow-hidden pointer-events-none">
            <!-- Center zero marker -->
            <div class="absolute left-1/2 top-0 bottom-0 w-0.5 -translate-x-1/2 bg-gray-600"></div>
            <!-- Dynamic fill from center -->
            <div
              class="absolute top-0 bottom-0 bg-indigo-500 rounded-full transition-all duration-75"
              :style="getBipolarFillStyle(brightness)"
            ></div>
          </div>
          <input
            v-model.number="brightness"
            type="range"
            min="-100"
            max="100"
            step="1"
            class="absolute inset-0 w-full h-full opacity-0 cursor-pointer"
            @input="applyChange"
          />
        </div>
      </div>

      <!-- Contrast -->
      <div class="space-y-1">
        <div class="flex justify-between text-gray-400 text-[11px]">
          <span class="cursor-pointer select-none" @dblclick="contrast = 0">{{ $t('adjustments.contrast') || 'Contrast' }}</span>
          <span class="font-mono" :class="contrast !== 0 ? 'text-indigo-400' : 'text-gray-400'">{{ formatCentered(contrast) }}</span>
        </div>
        <div class="relative flex items-center">
          <div class="w-full h-1.5 bg-gray-800 rounded-full relative overflow-hidden pointer-events-none">
            <div class="absolute left-1/2 top-0 bottom-0 w-0.5 -translate-x-1/2 bg-gray-600"></div>
            <div
              class="absolute top-0 bottom-0 bg-indigo-500 rounded-full transition-all duration-75"
              :style="getBipolarFillStyle(contrast)"
            ></div>
          </div>
          <input
            v-model.number="contrast"
            type="range"
            min="-100"
            max="100"
            step="1"
            class="absolute inset-0 w-full h-full opacity-0 cursor-pointer"
            @input="applyChange"
          />
        </div>
      </div>

      <!-- Highlights -->
      <div class="space-y-1">
        <div class="flex justify-between text-gray-400 text-[11px]">
          <span class="cursor-pointer select-none" @dblclick="highlights = 0">{{ $t('adjustments.highlights') || 'Highlights' }}</span>
          <span class="font-mono" :class="highlights !== 0 ? 'text-indigo-400' : 'text-gray-400'">{{ formatCentered(highlights) }}</span>
        </div>
        <div class="relative flex items-center">
          <div class="w-full h-1.5 bg-gray-800 rounded-full relative overflow-hidden pointer-events-none">
            <div class="absolute left-1/2 top-0 bottom-0 w-0.5 -translate-x-1/2 bg-gray-600"></div>
            <div
              class="absolute top-0 bottom-0 bg-indigo-500 rounded-full transition-all duration-75"
              :style="getBipolarFillStyle(highlights)"
            ></div>
          </div>
          <input
            v-model.number="highlights"
            type="range"
            min="-100"
            max="100"
            step="1"
            class="absolute inset-0 w-full h-full opacity-0 cursor-pointer"
            @input="applyChange"
          />
        </div>
      </div>

      <!-- Shadows -->
      <div class="space-y-1">
        <div class="flex justify-between text-gray-400 text-[11px]">
          <span class="cursor-pointer select-none" @dblclick="shadows = 0">{{ $t('adjustments.shadows') || 'Shadows' }}</span>
          <span class="font-mono" :class="shadows !== 0 ? 'text-indigo-400' : 'text-gray-400'">{{ formatCentered(shadows) }}</span>
        </div>
        <div class="relative flex items-center">
          <div class="w-full h-1.5 bg-gray-800 rounded-full relative overflow-hidden pointer-events-none">
            <div class="absolute left-1/2 top-0 bottom-0 w-0.5 -translate-x-1/2 bg-gray-600"></div>
            <div
              class="absolute top-0 bottom-0 bg-indigo-500 rounded-full transition-all duration-75"
              :style="getBipolarFillStyle(shadows)"
            ></div>
          </div>
          <input
            v-model.number="shadows"
            type="range"
            min="-100"
            max="100"
            step="1"
            class="absolute inset-0 w-full h-full opacity-0 cursor-pointer"
            @input="applyChange"
          />
        </div>
      </div>

      <!-- Whites -->
      <div class="space-y-1">
        <div class="flex justify-between text-gray-400 text-[11px]">
          <span class="cursor-pointer select-none" @dblclick="whites = 0">{{ $t('adjustments.whites') || 'Whites' }}</span>
          <span class="font-mono" :class="whites !== 0 ? 'text-indigo-400' : 'text-gray-400'">{{ formatCentered(whites) }}</span>
        </div>
        <div class="relative flex items-center">
          <div class="w-full h-1.5 bg-gray-800 rounded-full relative overflow-hidden pointer-events-none">
            <div class="absolute left-1/2 top-0 bottom-0 w-0.5 -translate-x-1/2 bg-gray-600"></div>
            <div
              class="absolute top-0 bottom-0 bg-indigo-500 rounded-full transition-all duration-75"
              :style="getBipolarFillStyle(whites)"
            ></div>
          </div>
          <input
            v-model.number="whites"
            type="range"
            min="-100"
            max="100"
            step="1"
            class="absolute inset-0 w-full h-full opacity-0 cursor-pointer"
            @input="applyChange"
          />
        </div>
      </div>

      <!-- Blacks -->
      <div class="space-y-1">
        <div class="flex justify-between text-gray-400 text-[11px]">
          <span class="cursor-pointer select-none" @dblclick="blacks = 0">{{ $t('adjustments.blacks') || 'Blacks' }}</span>
          <span class="font-mono" :class="blacks !== 0 ? 'text-indigo-400' : 'text-gray-400'">{{ formatCentered(blacks) }}</span>
        </div>
        <div class="relative flex items-center">
          <div class="w-full h-1.5 bg-gray-800 rounded-full relative overflow-hidden pointer-events-none">
            <div class="absolute left-1/2 top-0 bottom-0 w-0.5 -translate-x-1/2 bg-gray-600"></div>
            <div
              class="absolute top-0 bottom-0 bg-indigo-500 rounded-full transition-all duration-75"
              :style="getBipolarFillStyle(blacks)"
            ></div>
          </div>
          <input
            v-model.number="blacks"
            type="range"
            min="-100"
            max="100"
            step="1"
            class="absolute inset-0 w-full h-full opacity-0 cursor-pointer"
            @input="applyChange"
          />
        </div>
      </div>

      <!-- Saturation -->
      <div class="space-y-1">
        <div class="flex justify-between text-gray-400 text-[11px]">
          <span class="cursor-pointer select-none" @dblclick="saturation = 0">{{ $t('adjustments.saturation') || 'Saturation' }}</span>
          <span class="font-mono" :class="saturation !== 0 ? 'text-indigo-400' : 'text-gray-400'">{{ formatCentered(saturation) }}</span>
        </div>
        <div class="relative flex items-center">
          <div class="w-full h-1.5 bg-gray-800 rounded-full relative overflow-hidden pointer-events-none">
            <div class="absolute left-1/2 top-0 bottom-0 w-0.5 -translate-x-1/2 bg-gray-600"></div>
            <div
              class="absolute top-0 bottom-0 bg-indigo-500 rounded-full transition-all duration-75"
              :style="getBipolarFillStyle(saturation)"
            ></div>
          </div>
          <input
            v-model.number="saturation"
            type="range"
            min="-100"
            max="100"
            step="1"
            class="absolute inset-0 w-full h-full opacity-0 cursor-pointer"
            @input="applyChange"
          />
        </div>
      </div>

      <!-- Gamma -->
      <div class="space-y-1">
        <div class="flex justify-between text-gray-400 text-[11px]">
          <span class="cursor-pointer select-none" @dblclick="gamma = 0">{{ $t('adjustments.gamma') || 'Gamma' }}</span>
          <span class="font-mono" :class="gamma !== 0 ? 'text-indigo-400' : 'text-gray-400'">{{ formatCentered(gamma) }}</span>
        </div>
        <div class="relative flex items-center">
          <div class="w-full h-1.5 bg-gray-800 rounded-full relative overflow-hidden pointer-events-none">
            <div class="absolute left-1/2 top-0 bottom-0 w-0.5 -translate-x-1/2 bg-gray-600"></div>
            <div
              class="absolute top-0 bottom-0 bg-indigo-500 rounded-full transition-all duration-75"
              :style="getBipolarFillStyle(gamma)"
            ></div>
          </div>
          <input
            v-model.number="gamma"
            type="range"
            min="-100"
            max="100"
            step="1"
            class="absolute inset-0 w-full h-full opacity-0 cursor-pointer"
            @input="applyChange"
          />
        </div>
      </div>
    </div>

    <!-- Section 2: Color Temperature & Tint (Centered at 0 with Gradient) -->
    <div class="space-y-3 pt-2 border-t border-gray-800/60">
      <!-- Temperature -->
      <div class="space-y-1">
        <div class="flex justify-between text-gray-400 text-[11px]">
          <span class="cursor-pointer select-none" @dblclick="temperature = 0">{{ $t('adjustments.temperature') || 'Temperature' }}</span>
          <span class="font-mono" :class="temperature > 0 ? 'text-amber-400' : temperature < 0 ? 'text-cyan-400' : 'text-gray-400'">
            {{ formatCentered(temperature) }}
          </span>
        </div>
        <div class="relative flex items-center">
          <div class="w-full h-1.5 bg-gradient-to-r from-cyan-600 via-gray-700 to-amber-600 rounded-full relative overflow-hidden pointer-events-none">
            <div class="absolute left-1/2 top-0 bottom-0 w-0.5 -translate-x-1/2 bg-white/60"></div>
            <div
              class="absolute top-0 bottom-0 rounded-full transition-all duration-75"
              :class="temperature > 0 ? 'bg-amber-400/80' : 'bg-cyan-400/80'"
              :style="getBipolarFillStyle(temperature)"
            ></div>
          </div>
          <input
            v-model.number="temperature"
            type="range"
            min="-100"
            max="100"
            step="1"
            class="absolute inset-0 w-full h-full opacity-0 cursor-pointer"
            @input="applyChange"
          />
        </div>
      </div>

      <!-- Tint -->
      <div class="space-y-1">
        <div class="flex justify-between text-gray-400 text-[11px]">
          <span class="cursor-pointer select-none" @dblclick="tint = 0">{{ $t('adjustments.tint') || 'Tint' }}</span>
          <span class="font-mono" :class="tint > 0 ? 'text-pink-400' : tint < 0 ? 'text-emerald-400' : 'text-gray-400'">
            {{ formatCentered(tint) }}
          </span>
        </div>
        <div class="relative flex items-center">
          <div class="w-full h-1.5 bg-gradient-to-r from-emerald-600 via-gray-700 to-pink-600 rounded-full relative overflow-hidden pointer-events-none">
            <div class="absolute left-1/2 top-0 bottom-0 w-0.5 -translate-x-1/2 bg-white/60"></div>
            <div
              class="absolute top-0 bottom-0 rounded-full transition-all duration-75"
              :class="tint > 0 ? 'bg-pink-400/80' : 'bg-emerald-400/80'"
              :style="getBipolarFillStyle(tint)"
            ></div>
          </div>
          <input
            v-model.number="tint"
            type="range"
            min="-100"
            max="100"
            step="1"
            class="absolute inset-0 w-full h-full opacity-0 cursor-pointer"
            @input="applyChange"
          />
        </div>
      </div>
    </div>

    <!-- Section 3: Blur & Sharpen Detail Filters (0 to 100) -->
    <div class="space-y-3 pt-2 border-t border-gray-800/60">
      <!-- Blur -->
      <div class="space-y-1">
        <div class="flex justify-between text-gray-400 text-[11px]">
          <span class="cursor-pointer select-none" @dblclick="blurVal = 0">{{ $t('adjustments.blur') || 'Blur' }}</span>
          <span class="font-mono" :class="blurVal > 0 ? 'text-indigo-400' : 'text-gray-400'">{{ blurVal }}</span>
        </div>
        <div class="relative flex items-center">
          <div class="w-full h-1.5 bg-gray-800 rounded-full relative overflow-hidden pointer-events-none">
            <div
              class="h-full bg-indigo-500 rounded-full transition-all duration-75"
              :style="{ width: `${blurVal}%` }"
            ></div>
          </div>
          <input
            v-model.number="blurVal"
            type="range"
            min="0"
            max="100"
            step="1"
            class="absolute inset-0 w-full h-full opacity-0 cursor-pointer"
            @input="applyChange"
          />
        </div>
      </div>

      <!-- Sharpen -->
      <div class="space-y-1">
        <div class="flex justify-between text-gray-400 text-[11px]">
          <span class="cursor-pointer select-none" @dblclick="sharpenVal = 0">{{ $t('adjustments.sharpen') || 'Sharpen' }}</span>
          <span class="font-mono" :class="sharpenVal > 0 ? 'text-indigo-400' : 'text-gray-400'">{{ sharpenVal }}</span>
        </div>
        <div class="relative flex items-center">
          <div class="w-full h-1.5 bg-gray-800 rounded-full relative overflow-hidden pointer-events-none">
            <div
              class="h-full bg-indigo-500 rounded-full transition-all duration-75"
              :style="{ width: `${sharpenVal}%` }"
            ></div>
          </div>
          <input
            v-model.number="sharpenVal"
            type="range"
            min="0"
            max="100"
            step="1"
            class="absolute inset-0 w-full h-full opacity-0 cursor-pointer"
            @input="applyChange"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useTimelineStore } from '~/stores/timeline'

const timelineStore = useTimelineStore()
const selectedClip = computed(() => timelineStore.selectedClip)

// Helper to style dynamic fill from center 50%
function getBipolarFillStyle(val: number) {
  if (val > 0) {
    const w = (val / 200) * 100
    return { left: '50%', width: `${w}%` }
  } else if (val < 0) {
    const w = (-val / 200) * 100
    return { left: `${50 - w}%`, width: `${w}%` }
  }
  return { left: '50%', width: '0%' }
}

function formatCentered(val: number): string {
  if (val === 0) return '0'
  return val > 0 ? `+${val}` : `${val}`
}

// Bipolar mapped computed properties (-100 to +100)

// Brightness: -100 to +100 maps to -0.5 to +0.5
const brightness = computed({
  get: () => Math.round(((selectedClip.value?.color_grading?.brightness ?? 0.0) / 0.5) * 100),
  set: (val: number) => updateGrading({ brightness: Number((val * 0.005).toFixed(3)) }),
})

// Contrast: -100 to +100 maps to 0.4 to 1.8 (0 is 1.0)
const contrast = computed({
  get: () => {
    const raw = selectedClip.value?.color_grading?.contrast ?? 1.0
    return Math.round((raw - 1.0) * 125)
  },
  set: (val: number) => updateGrading({ contrast: Number((1.0 + val * 0.008).toFixed(3)) }),
})

// Highlights: -100 to +100 maps to -1.0 to 1.0
const highlights = computed({
  get: () => Math.round((selectedClip.value?.color_grading?.highlights ?? 0.0) * 100),
  set: (val: number) => updateGrading({ highlights: Number((val * 0.01).toFixed(3)) }),
})

// Shadows: -100 to +100 maps to -1.0 to 1.0
const shadows = computed({
  get: () => Math.round((selectedClip.value?.color_grading?.shadows ?? 0.0) * 100),
  set: (val: number) => updateGrading({ shadows: Number((val * 0.01).toFixed(3)) }),
})

// Whites: -100 to +100 maps to -1.0 to 1.0
const whites = computed({
  get: () => Math.round((selectedClip.value?.color_grading?.whites ?? 0.0) * 100),
  set: (val: number) => updateGrading({ whites: Number((val * 0.01).toFixed(3)) }),
})

// Blacks: -100 to +100 maps to -1.0 to 1.0
const blacks = computed({
  get: () => Math.round((selectedClip.value?.color_grading?.blacks ?? 0.0) * 100),
  set: (val: number) => updateGrading({ blacks: Number((val * 0.01).toFixed(3)) }),
})

// Saturation: -100 to +100 maps to 0.0 to 2.0 (0 is 1.0)
const saturation = computed({
  get: () => {
    const raw = selectedClip.value?.color_grading?.saturation ?? 1.0
    return Math.round((raw - 1.0) * 100)
  },
  set: (val: number) => updateGrading({ saturation: Number((1.0 + val * 0.01).toFixed(3)) }),
})

// Gamma: -100 to +100 maps to 0.5 to 1.75 (0 is 1.0)
const gamma = computed({
  get: () => {
    const raw = selectedClip.value?.color_grading?.gamma ?? 1.0
    return Math.round((raw - 1.0) * 133.3)
  },
  set: (val: number) => updateGrading({ gamma: Number((1.0 + val * 0.0075).toFixed(3)) }),
})

// Temperature: -100 to +100 maps to -0.8 to +0.8
const temperature = computed({
  get: () => Math.round(((selectedClip.value?.color_grading?.temperature ?? 0.0) / 0.8) * 100),
  set: (val: number) => updateGrading({ temperature: Number((val * 0.008).toFixed(3)) }),
})

// Tint: -100 to +100 maps to -0.8 to +0.8
const tint = computed({
  get: () => Math.round(((selectedClip.value?.color_grading?.tint ?? 0.0) / 0.8) * 100),
  set: (val: number) => updateGrading({ tint: Number((val * 0.008).toFixed(3)) }),
})

// Positive only: Blur (0 to 100)
const blurVal = computed({
  get: () => Math.round(selectedClip.value?.color_grading?.blur ?? 0),
  set: (val: number) => updateGrading({ blur: val }),
})

// Positive only: Sharpen (0 to 100)
const sharpenVal = computed({
  get: () => Math.round(selectedClip.value?.color_grading?.sharpen ?? 0),
  set: (val: number) => updateGrading({ sharpen: val }),
})

function updateGrading(patch: Record<string, any>) {
  if (!selectedClip.value) return
  if (!selectedClip.value.color_grading) {
    selectedClip.value.color_grading = {
      brightness: 0.0,
      contrast: 1.0,
      saturation: 1.0,
      gamma: 1.0,
      temperature: 0.0,
      tint: 0.0,
      highlights: 0.0,
      shadows: 0.0,
      whites: 0.0,
      blacks: 0.0,
      blur: 0,
      sharpen: 0,
    }
  }
  Object.assign(selectedClip.value.color_grading, patch)
}

function applyChange() {
  timelineStore.pushHistoryState('Update Color Adjustments')
}

function resetAdjustments() {
  if (!selectedClip.value) return
  selectedClip.value.color_grading = {
    brightness: 0.0,
    contrast: 1.0,
    saturation: 1.0,
    gamma: 1.0,
    temperature: 0.0,
    tint: 0.0,
    highlights: 0.0,
    shadows: 0.0,
    whites: 0.0,
    blacks: 0.0,
    blur: 0,
    sharpen: 0,
  }
  applyChange()
}
</script>
