<template>
  <header class="h-14 border-b border-gray-800 bg-gray-950/90 backdrop-blur px-4 flex items-center justify-between select-none z-30">
    <!-- Left: Branding & Project Title -->
    <div class="flex items-center gap-3">
      <div class="flex items-center gap-2">
        <div class="w-8 h-8 rounded-lg bg-gradient-to-br from-indigo-500 to-purple-600 flex items-center justify-center shadow-lg shadow-indigo-500/20">
          <UIcon name="i-heroicons-film" class="w-5 h-5 text-white" />
        </div>
        <span class="font-bold text-base tracking-tight bg-gradient-to-r from-white via-gray-200 to-gray-400 bg-clip-text text-transparent">
          {{ $t('app.title') }}
        </span>
      </div>

      <div class="h-5 w-px bg-gray-800 mx-1"></div>

      <!-- Project Name Input -->
      <div class="flex items-center gap-2 group">
        <input
          v-if="projectStore.currentProject"
          v-model="projectStore.currentProject.name"
          class="bg-transparent hover:bg-gray-900 focus:bg-gray-900 border border-transparent hover:border-gray-700 focus:border-indigo-500 rounded px-2 py-1 text-sm font-medium text-gray-200 focus:outline-none transition"
          @blur="saveProject"
        />
        <span v-if="projectStore.isSaving" class="text-xs text-indigo-400 animate-pulse flex items-center gap-1">
          <UIcon name="i-heroicons-arrow-path" class="w-3 h-3 animate-spin" />
          {{ $t('app.saving') }}
        </span>
        <span v-else class="text-xs text-gray-500 flex items-center gap-1">
          <UIcon name="i-heroicons-check" class="w-3 h-3 text-emerald-500" />
          {{ $t('app.saved') }}
        </span>
      </div>
    </div>

    <!-- Center: Quick Presets & Undo/Redo -->
    <div class="flex items-center gap-2">
      <!-- Undo / Redo -->
      <div class="flex items-center bg-gray-900 rounded-lg p-0.5 border border-gray-800">
        <UButton
          icon="i-heroicons-arrow-uturn-left"
          size="xs"
          color="neutral"
          variant="ghost"
          :disabled="timelineStore.undoStack.length === 0"
          title="Undo (Ctrl+Z)"
          @click="timelineStore.undo()"
        />
        <UButton
          icon="i-heroicons-arrow-uturn-right"
          size="xs"
          color="neutral"
          variant="ghost"
          :disabled="timelineStore.redoStack.length === 0"
          title="Redo (Ctrl+Y / Ctrl+Shift+Z)"
          @click="timelineStore.redo()"
        />
      </div>

      <!-- Aspect Ratio Preset Selector -->
      <UDropdownMenu :items="presetMenuItems">
        <UButton
          color="neutral"
          variant="soft"
          size="sm"
          icon="i-heroicons-rectangle-group"
          class="text-xs font-mono"
        >
          {{ projectStore.aspectRatio }} ({{ projectStore.canvasWidth }}x{{ projectStore.canvasHeight }})
        </UButton>
      </UDropdownMenu>

      <!-- Project Settings -->
      <UButton
        icon="i-heroicons-cog-6-tooth"
        size="sm"
        color="neutral"
        variant="ghost"
        :title="$t('project.title')"
        @click="projectStore.isSettingsOpen = true"
      />
    </div>

    <!-- Right: View Graph, Shortcuts Button, Render Export, Language -->
    <div class="flex items-center gap-2">
      <!-- Keyboard Shortcuts Button -->
      <UButton
        icon="i-heroicons-command-line"
        size="sm"
        color="neutral"
        variant="ghost"
        title="Keyboard Shortcuts (? / Ctrl+/)"
        @click="isShortcutsModalOpen = true"
      >
        <span class="hidden md:inline text-xs font-medium">{{ $t('shortcuts.button') || 'Shortcuts' }}</span>
        <UKbd size="xs" class="hidden lg:inline ml-1 font-mono text-[10px]">?</UKbd>
      </UButton>

      <!-- Filtergraph Visualizer -->
      <UButton
        icon="i-heroicons-cpu-chip"
        size="sm"
        color="neutral"
        variant="ghost"
        @click="projectStore.isGraphOpen = true"
      >
        {{ $t('app.view_graph') }}
      </UButton>

      <!-- Render & Export Button -->
      <UButton
        icon="i-heroicons-arrow-down-tray"
        size="sm"
        color="primary"
        class="font-semibold bg-gradient-to-r from-indigo-500 to-purple-600 hover:from-indigo-600 hover:to-purple-700 shadow-md shadow-indigo-500/20"
        @click="projectStore.isExportOpen = true"
      >
        {{ $t('app.export') }}
      </UButton>

      <div class="h-5 w-px bg-gray-800 mx-1"></div>

      <!-- Language Switcher -->
      <UButton
        size="xs"
        color="neutral"
        variant="ghost"
        class="font-mono text-xs"
        @click="toggleLanguage"
      >
        {{ currentLocale.toUpperCase() }}
      </UButton>
    </div>
  </header>
</template>

<script setup lang="ts">
import { useProjectStore, ASPECT_RATIO_PRESETS } from '~/stores/project'
import { useTimelineStore } from '~/stores/timeline'
import { isShortcutsModalOpen } from '~/composables/useGlobalShortcuts'

const { locale, setLocale } = useI18n()
const projectStore = useProjectStore()
const timelineStore = useTimelineStore()

const currentLocale = computed(() => locale.value)

function toggleLanguage() {
  const next = currentLocale.value === 'en' ? 'fa' : 'en'
  setLocale(next)
  if (typeof document !== 'undefined') {
    document.documentElement.setAttribute('dir', next === 'fa' ? 'rtl' : 'ltr')
  }
}

function saveProject() {
  projectStore.saveCurrentProject(timelineStore.toVideoSpec())
}

const presetMenuItems = computed(() => [
  ASPECT_RATIO_PRESETS.map((preset) => ({
    label: `${preset.name} (${preset.width}x${preset.height})`,
    icon: preset.ratio === '9:16' ? 'i-heroicons-device-phone-mobile' : 'i-heroicons-computer-desktop',
    onSelect: () => {
      projectStore.setPreset(preset)
      timelineStore.saveCurrentTimeline()
    },
  })),
])
</script>
