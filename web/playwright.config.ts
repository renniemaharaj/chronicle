import { existsSync } from 'node:fs';
import { resolve } from 'node:path';
import { defineConfig } from '@playwright/test';

export default defineConfig({
  testDir: './tests',
  fullyParallel: false,
  workers: 1,
  use: {
    baseURL: 'http://127.0.0.1:18080',
    viewport: { width: 1536, height: 1024 },
    launchOptions: { executablePath: existsSync('/usr/bin/chromium') ? '/usr/bin/chromium' : undefined },
  },
  webServer: {
    command: 'cd .. && go run ./cmd/chronicle',
    url: 'http://127.0.0.1:18080/api/health',
    env: {
      CHRONICLE_ADDR: '127.0.0.1:18080',
      CHRONICLE_WORKSPACE: resolve('test-results/workspace'),
    },
    reuseExistingServer: false,
    timeout: 30000,
  },
});
