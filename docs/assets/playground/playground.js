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
   * 1. TerralingsStorage: Client-Side Progress & Working Code Persistence
   * ==========================================================================
   */
  const TerralingsStorage = {
    state: null,
    saveTimeout: null,

    init(bundle) {
      let saved = null;
      try {
        const raw = localStorage.getItem(STORAGE_KEY);
        if (raw) {
          saved = JSON.parse(raw);
        }
      } catch (e) {
        console.warn("Failed to read Terralings state from localStorage:", e);
      }

      const totalExercises =
        bundle && bundle.exercises ? Object.keys(bundle.exercises).length : 56;

      if (!saved || saved.version !== 1 || !saved.exercises) {
        saved = {
          version: 1,
          lastActiveExerciseId: "primitives01",
          exercises: {},
          stats: {
            completedCount: 0,
            totalCount: totalExercises,
            completionPercentage: 0,
          },
        };
      }

      if (bundle && bundle.exercises) {
        for (const [id, ex] of Object.entries(bundle.exercises)) {
          if (!saved.exercises[id]) {
            saved.exercises[id] = {
              status: "not_started",
              userCode: ex.starter_code || "",
              hintsRevealed: 0,
            };
          }
        }
      }

      this.state = saved;
      this.recalculateStats(bundle);
      this.persist();
      return this.state;
    },

    recalculateStats(bundle) {
      if (!this.state || !this.state.exercises) return;
      let completed = 0;
      const total =
        bundle && bundle.exercises
          ? Object.keys(bundle.exercises).length
          : Object.keys(this.state.exercises).length;

      for (const exState of Object.values(this.state.exercises)) {
        if (exState.status === "completed") {
          completed++;
        }
      }

      this.state.stats = {
        completedCount: completed,
        totalCount: total || 1,
        completionPercentage:
          total > 0 ? Math.round((completed / total) * 100) : 0,
      };
    },

    persist() {
      if (!this.state) return;
      try {
        localStorage.setItem(STORAGE_KEY, JSON.stringify(this.state));
      } catch (e) {
        console.warn("Failed to write Terralings state to localStorage:", e);
      }
    },

    getExerciseState(exerciseId, defaultStarterCode = "") {
      if (!this.state)
        return {
          status: "not_started",
          userCode: defaultStarterCode,
          hintsRevealed: 0,
        };
      if (!this.state.exercises[exerciseId]) {
        this.state.exercises[exerciseId] = {
          status: "not_started",
          userCode: defaultStarterCode,
          hintsRevealed: 0,
        };
        this.persist();
      }
      return this.state.exercises[exerciseId];
    },

    saveExerciseCode(exerciseId, code) {
      if (!this.state) return;
      const exState = this.getExerciseState(exerciseId, code);
      exState.userCode = code;
      if (exState.status === "not_started") {
        exState.status = "in_progress";
      }
      exState.lastEvaluatedAt = new Date().toISOString();

      clearTimeout(this.saveTimeout);
      this.saveTimeout = setTimeout(() => {
        this.persist();
      }, 300);
    },

    markCompleted(exerciseId, bundle) {
      if (!this.state) return;
      const exState = this.getExerciseState(exerciseId);
      exState.status = "completed";
      exState.passedAt = new Date().toISOString();
      this.recalculateStats(bundle);
      this.persist();
    },

    setHintsRevealed(exerciseId, count) {
      if (!this.state) return;
      const exState = this.getExerciseState(exerciseId);
      exState.hintsRevealed = count;
      this.persist();
    },

    resetExercise(exerciseId, starterCode) {
      if (!this.state) return;
      this.state.exercises[exerciseId] = {
        status: "not_started",
        userCode: starterCode || "",
        hintsRevealed: 0,
      };
      this.persist();
    },

    resetAll(bundle) {
      const total =
        bundle && bundle.exercises ? Object.keys(bundle.exercises).length : 56;
      this.state = {
        version: 1,
        lastActiveExerciseId: "primitives01",
        exercises: {},
        stats: {
          completedCount: 0,
          totalCount: total,
          completionPercentage: 0,
        },
      };
      if (bundle && bundle.exercises) {
        for (const [id, ex] of Object.entries(bundle.exercises)) {
          this.state.exercises[id] = {
            status: "not_started",
            userCode: ex.starter_code || "",
            hintsRevealed: 0,
          };
        }
      }
      this.persist();
    },

    exportJson() {
      if (!this.state) return;
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
    },

    importJson(jsonText, bundle) {
      try {
        const parsed = JSON.parse(jsonText);
        if (parsed && parsed.exercises) {
          this.state = parsed;
          this.recalculateStats(bundle);
          this.persist();
          return true;
        }
      } catch (e) {
        console.error("Invalid progress JSON", e);
      }
      return false;
    },
  };

  /**
   * ==========================================================================
   * 2. State & DOM References
   * ==========================================================================
   */
  const state = {
    bundle: null,
    worker: null,
    workerReady: false,
    monacoLoaded: false,
    editor: null,
    diffEditor: null,
    originalModel: null,
    modifiedModel: null,
    currentExerciseId: "primitives01",
    revealedHints: 0,
    isDiffMode: false,
    isRunning: false,
    sidebarFilter: "all",
    searchQuery: "",
    expandedChapters: new Set(["01"]),
    container: null,
    elements: {},
  };

  function resolveAssetUrl(filename) {
    if (document.currentScript && document.currentScript.src) {
      return new URL(filename, document.currentScript.src).href;
    }
    const scripts = document.querySelectorAll('script[src*="playground.js"]');
    if (scripts.length > 0) {
      const src = scripts[scripts.length - 1].src;
      return new URL(filename, src).href;
    }
    return "assets/playground/" + filename;
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

  function getMonacoTheme() {
    const scheme = document.body.getAttribute("data-md-color-scheme");
    if (scheme === "slate") return "vs-dark";
    if (scheme === "default") return "vs";
    if (
      window.matchMedia &&
      window.matchMedia("(prefers-color-scheme: dark)").matches
    ) {
      return "vs-dark";
    }
    return "vs";
  }

  /**
   * ==========================================================================
   * 3. UI Layout Rendering & DOM Binding
   * ==========================================================================
   */
  function renderPlaygroundSkeleton(container) {
    container.innerHTML = `
      <div class="terralings-wrapper" id="terralings-workspace">
        <!-- Sidebar Navigation (320px) -->
        <aside class="terralings-sidebar" aria-label="Curriculum Sidebar">
          <div class="terralings-sidebar-header">
            <div class="terralings-brand">
              <span class="terralings-logo">🏗️</span>
              <h2>Terralings</h2>
            </div>
            <div class="terralings-progress-bar-container">
              <div class="terralings-progress-text">
                <span id="progress-count">0 / 56 Completed</span>
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
            <div id="hint-content"></div>
          </div>

          <!-- Monaco Editor Area -->
          <div class="terralings-editor-container" id="editor-container">
            <div id="monaco-editor" class="editor-pane"></div>
            <div id="monaco-diff" class="editor-pane" style="display: none;"></div>
          </div>

          <!-- Diagnostics Terminal -->
          <div class="terralings-terminal">
            <div class="terminal-header">
              <div class="terminal-title">
                <span class="terminal-dot"></span>
                <span>Diagnostics & Output</span>
              </div>
              <span id="terminal-status-badge" class="terminal-status loading">⚡ Initializing WebAssembly...</span>
            </div>
            <pre class="terminal-body" id="terminal-output">⚡ Starting Pyodide Python 3.12 WebAssembly Engine...</pre>
          </div>
        </main>
      </div>
    `;
  }

  function bindElements(container) {
    state.elements = {
      workspace: container.querySelector("#terralings-workspace"),
      syllabusTree: container.querySelector("#syllabus-tree"),
      progressCount: container.querySelector("#progress-count"),
      progressPercent: container.querySelector("#progress-percent"),
      progressFill: container.querySelector("#progress-fill"),
      chapterBadge: container.querySelector("#ex-chapter-badge"),
      exTitle: container.querySelector("#ex-title"),
      exStatusPill: container.querySelector("#ex-status-pill"),
      hintBadge: container.querySelector("#hint-count-badge"),
      hintDrawer: container.querySelector("#hint-drawer"),
      hintContent: container.querySelector("#hint-content"),
      monacoPane: container.querySelector("#monaco-editor"),
      diffPane: container.querySelector("#monaco-diff"),
      terminalStatus: container.querySelector("#terminal-status-badge"),
      terminalOutput: container.querySelector("#terminal-output"),
      searchInput: container.querySelector("#syllabus-search"),
      filterTabs: container.querySelectorAll(".filter-tab"),
      runBtn: container.querySelector("#btn-run"),
      hintBtn: container.querySelector("#btn-hint"),
      closeHintBtn: container.querySelector("#btn-close-hint"),
      resetCodeBtn: container.querySelector("#btn-reset-code"),
      diffBtn: container.querySelector("#btn-diff-view"),
      prevBtn: container.querySelector("#btn-prev-ex"),
      nextBtn: container.querySelector("#btn-next-ex"),
      fullscreenBtn: container.querySelector("#btn-fullscreen"),
      exportBtn: container.querySelector("#btn-export-json"),
      importBtn: container.querySelector("#btn-import-json"),
      resetAllBtn: container.querySelector("#btn-reset-all"),
      fileInput: container.querySelector("#file-import-input"),
    };
  }

  function updateStatus(stage, text) {
    const badge = state.elements.terminalStatus;
    if (badge) {
      badge.textContent = text;
      badge.className = `terminal-status ${stage}`;
    }
  }

  function updateGlobalProgress() {
    TerralingsStorage.recalculateStats(state.bundle);
    const stats = TerralingsStorage.state.stats;
    if (state.elements.progressCount) {
      state.elements.progressCount.textContent = `${stats.completedCount} / ${stats.totalCount} Completed`;
    }
    if (state.elements.progressPercent) {
      state.elements.progressPercent.textContent = `${stats.completionPercentage}%`;
    }
    if (state.elements.progressFill) {
      state.elements.progressFill.style.width = `${stats.completionPercentage}%`;
    }
  }

  function renderSyllabusTree() {
    const tree = state.elements.syllabusTree;
    if (!tree || !state.bundle) return;
    tree.innerHTML = "";

    state.bundle.chapters.forEach((ch) => {
      const matchingExercises = ch.exercises.filter((ex) => {
        const matchesSearch =
          !state.searchQuery ||
          ex.name.toLowerCase().includes(state.searchQuery) ||
          ex.title.toLowerCase().includes(state.searchQuery) ||
          ch.title.toLowerCase().includes(state.searchQuery);

        const exState =
          (TerralingsStorage.state &&
            TerralingsStorage.state.exercises[ex.name]) ||
          {};
        const status = exState.status || "not_started";

        let matchesStatus = true;
        if (state.sidebarFilter === "todo")
          matchesStatus = status !== "completed";
        if (state.sidebarFilter === "done")
          matchesStatus = status === "completed";

        return matchesSearch && matchesStatus;
      });

      if (
        matchingExercises.length === 0 &&
        (state.searchQuery || state.sidebarFilter !== "all")
      ) {
        return;
      }

      const completedInChapter = ch.exercises.filter((ex) => {
        const s =
          TerralingsStorage.state && TerralingsStorage.state.exercises[ex.name];
        return s && s.status === "completed";
      }).length;

      const isAllDone = completedInChapter === ch.exercises.length;
      const isExpanded = state.expandedChapters.has(ch.id);

      const chapterDiv = document.createElement("div");
      chapterDiv.className = `syllabus-chapter ${isExpanded ? "open" : ""}`;

      const chapterHeader = document.createElement("div");
      chapterHeader.className = "chapter-header";
      chapterHeader.innerHTML = `
        <div class="chapter-title-row">
          <span class="accordion-arrow">▶</span>
          <span class="chapter-number">${ch.id}</span>
          <span class="chapter-name">${ch.title}</span>
        </div>
        <span class="chapter-count-badge ${isAllDone ? "all-done" : ""}">
          ${completedInChapter}/${ch.exercises.length}
        </span>
      `;

      chapterHeader.onclick = () => {
        if (state.expandedChapters.has(ch.id)) {
          state.expandedChapters.delete(ch.id);
        } else {
          state.expandedChapters.add(ch.id);
        }
        renderSyllabusTree();
      };

      const exerciseList = document.createElement("div");
      exerciseList.className = "exercise-list";

      matchingExercises.forEach((ex) => {
        const exState =
          (TerralingsStorage.state &&
            TerralingsStorage.state.exercises[ex.name]) ||
          {};
        const status = exState.status || "not_started";
        const isActive = ex.name === state.currentExerciseId;

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
          selectExercise(ex.name);
        };

        exerciseList.appendChild(item);
      });

      chapterDiv.appendChild(chapterHeader);
      chapterDiv.appendChild(exerciseList);
      tree.appendChild(chapterDiv);
    });
  }

  function selectExercise(exId) {
    if (!state.bundle || !state.bundle.exercises[exId]) return;
    state.currentExerciseId = exId;
    TerralingsStorage.state.lastActiveExerciseId = exId;
    TerralingsStorage.persist();

    const ex = state.bundle.exercises[exId];
    const ch = state.bundle.chapters.find((c) => c.id === ex.chapter);
    if (ch) state.expandedChapters.add(ch.id);

    // Update Header
    if (state.elements.chapterBadge)
      state.elements.chapterBadge.textContent = ch
        ? `${ch.id}. ${ch.title}`
        : `Chapter ${ex.chapter}`;
    if (state.elements.exTitle)
      state.elements.exTitle.textContent = `${ex.name} — ${ex.title}`;

    const exState = TerralingsStorage.getExerciseState(ex.name, ex.starter_code);
    if (state.elements.exStatusPill) {
      if (exState.status === "completed") {
        state.elements.exStatusPill.className = "status-indicator completed";
        state.elements.exStatusPill.textContent = "✓ Completed";
      } else if (exState.status === "in_progress") {
        state.elements.exStatusPill.className = "status-indicator in_progress";
        state.elements.exStatusPill.textContent = "⏳ In Progress";
      } else {
        state.elements.exStatusPill.className = "status-indicator";
        state.elements.exStatusPill.textContent = "○ Not Started";
      }
    }

    state.revealedHints = exState.hintsRevealed || 0;
    renderHints();

    // Set Editor Code
    if (state.editor) {
      const codeToLoad = exState.userCode || ex.starter_code;
      state.editor.setValue(codeToLoad);
    }

    if (state.isDiffMode) {
      updateDiffModels();
    }

    renderSyllabusTree();
    updateGlobalProgress();
  }

  function renderHints() {
    const ex = state.bundle && state.bundle.exercises[state.currentExerciseId];
    if (!ex) return;

    if (state.elements.hintBadge) {
      state.elements.hintBadge.textContent = `${state.revealedHints}/${ex.hints.length}`;
    }

    if (!state.elements.hintContent) return;

    if (state.revealedHints === 0) {
      state.elements.hintDrawer.style.display = "none";
      state.elements.hintContent.innerHTML = "";
      return;
    }

    state.elements.hintDrawer.style.display = "block";
    let html = "";
    for (let i = 0; i < state.revealedHints && i < ex.hints.length; i++) {
      html += `
        <div class="hint-item">
          <strong>Tier ${i + 1}:</strong>
          <p>${escapeHtml(ex.hints[i])}</p>
        </div>
      `;
    }
    state.elements.hintContent.innerHTML = html;
  }

  function toggleHint() {
    const ex = state.bundle && state.bundle.exercises[state.currentExerciseId];
    if (!ex || ex.hints.length === 0) return;

    if (state.revealedHints >= ex.hints.length) {
      state.revealedHints = 0;
    } else {
      state.revealedHints++;
    }

    TerralingsStorage.setHintsRevealed(
      state.currentExerciseId,
      state.revealedHints
    );
    renderHints();
  }

  function updateDiffModels() {
    if (!window.monaco || !state.diffEditor) return;
    const ex = state.bundle && state.bundle.exercises[state.currentExerciseId];
    if (!ex) return;

    const userCode = state.editor
      ? state.editor.getValue()
      : ex.starter_code;
    const solutionCode = ex.solution_code || "";

    if (state.originalModel) state.originalModel.dispose();
    if (state.modifiedModel) state.modifiedModel.dispose();

    state.originalModel = window.monaco.editor.createModel(
      userCode,
      "terraform"
    );
    state.modifiedModel = window.monaco.editor.createModel(
      solutionCode,
      "terraform"
    );

    state.diffEditor.setModel({
      original: state.originalModel,
      modified: state.modifiedModel,
    });
  }

  function toggleDiffView() {
    if (!state.diffEditor || !state.editor) return;
    state.isDiffMode = !state.isDiffMode;

    if (state.isDiffMode) {
      updateDiffModels();
      state.elements.monacoPane.style.display = "none";
      state.elements.diffPane.style.display = "block";
      state.elements.diffBtn.textContent = "✕ Close Diff";
      state.elements.diffBtn.classList.add("primary-btn");
    } else {
      state.elements.diffPane.style.display = "none";
      state.elements.monacoPane.style.display = "block";
      state.elements.diffBtn.textContent = "🔍 Solution Diff";
      state.elements.diffBtn.classList.remove("primary-btn");
      state.editor.layout();
    }
  }

  function runCurrentExercise() {
    if (!state.worker || !state.bundle) return;
    const code = state.editor ? state.editor.getValue() : "";
    const ex = state.bundle.exercises[state.currentExerciseId];

    updateStatus("loading", "Validating HCL in WebAssembly...");
    if (state.elements.terminalOutput) {
      state.elements.terminalOutput.innerHTML = `<span class="term-dim">⚡ Validating exercise [${escapeHtml(state.currentExerciseId)}] in Pyodide WebAssembly...</span>`;
    }

    state.worker.postMessage({
      type: "RUN_EXERCISE",
      exerciseId: state.currentExerciseId,
      code: code,
      rules: (ex && ex.rules) || {},
    });
  }

  function handleRunResult(result) {
    const term = state.elements.terminalOutput;

    if (result.passed) {
      TerralingsStorage.markCompleted(state.currentExerciseId, state.bundle);
      updateStatus("success", "✓ PASSED");
      if (term) {
        term.innerHTML = `
<span class="term-pass">✓ SUCCESS: Exercise [${escapeHtml(result.exerciseId)}] passed all validation checks!</span>
<span class="term-dim">Execution time: ${result.durationMs.toFixed(1)} ms</span>

${escapeHtml(result.output || "")}
<span class="term-dim">Press </span><span class="term-pass">Alt + Right</span><span class="term-dim"> to proceed to next exercise!</span>
`;
      }
    } else {
      updateStatus("error", "✕ FAILED");
      if (term) {
        let errorMsg = result.error || "Validation check failed";
        if (result.line) errorMsg += ` (at line ${result.line})`;

        term.innerHTML = `
<span class="term-fail">✕ COMPILATION / VALIDATION FAILED</span>
<span class="term-fail">${escapeHtml(errorMsg)}</span>

${escapeHtml(result.output || "")}
<span class="term-dim">Tip: Press </span><span class="term-warn">💡 Hint (H)</span><span class="term-dim"> or compare with </span><span class="term-warn">🔍 Solution Diff</span>.
`;
      }
    }

    selectExercise(state.currentExerciseId);
  }

  function navigateNext() {
    if (!state.bundle) return;
    const keys = Object.keys(state.bundle.exercises);
    const idx = keys.indexOf(state.currentExerciseId);
    if (idx >= 0 && idx < keys.length - 1) {
      selectExercise(keys[idx + 1]);
    }
  }

  function navigatePrev() {
    if (!state.bundle) return;
    const keys = Object.keys(state.bundle.exercises);
    const idx = keys.indexOf(state.currentExerciseId);
    if (idx > 0) {
      selectExercise(keys[idx - 1]);
    }
  }

  function toggleFullscreen() {
    const ws = state.elements.workspace;
    if (!ws) return;
    ws.classList.toggle("terralings-fullscreen");
    setTimeout(() => {
      if (state.editor) state.editor.layout();
      if (state.diffEditor) state.diffEditor.layout();
    }, 150);
  }

  /**
   * ==========================================================================
   * 4. Monaco Editor Initialization
   * ==========================================================================
   */
  function registerHclLanguage(monaco) {
    if (monaco.languages.getLanguages().some((l) => l.id === "terraform")) {
      return;
    }

    monaco.languages.register({ id: "terraform", extensions: [".tf", ".hcl"] });
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
          [
            /[a-zA-Z_]\w*/,
            {
              cases: {
                "@keywords": "keyword",
                "@typeKeywords": "type",
                "@default": "identifier",
              },
            },
          ],
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
  }

  function initMonacoInstance() {
    const monaco = window.monaco;
    if (!monaco || !state.elements.monacoPane) return;

    registerHclLanguage(monaco);

    const theme = getMonacoTheme();
    const ex = state.bundle && state.bundle.exercises[state.currentExerciseId];
    const exState =
      ex &&
      TerralingsStorage.getExerciseState(state.currentExerciseId, ex.starter_code);
    const initialCode =
      (exState && exState.userCode) || (ex && ex.starter_code) || "";

    state.editor = monaco.editor.create(state.elements.monacoPane, {
      value: initialCode,
      language: "terraform",
      theme: theme,
      automaticLayout: true,
      minimap: { enabled: false },
      scrollBeyondLastLine: false,
      fontSize: 13,
      fontFamily:
        "ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace",
      lineNumbers: "on",
      tabSize: 2,
    });

    state.editor.onDidChangeModelContent(() => {
      const code = state.editor.getValue();
      TerralingsStorage.saveExerciseCode(state.currentExerciseId, code);
    });

    state.diffEditor = monaco.editor.createDiffEditor(state.elements.diffPane, {
      theme: theme,
      automaticLayout: true,
      readOnly: true,
      minimap: { enabled: false },
      fontSize: 13,
      fontFamily:
        "ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace",
    });

    state.monacoLoaded = true;

    // Observe MkDocs color scheme changes
    const observer = new MutationObserver(() => {
      if (window.monaco) {
        window.monaco.editor.setTheme(getMonacoTheme());
      }
    });
    observer.observe(document.body, {
      attributes: true,
      attributeFilter: ["data-md-color-scheme"],
    });
  }

  function loadMonacoScript() {
    if (window.monaco) {
      initMonacoInstance();
      return;
    }

    if (window.require && typeof window.require === "function") {
      window.require.config({
        paths: {
          vs: "https://cdnjs.cloudflare.com/ajax/libs/monaco-editor/0.45.0/min/vs",
        },
      });
      window.require(["vs/editor/editor.main"], function () {
        initMonacoInstance();
      });
      return;
    }

    const loaderScript = document.createElement("script");
    loaderScript.src =
      "https://cdnjs.cloudflare.com/ajax/libs/monaco-editor/0.45.0/min/vs/loader.min.js";
    loaderScript.onload = function () {
      window.require.config({
        paths: {
          vs: "https://cdnjs.cloudflare.com/ajax/libs/monaco-editor/0.45.0/min/vs",
        },
      });
      window.require(["vs/editor/editor.main"], function () {
        initMonacoInstance();
      });
    };
    document.head.appendChild(loaderScript);
  }

  /**
   * ==========================================================================
   * 5. Pyodide Web Worker Setup
   * ==========================================================================
   */
  async function initWorker() {
    const workerUrl = resolveAssetUrl("playground-worker.js");
    let worker;

    try {
      worker = new Worker(workerUrl);
    } catch (e) {
      try {
        const resp = await fetch(workerUrl);
        const code = await resp.text();
        const blob = new Blob([code], { type: "application/javascript" });
        worker = new Worker(URL.createObjectURL(blob));
      } catch (err) {
        updateStatus("error", "Failed to spawn Web Worker: " + err.message);
        return;
      }
    }

    state.worker = worker;

    worker.onmessage = function (e) {
      const msg = e.data;
      if (!msg) return;

      if (msg.type === "STATUS") {
        updateStatus(msg.stage, msg.message);
        if (msg.stage === "ready") {
          state.workerReady = true;
          const ex =
            state.bundle && state.bundle.exercises[state.currentExerciseId];
          if (state.elements.terminalOutput && ex) {
            state.elements.terminalOutput.innerHTML = `
<span class="term-pass">🚀 Terralings WebAssembly Sandbox Ready (Python 3.12 + Pyodide)</span>
<span class="term-dim">Active Exercise: </span><strong>${escapeHtml(ex.name)}</strong> — ${escapeHtml(ex.title)}
<span class="term-dim">Click </span><span class="term-pass">▶ Run Solution</span><span class="term-dim"> (Ctrl+Enter) to evaluate HCL.</span>
`;
          }
        }
      } else if (msg.type === "RUN_RESULT") {
        handleRunResult(msg);
      }
    };

    worker.onerror = function (err) {
      updateStatus("error", "Worker error: " + (err.message || "Unknown error"));
    };

    worker.postMessage({
      type: "INIT",
      bundle: state.bundle,
    });
  }

  /**
   * ==========================================================================
   * 6. Events & Main Entry Point
   * ==========================================================================
   */
  function attachEventHandlers() {
    const el = state.elements;

    el.runBtn?.addEventListener("click", runCurrentExercise);
    el.hintBtn?.addEventListener("click", toggleHint);
    el.closeHintBtn?.addEventListener("click", () => {
      state.revealedHints = 0;
      TerralingsStorage.setHintsRevealed(state.currentExerciseId, 0);
      renderHints();
    });
    el.diffBtn?.addEventListener("click", toggleDiffView);
    el.prevBtn?.addEventListener("click", navigatePrev);
    el.nextBtn?.addEventListener("click", navigateNext);
    el.fullscreenBtn?.addEventListener("click", toggleFullscreen);

    el.resetCodeBtn?.addEventListener("click", () => {
      const ex = state.bundle && state.bundle.exercises[state.currentExerciseId];
      if (!ex) return;
      if (confirm(`Reset ${ex.name} to starter template?`)) {
        TerralingsStorage.resetExercise(state.currentExerciseId, ex.starter_code);
        selectExercise(state.currentExerciseId);
      }
    });

    el.exportBtn?.addEventListener("click", () => TerralingsStorage.exportJson());
    el.importBtn?.addEventListener("click", () => el.fileInput?.click());
    el.fileInput?.addEventListener("change", (e) => {
      const file = e.target.files?.[0];
      if (!file) return;
      const reader = new FileReader();
      reader.onload = (evt) => {
        const ok = TerralingsStorage.importJson(evt.target?.result, state.bundle);
        if (ok) {
          selectExercise(TerralingsStorage.state.lastActiveExerciseId || "primitives01");
          alert("Progress restored successfully!");
        } else {
          alert("Failed to parse progress JSON file.");
        }
      };
      reader.readAsText(file);
    });

    el.resetAllBtn?.addEventListener("click", () => {
      if (confirm("Reset all progress across all 56 exercises? This cannot be undone.")) {
        TerralingsStorage.resetAll(state.bundle);
        selectExercise("primitives01");
      }
    });

    el.searchInput?.addEventListener("input", (e) => {
      state.searchQuery = e.target.value.toLowerCase().trim();
      renderSyllabusTree();
    });

    el.filterTabs.forEach((tab) => {
      tab.addEventListener("click", (e) => {
        el.filterTabs.forEach((t) => t.classList.remove("active"));
        e.target.classList.add("active");
        state.sidebarFilter = e.target.dataset.filter;
        renderSyllabusTree();
      });
    });

    window.addEventListener("keydown", (e) => {
      if ((e.ctrlKey || e.metaKey) && e.key === "Enter") {
        e.preventDefault();
        runCurrentExercise();
      } else if (e.altKey && e.key === "ArrowRight") {
        e.preventDefault();
        navigateNext();
      } else if (e.altKey && e.key === "ArrowLeft") {
        e.preventDefault();
        navigatePrev();
      } else if (e.key === "F11") {
        e.preventDefault();
        toggleFullscreen();
      }
    });
  }

  async function initPlayground() {
    const container = document.getElementById("terralings-app");
    if (!container) return;
    state.container = container;

    renderPlaygroundSkeleton(container);
    bindElements(container);
    attachEventHandlers();

    try {
      updateStatus("loading", "⚡ Loading 56-exercise curriculum bundle...");
      const bundleUrl = resolveAssetUrl("playground-bundle.json");
      const resp = await fetch(bundleUrl);
      state.bundle = await resp.json();

      TerralingsStorage.init(state.bundle);
      const startExId =
        TerralingsStorage.state.lastActiveExerciseId || "primitives01";
      state.currentExerciseId = state.bundle.exercises[startExId]
        ? startExId
        : "primitives01";

      updateGlobalProgress();
      renderSyllabusTree();
      selectExercise(state.currentExerciseId);

      loadMonacoScript();
      await initWorker();
    } catch (err) {
      updateStatus("error", "Initialization failed: " + err.message);
      if (state.elements.terminalOutput) {
        state.elements.terminalOutput.innerHTML = `<span class="term-fail">❌ Failed to initialize playground: ${escapeHtml(err.message)}</span>`;
      }
    }
  }

  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", initPlayground);
  } else {
    initPlayground();
  }

  if (typeof window.document$ !== "undefined") {
    window.document$.subscribe(initPlayground);
  }
})();
