<template>
  <div class="space-y-4 text-xs">
    <!-- Font Size & Color -->
    <div class="grid grid-cols-2 gap-3">
      <div>
        <label class="block text-gray-400 mb-1">{{ $t('inspector.subtitle.font_size') }}</label>
        <input
          v-model.number="fontSize"
          type="number"
          min="12"
          max="120"
          class="w-full bg-gray-950 border border-gray-800 rounded px-2 py-1 font-mono text-gray-300 focus:outline-none focus:border-indigo-500"
          @change="updateSubtitle"
        />
      </div>

      <div>
        <label class="block text-gray-400 mb-1">{{ $t('inspector.subtitle.text_color') }}</label>
        <div class="flex items-center gap-1.5">
          <input
            v-model="textColor"
            type="color"
            class="w-7 h-7 rounded border border-gray-700 bg-transparent cursor-pointer"
            @input="updateSubtitle"
          />
          <input
            v-model="textColor"
            type="text"
            class="w-full bg-gray-950 border border-gray-800 rounded px-1.5 py-1 font-mono text-gray-300 text-[11px]"
            @change="updateSubtitle"
          />
        </div>
      </div>
    </div>

    <!-- Outline Options -->
    <div class="grid grid-cols-2 gap-3">
      <div>
        <label class="block text-gray-400 mb-1">{{ $t('inspector.subtitle.outline_width') }}</label>
        <input
          v-model.number="outlineWidth"
          type="number"
          min="0"
          max="10"
          class="w-full bg-gray-950 border border-gray-800 rounded px-2 py-1 font-mono text-gray-300"
          @change="updateSubtitle"
        />
      </div>

      <div>
        <label class="block text-gray-400 mb-1">{{ $t('inspector.subtitle.outline_color') }}</label>
        <div class="flex items-center gap-1.5">
          <input
            v-model="outlineColor"
            type="color"
            class="w-7 h-7 rounded border border-gray-700 bg-transparent cursor-pointer"
            @input="updateSubtitle"
          />
          <input
            v-model="outlineColor"
            type="text"
            class="w-full bg-gray-950 border border-gray-800 rounded px-1.5 py-1 font-mono text-gray-300 text-[11px]"
            @change="updateSubtitle"
          />
        </div>
      </div>
    </div>

    <!-- Background Box Toggle -->
    <div class="flex items-center justify-between pt-1">
      <span class="text-gray-400">{{ $t('inspector.subtitle.box') }}</span>
      <USwitch v-model="hasBox" @update:model-value="updateSubtitle" />
    </div>

    <!-- Subtitle Cues List -->
    <div class="pt-3 border-t border-gray-800 space-y-2">
      <div class="flex items-center justify-between">
        <span class="font-medium text-gray-300">{{ $t('inspector.subtitle.cues') }}</span>
        <UButton
          size="2xs"
          color="neutral"
          variant="soft"
          icon="i-heroicons-plus"
          @click="addCue"
        >
          Add Cue
        </UButton>
      </div>

      <div class="space-y-2 max-h-48 overflow-y-auto pr-1">
        <div
          v-for="(cue, idx) in cues"
          :key="idx"
          class="p-2 bg-gray-950 rounded border border-gray-800 space-y-1.5"
        >
          <div class="flex items-center justify-between gap-2">
            <div class="flex items-center gap-1 font-mono text-[10px] text-gray-400">
              <input
                v-model.number="cue.start"
                type="number"
                step="0.1"
                class="w-12 bg-gray-900 border border-gray-700 rounded px-1 py-0.5"
                @change="updateSubtitle"
              />
              <span>s -</span>
              <input
                v-model.number="cue.end"
                type="number"
                step="0.1"
                class="w-12 bg-gray-900 border border-gray-700 rounded px-1 py-0.5"
                @change="updateSubtitle"
              />
              <span>s</span>
            </div>

            <button
              class="text-gray-500 hover:text-red-400 p-0.5"
              @click="removeCue(idx)"
            >
              <UIcon name="i-heroicons-trash" class="w-3.5 h-3.5" />
            </button>
          </div>

          <input
            v-model="cue.text"
            type="text"
            placeholder="Subtitle text..."
            class="w-full bg-gray-900 border border-gray-700 rounded px-2 py-1 text-xs text-gray-200 focus:outline-none focus:border-indigo-500"
            @change="updateSubtitle"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useTimelineStore } from '~/stores/timeline'
import type { SubtitleCueSpec } from '~/types/spec'

const timelineStore = useTimelineStore()
const selectedClip = computed(() => timelineStore.selectedClip)

const fontSize = ref(32)
const textColor = ref('#FFFFFF')
const outlineWidth = ref(2)
const outlineColor = ref('#000000')
const hasBox = ref(false)
const cues = ref<SubtitleCueSpec[]>([])

watch(selectedClip, (clip) => {
  if (clip?.subtitles) {
    fontSize.value = clip.subtitles.font_size || 32
    textColor.value = clip.subtitles.font_color || '#FFFFFF'
    outlineWidth.value = clip.subtitles.outline_width ?? 2
    outlineColor.value = clip.subtitles.outline_color || '#000000'
    hasBox.value = clip.subtitles.box ?? false
    cues.value = clip.subtitles.cues || []
  }
}, { immediate: true })

function addCue() {
  cues.value.push({
    start: 0,
    end: 3.0,
    text: 'New Subtitle Caption',
  })
  updateSubtitle()
}

function removeCue(index: number) {
  cues.value.splice(index, 1)
  updateSubtitle()
}

function updateSubtitle() {
  if (!selectedClip.value) return
  selectedClip.value.subtitles = {
    font_size: fontSize.value,
    font_color: textColor.value,
    outline_width: outlineWidth.value,
    outline_color: outlineColor.value,
    box: hasBox.value,
    cues: cues.value,
  }
  timelineStore.pushHistoryState('Update Subtitles')
}
</script>
