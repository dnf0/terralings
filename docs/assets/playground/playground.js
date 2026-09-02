/**
 * Terralings In-Browser Interactive Learning Platform.
 *
 * Implements:
 * 1. TerralingsStorage: Zero-backend localStorage state persistence, auto-save,
 *    progress calculation, and JSON backup export/import.
 * 2. TerralingsUI: Monaco Editor integration, Pyodide Web Worker communication,
 *    320px collapsible syllabus sidebar, search filter, progressive hints, diff viewer,
 *    fullscreen mode, and keyboard shortcuts.
 */

/* global monaco, require */

const STORAGE_KEY = "terralings_learning_state_v1";

// ============================================================================
// 1. Storage Manager (localStorage)
// ============================================================================

class TerralingsStorage {
  constructor(bundle) {
    this.bundle = bundle;
    this.state = this.loadState();
  }

  loadState() {
    try {
      const raw = localStorage.getItem(STORAGE_KEY);
      if (raw) {
        const parsed = JSON.parse(raw);
        if (parsed && parsed.version === 1 && parsed.exercises) {
          return parsed;
        }
      }
    } catch (e) {
      console.warn("Failed to read Terralings state from localStorage", e);
    }

    return this.initDefaultState();
  }

  initDefaultState() {
    const state = {
      version: 1,
      lastActiveExerciseId: "primitives01",
      exercises: {},
      stats: {
        completedCount: 0,
        totalCount: this.bundle.stats.totalExercises || 56,
        completionPercentage: 0,
      },
    };

    for (const exId in this.bundle.exercises) {
      state.exercises[exId] = {
        status: "not_started",
        userCode: this.bundle.exercises[exId].starterCode,
        hintsRevealed: 0,
      };
    }

    this.saveState(state);
    return state;
  }

  saveState(stateToSave = null) {
    if (stateToSave) {
      this.state = stateToSave;
    }
    this.recalculateStats();
    try {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(this.state));
    } catch (e) {
      console.error("Failed to write Terralings state to localStorage", e);
    }
  }

  recalculateStats() {
    let completed = 0;
    const total = Object.keys(this.bundle.exercises).length;

    for (const exId in this.bundle.exercises) {
      if (
        this.state.exercises[exId] &&
        this.state.exercises[exId].status === "completed"
      ) {
        completed++;
      }
    }

    this.state.stats = {
      completedCount: completed,
      totalCount: total,
      completionPercentage:
        total > 0 ? Math.round((completed / total) * 100) : 0,
    };
  }

  getExerciseCode(exId) {
    if (this.state.exercises[exId] && this.state.exercises[exId].userCode) {
      return this.state.exercises[exId].userCode;
    }
    return this.bundle.exercises[exId]
      ? this.bundle.exercises[exId].starterCode
      : "";
  }

  saveExerciseCode(exId, code) {
    if (!this.state.exercises[exId]) {
      this.state.exercises[exId] = {
        status: "in_progress",
        userCode: code,
        hintsRevealed: 0,
      };
    } else {
      this.state.exercises[exId].userCode = code;
      if (this.state.exercises[exId].status === "not_started") {
        this.state.exercises[exId].status = "in_progress";
      }
    }
    this.state.lastActiveExerciseId = exId;
    this.saveState();
  }

  markCompleted(exId) {
    if (!this.state.exercises[exId]) {
      this.state.exercises[exId] = {
        status: "completed",
        userCode: this.getExerciseCode(exId),
        hintsRevealed: 0,
        passedAt: new Date().toISOString(),
      };
    } else {
      this.state.exercises[exId].status = "completed";
      this.state.exercises[exId].passedAt = new Date().toISOString();
    }
    this.saveState();
  }

  getHintsRevealed(exId) {
    return (
      (this.state.exercises[exId] &&
        this.state.exercises[exId].hintsRevealed) ||
      0
    );
  }

  revealNextHint(exId) {
    if (!this.state.exercises[exId]) {
      this.state.exercises[exId] = {
        status: "in_progress",
        userCode: this.getExerciseCode(exId),
        hintsRevealed: 1,
      };
    } else {
      const current = this.state.exercises[exId].hintsRevealed || 0;
      this.state.exercises[exId].hintsRevealed = current + 1;
    }
    this.saveState();
    return this.state.exercises[exId].hintsRevealed;
  }

  resetExercise(exId) {
    const starter = this.bundle.exercises[exId]
      ? this.bundle.exercises[exId].starterCode
      : "";
    this.state.exercises[exId] = {
      status: "not_started",
      userCode: starter,
      hintsRevealed: 0,
    };
    this.saveState();
    return starter;
  }

  exportProgress() {
    const dataStr =
      "data:text/json;charset=utf-8," +
      encodeURIComponent(JSON.stringify(this.state, null, 2));
    const downloadAnchor = document.createElement("a");
    const dateStr = new Date().toISOString().slice(0, 10);
    downloadAnchor.setAttribute("href", dataStr);
    downloadAnchor.setAttribute(
      "download",
      `terralings-progress-${dateStr}.json`
    );
    document.body.appendChild(downloadAnchor);
    downloadAnchor.click();
    downloadAnchor.remove();
  }

  importProgress(jsonString) {
    try {
      const imported = JSON.parse(jsonString);
      if (imported && imported.version === 1 && imported.exercises) {
        this.state = imported;
        this.saveState();
        return true;
      }
    } catch (e) {
      console.error("Invalid progress JSON", e);
    }
    return false;
  }

  resetAllProgress() {
    localStorage.removeItem(STORAGE_KEY);
    this.state = this.initDefaultState();
  }
}

