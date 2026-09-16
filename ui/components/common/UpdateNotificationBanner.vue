<template>
  <Transition
    enter-active-class="transition duration-300 ease-out"
    enter-from-class="transform -translate-y-4 opacity-0"
    enter-to-class="transform translate-y-0 opacity-100"
    leave-active-class="transition duration-200 ease-in"
    leave-from-class="transform translate-y-0 opacity-100"
    leave-to-class="transform -translate-y-4 opacity-0"
  >
    <div
      v-if="hasUpdate && !isUpdateBannerDismissed"
      class="bg-gradient-to-r from-indigo-900/90 via-purple-900/90 to-indigo-900/90 backdrop-blur border-b border-indigo-500/30 px-4 py-1.5 flex items-center justify-between text-xs text-white z-40 relative shadow-lg"
    >
      <div class="flex items-center gap-2.5">
        <span class="flex h-2 w-2 relative">
          <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-amber-400 opacity-75"></span>
          <span class="relative inline-flex rounded-full h-2 w-2 bg-amber-500"></span>
        </span>
        <span class="font-medium">
          {{ $t('about.update_banner', { version: latestVersion }) || `Vidonex Studio v${latestVersion} is available!` }}
        </span>
      </div>

      <div class="flex items-center gap-2">
        <UButton
          size="xs"
          color="neutral"
          variant="ghost"
          class="text-xs text-indigo-200 hover:text-white"
          @click="openAbout"
        >
          {{ $t('about.view_details') || 'What\'s New' }}
        </UButton>

        <UButton
          size="xs"
          color="primary"
          icon="i-heroicons-arrow-down-tray"
          class="font-semibold shadow-sm text-xs bg-indigo-500 hover:bg-indigo-600"
          @click="openReleaseUrl"
        >
          {{ $t('about.update_now') || 'Update Now' }}
        </UButton>

        <UButton
          size="xs"
          color="neutral"
          variant="ghost"
          icon="i-heroicons-x-mark"
          class="text-indigo-300 hover:text-white p-1"
          @click="isUpdateBannerDismissed = true"
        />
      </div>
    </div>
  </Transition>
</template>

<script setup lang="ts">
import { useAppUpdate, isAboutModalOpen, isUpdateBannerDismissed } from '~/composables/useAppUpdate'

const { hasUpdate, latestVersion, openReleaseUrl } = useAppUpdate()

function openAbout() {
  isAboutModalOpen.value = true
}
</script>
