import * as vscode from 'vscode';
import * as fs from 'fs';
import { TerralingsTreeDataProvider } from './treeProvider';
import { TerralingsStatusBar } from './statusBar';
import { findStateJsonPath } from './pathUtils';

/**
 * Initializes multi-layer state synchronization for Terralings:
 * 1. Direct Node fs.watchFile polling on `.terralings/state.json` (bypasses gitignore & hidden folder watcher exclusions)
 * 2. VS Code FileSystemWatcher on state and exercise source files
 * 3. Heartbeat polling interval (1s) to immediately reflect terminal background runs (watch/tui/run)
 * 4. VS Code lifecycle triggers (file saves, tab changes, window focus events)
 *
 * Debounces UI tree & status bar updates to prevent UI thrashing.
 */
export function initStateWatcher(
  treeProvider: TerralingsTreeDataProvider,
  statusBar: TerralingsStatusBar
): vscode.Disposable {
  const disposables: vscode.Disposable[] = [];
  let debounceTimer: NodeJS.Timeout | undefined;
  let pollingTimer: NodeJS.Timeout | undefined;
  let activeWatchedStateFile: string | undefined;

  const triggerUpdate = () => {
    if (debounceTimer) {
      clearTimeout(debounceTimer);
    }
    debounceTimer = setTimeout(() => {
      treeProvider.refresh();
      const progress = treeProvider.getProgress();
      statusBar.update(progress.completed, progress.total);
    }, 80);
  };

  // 1. VS Code workspace watchers
  const stateWatcher = vscode.workspace.createFileSystemWatcher('**/.terralings/state.json');
  disposables.push(
    stateWatcher,
    stateWatcher.onDidChange(triggerUpdate),
    stateWatcher.onDidCreate(triggerUpdate),
    stateWatcher.onDidDelete(triggerUpdate)
  );

  const exerciseWatcher = vscode.workspace.createFileSystemWatcher('**/*.{tf,hcl}');
  disposables.push(
    exerciseWatcher,
    exerciseWatcher.onDidChange(triggerUpdate),
    exerciseWatcher.onDidCreate(triggerUpdate),
    exerciseWatcher.onDidDelete(triggerUpdate)
  );

  // 2. VS Code lifecycle event triggers
  disposables.push(
    vscode.workspace.onDidSaveTextDocument((doc) => {
      const name = doc.fileName;
      if (name.endsWith('.tf') || name.endsWith('.hcl') || name.endsWith('.json')) {
        triggerUpdate();
      }
    }),
    vscode.window.onDidChangeActiveTextEditor(() => {
      triggerUpdate();
    }),
    vscode.window.onDidChangeWindowState((state) => {
      if (state.focused) {
        triggerUpdate();
      }
    })
  );

  // 3. Direct fs.watchFile tracking on authoritative state.json
  const setupDirectFileWatch = () => {
    const stateFile = findStateJsonPath();
    if (stateFile && stateFile !== activeWatchedStateFile) {
      if (activeWatchedStateFile) {
        try {
          fs.unwatchFile(activeWatchedStateFile);
        } catch {
          // ignore unwatch errors
        }
      }
      activeWatchedStateFile = stateFile;
      try {
        fs.watchFile(stateFile, { interval: 500 }, (curr, prev) => {
          if (curr.mtimeMs !== prev.mtimeMs || curr.size !== prev.size) {
            triggerUpdate();
          }
        });
      } catch {
        // ignore watch errors
      }
    }
  };

  setupDirectFileWatch();

  // 4. Lightweight polling heartbeat (1s) to guarantee zero desync from external terminal runs
  let lastMtime = 0;
  pollingTimer = setInterval(() => {
    setupDirectFileWatch();
    if (activeWatchedStateFile && fs.existsSync(activeWatchedStateFile)) {
      try {
        const stat = fs.statSync(activeWatchedStateFile);
        if (stat.mtimeMs !== lastMtime) {
          lastMtime = stat.mtimeMs;
          triggerUpdate();
        }
      } catch {
        // ignore stat errors
      }
    }
  }, 1000);

  const compositeDisposable = vscode.Disposable.from(
    ...disposables,
    new vscode.Disposable(() => {
      if (debounceTimer) {
        clearTimeout(debounceTimer);
        debounceTimer = undefined;
      }
      if (pollingTimer) {
        clearInterval(pollingTimer);
        pollingTimer = undefined;
      }
      if (activeWatchedStateFile) {
        try {
          fs.unwatchFile(activeWatchedStateFile);
        } catch {
          // ignore
        }
        activeWatchedStateFile = undefined;
      }
    })
  );

  return compositeDisposable;
}
