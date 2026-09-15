<template>
  <div v-if="activeTarget" class="space-y-4 text-xs">
    <!-- Header with Reset Button -->
    <div class="flex items-center justify-between pb-1 border-b border-gray-800/60">
      <span class="text-[11px] font-medium text-gray-400 uppercase tracking-wider flex items-center gap-1.5">
        <UIcon name="i-heroicons-sparkles" class="w-3.5 h-3.5 text-indigo-400" />
        {{ $t('inspector.waveform.title') }}
      </span>
      <button
        type="button"
        class="text-[10px] text-gray-500 hover:text-indigo-400 transition cursor-pointer"
        @click="resetToDefault"
      >
        {{ $t('adjustments.reset') || 'Reset' }}
      </button>
    </div>

    <!-- Section 1: Visual Mode Selector (5 Modes) -->
    <div class="space-y-1.5">
      <label class="block text-gray-400 text-[11px] font-medium">{{ $t('inspector.waveform.mode') }}</label>
      <div class="grid grid-cols-5 gap-1 bg-gray-900/90 p-1 rounded-lg border border-gray-800">
        <button
          v-for="m in modesList"
          :key="m.id"
          type="button"
          class="flex flex-col items-center justify-center p-1.5 rounded transition cursor-pointer text-center group"
          :class="mode === m.id ? 'bg-indigo-600 text-white shadow-sm' : 'text-gray-400 hover:text-gray-200 hover:bg-gray-800/60'"
          :title="m.label"
          @click="selectMode(m.id)"
        >
          <UIcon :name="m.icon" class="w-4 h-4 mb-0.5" />
          <span class="text-[9px] truncate max-w-full font-sans leading-tight">{{ m.shortLabel }}</span>
        </button>
      </div>
    </div>

    <!-- Section 2: Quick Style Presets -->
    <div class="space-y-1.5">
      <label class="block text-gray-400 text-[11px] font-medium">{{ $t('inspector.waveform.presets_title') }}</label>
      <div class="grid grid-cols-2 gap-1.5">
        <button
          v-for="preset in stylePresets"
          :key="preset.id"
          type="button"
          class="flex items-center gap-2 p-1.5 rounded border border-gray-800 bg-gray-900/60 hover:bg-gray-800/80 hover:border-gray-700 transition cursor-pointer text-left"
          @click="applyPreset(preset)"
        >
          <div class="w-4 h-4 rounded-full flex-shrink-0 shadow-inner flex items-center justify-center overflow-hidden" :style="{ background: `linear-gradient(135deg, ${preset.color}, ${preset.secondaryColor})` }"></div>
          <div class="flex-1 min-w-0">
            <div class="text-[10px] font-medium text-gray-300 truncate">{{ preset.name }}</div>
          </div>
        </button>
      </div>
    </div>

    <!-- Section 3: Dimensions & Alignment -->
    <div class="space-y-3 pt-2 border-t border-gray-800/60">
      <div class="flex items-center justify-between text-gray-400 text-[11px] font-medium">
        <span>{{ $t('inspector.waveform.dimensions') }}</span>
        <div class="flex items-center gap-1">
          <!-- Quick Alignment Buttons -->
          <button
            class="px-1.5 py-0.5 rounded bg-gray-800 hover:bg-gray-700 text-gray-300 text-[10px] transition"
            title="Center Canvas"
            @click="alignCenter"
          >
            Center
          </button>
          <button
            class="px-1.5 py-0.5 rounded bg-gray-800 hover:bg-gray-700 text-gray-300 text-[10px] transition"
            title="Lower Third"
            @click="alignLowerThird"
          >
            Lower ⅓
          </button>
        </div>
      </div>

      <!-- Width Slider -->
      <div class="space-y-1">
        <div class="flex justify-between text-gray-400 text-[11px]">
          <span>{{ $t('inspector.waveform.width') }}</span>
          <span class="font-mono text-gray-300">{{ width }} px</span>
        </div>
        <input
          v-model.number="width"
          type="range"
          min="150"
          :max="projectStore.canvasWidth || 1920"
          step="10"
          class="w-full h-1 bg-gray-800 rounded appearance-none accent-indigo-500 cursor-pointer"
          @input="commitChanges"
        />
      </div>

      <!-- Height Slider -->
      <div class="space-y-1">
        <div class="flex justify-between text-gray-400 text-[11px]">
          <span>{{ $t('inspector.waveform.height') }}</span>
          <span class="font-mono text-gray-300">{{ height }} px</span>
        </div>
        <input
          v-model.number="height"
          type="range"
          min="40"
          max="600"
          step="5"
          class="w-full h-1 bg-gray-800 rounded appearance-none accent-indigo-500 cursor-pointer"
          @input="commitChanges"
        />
      </div>

      <!-- Position X (Centered at 0) -->
      <div class="space-y-1">
        <div class="flex justify-between text-gray-400 text-[11px]">
          <span class="cursor-pointer select-none" @dblclick="positionX = 0">{{ $t('inspector.waveform.position_x') }}</span>
          <span class="font-mono" :class="positionX !== 0 ? 'text-indigo-400' : 'text-gray-400'">{{ positionX === 0 ? '0 (Center)' : (positionX > 0 ? `+${positionX}` : positionX) }}</span>
        </div>
        <input
          v-model.number="positionX"
          type="range"
          :min="-(projectStore.canvasWidth || 1920) / 2"
          :max="(projectStore.canvasWidth || 1920) / 2"
          step="5"
          class="w-full h-1 bg-gray-800 rounded appearance-none accent-indigo-500 cursor-pointer"
          @input="commitChanges"
        />
      </div>

      <!-- Position Y (Centered at 0) -->
      <div class="space-y-1">
        <div class="flex justify-between text-gray-400 text-[11px]">
          <span class="cursor-pointer select-none" @dblclick="positionY = 0">{{ $t('inspector.waveform.position_y') }}</span>
          <span class="font-mono" :class="positionY !== 0 ? 'text-indigo-400' : 'text-gray-400'">{{ positionY === 0 ? '0 (Center)' : (positionY > 0 ? `+${positionY}` : positionY) }}</span>
        </div>
        <input
          v-model.number="positionY"
          type="range"
          :min="-(projectStore.canvasHeight || 1080) / 2"
          :max="(projectStore.canvasHeight || 1080) / 2"
          step="5"
          class="w-full h-1 bg-gray-800 rounded appearance-none accent-indigo-500 cursor-pointer"
          @input="commitChanges"
        />
      </div>
    </div>

    <!-- Section 4: Dual Colors & Neon Glow -->
    <div class="space-y-3 pt-2 border-t border-gray-800/60">
      <h4 class="text-gray-400 text-[11px] font-medium uppercase tracking-wider">Colors & Styling</h4>

      <!-- Primary and Secondary Color -->
      <div class="grid grid-cols-2 gap-3">
        <!-- Primary Color -->
        <div class="space-y-1">
          <label class="block text-gray-400 text-[11px]">{{ $t('inspector.waveform.color') }}</label>
          <div class="flex items-center gap-1.5 bg-gray-900 border border-gray-800 rounded px-2 py-1">
            <input
              v-model="color"
              type="color"
              class="w-5 h-5 rounded border-0 bg-transparent cursor-pointer p-0"
              @input="commitChanges"
            />
            <input
              v-model="color"
              type="text"
              class="flex-1 bg-transparent font-mono text-[11px] text-gray-200 focus:outline-none uppercase"
              @change="commitChanges"
            />
          </div>
        </div>

        <!-- Secondary Color -->
        <div class="space-y-1">
          <label class="block text-gray-400 text-[11px]">{{ $t('inspector.waveform.secondary_color') }}</label>
          <div class="flex items-center gap-1.5 bg-gray-900 border border-gray-800 rounded px-2 py-1">
            <input
              v-model="secondaryColor"
              type="color"
              class="w-5 h-5 rounded border-0 bg-transparent cursor-pointer p-0"
              @input="commitChanges"
            />
            <input
              v-model="secondaryColor"
              type="text"
              class="flex-1 bg-transparent font-mono text-[11px] text-gray-200 focus:outline-none uppercase"
              @change="commitChanges"
            />
          </div>
        </div>
      </div>

      <!-- Glow Toggle & Intensity -->
      <div class="flex items-center justify-between">
        <label class="text-gray-400 text-[11px] flex items-center gap-1.5 cursor-pointer">
          <input
            v-model="isGlow"
            type="checkbox"
            class="rounded border-gray-700 bg-gray-900 text-indigo-600 focus:ring-indigo-500"
            @change="commitChanges"
          />
          <span>{{ $t('inspector.waveform.glow') }}</span>
        </label>
        <span v-if="isGlow" class="font-mono text-indigo-400 text-[10px]">{{ glowRadius }}px</span>
      </div>
      <div v-if="isGlow" class="space-y-1">
        <input
          v-model.number="glowRadius"
          type="range"
          min="2"
          max="24"
          step="1"
          class="w-full h-1 bg-gray-800 rounded appearance-none accent-indigo-500 cursor-pointer"
          @input="commitChanges"
        />
      </div>

      <!-- Opacity Slider -->
      <div class="space-y-1">
        <div class="flex justify-between text-gray-400 text-[11px]">
          <span>Opacity</span>
          <span class="font-mono text-gray-300">{{ Math.round(opacity * 100) }}%</span>
        </div>
        <input
          v-model.number="opacity"
          type="range"
          min="0.1"
          max="1.0"
          step="0.05"
          class="w-full h-1 bg-gray-800 rounded appearance-none accent-indigo-500 cursor-pointer"
          @input="commitChanges"
        />
      </div>
    </div>

    <!-- Section 5: Mode-Specific Customizations -->
    <div class="space-y-3 pt-2 border-t border-gray-800/60">
      <h4 class="text-gray-400 text-[11px] font-medium uppercase tracking-wider">Mode Fine-Tuning</h4>

      <!-- If Bars / Spectrum / Dots: Density Slider -->
      <div v-if="mode === 'peak_to_peak' || mode === 'bars' || mode === 'spectrum' || mode === 'dots' || mode === 'circular'" class="space-y-1">
        <div class="flex justify-between text-gray-400 text-[11px]">
          <span>{{ $t('inspector.waveform.density') }}</span>
          <span class="font-mono text-gray-300">{{ density }}</span>
        </div>
        <input
          v-model.number="density"
          type="range"
          min="16"
          max="96"
          step="4"
          class="w-full h-1 bg-gray-800 rounded appearance-none accent-indigo-500 cursor-pointer"
          @input="commitChanges"
        />
      </div>

      <!-- If Bars or Spectrum: Corner Roundness -->
      <div v-if="mode === 'peak_to_peak' || mode === 'bars' || mode === 'spectrum'" class="space-y-1">
        <div class="flex justify-between text-gray-400 text-[11px]">
          <span>{{ $t('inspector.waveform.roundness') }}</span>
          <span class="font-mono text-gray-300">{{ roundness }} px</span>
        </div>
        <input
          v-model.number="roundness"
          type="range"
          min="0"
          max="12"
          step="1"
          class="w-full h-1 bg-gray-800 rounded appearance-none accent-indigo-500 cursor-pointer"
          @input="commitChanges"
        />
      </div>

      <!-- If Wave: Line Width & Area Fill -->
      <div v-if="mode === 'wave' || mode === 'line'" class="space-y-2">
        <div class="space-y-1">
          <div class="flex justify-between text-gray-400 text-[11px]">
            <span>{{ $t('inspector.waveform.line_width') }}</span>
            <span class="font-mono text-gray-300">{{ lineWidth }} px</span>
          </div>
          <input
            v-model.number="lineWidth"
            type="range"
            min="1"
            max="8"
            step="1"
            class="w-full h-1 bg-gray-800 rounded appearance-none accent-indigo-500 cursor-pointer"
            @input="commitChanges"
          />
        </div>

        <label class="text-gray-400 text-[11px] flex items-center gap-1.5 cursor-pointer">
          <input
            v-model="isFillArea"
            type="checkbox"
            class="rounded border-gray-700 bg-gray-900 text-indigo-600 focus:ring-indigo-500"
            @change="commitChanges"
          />
          <span>{{ $t('inspector.waveform.fill_area') }}</span>
        </label>
      </div>

      <!-- If Circular: Inner Radius -->
      <div v-if="mode === 'circular'" class="space-y-1">
        <div class="flex justify-between text-gray-400 text-[11px]">
          <span>{{ $t('inspector.waveform.inner_radius') }}</span>
          <span class="font-mono text-gray-300">{{ Math.round(innerRadius * 100) }}%</span>
        </div>
        <input
          v-model.number="innerRadius"
          type="range"
          min="0.2"
          max="0.75"
          step="0.05"
          class="w-full h-1 bg-gray-800 rounded appearance-none accent-indigo-500 cursor-pointer"
          @input="commitChanges"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useTimelineStore } from '~/stores/timeline'
