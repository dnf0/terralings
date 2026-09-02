---
title: Interactive Browser Playground
description: Learn Terraform and OpenTofu directly in your browser with Pyodide WebAssembly and client-side state persistence.
---

# 🏗️ Interactive In-Browser Learning Platform

Practice Terraform and OpenTofu exercises directly in your browser with **zero local installations, zero CLI dependencies, and zero cloud credentials**.

The playground runs an in-memory **HCL AST parser and rule validator** compiled to **WebAssembly (Pyodide)**, with local state persistence, progressive hints, interactive diffs, and an embedded Monaco editor.

---

<div id="terralings-app">
  <div class="terralings-loading-screen">
    <div class="spinner"></div>
    <p>⚡ Initializing Python 3.12 WebAssembly Runtime & Terralings Curriculum...</p>
  </div>
</div>

---

## ⌨️ Keyboard Shortcuts

| Shortcut | Action | Description |
|---|---|---|
| <kbd>Ctrl</kbd> + <kbd>Enter</kbd> / <kbd>⌘</kbd> + <kbd>Enter</kbd> | **Run Solution** | Evaluates active exercise in Pyodide WebAssembly |
| <kbd>Alt</kbd> + <kbd>→</kbd> | **Next Exercise** | Navigates to the subsequent exercise in the curriculum |
| <kbd>Alt</kbd> + <kbd>←</kbd> | **Previous Exercise** | Navigates to the preceding exercise |
| <kbd>H</kbd> | **Reveal Hint** | Unlocks the next progressive hint tier |
| <kbd>F11</kbd> | **Toggle Fullscreen** | Expands playground into a full-bleed 100vw × 100vh distraction-free workspace |

---

## 💾 Client-Side State & Data Portability

Your learning progress, completed exercises, and editor drafts are stored locally in your browser's `localStorage` (`terralings_learning_state_v1`).

- **Auto-Save:** Drafts are saved automatically as you type (debounced 300ms).
- **Export Progress:** Click **💾 Export** in the syllabus sidebar to download your complete progress snapshot as `terralings-progress-YYYY-MM-DD.json`.
- **Import Progress:** Click **📂 Import** to restore your progress across different browsers or machines.
- **Reset:** Use **↺ Reset Code** to restore the current exercise starter template, or **🗑️ Reset** to clear all stored progress.
