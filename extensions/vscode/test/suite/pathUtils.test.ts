import * as assert from 'assert';
import * as path from 'path';
import * as fs from 'fs';
import * as os from 'os';
import { getEffectiveWorkspaceRoot, resolveExercisePath } from '../../src/pathUtils';

suite('PathUtils Test Suite', () => {
  const tempDirs: string[] = [];

  setup(() => {
    // create clean temp dir for test isolation
  });

  teardown(() => {
    for (const dir of tempDirs) {
      if (fs.existsSync(dir)) {
        fs.rmSync(dir, { recursive: true, force: true });
      }
    }
    tempDirs.length = 0;
  });

  function createTempDir(): string {
    const d = fs.mkdtempSync(path.join(os.tmpdir(), 'terralings-path-test-'));
    tempDirs.push(d);
    return d;
  }

  test('getEffectiveWorkspaceRoot returns a non-empty string and never returns root slash', () => {
    const root = getEffectiveWorkspaceRoot('terralings');
    assert.ok(root, 'workspace root should not be empty');
    assert.notStrictEqual(root, '/', 'workspace root should never be root slash');
  });

  test('resolveExercisePath resolves direct file in workspace root', () => {
    const root = createTempDir();
    const exDir = path.join(root, 'exercises', '01_primitives');
    fs.mkdirSync(exDir, { recursive: true });
    const exFile = path.join(exDir, 'primitives01.tf');
    fs.writeFileSync(exFile, 'terraform {}\n');

    const resolved = resolveExercisePath('exercises/01_primitives/primitives01.tf', root);
    assert.strictEqual(resolved, exFile);
  });

  test('resolveExercisePath resolves when workspace is inside exercises subfolder', () => {
    const root = createTempDir();
    const exercisesDir = path.join(root, 'exercises');
    const chapterDir = path.join(exercisesDir, '01_primitives');
    fs.mkdirSync(chapterDir, { recursive: true });
    const exFile = path.join(chapterDir, 'primitives01.tf');
    fs.writeFileSync(exFile, 'terraform {}\n');

    // Opened in root/exercises
    const resolvedFromExercises = resolveExercisePath('exercises/01_primitives/primitives01.tf', exercisesDir);
    assert.strictEqual(resolvedFromExercises, exFile);

    // Opened in root/exercises/01_primitives
    const resolvedFromChapter = resolveExercisePath('exercises/01_primitives/primitives01.tf', chapterDir);
    assert.strictEqual(resolvedFromChapter, exFile);
  });

  test('resolveExercisePath resolves directory to main.tf if present', () => {
    const root = createTempDir();
    const exDir = path.join(root, 'exercises', '06_modules', 'modules01');
    fs.mkdirSync(exDir, { recursive: true });
    const mainTf = path.join(exDir, 'main.tf');
    fs.writeFileSync(mainTf, 'module "foo" {}\n');

    const resolved = resolveExercisePath('exercises/06_modules/modules01', root);
    assert.strictEqual(resolved, mainTf);
  });

  test('resolveExercisePath handles existing absolute paths directly', () => {
    const root = createTempDir();
    const exFile = path.join(root, 'custom.tf');
    fs.writeFileSync(exFile, 'terraform {}\n');

    const resolved = resolveExercisePath(exFile, root);
    assert.strictEqual(resolved, exFile);
  });

  test('resolveExercisePath falls back safely when file is not yet initialized', () => {
    const root = createTempDir();
    const rel = 'exercises/99_future/future01.tf';
    const resolved = resolveExercisePath(rel, root);
    assert.strictEqual(resolved, path.resolve(root, rel));
  });
});
