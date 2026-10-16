import { spawnSync } from 'node:child_process';
import { readFileSync, writeFileSync } from 'node:fs';

const result = spawnSync('npm', ['ci', '--prefix', 'web', '--no-audit', '--no-fund'], { stdio: 'inherit' });
if (result.status !== 0) process.exit(result.status ?? 1);

// Generated locks retain their complete data in compact machine-readable form.
const path = new URL('../web/package-lock.json', import.meta.url);
writeFileSync(path, `${JSON.stringify(JSON.parse(readFileSync(path, 'utf8')))}\n`);
