import { spawn } from 'node:child_process';

const detached = process.platform !== 'win32';
const children = [
  spawn('go', ['run', './cmd/chronicle'], { stdio: 'inherit', detached }),
  spawn('npm', ['--prefix', 'web', 'run', 'dev'], { stdio: 'inherit', detached })
];
let stopping = false;

function stop(code = 0) {
  if (stopping) return;
  stopping = true;
  for (const child of children) {
    try {
      if (detached) process.kill(-child.pid, 'SIGTERM');
      else child.kill('SIGTERM');
    } catch {
      // A process may already have exited before its sibling stopped.
    }
  }
  process.exitCode = code;
}

for (const child of children) {
  child.on('error', (error) => {
    process.stderr.write(`${error.message}\n`);
    stop(1);
  });
  child.on('exit', (code) => stop(code ?? 0));
}
process.on('SIGINT', () => stop());
process.on('SIGTERM', () => stop());
