// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2025-01-01',
  future: {
    compatibilityVersion: 4,
  },
  modules: [
    '@nuxt/ui',
    '@nuxtjs/i18n',
    '@pinia/nuxt',
    '@vueuse/nuxt',
  ],
  devtools: { enabled: true },
  ssr: false, // SPA Mode for smooth client-side real-time timeline editing and canvas gizmo
  app: {
    head: {
      title: 'Vidonyx Studio - Video Composition Workstation',
      meta: [
        { name: 'description', content: 'Professional Declarative Video Composition & FFmpeg Filtergraph Editor' },
        { name: 'viewport', content: 'width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no' },
      ],
      link: [
        { rel: 'preconnect', href: 'https://fonts.googleapis.com' },
        { rel: 'preconnect', href: 'https://fonts.gstatic.com', crossorigin: '' },
        { rel: 'stylesheet', href: 'https://fonts.googleapis.com/css2?family=Vazirmatn:wght@300;400;500;600;700;800&family=JetBrains+Mono:wght@400;500;600&family=Inter:wght@300;400;500;600;700&display=swap' },
      ],
    },
  },
  css: ['~/assets/css/main.css'],
  i18n: {
    locales: [
      { code: 'en', name: 'English', dir: 'ltr', file: 'en.json' },
      { code: 'fa', name: 'فارسی', dir: 'rtl', file: 'fa.json' },
    ],
    defaultLocale: 'fa',
    lazy: false,
    langDir: 'locales',
    strategy: 'no_prefix',
  },
})
