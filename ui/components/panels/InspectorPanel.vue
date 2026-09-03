<template>
  <div class="h-full flex flex-col bg-gray-950 border-l border-gray-800 select-none overflow-hidden">
    <!-- Header -->
    <div class="h-10 border-b border-gray-800 px-4 flex items-center justify-between flex-shrink-0">
      <span class="text-xs font-semibold text-gray-200 tracking-wide uppercase">
        {{ selectedClip ? $t('inspector.title') : (selectedTrack ? 'Track Properties' : $t('inspector.title')) }}
      </span>
      <span v-if="selectedClip" class="text-[10px] font-mono text-indigo-400 bg-indigo-500/10 px-1.5 py-0.5 rounded border border-indigo-500/20">
        {{ selectedClip.id }}
      </span>
    </div>

    <!-- Inspector Body -->
    <div class="flex-1 overflow-y-auto p-4 space-y-6">
      <!-- Empty State -->
      <div v-if="!selectedClip && !selectedTrack" class="h-full flex flex-col items-center justify-center text-center text-gray-500 gap-2 py-16">
        <UIcon name="i-heroicons-adjustments-horizontal" class="w-8 h-8 opacity-40" />
        <p class="text-xs">{{ $t('inspector.empty') }}</p>
      </div>

      <!-- Clip Properties -->
      <template v-else-if="selectedClip">
        <!-- Section: Timing & Duration -->
        <div class="space-y-3">
          <h4 class="text-xs font-medium text-gray-400 uppercase tracking-wider">{{ $t('inspector.timing') }}</h4>
          <div class="grid grid-cols-2 gap-3 text-xs">
            <div>
              <label class="block text-gray-500 mb-1">{{ $t('inspector.start_time') }} (s)</label>
              <input
                v-model.number="selectedClip.start"
                type="number"
                step="0.1"
                min="0"
                class="w-full bg-gray-900 border border-gray-800 rounded px-2 py-1.5 font-mono text-gray-200 focus:outline-none focus:border-indigo-500"
                @change="commitChange('Update Start Time')"
              />
            </div>
            <div>
              <label class="block text-gray-500 mb-1">{{ $t('inspector.duration') }} (s)</label>
              <input
                v-model.number="selectedClip.duration"
                type="number"
                step="0.1"
                min="0.1"
                class="w-full bg-gray-900 border border-gray-800 rounded px-2 py-1.5 font-mono text-gray-200 focus:outline-none focus:border-indigo-500"
                @change="commitChange('Update Duration')"
              />
            </div>
          </div>
        </div>

        <!-- Section: Transform (Position, Scale, Rotation, Opacity) -->
        <div v-if="isVisualClip" class="space-y-3 pt-3 border-t border-gray-800/80">
          <h4 class="text-xs font-medium text-gray-400 uppercase tracking-wider">{{ $t('inspector.transform') }}</h4>

          <!-- Position X & Y -->
          <div class="grid grid-cols-2 gap-3 text-xs">
            <div>
              <label class="block text-gray-500 mb-1">Position X</label>
              <input
                v-model.number="clipPositionX"
                type="number"
                class="w-full bg-gray-900 border border-gray-800 rounded px-2 py-1.5 font-mono text-gray-200 focus:outline-none focus:border-indigo-500"
                @change="updatePosition"
              />
            </div>
            <div>
              <label class="block text-gray-500 mb-1">Position Y</label>
              <input
                v-model.number="clipPositionY"
                type="number"
                class="w-full bg-gray-900 border border-gray-800 rounded px-2 py-1.5 font-mono text-gray-200 focus:outline-none focus:border-indigo-500"
                @change="updatePosition"
              />
            </div>
          </div>

          <!-- Alignment Quick Select -->
          <div class="text-xs">
            <label class="block text-gray-500 mb-1">{{ $t('inspector.alignment') }}</label>
            <USelect
              v-model="selectedClip.alignment"
              :items="alignmentOptions"
              size="sm"
              @change="commitChange('Update Alignment')"
            />
          </div>

          <!-- Scale & Rotation -->
          <div class="grid grid-cols-2 gap-3 text-xs">
            <div>
              <div class="flex justify-between text-gray-500 mb-1">
                <span>{{ $t('inspector.scale') }}</span>
                <span class="font-mono text-gray-300">{{ selectedClip.scale ?? 1.0 }}x</span>
              </div>
              <input
                v-model.number="selectedClip.scale"
                type="range"
                min="0.1"
                max="3.0"
                step="0.05"
                class="w-full h-1 bg-gray-800 rounded appearance-none accent-indigo-500"
                @input="commitChange('Update Scale')"
              />
            </div>
            <div>
              <div class="flex justify-between text-gray-500 mb-1">
                <span>{{ $t('inspector.rotation') }}</span>
                <span class="font-mono text-gray-300">{{ selectedClip.rotation ?? 0 }}°</span>
              </div>
              <input
                v-model.number="selectedClip.rotation"
                type="range"
                min="-180"
                max="180"
                step="1"
                class="w-full h-1 bg-gray-800 rounded appearance-none accent-indigo-500"
                @input="commitChange('Update Rotation')"
              />
            </div>
          </div>

          <!-- Opacity Slider -->
          <div class="text-xs">
            <div class="flex justify-between text-gray-500 mb-1">
              <span>{{ $t('inspector.opacity') }}</span>
              <span class="font-mono text-gray-300">{{ Math.round((selectedClip.opacity ?? 1.0) * 100) }}%</span>
            </div>
            <input
              v-model.number="selectedClip.opacity"
              type="range"
              min="0"
              max="1"
              step="0.05"
              class="w-full h-1 bg-gray-800 rounded appearance-none accent-indigo-500"
              @input="commitChange('Update Opacity')"
            />
          </div>

          <!-- Blend Mode Select -->
          <div class="text-xs">
            <label class="block text-gray-500 mb-1">{{ $t('inspector.blend_mode') }}</label>
            <USelect
              v-model="clipBlendMode"
              :items="blendModeOptions"
              size="sm"
              @change="commitChange('Update Blend Mode')"
            />
          </div>
        </div>

        <!-- Section: Audio Volume & Pan -->
        <div v-if="isAudioClip" class="space-y-3 pt-3 border-t border-gray-800/80">
          <h4 class="text-xs font-medium text-gray-400 uppercase tracking-wider">{{ $t('inspector.audio') }}</h4>
          <div class="text-xs">
            <div class="flex justify-between text-gray-500 mb-1">
              <span>{{ $t('inspector.volume') }}</span>
              <span class="font-mono text-gray-300">{{ Math.round((selectedClip.volume ?? 1.0) * 100) }}%</span>
            </div>
            <input
              v-model.number="selectedClip.volume"
              type="range"
              min="0"
              max="2.0"
              step="0.05"
              class="w-full h-1 bg-gray-800 rounded appearance-none accent-indigo-500"
              @input="commitChange('Update Volume')"
            />
          </div>
        </div>

        <!-- Section: Chroma Keying Panel (Green Screen) -->
        <div v-if="isVisualClip" class="pt-3 border-t border-gray-800/80">
          <h4 class="text-xs font-medium text-gray-400 uppercase tracking-wider mb-2">{{ $t('inspector.chroma_key.title') }}</h4>
          <PanelsChromaKeyInspector />
        </div>

        <!-- Section: Waveform Visualizer Panel -->
        <div v-if="selectedClipTrackKind === 'waveform'" class="pt-3 border-t border-gray-800/80">
          <h4 class="text-xs font-medium text-gray-400 uppercase tracking-wider mb-2">{{ $t('inspector.waveform.title') }}</h4>
          <PanelsWaveformInspector />
        </div>

        <!-- Section: Subtitles Caption Panel -->
        <div v-if="selectedClipTrackKind === 'subtitle'" class="pt-3 border-t border-gray-800/80">
          <h4 class="text-xs font-medium text-gray-400 uppercase tracking-wider mb-2">{{ $t('inspector.subtitle.title') }}</h4>
          <PanelsSubtitleInspector />
        </div>

        <!-- Section: Ken Burns Keyframe Animation Panel -->
        <div v-if="isVisualClip" class="pt-3 border-t border-gray-800/80">
          <h4 class="text-xs font-medium text-gray-400 uppercase tracking-wider mb-2">{{ $t('inspector.motion.title') }}</h4>
          <PanelsKeyframeInspector />
        </div>
      </template>

      <!-- Track Ducking Properties (When track selected) -->
      <template v-else-if="selectedTrack">
        <div class="space-y-3">
          <h4 class="text-xs font-medium text-gray-400 uppercase tracking-wider">{{ $t('inspector.ducking.title') }}</h4>
          <PanelsDuckingInspector />
        </div>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useTimelineStore } from '~/stores/timeline'

