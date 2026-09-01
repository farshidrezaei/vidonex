<template>
  <div
    class="absolute top-1/2 -translate-y-1/2 -translate-x-1/2 z-20 group"
    :style="{ left: `${positionX}px` }"
  >
    <UDropdownMenu :items="transitionMenuItems">
      <button
        class="w-5 h-5 rounded-full flex items-center justify-center transition shadow-md"
        :class="transition ? 'bg-amber-500 text-black' : 'bg-gray-800 hover:bg-gray-700 text-gray-400 border border-gray-700'"
        :title="transition ? `Transition: ${transition.type} (${transition.duration}s)` : $t('timeline.add_transition')"
      >
        <UIcon :name="transition ? 'i-heroicons-sparkles' : 'i-heroicons-plus'" class="w-3 h-3" />
      </button>
    </UDropdownMenu>
  </div>
</template>

<script setup lang="ts">
import { useTimelineStore } from '~/stores/timeline'
import type { TransitionSpec, ClipSpec } from '~/types/spec'

const props = defineProps<{
  trackId: string
  fromClip: ClipSpec
  toClip: ClipSpec
  transition?: TransitionSpec
}>()

const { t } = useI18n()
const timelineStore = useTimelineStore()

const positionX = computed(() => {
  const fromEnd = (Number(props.fromClip.start) || 0) + (Number(props.fromClip.duration) || 0)
  return fromEnd * timelineStore.pixelsPerSecond
})

const transitionMenuItems = computed(() => [
  [
    {
      label: t('timeline.transitions.dissolve'),
      icon: 'i-heroicons-sparkles',
      onSelect: () => addOrUpdateTransition('dissolve'),
    },
    {
      label: t('timeline.transitions.fade'),
      icon: 'i-heroicons-eye-slash',
      onSelect: () => addOrUpdateTransition('fade'),
    },
    {
      label: t('timeline.transitions.wipeleft'),
      icon: 'i-heroicons-arrow-left',
      onSelect: () => addOrUpdateTransition('wipeleft'),
    },
    {
      label: t('timeline.transitions.wiperight'),
      icon: 'i-heroicons-arrow-right',
      onSelect: () => addOrUpdateTransition('wiperight'),
    },
    {
      label: t('timeline.transitions.circleopen'),
      icon: 'i-heroicons-lifebuoy',
      onSelect: () => addOrUpdateTransition('circleopen'),
    },
  ],
  ...(props.transition
    ? [
        [
          {
            label: 'Remove Transition',
            icon: 'i-heroicons-trash',
            onSelect: removeTransition,
          },
        ],
      ]
    : []),
])

function addOrUpdateTransition(type: string) {
  timelineStore.addTransition(props.trackId, props.fromClip.id, props.toClip.id, type, 1.0)
}

function removeTransition() {
  const track = timelineStore.tracks.find((t) => t.id === props.trackId)
  if (track && track.transitions) {
    track.transitions = track.transitions.filter((t) => !(t.from === props.fromClip.id && t.to === props.toClip.id))
    timelineStore.pushHistoryState('Remove transition')
  }
}
</script>
