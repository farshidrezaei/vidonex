<template>
  <UModal
    v-model:open="isOpen"
    :title="$t('shortcuts.title') || 'Keyboard Shortcuts'"
    :ui="{ content: 'max-w-2xl' }"
  >
    <template #body>
      <div class="space-y-4">
        <!-- Search Input -->
        <div class="relative">
          <UInput
            v-model="searchQuery"
            icon="i-heroicons-magnifying-glass"
            placeholder="Search shortcut (e.g. undo, split, play, zoom)..."
            size="sm"
            class="w-full"
          />
        </div>

        <!-- Shortcuts Categories List -->
        <div class="max-h-[60vh] overflow-y-auto space-y-5 pr-1 text-sm">
          <div
            v-for="cat in filteredCategories"
            :key="cat.name"
            class="bg-gray-950/60 rounded-xl p-3 border border-gray-800/80"
          >
            <div class="flex items-center gap-2 mb-2.5 pb-1.5 border-b border-gray-800">
              <UIcon :name="cat.icon" class="w-4 h-4 text-indigo-400" />
              <h4 class="text-xs font-bold uppercase tracking-wider text-gray-200">{{ cat.name }}</h4>
            </div>

            <div class="space-y-2">
              <div
                v-for="item in cat.shortcuts"
                :key="item.description"
                class="flex items-center justify-between py-1 px-1 rounded hover:bg-gray-900/60 transition"
              >
                <div class="flex flex-col">
                  <span class="text-xs font-medium text-gray-200">{{ item.description }}</span>
                  <span v-if="item.hint" class="text-[10px] text-gray-500">{{ item.hint }}</span>
                </div>

                <div class="flex items-center gap-1">
                  <template v-for="(k, idx) in item.keys" :key="idx">
                    <span v-if="k === '+' || k === 'or'" class="text-[10px] text-gray-500 mx-0.5">{{ k }}</span>
                    <UKbd v-else size="sm" class="font-mono">{{ k }}</UKbd>
                  </template>
                </div>
              </div>
            </div>
          </div>

          <div v-if="filteredCategories.length === 0" class="text-center py-8 text-gray-500 text-xs">
            No shortcuts found matching "{{ searchQuery }}"
          </div>
        </div>
      </div>
    </template>
  </UModal>
</template>

<script setup lang="ts">
const props = defineProps<{
  modelValue: boolean
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
}>()

const isOpen = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val),
})

const searchQuery = ref('')

const categories = [
  {
    name: 'History & Timeline Editing',
    icon: 'i-heroicons-scissors',
    shortcuts: [
      { description: 'Undo Last Action', hint: 'Reverts latest edit', keys: ['Ctrl', '+', 'Z'] },
      { description: 'Redo Action', hint: 'Re-applies reverted action', keys: ['Ctrl', '+', 'Y', 'or', 'Ctrl', '+', 'Shift', '+', 'Z'] },
      { description: 'Duplicate Selected Clip', hint: 'Clones without overlap', keys: ['Ctrl', '+', 'D'] },
      { description: 'Split Clip at Playhead', hint: 'Cuts clip under red needle', keys: ['S', 'or', 'C'] },
      { description: 'Delete Selected Clip', hint: 'Removes active clip', keys: ['Delete', 'or', 'Backspace'] },
      { description: 'Deselect Active Clip', hint: 'Clears selection', keys: ['Esc'] },
    ],
  },
  {
    name: 'Canvas & Viewport Navigation',
    icon: 'i-heroicons-cursor-arrow-rays',
    shortcuts: [
      { description: 'Hand / Pan Tool', hint: 'Pan viewport canvas freely', keys: ['Space', '+', 'Drag'] },
      { description: 'Viewport Zoom', hint: 'Smoothly zoom in / out', keys: ['Mouse Wheel'] },
      { description: 'Reset Zoom & Pan', hint: 'Return to fit view', keys: ['Double Click'] },
      { description: 'Nudge Element (1px)', hint: 'Precise movement in canvas', keys: ['↑', '↓', '←', '→'] },
      { description: 'Nudge Element (10px)', hint: 'Fast movement in canvas', keys: ['Shift', '+', '↑ / ↓ / ← / →'] },
    ],
  },
  {
    name: 'Playback & Navigation',
    icon: 'i-heroicons-play-circle',
    shortcuts: [
      { description: 'Step 1 Frame Backward', hint: 'Precision frame jump left (no clip selected)', keys: ['←', 'or', 'J'] },
      { description: 'Step 1 Frame Forward', hint: 'Precision frame jump right (no clip selected)', keys: ['→', 'or', 'L'] },
      { description: 'Jump to Beginning (0s)', hint: 'Seeks to start of timeline', keys: ['Home'] },
      { description: 'Toggle Magnetic Snapping', hint: 'Enables/disables magnet alignment', keys: ['N'] },
    ],
  },
  {
    name: 'Zoom & Cheatsheet',
    icon: 'i-heroicons-magnifying-glass-plus',
    shortcuts: [
      { description: 'Zoom In Timeline', hint: 'Expands timeline scale', keys: ['+', 'or', '='] },
      { description: 'Zoom Out Timeline', hint: 'Shrinks timeline scale', keys: ['-'] },
      { description: 'Open Shortcuts Cheatsheet', hint: 'Shows this guide', keys: ['?', 'or', 'Ctrl', '+', '/'] },
    ],
  },
]

const filteredCategories = computed(() => {
  const q = searchQuery.value.trim().toLowerCase()
  if (!q) return categories

  return categories
    .map((cat) => {
      const filtered = cat.shortcuts.filter(
        (s) =>
          s.description.toLowerCase().includes(q) ||
          s.keys.some((k) => k.toLowerCase().includes(q)) ||
          (s.hint && s.hint.toLowerCase().includes(q))
      )
      return { ...cat, shortcuts: filtered }
    })
    .filter((cat) => cat.shortcuts.length > 0)
})
</script>