const timelineStore = useTimelineStore()

const selectedClip = computed(() => timelineStore.selectedClip)
const selectedTrack = computed(() => timelineStore.selectedTrack)

const selectedClipTrackKind = computed(() => {
  if (!timelineStore.selectedClipId) return null
  for (const track of timelineStore.tracks) {
    if (track.clips?.some((c) => c.id === timelineStore.selectedClipId)) {
      return track.kind
    }
  }
  return null
})

const isVisualClip = computed(() => {
  const k = selectedClipTrackKind.value
  return k === 'video' || k === 'overlay' || k === 'waveform'
})

const isAudioClip = computed(() => {
  const k = selectedClipTrackKind.value
  return k === 'audio' || k === 'video'
})

const clipPositionX = computed({
  get: () => selectedClip.value?.position?.x || 0,
  set: (val) => {
    if (selectedClip.value) {
      if (!selectedClip.value.position) selectedClip.value.position = { x: 0, y: 0 }
      selectedClip.value.position.x = val
    }
  },
})

const clipPositionY = computed({
  get: () => selectedClip.value?.position?.y || 0,
  set: (val) => {
    if (selectedClip.value) {
      if (!selectedClip.value.position) selectedClip.value.position = { x: 0, y: 0 }
      selectedClip.value.position.y = val
    }
  },
})

