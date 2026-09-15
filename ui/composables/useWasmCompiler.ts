export interface WasmValidationResult {
  valid: boolean
  error?: string
  duration?: number
  tracks?: number
}

export interface WasmCompilationResult {
  success: boolean
  error?: string
  command?: string
  filterComplex?: string
  mermaid?: string
  inputs?: string[]
  inputFiles?: string[]
  duration?: number
}

declare global {
  interface Window {
    Go?: any
    vidonex?: {
      validate: (content: string) => WasmValidationResult
      compile: (content: string, outputPath?: string) => WasmCompilationResult
      mermaid: (content: string) => { success: boolean; mermaid?: string; error?: string }
      version: () => string
    }
  }
}

export const isWasmLoaded = ref(false)
export const isWasmLoading = ref(false)
export const wasmLoadError = ref<string | null>(null)
export const wasmVersion = ref<string>('')

export function useWasmCompiler() {
  async function initWasm(): Promise<boolean> {
    if (isWasmLoaded.value && window.vidonex) return true
    if (isWasmLoading.value) return false

    if (typeof window === 'undefined') return false

    isWasmLoading.value = true
    wasmLoadError.value = null

    try {
      // 1. Load Go runtime glue script if not already present
      if (!window.Go) {
        await new Promise<void>((resolve, reject) => {
          const script = document.createElement('script')
          script.src = '/wasm_exec.js'
          script.onload = () => resolve()
          script.onerror = () => reject(new Error('Failed to load /wasm_exec.js'))
          document.head.appendChild(script)
        })
      }

      // 2. Instantiate WebAssembly binary
      const go = new window.Go()
      const wasmResponse = await fetch('/vidonex.wasm')
      if (!wasmResponse.ok) {
        throw new Error(`Failed to fetch /vidonex.wasm: ${wasmResponse.statusText}`)
      }

      const wasmBuffer = await wasmResponse.arrayBuffer()
      const wasmModule = await WebAssembly.instantiate(wasmBuffer, go.importObject)

      // 3. Run WebAssembly process in background
      go.run(wasmModule.instance)

      if (window.vidonex) {
        isWasmLoaded.value = true
        wasmVersion.value = window.vidonex.version()
        return true
      } else {
        throw new Error('WebAssembly initialized but globalThis.vidonex is undefined')
      }
    } catch (err: any) {
      wasmLoadError.value = err?.message || String(err)
      return false
    } finally {
      isWasmLoading.value = false
    }
  }

  function validate(specContent: string): WasmValidationResult {
    if (!window.vidonex) {
      return { valid: false, error: 'WebAssembly engine not initialized' }
    }
    return window.vidonex.validate(specContent)
  }

  function compile(specContent: string, outputPath = 'output.mp4'): WasmCompilationResult {
    if (!window.vidonex) {
      return { success: false, error: 'WebAssembly engine not initialized' }
    }
    return window.vidonex.compile(specContent, outputPath)
  }

  function mermaid(specContent: string): { success: boolean; mermaid?: string; error?: string } {
    if (!window.vidonex) {
      return { success: false, error: 'WebAssembly engine not initialized' }
    }
    return window.vidonex.mermaid(specContent)
  }

  return {
    isWasmLoaded,
    isWasmLoading,
    wasmLoadError,
    wasmVersion,
    initWasm,
    validate,
    compile,
    mermaid,
  }
}
