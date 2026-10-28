import { readFile, readdir, writeFile } from 'node:fs/promises';
import { join } from 'node:path';
import { createRequire } from 'node:module';
import { javascriptLayout } from './readable-layout.mjs';
import { svelteLayout } from './readable-svelte.mjs';

const require = createRequire(new URL('../web/package.json', import.meta.url));
const typescript = require('typescript');
const svelte = require('svelte/compiler');
const check = process.argv.includes('--check');
const ignored = new Set([
  'node_modules',
  '.git',
  '.chronicle',
  '.svelte-kit',
  'build',
  'test-results',
  'playwright-report',
]);

async function visit(directory) {
  for (const entry of await readdir(directory, { withFileTypes: true })) {
    if (ignored.has(entry.name)) {
      continue;
    }

    const path = join(directory, entry.name);
    if (entry.isDirectory()) {
      await visit(path);
      continue;
    }

    if (!/\.(?:svelte|ts|mjs|js)$/.test(path)) {
      continue;
    }

    const source = await readFile(path, 'utf8');
    const result = path.endsWith('.svelte')
      ? svelteLayout(svelte, source).replace(
          /(<script\b[^>]*>)([\s\S]*?)(<\/script>)/g,
          (_, open, script, close) => {
            return open + javascriptLayout(typescript, script) + close;
          },
        )
      : javascriptLayout(typescript, source);

    if (source === result) {
      continue;
    }

    if (check) {
      process.stderr.write(`${path}: run npm run format\n`);
      process.exitCode = 1;
    } else {
      await writeFile(path, result);
    }
  }
}

await visit('.');
