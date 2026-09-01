<template>
  <UModal
    v-model:open="projectStore.isSettingsOpen"
    :title="$t('project.title')"
    :ui="{ content: 'sm:max-w-md bg-gray-950 border border-gray-800' }"
  >
    <template #body>
      <div v-if="projectStore.currentProject" class="space-y-4 text-xs">
        <!-- Project Name -->
        <div>
          <label class="block text-gray-400 mb-1">{{ $t('project.name') }}</label>
          <input
            v-model="projectStore.currentProject.name"
            type="text"
            class="w-full bg-gray-900 border border-gray-800 rounded px-2.5 py-1.5 text-gray-200 focus:outline-none focus:border-indigo-500"
          />
        </div>

        <!-- Description -->
        <div>
          <label class="block text-gray-400 mb-1">{{ $t('project.description') }}</label>
          <textarea
            v-model="projectStore.currentProject.description"
            rows="2"
            class="w-full bg-gray-900 border border-gray-800 rounded px-2.5 py-1.5 text-gray-200 focus:outline-none focus:border-indigo-500"
          ></textarea>
        </div>

        <!-- Width & Height -->
        <div class="grid grid-cols-2 gap-3">
          <div>
            <label class="block text-gray-400 mb-1">{{ $t('project.width') }} (px)</label>
            <input
              v-model.number="projectStore.currentProject.width"
              type="number"
              class="w-full bg-gray-900 border border-gray-800 rounded px-2.5 py-1.5 font-mono text-gray-200"
            />
          </div>
          <div>
            <label class="block text-gray-400 mb-1">{{ $t('project.height') }} (px)</label>
            <input
              v-model.number="projectStore.currentProject.height"
              type="number"
              class="w-full bg-gray-900 border border-gray-800 rounded px-2.5 py-1.5 font-mono text-gray-200"
            />
          </div>
        </div>

        <!-- Frame Rate & Background Color -->
        <div class="grid grid-cols-2 gap-3">
          <div>
            <label class="block text-gray-400 mb-1">{{ $t('project.frame_rate') }}</label>
            <USelect
              v-model.number="projectStore.currentProject.frame_rate"
              :items="fpsOptions"
              size="sm"
            />
          </div>

          <div>
            <label class="block text-gray-400 mb-1">{{ $t('project.background_color') }}</label>
            <div class="flex items-center gap-2">
              <input
                v-model="projectStore.currentProject.background_color"
                type="color"
                class="w-8 h-8 rounded border border-gray-700 bg-transparent cursor-pointer"
              />
              <input
                v-model="projectStore.currentProject.background_color"
                type="text"
                class="w-full bg-gray-900 border border-gray-800 rounded px-2 py-1 font-mono text-gray-300 focus:outline-none"
              />
            </div>
          </div>
        </div>

        <UButton
          size="sm"
          color="primary"
          block
          class="font-semibold mt-4"
          @click="saveAndClose"
        >
          {{ $t('app.save') }}
        </UButton>
      </div>
    </template>
  </UModal>
</template>

<script setup lang="ts">
import { useProjectStore } from '~/stores/project'
import { useTimelineStore } from '~/stores/timeline'

const projectStore = useProjectStore()
const timelineStore = useTimelineStore()

const fpsOptions = [
  { label: '23.976 FPS (Cinema NTSC)', value: 23.976 },
  { label: '24 FPS (Film)', value: 24 },
  { label: '25 FPS (PAL Broadcast)', value: 25 },
  { label: '29.97 FPS (NTSC TV)', value: 29.97 },
  { label: '30 FPS (Standard Web)', value: 30 },
  { label: '60 FPS (Smooth / Gaming)', value: 60 },
]

function saveAndClose() {
  projectStore.saveCurrentProject(timelineStore.toVideoSpec())
  projectStore.isSettingsOpen = false
}
</script>
