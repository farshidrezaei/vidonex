<template>
  <div class="space-y-4 text-xs">
    <!-- Enable Ken Burns -->
    <div class="flex items-center justify-between">
      <span class="font-medium text-gray-300">{{ $t('inspector.motion.enable') }}</span>
      <USwitch v-model="isEnabled" @update:model-value="toggleKenBurns" />
    </div>

    <div v-if="isEnabled" class="space-y-3 pt-2 border-t border-gray-800">
      <!-- Start / End Scale -->
      <div class="grid grid-cols-2 gap-3">
        <div>
          <label class="block text-gray-400 mb-1">{{ $t('inspector.motion.start_scale') }}</label>
          <input
            v-model.number="startScale"
            type="number"
            step="0.05"
            min="0.1"
            max="5.0"
            class="w-full bg-gray-950 border border-gray-800 rounded px-2 py-1 font-mono text-gray-300"
            @change="updateKenBurns"
          />
        </div>

        <div>
          <label class="block text-gray-400 mb-1">{{ $t('inspector.motion.end_scale') }}</label>
          <input
            v-model.number="endScale"
            type="number"
            step="0.05"
            min="0.1"
            max="5.0"
            class="w-full bg-gray-950 border border-gray-800 rounded px-2 py-1 font-mono text-gray-300"
            @change="updateKenBurns"
          />
        </div>
      </div>

      <!-- Easing Curve -->
      <div>
        <label class="block text-gray-400 mb-1">{{ $t('inspector.motion.easing') }}</label>
        <USelect
          v-model="easing"
          :items="easingOptions"
          size="sm"
          @change="updateKenBurns"
        />
      </div>

      <!-- Visual Bezier Curve Editor -->
      <PanelsCurveEditor :easing="easing" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { useTimelineStore } from '~/stores/timeline'

const { t } = useI18n()
const timelineStore = useTimelineStore()
const selectedClip = computed(() => timelineStore.selectedClip)

const isEnabled = ref(false)
const startScale = ref(1.0)
const endScale = ref(1.25)
const easing = ref('easeInOutQuad')

const easingOptions = computed(() => [
  { label: t('inspector.motion.easings.linear'), value: 'linear' },
  { label: t('inspector.motion.easings.easeInQuad'), value: 'easeInQuad' },
  { label: t('inspector.motion.easings.easeOutQuad'), value: 'easeOutQuad' },
  { label: t('inspector.motion.easings.easeInOutQuad'), value: 'easeInOutQuad' },
  { label: t('inspector.motion.easings.easeInOutCubic'), value: 'easeInOutCubic' },
])

watch(selectedClip, (clip) => {
  if (clip?.ken_burns) {
    isEnabled.value = true
    startScale.value = clip.ken_burns.start_scale || 1.0
    endScale.value = clip.ken_burns.end_scale || 1.25
    easing.value = clip.ken_burns.easing || 'easeInOutQuad'
  } else {
    isEnabled.value = false
  }
}, { immediate: true })

function toggleKenBurns(enabled: boolean) {
  if (!selectedClip.value) return
  if (enabled) {
    selectedClip.value.ken_burns = {
      start_scale: startScale.value,
      end_scale: endScale.value,
      easing: easing.value,
    }
  } else {
    delete selectedClip.value.ken_burns
  }
  timelineStore.pushHistoryState('Toggle Ken Burns Animation')
}

function updateKenBurns() {
  if (!selectedClip.value || !isEnabled.value) return
  selectedClip.value.ken_burns = {
    start_scale: startScale.value,
    end_scale: endScale.value,
    easing: easing.value,
  }
}
</script>