// ============================================================================
// 2. UI Controller & Monaco Integration
// ============================================================================

class TerralingsUI {
  constructor() {
    this.bundle = null;
    this.storage = null;
    this.currentExerciseId = null;
    this.monacoEditor = null;
    this.diffEditor = null;
    this.worker = null;
    this.saveTimeout = null;
    this.isDiffMode = false;
    this.filterStatus = "all";
    this.searchQuery = "";
  }

  async init() {
    this.renderLoadingPlaceholder("Loading Terralings Curriculum Bundle...");

    try {
      // Find relative or absolute path for bundle
      const bundlePath = this.getAssetPath("playground-bundle.json");
      const resp = await fetch(bundlePath);
      this.bundle = await resp.json();
    } catch (err) {
      console.error("Error loading bundle:", err);
      this.renderErrorPlaceholder(
        "Failed to load curriculum bundle. Please refresh or check connection."
      );
      return;
    }

    this.storage = new TerralingsStorage(this.bundle);
    this.currentExerciseId =
      this.storage.state.lastActiveExerciseId || "primitives01";
    if (!this.bundle.exercises[this.currentExerciseId]) {
      this.currentExerciseId = Object.keys(this.bundle.exercises)[0];
    }

    this.initWorker();
    this.renderLayout();
    this.initMonaco();
    this.bindEvents();
  }

  getAssetPath(fileName) {
    // MkDocs material site location detection
    const currentHref = window.location.pathname;
    if (currentHref.includes("/playground")) {
      return `../assets/playground/${fileName}`;
    }
    return `./assets/playground/${fileName}`;
  }

  initWorker() {
    const workerPath = this.getAssetPath("playground-worker.js");
    this.worker = new Worker(workerPath);

    this.worker.onmessage = (e) => {
      const msg = e.data;
      if (!msg) return;

      if (msg.type === "STATUS") {
        this.updateStatusBadge(msg.message, msg.stage);
      } else if (msg.type === "RUN_RESULT") {
        this.handleRunResult(msg);
      }
    };

    this.worker.postMessage({
      type: "INIT",
      bundle: this.bundle,
    });
  }

