<template>
  <div v-if="selectedClip" class="space-y-3.5 text-xs">
    <!-- Header with Reset Button -->
    <div class="flex items-center justify-between">
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

    <!-- Brightness Slider -->
    <div class="space-y-1">
      <div class="flex justify-between text-gray-400 text-[11px]">
        <span>{{ $t('adjustments.brightness') || 'Brightness' }}</span>
        <span class="font-mono text-gray-300">{{ formatPercent(brightness, 0, 1) }}</span>
      </div>
      <input
        v-model.number="brightness"
        type="range"
        min="-0.5"
        max="0.5"
        step="0.01"
        class="w-full h-1 bg-gray-800 rounded appearance-none accent-indigo-500 cursor-pointer"
        @input="applyChange"
      />
    </div>

    <!-- Contrast Slider -->
    <div class="space-y-1">
      <div class="flex justify-between text-gray-400 text-[11px]">
        <span>{{ $t('adjustments.contrast') || 'Contrast' }}</span>
        <span class="font-mono text-gray-300">{{ formatMultiplier(contrast) }}</span>
      </div>
      <input
        v-model.number="contrast"
        type="range"
        min="0.2"
        max="2.5"
        step="0.05"
        class="w-full h-1 bg-gray-800 rounded appearance-none accent-indigo-500 cursor-pointer"
        @input="applyChange"
      />
    </div>

    <!-- Saturation Slider -->
    <div class="space-y-1">
      <div class="flex justify-between text-gray-400 text-[11px]">
        <span>{{ $t('adjustments.saturation') || 'Saturation' }}</span>
        <span class="font-mono text-gray-300">{{ formatMultiplier(saturation) }}</span>
      </div>
      <input
        v-model.number="saturation"
        type="range"
        min="0.0"
        max="2.5"
        step="0.05"
        class="w-full h-1 bg-gray-800 rounded appearance-none accent-indigo-500 cursor-pointer"
        @input="applyChange"
      />
    </div>

    <!-- Gamma Slider -->
    <div class="space-y-1">
      <div class="flex justify-between text-gray-400 text-[11px]">
        <span>{{ $t('adjustments.gamma') || 'Gamma' }}</span>
        <span class="font-mono text-gray-300">{{ formatMultiplier(gamma) }}</span>
      </div>
      <input
        v-model.number="gamma"
        type="range"
        min="0.2"
        max="2.5"
        step="0.05"
        class="w-full h-1 bg-gray-800 rounded appearance-none accent-indigo-500 cursor-pointer"
        @input="applyChange"
      />
    </div>

    <!-- Temperature (Warmth) Slider -->
    <div class="space-y-1">
      <div class="flex justify-between text-gray-400 text-[11px]">
        <span>{{ $t('adjustments.temperature') || 'Temperature (Warmth)' }}</span>
        <span class="font-mono" :class="temperature > 0 ? 'text-amber-400' : temperature < 0 ? 'text-cyan-400' : 'text-gray-300'">
          {{ formatTemperature(temperature) }}
        </span>
      </div>
      <input
        v-model.number="temperature"
        type="range"
        min="-0.8"
        max="0.8"
        step="0.02"
        class="w-full h-1 bg-gradient-to-r from-cyan-600 via-gray-700 to-amber-600 rounded appearance-none accent-amber-400 cursor-pointer"
        @input="applyChange"
      />
    </div>

    <!-- Tint (Green to Magenta) Slider -->
    <div class="space-y-1">
      <div class="flex justify-between text-gray-400 text-[11px]">
        <span>{{ $t('adjustments.tint') || 'Tint' }}</span>
        <span class="font-mono" :class="tint > 0 ? 'text-pink-400' : tint < 0 ? 'text-emerald-400' : 'text-gray-300'">
          {{ formatTint(tint) }}
        </span>
      </div>
      <input
        v-model.number="tint"
        type="range"
        min="-0.8"
        max="0.8"
        step="0.02"
        class="w-full h-1 bg-gradient-to-r from-emerald-600 via-gray-700 to-pink-600 rounded appearance-none accent-pink-400 cursor-pointer"
        @input="applyChange"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { useTimelineStore } from '~/stores/timeline'

const timelineStore = useTimelineStore()
const selectedClip = computed(() => timelineStore.selectedClip)

const brightness = computed({
  get: () => selectedClip.value?.color_grading?.brightness ?? 0.0,
  set: (val: number) => updateGrading({ brightness: val }),
})

const contrast = computed({
  get: () => selectedClip.value?.color_grading?.contrast ?? 1.0,
  set: (val: number) => updateGrading({ contrast: val }),
})

const saturation = computed({
  get: () => selectedClip.value?.color_grading?.saturation ?? 1.0,
  set: (val: number) => updateGrading({ saturation: val }),
})

const gamma = computed({
  get: () => selectedClip.value?.color_grading?.gamma ?? 1.0,
  set: (val: number) => updateGrading({ gamma: val }),
})

const temperature = computed({
  get: () => selectedClip.value?.color_grading?.temperature ?? 0.0,
  set: (val: number) => updateGrading({ temperature: val }),
})

const tint = computed({
  get: () => selectedClip.value?.color_grading?.tint ?? 0.0,
  set: (val: number) => updateGrading({ tint: val }),
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
  }
  applyChange()
}

function formatPercent(val: number, base: number, scale: number): string {
  const diff = Math.round((val - base) * scale * 100)
  return diff > 0 ? `+${diff}%` : `${diff}%`
}

function formatMultiplier(val: number): string {
  return `${val.toFixed(2)}x`
}

function formatTemperature(val: number): string {
  if (Math.abs(val) < 0.01) return '0 (Neutral)'
  return val > 0 ? `+${Math.round(val * 100)} (Warm)` : `${Math.round(val * 100)} (Cool)`
}

function formatTint(val: number): string {
  if (Math.abs(val) < 0.01) return '0 (Neutral)'
  return val > 0 ? `+${Math.round(val * 100)} (Magenta)` : `${Math.round(val * 100)} (Green)`
}
</script>
