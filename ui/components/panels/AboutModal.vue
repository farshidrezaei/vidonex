<template>
  <UModal
    v-model:open="isOpen"
    :title="$t('about.title') || 'About Vidonex Studio'"
    :ui="{
      content: 'sm:max-w-xl bg-gray-950 border border-gray-800 shadow-2xl rounded-2xl overflow-hidden'
    }"
  >
    <template #body>
      <div class="space-y-5 text-gray-200">
        <!-- App Header Card -->
        <div class="flex items-center gap-4 p-4 rounded-xl bg-gradient-to-br from-indigo-950/40 via-gray-900 to-gray-950 border border-indigo-500/20">
          <img src="/logo.png" alt="Vidonex Logo" class="w-14 h-14 rounded-2xl object-contain shadow-lg shadow-indigo-500/30 ring-1 ring-white/10" />
          <div class="flex-1">
            <div class="flex items-center gap-2.5">
              <h2 class="text-xl font-extrabold tracking-tight bg-gradient-to-r from-white via-gray-100 to-indigo-200 bg-clip-text text-transparent">
                Vidonex Studio
              </h2>
              <UBadge color="primary" variant="subtle" size="sm" class="font-mono text-[11px] font-semibold">
                v{{ currentVersion }}
              </UBadge>
            </div>
            <p class="text-xs text-gray-400 mt-0.5">
              {{ $t('app.tagline') || 'Declarative Video Composition & FFmpeg Filtergraph Workstation' }}
            </p>
          </div>
        </div>

        <!-- Update Status Card -->
        <div class="p-4 rounded-xl border transition-colors" :class="updateCardClasses">
          <div class="flex items-start justify-between gap-3">
            <div class="flex items-start gap-3 flex-1">
              <div class="mt-0.5 p-1.5 rounded-lg shrink-0" :class="updateIconBgClasses">
                <UIcon :name="updateIconName" class="w-5 h-5" :class="updateIconColorClasses" />
              </div>
              <div class="space-y-1.5 flex-1">
                <div class="flex items-center gap-2">
                  <h3 class="text-sm font-semibold text-white">
                    {{ updateStatusTitle }}
                  </h3>
                  <UBadge
                    v-if="hasUpdate && !isUpgradeComplete"
                    color="warning"
                    variant="solid"
                    size="xs"
                    class="animate-pulse text-[10px] font-mono font-bold"
                  >
                    {{ $t('about.new_badge') || 'NEW' }} v{{ latestVersion }}
                  </UBadge>
                </div>
                <p class="text-xs text-gray-400 leading-relaxed">
                  {{ updateStatusDescription }}
                </p>

                <!-- In-App Upgrade Progress Bar -->
                <div v-if="isUpgrading" class="mt-3 space-y-1.5 bg-gray-900/90 p-3 rounded-lg border border-indigo-500/30">
                  <div class="flex items-center justify-between text-xs">
                    <span class="text-indigo-300 font-medium flex items-center gap-1.5">
                      <UIcon name="i-heroicons-arrow-path" class="w-3.5 h-3.5 animate-spin text-indigo-400" />
                      {{ currentProgressLabel }}
                    </span>
                    <span class="font-mono text-indigo-200 font-bold">{{ upgradeProgress }}%</span>
                  </div>
                  <div class="w-full bg-gray-800 rounded-full h-2 overflow-hidden">
                    <div
                      class="bg-gradient-to-r from-indigo-500 to-sky-400 h-2 rounded-full transition-all duration-300"
                      :style="{ width: `${upgradeProgress}%` }"
                    />
                  </div>
                  <p class="text-[11px] text-gray-400 font-mono truncate">
                    {{ upgradeMessage }}
                  </p>
                </div>

                <!-- Upgrade Completed Box -->
                <div v-if="isUpgradeComplete" class="mt-3 p-3 rounded-lg bg-emerald-950/40 border border-emerald-500/40 space-y-2">
                  <div class="flex items-center gap-2 text-xs font-semibold text-emerald-300">
                    <UIcon name="i-heroicons-check-badge" class="w-4 h-4 text-emerald-400" />
                    <span>{{ $t('about.update_complete') || 'Update applied successfully!' }}</span>
                  </div>
                  <p class="text-[11px] text-emerald-200/80 leading-relaxed">
                    {{ $t('about.restart_notice') || 'New version installed. Please reload or restart.' }}
                  </p>
                  <UButton
                    color="success"
                    size="xs"
                    icon="i-heroicons-arrow-path"
                    class="font-semibold shadow-sm"
                    @click="reloadStudio"
                  >
                    {{ $t('about.refresh_page') || 'Reload Studio' }}
                  </UButton>
                </div>

                <!-- Upgrade Error Box -->
                <div v-if="upgradeError" class="mt-3 p-3 rounded-lg bg-rose-950/40 border border-rose-500/40 text-xs text-rose-300 space-y-1">
                  <div class="flex items-center gap-1.5 font-semibold">
                    <UIcon name="i-heroicons-exclamation-triangle" class="w-4 h-4 text-rose-400" />
                    <span>Failed to install update</span>
                  </div>
                  <p class="text-[11px] text-rose-300/80 font-mono break-all">{{ upgradeError }}</p>
                </div>

                <!-- Release Notes Preview -->
                <div v-if="releaseNotes && hasUpdate && !isUpgrading && !isUpgradeComplete" class="mt-2 text-xs text-gray-300 bg-gray-900/80 p-2.5 rounded-lg border border-gray-800/80 max-h-24 overflow-y-auto font-mono text-[11px]">
                  {{ releaseNotes }}
                </div>
              </div>
            </div>

            <!-- Action: In-App Upgrade or Manual Check -->
            <div class="flex flex-col items-end gap-1.5 shrink-0">
              <template v-if="hasUpdate && !isUpgradeComplete">
                <UButton
                  icon="i-heroicons-bolt"
                  color="primary"
                  size="sm"
                  :loading="isUpgrading"
                  class="font-semibold shadow-md shadow-indigo-500/30"
                  @click="confirmAndUpgrade"
                >
                  {{ $t('about.install_update') || 'Install Update' }}
                </UButton>
                <UButton
                  color="neutral"
                  variant="ghost"
                  size="xs"
                  class="text-[11px] text-gray-400 hover:text-gray-200"
                  @click="openReleaseUrl"
                >
                  {{ $t('about.view_release_page') || 'GitHub' }}
                </UButton>
              </template>
              <template v-else-if="!isUpgradeComplete">
                <UButton
                  :loading="isChecking"
                  icon="i-heroicons-arrow-path"
                  color="neutral"
                  variant="subtle"
                  size="xs"
                  @click="checkForUpdates(true)"
                >
                  {{ $t('about.check_updates') || 'Check for Updates' }}
                </UButton>
                <span v-if="lastCheckedAt" class="text-[10px] text-gray-500 font-mono">
                  {{ formatCheckedTime }}
                </span>
              </template>
            </div>
          </div>
        </div>


        <!-- System Runtime Diagnostics -->
        <div class="space-y-2">
          <div class="flex items-center justify-between">
            <h4 class="text-xs font-bold uppercase tracking-wider text-gray-400 flex items-center gap-1.5">
              <UIcon name="i-heroicons-cpu-chip" class="w-4 h-4 text-indigo-400" />
              <span>{{ $t('about.system_specs') || 'System & Engine Diagnostics' }}</span>
            </h4>
          </div>

          <div class="grid grid-cols-2 gap-2 text-xs font-mono">
            <div class="bg-gray-900/60 p-2.5 rounded-lg border border-gray-800/70 flex flex-col justify-between">
              <span class="text-gray-500 text-[11px]">{{ $t('about.os_arch') || 'Operating System / Arch' }}</span>
              <span class="text-gray-200 font-medium truncate mt-1">
                {{ systemDiagnostics?.os || 'linux' }} / {{ systemDiagnostics?.architecture || 'amd64' }} ({{ systemDiagnostics?.cpu_count || 1 }} {{ $t('about.cores') || 'Cores' }})
              </span>
            </div>

            <div class="bg-gray-900/60 p-2.5 rounded-lg border border-gray-800/70 flex flex-col justify-between">
              <span class="text-gray-500 text-[11px]">{{ $t('about.hardware_accel') || 'Hardware Acceleration' }}</span>
              <span class="text-emerald-400 font-medium truncate mt-1 flex items-center gap-1">
                <UIcon name="i-heroicons-bolt" class="w-3.5 h-3.5 shrink-0" />
                {{ hardwareAcceleratorLabel }}
              </span>
            </div>

            <div class="col-span-2 bg-gray-900/60 p-2.5 rounded-lg border border-gray-800/70 space-y-1.5">
              <div class="flex items-center justify-between">
                <span class="text-gray-500 text-[11px]">{{ $t('about.ffmpeg_engine') || 'FFmpeg Engine' }}</span>
                <span class="text-[10px] text-emerald-400 font-sans font-medium flex items-center gap-1">
                  <span class="w-1.5 h-1.5 rounded-full bg-emerald-400 animate-pulse"></span>
                  {{ $t('about.ffmpeg_active') || 'Active' }}
                </span>
              </div>
              <p class="text-gray-300 text-[11px] truncate" :title="systemDiagnostics?.ffmpeg_version">
                {{ systemDiagnostics?.ffmpeg_version || ($t('about.ffmpeg_detected') || 'FFmpeg system binary detected') }}
              </p>
            </div>
          </div>
        </div>

        <!-- Links & Copyright -->
        <div class="pt-3 border-t border-gray-800/80 flex items-center justify-between text-xs text-gray-500">
          <div class="flex items-center gap-3">
            <a
              href="https://github.com/farshidrezaei/vidonex"
              target="_blank"
              rel="noopener noreferrer"
              class="hover:text-gray-300 transition flex items-center gap-1"
            >
              <UIcon name="i-heroicons-code-bracket" class="w-3.5 h-3.5" />
              <span>{{ $t('about.github') || 'GitHub' }}</span>
            </a>
            <a
              href="https://farshidrezaei.github.io/vidonex/"
              target="_blank"
              rel="noopener noreferrer"
              class="hover:text-gray-300 transition flex items-center gap-1"
            >
              <UIcon name="i-heroicons-book-open" class="w-3.5 h-3.5" />
              <span>{{ $t('about.documentation') || 'Documentation' }}</span>
            </a>
            <a
              href="https://github.com/farshidrezaei/vidonex/blob/main/CHANGELOG.md"
              target="_blank"
              rel="noopener noreferrer"
              class="hover:text-gray-300 transition flex items-center gap-1"
            >
              <UIcon name="i-heroicons-document-text" class="w-3.5 h-3.5" />
              <span>{{ $t('about.changelog') || 'Changelog' }}</span>
            </a>
          </div>

          <span>{{ $t('about.license') || 'MIT License' }}</span>
        </div>
      </div>
    </template>
  </UModal>