  renderLayout() {
    const container = document.getElementById("terralings-app");
    if (!container) return;

    container.innerHTML = `
      <div class="terralings-wrapper" id="terralings-workspace">
        <!-- Sidebar Navigation (320px) -->
        <aside class="terralings-sidebar">
          <div class="terralings-sidebar-header">
            <div class="terralings-brand">
              <span class="terralings-logo">🏗️</span>
              <h2>Terralings</h2>
            </div>
            <div class="terralings-progress-bar-container">
              <div class="terralings-progress-text">
                <span id="progress-count">0 / ${this.bundle.stats.totalExercises} Completed</span>
                <span id="progress-percent" class="progress-badge">0%</span>
              </div>
              <div class="terralings-progress-track">
                <div id="progress-fill" class="terralings-progress-fill" style="width: 0%"></div>
              </div>
            </div>
            <div class="terralings-storage-actions">
              <button id="btn-export-json" title="Export Progress (JSON)" class="icon-btn">💾 Export</button>
              <button id="btn-import-json" title="Import Progress (JSON)" class="icon-btn">📂 Import</button>
              <button id="btn-reset-all" title="Reset All Progress" class="icon-btn danger-icon">🗑️ Reset</button>
              <input type="file" id="file-import-input" accept=".json" style="display:none;" />
            </div>
          </div>

          <div class="terralings-search-box">
            <input type="text" id="syllabus-search" placeholder="🔍 Search exercises, concepts..." />
            <div class="terralings-filter-tabs">
              <button class="filter-tab active" data-filter="all">All</button>
              <button class="filter-tab" data-filter="todo">To Do</button>
              <button class="filter-tab" data-filter="done">Done</button>
            </div>
          </div>

          <div class="terralings-syllabus-tree" id="syllabus-tree">
            <!-- Dynamic Chapter Accordions -->
          </div>
        </aside>

        <!-- Main Code & Diagnostics Workspace -->
        <main class="terralings-main">
          <!-- Workspace Header -->
          <header class="terralings-header">
            <div class="terralings-exercise-info">
              <span class="chapter-badge" id="ex-chapter-badge">Chapter 01</span>
              <h3 id="ex-title">primitives01</h3>
              <span class="status-indicator" id="ex-status-pill">⏳ In Progress</span>
            </div>
            <div class="terralings-header-actions">
              <button id="btn-prev-ex" class="action-btn secondary-btn" title="Previous Exercise (Alt+Left)">← Prev</button>
              <button id="btn-next-ex" class="action-btn secondary-btn" title="Next Exercise (Alt+Right)">Next →</button>
              <button id="btn-fullscreen" class="action-btn secondary-btn" title="Toggle Fullscreen (F11)">⛶ Fullscreen</button>
            </div>
          </header>

          <!-- Action Toolbar -->
          <div class="terralings-toolbar">
            <button id="btn-run" class="action-btn primary-btn" title="Run Solution (Ctrl+Enter)">
              ▶ Run Solution
            </button>
            <button id="btn-hint" class="action-btn secondary-btn" title="Reveal Hint (H)">
              💡 Hint (<span id="hint-count-badge">0</span>)
            </button>
            <button id="btn-reset-code" class="action-btn secondary-btn" title="Reset to Starter Template">
              ↺ Reset Code
            </button>
            <button id="btn-diff-view" class="action-btn secondary-btn" title="Compare with Reference Solution">
              🔍 Solution Diff
            </button>
          </div>

          <!-- Hint Drawer -->
          <div id="hint-drawer" class="terralings-hint-drawer" style="display: none;">
            <div class="hint-header">
              <h4>💡 Progressive Hints</h4>
              <button id="btn-close-hint" class="close-btn">×</button>
            </div>
            <div id="hint-content" class="hint-body"></div>
          </div>

          <!-- Editor Container -->
          <div class="terralings-editor-container">
            <div id="monaco-code-editor" class="editor-pane"></div>
            <div id="monaco-diff-editor" class="editor-pane" style="display: none;"></div>
          </div>

          <!-- Terminal Diagnostics Pane -->
          <section class="terralings-terminal">
            <div class="terminal-header">
              <div class="terminal-title">
                <span class="terminal-dot"></span>
                <span>Diagnostics & Output</span>
              </div>
              <div class="terminal-status" id="terminal-status-badge">⚡ Pyodide Ready</div>
            </div>
            <pre class="terminal-body" id="terminal-output">Ready. Press ▶ Run Solution or Ctrl+Enter to validate your HCL code.</pre>
          </section>
        </main>
      </div>
    `;

    this.renderSyllabusTree();
    this.updateGlobalProgress();
  }

