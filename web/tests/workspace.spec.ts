import { expect, test } from '@playwright/test';
import { readFile } from 'node:fs/promises';

async function sample(name: string) {
  const project = JSON.parse(await readFile(`../internal/fixtures/${name}.json`, 'utf8'));
  project.fixture = false;
  return project;
}

test('starts empty with readable typography, the favicon mark, and both themes', async ({ page }) => {
  await page.goto('/');
  await expect(page.getByRole('heading', { name: 'Begin with a subject.' })).toBeVisible();
  await expect(page.locator('.brand-mark img')).toHaveAttribute('src', '/favicon.svg');
  await expect(page.locator('.svelte-flow__node')).toHaveCount(0);
  const sizes = await page.locator('.app-shell').evaluate((shell) => {
    const elements = [...shell.querySelectorAll('*')].filter((node) =>
      [...node.childNodes].some((child) => child.nodeType === Node.TEXT_NODE && child.textContent?.trim()),
    );
    return elements.map((node) => parseFloat(getComputedStyle(node).fontSize));
  });
  expect(Math.min(...sizes)).toBeGreaterThanOrEqual(10);
  expect(Math.max(...sizes)).toBeLessThanOrEqual(27);
  await page.getByLabel('Toggle theme').click();
  await expect(page.locator('.app-shell')).toHaveClass(/light/);
  await expect(page.locator('.app-shell')).toHaveCSS('color', 'rgb(24, 43, 67)');
  await page.reload();
  await expect(page.locator('.app-shell')).toHaveClass(/light/);
  await page.getByRole('button', { name: 'New Research', exact: true }).click();
  await page.getByLabel('Research subject', { exact: true }).fill('A synthetic institution');
  await page.locator('.research-subject button').click();
  await expect(page.getByRole('alert')).toContainText('configure a model');
  await expect(page.locator('.run-metrics')).toHaveCount(0);
});

test('assistant tools publish a project live for graph exploration and evidence inspection', async ({
  page,
  request,
}) => {
  const errors: string[] = [];
  page.on('pageerror', (error) => errors.push(error.message));
  await page.goto('/');
  await expect(page.locator('.empty-workspace')).toBeVisible();
  const project = await sample('person');
  const response = await request.post('/api/mcp', {
    headers: { Accept: 'application/json, text/event-stream', 'MCP-Protocol-Version': '2025-11-25' },
    data: {
      jsonrpc: '2.0',
      id: 1,
      method: 'tools/call',
      params: { name: 'chronicle_save_project', arguments: { project, rebuildNarrative: true } },
    },
  });
  expect(response.ok()).toBe(true);
  expect((await response.json()).result.isError).not.toBe(true);
  await expect(page.locator('.subject-info h1')).toHaveText('Alex Meridian');
  await expect(page.locator('.svelte-flow__node').first()).toBeVisible();
  const count = await page.locator('.svelte-flow__node').count();
  await page.getByLabel('Temporal resolution').selectOption('Century');
  await expect.poll(() => page.locator('.svelte-flow__node').count()).toBeLessThan(count);
  await page.getByLabel('Temporal resolution').selectOption('Year');
  await page.getByLabel('Graph view').selectOption('Parallel');
  await expect(page.locator('.conflict-lane-banner')).toContainText('Parallel');
  await page.getByRole('button', { name: 'Conflicts (1)', exact: true }).click();
  await expect(page.locator('.conflict-card')).toContainText('1880 and 1830');
  await page.getByRole('button', { name: 'Explore this interpretation' }).first().click();
  await expect(page.getByRole('button', { name: 'Active interpretation' })).toBeVisible();
  await page.getByLabel('Graph view').selectOption('Shadows');
  await page.getByRole('button', { name: 'Details', exact: true }).click();
  await page.getByLabel('Research note', { exact: true }).fill('Synthetic browser test annotation');
  await page.getByLabel('Save research note').click();
  await expect(page.locator('.annotation-box')).toContainText('Synthetic browser test annotation');
  await page.getByRole('button', { name: 'Save project', exact: true }).click();
  await expect(page.getByRole('status')).toContainText('saved');
  await page.getByLabel('Toggle theme').click();
  await expect(page.locator('.phase-ruler')).toHaveCSS('background-color', 'rgb(249, 251, 254)');
  const card = page.locator('.temporal-node').first();
  await expect(card).toHaveCSS('color', 'rgb(24, 43, 67)');
  await page.screenshot({ path: 'test-results/workspace-light.png', fullPage: true });
  await page.getByRole('button', { name: 'Exports', exact: true }).click();
  const downloadPromise = page.waitForEvent('download');
  await page.getByRole('button', { name: 'Project JSON Portable format v1' }).click();
  const download = await downloadPromise;
  const exported = JSON.parse(await readFile((await download.path())!, 'utf8'));
  expect(exported.formatVersion).toBe(1);
  expect(exported.evidence.length).toBeGreaterThan(0);
  expect(exported.conflicts.length).toBeGreaterThan(0);
  await page.getByLabel('Close dialog').click();
  await page.getByRole('button', { name: 'Library', exact: true }).click();
  await page.locator('input[type=file]').setInputFiles({
    name: 'family.json',
    mimeType: 'application/json',
    buffer: Buffer.from(JSON.stringify(await sample('family'))),
  });
  await expect(page.locator('.subject-info h1')).toHaveText('Meridian family');
  await page.getByRole('button', { name: 'Family & Relationships' }).click();
  await expect(page.locator('.svelte-flow__node')).toHaveCount(2);
  await page.getByRole('button', { name: 'Conflicts & Disagreements' }).click();
  await expect(page.locator('.conflict-card')).toContainText('1880 and 1830');
  expect(errors).toEqual([]);
});
