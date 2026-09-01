<template>
  <UModal
    v-model:open="projectStore.isGraphOpen"
    :title="$t('graph.title')"
    :description="$t('graph.subtitle')"
    :ui="{ content: 'sm:max-w-4xl bg-gray-950 border border-gray-800' }"
  >
    <template #body>
      <div class="space-y-4">
        <!-- Loading State -->
        <div v-if="isLoading" class="flex flex-col items-center justify-center py-16 gap-3">
          <UIcon name="i-heroicons-arrow-path" class="w-8 h-8 text-indigo-500 animate-spin" />
          <span class="text-xs text-gray-400">Compiling Filtergraph DAG...</span>
        </div>

        <!-- Error State -->
        <div v-else-if="error" class="p-4 bg-red-950/40 border border-red-800/80 rounded-lg text-xs text-red-300 font-mono">
          {{ error }}
        </div>

        <!-- Graph Container -->
        <div v-else class="space-y-4">
          <div
            ref="graphContainerRef"
            class="bg-gray-900/80 border border-gray-800 rounded-xl p-4 overflow-x-auto min-h-[280px] flex items-center justify-center"
          ></div>

          <!-- Raw Mermaid Code -->
          <div class="space-y-2">
            <div class="flex items-center justify-between text-xs text-gray-400">
              <span>Mermaid Definition</span>
              <UButton
                size="xs"
                color="neutral"
                variant="soft"
                icon="i-heroicons-clipboard-document"
                @click="copyMermaidCode"
              >
                {{ isCopied ? 'Copied!' : $t('graph.copy_mermaid') }}
              </UButton>
            </div>
            <pre class="bg-gray-950 p-3 rounded-lg border border-gray-800 text-[11px] font-mono text-gray-300 overflow-x-auto max-h-40 leading-relaxed">{{ mermaidCode }}</pre>
          </div>
        </div>
      </div>
    </template>
  </UModal>
</template>

<script setup lang="ts">
import mermaid from 'mermaid'
import { useProjectStore } from '~/stores/project'
import { useTimelineStore } from '~/stores/timeline'

const projectStore = useProjectStore()
const timelineStore = useTimelineStore()

const isLoading = ref(false)
const error = ref<string | null>(null)
const mermaidCode = ref('')
const isCopied = ref(false)
const graphContainerRef = ref<HTMLDivElement | null>(null)

watch(() => projectStore.isGraphOpen, async (isOpen) => {
  if (isOpen) {
    await renderGraph()
  }
})

async function renderGraph() {
  isLoading.value = true
  error.value = null
  mermaidCode.value = ''

  try {
    const spec = timelineStore.toVideoSpec()
    const response = await $fetch<{ format: string; content: string }>('/api/spec/graph', {
      method: 'POST',
      body: {
        specification: spec,
        format: 'mermaid',
      },
    })

    mermaidCode.value = response.content

    await nextTick()
    if (graphContainerRef.value) {
      mermaid.initialize({
        startOnLoad: false,
        theme: 'dark',
        themeVariables: {
          darkMode: true,
          background: '#111827',
          primaryColor: '#6366F1',
          primaryTextColor: '#F3F4F6',
          lineColor: '#818CF8',
        },
      })
      const id = `mermaid_${Date.now()}`
      const { svg } = await mermaid.render(id, response.content)
      graphContainerRef.value.innerHTML = svg
    }
  } catch (err: any) {
    console.error('Failed generating graph', err)
    error.value = err?.data?.message || err?.message || 'Failed generating graph'
  } finally {
    isLoading.value = false
  }
}

async function copyMermaidCode() {
  await navigator.clipboard.writeText(mermaidCode.value)
  isCopied.value = true
  setTimeout(() => {
    isCopied.value = false
  }, 2000)
}
</script>
