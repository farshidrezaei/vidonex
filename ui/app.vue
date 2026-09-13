<template>
  <UApp>
    <div class="h-screen w-screen flex flex-col bg-[#0b0f19] text-gray-100 overflow-hidden font-sans">
      <!-- Top Header -->
      <LayoutAppHeader />

      <!-- Studio Workspace with Nested Resizable USplitter Panels -->
      <USplitter
        id="studio-vertical-splitter"
        auto-save-id="vidonex-studio-vertical"
        orientation="vertical"
        :items="verticalPanels"
        class="flex-1 overflow-hidden"
        :ui="{
          handle: 'group relative bg-gray-800/80 hover:bg-indigo-500/60 data-[state=drag]:bg-indigo-500 transition-colors cursor-row-resize h-1.5 flex items-center justify-center z-20 before:content-[\'\'] before:absolute before:-top-1 before:-bottom-1 before:inset-x-0'
        }"
      >
        <template #resize-handle>
          <div class="w-8 h-0.5 rounded-full bg-gray-600/70 group-hover:bg-indigo-200 transition-colors"></div>
        </template>

        <!-- Upper Workspace: Horizontal Splitter for Media, Viewport Preview, and Inspector -->
        <template #workspace>
          <USplitter
            id="studio-horizontal-splitter"
            auto-save-id="vidonex-studio-horizontal"
            orientation="horizontal"
            :items="horizontalPanels"
            class="h-full w-full overflow-hidden"
            :ui="{
              handle: 'group relative bg-gray-800/80 hover:bg-indigo-500/60 data-[state=drag]:bg-indigo-500 transition-colors cursor-col-resize w-1.5 flex items-center justify-center z-20 before:content-[\'\'] before:absolute before:-left-1 before:-right-1 before:inset-y-0'
            }"
          >
            <template #resize-handle>
              <div class="h-8 w-0.5 rounded-full bg-gray-600/70 group-hover:bg-indigo-200 transition-colors"></div>
            </template>

            <!-- Left Panel: Media Asset Library -->
            <template #media>
              <aside class="h-full w-full flex flex-col overflow-hidden">
                <PanelsMediaLibrary />
              </aside>
            </template>

            <!-- Center Panel: Interactive Canvas Viewport & Controls -->
            <template #preview>
              <main class="h-full w-full flex flex-col overflow-hidden">
                <ViewportCanvasViewport />
              </main>
            </template>

            <!-- Right Panel: Properties Inspector -->
            <template #inspector>
              <aside class="h-full w-full flex flex-col overflow-hidden">
                <PanelsInspectorPanel />
              </aside>
            </template>
          </USplitter>
        </template>

        <!-- Bottom Panel: Pro Multi-Track Timeline -->
        <template #timeline>
          <section class="h-full w-full flex flex-col overflow-hidden z-10">
            <TimelineContainer />
          </section>
        </template>
      </USplitter>

      <!-- Bottom Footer Status Bar -->
      <LayoutAppFooter />

      <!-- Modals & Dialogs -->
      <PanelsProjectSettingsModal />
      <PanelsGraphVisualizerModal />
      <PanelsRenderExportModal />
      <PanelsProbeDetailsModal />
      <PanelsKeyboardShortcutsModal v-model="isShortcutsModalOpen" />
    </div>
  </UApp>
</template>

<script setup lang="ts">
import type { SplitterItem } from '@nuxt/ui'
import { useProjectStore } from '~/stores/project'
import { useMediaStore } from '~/stores/media'
import { useTimelineStore } from '~/stores/timeline'
import { useWebSocket } from '~/composables/useWebSocket'
import { useGlobalShortcuts, isShortcutsModalOpen } from '~/composables/useGlobalShortcuts'

const { locale } = useI18n()
const projectStore = useProjectStore()
const mediaStore = useMediaStore()
const timelineStore = useTimelineStore()
const { connect } = useWebSocket()
const { registerGlobalListeners, unregisterGlobalListeners } = useGlobalShortcuts()

// Outer vertical layout splitter panels (Workspace vs Timeline)
const verticalPanels = ref<SplitterItem[]>([
  {
    slot: 'workspace',
    defaultSize: 65,
    minSize: 30,
  },
  {
    slot: 'timeline',
    defaultSize: 35,
    minSize: 15,
    collapsible: true,
    collapsedSize: 4,
  },
])

// Inner horizontal layout splitter panels (Media Library vs Preview Canvas vs Inspector)
const horizontalPanels = ref<SplitterItem[]>([
  {
    slot: 'media',
    defaultSize: 20,
    minSize: 12,
    collapsible: true,
  },
  {
    slot: 'preview',
    defaultSize: 55,
    minSize: 30,
  },
  {
    slot: 'inspector',
    defaultSize: 25,
    minSize: 16,
    collapsible: true,
  },
])

// Synchronize document direction (RTL for 'fa', LTR for 'en') reactively on page load, refresh & language toggle
useHead({
  htmlAttrs: {
    lang: computed(() => locale.value),
    dir: computed(() => (locale.value === 'fa' ? 'rtl' : 'ltr')),
  },
})

watch(
  () => locale.value,
  (newLocale) => {
    if (typeof document !== 'undefined') {
      document.documentElement.setAttribute('dir', newLocale === 'fa' ? 'rtl' : 'ltr')
      document.documentElement.setAttribute('lang', newLocale)
    }
  },
  { immediate: true }
)

onMounted(async () => {
  // Register global shortcuts capture listener
  registerGlobalListeners()

  // Connect live WebSocket hub
  connect()

  // Fetch or create initial project
  await projectStore.fetchProjects()

  if (projectStore.projectsList.length > 0) {
    await projectStore.loadProject(projectStore.projectsList[0].id)
  } else {
    await projectStore.createProject('My First Composition', 1920, 1080, 30.0)
  }

  if (projectStore.currentProject) {
    await mediaStore.fetchAssets(projectStore.currentProject.id)
    const specData = projectStore.currentProject.specification
    if (typeof specData === 'string') {
      try {
        const parsed = JSON.parse(specData)
        timelineStore.loadFromSpec(parsed)
      } catch (err) {
        console.error('Failed parsing project spec string:', err)
      }
    } else if (specData && typeof specData === 'object') {
      timelineStore.loadFromSpec(specData as any)
    }
  }
})

onUnmounted(() => {
  unregisterGlobalListeners()
})
</script>
