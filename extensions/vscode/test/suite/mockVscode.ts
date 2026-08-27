import * as path from 'path';

let isMocked = false;

export const mockState: {
  config: Record<string, unknown>;
  workspaceRoot: string;
  activeFileName?: string;
} = {
  config: {},
  workspaceRoot: path.resolve(__dirname, '../../../../'),
  activeFileName: undefined
};

export function setMockConfig(key: string, value: unknown) {
  mockState.config[key] = value;
}

export function resetMockState() {
  mockState.config = {};
  mockState.workspaceRoot = path.resolve(__dirname, '../../../../');
  mockState.activeFileName = undefined;
}

export function setupMockVscode() {
  if (isMocked) {
    return;
  }

  // eslint-disable-next-line @typescript-eslint/no-var-requires
  const Module = require('module');
  const originalLoad = Module._load;

  class Position {
    constructor(public line: number, public character: number) {}
  }

  class Range {
    constructor(public start: Position, public end: Position) {}
  }

  class Location {
    constructor(public uri: unknown, public rangeOrPosition: unknown) {}
  }

  class Diagnostic {
    constructor(public range: Range, public message: string, public severity?: unknown) {}
  }

  class CancellationError extends Error {
    constructor() {
      super('Canceled');
      this.name = 'CancellationError';
    }
  }

  class Disposable {
    static from(...disposables: { dispose(): unknown }[]) {
      return new Disposable(() => disposables.forEach((d) => d.dispose()));
    }
    constructor(private callOnDispose: () => unknown) {}
    dispose() {
      if (this.callOnDispose) {
        this.callOnDispose();
      }
    }
  }

  const mockVscode = {
    version: '1.85.0',
    Position,
    Range,
    Location,
    Diagnostic,
    CancellationError,
    Disposable,
    DiagnosticSeverity: { Error: 0, Warning: 1, Information: 2, Hint: 3 },
    DiagnosticTag: { Unnecessary: 1, Deprecated: 2 },
    TreeItemCollapsibleState: {
      None: 0,
      Collapsed: 1,
      Expanded: 2
    },
    TreeItem: class TreeItem {
      public label: string | { label: string };
      public collapsibleState?: number;
      public description?: string;
      public tooltip?: string;
      public iconPath?: unknown;
      public contextValue?: string;
      public id?: string;
      public command?: unknown;

      constructor(label: string | { label: string }, collapsibleState?: number) {
        this.label = label;
        this.collapsibleState = collapsibleState ?? 0;
      }
    },
    ThemeIcon: class ThemeIcon {
      public id: string;
      public color?: unknown;

      constructor(id: string, color?: unknown) {
        this.id = id;
        this.color = color;
      }
    },
    ThemeColor: class ThemeColor {
      public id: string;

      constructor(id: string) {
        this.id = id;
      }
    },
    EventEmitter: class EventEmitter<T> {
      private listeners: Array<(e: T) => void> = [];

      public event = (listener: (e: T) => void) => {
        this.listeners.push(listener);
        return { dispose: () => {} };
      };

      public fire(data: T): void {
        for (const listener of this.listeners) {
          listener(data);
        }
      }

      public dispose(): void {
        this.listeners = [];
      }
    },
    Uri: {
      file: (filePath: string) => ({
        fsPath: path.resolve(filePath),
        path: path.resolve(filePath),
        scheme: 'file',
        toString: () => `file://${path.resolve(filePath)}`
      }),
      parse: (uriStr: string) => ({
        fsPath: uriStr,
        path: uriStr,
        scheme: 'file',
        toString: () => uriStr
      })
    },
    StatusBarAlignment: {
      Left: 1,
      Right: 2
    },
    CompletionItem: class CompletionItem {
      constructor(public label: unknown, public kind?: unknown) {}
    },
    CompletionItemKind: {
      Text: 0,
      Method: 1,
      Function: 2,
      Constructor: 3,
      Field: 4,
      Variable: 5,
      Class: 6,
      Interface: 7,
      Module: 8,
      Property: 9,
      Unit: 10,
      Value: 11,
      Enum: 12,
      Keyword: 13,
      Snippet: 14,
      Color: 15,
      File: 16,
      Reference: 17,
      Folder: 18,
      EnumMember: 19,
      Constant: 20,
      Struct: 21,
      Event: 22,
      Operator: 23,
      TypeParameter: 24
    },
    CodeAction: class CodeAction {
      constructor(public title: string, public kind?: unknown) {}
    },
    CodeActionKind: {
      Empty: '',
      QuickFix: 'quickfix',
      Refactor: 'refactor',
      RefactorExtract: 'refactor.extract',
      RefactorInline: 'refactor.inline',
      RefactorRewrite: 'refactor.rewrite',
      Source: 'source',
      SourceOrganizeImports: 'source.organizeImports',
      SourceFixAll: 'source.fixAll'
    },
    CodeLens: class CodeLens {
      constructor(public range: Range, public command?: unknown) {}
    },
    DocumentHighlight: class DocumentHighlight {
      constructor(public range: Range, public kind?: unknown) {}
    },
    DocumentHighlightKind: { Text: 0, Read: 1, Write: 2 },
    DocumentSymbol: class DocumentSymbol {
      constructor(
        public name: string,
        public detail: string,
        public kind: unknown,
        public range: Range,
        public selectionRange: Range
      ) {}
    },
    SymbolInformation: class SymbolInformation {
      constructor(
        public name: string,
        public kind: unknown,
        public containerNameOrRange: unknown,
        public uri?: unknown,
        public range?: Range
      ) {}
    },
    SymbolKind: {
      File: 0,
      Module: 1,
      Namespace: 2,
      Package: 3,
      Class: 4,
      Method: 5,
      Property: 6,
      Field: 7,
      Constructor: 8,
      Enum: 9,
      Interface: 10,
      Function: 11,
      Variable: 12,
      Constant: 13,
      String: 14,
      Number: 15,
      Boolean: 16,
      Array: 17,
      Object: 18,
      Key: 19,
      Null: 20,
      EnumMember: 21,
      Struct: 22,
      Event: 23,
      Operator: 24,
      TypeParameter: 25
    },
    Hover: class Hover {
      constructor(public contents: unknown, public range?: Range) {}
    },
    SignatureHelp: class SignatureHelp {
      public signatures: unknown[] = [];
      public activeSignature = 0;
      public activeParameter = 0;
    },
    SignatureInformation: class SignatureInformation {
      constructor(public label: string, public documentation?: unknown) {}
    },
    ParameterInformation: class ParameterInformation {
      constructor(public label: unknown, public documentation?: unknown) {}
    },
    TextEdit: class TextEdit {
      static replace(range: Range, newText: string) {
        return new TextEdit(range, newText);
      }
      static insert(position: Position, newText: string) {
        return new TextEdit(new Range(position, position), newText);
      }
      static delete(range: Range) {
        return new TextEdit(range, '');
      }
      constructor(public range: Range, public newText: string) {}
    },
    WorkspaceEdit: class WorkspaceEdit {
      set() {}
      replace() {}
      insert() {}
      delete() {}
    },
    ColorInformation: class ColorInformation {
      constructor(public range: Range, public color: unknown) {}
    },
    ColorPresentation: class ColorPresentation {
      constructor(public label: string) {}
    },
    FoldingRange: class FoldingRange {
      constructor(public start: number, public end: number, public kind?: unknown) {}
    },
    FoldingRangeKind: { Comment: 1, Imports: 2, Region: 3 },
    SelectionRange: class SelectionRange {
      constructor(public range: Range, public parent?: SelectionRange) {}
    },
    CallHierarchyItem: class CallHierarchyItem {
      constructor(
        public kind: unknown,
        public name: string,
        public detail: string,
        public uri: unknown,
        public range: Range,
        public selectionRange: Range
      ) {}
    },
    CallHierarchyIncomingCall: class CallHierarchyIncomingCall {
      constructor(public from: unknown, public fromRanges: Range[]) {}
    },
    CallHierarchyOutgoingCall: class CallHierarchyOutgoingCall {
      constructor(public to: unknown, public fromRanges: Range[]) {}
    },
    TypeHierarchyItem: class TypeHierarchyItem {
      constructor(
        public kind: unknown,
        public name: string,
        public detail: string,
        public uri: unknown,
        public range: Range,
        public selectionRange: Range
      ) {}
    },
    InlayHint: class InlayHint {
      constructor(public position: Position, public label: unknown, public kind?: unknown) {}
    },
    InlayHintKind: { Type: 1, Parameter: 2 },
    InlineValueText: class InlineValueText {
      constructor(public range: Range, public text: string) {}
    },
    InlineValueVariableLookup: class InlineValueVariableLookup {
      constructor(
        public range: Range,
        public variableName?: string,
        public caseSensitiveLookup?: boolean
      ) {}
    },
    InlineValueEvaluatableExpression: class InlineValueEvaluatableExpression {
      constructor(public range: Range, public expression?: string) {}
    },
    SemanticTokens: class SemanticTokens {
      constructor(public data: Uint32Array, public resultId?: string) {}
    },
    SemanticTokensBuilder: class SemanticTokensBuilder {
      push() {}
      build() {
        return new mockVscode.SemanticTokens(new Uint32Array());
      }
    },
    SemanticTokensLegend: class SemanticTokensLegend {
      constructor(public tokenTypes: string[], public tokenModifiers: string[] = []) {}
    },
    ProgressLocation: { SourceControl: 1, Window: 10, Notification: 15 },
    CancellationTokenSource: class CancellationTokenSource {
      token = {
        isCancellationRequested: false,
        onCancellationRequested: () => ({ dispose: () => {} })
      };
      cancel() {
        this.token.isCancellationRequested = true;
      }
      dispose() {}
    },
    DocumentLink: class DocumentLink {
      constructor(public range: Range, public target?: unknown) {}
    },
    SnippetString: class SnippetString {
      constructor(public value?: string) {}
      appendText(t: string) {
        this.value = (this.value || '') + t;
        return this;
      }
      appendTabstop() {
        return this;
      }
      appendPlaceholder() {
        return this;
      }
    },
    MarkdownString: class MarkdownString {
      constructor(public value?: string, public isTrusted?: boolean) {}
      appendText(t: string) {
        this.value = (this.value || '') + t;
        return this;
      }
      appendMarkdown(m: string) {
        this.value = (this.value || '') + m;
        return this;
      }
      appendCodeblock(c: string, l?: string) {
        this.value = (this.value || '') + '\n```' + (l || '') + '\n' + c + '\n```\n';
        return this;
      }
    },
    RelativePattern: class RelativePattern {
      constructor(public base: unknown, public pattern: string) {}
    },
    InlineCompletionItem: class InlineCompletionItem {
      constructor(public insertText: string, public range?: Range, public command?: unknown) {}
    },
    InlineCompletionList: class InlineCompletionList {
      constructor(public items: unknown[]) {}
    },
    InlineCompletionTriggerKind: { Invoke: 0, Automatic: 1 },
    FileChangeType: { Created: 1, Changed: 2, Deleted: 3 },
    TextDocumentSaveReason: { Manual: 1, AfterDelay: 2, FocusOut: 3 },
    env: {
      appName: 'Visual Studio Code',
      language: 'en',
      openExternal: async () => true
    },
    languages: {
      match: () => 1,
      createDiagnosticCollection: () => ({
        set: () => {},
        delete: () => {},
        clear: () => {},
        dispose: () => {}
      }),
      registerDefinitionProvider: () => ({ dispose: () => {} }),
      registerRenameProvider: () => ({ dispose: () => {} }),
      registerHoverProvider: () => ({ dispose: () => {} }),
      registerDocumentFormattingEditProvider: () => ({ dispose: () => {} }),
      registerDocumentRangeFormattingEditProvider: () => ({ dispose: () => {} }),
      registerOnTypeFormattingEditProvider: () => ({ dispose: () => {} }),
      registerCallHierarchyProvider: () => ({ dispose: () => {} }),
      registerTypeHierarchyProvider: () => ({ dispose: () => {} }),
      registerColorProvider: () => ({ dispose: () => {} }),
      registerCodeActionsProvider: () => ({ dispose: () => {} }),
      registerCodeLensProvider: () => ({ dispose: () => {} }),
      registerDocumentHighlightProvider: () => ({ dispose: () => {} }),
      registerDocumentSymbolProvider: () => ({ dispose: () => {} }),
      registerWorkspaceSymbolProvider: () => ({ dispose: () => {} }),
      registerImplementationProvider: () => ({ dispose: () => {} }),
      registerTypeDefinitionProvider: () => ({ dispose: () => {} }),
      registerDeclarationProvider: () => ({ dispose: () => {} }),
      registerFoldingRangeProvider: () => ({ dispose: () => {} }),
      registerSelectionRangeProvider: () => ({ dispose: () => {} }),
      registerDocumentLinkProvider: () => ({ dispose: () => {} }),
      registerInlayHintsProvider: () => ({ dispose: () => {} }),
      registerInlineValuesProvider: () => ({ dispose: () => {} }),
      registerDocumentSemanticTokensProvider: () => ({ dispose: () => {} }),
      registerDocumentRangeSemanticTokensProvider: () => ({ dispose: () => {} })
    },
    workspace: {
      textDocuments: [] as unknown[],
      onDidChangeConfiguration: () => ({ dispose: () => {} }),
      onDidOpenTextDocument: () => ({ dispose: () => {} }),
      onDidCloseTextDocument: () => ({ dispose: () => {} }),
      onDidChangeTextDocument: () => ({ dispose: () => {} }),
      onDidSaveTextDocument: () => ({ dispose: () => {} }),
      onWillSaveTextDocument: () => ({ dispose: () => {} }),
      onDidChangeWorkspaceFolders: () => ({ dispose: () => {} }),
      get workspaceFolders() {
        if (!mockState.workspaceRoot) {
          return undefined;
        }
        return [
          {
            uri: {
              fsPath: mockState.workspaceRoot,
              path: mockState.workspaceRoot,
              scheme: 'file',
              toString: () => `file://${mockState.workspaceRoot}`
            },
            name: 'terralings',
            index: 0
          }
        ];
      },
      getConfiguration: (section?: string) => ({
        get: <T>(key: string, defaultValue?: T): T | undefined => {
          const fullKey = section ? `${section}.${key}` : key;
          if (fullKey in mockState.config) {
            return mockState.config[fullKey] as T;
          }
          if (key in mockState.config) {
            return mockState.config[key] as T;
          }
          return defaultValue;
        },
        has: (key: string) => {
          const fullKey = section ? `${section}.${key}` : key;
          return fullKey in mockState.config || key in mockState.config;
        },
        inspect: () => undefined,
        update: async () => {}
      }),
      createFileSystemWatcher: () => ({
        onDidChange: () => ({ dispose: () => {} }),
        onDidCreate: () => ({ dispose: () => {} }),
        onDidDelete: () => ({ dispose: () => {} }),
        dispose: () => {}
      })
    },
    window: {
      get activeTextEditor() {
        if (!mockState.activeFileName) {
          return undefined;
        }
        return {
          document: {
            fileName: mockState.activeFileName,
            uri: {
              fsPath: mockState.activeFileName,
              path: mockState.activeFileName,
              scheme: 'file'
            }
          }
        };
      },
      terminals: [] as unknown[],
      withProgress: async (_opts: unknown, task: (progress: unknown, token: unknown) => Promise<unknown>) => {
        return task({ report: () => {} }, { isCancellationRequested: false });
      },
      showTextDocument: async () => undefined,
      createOutputChannel: (name: string) => ({
        name,
        append: () => {},
        appendLine: () => {},
        clear: () => {},
        show: () => {},
        hide: () => {},
        dispose: () => {}
      }),
      createTerminal: (options?: { name?: string }) => ({
        name: options?.name || 'terminal',
        show: () => {},
        hide: () => {},
        sendText: () => {},
        dispose: () => {}
      }),
      createStatusBarItem: () => ({
        alignment: 2,
        priority: 100,
        text: '',
        tooltip: '',
        command: '',
        show: () => {},
        hide: () => {},
        dispose: () => {}
      }),
      showInformationMessage: async () => undefined,
      showWarningMessage: async () => undefined,
      showErrorMessage: async () => undefined,
      showInputBox: async () => undefined,
      showQuickPick: async () => undefined,
      registerTreeDataProvider: () => ({ dispose: () => {} })
    },
    commands: {
      registerCommand: () => ({ dispose: () => {} }),
      executeCommand: async () => undefined
    }
  };

  Module._load = function (request: string, _parent: unknown, _isMain: boolean) {
    if (request === 'vscode') {
      return mockVscode;
    }
    return originalLoad.apply(this, arguments);
  };

  isMocked = true;
  return mockVscode;
}
