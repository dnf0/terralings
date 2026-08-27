import * as cp from 'child_process';
import * as vscode from 'vscode';
import * as path from 'path';
import { findTerralingsBinary } from './lspClient';

let outputChannel: vscode.OutputChannel | undefined;

/**
 * Returns or creates the shared Terralings output channel.
 */
export function getOutputChannel(): vscode.OutputChannel {
  if (!outputChannel) {
    outputChannel = vscode.window.createOutputChannel('Terralings');
  }
  return outputChannel;
}

/**
 * Disposes the shared Terralings output channel.
 */
export function disposeOutputChannel(): void {
  if (outputChannel) {
    outputChannel.dispose();
    outputChannel = undefined;
  }
}

/**
 * Resolves the path to the terralings executable for CLI commands.
 */
export function getTerralingsBinary(): string {
  const workspaceRoot = vscode.workspace.workspaceFolders?.[0]?.uri.fsPath;
  return findTerralingsBinary(workspaceRoot);
}

/**
 * Normalizes an exercise identifier from a filepath, bare name, or currently active editor document.
 */
export function normalizeExerciseName(input?: string): string | undefined {
  if (!input) {
    const activeDoc = vscode.window.activeTextEditor?.document;
    if (activeDoc) {
      const fileName = activeDoc.fileName;
      if (fileName.endsWith('.tf') || fileName.endsWith('.hcl') || fileName.includes('exercises')) {
        input = fileName;
      }
    }
  }
  if (!input) {
    return undefined;
  }
  const trimmed = input.trim();
  const base = path.basename(trimmed);
  const ext = path.extname(base);
  if (ext === '.tf' || ext === '.hcl' || ext === '.json') {
    return base.slice(0, -ext.length);
  }
  return base;
}

/**
 * Checks for an existing active terminal with the given name or creates a new one,
 * then sends the terralings command to it and reveals the terminal.
 */
export function openTerminalCommand(name: string, args: string[]): vscode.Terminal {
  const workspaceRoot = vscode.workspace.workspaceFolders?.[0]?.uri.fsPath;
  let terminal = vscode.window.terminals.find((t) => t.name === name);
  if (!terminal) {
    terminal = vscode.window.createTerminal({
      name,
      cwd: workspaceRoot
    });
  }

  const bin = getTerralingsBinary();
  const isWindows = process.platform === 'win32';
  let cmd: string;
  if (isWindows && bin.includes(' ')) {
    cmd = `& "${bin}" ${args.join(' ')}`;
  } else if (bin.includes(' ')) {
    cmd = `"${bin}" ${args.join(' ')}`;
  } else {
    cmd = `${bin} ${args.join(' ')}`;
  }

  terminal.show();
  terminal.sendText(cmd);
  return terminal;
}

/**
 * Launches `terralings watch` in a dedicated terminal.
 */
export function runWatchTerminal(): vscode.Terminal {
  return openTerminalCommand('Terralings Watch', ['watch']);
}

/**
 * Launches `terralings tui` in a dedicated terminal.
 */
export function runTuiTerminal(): vscode.Terminal {
  return openTerminalCommand('Terralings TUI', ['tui']);
}

/**
 * Executes a terralings CLI command, writing formatted output to the shared OutputChannel.
 */
export function runCliCommand(args: string[], title: string): Promise<string> {
  const channel = getOutputChannel();
  const binary = getTerralingsBinary();
  const workspaceRoot = vscode.workspace.workspaceFolders?.[0]?.uri.fsPath;

  const timestamp = new Date().toISOString().replace('T', ' ').slice(0, 19);
  channel.appendLine(`\n=== [${timestamp}] ${title} ===`);
  channel.appendLine(`$ ${binary} ${args.join(' ')}\n`);

  return new Promise<string>((resolve, reject) => {
    cp.execFile(
      binary,
      args,
      {
        cwd: workspaceRoot,
        env: { ...process.env, NO_COLOR: '1' }
      },
      (error, stdout, stderr) => {
        if (stdout) {
          channel.append(stdout);
        }
        if (stderr) {
          channel.append(stderr);
        }
        if (error) {
          channel.appendLine(`\n[Command failed with exit code ${error.code ?? 1}]`);
          const combined = (stdout + '\n' + stderr).trim();
          reject(new Error(combined || error.message));
        } else {
          channel.appendLine(`\n[Command completed successfully]`);
          resolve(stdout.trim());
        }
      }
    );
  });
}