const alignmentOptions = [
  { label: 'Center', value: 'center' },
  { label: 'Top Left', value: 'top-left' },
  { label: 'Top Right', value: 'top-right' },
  { label: 'Bottom Left', value: 'bottom-left' },
  { label: 'Bottom Right', value: 'bottom-right' },
]

const { t } = useI18n()

const clipBlendMode = computed({
  get: () => selectedClip.value?.blend_mode || 'normal',
  set: (val: string) => {
    if (selectedClip.value) {
      selectedClip.value.blend_mode = val
    }
  },
})

const blendModeOptions = computed(() => [
  { label: t('inspector.blend_modes.normal'), value: 'normal' },
  { label: t('inspector.blend_modes.multiply'), value: 'multiply' },
  { label: t('inspector.blend_modes.screen'), value: 'screen' },
  { label: t('inspector.blend_modes.overlay'), value: 'overlay' },
  { label: t('inspector.blend_modes.add'), value: 'add' },
  { label: t('inspector.blend_modes.darken'), value: 'darken' },
  { label: t('inspector.blend_modes.lighten'), value: 'lighten' },
  { label: t('inspector.blend_modes.difference'), value: 'difference' },
  { label: t('inspector.blend_modes.hardlight'), value: 'hardlight' },
  { label: t('inspector.blend_modes.softlight'), value: 'softlight' },
])

function updatePosition() {
  timelineStore.pushHistoryState('Update Clip Position')
}

function commitChange(label: string) {
  timelineStore.pushHistoryState(label)
}
</script>