</template>

<script setup lang="ts">
import { isAboutModalOpen, useAppUpdate } from '~/composables/useAppUpdate'

const { t } = useI18n()
const {
  isChecking,
  hasUpdate,
  currentVersion,
  latestVersion,
  releaseNotes,
  systemDiagnostics,
  lastCheckedAt,
  isUpgrading,
  upgradeProgress,
  upgradeStage,
  upgradeMessage,
  upgradeError,
  isUpgradeComplete,
  checkForUpdates,
  startInAppUpdate,
  openReleaseUrl,
} = useAppUpdate()

const isOpen = computed({
  get: () => isAboutModalOpen.value,
  set: (val: boolean) => {
    isAboutModalOpen.value = val
  },
})

const updateCardClasses = computed(() => {
  if (isUpgradeComplete.value) {
    return 'bg-emerald-950/20 border-emerald-500/40'
  }
  if (hasUpdate.value) {
    return 'bg-amber-950/20 border-amber-500/30'
  }
  return 'bg-gray-900/40 border-gray-800/80'
})

const updateIconBgClasses = computed(() => {
  if (isUpgradeComplete.value) return 'bg-emerald-500/20'
  if (hasUpdate.value) return 'bg-amber-500/20'
  return 'bg-emerald-500/20'
})

const updateIconName = computed(() => {
  if (isUpgradeComplete.value) return 'i-heroicons-check-circle'
  if (hasUpdate.value) return 'i-heroicons-arrow-up-circle'
  return 'i-heroicons-check-circle'
})

