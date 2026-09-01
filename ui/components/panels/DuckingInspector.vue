<template>
  <div class="space-y-4 text-xs">
    <!-- Enable Ducking -->
    <div class="flex items-center justify-between">
      <span class="font-medium text-gray-300">{{ $t('inspector.ducking.enable') }}</span>
      <USwitch v-model="isEnabled" @update:model-value="toggleDucking" />
    </div>

    <div v-if="isEnabled" class="space-y-3 pt-2 border-t border-gray-800">
      <!-- Target Voiceover Track -->
      <div>
        <label class="block text-gray-400 mb-1">{{ $t('inspector.ducking.target_track') }}</label>
        <USelect
          v-model="targetTrack"
          :items="availableVoiceTracks"
          size="sm"
          @change="updateDucking"
        />
      </div>

      <!-- Threshold Slider -->
      <div>
        <div class="flex justify-between text-gray-400 mb-1">
          <span>{{ $t('inspector.ducking.threshold') }}</span>
          <span class="font-mono text-gray-300">{{ threshold }} dB</span>
        </div>
        <input
          v-model.number="threshold"
          type="range"
          min="-60"
          max="-5"
          step="1"
          class="w-full h-1 bg-gray-800 rounded appearance-none accent-indigo-500"
          @input="updateDucking"
        />
      </div>

      <!-- Attenuation / Duck Amount -->
      <div>
        <div class="flex justify-between text-gray-400 mb-1">
          <span>{{ $t('inspector.ducking.duck_amount') }}</span>
          <span class="font-mono text-gray-300">{{ duckAmount }} dB</span>
        </div>
        <input
          v-model.number="duckAmount"
          type="range"
          min="-40"
          max="-6"
          step="1"
          class="w-full h-1 bg-gray-800 rounded appearance-none accent-indigo-500"
          @input="updateDucking"
        />
      </div>

      <!-- Attack Time -->
      <div>
        <div class="flex justify-between text-gray-400 mb-1">
          <span>{{ $t('inspector.ducking.attack') }}</span>
          <span class="font-mono text-gray-300">{{ attackMs }} ms</span>
        </div>
        <input
          v-model.number="attackMs"
          type="range"
          min="10"
          max="500"
          step="10"
          class="w-full h-1 bg-gray-800 rounded appearance-none accent-indigo-500"
          @input="updateDucking"
        />
      </div>

      <!-- Release Time -->
      <div>
        <div class="flex justify-between text-gray-400 mb-1">
          <span>{{ $t('inspector.ducking.release') }}</span>
          <span class="font-mono text-gray-300">{{ releaseMs }} ms</span>
        </div>
        <input
          v-model.number="releaseMs"
          type="range"
          min="100"
          max="3000"
          step="50"
          class="w-full h-1 bg-gray-800 rounded appearance-none accent-indigo-500"
          @input="updateDucking"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useTimelineStore } from '~/stores/timeline'

const timelineStore = useTimelineStore()
const selectedTrack = computed(() => timelineStore.selectedTrack)

const isEnabled = ref(false)
const targetTrack = ref('')
const threshold = ref(-25)
const duckAmount = ref(-14)
const attackMs = ref(50)
const releaseMs = ref(500)

const availableVoiceTracks = computed(() => {
  return timelineStore.tracks
    .filter((t) => t.id !== selectedTrack.value?.id && (t.kind === 'audio' || t.kind === 'video'))
    .map((t) => ({ label: `${t.kind.toUpperCase()} (${t.id})`, value: t.id }))
})

watch(selectedTrack, (track) => {
  if (track?.ducking || track?.duck_under) {
    isEnabled.value = true
    targetTrack.value = track.duck_under || availableVoiceTracks.value[0]?.value || ''
    threshold.value = track.ducking?.threshold_db ?? -25
    duckAmount.value = track.ducking?.duck_db ?? -14
    attackMs.value = track.ducking?.attack_ms ?? 50
    releaseMs.value = track.ducking?.release_ms ?? 500
  } else {
    isEnabled.value = false
    targetTrack.value = availableVoiceTracks.value[0]?.value || ''
  }
}, { immediate: true })

function toggleDucking(enabled: boolean) {
  if (!selectedTrack.value) return
  if (enabled) {
    selectedTrack.value.duck_under = targetTrack.value
    selectedTrack.value.ducking = {
      threshold_db: threshold.value,
      duck_db: duckAmount.value,
      attack_ms: attackMs.value,
      release_ms: releaseMs.value,
    }
  } else {
    delete selectedTrack.value.duck_under
    delete selectedTrack.value.ducking
  }
  timelineStore.pushHistoryState('Toggle Sidechain Ducking')
}

function updateDucking() {
  if (!selectedTrack.value || !isEnabled.value) return
  selectedTrack.value.duck_under = targetTrack.value
  selectedTrack.value.ducking = {
    threshold_db: threshold.value,
    duck_db: duckAmount.value,
    attack_ms: attackMs.value,
    release_ms: releaseMs.value,
  }
}
</script>