  renderSyllabusTree() {
    const tree = document.getElementById("syllabus-tree");
    if (!tree) return;

    tree.innerHTML = "";

    this.bundle.chapters.forEach((ch) => {
      // Filter exercises in chapter
      const matchingExercises = ch.exercises.filter((ex) => {
        const matchesSearch =
          !this.searchQuery ||
          ex.name.toLowerCase().includes(this.searchQuery) ||
          ex.title.toLowerCase().includes(this.searchQuery) ||
          ch.title.toLowerCase().includes(this.searchQuery);

        const exState = this.storage.state.exercises[ex.name] || {};
        const status = exState.status || "not_started";

        let matchesStatus = true;
        if (this.filterStatus === "todo") matchesStatus = status !== "completed";
        if (this.filterStatus === "done") matchesStatus = status === "completed";

        return matchesSearch && matchesStatus;
      });

      if (matchingExercises.length === 0 && (this.searchQuery || this.filterStatus !== "all")) {
        return;
      }

      const completedInChapter = ch.exercises.filter((ex) => {
        const s = this.storage.state.exercises[ex.name];
        return s && s.status === "completed";
      }).length;

      const isCurrentChapter = ch.exercises.some((e) => e.name === this.currentExerciseId);

      const chapterDiv = document.createElement("div");
      chapterDiv.className = `syllabus-chapter ${isCurrentChapter ? "open" : ""}`;

      const chapterHeader = document.createElement("div");
      chapterHeader.className = "chapter-header";
      chapterHeader.innerHTML = `
        <div class="chapter-title-row">
          <span class="accordion-arrow">▸</span>
          <span class="chapter-number">${String(ch.number).padStart(2, "0")}.</span>
          <span class="chapter-name">${ch.title}</span>
        </div>
        <span class="chapter-count-badge ${completedInChapter === ch.exercises.length ? "all-done" : ""}">${completedInChapter}/${ch.exercises.length} ✓</span>
      `;

      chapterHeader.onclick = () => {
        chapterDiv.classList.toggle("open");
      };

      const exerciseList = document.createElement("div");
      exerciseList.className = "exercise-list";

      matchingExercises.forEach((ex) => {
        const exState = this.storage.state.exercises[ex.name] || {};
        const status = exState.status || "not_started";
        const isActive = ex.name === this.currentExerciseId;

        const item = document.createElement("div");
        item.className = `exercise-item ${isActive ? "active" : ""} ${status}`;

        let statusIcon = "○";
        if (status === "completed") statusIcon = "✓";
        else if (status === "in_progress") statusIcon = "⏳";

        item.innerHTML = `
          <span class="exercise-status-icon">${statusIcon}</span>
          <span class="exercise-name">${ex.name}</span>
          <span class="exercise-label">${ex.title}</span>
        `;

        item.onclick = (e) => {
          e.stopPropagation();
          this.switchExercise(ex.name);
        };

        exerciseList.appendChild(item);
      });

      chapterDiv.appendChild(chapterHeader);
      chapterDiv.appendChild(exerciseList);
      tree.appendChild(chapterDiv);
    });
  }

  updateGlobalProgress() {
    this.storage.recalculateStats();
    const stats = this.storage.state.stats;

    const countEl = document.getElementById("progress-count");
    const percentEl = document.getElementById("progress-percent");
    const fillEl = document.getElementById("progress-fill");

    if (countEl)
      countEl.textContent = `${stats.completedCount} / ${stats.totalCount} Completed`;
    if (percentEl) percentEl.textContent = `${stats.completionPercentage}%`;
    if (fillEl) fillEl.style.width = `${stats.completionPercentage}%`;
  }

