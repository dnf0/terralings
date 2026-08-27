import { run } from './suite';

async function main(): Promise<void> {
  try {
    console.log('Running Terralings VS Code Extension Mocha Test Suite...\n');
    await run();
    console.log('\nAll VS Code extension tests passed successfully.');
  } catch (err) {
    console.error('\nExtension test suite execution failed:', err);
    process.exit(1);
  }
}

main();
