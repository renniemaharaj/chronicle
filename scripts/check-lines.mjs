import { readdir, readFile } from 'node:fs/promises';
import { dirname, join, relative } from 'node:path';
import { fileURLToPath } from 'node:url';

const root = dirname(dirname(fileURLToPath(import.meta.url)));
const ignored = new Set(['node_modules', '.git', '.svelte-kit', 'build', '.chronicle', 'test-results', 'playwright-report']);
const failures = [];
let count = 0;

async function visit(directory) {
  for (const entry of await readdir(directory, { withFileTypes: true })) {
    if (ignored.has(entry.name)) continue;
    // Dependency lockfiles are generated, not authored source.
    if (entry.name === 'package-lock.json') continue;
    const path = join(directory, entry.name);
    if (entry.isDirectory()) {
      await visit(path);
      continue;
    }
    const content = await readFile(path);
    if (content.includes(0)) continue;
    count++;
    const lines = content.toString().split('\n').length - (content.at(-1) === 10 ? 1 : 0);
    if (lines >= 150) failures.push(`${relative(root, path)}: ${lines} lines`);
  }
}

await visit(root);
if (failures.length) {
  process.stderr.write(`Files must have fewer than 150 lines:\n${failures.join('\n')}\n`);
  process.exitCode = 1;
} else {
  process.stdout.write(`All ${count} project text files have fewer than 150 lines.\n`);
}
