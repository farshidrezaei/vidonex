<template>
  <header class="h-14 border-b border-gray-800 bg-gray-950/90 backdrop-blur px-4 flex items-center justify-between select-none z-30">
    <!-- Left: Branding & Project Title -->
    <div class="flex items-center gap-3">
      <div class="flex items-center gap-2">
        <img src="/logo.png" alt="Vidonex Logo" class="w-8 h-8 rounded-lg object-contain shadow-lg shadow-indigo-500/20" />
        <span class="font-bold text-base tracking-tight bg-gradient-to-r from-white via-gray-200 to-gray-400 bg-clip-text text-transparent">
          {{ $t('app.title') }}
        </span>
      </div>

      <div class="h-5 w-px bg-gray-800 mx-1"></div>

      <!-- Project Switcher & Name Input -->
      <div class="flex items-center gap-1.5 group">
        <UDropdownMenu :items="projectMenuItems">
          <UButton
            color="neutral"
            variant="ghost"
            size="xs"
            class="p-1 text-gray-400 hover:text-white"
            :title="$t('project.title')"
          >
            <UIcon name="i-heroicons-folder" class="w-4 h-4 text-indigo-400" />
            <UIcon name="i-heroicons-chevron-down" class="w-3 h-3 text-gray-500" />
          </UButton>
        </UDropdownMenu>

        <input
          v-if="projectStore.currentProject"
          v-model="projectStore.currentProject.name"
          class="bg-transparent hover:bg-gray-900 focus:bg-gray-900 border border-transparent hover:border-gray-700 focus:border-indigo-500 rounded px-2 py-1 text-sm font-medium text-gray-200 focus:outline-none transition max-w-[180px] truncate"
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
      <!-- Command Palette Button -->
      <UButton
        icon="i-heroicons-magnifying-glass"
        size="sm"
        color="neutral"
        variant="ghost"
        :title="$t('command_palette.button') || 'Command Palette (Ctrl+K)'"
        @click="isCommandPaletteOpen = true"
      >
        <span class="hidden xl:inline text-xs font-medium">{{ $t('command_palette.button') || 'Commands' }}</span>
        <UKbd size="xs" class="hidden sm:inline ml-1 font-mono text-[10px]">⌘K</UKbd>
      </UButton>

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

      <!-- About / Updates Button -->
      <div class="relative">
        <UButton
          icon="i-heroicons-information-circle"
          size="sm"
          color="neutral"
          variant="ghost"
          :title="$t('about.title') || 'About Vidonex & Updates'"
          @click="isAboutModalOpen = true"
        >
          <span class="hidden lg:inline text-xs font-medium">{{ $t('about.button') || 'About' }}</span>
        </UButton>
        <!-- Update Available Indicator Dot -->
        <span
          v-if="hasUpdate"
          class="absolute -top-0.5 -right-0.5 flex h-2 w-2"
          :title="$t('about.update_available')"
        >
          <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-amber-400 opacity-75"></span>
          <span class="relative inline-flex rounded-full h-2 w-2 bg-amber-500"></span>
        </span>
      </div>

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
import { isShortcutsModalOpen, isCommandPaletteOpen } from '~/composables/useGlobalShortcuts'
import { useConfirmDialog } from '~/composables/useConfirmDialog'
import { isAboutModalOpen, useAppUpdate } from '~/composables/useAppUpdate'

const { locale, setLocale, t } = useI18n()
const { confirm } = useConfirmDialog()
const projectStore = useProjectStore()
const timelineStore = useTimelineStore()
const { hasUpdate } = useAppUpdate()

const currentLocale = computed(() => locale.value)

const projectMenuItems = computed(() => {
  const currentId = projectStore.currentProject?.id
  const projects = projectStore.projectsList.map((p) => ({
    label: p.name || 'Untitled Project',
    icon: p.id === currentId ? 'i-heroicons-check' : 'i-heroicons-film',
    color: p.id === currentId ? ('primary' as const) : undefined,
    onSelect: async () => {
      if (p.id !== currentId) {
        await projectStore.switchProject(p.id)
      }
    },
  }))

  return [
    projects,
    [
      {
        label: t('app.new_project') || 'New Project...',
        icon: 'i-heroicons-plus-circle',
        color: 'primary' as const,
        onSelect: () => {
          projectStore.isNewProjectOpen = true
        },
      },
      ...(projectStore.projectsList.length > 1
        ? [
            {
              label: t('confirm.delete_project_title') || 'Delete Project',
              icon: 'i-heroicons-trash',
              color: 'error' as const,
              onSelect: async () => {
                if (!projectStore.currentProject) return
                const currentName = projectStore.currentProject.name
                const ok = await confirm({
                  title: t('confirm.delete_project_title'),
                  message: t('confirm.delete_project_message', { name: currentName }),
                  confirmText: t('confirm.delete'),
                  cancelText: t('confirm.cancel'),
                  isDanger: true,
                })
                if (ok) {
                  await projectStore.deleteProject(projectStore.currentProject.id)
                }
              },
            },
          ]
        : []),
    ],
  ]
})

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
