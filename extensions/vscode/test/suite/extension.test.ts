import * as assert from 'assert';
import * as path from 'path';
import * as fs from 'fs';
import { setupMockVscode, setMockConfig, resetMockState, mockState } from './mockVscode';

// Ensure mock is initialized before imports
setupMockVscode();

import { findTerralingsBinary } from '../../src/lspClient';
import { normalizeExerciseName } from '../../src/cliRunner';
import * as extension from '../../src/extension';

suite('Terralings Extension Test Suite', () => {
  setup(() => {
    resetMockState();
  });

  teardown(() => {
    resetMockState();
  });

  suite('Walkthrough Documentation', () => {
    const extensionRoot = path.resolve(__dirname, '../../../');
    const walkthroughDir = path.join(extensionRoot, 'walkthrough');

    const expectedFiles = [
      '01_welcome.md',
      '02_anatomy.md',
      '03_watch.md',
      '04_tui_hints.md',
      '05_editor_lsp.md'
    ];

    test('walkthrough directory exists', () => {
      assert.ok(fs.existsSync(walkthroughDir), `Walkthrough directory should exist at ${walkthroughDir}`);
    });

    test('all 5 walkthrough markdown files exist and are non-empty', () => {
      for (const filename of expectedFiles) {
        const filePath = path.join(walkthroughDir, filename);
        assert.ok(fs.existsSync(filePath), `Walkthrough file ${filename} should exist at ${filePath}`);

        const content = fs.readFileSync(filePath, 'utf-8');
        assert.ok(content.trim().length > 0, `Walkthrough file ${filename} should not be empty`);
      }
    });

    test('walkthrough contents contain expected key sections', () => {
      const welcome = fs.readFileSync(path.join(walkthroughDir, '01_welcome.md'), 'utf-8');
      assert.ok(welcome.toLowerCase().includes('welcome'), '01_welcome.md should contain welcome content');

      const anatomy = fs.readFileSync(path.join(walkthroughDir, '02_anatomy.md'), 'utf-8');
      assert.ok(anatomy.includes('TODO') || anatomy.toLowerCase().includes('anatomy'), '02_anatomy.md should explain exercise structure');

      const watch = fs.readFileSync(path.join(walkthroughDir, '03_watch.md'), 'utf-8');
      assert.ok(watch.toLowerCase().includes('watch'), '03_watch.md should describe watch mode');

      const tui = fs.readFileSync(path.join(walkthroughDir, '04_tui_hints.md'), 'utf-8');
      assert.ok(tui.toLowerCase().includes('tui') || tui.toLowerCase().includes('hint'), '04_tui_hints.md should describe TUI/hints');

      const lsp = fs.readFileSync(path.join(walkthroughDir, '05_editor_lsp.md'), 'utf-8');
      assert.ok(lsp.toLowerCase().includes('lsp') || lsp.toLowerCase().includes('diagnostic'), '05_editor_lsp.md should describe LSP');
    });
  });

  suite('findTerralingsBinary', () => {
    test('resolves default terralings binary when no configuration is set', () => {
      const binary = findTerralingsBinary();
      assert.ok(typeof binary === 'string' && binary.length > 0, 'Binary should resolve to a non-empty string');
    });

    test('respects custom binaryPath configuration setting', () => {
      const customPath = '/usr/local/custom/terralings';
      setMockConfig('terralings.binaryPath', customPath);
      const binary = findTerralingsBinary();
      assert.strictEqual(binary, customPath, 'Custom binaryPath setting should be returned');
    });

    test('resolves relative binary path within workspace root', () => {
      const binary = findTerralingsBinary('/tmp/test-workspace');
      assert.ok(typeof binary === 'string' && binary.length > 0);
    });
  });

  suite('normalizeExerciseName', () => {
    test('handles bare exercise name', () => {
      assert.strictEqual(normalizeExerciseName('primitives01'), 'primitives01');
      assert.strictEqual(normalizeExerciseName('syntax01'), 'syntax01');
      assert.strictEqual(normalizeExerciseName('state03'), 'state03');
    });

    test('handles bare name with leading/trailing whitespace', () => {
      assert.strictEqual(normalizeExerciseName('  primitives01  '), 'primitives01');
    });

    test('handles full and relative file paths with .tf extension', () => {
      assert.strictEqual(
        normalizeExerciseName('exercises/01_primitives/primitives01.tf'),
        'primitives01'
      );
      assert.strictEqual(
        normalizeExerciseName('/home/user/repo/exercises/03_collections/collections02.tf'),
        'collections02'
      );
      assert.strictEqual(normalizeExerciseName('variables01.tf'), 'variables01');
    });

    test('handles .hcl and .json extensions', () => {
      assert.strictEqual(
        normalizeExerciseName('exercises/04_functions/functions01.hcl'),
        'functions01'
      );
      assert.strictEqual(normalizeExerciseName('state.json'), 'state');
    });

    test('returns undefined for empty or invalid input without active editor', () => {
      assert.strictEqual(normalizeExerciseName(''), undefined);
      assert.strictEqual(normalizeExerciseName(undefined), undefined);
    });

    test('infers exercise name from active text editor when input is not provided', () => {
      mockState.activeFileName = '/repo/exercises/02_variables/variables03.tf';
      assert.strictEqual(normalizeExerciseName(), 'variables03');
    });
  });

  suite('Extension Module Exports', () => {
    test('exports activate and deactivate lifecycle functions', () => {
      assert.strictEqual(typeof extension.activate, 'function');
      assert.strictEqual(typeof extension.deactivate, 'function');
    });
  });
});