import { useProjectStore } from '~/stores/project'
import type { WaveformSpec } from '~/types/spec'

const { t } = useI18n()
const timelineStore = useTimelineStore()
const projectStore = useProjectStore()

const selectedClip = computed(() => timelineStore.selectedClip)
const selectedTrack = computed(() => timelineStore.selectedTrack)

// Target object holding waveform settings (clip or track)
const activeTarget = computed(() => {
  return selectedClip.value || (selectedTrack.value?.kind === 'waveform' ? selectedTrack.value : null)
})

const waveformConfig = computed<WaveformSpec>(() => {
  if (selectedClip.value?.waveform) return selectedClip.value.waveform
  if (selectedTrack.value?.waveform) return selectedTrack.value.waveform
  return {}
})

// Visual Modes List
const modesList = computed(() => [
  { id: 'peak_to_peak', label: t('inspector.waveform.modes.peak_to_peak'), shortLabel: 'Bars', icon: 'i-heroicons-chart-bar' },
  { id: 'spectrum', label: t('inspector.waveform.modes.spectrum'), shortLabel: 'EQ', icon: 'i-heroicons-bars-4' },
  { id: 'wave', label: t('inspector.waveform.modes.wave'), shortLabel: 'Wave', icon: 'i-heroicons-variable' },
  { id: 'circular', label: t('inspector.waveform.modes.circular'), shortLabel: 'Radial', icon: 'i-heroicons-lifebuoy' },
  { id: 'dots', label: t('inspector.waveform.modes.dots'), shortLabel: 'Dots', icon: 'i-heroicons-squares-2x2' },
])

