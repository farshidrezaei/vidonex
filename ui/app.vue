<template>
  <UApp>
    <div class="h-screen w-screen flex flex-col bg-[#0b0f19] text-gray-100 overflow-hidden font-sans">
      <!-- Top Header -->
      <LayoutAppHeader />

      <!-- Main Workspace (Center Split) -->
      <div class="flex-1 flex overflow-hidden">
        <!-- Left Panel: Media Asset Library -->
        <aside class="w-72 flex-shrink-0 flex flex-col overflow-hidden">
          <PanelsMediaLibrary />
        </aside>

        <!-- Center Panel: Interactive Canvas Viewport & Controls -->
        <main class="flex-1 flex flex-col overflow-hidden border-x border-gray-800">
          <ViewportCanvasViewport />
        </main>

        <!-- Right Panel: Properties Inspector -->
        <aside class="w-80 flex-shrink-0 flex flex-col overflow-hidden">
          <PanelsInspectorPanel />
        </aside>
      </div>

      <!-- Bottom Panel: Pro Multi-Track Timeline -->
      <section class="h-72 flex-shrink-0 border-t border-gray-800 flex flex-col overflow-hidden z-10">
        <TimelineContainer />
      </section>

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
import { useProjectStore } from '~/stores/project'
import { useMediaStore } from '~/stores/media'
import { useTimelineStore } from '~/stores/timeline'
import { useWebSocket } from '~/composables/useWebSocket'
import { useGlobalShortcuts, isShortcutsModalOpen } from '~/composables/useGlobalShortcuts'

const projectStore = useProjectStore()
const mediaStore = useMediaStore()
const timelineStore = useTimelineStore()
const { connect } = useWebSocket()
const { registerGlobalListeners, unregisterGlobalListeners } = useGlobalShortcuts()

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
    if (typeof projectStore.currentProject.specification === 'string') {
      try {
        const parsed = JSON.parse(projectStore.currentProject.specification)
        timelineStore.loadFromSpec(parsed)
      } catch {
        // Use initial tracks
      }
    }
  }
})

onUnmounted(() => {
  unregisterGlobalListeners()
})
</script>
