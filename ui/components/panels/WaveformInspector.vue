<template>
  <div class="space-y-4 text-xs">
    <!-- Visualizer Mode -->
    <div>
      <label class="block text-gray-400 mb-1">{{ $t('inspector.waveform.mode') }}</label>
      <USelect
        v-model="mode"
        :items="modeOptions"
        size="sm"
        @change="updateWaveform"
      />
    </div>

    <!-- Waveform Color -->
    <div>
      <label class="block text-gray-400 mb-1">{{ $t('inspector.waveform.color') }}</label>
      <div class="flex items-center gap-2">
        <input
          v-model="color"
          type="color"
          class="w-8 h-8 rounded border border-gray-700 bg-transparent cursor-pointer"
          @input="updateWaveform"
        />
        <input
          v-model="color"
          type="text"
          class="flex-1 bg-gray-950 border border-gray-800 rounded px-2 py-1 font-mono text-gray-300 focus:outline-none focus:border-indigo-500"
          @change="updateWaveform"
        />
      </div>
    </div>

    <!-- Height Slider -->
    <div>
      <div class="flex justify-between text-gray-400 mb-1">
        <span>{{ $t('inspector.waveform.height') }}</span>
        <span class="font-mono text-gray-300">{{ height }} px</span>
      </div>
      <input
        v-model.number="height"
        type="range"
        min="40"
        max="400"
        step="10"
        class="w-full h-1 bg-gray-800 rounded appearance-none accent-indigo-500"
        @input="updateWaveform"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { useTimelineStore } from '~/stores/timeline'

const { t } = useI18n()
const timelineStore = useTimelineStore()
const selectedClip = computed(() => timelineStore.selectedClip)

const mode = ref('peak_to_peak')
const color = ref('#6366F1')
const height = ref(120)

const modeOptions = computed(() => [
  { label: t('inspector.waveform.modes.peak_to_peak'), value: 'peak_to_peak' },
  { label: t('inspector.waveform.modes.bars'), value: 'bars' },
  { label: t('inspector.waveform.modes.circular'), value: 'circular' },
  { label: t('inspector.waveform.modes.wave'), value: 'wave' },
])

watch(selectedClip, (clip) => {
  if (clip?.waveform) {
    mode.value = clip.waveform.mode || 'peak_to_peak'
    color.value = clip.waveform.color || '#6366F1'
    height.value = clip.waveform.height || 120
  }
}, { immediate: true })

function updateWaveform() {
  if (!selectedClip.value) return
  selectedClip.value.waveform = {
    mode: mode.value,
    color: color.value,
    height: height.value,
  }
  timelineStore.pushHistoryState('Update Waveform Visualizer')
}
</script>