// Style Presets
const stylePresets = [
  { id: 'cyberpunk', name: 'Cyberpunk Neon', mode: 'peak_to_peak', color: '#00E6B4', secondaryColor: '#6366F1', glow: true, glowRadius: 10, density: 48, roundness: 6 },
  { id: 'minimal', name: 'Podcast Minimal', mode: 'spectrum', color: '#F43F5E', secondaryColor: '#FB923C', glow: false, glowRadius: 4, density: 40, roundness: 4 },
  { id: 'retro', name: 'Studio Oscilloscope', mode: 'wave', color: '#10B981', secondaryColor: '#06B6D4', glow: true, glowRadius: 12, density: 48, lineWidth: 3, fillArea: true },
  { id: 'glow', name: 'Spotify Radial', mode: 'circular', color: '#A855F7', secondaryColor: '#EC4899', glow: true, glowRadius: 10, density: 64, innerRadius: 0.45 },
]

// State variables
const mode = ref('peak_to_peak')
const color = ref('#00E6B4')
const secondaryColor = ref('#6366F1')
const width = ref(1200)
const height = ref(180)
const positionX = ref(0)
const positionY = ref(0)
const opacity = ref(1.0)
const isGlow = ref(true)
const glowRadius = ref(8)
const density = ref(48)
const roundness = ref(6)
const lineWidth = ref(3)
const isFillArea = ref(false)
const innerRadius = ref(0.45)

