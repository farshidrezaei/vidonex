<template>
  <UModal
    v-model:open="projectStore.isNewProjectOpen"
    :title="$t('app.new_project') || 'New Project'"
    :ui="{ content: 'sm:max-w-md bg-gray-950 border border-gray-800' }"
  >
    <template #body>
      <form class="space-y-4 text-xs" @submit.prevent="handleCreate">
        <!-- Project Name -->
        <div>
          <label class="block text-gray-400 mb-1 font-medium">{{ $t('project.name') || 'Project Name' }}</label>
          <input
            v-model="projectName"
            type="text"
            required
            placeholder="e.g. My New Composition"
            class="w-full bg-gray-900 border border-gray-800 focus:border-indigo-500 rounded px-2.5 py-1.5 text-gray-200 focus:outline-none transition"
          />
        </div>

        <!-- Presets -->
        <div>
          <label class="block text-gray-400 mb-1.5 font-medium">{{ $t('project.aspect_ratio') || 'Preset / Ratio' }}</label>
          <div class="grid grid-cols-2 gap-2">
            <button
              v-for="preset in ASPECT_RATIO_PRESETS"
              :key="preset.id"
              type="button"
              class="p-2 rounded-lg border text-left flex items-center gap-2 transition"
              :class="
                selectedPresetId === preset.id
                  ? 'border-indigo-500 bg-indigo-500/10 text-indigo-200'
                  : 'border-gray-800 hover:border-gray-700 bg-gray-900/50 text-gray-400 hover:text-gray-300'
              "
              @click="applyPreset(preset)"
            >
              <UIcon :name="preset.icon" class="w-4 h-4 shrink-0 text-indigo-400" />
              <div class="truncate">
                <div class="font-medium text-[11px] truncate">{{ preset.name.split(' (')[0] }}</div>
                <div class="text-[10px] text-gray-500 font-mono">{{ preset.width }}x{{ preset.height }}</div>
              </div>
            </button>
          </div>
        </div>

        <!-- Dimensions -->
        <div class="grid grid-cols-2 gap-3">
          <div>
            <label class="block text-gray-400 mb-1 font-medium">{{ $t('project.width') || 'Width' }} (px)</label>
            <input
              v-model.number="width"
              type="number"
              min="320"
              max="7680"
              step="2"
              required
              class="w-full bg-gray-900 border border-gray-800 rounded px-2.5 py-1.5 font-mono text-gray-200 focus:border-indigo-500 focus:outline-none"
            />
          </div>
          <div>
            <label class="block text-gray-400 mb-1 font-medium">{{ $t('project.height') || 'Height' }} (px)</label>
            <input
              v-model.number="height"
              type="number"
              min="240"
              max="4320"
              step="2"
              required
              class="w-full bg-gray-900 border border-gray-800 rounded px-2.5 py-1.5 font-mono text-gray-200 focus:border-indigo-500 focus:outline-none"
            />
          </div>
        </div>

        <!-- Frame Rate -->
        <div>
          <label class="block text-gray-400 mb-1 font-medium">{{ $t('project.frame_rate') || 'Frame Rate' }}</label>
          <USelect
            v-model.number="frameRate"
            :items="fpsOptions"
            size="sm"
          />
        </div>
      </form>
    </template>

    <template #footer>
      <div class="flex items-center justify-end gap-2 w-full">
        <UButton
          color="neutral"
          variant="ghost"
          size="sm"
          @click="projectStore.isNewProjectOpen = false"
        >
          {{ $t('confirm.cancel') || 'Cancel' }}
        </UButton>
        <UButton
          color="primary"
          size="sm"
          icon="i-heroicons-plus"
          class="font-semibold"
          :loading="isCreating"
          @click="handleCreate"
        >
          {{ $t('app.new_project') || 'Create Project' }}
        </UButton>
      </div>
    </template>
  </UModal>
</template>

<script setup lang="ts">
import { useProjectStore, ASPECT_RATIO_PRESETS } from '~/stores/project'
import type { AspectRatioPreset } from '~/types/project'

const projectStore = useProjectStore()

const projectName = ref('New Composition')
const selectedPresetId = ref('youtube_16_9')
const width = ref(1920)
const height = ref(1080)
const frameRate = ref(30)
const isCreating = ref(false)

const fpsOptions = [
  { label: '24 FPS (Film)', value: 24 },
  { label: '25 FPS (PAL Broadcast)', value: 25 },
  { label: '30 FPS (Standard Web)', value: 30 },
  { label: '60 FPS (Smooth / Gaming)', value: 60 },
]

function applyPreset(preset: AspectRatioPreset) {
  selectedPresetId.value = preset.id
  width.value = preset.width
  height.value = preset.height
}

async function handleCreate() {
  if (!projectName.value.trim()) return

  isCreating.value = true
  try {
    const created = await projectStore.createProject(
      projectName.value.trim(),
      width.value,
      height.value,
      frameRate.value
    )
    if (created) {
      await projectStore.switchProject(created.id)
      projectStore.isNewProjectOpen = false
      projectName.value = 'New Composition'
    }
  } finally {
    isCreating.value = false
  }
}
</script>
