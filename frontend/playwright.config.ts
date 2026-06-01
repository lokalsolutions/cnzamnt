import path from 'node:path';
import { fileURLToPath } from 'node:url';
import { defineConfig, devices } from '@playwright/test';

const apiURL = process.env.E2E_API_URL ?? 'http://127.0.0.1:18090';
const appURL = process.env.E2E_APP_URL ?? 'http://127.0.0.1:15173';
const currentDir = path.dirname(fileURLToPath(import.meta.url));
const repoRoot = path.resolve(currentDir, '..');
const backendDir = path.join(repoRoot, 'backend');
const e2eDBPath = path.join(backendDir, 'data', `e2e-cnzamnt-${process.pid}.db`);

export default defineConfig({
  testDir: './e2e',
  timeout: 30_000,
  expect: {
    timeout: 10_000,
  },
  fullyParallel: false,
  reporter: [['list']],
  use: {
    baseURL: appURL,
    trace: 'retain-on-failure',
    ...devices['Pixel 5'],
    channel: process.env.PLAYWRIGHT_CHROME_CHANNEL ?? 'chrome',
  },
  webServer: [
    {
      command:
        `env GOCACHE=${path.join(repoRoot, 'backend', '.gocache')} CNZAMNT_ADDR=127.0.0.1:18090 CNZAMNT_DB_PATH=${e2eDBPath} go run ./cmd/server`,
      cwd: backendDir,
      url: `${apiURL}/health`,
      reuseExistingServer: false,
      timeout: 60_000,
    },
    {
      command: `env VITE_PROXY_API_TARGET=${apiURL} npm run dev -- --host 127.0.0.1 --port 15173 --strictPort`,
      url: appURL,
      reuseExistingServer: false,
      timeout: 60_000,
    },
  ],
});
