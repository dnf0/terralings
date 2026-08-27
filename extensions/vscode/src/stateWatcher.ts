import * as vscode from 'vscode';
import { TerralingsTreeDataProvider } from './treeProvider';
import { TerralingsStatusBar } from './statusBar';

/**
 * Initializes file system watchers on `.terralings/state.json` and exercise files,
 * debouncing updates before refreshing the curriculum tree provider and status bar.
 */
export function initStateWatcher(
  treeProvider: TerralingsTreeDataProvider,
  statusBar: TerralingsStatusBar,
  context: vscode.ExtensionContext
): vscode.Disposable {
  const disposables: vscode.Disposable[] = [];
  let debounceTimer: NodeJS.Timeout | undefined;

  const triggerUpdate = () => {
    if (debounceTimer) {
      clearTimeout(debounceTimer);
    }
    debounceTimer = setTimeout(() => {
      treeProvider.refresh();
      const progress = treeProvider.getProgress();
      statusBar.update(progress.completed, progress.total);
    }, 100);
  };

  // Watch .terralings/state.json for real-time progress transitions
  const stateWatcher = vscode.workspace.createFileSystemWatcher('**/.terralings/state.json');
  stateWatcher.onDidChange(triggerUpdate);
  stateWatcher.onDidCreate(triggerUpdate);
  stateWatcher.onDidDelete(triggerUpdate);
  disposables.push(stateWatcher);

  // Watch exercise files for user code modifications
  const exerciseWatcher = vscode.workspace.createFileSystemWatcher('**/exercises/**/*.{tf,hcl}');
  exerciseWatcher.onDidChange(triggerUpdate);
  exerciseWatcher.onDidCreate(triggerUpdate);
  exerciseWatcher.onDidDelete(triggerUpdate);
  disposables.push(exerciseWatcher);

  const compositeDisposable = vscode.Disposable.from(
    ...disposables,
    new vscode.Disposable(() => {
      if (debounceTimer) {
        clearTimeout(debounceTimer);
        debounceTimer = undefined;
      }
    })
  );

  context.subscriptions.push(compositeDisposable);
  return compositeDisposable;
}