// Synchronize state from active target
watch(
  waveformConfig,
  (cfg) => {
    mode.value = cfg.mode || 'peak_to_peak'
    color.value = cfg.color || '#00E6B4'
    secondaryColor.value = cfg.secondary_color || '#6366F1'
    width.value = Number(cfg.width) || Math.round((projectStore.canvasWidth || 1920) * 0.75)
    height.value = Number(cfg.height) || 180
    positionX.value = selectedClip.value?.position?.x ?? (cfg.position?.x ?? 0)
    positionY.value = selectedClip.value?.position?.y ?? (cfg.position?.y ?? 0)
    opacity.value = selectedClip.value?.opacity ?? (cfg.opacity ?? 1.0)
    isGlow.value = cfg.glow ?? true
    glowRadius.value = Number(cfg.glow_radius) || 8
    density.value = Number(cfg.bar_density) || 48
    roundness.value = Number(cfg.bar_roundness) ?? 6
    lineWidth.value = Number(cfg.line_width) || 3
    isFillArea.value = cfg.fill_area ?? false
    innerRadius.value = Number(cfg.inner_radius) || 0.45
  },
  { immediate: true, deep: true }
)

function selectMode(newMode: string) {
  mode.value = newMode
  if (newMode === 'circular') {
    height.value = Math.min(width.value, 400)
  }
  commitChanges()
}

