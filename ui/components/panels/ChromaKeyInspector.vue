<template>
  <div class="space-y-4 text-xs">
    <!-- Enable Chroma Key -->
    <div class="flex items-center justify-between">
      <span class="font-medium text-gray-300">{{ $t('inspector.chroma_key.enable') }}</span>
      <USwitch v-model="isEnabled" @update:model-value="toggleChromaKey" />
    </div>

    <div v-if="isEnabled" class="space-y-3 pt-2 border-t border-gray-800">
      <!-- Key Color Picker -->
      <div>
        <label class="block text-gray-400 mb-1">{{ $t('inspector.chroma_key.key_color') }}</label>
        <div class="flex items-center gap-2">
          <input
            v-model="keyColor"
            type="color"
            class="w-8 h-8 rounded border border-gray-700 bg-transparent cursor-pointer"
            @input="updateChromaKey"
          />
          <input
            v-model="keyColor"
            type="text"
            class="flex-1 bg-gray-950 border border-gray-800 rounded px-2 py-1 font-mono text-gray-300 focus:outline-none focus:border-indigo-500"
            @change="updateChromaKey"
          />
        </div>
      </div>

      <!-- Color Similarity Slider -->
      <div>
        <div class="flex justify-between text-gray-400 mb-1">
          <span>{{ $t('inspector.chroma_key.similarity') }}</span>
          <span class="font-mono text-gray-300">{{ similarity }}</span>
        </div>
        <input
          v-model.number="similarity"
          type="range"
          min="0.01"
          max="1.0"
          step="0.01"
          class="w-full h-1 bg-gray-800 rounded appearance-none accent-indigo-500"
          @input="updateChromaKey"
        />
      </div>

      <!-- Blend Edge Slider -->
      <div>
        <div class="flex justify-between text-gray-400 mb-1">
          <span>{{ $t('inspector.chroma_key.blend') }}</span>
          <span class="font-mono text-gray-300">{{ blend }}</span>
        </div>
        <input
          v-model.number="blend"
          type="range"
          min="0.0"
          max="0.5"
          step="0.01"
          class="w-full h-1 bg-gray-800 rounded appearance-none accent-indigo-500"
          @input="updateChromaKey"
        />
      </div>

      <!-- Despill Toggle -->
      <div class="flex items-center justify-between pt-1">
        <span class="text-gray-400">{{ $t('inspector.chroma_key.despill') }}</span>
        <USwitch v-model="despill" @update:model-value="updateChromaKey" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useTimelineStore } from '~/stores/timeline'

const timelineStore = useTimelineStore()
const selectedClip = computed(() => timelineStore.selectedClip)

const isEnabled = ref(false)
const keyColor = ref('#00FF00')
const similarity = ref(0.3)
const blend = ref(0.1)
const despill = ref(true)

watch(selectedClip, (clip) => {
  if (clip?.chroma_key) {
    isEnabled.value = true
    keyColor.value = clip.chroma_key.color || '#00FF00'
    similarity.value = clip.chroma_key.similarity ?? 0.3
    blend.value = clip.chroma_key.blend ?? 0.1
    despill.value = clip.chroma_key.despill ?? true
  } else {
    isEnabled.value = false
    keyColor.value = '#00FF00'
    similarity.value = 0.3
    blend.value = 0.1
    despill.value = true
  }
}, { immediate: true })

function toggleChromaKey(enabled: boolean) {
  if (!selectedClip.value) return
  if (enabled) {
    selectedClip.value.chroma_key = {
      color: keyColor.value,
      similarity: similarity.value,
      blend: blend.value,
      despill: despill.value,
    }
  } else {
    delete selectedClip.value.chroma_key
  }
  timelineStore.pushHistoryState('Toggle Chroma Key')
}

function updateChromaKey() {
  if (!selectedClip.value || !isEnabled.value) return
  selectedClip.value.chroma_key = {
    color: keyColor.value,
    similarity: similarity.value,
    blend: blend.value,
    despill: despill.value,
  }
}
</script>
