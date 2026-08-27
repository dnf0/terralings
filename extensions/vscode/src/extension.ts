import * as vscode from 'vscode';
import * as path from 'path';
import * as fs from 'fs';
import {
  startLspClient,
  stopLspClient
} from './lspClient';
import {
  getOutputChannel,
  disposeOutputChannel,
  runWatchTerminal,
  runTuiTerminal,
  runDoctor,
  runHint,
  runReset,
  runCurrentExercise
} from './cliRunner';
import { TerralingsStatusBar } from './statusBar';
import {
  TerralingsTreeDataProvider,
  ALL_EXERCISES,
  findExercise,
  resolveExerciseUri,
  findChapter
} from './treeProvider';
import { initStateWatcher } from './stateWatcher';

/**
 * Helper to extract an exercise name or target path from command parameters or tree items.
 */
function extractExerciseTarget(item?: unknown): string | undefined {
  if (!item) {
    return undefined;
  }
  if (typeof item === 'string') {
    return item;
  }
  if (typeof item === 'object') {
    if (
      'exercise' in item &&
      item.exercise &&
      typeof item.exercise === 'object' &&
      'name' in item.exercise
    ) {
      return String((item.exercise as { name: string }).name);
    }
    if ('name' in item && typeof (item as { name: unknown }).name === 'string') {
      return (item as { name: string }).name;
    }
    if ('fsPath' in item && typeof (item as { fsPath: unknown }).fsPath === 'string') {
      return (item as { fsPath: string }).fsPath;
    }
    if ('path' in item && typeof (item as { path: unknown }).path === 'string') {
      return (item as { path: string }).path;
    }
  }
  return undefined;
}

/**
 * Checks configuration and workspace state to conditionally launch the welcome walkthrough.
 */
function checkAutoOpenWalkthrough(context: vscode.ExtensionContext): void {
  const config = vscode.workspace.getConfiguration('terralings');
  const autoOpen = config.get<boolean>('autoOpenWalkthrough', true);
  if (!autoOpen) {
    return;
  }

  const hasSeenWalkthrough = context.globalState.get<boolean>(
    'terralings.hasSeenWalkthrough',
    false
  );
  if (hasSeenWalkthrough) {
    return;
  }

  const workspaceRoot = vscode.workspace.workspaceFolders?.[0]?.uri.fsPath;
  const statePath = workspaceRoot ? path.join(workspaceRoot, '.terralings', 'state.json') : undefined;
  const stateExists = statePath ? fs.existsSync(statePath) : false;

  if (!stateExists) {
    context.globalState.update('terralings.hasSeenWalkthrough', true);
    vscode.commands.executeCommand(
      'workbench.action.openWalkthrough',
      'dnf0.terralings-vscode#terralings.welcome',
      false
    );
  }
}

/**
 * Extension activation entry point.
 */
export async function activate(context: vscode.ExtensionContext): Promise<void> {
  // 1. Initialize output channel
  const outputChannel = getOutputChannel();
  outputChannel.appendLine('[Terralings] Initializing extension...');

  // 2. Start LSP client
  startLspClient(context).catch((err) => {
    outputChannel.appendLine(
      `[Terralings] LSP client startup error: ${err instanceof Error ? err.message : String(err)}`
    );
  });

  // 3. Initialize Tree Data Provider
  const treeDataProvider = new TerralingsTreeDataProvider();
  context.subscriptions.push(treeDataProvider);
  context.subscriptions.push(
    vscode.window.registerTreeDataProvider('terralings-tree', treeDataProvider)
  );

  // 4. Initialize Status Bar and State Watcher
  const statusBar = new TerralingsStatusBar();
  context.subscriptions.push(statusBar);

  const stateWatcherDisposable = initStateWatcher(treeDataProvider, statusBar, context);
  context.subscriptions.push(stateWatcherDisposable);

  // Initial status bar sync
  const initialProgress = treeDataProvider.getProgress();
  statusBar.update(initialProgress.completed, initialProgress.total);

  // 5. Register Extension Commands
  context.subscriptions.push(
    vscode.commands.registerCommand('terralings.openExercise', async (item?: unknown) => {
      let targetName = extractExerciseTarget(item);

      if (!targetName) {
        // Show QuickPick with all exercises
        const items = ALL_EXERCISES.map((ex) => {
          const status = treeDataProvider.getExerciseStatus(ex.name);
          const icon =
            status === 'passed'
              ? '✅'
              : status === 'failed'
              ? '❌'
              : status === 'in_progress'
              ? '▶️'
              : '⚪';
          const chapter = findChapter(ex.chapterName);
          return {
            label: `${icon} ${ex.name}`,
            description: ex.title,
            detail: `${chapter ? `${chapter.number}. ${chapter.title}` : ex.chapterName} • ${ex.path}`,
            exercise: ex
          };
        });

        const selected = await vscode.window.showQuickPick(items, {
          placeHolder: 'Select an exercise to open...',
          matchOnDescription: true,
          matchOnDetail: true
        });

        if (selected) {
          targetName = selected.exercise.name;
        }
      }

      if (targetName) {
        const found = findExercise(targetName);
        const filePath = found ? found.path : targetName;
        const uri = resolveExerciseUri(filePath);
        await vscode.commands.executeCommand('vscode.open', uri);
      }
    })
  );

  context.subscriptions.push(
    vscode.commands.registerCommand('terralings.runCurrent', async (item?: unknown) => {
      const target = extractExerciseTarget(item);
      await runCurrentExercise(target);
    })
  );

  context.subscriptions.push(
    vscode.commands.registerCommand('terralings.watch', () => {
      runWatchTerminal();
    })
  );

  context.subscriptions.push(
    vscode.commands.registerCommand('terralings.tui', () => {
      runTuiTerminal();
    })
  );

  context.subscriptions.push(
    vscode.commands.registerCommand('terralings.hint', async (item?: unknown) => {
      const target = extractExerciseTarget(item);
      await runHint(target);
    })
  );

  context.subscriptions.push(
    vscode.commands.registerCommand('terralings.reset', async (item?: unknown) => {
      const target = extractExerciseTarget(item);
      await runReset(target);
    })
  );

  context.subscriptions.push(
    vscode.commands.registerCommand('terralings.doctor', async () => {
      await runDoctor();
    })
  );

  context.subscriptions.push(
    vscode.commands.registerCommand('terralings.tour', () => {
      vscode.commands.executeCommand(
        'workbench.action.openWalkthrough',
        'dnf0.terralings-vscode#terralings.welcome',
        false
      );
    })
  );

  context.subscriptions.push(
    vscode.commands.registerCommand('terralings.refreshTree', () => {
      treeDataProvider.refresh();
      const progress = treeDataProvider.getProgress();
      statusBar.update(progress.completed, progress.total);
    })
  );

  // 6. Check and trigger walkthrough if appropriate
  checkAutoOpenWalkthrough(context);

  outputChannel.appendLine('[Terralings] Extension activation complete.');
}

/**
 * Extension deactivation lifecycle hook.
 */
export async function deactivate(): Promise<void> {
  await stopLspClient();
  disposeOutputChannel();
}
