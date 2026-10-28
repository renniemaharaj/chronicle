import { spawn } from 'node:child_process';

const detached = process.platform !== 'win32';
const backend = spawn('go', ['run', './cmd/chronicle'], {
  stdio: 'inherit',
  detached,
});

backend.on('error', (error) => {
  process.stderr.write(`${error.message}\n`);
  process.exitCode = 1;
});

backend.on('exit', (code) => {
  process.exitCode = code ?? 1;
});

for (const signal of ['SIGINT', 'SIGTERM']) {
  process.on(signal, () => {
    try {
      // go run starts a server process that must receive shutdown signals too.
      if (detached) {
        process.kill(-backend.pid, signal);
      } else {
        backend.kill(signal);
      }
    } catch {
      // The backend may have exited before the signal arrived.
    }
  });
}
