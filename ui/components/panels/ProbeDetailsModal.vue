<template>
  <UModal
    v-model:open="mediaStore.isProbeModalOpen"
    :title="$t('media.probe_details')"
    :ui="{ content: 'sm:max-w-lg bg-gray-950 border border-gray-800' }"
  >
    <template #body>
      <div v-if="asset" class="space-y-3 text-xs">
        <div class="p-3 bg-gray-900/60 rounded-lg border border-gray-800 space-y-2 font-mono">
          <div class="flex justify-between py-1 border-b border-gray-800">
            <span class="text-gray-500">{{ $t('media.file_name') }}:</span>
            <span class="text-gray-200 truncate max-w-[240px]">{{ asset.file_name }}</span>
          </div>
          <div class="flex justify-between py-1 border-b border-gray-800">
            <span class="text-gray-500">{{ $t('media.duration') }}:</span>
            <span class="text-indigo-400">{{ asset.duration_seconds }}s</span>
          </div>
          <div v-if="asset.width && asset.height" class="flex justify-between py-1 border-b border-gray-800">
            <span class="text-gray-500">{{ $t('media.resolution') }}:</span>
            <span class="text-gray-200">{{ asset.width }} x {{ asset.height }}</span>
          </div>
          <div v-if="asset.frame_rate" class="flex justify-between py-1 border-b border-gray-800">
            <span class="text-gray-500">{{ $t('media.fps') }}:</span>
            <span class="text-gray-200">{{ asset.frame_rate }} FPS</span>
          </div>
          <div v-if="asset.sample_rate" class="flex justify-between py-1 border-b border-gray-800">
            <span class="text-gray-500">{{ $t('media.sample_rate') }}:</span>
            <span class="text-gray-200">{{ asset.sample_rate }} Hz</span>
          </div>
          <div v-if="asset.channels" class="flex justify-between py-1 border-b border-gray-800">
            <span class="text-gray-500">{{ $t('media.channels') }}:</span>
            <span class="text-gray-200">{{ asset.channels }} Channels</span>
          </div>
          <div class="flex justify-between py-1">
            <span class="text-gray-500">{{ $t('media.size') }}:</span>
            <span class="text-gray-200">{{ asset.file_size_bytes }} bytes</span>
          </div>
        </div>
      </div>
    </template>
  </UModal>
</template>

<script setup lang="ts">
import { useMediaStore } from '~/stores/media'

const mediaStore = useMediaStore()
const asset = computed(() => mediaStore.activeProbeAsset)
</script>