  initMonaco() {
    // Configure Monaco Environment for AMD Loader
    if (typeof require !== "undefined") {
      require.config({
        paths: {
          vs: "https://cdnjs.cloudflare.com/ajax/libs/monaco-editor/0.45.0/min/vs",
        },
      });

      require(["vs/editor/editor.main"], () => {
        // Register HCL / Terraform language tokens
        monaco.languages.register({ id: "terraform" });
        monaco.languages.setMonarchTokensProvider("terraform", {
          keywords: [
            "terraform",
            "resource",
            "data",
            "variable",
            "locals",
            "output",
            "module",
            "moved",
            "check",
            "import",
            "removed",
            "provider",
          ],
          typeKeywords: [
            "string",
            "number",
            "bool",
            "list",
            "map",
            "set",
            "object",
            "tuple",
            "any",
          ],
          operators: ["=", "==", "!=", "<=", ">=", "<", ">", "&&", "||", "!", "?"],
          tokenizer: {
            root: [
              [/[a-zA-Z_]\w*/, { cases: { "@keywords": "keyword", "@typeKeywords": "type", "@default": "identifier" } }],
              [/"([^"\\]|\\.)*"/, "string"],
              [/<<-?\s*(\w+)/, "string.heredoc"],
              [/#.*$/, "comment"],
              [/\/\/.*$/, "comment"],
              [/\/\*/, "comment", "@comment"],
              [/[{}()\[\]]/, "@brackets"],
              [/[=><!~?:]+/, "operator"],
              [/\d+/, "number"],
            ],
            comment: [
              [/[^\/*]+/, "comment"],
              [/\*\//, "comment", "@pop"],
              [/[\/*]/, "comment"],
            ],
          },
        });

        const isDark =
          document.body.getAttribute("data-md-color-scheme") === "slate" ||
          window.matchMedia("(prefers-color-scheme: dark)").matches;

        const codeContainer = document.getElementById("monaco-code-editor");
        if (codeContainer) {
          this.monacoEditor = monaco.editor.create(codeContainer, {
            value: this.storage.getExerciseCode(this.currentExerciseId),
            language: "terraform",
            theme: isDark ? "vs-dark" : "vs",
            automaticLayout: true,
            fontSize: 14,
            minimap: { enabled: false },
            lineNumbers: "on",
            scrollBeyondLastLine: false,
            tabSize: 2,
          });

          this.monacoEditor.onDidChangeModelContent(() => {
            clearTimeout(this.saveTimeout);
            this.saveTimeout = setTimeout(() => {
              const code = this.monacoEditor.getValue();
              this.storage.saveExerciseCode(this.currentExerciseId, code);
              this.renderSyllabusTree();
            }, 300);
          });
        }

        const diffContainer = document.getElementById("monaco-diff-editor");
        if (diffContainer) {
          this.diffEditor = monaco.editor.createDiffEditor(diffContainer, {
            theme: isDark ? "vs-dark" : "vs",
            automaticLayout: true,
            fontSize: 14,
            readOnly: true,
            renderSideBySide: true,
          });
        }

        this.updateExerciseView();
      });
    }
  }

  switchExercise(exId) {
    if (!this.bundle.exercises[exId]) return;

    this.currentExerciseId = exId;
    this.storage.state.lastActiveExerciseId = exId;
    this.storage.saveState();

    if (this.isDiffMode) {
      this.toggleDiffMode(false);
    }

    if (this.monacoEditor) {
      this.monacoEditor.setValue(this.storage.getExerciseCode(exId));
    }

    this.updateExerciseView();
    this.renderSyllabusTree();
    this.hideHints();
  }

  updateExerciseView() {
    const ex = this.bundle.exercises[this.currentExerciseId];
    if (!ex) return;

    const chBadge = document.getElementById("ex-chapter-badge");
    const titleEl = document.getElementById("ex-title");
    const statusPill = document.getElementById("ex-status-pill");
    const hintBadge = document.getElementById("hint-count-badge");

    if (chBadge)
      chBadge.textContent = `Chapter ${String(ex.chapterNumber).padStart(2, "0")}: ${ex.chapterTitle}`;
    if (titleEl) titleEl.textContent = `${ex.name} - ${ex.title}`;

    const exState = this.storage.state.exercises[ex.name] || {};
    const status = exState.status || "not_started";

    if (statusPill) {
      statusPill.className = `status-indicator ${status}`;
      statusPill.textContent =
        status === "completed"
          ? "✓ Completed"
          : status === "in_progress"
          ? "⏳ In Progress"
          : "○ Not Started";
    }

    const revealed = this.storage.getHintsRevealed(ex.name);
    if (hintBadge) {
      hintBadge.textContent = `${revealed}/${ex.hints.length}`;
    }
  }

  runValidation() {
    if (!this.worker) return;

    const code = this.monacoEditor ? this.monacoEditor.getValue() : "";
    this.updateTerminal("Validating HCL in WebAssembly...", "loading");

    this.worker.postMessage({
      type: "RUN_EXERCISE",
      exerciseId: this.currentExerciseId,
      code: code,
    });
  }

  handleRunResult(result) {
    const term = document.getElementById("terminal-output");
    const badge = document.getElementById("terminal-status-badge");

    if (result.passed) {
      this.storage.markCompleted(this.currentExerciseId);
      this.updateGlobalProgress();
      this.renderSyllabusTree();
      this.updateExerciseView();

      if (badge) {
        badge.className = "terminal-status success";
        badge.textContent = `✓ Passed in ${result.durationMs}ms`;
      }

      if (term) {
        term.innerHTML = `<span class="term-pass">${result.output}</span>\n\n🎉 <strong style="color:#4ade80;">Exercise passed!</strong> Press Next → or Alt+Right to proceed to next exercise.`;
      }
    } else {
      if (badge) {
        badge.className = "terminal-status error";
        badge.textContent = `✕ Failed in ${result.durationMs}ms`;
      }

      if (term) {
        let errLineMsg = result.line ? ` [Line ${result.line}]` : "";
        term.innerHTML = `<span class="term-fail">✕ Validation Error${errLineMsg}:</span>\n${result.error}\n\n<span class="term-dim">${result.output}</span>`;
      }
    }
  }

  revealNextHint() {
    const ex = this.bundle.exercises[this.currentExerciseId];
    if (!ex || !ex.hints || ex.hints.length === 0) return;

    const revealedCount = this.storage.revealNextHint(this.currentExerciseId);
    this.showHints(revealedCount);
    this.updateExerciseView();
  }

  showHints(count) {
    const ex = this.bundle.exercises[this.currentExerciseId];
    const drawer = document.getElementById("hint-drawer");
    const content = document.getElementById("hint-content");

    if (!drawer || !content) return;

    drawer.style.display = "block";
    let html = "";

    ex.hints.slice(0, count).forEach((hint, idx) => {
      html += `
        <div class="hint-item">
          <strong>Hint Tier ${idx + 1}:</strong>
          <p>${hint}</p>
        </div>
      `;
    });

    if (count < ex.hints.length) {
      html += `
        <button id="btn-reveal-more-hints" class="action-btn secondary-btn" style="margin-top:8px;">
          💡 Reveal Next Hint (${count + 1}/${ex.hints.length})
        </button>
      `;
    }

    content.innerHTML = html;

    const moreBtn = document.getElementById("btn-reveal-more-hints");
    if (moreBtn) {
      moreBtn.onclick = () => this.revealNextHint();
    }
  }

  hideHints() {
    const drawer = document.getElementById("hint-drawer");
    if (drawer) drawer.style.display = "none";
  }

  toggleDiffMode(forceState = null) {
    this.isDiffMode = forceState !== null ? forceState : !this.isDiffMode;

    const codeContainer = document.getElementById("monaco-code-editor");
    const diffContainer = document.getElementById("monaco-diff-editor");
    const diffBtn = document.getElementById("btn-diff-view");

    if (this.isDiffMode) {
      if (codeContainer) codeContainer.style.display = "none";
      if (diffContainer) diffContainer.style.display = "block";
      if (diffBtn) diffBtn.textContent = "✎ Code Editor";

      const currentCode = this.storage.getExerciseCode(this.currentExerciseId);
      const solutionCode =
        this.bundle.exercises[this.currentExerciseId].solutionCode;

      if (this.diffEditor) {
        this.diffEditor.setModel({
          original: monaco.editor.createModel(currentCode, "terraform"),
          modified: monaco.editor.createModel(solutionCode, "terraform"),
        });
      }
    } else {
      if (codeContainer) codeContainer.style.display = "block";
      if (diffContainer) diffContainer.style.display = "none";
      if (diffBtn) diffBtn.textContent = "🔍 Solution Diff";
    }
  }

  resetCurrentCode() {
    if (confirm("Reset current exercise to starter template? Your edits will be replaced.")) {
      const starter = this.storage.resetExercise(this.currentExerciseId);
      if (this.monacoEditor) {
        this.monacoEditor.setValue(starter);
      }
      this.updateExerciseView();
      this.renderSyllabusTree();
      this.updateTerminal("Reset code to original template.", "ready");
    }
  }

  toggleFullscreen() {
    const ws = document.getElementById("terralings-workspace");
    if (ws) {
      ws.classList.toggle("terralings-fullscreen");
      if (this.monacoEditor) this.monacoEditor.layout();
      if (this.diffEditor) this.diffEditor.layout();
    }
  }

  navigateNext() {
    const allIds = Object.keys(this.bundle.exercises);
    const idx = allIds.indexOf(this.currentExerciseId);
    if (idx >= 0 && idx < allIds.length - 1) {
      this.switchExercise(allIds[idx + 1]);
    }
  }

  navigatePrev() {
    const allIds = Object.keys(this.bundle.exercises);
    const idx = allIds.indexOf(this.currentExerciseId);
    if (idx > 0) {
      this.switchExercise(allIds[idx - 1]);
    }
  }

  updateTerminal(text, statusStage = "ready") {
    const term = document.getElementById("terminal-output");
    const badge = document.getElementById("terminal-status-badge");
    if (term) term.textContent = text;
    if (badge) {
      badge.className = `terminal-status ${statusStage}`;
      badge.textContent = statusStage === "loading" ? "⏳ Running..." : "⚡ Ready";
    }
  }

  updateStatusBadge(msg, stage) {
    const badge = document.getElementById("terminal-status-badge");
    if (badge) {
      badge.textContent = msg;
      badge.className = `terminal-status ${stage === "ready" ? "success" : "loading"}`;
    }
  }

  renderLoadingPlaceholder(message) {
    const container = document.getElementById("terralings-app");
    if (container) {
      container.innerHTML = `
        <div class="terralings-loading-screen">
          <div class="spinner"></div>
          <p>${message}</p>
        </div>
      `;
    }
  }

  renderErrorPlaceholder(message) {
    const container = document.getElementById("terralings-app");
    if (container) {
      container.innerHTML = `
        <div class="terralings-error-screen">
          <p>⚠️ ${message}</p>
        </div>
      `;
    }
  }

  bindEvents() {
    // Toolbar buttons
    document.getElementById("btn-run")?.addEventListener("click", () => this.runValidation());
    document.getElementById("btn-hint")?.addEventListener("click", () => this.revealNextHint());
    document.getElementById("btn-close-hint")?.addEventListener("click", () => this.hideHints());
    document.getElementById("btn-reset-code")?.addEventListener("click", () => this.resetCurrentCode());
    document.getElementById("btn-diff-view")?.addEventListener("click", () => this.toggleDiffMode());
    document.getElementById("btn-fullscreen")?.addEventListener("click", () => this.toggleFullscreen());
    document.getElementById("btn-next-ex")?.addEventListener("click", () => this.navigateNext());
    document.getElementById("btn-prev-ex")?.addEventListener("click", () => this.navigatePrev());

    // Storage Actions
    document.getElementById("btn-export-json")?.addEventListener("click", () => this.storage.exportProgress());

    const fileInput = document.getElementById("file-import-input");
    document.getElementById("btn-import-json")?.addEventListener("click", () => {
      fileInput?.click();
    });

    fileInput?.addEventListener("change", (e) => {
      const file = e.target.files[0];
      if (!file) return;
      const reader = new FileReader();
      reader.onload = (evt) => {
        if (this.storage.importProgress(evt.target.result)) {
          alert("✓ Progress imported successfully!");
          this.updateGlobalProgress();
          this.switchExercise(this.storage.state.lastActiveExerciseId || "primitives01");
        } else {
          alert("✕ Failed to import progress file. Invalid JSON structure.");
        }
      };
      reader.readAsText(file);
    });

    document.getElementById("btn-reset-all")?.addEventListener("click", () => {
      if (confirm("Reset ALL progress across all 56 exercises? This cannot be undone.")) {
        this.storage.resetAllProgress();
        this.updateGlobalProgress();
        this.switchExercise("primitives01");
      }
    });

    // Search and Filters
    const searchInput = document.getElementById("syllabus-search");
    searchInput?.addEventListener("input", (e) => {
      this.searchQuery = e.target.value.toLowerCase().trim();
      this.renderSyllabusTree();
    });

    document.querySelectorAll(".filter-tab").forEach((tab) => {
      tab.addEventListener("click", (e) => {
        document.querySelectorAll(".filter-tab").forEach((t) => t.classList.remove("active"));
        e.target.classList.add("active");
        this.filterStatus = e.target.dataset.filter;
        this.renderSyllabusTree();
      });
    });

    // Global Keybindings
    window.addEventListener("keydown", (e) => {
      if ((e.ctrlKey || e.metaKey) && e.key === "Enter") {
        e.preventDefault();
        this.runValidation();
      } else if (e.altKey && e.key === "ArrowRight") {
        e.preventDefault();
        this.navigateNext();
      } else if (e.altKey && e.key === "ArrowLeft") {
        e.preventDefault();
        this.navigatePrev();
      } else if (e.key === "F11") {
        e.preventDefault();
        this.toggleFullscreen();
      }
    });
  }
}

// Auto-initialize when page loads
document.addEventListener("DOMContentLoaded", () => {
  if (document.getElementById("terralings-app")) {
    const app = new TerralingsUI();
    app.init();
    window.terralingsApp = app;
  }
});
