<template>
  <UModal
    v-model:open="isOpen"
    :ui="{
      content: 'sm:max-w-xl p-0 overflow-hidden bg-gray-950 border border-gray-800 shadow-2xl rounded-xl'
    }"
  >
    <template #content>
      <UCommandPalette
        :groups="commandGroups"
        :placeholder="$t('command_palette.placeholder') || 'Type a command or search actions...'"
        :autofocus="true"
        class="min-h-[320px] max-h-[480px]"
        @update:model-value="handleSelect"
      />
    </template>
  </UModal>
</template>

<script setup lang="ts">
import { isCommandPaletteOpen, isShortcutsModalOpen } from '~/composables/useGlobalShortcuts'
import { useTimelineStore } from '~/stores/timeline'
import { usePlaybackStore } from '~/stores/playback'
import { useProjectStore } from '~/stores/project'
import { resetViewportZoom } from '~/composables/useViewportZoom'
import { isAboutModalOpen, useAppUpdate } from '~/composables/useAppUpdate'

const { t } = useI18n()
const { checkForUpdates } = useAppUpdate()
const timelineStore = useTimelineStore()
const playbackStore = usePlaybackStore()
const projectStore = useProjectStore()

const isOpen = computed({
  get: () => isCommandPaletteOpen.value,
  set: (val: boolean) => {
    isCommandPaletteOpen.value = val
  },
})

interface ActionCommandItem {
  id: string
  label: string
  description?: string
  icon?: string
  kbds?: string[]
  action: () => void
}