/**
 * Runs `terralings doctor` and displays diagnostic results in the output channel.
 */
export async function runDoctor(): Promise<void> {
  const channel = getOutputChannel();
  channel.show(true);
  try {
    await runCliCommand(['doctor'], 'Environment Diagnostics (Doctor)');
    vscode.window.showInformationMessage('Terralings: Environment diagnostics passed successfully.');
  } catch (err) {
    vscode.window.showWarningMessage(
      'Terralings: Environment diagnostics reported warnings or missing dependencies. See output channel for details.'
    );
  }
}

/**
 * Shows the hint for an exercise in an information notification and the output channel.
 */
export async function runHint(exerciseName?: string): Promise<void> {
  let target = normalizeExerciseName(exerciseName);
  if (!target) {
    const input = await vscode.window.showInputBox({
      prompt: 'Enter the exercise name for hints (e.g. syntax01)',
      placeHolder: 'syntax01'
    });
    if (input) {
      target = normalizeExerciseName(input.trim());
    }
  }
  if (!target) {
    return;
  }

  const channel = getOutputChannel();
  try {
    const hint = await runCliCommand(['hint', target], `Hint for ${target}`);
    channel.show(true);
    vscode.window.showInformationMessage(`💡 Hint for ${target}:\n\n${hint}`);
  } catch (err) {
    channel.show(true);
    vscode.window.showErrorMessage(
      `Failed to retrieve hint for "${target}": ${err instanceof Error ? err.message : String(err)}`
    );
  }
}

/**
 * Prompts user for confirmation and resets an exercise template to its initial state.
 */
export async function runReset(exerciseName?: string): Promise<void> {
  let target = normalizeExerciseName(exerciseName);
  if (!target) {
    const input = await vscode.window.showInputBox({
      prompt: 'Enter the exercise name to reset (e.g. syntax01)',
      placeHolder: 'syntax01'
    });
    if (input) {
      target = normalizeExerciseName(input.trim());
    }
  }
  if (!target) {
    return;
  }

  const confirmation = await vscode.window.showWarningMessage(
    `Are you sure you want to reset exercise "${target}"? All your changes in this exercise will be replaced with the original starting template.`,
    { modal: true },
    'Reset Exercise'
  );

  if (confirmation !== 'Reset Exercise') {
    return;
  }

  const channel = getOutputChannel();
  try {
    await runCliCommand(['reset', target, '-f'], `Reset Exercise ${target}`);
    vscode.window.showInformationMessage(`Exercise "${target}" has been reset to its starting template.`);
  } catch (err) {
    channel.show(true);
    vscode.window.showErrorMessage(
      `Failed to reset exercise "${target}": ${err instanceof Error ? err.message : String(err)}`
    );
  }
}

/**
 * Verifies the specified or active exercise file.
 */
export async function runCurrentExercise(exercisePathOrName?: string): Promise<void> {
  const target = normalizeExerciseName(exercisePathOrName);
  if (!target) {
    vscode.window.showWarningMessage(
      'No active Terralings exercise found. Open an exercise file (.tf) to verify.'
    );
    return;
  }

  const channel = getOutputChannel();
  try {
    await runCliCommand(['run', target], `Verify Exercise ${target}`);
    vscode.window.showInformationMessage(`✅ Exercise "${target}" passed successfully!`);
  } catch (err) {
    channel.show(true);
    vscode.window.showErrorMessage(
      `❌ Exercise "${target}" failed verification. See Terralings output channel for details.`
    );
  }
}

