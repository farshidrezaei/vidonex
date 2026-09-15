<template>
  <UModal
    v-model:open="isOpen"
    :title="options.title"
    :ui="{ content: 'sm:max-w-md bg-gray-950 border border-gray-800' }"
  >
    <template #body>
      <div class="flex items-start gap-4 py-2">
        <div
          class="w-10 h-10 rounded-xl flex items-center justify-center shrink-0"
          :class="options.isDanger ? 'bg-red-500/10 text-red-400 border border-red-500/20' : 'bg-indigo-500/10 text-indigo-400 border border-indigo-500/20'"
        >
          <UIcon :name="options.icon || 'i-heroicons-exclamation-triangle'" class="w-5 h-5" />
        </div>
        <div class="flex-1 text-sm text-gray-300 leading-relaxed pt-1">
          {{ options.message }}
        </div>
      </div>
    </template>

    <template #footer>
      <div class="flex items-center justify-end gap-2 w-full">
        <UButton
          color="neutral"
          variant="ghost"
          size="sm"
          @click="onCancel"
        >
          {{ options.cancelText || $t('confirm.cancel') }}
        </UButton>
        <UButton
          :color="options.isDanger ? 'error' : 'primary'"
          size="sm"
          class="font-semibold"
          @click="onConfirm"
        >
          {{ options.confirmText || $t('confirm.delete') }}
        </UButton>
      </div>
    </template>
  </UModal>
</template>

<script setup lang="ts">
import { useConfirmDialog } from '~/composables/useConfirmDialog'

const { isOpen, options, onConfirm, onCancel } = useConfirmDialog()
</script>