function applyPreset(preset: typeof stylePresets[0]) {
  mode.value = preset.mode
  color.value = preset.color
  secondaryColor.value = preset.secondaryColor
  isGlow.value = preset.glow
  glowRadius.value = preset.glowRadius
  density.value = preset.density
  if (preset.roundness !== undefined) roundness.value = preset.roundness
  if (preset.lineWidth !== undefined) lineWidth.value = preset.lineWidth
  if (preset.fillArea !== undefined) isFillArea.value = preset.fillArea
  if (preset.innerRadius !== undefined) innerRadius.value = preset.innerRadius
  commitChanges()
  timelineStore.pushHistoryState(`Apply Waveform Preset ${preset.name}`)
}

function alignCenter() {
  positionX.value = 0
  positionY.value = 0
  commitChanges()
}

function alignLowerThird() {
  positionX.value = 0
  positionY.value = Math.round((projectStore.canvasHeight || 1080) * 0.28)
  commitChanges()
}

function commitChanges() {
  const patch: WaveformSpec = {
    mode: mode.value as any,
    color: color.value,
    secondary_color: secondaryColor.value,
    width: width.value,
    height: height.value,
    position: { x: positionX.value, y: positionY.value, alignment: 'center' },
    opacity: opacity.value,
    glow: isGlow.value,
    glow_radius: glowRadius.value,
    bar_density: density.value,
    bar_roundness: roundness.value,
    line_width: lineWidth.value,
    fill_area: isFillArea.value,
    inner_radius: innerRadius.value,
  }

  if (selectedClip.value) {
    if (!selectedClip.value.waveform) selectedClip.value.waveform = {}
    Object.assign(selectedClip.value.waveform, patch)
    if (!selectedClip.value.position) selectedClip.value.position = { x: 0, y: 0, alignment: 'center' }
    selectedClip.value.position.x = positionX.value
    selectedClip.value.position.y = positionY.value
    selectedClip.value.opacity = opacity.value
  } else if (selectedTrack.value) {
    if (!selectedTrack.value.waveform) selectedTrack.value.waveform = {}
    Object.assign(selectedTrack.value.waveform, patch)
  }

  timelineStore.pushHistoryState('Update Waveform Visualizer')
}

function resetToDefault() {
  mode.value = 'peak_to_peak'
  color.value = '#00E6B4'
  secondaryColor.value = '#6366F1'
  width.value = Math.round((projectStore.canvasWidth || 1920) * 0.75)
  height.value = 180
  positionX.value = 0
  positionY.value = 0
  opacity.value = 1.0
  isGlow.value = true
  glowRadius.value = 8
  density.value = 48
  roundness.value = 6
  lineWidth.value = 3
  isFillArea.value = false
  innerRadius.value = 0.45
  commitChanges()
}
</script>
