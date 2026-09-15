export interface ConfirmOptions {
  title: string
  message: string
  confirmText?: string
  cancelText?: string
  isDanger?: boolean
  icon?: string
}

const isOpen = ref(false)
const options = ref<ConfirmOptions>({
  title: '',
  message: '',
  confirmText: 'Delete',
  cancelText: 'Cancel',
  isDanger: true,
  icon: 'i-heroicons-exclamation-triangle',
})

let resolvePromise: ((value: boolean) => void) | null = null

export function useConfirmDialog() {
  function confirm(opts: ConfirmOptions): Promise<boolean> {
    options.value = {
      confirmText: 'Delete',
      cancelText: 'Cancel',
      isDanger: true,
      icon: 'i-heroicons-exclamation-triangle',
      ...opts,
    }
    isOpen.value = true

    return new Promise<boolean>((resolve) => {
      resolvePromise = resolve
    })
  }

  function onConfirm() {
    isOpen.value = false
    if (resolvePromise) {
      resolvePromise(true)
      resolvePromise = null
    }
  }

  function onCancel() {
    isOpen.value = false
    if (resolvePromise) {
      resolvePromise(false)
      resolvePromise = null
    }
  }

  return {
    isOpen,
    options,
    confirm,
    onConfirm,
    onCancel,
  }
}
