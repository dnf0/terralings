import * as vscode from 'vscode';
import * as path from 'path';
import * as fs from 'fs';
import {
  LanguageClient,
  LanguageClientOptions,
  ServerOptions,
  TransportKind
} from 'vscode-languageclient/node';
import { getEffectiveWorkspaceRoot } from './pathUtils';

let client: LanguageClient | undefined;
let fileWatcher: vscode.FileSystemWatcher | undefined;

/**
 * Resolves the path to the terralings executable.
 * Priority:
 * 1. terralings.binaryPath configuration in VS Code settings.
 * 2. Workspace root ./bin/terralings.
 * 3. System $PATH lookup.
 * 4. Fallback to 'terralings'.
 */
export function findTerralingsBinary(workspaceRoot?: string): string {
  const config = vscode.workspace.getConfiguration('terralings');
  const configuredPath = config.get<string>('binaryPath')?.trim();

  const isWindows = process.platform === 'win32';
  const binName = isWindows ? 'terralings.exe' : 'terralings';

  // 1. Check custom binaryPath setting if provided and not default
  if (configuredPath && configuredPath !== 'terralings' && configuredPath.length > 0) {
    if (path.isAbsolute(configuredPath)) {
      if (fs.existsSync(configuredPath)) {
        return configuredPath;
      }
    } else if (workspaceRoot) {
      const resolved = path.join(workspaceRoot, configuredPath);
      if (fs.existsSync(resolved)) {
        return resolved;
      }
    }
    return configuredPath;
  }

  // 2. Check workspace root ./bin/terralings
  const root = workspaceRoot ?? getEffectiveWorkspaceRoot();
  if (root) {
    const workspaceBin = path.join(root, 'bin', binName);
    if (fs.existsSync(workspaceBin)) {
      return workspaceBin;
    }
  }

  // 3. Search system $PATH
  const envPath = process.env.PATH || '';
  const pathDirs = envPath.split(path.delimiter);
  for (const dir of pathDirs) {
    if (!dir) {
      continue;
    }
    const candidate = path.join(dir, binName);
    try {
      if (fs.existsSync(candidate)) {
        return candidate;
      }
    } catch {
      // Ignore directory stat errors
    }
  }

  // 4. Default fallback
  return configuredPath || 'terralings';
}

/**
 * Starts the Terralings LSP client and connects via stdio to `terralings lsp`.
 */
export async function startLspClient(
  context: vscode.ExtensionContext
): Promise<LanguageClient | undefined> {
  const config = vscode.workspace.getConfiguration('terralings');
  const enableLsp = config.get<boolean>('enableLsp', true);
  if (!enableLsp) {
    return undefined;
  }

  if (client) {
    return client;
  }

  const workspaceRoot = getEffectiveWorkspaceRoot();
  const binaryPath = findTerralingsBinary(workspaceRoot);

  const serverOptions: ServerOptions = {
    command: binaryPath,
    args: ['lsp'],
    transport: TransportKind.stdio,
    options: { cwd: workspaceRoot }
  };

  if (!fileWatcher) {
    fileWatcher = vscode.workspace.createFileSystemWatcher('**/*.{tf,hcl,json}');
  }

  const clientOptions: LanguageClientOptions = {
    documentSelector: [
      { scheme: 'file', language: 'terraform' },
      { scheme: 'file', language: 'hcl' },
      { scheme: 'file', pattern: '**/*.{tf,tfvars,hcl}' }
    ],
    synchronize: {
      fileEvents: fileWatcher
    }
  };

  client = new LanguageClient(
    'terralings-lsp',
    'Terralings Language Server',
    serverOptions,
    clientOptions
  );

  try {
    await client.start();
    context.subscriptions.push({
      dispose: () => {
        stopLspClient().catch(console.error);
      }
    });
    return client;
  } catch (err) {
    vscode.window.showErrorMessage(
      `Failed to start Terralings LSP client: ${err instanceof Error ? err.message : String(err)}`
    );
    client = undefined;
    if (fileWatcher) {
      fileWatcher.dispose();
      fileWatcher = undefined;
    }
    return undefined;
  }
}

/**
 * Stops the running Terralings LSP client if active.
 */
export async function stopLspClient(): Promise<void> {
  if (fileWatcher) {
    fileWatcher.dispose();
    fileWatcher = undefined;
  }
  if (!client) {
    return;
  }
  const runningClient = client;
  client = undefined;
  await runningClient.stop();
}

/**
 * Returns the active LanguageClient instance if started.
 */
export function getLspClient(): LanguageClient | undefined {
  return client;
}
