import { readFile } from 'node:fs/promises';
import { chromium } from '../web/node_modules/playwright-core/index.mjs';
import { backend } from './reel-backend.mjs';

export async function openMapBrowser(project) {
  const browser = await chromium.launch({
    executablePath: process.env.CHROME_RENDERER ?? '/usr/bin/chromium',
    headless: true,
  });
  const page = await browser.newPage({
    viewport: {
      width: 1080,
      height: 1920,
    },
    reducedMotion: 'no-preference',
  });
  await page.addInitScript((saved) => {
    localStorage.setItem(
      'chronicle-workspace',
      JSON.stringify({
        projectId: saved.id,
        projects: [saved],
      }),
    );
  }, project);
  await page.goto(`${backend}/?capture=story`);
  await page.locator('.status-indicator.online').waitFor();
  await page
    .getByRole('button', {
      name: 'Story',
      exact: true,
    })
    .click();

  const follow = page.getByRole('button', {
    name: 'Follow location',
    exact: true,
  });
  if ((await follow.getAttribute('aria-pressed')) !== 'true') {
    await follow.click();
  }

  await page
    .getByRole('button', {
      name: 'Reel and Shorts preview',
      exact: true,
    })
    .click();

  const style = await readFile(new URL('../web/src/lib/styles/reel-capture.css', import.meta.url), 'utf8');
  await page.addStyleTag({ content: style });
  await page.bringToFront();

  return { page, browser };
}
