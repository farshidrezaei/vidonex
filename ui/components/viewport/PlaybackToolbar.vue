<template>
  <div class="h-12 bg-gray-950/80 border-t border-gray-800 flex items-center justify-between px-4 select-none z-20">
    <!-- Left: Timecode Display -->
    <div class="flex items-center gap-2 font-mono text-xs text-gray-300">
      <span class="px-2 py-1 bg-gray-900 rounded border border-gray-800 text-indigo-400 font-semibold">
        {{ formattedTime }}
      </span>
      <span class="text-gray-600">/</span>
      <span class="text-gray-500">
        {{ formattedDuration }}
      </span>
    </div>

    <!-- Center: Playback Transport Controls -->
    <div class="flex items-center gap-1">
      <!-- Jump to Start -->
      <UButton
        icon="i-heroicons-backward"
        size="xs"
        color="gray"
        variant="ghost"
        :title="$t('viewport.jump_start')"
        @click="playbackStore.jumpToStart"
      />

      <!-- Step Backward 1 Frame -->
      <UButton
        icon="i-heroicons-chevron-left"
        size="xs"
        color="gray"
        variant="ghost"
        :title="$t('viewport.step_backward')"
        @click="playbackStore.stepFrame(-1, projectStore.frameRate)"
      />

      <!-- Play / Pause Main Button -->
      <UButton
        :icon="playbackStore.isPlaying ? 'i-heroicons-pause' : 'i-heroicons-play'"
        size="sm"
        color="primary"
        class="w-9 h-9 rounded-full flex items-center justify-center bg-gradient-to-r from-indigo-500 to-purple-600 shadow-md shadow-indigo-500/20"
        :title="playbackStore.isPlaying ? $t('viewport.pause') : $t('viewport.play')"
        @click="playbackStore.togglePlay"
      />

      <!-- Step Forward 1 Frame -->
      <UButton
        icon="i-heroicons-chevron-right"
        size="xs"
        color="gray"
        variant="ghost"
        :title="$t('viewport.step_forward')"
        @click="playbackStore.stepFrame(1, projectStore.frameRate)"
      />

      <!-- Jump to End -->
      <UButton
        icon="i-heroicons-forward"
        size="xs"
        color="gray"
        variant="ghost"
        :title="$t('viewport.jump_end')"
        @click="playbackStore.jumpToEnd"
      />
    </div>

    <!-- Right: Loop & Volume Controls -->
    <div class="flex items-center gap-3">
      <!-- Loop Toggle -->
      <button
        class="p-1 rounded transition"
        :class="playbackStore.isLooping ? 'text-indigo-400 bg-indigo-500/10' : 'text-gray-500 hover:text-gray-300'"
        :title="$t('viewport.loop')"
        @click="playbackStore.isLooping = !playbackStore.isLooping"
      >
        <UIcon name="i-heroicons-arrow-path" class="w-4 h-4" />
      </button>

      <!-- Volume Slider -->
      <div class="flex items-center gap-1.5 group">
        <button
          class="text-gray-400 hover:text-gray-200 transition"
          @click="toggleMute"
        >
          <UIcon :name="playbackStore.isMuted || playbackStore.volume === 0 ? 'i-heroicons-speaker-x-mark' : 'i-heroicons-speaker-wave'" class="w-4 h-4" />
        </button>

        <input
          v-model.number="playbackStore.volume"
          type="range"
          min="0"
          max="1"
          step="0.05"
          class="w-16 h-1 bg-gray-800 rounded appearance-none accent-indigo-500 cursor-pointer"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { usePlaybackStore } from '~/stores/playback'
import { useProjectStore } from '~/stores/project'

const playbackStore = usePlaybackStore()
const projectStore = useProjectStore()

const formattedTime = computed(() => playbackStore.formatTimecode(playbackStore.currentTime, projectStore.frameRate))
const formattedDuration = computed(() => playbackStore.formatTimecode(playbackStore.duration, projectStore.frameRate))

function toggleMute() {
  playbackStore.isMuted = !playbackStore.isMuted
}
</script>
