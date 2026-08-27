import * as assert from 'assert';
import { setupMockVscode, resetMockState } from './mockVscode';

// Ensure mock is initialized before imports
setupMockVscode();

import {
  CURRICULUM,
  ALL_EXERCISES,
  findChapter,
  findExercise,
  resolveExerciseUri,
  TerralingsTreeDataProvider,
  TerralingsTreeItem
} from '../../src/treeProvider';

suite('Terralings TreeDataProvider & Curriculum Test Suite', () => {
  let provider: TerralingsTreeDataProvider;

  setup(() => {
    resetMockState();
    provider = new TerralingsTreeDataProvider();
  });

  teardown(() => {
    resetMockState();
    provider.dispose();
  });

  suite('Curriculum Structure & Integrity', () => {
    test('CURRICULUM contains exactly 13 chapters', () => {
      assert.strictEqual(CURRICULUM.length, 13, 'Curriculum must have exactly 13 chapters');
    });

    test('ALL_EXERCISES contains exactly 56 exercises', () => {
      assert.strictEqual(ALL_EXERCISES.length, 56, 'ALL_EXERCISES array must have exactly 56 exercises');

      const totalExercisesInChapters = CURRICULUM.reduce(
        (sum, ch) => sum + ch.exercises.length,
        0
      );
      assert.strictEqual(
        totalExercisesInChapters,
        56,
        'Sum of exercises in all chapters must equal 56'
      );
    });

    test('chapter progression is strictly ordered 1 to 13', () => {
      CURRICULUM.forEach((chapter, index) => {
        const expectedNumber = index + 1;
        assert.strictEqual(
          chapter.number,
          expectedNumber,
          `Chapter at index ${index} should have number ${expectedNumber}`
        );
        assert.ok(chapter.name.length > 0, `Chapter ${expectedNumber} must have a non-empty name`);
        assert.ok(chapter.title.length > 0, `Chapter ${expectedNumber} must have a non-empty title`);
        assert.ok(
          chapter.description.length > 0,
          `Chapter ${expectedNumber} must have a non-empty description`
        );
        assert.ok(
          chapter.exercises.length > 0,
          `Chapter ${expectedNumber} must have at least 1 exercise`
        );
      });
    });

    test('exercise definitions contain all required fields and valid paths', () => {
      const exerciseNames = new Set<string>();

      ALL_EXERCISES.forEach((exercise) => {
        assert.ok(exercise.name.length > 0, 'Exercise name must not be empty');
        assert.ok(exercise.title.length > 0, 'Exercise title must not be empty');
        assert.ok(exercise.path.startsWith('exercises/'), `Path for ${exercise.name} must start with exercises/`);
        assert.ok(exercise.chapterName.length > 0, `Chapter name for ${exercise.name} must not be empty`);

        assert.ok(
          !exerciseNames.has(exercise.name),
          `Exercise name ${exercise.name} must be unique across the curriculum`
        );
        exerciseNames.add(exercise.name);
      });
    });
  });

  suite('Curriculum Lookup Helpers', () => {
    test('findChapter retrieves existing chapters and returns undefined for unknown', () => {
      const firstChapter = findChapter('01_primitives');
      assert.ok(firstChapter, 'Should find 01_primitives');
      assert.strictEqual(firstChapter?.number, 1);
      assert.strictEqual(firstChapter?.title, 'HCL Foundations & Core Primitives');

      const lastChapter = findChapter('13_governance');
      assert.ok(lastChapter, 'Should find 13_governance');
      assert.strictEqual(lastChapter?.number, 13);

      const unknownChapter = findChapter('99_nonexistent');
      assert.strictEqual(unknownChapter, undefined);
    });

    test('findExercise retrieves existing exercises and returns undefined for unknown', () => {
      const firstExercise = findExercise('primitives01');
      assert.ok(firstExercise, 'Should find primitives01');
      assert.strictEqual(firstExercise?.name, 'primitives01');
      assert.strictEqual(firstExercise?.chapterName, '01_primitives');

      const lastExercise = findExercise('gov03');
      assert.ok(lastExercise, 'Should find gov03');
      assert.strictEqual(lastExercise?.name, 'gov03');
      assert.strictEqual(lastExercise?.chapterName, '13_governance');

      const unknownExercise = findExercise('nonexistent_exercise');
      assert.strictEqual(unknownExercise, undefined);
    });

    test('resolveExerciseUri resolves file URI', () => {
      const uri = resolveExerciseUri('exercises/01_primitives/primitives01.tf');
      assert.ok(uri, 'Should return a valid URI');
      assert.ok(uri.fsPath.includes('primitives01.tf'));
    });
  });

  suite('TerralingsTreeDataProvider', () => {
    test('getChildren() with no element returns all 13 chapter tree items', () => {
      const rootChildren = provider.getChildren() as TerralingsTreeItem[];
      assert.ok(Array.isArray(rootChildren));
      assert.strictEqual(rootChildren.length, 13, 'Root must return 13 chapters');

      rootChildren.forEach((item, idx) => {
        assert.strictEqual(item.itemType, 'chapter');
        assert.ok(item.chapter);
        assert.strictEqual(item.chapter?.number, idx + 1);
        assert.strictEqual(item.collapsibleState, 1, 'Chapters should be Collapsed');
        assert.strictEqual(item.contextValue, 'chapter');
        assert.strictEqual(item.id, `chapter-${item.chapter?.name}`);
      });
    });

    test('getChildren(chapter) returns all exercises belonging to that chapter', () => {
      const rootChildren = provider.getChildren() as TerralingsTreeItem[];
      const firstChapterItem = rootChildren[0];
      assert.ok(firstChapterItem.chapter);

      const chapterExercises = provider.getChildren(firstChapterItem) as TerralingsTreeItem[];
      assert.ok(Array.isArray(chapterExercises));
      assert.strictEqual(
        chapterExercises.length,
        firstChapterItem.chapter.exercises.length,
        'Should return exact number of exercises for chapter 1'
      );

      chapterExercises.forEach((item, idx) => {
        const expectedExercise = firstChapterItem.chapter!.exercises[idx];
        assert.strictEqual(item.itemType, 'exercise');
        assert.strictEqual(item.exercise?.name, expectedExercise.name);
        assert.strictEqual(item.collapsibleState, 0, 'Exercises should be None collapsible');
        assert.strictEqual(item.contextValue, 'exercise');
        assert.strictEqual(item.id, `exercise-${expectedExercise.name}`);
        assert.ok(item.command, 'Exercise item should have open command attached');
      });
    });

    test('getChildren(exercise) returns empty array', () => {
      const rootChildren = provider.getChildren() as TerralingsTreeItem[];
      const chapterExercises = provider.getChildren(rootChildren[0]) as TerralingsTreeItem[];
      const leafChildren = provider.getChildren(chapterExercises[0]) as TerralingsTreeItem[];
      assert.deepStrictEqual(leafChildren, []);
    });

    test('getTreeItem returns the passed element', () => {
      const rootChildren = provider.getChildren() as TerralingsTreeItem[];
      const firstItem = rootChildren[0];
      assert.strictEqual(provider.getTreeItem(firstItem), firstItem);
    });

    test('getProgress returns accurate counts and percentage', () => {
      const progress = provider.getProgress();
      assert.strictEqual(progress.total, 56);
      assert.strictEqual(typeof progress.completed, 'number');
      assert.strictEqual(typeof progress.percentage, 'number');
      assert.ok(progress.percentage >= 0 && progress.percentage <= 100);
    });
  });

  suite('TerralingsTreeItem Status Badge Rendering', () => {
    const sampleChapter = CURRICULUM[0];
    const sampleExercise = sampleChapter.exercises[0];

    test('renders chapter item correctly with completed count', () => {
      const chapterItem = new TerralingsTreeItem(
        'chapter',
        sampleChapter,
        undefined,
        undefined,
        2,
        5
      );
      assert.strictEqual(chapterItem.contextValue, 'chapter');
      assert.strictEqual(chapterItem.description, '2/5');
      assert.strictEqual(chapterItem.id, `chapter-${sampleChapter.name}`);
    });

    test('renders exercise item with status: passed', () => {
      const item = new TerralingsTreeItem('exercise', sampleChapter, sampleExercise, 'passed');
      assert.strictEqual(item.description, 'Passed');
      const icon = item.iconPath as { id: string; color?: { id: string } };
      assert.strictEqual(icon.id, 'pass-filled');
    });

    test('renders exercise item with status: failed', () => {
      const item = new TerralingsTreeItem('exercise', sampleChapter, sampleExercise, 'failed');
      assert.strictEqual(item.description, 'Failed');
      const icon = item.iconPath as { id: string; color?: { id: string } };
      assert.strictEqual(icon.id, 'error');
    });

    test('renders exercise item with status: in_progress', () => {
      const item = new TerralingsTreeItem('exercise', sampleChapter, sampleExercise, 'in_progress');
      assert.strictEqual(item.description, 'In Progress');
      const icon = item.iconPath as { id: string; color?: { id: string } };
      assert.strictEqual(icon.id, 'play');
    });

    test('renders exercise item with status: not_started', () => {
      const item = new TerralingsTreeItem('exercise', sampleChapter, sampleExercise, 'not_started');
      assert.strictEqual(item.description, '');
      const icon = item.iconPath as { id: string };
      assert.strictEqual(icon.id, 'circle-outline');
    });
  });
});
