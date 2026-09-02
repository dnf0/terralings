/**
 * Terralings WebAssembly Playground UI Controller & State Engine
 *
 * Full 56-exercise browser learning environment powered by Pyodide WebAssembly.
 * Features client-side localStorage persistence, interactive split-pane syllabus sidebar,
 * real-time search & filters, progressive hints, side-by-side solution diffs, and progress backup.
 */

(function () {
  "use strict";

  const STORAGE_KEY = "terralings_learning_state_v1";

  /**
   * ==========================================================================
   * 1. Storage & State Management Layer
   * ==========================================================================
   */
  const TerralingsStorage = {
    state: {
      exercises: {},
      completedCount: 0,
      lastActiveExerciseId: null,
      lastSyncTimestamp: null,
    },

    init() {
      try {
        const raw = localStorage.getItem(STORAGE_KEY);
        if (raw) {
          const parsed = JSON.parse(raw);
          if (parsed && typeof parsed === "object") {
            this.state = Object.assign(this.state, parsed);
          }
        }
      } catch (err) {
        console.warn("Could not load Terralings state from localStorage:", err);
      }
      this.recalculateTotals();
    },

    persist() {
      try {
        this.state.lastSyncTimestamp = new Date().toISOString();
        localStorage.setItem(STORAGE_KEY, JSON.stringify(this.state));
      } catch (err) {
        console.warn("Could not persist Terralings state to localStorage:", err);
      }
    },

    getExerciseState(exerciseId, starterCode = "") {
      if (!this.state.exercises[exerciseId]) {
        this.state.exercises[exerciseId] = {
          userCode: starterCode,
          status: "unstarted", // 'unstarted' | 'in_progress' | 'completed'
          completedAt: null,
          attempts: 0,
          hintsRevealed: 0,
        };
      }
      return this.state.exercises[exerciseId];
    },

    saveDraft(exerciseId, code) {
      const ex = this.getExerciseState(exerciseId);
      ex.userCode = code;
      if (ex.status === "unstarted") {
        ex.status = "in_progress";
      }
      this.persist();
    },

    markCompleted(exerciseId) {
      const ex = this.getExerciseState(exerciseId);
      ex.status = "completed";
      if (!ex.completedAt) {
        ex.completedAt = new Date().toISOString();
      }
      ex.attempts = (ex.attempts || 0) + 1;
      this.recalculateTotals();
      this.persist();
    },

    recordAttempt(exerciseId) {
      const ex = this.getExerciseState(exerciseId);
      ex.attempts = (ex.attempts || 0) + 1;
      if (ex.status === "unstarted") {
        ex.status = "in_progress";
      }
      this.persist();
    },

    revealHint(exerciseId) {
      const ex = this.getExerciseState(exerciseId);
      ex.hintsRevealed = (ex.hintsRevealed || 0) + 1;
      this.persist();
      return ex.hintsRevealed;
    },

    resetExercise(exerciseId, starterCode = "") {
      this.state.exercises[exerciseId] = {
        userCode: starterCode,
        status: "unstarted",
        completedAt: null,
        attempts: 0,
        hintsRevealed: 0,
      };
      this.recalculateTotals();
      this.persist();
    },

    resetAll() {
      this.state.exercises = {};
      this.state.completedCount = 0;
      this.state.lastActiveExerciseId = null;
      this.persist();
    },

    recalculateTotals() {
      let count = 0;
      for (const key of Object.keys(this.state.exercises || {})) {
        if (this.state.exercises[key].status === "completed") {
          count++;
        }
      }
      this.state.completedCount = count;
    },

    exportJson() {
      return JSON.stringify(this.state, null, 2);
    },

    importJson(jsonString) {
      try {
        const parsed = JSON.parse(jsonString);
        if (!parsed || typeof parsed !== "object") {
          throw new Error("Invalid format: expected JSON object");
        }
        if (parsed.exercises && typeof parsed.exercises === "object") {
          this.state = parsed;
          this.recalculateTotals();
          this.persist();
          return true;
        }
        throw new Error("Invalid schema: missing 'exercises' root key");
      } catch (err) {
        alert("Failed to import progress: " + err.message);
        return false;
      }
    },
  };

  /**
   * ==========================================================================
   * 2. Playground Runtime State
   * ==========================================================================
   */
  const state = {
    bundle: null,
    currentExerciseId: null,
    worker: null,
    workerReady: false,
    editor: null,
    diffEditor: null,
    isDiffMode: false,
    expandedChapters: new Set(),
    activeFilter: "all",
    searchQuery: "",
    revealedHints: 0,
    saveDebounceTimer: null,
    elements: {},
  };

  /**
   * Resolve URL relative to playground.js
   */
  function resolveAssetUrl(filename) {
    if (document.currentScript && document.currentScript.src) {
      return new URL(filename, document.currentScript.src).href;
    }
    const scripts = document.querySelectorAll("script");
    for (const s of scripts) {
      if (s.src && s.src.includes("playground.js")) {
        return new URL(filename, s.src).href;
      }
    }
    return filename;
  }

  /**
   * ==========================================================================
   * 3. DOM Construction & Skeleton Rendering
   * ==========================================================================
   */
  function renderPlaygroundSkeleton(container) {
    if (container.querySelector(".playground-split-layout")) {
      return;
    }

    container.innerHTML = `
      <div class="playground-split-layout">
        <!-- Sidebar: Curriculum & Progress Explorer -->
        <aside class="playground-sidebar" aria-label="Terralings Curriculum Sidebar">
          <div class="sidebar-header">
            <div class="sidebar-title-row">
              <span class="sidebar-title">🏗️ Curriculum</span>
              <div class="sidebar-actions">
                <button id="pg-btn-export" class="sidebar-icon-btn" title="Export Progress (JSON)">📥</button>
                <button id="pg-btn-import" class="sidebar-icon-btn" title="Import Progress (JSON)">📤</button>
                <button id="pg-btn-reset-all" class="sidebar-icon-btn sidebar-icon-danger" title="Reset All Progress">🗑️</button>
                <input type="file" id="pg-file-import" accept=".json" style="display:none;" />
              </div>
            </div>

            <!-- Global Progress Bar -->
            <div class="sidebar-progress-container">
              <div class="sidebar-progress-labels">
                <span id="pg-progress-text" class="sidebar-progress-text">0 / 56 Completed</span>
                <span id="pg-progress-pct" class="sidebar-progress-pct">0%</span>
              </div>
              <div class="sidebar-progress-track">
                <div id="pg-progress-fill" class="sidebar-progress-fill" style="width: 0%;"></div>
              </div>
            </div>

            <!-- Search & Filters -->
            <div class="sidebar-search-row">
              <input type="text" id="pg-search-input" class="sidebar-search-input" placeholder="Search exercises, concepts..." />
            </div>
            <div class="sidebar-filter-tabs">
              <button class="filter-tab active" data-filter="all">All</button>
              <button class="filter-tab" data-filter="incomplete">To Do</button>
              <button class="filter-tab" data-filter="completed">Done</button>
            </div>
          </div>

          <!-- Syllabus Chapters Tree -->
          <div id="pg-syllabus-tree" class="sidebar-syllabus-tree">
            <div class="sidebar-loading-placeholder">⚡ Loading curriculum syllabus...</div>
          </div>
        </aside>

        <!-- Main Workspace: Editor, Controls & Diagnostics -->
        <main class="playground-main-workspace">
          <!-- Exercise Breadcrumb & Top Bar -->
          <div class="workspace-top-bar">
            <div class="workspace-meta-left">
              <span id="pg-chapter-badge" class="chapter-badge">Chapter 01</span>
              <h2 id="pg-exercise-title" class="exercise-title">Loading exercise...</h2>
            </div>
            <div class="workspace-meta-right">
              <div class="nav-stepper">
                <button id="pg-prev-btn" class="nav-btn" title="Previous Exercise (Alt+Left)">← Prev</button>
                <button id="pg-next-btn" class="nav-btn" title="Next Exercise (Alt+Right)">Next →</button>
              </div>
              <button id="pg-fullscreen-btn" class="nav-btn" title="Toggle Fullscreen (F11)">⛶ Fullscreen</button>
              <div id="playground-status" class="playground-status-pill status-loading" title="WebAssembly Engine Status">
                <span class="status-dot"></span>
                <span class="status-text">⚡ Starting Python Wasm...</span>
              </div>
            </div>
          </div>

          <!-- Action Toolbar -->
          <div class="playground-toolbar">
            <button id="playground-run-btn" class="playground-btn playground-btn-primary" title="Execute validation in Pyodide Wasm (Ctrl+Enter)">
              <span class="btn-icon">▶</span>
              <span>Run Solution</span>
              <span class="playground-btn-kbd">Ctrl+Enter</span>
            </button>
            <button id="playground-hint-btn" class="playground-btn" title="Reveal hints step-by-step (H)">
              <span class="btn-icon">💡</span>
              <span class="hint-label">Reveal Hint</span>
            </button>
            <button id="playground-reset-btn" class="playground-btn" title="Reset current file to starter template">
              <span class="btn-icon">↺</span>
              <span>Reset Code</span>
            </button>
            <button id="playground-diff-btn" class="playground-btn" title="Compare side-by-side with reference solution">
              <span class="btn-icon">🔍</span>
              <span class="diff-label">Compare Solution</span>
            </button>
          </div>

          <!-- Progressive Hints Card -->
          <div id="playground-hints" class="playground-hints-card" aria-live="polite"></div>

          <!-- Editor & Terminal Panes -->
          <div class="playground-workspace">
            <div class="playground-editor-pane">
              <div id="playground-editor"></div>
              <div id="playground-diff-editor"></div>
            </div>
            <div class="playground-output-pane">
              <div class="playground-output-header">
                <div class="playground-output-title">
                  <span class="playground-output-title-dot"></span>
                  <span>Terminal Diagnostics</span>
                </div>
                <div id="playground-output-meta" class="playground-output-meta">Pyodide Wasm Sandbox</div>
              </div>
              <pre id="playground-output">⚡ Initializing Python 3.12 WebAssembly Runtime and Monaco Editor...</pre>
            </div>
          </div>
        </main>
      </div>
    `;
  }

  /**
   * Bind cached DOM references.
   */
  function bindElements(container) {
    state.elements = {
      sidebarTree: container.querySelector("#pg-syllabus-tree"),
      progressText: container.querySelector("#pg-progress-text"),
      progressPct: container.querySelector("#pg-progress-pct"),
      progressFill: container.querySelector("#pg-progress-fill"),
      searchInput: container.querySelector("#pg-search-input"),
      filterTabs: container.querySelectorAll(".filter-tab"),
      exportBtn: container.querySelector("#pg-btn-export"),
      importBtn: container.querySelector("#pg-btn-import"),
      importFile: container.querySelector("#pg-file-import"),
      resetAllBtn: container.querySelector("#pg-btn-reset-all"),
      chapterBadge: container.querySelector("#pg-chapter-badge"),
      exerciseTitle: container.querySelector("#pg-exercise-title"),
      prevBtn: container.querySelector("#pg-prev-btn"),
      nextBtn: container.querySelector("#pg-next-btn"),
      fullscreenBtn: container.querySelector("#pg-fullscreen-btn"),
      status: container.querySelector("#playground-status"),
      statusText: container.querySelector("#playground-status .status-text"),
      runBtn: container.querySelector("#playground-run-btn"),
      resetBtn: container.querySelector("#playground-reset-btn"),
      hintBtn: container.querySelector("#playground-hint-btn"),
      hintLabel: container.querySelector("#playground-hint-btn .hint-label"),
      diffBtn: container.querySelector("#playground-diff-btn"),
      diffLabel: container.querySelector("#playground-diff-btn .diff-label"),
      hintsCard: container.querySelector("#playground-hints"),
      workspace: container.querySelector(".playground-workspace"),
      editorContainer: container.querySelector("#playground-editor"),
      diffContainer: container.querySelector("#playground-diff-editor"),
      output: container.querySelector("#playground-output"),
      outputMeta: container.querySelector("#playground-output-meta"),
    };
  }

  function updateStatus(stage, message) {
    const el = state.elements.status;
    const txt = state.elements.statusText;
    if (!el || !txt) return;

    el.className = "playground-status-pill";
    if (stage === "ready") {
      el.classList.add("status-ready");
      state.workerReady = true;
    } else if (stage === "running") {
      el.classList.add("status-running");
    } else if (stage === "error") {
      el.classList.add("status-error");
    } else {
      el.classList.add("status-loading");
    }
    txt.textContent = message;
  }

  /**
   * ==========================================================================
   * 4. Web Worker (Pyodide Wasm) Controller
   * ==========================================================================
   */
  function initWorker() {
    const workerUrl = resolveAssetUrl("playground-worker.js");
    try {
      state.worker = new Worker(workerUrl);
    } catch (e) {
      // Cross-origin fallback via Blob URL
      fetch(workerUrl)
        .then((r) => r.text())
        .then((code) => {
          const blob = new Blob([code], { type: "application/javascript" });
          state.worker = new Worker(URL.createObjectURL(blob));
          attachWorkerListeners();
          if (state.bundle) {
            state.worker.postMessage({ type: "INIT", bundle: state.bundle });
          }
        })
        .catch((err) => {
          updateStatus("error", "Failed to start Web Worker: " + err.message);
        });
      return;
    }
    attachWorkerListeners();
  }

  function attachWorkerListeners() {
    if (!state.worker) return;

    state.worker.onmessage = function (e) {
      const msg = e.data;
      if (!msg || !msg.type) return;

      if (msg.type === "STATUS") {
        updateStatus(msg.stage, msg.message);
        if (msg.stage === "ready") {
          const runBtn = state.elements.runBtn;
          if (runBtn) runBtn.disabled = false;
        }
      } else if (msg.type === "RUN_RESULT") {
        handleRunResult(msg);
      }
    };

    state.worker.onerror = function (err) {
      updateStatus("error", "WebWorker Error: " + (err.message || String(err)));
    };
  }

  /**
   * ==========================================================================
   * 5. Syllabus Tree & Search Filtering
   * ==========================================================================
   */
  function renderSyllabusTree() {
    const treeEl = state.elements.sidebarTree;
    if (!treeEl || !state.bundle) return;

    const chapters = state.bundle.chapters || [];
    const filter = state.activeFilter;
    const query = (state.searchQuery || "").trim().toLowerCase();

    let html = "";
    let matchingExerciseCount = 0;

    for (const chapter of chapters) {
      const chapterExercises = (chapter.exercise_ids || [])
        .map((id) => state.bundle.exercises[id])
        .filter(Boolean);

      // Calculate chapter progress
      let chCompleted = 0;
      for (const ex of chapterExercises) {
        const exState = TerralingsStorage.getExerciseState(ex.id, ex.starter_code);
        if (exState.status === "completed") chCompleted++;
      }

      // Filter exercises by search query and completion filter
      const visibleExercises = chapterExercises.filter((ex) => {
        const exState = TerralingsStorage.getExerciseState(ex.id, ex.starter_code);
        if (filter === "completed" && exState.status !== "completed") return false;
        if (filter === "incomplete" && exState.status === "completed") return false;

        if (query) {
          const matchTitle = ex.title.toLowerCase().includes(query);
          const matchId = ex.id.toLowerCase().includes(query);
          const matchChapter = chapter.title.toLowerCase().includes(query);
          return matchTitle || matchId || matchChapter;
        }
        return true;
      });

      if (query && visibleExercises.length === 0) {
        continue;
      }

      matchingExerciseCount += visibleExercises.length;
      const isExpanded = query ? true : state.expandedChapters.has(chapter.number);
      const isChapterComplete = chCompleted === chapterExercises.length && chapterExercises.length > 0;

      html += `
        <div class="chapter-group ${isExpanded ? "expanded" : ""}" data-chapter-num="${chapter.number}">
          <div class="chapter-header" data-toggle-chapter="${chapter.number}">
            <div class="chapter-header-title">
              <span class="chapter-chevron">▸</span>
              <span class="chapter-num">${String(chapter.number).padStart(2, "0")}.</span>
              <span class="chapter-name" title="${escapeHtml(chapter.title)}">${escapeHtml(chapter.title)}</span>
            </div>
            <span class="chapter-badge-count ${isChapterComplete ? "complete" : ""}">
              ${chCompleted}/${chapterExercises.length} ${isChapterComplete ? "✓" : ""}
            </span>
          </div>
          <div class="chapter-exercise-list">
      `;

      for (const ex of visibleExercises) {
        const exState = TerralingsStorage.getExerciseState(ex.id, ex.starter_code);
        const isActive = ex.id === state.currentExerciseId;
        const status = exState.status;

        let statusIcon = "○";
        let statusClass = "status-unstarted";
        if (status === "completed") {
          statusIcon = "✓";
          statusClass = "status-done";
        } else if (status === "in_progress") {
          statusIcon = "⏳";
          statusClass = "status-progress";
        }

        html += `
          <div class="exercise-item ${isActive ? "active" : ""} ${statusClass}" data-exercise-id="${ex.id}">
            <span class="exercise-status-icon">${statusIcon}</span>
            <div class="exercise-item-content">
              <div class="exercise-item-title">
                <span class="exercise-item-id">${escapeHtml(ex.id)}:</span> ${escapeHtml(ex.title)}
              </div>
            </div>
          </div>
        `;
      }

      html += `
          </div>
        </div>
      `;
    }

    if (matchingExerciseCount === 0) {
      html = `<div class="sidebar-empty">No exercises found matching "${escapeHtml(query)}"</div>`;
    }

    treeEl.innerHTML = html;

    // Attach chapter toggle handlers
    treeEl.querySelectorAll("[data-toggle-chapter]").forEach((el) => {
      el.addEventListener("click", () => {
        const num = parseInt(el.getAttribute("data-toggle-chapter"), 10);
        if (state.expandedChapters.has(num)) {
          state.expandedChapters.delete(num);
        } else {
          state.expandedChapters.add(num);
        }
        renderSyllabusTree();
      });
    });

    // Attach exercise select handlers
    treeEl.querySelectorAll(".exercise-item").forEach((el) => {
      el.addEventListener("click", () => {
        const exId = el.getAttribute("data-exercise-id");
        if (exId) {
          selectExercise(exId);
        }
      });
    });
  }

  /**
   * Select and load an exercise into workspace.
   */
  function selectExercise(exerciseId) {
    if (!state.bundle || !state.bundle.exercises[exerciseId]) return;

    state.currentExerciseId = exerciseId;
    const ex = state.bundle.exercises[exerciseId];
    TerralingsStorage.state.lastActiveExerciseId = exerciseId;
    TerralingsStorage.persist();

    // Update browser URL query parameter for shareable deep-linking
    if (typeof window !== "undefined" && window.history && window.history.replaceState) {
      try {
        const url = new URL(window.location.href);
        url.searchParams.set("exercise", exerciseId);
        window.history.replaceState({}, "", url.toString());
      } catch (e) {
        // Ignore if restricted
      }
    }

    // Auto-expand current chapter
    if (ex.chapter_number) {
      state.expandedChapters.add(ex.chapter_number);
    }

    // Update Header
    if (state.elements.chapterBadge) {
      state.elements.chapterBadge.textContent = `Chapter ${String(ex.chapter_number || 1).padStart(2, "0")}: ${ex.chapter_title || ex.chapter}`;
    }
    if (state.elements.exerciseTitle) {
      state.elements.exerciseTitle.textContent = `${ex.id} — ${ex.title}`;
    }

    // Load saved user code from localStorage
    const savedState = TerralingsStorage.getExerciseState(exerciseId, ex.starter_code);
    state.revealedHints = savedState.hintsRevealed || 0;
    renderHints();

    // Update Monaco editor code
    if (state.editor) {
      state.editor.setValue(savedState.userCode || ex.starter_code || "");
    }

    // If diff editor is active, update diff models
    if (state.isDiffMode && state.diffEditor && window.monaco) {
      updateDiffModels();
    }

    // Update Next/Prev buttons
    updateStepperButtons();

    // Refresh syllabus tree highlight
    renderSyllabusTree();

    // Welcome terminal message
    if (state.elements.output) {
      state.elements.output.innerHTML = `
<span class="term-banner-info">📚 Loaded exercise: <strong>${escapeHtml(ex.id)}</strong> — ${escapeHtml(ex.title)}</span>
<span class="term-dim">Fix the HCL code in the editor, then click </span><span class="term-pass">▶ Run Solution</span><span class="term-dim"> (Ctrl+Enter).</span>
`;
    }
  }

  function getOrderedExerciseList() {
    if (!state.bundle || !state.bundle.chapters) return [];
    const list = [];
    for (const ch of state.bundle.chapters) {
      if (ch.exercise_ids) {
        list.push(...ch.exercise_ids);
      }
    }
    return list;
  }

  function updateStepperButtons() {
    const list = getOrderedExerciseList();
    const idx = list.indexOf(state.currentExerciseId);

    if (state.elements.prevBtn) {
      state.elements.prevBtn.disabled = idx <= 0;
    }
    if (state.elements.nextBtn) {
      state.elements.nextBtn.disabled = idx < 0 || idx >= list.length - 1;
    }
  }

  function goToPreviousExercise() {
    const list = getOrderedExerciseList();
    const idx = list.indexOf(state.currentExerciseId);
    if (idx > 0) {
      selectExercise(list[idx - 1]);
    }
  }

  function goToNextExercise() {
    const list = getOrderedExerciseList();
    const idx = list.indexOf(state.currentExerciseId);
    if (idx >= 0 && idx < list.length - 1) {
      selectExercise(list[idx + 1]);
    }
  }

  /**
   * Render hint drawer cards.
   */
  function renderHints() {
    const card = state.elements.hintsCard;
    const label = state.elements.hintLabel;
    if (!card) return;

    const ex = state.bundle && state.bundle.exercises[state.currentExerciseId];
    const hints = (ex && ex.hints) || [];

    if (state.revealedHints <= 0 || hints.length === 0) {
      card.className = "playground-hints-card";
      card.innerHTML = "";
      if (label) label.textContent = "Reveal Hint";
      return;
    }

    card.className = "playground-hints-card hints-visible";
    let html = "";
    for (let i = 0; i < state.revealedHints && i < hints.length; i++) {
      html += `
        <div class="playground-hint-item">
          <span class="playground-hint-badge">Hint ${i + 1}/${hints.length}</span>
          <span class="playground-hint-text">${escapeHtml(hints[i])}</span>
        </div>
      `;
    }

    card.innerHTML = html;
    if (label) {
      if (state.revealedHints >= hints.length) {
        label.textContent = `All Hints Revealed (${hints.length}/${hints.length})`;
      } else {
        label.textContent = `Next Hint (${state.revealedHints}/${hints.length})`;
      }
    }
  }

  function triggerRevealHint() {
    const ex = state.bundle && state.bundle.exercises[state.currentExerciseId];
    if (!ex || !ex.hints || ex.hints.length === 0) {
      alert("No hints available for this exercise.");
      return;
    }

    if (state.revealedHints < ex.hints.length) {
      state.revealedHints = TerralingsStorage.revealHint(state.currentExerciseId);
      renderHints();
    } else {
      // Toggle off if already all revealed
      state.revealedHints = 0;
      renderHints();
    }
  }

  /**
   * ==========================================================================
   * 6. Diff Mode Controller
   * ==========================================================================
   */
  function toggleDiffMode() {
    state.isDiffMode = !state.isDiffMode;
    const btn = state.elements.diffBtn;
    const label = state.elements.diffLabel;
    const workspace = state.elements.workspace;

    if (state.isDiffMode) {
      if (btn) btn.classList.add("btn-active");
      if (label) label.textContent = "Hide Solution Diff";
      if (workspace) workspace.classList.add("diff-active");
      initDiffEditor();
    } else {
      if (btn) btn.classList.remove("btn-active");
      if (label) label.textContent = "Compare Solution";
      if (workspace) workspace.classList.remove("diff-active");
      if (state.editor) state.editor.layout();
    }
  }

  function initDiffEditor() {
    if (!window.monaco || !state.elements.diffContainer) return;

    if (!state.diffEditor) {
      const isDark = document.documentElement.getAttribute("data-theme") !== "light";
      state.diffEditor = window.monaco.editor.createDiffEditor(state.elements.diffContainer, {
        theme: isDark ? "vs-dark" : "vs",
        readOnly: true,
        renderSideBySide: true,
        automaticLayout: true,
        minimap: { enabled: false },
        fontSize: 13,
      });
    }

    updateDiffModels();
  }

  function updateDiffModels() {
    if (!state.diffEditor || !window.monaco) return;
    const ex = state.bundle && state.bundle.exercises[state.currentExerciseId];
    if (!ex) return;

    const userCode = state.editor ? state.editor.getValue() : ex.starter_code;
    const refCode = ex.solution_code || "# No reference solution available";

    const originalModel = window.monaco.editor.createModel(userCode, "hcl");
    const modifiedModel = window.monaco.editor.createModel(refCode, "hcl");
    state.diffEditor.setModel({ original: originalModel, modified: modifiedModel });
  }

  /**
   * ==========================================================================
   * 7. Execution & Result Handling
   * ==========================================================================
   */
  function runCurrentExercise() {
    if (!state.workerReady || !state.worker) {
      alert("WebAssembly Python runtime is still starting, please wait a moment.");
      return;
    }

    const ex = state.bundle && state.bundle.exercises[state.currentExerciseId];
    if (!ex) return;

    const userCode = state.editor ? state.editor.getValue() : "";

    updateStatus("running", "⏳ Validating HCL AST & rules...");
    if (state.elements.runBtn) state.elements.runBtn.disabled = true;

    if (state.elements.output) {
      state.elements.output.innerHTML = `<span class="term-dim">⚡ Evaluating ${escapeHtml(ex.id)} in WebAssembly sandbox...</span>`;
    }

    TerralingsStorage.recordAttempt(state.currentExerciseId);

    state.worker.postMessage({
      type: "RUN_EXERCISE",
      exerciseId: ex.id,
      code: userCode,
      rules: ex.rules || {},
    });
  }

  function handleRunResult(res) {
    updateStatus("ready", "✅ Ready! WebAssembly HCL engine loaded.");
    if (state.elements.runBtn) state.elements.runBtn.disabled = false;

    const out = state.elements.output;
    const ex = state.bundle && state.bundle.exercises[res.exerciseId];
    if (!out || !ex) return;

    if (res.passed) {
      TerralingsStorage.markCompleted(res.exerciseId);
      updateGlobalProgressBar();
      renderSyllabusTree();

      const list = getOrderedExerciseList();
      const idx = list.indexOf(res.exerciseId);
      const hasNext = idx >= 0 && idx < list.length - 1;

      out.innerHTML = `
<span class="term-banner-pass">✓ EXERCISE PASSED: ${escapeHtml(ex.id)}</span>
<span class="term-pass">${escapeHtml(res.output || "All rules and checks passed successfully!")}</span>
<span class="term-dim">Execution time: ${res.durationMs ? res.durationMs.toFixed(2) : "0.5"}ms (Wasm in-memory evaluation)</span>

${hasNext ? `<button id="pg-next-ex-inline-btn" class="term-inline-btn">Next Exercise →</button>` : `<span class="term-pass">🎉 You have completed all exercises in this track!</span>`}
`;

      const nextBtn = document.getElementById("pg-next-ex-inline-btn");
      if (nextBtn) {
        nextBtn.addEventListener("click", goToNextExercise);
      }
    } else {
      out.innerHTML = `
<span class="term-banner-fail">✕ EXERCISE FAILED: ${escapeHtml(ex.id)}</span>
<span class="term-fail">${escapeHtml(res.output || res.error || "Validation failed.")}</span>
<span class="term-dim">Review the syntax errors or missing attributes above, edit the code, and try again.</span>
`;
    }
  }

  function updateGlobalProgressBar() {
    TerralingsStorage.recalculateTotals();
    const total = (state.bundle && state.bundle.total_exercises) || 56;
    const completed = TerralingsStorage.state.completedCount || 0;
    const pct = total > 0 ? Math.round((completed / total) * 100) : 0;

    if (state.elements.progressText) {
      state.elements.progressText.textContent = `${completed} / ${total} Completed`;
    }
    if (state.elements.progressPct) {
      state.elements.progressPct.textContent = `${pct}%`;
    }
    if (state.elements.progressFill) {
      state.elements.progressFill.style.width = `${pct}%`;
    }
  }

  /**
   * ==========================================================================
   * 8. Monaco Editor Initialization & Tokenizer
   * ==========================================================================
   */
  function registerHclLanguage() {
    if (!window.monaco) return;
    if (window.monaco.languages.getLanguages().some((l) => l.id === "hcl")) return;

    window.monaco.languages.register({ id: "hcl" });
    window.monaco.languages.setMonarchTokensProvider("hcl", {
      defaultToken: "",
      tokenPostfix: ".hcl",
      keywords: [
        "resource", "data", "variable", "output", "locals", "module", "terraform",
        "provider", "check", "import", "moved", "true", "false", "null",
        "for", "in", "dynamic", "content", "each", "count", "self", "var", "local"
      ],
      typeKeywords: ["string", "number", "bool", "list", "map", "set", "object", "tuple", "any"],
      operators: ["=", ">", "<", "!", "~", "?", ":", "==", "<=", ">=", "!=", "&&", "||", "+", "-", "*", "/", "%"],
      symbols: /[=><!~?:&|+\-*\/\^%]+/,
      tokenizer: {
        root: [
          [/[a-zA-Z_]\w*/, {
            cases: {
              "@keywords": "keyword",
              "@typeKeywords": "type",
              "@default": "identifier",
            },
          }],
          { include: "@whitespace" },
          [/[{}()\[\]]/, "@brackets"],
          [/@symbols/, { cases: { "@operators": "operator", "@default": "" } }],
          [/\d*\.\d+([eE][\-+]?\d+)?/, "number.float"],
          [/\d+/, "number"],
          [/[;,.]/, "delimiter"],
          [/"([^"\\]|\\.)*$/, "string.invalid"],
          [/"/, { token: "string.quote", bracket: "@open", next: "@string" }],
        ],
        string: [
          [/[^\\"]+/, "string"],
          [/"/, { token: "string.quote", bracket: "@close", next: "@pop" }],
        ],
        whitespace: [
          [/[ \t\r\n]+/, "white"],
          [/\/\*/, "comment", "@comment"],
          [/\/\/.*$/, "comment"],
          [/#.*$/, "comment"],
        ],
        comment: [
          [/[^\/*]+/, "comment"],
          [/\/\*/, "comment", "@push"],
          ["\\*/", "comment", "@pop"],
          [/[\/*]/, "comment"],
        ],
      },
    });
  }

  function initMonacoEditor() {
    if (!window.monaco || !state.elements.editorContainer) return;

    registerHclLanguage();

    const isDark = document.documentElement.getAttribute("data-theme") !== "light";
    state.editor = window.monaco.editor.create(state.elements.editorContainer, {
      value: "# Loading...",
      language: "hcl",
      theme: isDark ? "vs-dark" : "vs",
      automaticLayout: true,
      fontSize: 13,
      minimap: { enabled: false },
      scrollBeyondLastLine: false,
      tabSize: 2,
      lineNumbers: "on",
    });

    state.editor.onDidChangeModelContent(() => {
      if (!state.currentExerciseId) return;
      const code = state.editor.getValue();
      clearTimeout(state.saveDebounceTimer);
      state.saveDebounceTimer = setTimeout(() => {
        TerralingsStorage.saveDraft(state.currentExerciseId, code);
        renderSyllabusTree();
      }, 300);
    });
  }

  function loadMonacoScript(callback) {
    if (window.monaco) {
      callback();
      return;
    }

    if (window.require && typeof window.require === "function" && window.require.config) {
      window.require.config({
        paths: { vs: "https://cdnjs.cloudflare.com/ajax/libs/monaco-editor/0.45.0/min/vs" },
      });
      window.require(["vs/editor/editor.main"], () => {
        callback();
      });
      return;
    }

    const script = document.createElement("script");
    script.src = "https://cdnjs.cloudflare.com/ajax/libs/monaco-editor/0.45.0/min/vs/loader.min.js";
    script.onload = () => {
      window.require.config({
        paths: { vs: "https://cdnjs.cloudflare.com/ajax/libs/monaco-editor/0.45.0/min/vs" },
      });
      window.require(["vs/editor/editor.main"], () => {
        callback();
      });
    };
    document.head.appendChild(script);
  }

  /**
   * ==========================================================================
   * 9. Event Listeners & Keyboard Shortcuts
   * ==========================================================================
   */
  function setupEventListeners(container) {
    // Run button
    state.elements.runBtn.addEventListener("click", runCurrentExercise);

    // Prev / Next buttons
    state.elements.prevBtn.addEventListener("click", goToPreviousExercise);
    state.elements.nextBtn.addEventListener("click", goToNextExercise);

    // Hint button
    state.elements.hintBtn.addEventListener("click", triggerRevealHint);

    // Reset code button
    state.elements.resetBtn.addEventListener("click", () => {
      const ex = state.bundle && state.bundle.exercises[state.currentExerciseId];
      if (!ex) return;
      if (confirm(`Reset code for ${ex.id} to original starter template?`)) {
        TerralingsStorage.resetExercise(ex.id, ex.starter_code);
        if (state.editor) state.editor.setValue(ex.starter_code || "");
        state.revealedHints = 0;
        renderHints();
        renderSyllabusTree();
      }
    });

    // Diff button
    state.elements.diffBtn.addEventListener("click", toggleDiffMode);

    // Fullscreen toggle
    state.elements.fullscreenBtn.addEventListener("click", () => {
      if (!document.fullscreenElement) {
        if (container.requestFullscreen) {
          container.requestFullscreen();
        } else {
          container.classList.toggle("is-fullscreen");
        }
      } else {
        if (document.exitFullscreen) {
          document.exitFullscreen();
        }
      }
    });

    // Search input
    state.elements.searchInput.addEventListener("input", (e) => {
      state.searchQuery = e.target.value;
      renderSyllabusTree();
    });

    // Filter tabs
    state.elements.filterTabs.forEach((tab) => {
      tab.addEventListener("click", () => {
        state.elements.filterTabs.forEach((t) => t.classList.remove("active"));
        tab.classList.add("active");
        state.activeFilter = tab.getAttribute("data-filter") || "all";
        renderSyllabusTree();
      });
    });

    // Export progress JSON
    state.elements.exportBtn.addEventListener("click", () => {
      const json = TerralingsStorage.exportJson();
      const blob = new Blob([json], { type: "application/json" });
      const url = URL.createObjectURL(blob);
      const a = document.createElement("a");
      a.href = url;
      a.download = `terralings-progress-${new Date().toISOString().slice(0, 10)}.json`;
      a.click();
      URL.revokeObjectURL(url);
    });

    // Import progress JSON
    state.elements.importBtn.addEventListener("click", () => {
      if (state.elements.importFile) {
        state.elements.importFile.click();
      }
    });

    state.elements.importFile.addEventListener("change", (e) => {
      const file = e.target.files && e.target.files[0];
      if (!file) return;
      const reader = new FileReader();
      reader.onload = (evt) => {
        if (TerralingsStorage.importJson(evt.target.result)) {
          alert("Progress imported successfully!");
          updateGlobalProgressBar();
          selectExercise(state.currentExerciseId);
        }
      };
      reader.readAsText(file);
      e.target.value = "";
    });

    // Reset All button
    state.elements.resetAllBtn.addEventListener("click", () => {
      if (confirm("Are you sure you want to RESET ALL Terralings progress? This action cannot be undone.")) {
        TerralingsStorage.resetAll();
        updateGlobalProgressBar();
        selectExercise(state.currentExerciseId);
      }
    });

    // Global keyboard shortcuts
    window.addEventListener("keydown", (e) => {
      if ((e.metaKey || e.ctrlKey) && e.key === "Enter") {
        e.preventDefault();
        runCurrentExercise();
      } else if (e.altKey && e.key === "ArrowLeft") {
        e.preventDefault();
        goToPreviousExercise();
      } else if (e.altKey && e.key === "ArrowRight") {
        e.preventDefault();
        goToNextExercise();
      } else if (e.key === "F11") {
        e.preventDefault();
        state.elements.fullscreenBtn.click();
      }
    });

    // Theme observer
    const observer = new MutationObserver(() => {
      if (window.monaco && state.editor) {
        const isDark = document.documentElement.getAttribute("data-theme") !== "light";
        window.monaco.editor.setTheme(isDark ? "vs-dark" : "vs");
      }
    });
    observer.observe(document.documentElement, {
      attributes: true,
      attributeFilter: ["data-theme", "data-md-color-scheme"],
    });
  }

  function escapeHtml(str) {
    if (!str) return "";
    return String(str)
      .replace(/&/g, "&amp;")
      .replace(/</g, "&lt;")
      .replace(/>/g, "&gt;")
      .replace(/"/g, "&quot;")
      .replace(/'/g, "&#039;");
  }

  /**
   * ==========================================================================
   * 10. Bootstrap Application
   * ==========================================================================
   */
  async function initPlayground() {
    const container = document.getElementById("terralings-playground") || document.getElementById("terralings-app");
    if (!container) return;

    TerralingsStorage.init();
    renderPlaygroundSkeleton(container);
    bindElements(container);
    setupEventListeners(container);

    // Initialize Monaco Editor
    loadMonacoScript(() => {
      initMonacoEditor();
      if (state.bundle && state.currentExerciseId) {
        selectExercise(state.currentExerciseId);
      }
    });

    // Start background Pyodide Web Worker
    initWorker();

    // Fetch curriculum bundle
    try {
      const bundleUrl = resolveAssetUrl("playground-bundle.json") + "?t=" + Date.now();
      const resp = await fetch(bundleUrl);
      if (!resp.ok) {
        throw new Error(`HTTP ${resp.status} ${resp.statusText}`);
      }
      state.bundle = await resp.json();

      // Pass bundle validator code to Worker
      if (state.worker) {
        state.worker.postMessage({ type: "INIT", bundle: state.bundle });
      }

      updateGlobalProgressBar();

      // Determine initial active exercise (URL query param -> saved active -> first exercise)
      const urlParams = new URLSearchParams(window.location.search);
      const urlExercise = urlParams.get("exercise");
      const ordered = getOrderedExerciseList();
      let startExId = ordered[0] || "01_syntax_basics_01";

      if (urlExercise && state.bundle.exercises[urlExercise]) {
        startExId = urlExercise;
      } else if (TerralingsStorage.state.lastActiveExerciseId && state.bundle.exercises[TerralingsStorage.state.lastActiveExerciseId]) {
        startExId = TerralingsStorage.state.lastActiveExerciseId;
      }

      selectExercise(startExId);
    } catch (err) {
      updateStatus("error", "Failed to load curriculum bundle: " + err.message);
    }
  }

  // Auto-start when DOM is ready
  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", initPlayground);
  } else {
    initPlayground();
  }

  if (window.document$ && typeof window.document$.subscribe === "function") {
    window.document$.subscribe(() => {
      initPlayground();
    });
  }
})();
