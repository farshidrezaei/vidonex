# GitHub Repository SEO & Optimization Guide

To maximize discoverability and SEO rankings for **Vidonex** on GitHub and search engines (Google, DuckDuckGo, Bing), configure the repository's metadata in the GitHub web interface as follows:

---

## 1. Repository "About" Section (Main Page Sidebar)

Click the gear icon ⚙️ next to "About" on the repository main page:

- **Description**:
  > 🎬 The Next-Gen Non-Linear Video Editor & Declarative FFmpeg Filtergraph Compiler Engine in Go + Modern Desktop Studio (Vue 3 / Nuxt UI)

- **Website**:
  > https://farshidrezaei.github.io/vidonex/

- **Topics (Tags)** — *Enter all 20 tags to maximize search index coverage*:
  ```text
  video-editor
  ffmpeg
  golang
  video-processing
  remotion-alternative
  filtergraph
  motion-graphics
  video-generation
  timeline-editor
  vue3
  nuxt-ui
  wails
  video-automation
  render-engine
  creative-tools
  audio-ducking
  chromakey
  audiogram
  desktop-app
  declarative-ui
  ```

- **Include in the home page**:
  - [x] Releases
  - [x] Packages
  - [x] Environments (GitHub Pages)

---

## 2. GitHub Social Preview (OpenGraph Banner)

1. Go to **Settings** -> **General** -> **Social preview**.
2. Click **Edit** -> **Upload an image**.
3. Upload `.github/assets/banner.png` (or `.github/assets/studio-preview.png`).
   - Dimensions: 1280 × 640 px.

---

## 3. GitHub Pages Settings

1. Go to **Settings** -> **Pages**.
2. **Source**: Select **GitHub Actions** (The `.github/workflows/docs.yml` workflow will automatically build and publish the docs).
3. The documentation will be live at:
   `https://farshidrezaei.github.io/vidonex/`

---

## 4. GitHub Discussions & Features

Go to **Settings** → **General** → **Features**:
- [x] **Discussions**: Enable to foster a community and encourage Q&A.
  - Recommended categories: *Announcements*, *General*, *Ideas & Feature Requests*, *Show and Tell* (where users show videos they generated).
- [x] **Issues**: Enable (templates are already configured in `.github/ISSUE_TEMPLATE/`).
