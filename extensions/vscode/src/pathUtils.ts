import * as vscode from 'vscode';
import * as fs from 'fs';
import * as path from 'path';
import * as os from 'os';

/**
 * Returns a valid, writable workspace directory.
 * If no folder is open in VS Code, checks standard locations (~/repos/terralings, ~/terralings)
 * before falling back safely to os.homedir(). Never returns '/' or read-only system root.
 */
export function getEffectiveWorkspaceRoot(repoName: string = 'terralings'): string {
  const wsFolder = vscode.workspace.workspaceFolders?.[0]?.uri.fsPath;
  if (wsFolder) {
    return wsFolder;
  }

  // Check known existing terralings directories
  const home = os.homedir();
  const candidates = [
    path.join(home, 'repos', repoName),
    path.join(home, repoName),
    path.join(home, 'Developer', repoName),
  ];

  for (const candidate of candidates) {
    if (fs.existsSync(candidate)) {
      try {
        if (fs.statSync(candidate).isDirectory()) {
          return candidate;
        }
      } catch {
        // ignore access errors and continue
      }
    }
  }

  return home;
}

/**
 * Robustly resolves an exercise relative or absolute path against the current workspace
 * or common locations, checking:
 * 1. Absolute paths
 * 2. Direct resolution against effective workspace root (e.g. root/exercises/01_primitives/...)
 * 3. Resolution with "exercises/" prefix added or stripped
 * 4. Checking multi-file module directory main.tf files (e.g. exercises/08_modules/module01)
 * 5. Ascending parent directory traversal up to 6 levels (when workspace is a subfolder)
 * 6. Checking standard ~/repos/terralings and ~/terralings locations
 */
export function resolveExercisePath(exPath: string, workspaceRoot?: string): string {
  const root = workspaceRoot || getEffectiveWorkspaceRoot();

  // 1. If path is already absolute and exists on disk
  if (path.isAbsolute(exPath)) {
    if (fs.existsSync(exPath)) {
      const stat = fs.statSync(exPath);
      if (stat.isDirectory()) {
        const mainTf = path.join(exPath, 'main.tf');
        if (fs.existsSync(mainTf)) {
          return mainTf;
        }
      }
      return exPath;
    }
  }

  const checkCandidate = (candidatePath: string): string | undefined => {
    if (fs.existsSync(candidatePath)) {
      try {
        const stat = fs.statSync(candidatePath);
        if (stat.isDirectory()) {
          const mainTf = path.join(candidatePath, 'main.tf');
          if (fs.existsSync(mainTf)) {
            return mainTf;
          }
        }
        return candidatePath;
      } catch {
        return candidatePath;
      }
    }
    return undefined;
  };

  // 2. Direct resolve with root
  const directPath = path.resolve(root, exPath);
  const directResolved = checkCandidate(directPath);
  if (directResolved) {
    return directResolved;
  }

  // 3. Try prepending "exercises/" if not present
  if (!exPath.startsWith('exercises/') && !exPath.startsWith('exercises\\')) {
    const withExercises = path.resolve(root, 'exercises', exPath);
    const withExercisesResolved = checkCandidate(withExercises);
    if (withExercisesResolved) {
      return withExercisesResolved;
    }
  }

  // 4. Try stripping "exercises/" prefix if present
  if (exPath.startsWith('exercises/') || exPath.startsWith('exercises\\')) {
    const stripped = exPath.replace(/^exercises[/\\]/, '');
    const strippedPath = path.resolve(root, stripped);
    const strippedResolved = checkCandidate(strippedPath);
    if (strippedResolved) {
      return strippedResolved;
    }
  }

  // 5. If root is in a subfolder (e.g. exercises/ or exercises/01_primitives/), search parent directories up to 6 levels
  let cur = root;
  for (let i = 0; i < 6; i++) {
    const candidateFull = path.resolve(cur, exPath);
    const fullRes = checkCandidate(candidateFull);
    if (fullRes) {
      return fullRes;
    }

    if (!exPath.startsWith('exercises/') && !exPath.startsWith('exercises\\')) {
      const candidateWithEx = path.resolve(cur, 'exercises', exPath);
      const withExRes = checkCandidate(candidateWithEx);
      if (withExRes) {
        return withExRes;
      }
    } else {
      const stripped = exPath.replace(/^exercises[/\\]/, '');
      const candidateStripped = path.resolve(cur, stripped);
      const strippedRes = checkCandidate(candidateStripped);
      if (strippedRes) {
        return strippedRes;
      }
    }

    const parent = path.dirname(cur);
    if (parent === cur) {
      break;
    }
    cur = parent;
  }

  // 6. Check standard candidate directories in home
  const standardLocations = [
    path.join(os.homedir(), 'repos', 'terralings'),
    path.join(os.homedir(), 'terralings'),
    path.join(os.homedir(), 'Developer', 'terralings'),
  ];

  for (const loc of standardLocations) {
    const cand = path.resolve(loc, exPath);
    const candRes = checkCandidate(cand);
    if (candRes) {
      return candRes;
    }
  }

  // 7. Default fallback to direct path
  return directPath;
}

/**
 * Locates the authoritative `.terralings/state.json` file across:
 * 1. Active editor document directory & parent hierarchy (up to 8 levels)
 * 2. Open workspace folders & parent hierarchies (up to 6 levels)
 * 3. Effective workspace root & parent hierarchy (up to 6 levels)
 * 4. Standard home directory repositories (~/repos/terralings, ~/terralings, etc.)
 */
export function findStateJsonPath(workspaceRoot?: string): string | undefined {
  const checkHierarchy = (startDir: string, maxDepth: number = 8): string | undefined => {
    let cur = startDir;
    for (let i = 0; i < maxDepth; i++) {
      const candidate = path.join(cur, '.terralings', 'state.json');
      if (fs.existsSync(candidate)) {
        return candidate;
      }
      const parent = path.dirname(cur);
      if (parent === cur) {
        break;
      }
      cur = parent;
    }
    return undefined;
  };

  // If explicit root provided, search it first
  if (workspaceRoot) {
    const foundInRoot = checkHierarchy(workspaceRoot, 8);
    if (foundInRoot) {
      return foundInRoot;
    }
  }

  // 1. Check active editor document parent hierarchy
  const activeFile = vscode.window.activeTextEditor?.document?.fileName;
  if (activeFile) {
    const foundInDoc = checkHierarchy(path.dirname(activeFile), 8);
    if (foundInDoc) {
      return foundInDoc;
    }
  }

  // 2. Check open workspace folders
  const wsFolders = vscode.workspace.workspaceFolders;
  if (wsFolders && wsFolders.length > 0) {
    for (const folder of wsFolders) {
      const foundInWs = checkHierarchy(folder.uri.fsPath, 6);
      if (foundInWs) {
        return foundInWs;
      }
    }
  }

  // 3. Check effective root and standard home paths only if explicit root wasn't passed
  if (!workspaceRoot) {
    const root = getEffectiveWorkspaceRoot();
    const foundInEffective = checkHierarchy(root, 6);
    if (foundInEffective) {
      return foundInEffective;
    }

    // 4. Check standard home paths
    const standardLocations = [
      path.join(os.homedir(), 'repos', 'terralings', '.terralings', 'state.json'),
      path.join(os.homedir(), 'terralings', '.terralings', 'state.json'),
      path.join(os.homedir(), 'Developer', 'terralings', '.terralings', 'state.json'),
    ];
    for (const loc of standardLocations) {
      if (fs.existsSync(loc)) {
        return loc;
      }
    }
  }

  return undefined;
}