const commandGroups = computed(() => [
  {
    id: 'timeline-editing',
    label: t('command_palette.groups.editing') || 'Timeline Editing',
    items: [
      {
        id: 'split-clip',
        label: t('command_palette.actions.split_clip') || 'Split Clip at Playhead',
        description: t('command_palette.actions.split_clip_desc') || 'Cut selected clip at current playhead position',
        icon: 'i-heroicons-scissors',
        kbds: ['S'],
        action: () => {
          timelineStore.splitClipAtPlayhead(playbackStore.currentTime)
        },
      },
      {
        id: 'duplicate-clip',
        label: t('command_palette.actions.duplicate_clip') || 'Duplicate Selected Clip',
        description: t('command_palette.actions.duplicate_clip_desc') || 'Clone selected clip to track',
        icon: 'i-heroicons-document-duplicate',
        kbds: ['Ctrl', 'D'],
        action: () => {
          if (timelineStore.selectedClipId) {
            timelineStore.duplicateClip(timelineStore.selectedClipId)
          }
        },
      },
      {
        id: 'delete-clip',
        label: t('command_palette.actions.delete_clip') || 'Delete Selected Clip',
        description: t('command_palette.actions.delete_clip_desc') || 'Remove clip from timeline',
        icon: 'i-heroicons-trash',
        kbds: ['Del'],
        action: () => {
          if (timelineStore.selectedClipId) {
            timelineStore.removeClip(timelineStore.selectedClipId)
          }
        },
      },
      {
        id: 'undo',
        label: t('app.undo') || 'Undo',
        description: t('command_palette.actions.undo_desc') || 'Revert last timeline change',
        icon: 'i-heroicons-arrow-uturn-left',
        kbds: ['Ctrl', 'Z'],
        action: () => {
          timelineStore.undo()
        },
      },
      {
        id: 'redo',
        label: t('app.redo') || 'Redo',
        description: t('command_palette.actions.redo_desc') || 'Reapply last undone change',
        icon: 'i-heroicons-arrow-uturn-right',
        kbds: ['Ctrl', 'Y'],
        action: () => {
          timelineStore.redo()
        },
      },
      {
        id: 'toggle-snapping',
        label: t('command_palette.actions.toggle_snapping') || 'Toggle Magnetic Snapping',
        description: timelineStore.isSnappingEnabled ? 'Currently Enabled' : 'Currently Disabled',
        icon: 'i-heroicons-arrows-pointing-in',
        kbds: ['N'],
        action: () => {
          timelineStore.isSnappingEnabled = !timelineStore.isSnappingEnabled
        },
      },
    ],
  },
  {
    id: 'tracks',
    label: t('command_palette.groups.tracks') || 'Track Management',
    items: [
      {
        id: 'add-video-track',
        label: t('command_palette.actions.add_video_track') || 'Add Video Track',
        description: 'New overlay video track on top',
        icon: 'i-heroicons-film',
        action: () => {
          timelineStore.addTrack('video')
        },
      },
      {
        id: 'add-audio-track',
        label: t('command_palette.actions.add_audio_track') || 'Add Audio Track',
        description: 'New background music / voiceover track',
        icon: 'i-heroicons-musical-note',
        action: () => {
          timelineStore.addTrack('audio')
        },
      },
      {
        id: 'add-waveform-track',
        label: t('command_palette.actions.add_waveform_track') || 'Add Audiogram Waveform Track',
        description: 'Animated visual waveform overlay track',
        icon: 'i-heroicons-chart-bar',
        action: () => {
          timelineStore.addTrack('waveform')
        },
      },
    ],
  },
  {
    id: 'aspect-ratios',
    label: t('command_palette.groups.aspect_ratio') || 'Canvas Aspect Ratio Presets',
    items: [
      {
        id: 'aspect-16-9',
        label: '16:9 Landscape (1920x1080)',
        description: 'YouTube, Web, Broadcast Full HD',
        icon: 'i-heroicons-tv',
        action: () => {
          projectStore.setAspectRatio('16:9')
        },
      },
      {
        id: 'aspect-9-16',
        label: '9:16 Vertical (1080x1920)',
        description: 'TikTok, Instagram Reels, YouTube Shorts',
        icon: 'i-heroicons-device-phone-mobile',
        action: () => {
          projectStore.setAspectRatio('9:16')
        },
      },
      {
        id: 'aspect-1-1',
        label: '1:1 Square (1080x1080)',
        description: 'Instagram Square Posts & Carousels',
        icon: 'i-heroicons-stop',
        action: () => {
          projectStore.setAspectRatio('1:1')
        },
      },
      {
        id: 'aspect-21-9',
        label: '21:9 Ultrawide Cinema (2560x1080)',
        description: 'Cinematic Anamorphic Widescreen',
        icon: 'i-heroicons-rectangle-group',
        action: () => {
          projectStore.setAspectRatio('21:9')
        },
      },
    ],
  },
  {
    id: 'playback',
    label: t('command_palette.groups.playback') || 'Playback & Viewport',
    items: [
      {
        id: 'toggle-play',
        label: playbackStore.isPlaying
          ? (t('command_palette.actions.pause_video') || 'Pause Playback')
          : (t('command_palette.actions.play_video') || 'Play Video'),
        icon: playbackStore.isPlaying ? 'i-heroicons-pause' : 'i-heroicons-play',
        kbds: ['Space'],
        action: () => {
          playbackStore.togglePlay()
        },
      },
      {
        id: 'jump-start',
        label: t('command_palette.actions.jump_start') || 'Jump to Start',
        description: t('command_palette.actions.jump_start_desc') || 'Move playhead to 00:00:00',
        icon: 'i-heroicons-backward',
        kbds: ['Home'],
        action: () => {
          playbackStore.jumpToStart()
        },
      },
      {
        id: 'jump-end',
        label: t('command_palette.actions.jump_end') || 'Jump to End of Timeline',
        description: t('command_palette.actions.jump_end_desc') || 'Move playhead to end of timeline',
        icon: 'i-heroicons-forward',
        kbds: ['End'],
        action: () => {
          playbackStore.jumpToEnd()
        },
      },
      {
        id: 'reset-zoom',
        label: t('command_palette.actions.fit_screen') || 'Fit Viewport to Screen',
        description: t('command_palette.actions.fit_screen_desc') || 'Reset canvas zoom and pan to fit',
        icon: 'i-heroicons-arrows-pointing-out',
        kbds: ['Shift', 'Z'],
        action: () => {
          resetViewportZoom()
        },
      },
    ],
  },
  {
    id: 'project-export',
    label: t('command_palette.groups.project') || 'Project & Modals',
    items: [
      {
        id: 'render-export',
        label: t('command_palette.actions.render_export') || 'Render & Export Video...',
        description: t('command_palette.actions.render_export_desc') || 'Compile and encode composition with FFmpeg',
        icon: 'i-heroicons-arrow-up-tray',
        kbds: ['Ctrl', 'E'],
        action: () => {
          projectStore.isExportOpen = true
        },
      },
      {
        id: 'new-project',
        label: t('app.new_project') || 'New Project...',
        description: 'Create a new video composition project',
        icon: 'i-heroicons-plus-circle',
        action: () => {
          projectStore.isNewProjectOpen = true
        },
      },
      {
        id: 'project-settings',
        label: t('command_palette.actions.project_settings') || 'Project Settings...',
        description: t('command_palette.actions.project_settings_desc') || 'Adjust resolution, frame rate and background color',
        icon: 'i-heroicons-cog-6-tooth',
        action: () => {
          projectStore.isSettingsOpen = true
        },
      },
      {
        id: 'shortcuts-modal',
        label: t('command_palette.actions.keyboard_shortcuts') || 'Keyboard Shortcuts Reference...',
        description: t('command_palette.actions.keyboard_shortcuts_desc') || 'View complete hotkey cheat sheet',
        icon: 'i-heroicons-command-line',
        action: () => {
          isShortcutsModalOpen.value = true
        },
      },
      {
        id: 'about-modal',
        label: t('about.title') || 'About Vidonex Studio...',
        description: 'View engine version, system diagnostics and check for updates',
        icon: 'i-heroicons-information-circle',
        action: () => {
          isAboutModalOpen.value = true
        },
      },
      {
        id: 'check-updates',
        label: t('about.check_updates') || 'Check for Updates...',
        description: 'Check for the latest Vidonex release on GitHub',
        icon: 'i-heroicons-arrow-path',
        action: () => {
          isAboutModalOpen.value = true
          checkForUpdates(true)
        },
      },
      {
        id: 'graph-inspector',
        label: t('command_palette.actions.filtergraph') || 'Filtergraph Visualizer...',
        description: t('command_palette.actions.filtergraph_desc') || 'Inspect live DAG flowchart in Mermaid.js',
        icon: 'i-heroicons-cpu-chip',
        action: () => {
          projectStore.isGraphOpen = true
        },
      },
    ],
  },
])

function handleSelect(item: any) {
  if (!item) return
  isOpen.value = false
  if (typeof item.action === 'function') {
    item.action()
  }
}
</script>
