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
        <div v-else-if="error" class="space-y-4">
          <div class="p-4 bg-red-950/40 border border-red-800/80 rounded-lg text-xs text-red-300 font-mono">
            <div class="font-semibold mb-1">Failed to render graph:</div>
            <div>{{ error }}</div>
          </div>
          <div v-if="mermaidCode" class="space-y-2">
            <div class="flex items-center justify-between text-xs text-gray-400">
              <span>Mermaid Definition</span>
              <div class="flex items-center gap-2">
                <UButton
                  size="xs"
                  color="neutral"
                  variant="soft"
                  icon="i-heroicons-arrow-path"
                  :loading="isLoading"
                  @click="renderGraph"
                >
                  {{ $t('graph.refresh') }}
                </UButton>
                <UButton
                  size="xs"
                  color="neutral"
                  variant="soft"
                  icon="i-heroicons-clipboard-document"
                  @click="copyMermaidCode"
                >
                  {{ isCopied ? $t('graph.copied') : $t('graph.copy_mermaid') }}
                </UButton>
              </div>
            </div>
            <pre class="bg-gray-950 p-3 rounded-lg border border-gray-800 text-[11px] font-mono text-gray-300 overflow-x-auto max-h-40 leading-relaxed">{{ mermaidCode }}</pre>
          </div>
        </div>

        <!-- Graph Container -->
        <div v-else class="space-y-4">
          <div
            ref="graphContainerRef"
            dir="ltr"
            class="bg-gray-900/80 border border-gray-800 rounded-xl p-4 overflow-x-auto min-h-[320px] flex items-center [&>svg]:m-auto [&>svg]:shrink-0"
          ></div>

          <!-- Raw Mermaid Code & Actions -->
          <div class="space-y-2">
            <div class="flex items-center justify-between text-xs text-gray-400">
              <span>Mermaid Definition</span>
              <div class="flex items-center gap-2">
                <UButton
                  size="xs"
                  color="neutral"
                  variant="soft"
                  icon="i-heroicons-arrow-path"
                  :loading="isLoading"
                  @click="renderGraph"
                >
                  {{ $t('graph.refresh') }}
                </UButton>
                <UButton
                  size="xs"
                  color="neutral"
                  variant="soft"
                  icon="i-heroicons-arrow-down-tray"
                  @click="downloadSvg"
                >
                  {{ $t('graph.export_svg') }}
                </UButton>
                <UButton
                  size="xs"
                  color="neutral"
                  variant="soft"
                  icon="i-heroicons-clipboard-document"
                  @click="copyMermaidCode"
                >
                  {{ isCopied ? $t('graph.copied') : $t('graph.copy_mermaid') }}
                </UButton>
              </div>
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
  if (graphContainerRef.value) {
    graphContainerRef.value.innerHTML = ''
  }

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

    // Crucial: Set isLoading to false before DOM update so graphContainerRef is rendered in the DOM
    isLoading.value = false

    await nextTick()

    if (!graphContainerRef.value) {
      await new Promise((resolve) => setTimeout(resolve, 50))
    }

    if (graphContainerRef.value && response.content) {
      mermaid.initialize({
        startOnLoad: false,
        securityLevel: 'loose',
        theme: 'dark',
        flowchart: {
          useMaxWidth: false,
          htmlLabels: true,
          curve: 'basis',
        },
        themeVariables: {
          darkMode: true,
          background: '#111827',
          primaryColor: '#6366F1',
          primaryTextColor: '#F3F4F6',
          lineColor: '#818CF8',
          secondaryColor: '#1F2937',
          tertiaryColor: '#374151',
          edgeLabelBackground: '#1F2937',
        },
      })

      const id = `mermaid_${Date.now()}_${Math.random().toString(36).substring(2, 7)}`
      const { svg, bindFunctions } = await mermaid.render(id, response.content)
      graphContainerRef.value.innerHTML = svg
      if (bindFunctions) {
        bindFunctions(graphContainerRef.value)
      }
    }
  } catch (err: any) {
    console.error('Failed generating graph', err)
    error.value = err?.data?.message || err?.message || 'Failed generating graph'
  } finally {
    isLoading.value = false
  }
}

async function copyMermaidCode() {
  if (!mermaidCode.value) return
  await navigator.clipboard.writeText(mermaidCode.value)
  isCopied.value = true
  setTimeout(() => {
    isCopied.value = false
  }, 2000)
}

function downloadSvg() {
  if (!graphContainerRef.value) return
  const svgElement = graphContainerRef.value.querySelector('svg')
  if (!svgElement) return
  const svgData = new XMLSerializer().serializeToString(svgElement)
  const blob = new Blob([svgData], { type: 'image/svg+xml;charset=utf-8' })
  const url = URL.createObjectURL(blob)
  const downloadLink = document.createElement('a')
  downloadLink.href = url
  downloadLink.download = `vidonyx-filtergraph-${Date.now()}.svg`
  document.body.appendChild(downloadLink)
  downloadLink.click()
  document.body.removeChild(downloadLink)
  URL.revokeObjectURL(url)
}
</script>