const updateIconColorClasses = computed(() => {
  if (isUpgradeComplete.value) return 'text-emerald-400'
  if (hasUpdate.value) return 'text-amber-400'
  return 'text-emerald-400'
})

const updateStatusTitle = computed(() => {
  if (isUpgradeComplete.value) {
    return t('about.update_complete') || 'Update Complete'
  }
  if (hasUpdate.value) {
    return t('about.update_available') || `Update Available: v${latestVersion.value}`
  }
  return t('about.up_to_date') || 'Vidonex is up to date'
})

const updateStatusDescription = computed(() => {
  if (isUpgradeComplete.value) {
    return t('about.restart_notice') || 'Vidonex has been upgraded to the latest version.'
  }
  if (hasUpdate.value) {
    return t('about.update_available_desc') || `A new release (v${latestVersion.value}) is available with improvements and new features.`
  }
  return t('about.up_to_date_desc') || `You are running the latest version v${currentVersion.value}.`
})

const currentProgressLabel = computed(() => {
  switch (upgradeStage.value) {
    case 'downloading':
      return t('about.downloading') || 'Downloading update...'
    case 'extracting':
      return t('about.extracting') || 'Extracting files...'
    case 'replacing':
      return t('about.replacing') || 'Applying update...'
    default:
      return 'Processing update...'
  }
})

function confirmAndUpgrade() {
  const confirmMsg = (t('about.confirm_update_desc') || 'Do you want to automatically download and install Vidonex v{version}?').replace('{version}', latestVersion.value)
  if (typeof window !== 'undefined' && window.confirm(confirmMsg)) {
    startInAppUpdate()
  }
}

function reloadStudio() {
  if (typeof window !== 'undefined') {
    window.location.reload()
  }
}

const hardwareAcceleratorLabel = computed(() => {
  const accels = systemDiagnostics.value?.available_accelerators || []
  if (accels.length > 0) {
    return accels.map((a) => a.toUpperCase()).join(', ')
  }
  return t('about.software_cpu') || 'Software (CPU)'
})

const formatCheckedTime = computed(() => {
  if (!lastCheckedAt.value) return ''
  return lastCheckedAt.value.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
})
</script>
