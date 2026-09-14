import { defineConfig } from 'vitepress'

export default defineConfig({
  title: 'Vidonex',
  description: 'Declarative Video Composition & FFmpeg Filtergraph Compiler Engine in Go',
  base: '/vidonex/',
  head: [
    ['link', { rel: 'icon', type: 'image/png', href: '/vidonex/logo.png' }],
    ['meta', { name: 'theme-color', content: '#6366f1' }],
    ['meta', { property: 'og:type', content: 'website' }],
    ['meta', { property: 'og:locale', content: 'en' }],
    ['meta', { property: 'og:title', content: 'Vidonex - Modern Video Engine in Go' }],
    ['meta', { property: 'og:site_name', content: 'Vidonex' }],
    ['meta', { property: 'og:image', content: 'https://raw.githubusercontent.com/farshidrezaei/vidonex/main/.github/assets/banner.png' }],
    ['meta', { property: 'og:description', content: 'Declarative Video Composition & FFmpeg Filtergraph Compiler Engine in Go with Modern Web & Desktop Studio' }],
    ['meta', { name: 'twitter:card', content: 'summary_large_image' }],
  ],
  themeConfig: {
    logo: '/logo.png',
    siteTitle: 'Vidonex',
    search: {
      provider: 'local'
    },
    nav: [
      { text: 'Guide', link: '/guide/introduction' },
      { text: 'API Reference', link: '/api/' },
      { text: 'Recipes', link: '/recipes/' },
      { text: 'Studio', link: '/guide/desktop-studio' },
      {
        text: 'v1.0.0',
        items: [
          { text: 'Release Notes', link: 'https://github.com/farshidrezaei/vidonex/releases' },
          { text: 'Contributing', link: 'https://github.com/farshidrezaei/vidonex/blob/main/CONTRIBUTING.md' },
          { text: 'Roadmap & Architecture', link: '/guide/architecture' }
        ]
      }
    ],
    sidebar: {
      '/guide/': [
        {
          text: 'Getting Started',
          items: [
            { text: 'Introduction', link: '/guide/introduction' },
            { text: 'Quick Start', link: '/guide/quick-start' },
            { text: 'Architecture & DAG IR', link: '/guide/architecture' },
            { text: 'CLI Automation', link: '/guide/cli' },
            { text: 'Declarative YAML Spec', link: '/guide/yaml-spec' },
            { text: 'Desktop Studio Workstation', link: '/guide/desktop-studio' }
          ]
        }
      ],
      '/api/': [
        {
          text: 'Go SDK Reference',
          items: [
            { text: 'API Overview', link: '/api/' },
            { text: 'Timeline & AST', link: '/api/timeline' },
            { text: 'Compiler Engine', link: '/api/compiler' },
            { text: 'Filtergraph & Optimization', link: '/api/filtergraph' },
            { text: 'Visual Effects', link: '/api/effects' },
            { text: 'Audio & Sidechain Ducking', link: '/api/ducking-and-audio' },
            { text: 'ChromaKey & Despill', link: '/api/chromakey' }
          ]
        }
      ],
      '/recipes/': [
        {
          text: 'Runnable Recipes',
          items: [
            { text: 'Recipe Catalog', link: '/recipes/' },
            { text: 'Picture-in-Picture (PiP)', link: '/recipes/picture-in-picture' },
            { text: 'Sidechain Audio Ducking', link: '/recipes/audio-ducking' },
            { text: 'Animated Neon Waveforms', link: '/recipes/animated-waveforms' },
            { text: 'Green Screen & Despill', link: '/recipes/green-screen' },
            { text: 'Burn-in Subtitles & Captions', link: '/recipes/subtitles' }
          ]
        }
      ]
    },
    socialLinks: [
      { icon: 'github', link: 'https://github.com/farshidrezaei/vidonex' }
    ],
    footer: {
      message: 'Released under the MIT License.',
      copyright: 'Copyright © 2026 Vidonex Contributors'
    },
    editLink: {
      pattern: 'https://github.com/farshidrezaei/vidonex/edit/main/docs/:path',
      text: 'Edit this page on GitHub'
    }
  }
})
