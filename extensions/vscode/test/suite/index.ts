import * as path from 'path';
import * as fs from 'fs';
import Mocha = require('mocha');
import { setupMockVscode } from './mockVscode';

export function run(): Promise<void> {
  // Ensure VS Code mock is initialized before running tests
  setupMockVscode();

  const mocha = new Mocha({
    ui: 'tdd',
    color: true,
    timeout: 10000
  });

  const testsRoot = path.resolve(__dirname);

  return new Promise((resolve, reject) => {
    try {
      const files = fs
        .readdirSync(testsRoot)
        .filter((file) => file.endsWith('.test.js') || (file.endsWith('.test.ts') && !file.endsWith('.d.ts')));

      for (const file of files) {
        mocha.addFile(path.join(testsRoot, file));
      }

      mocha.run((failures: number) => {
        if (failures > 0) {
          reject(new Error(`${failures} test(s) failed.`));
        } else {
          resolve();
        }
      });
    } catch (err) {
      reject(err);
    }
  });
}
