import { expect, test } from '@playwright/test';
import { readFile } from 'node:fs/promises';

async function sample(name: string) {
  const project = JSON.parse(await readFile(`../internal/fixtures/${name}.json`, 'utf8'));
  project.fixture = false;

  return project;
}

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
    headers: {
      Accept: 'application/json, text/event-stream',
      'MCP-Protocol-Version': '2025-11-25',
    },
    data: {
      jsonrpc: '2.0',
      id: 1,
      method: 'tools/call',
      params: {
        name: 'chronicle_save_project',
        arguments: {
          project,
          rebuildNarrative: true,
        },
      },
    },
  });
  expect(response.ok()).toBe(true);
  expect((await response.json()).result.isError).not.toBe(true);
  await expect(page.locator('.subject-info h1')).toHaveText('Alex Meridian');
  await expect(page.locator('.svelte-flow__node').first()).toBeVisible();

  const floor = await page.locator('.svelte-flow__viewport').evaluate((viewport) => {
    const zoom = new DOMMatrix(getComputedStyle(viewport).transform).a;
    const labels = [
      ...viewport.querySelectorAll('.temporal-node :is(strong, small, .node-kind), .svelte-flow__edge-label'),
    ];

    return Math.min(...labels.map((label) => parseFloat(getComputedStyle(label).fontSize) * zoom));
  });
  expect(floor).toBeGreaterThanOrEqual(13);
  await expect(page.locator('.temporal-node').first()).toHaveCSS('color', 'rgb(231, 237, 248)');

  const count = await page.locator('.svelte-flow__node').count();
  await page.getByLabel('Temporal resolution').selectOption('Century');
  await expect.poll(() => page.locator('.svelte-flow__node').count()).toBeLessThan(count);
  await page.getByLabel('Temporal resolution').selectOption('Year');
  await page
    .getByRole('combobox', {
      name: 'Graph view',
      exact: true,
    })
    .selectOption('Parallel');
  await expect(page.locator('.conflict-lane-banner')).toContainText('Parallel');
  await page
    .getByRole('tab', {
      name: 'Conflicts (1)',
      exact: true,
    })
    .click();
  await expect(page.locator('.conflict-card')).toContainText('1880 and 1830');
  await page.getByRole('button', { name: 'Explore this interpretation' }).first().click();
  await expect(page.getByRole('button', { name: 'Active interpretation' })).toBeVisible();
  await page
    .getByRole('combobox', {
      name: 'Graph view',
      exact: true,
    })
    .selectOption('Shadows');
  await page
    .getByRole('tab', {
      name: 'Details',
      exact: true,
    })
    .click();
  await page.getByLabel('Research note', { exact: true }).fill('Synthetic browser test annotation');
  await page.getByLabel('Save research note').click();
  await expect(page.locator('.annotation-box')).toContainText('Synthetic browser test annotation');
  await page
    .getByRole('button', {
      name: 'Save project',
      exact: true,
    })
    .click();
  await expect(page.locator('.toast')).toContainText('saved');
  await page.getByLabel('Toggle theme').click();
  await expect(page.locator('.phase-ruler')).toHaveCSS('background-color', 'rgb(249, 251, 254)');

  const card = page.locator('.temporal-node').first();
  await expect(card).toHaveCSS('color', 'rgb(24, 43, 67)');
  await page.screenshot({
    path: 'test-results/workspace-light.png',
    fullPage: true,
  });
  await page
    .getByRole('button', {
      name: 'Exports',
      exact: true,
    })
    .click();

  const downloadPromise = page.waitForEvent('download');
  await page.getByRole('button', { name: 'Project JSON Portable format v1' }).click();

  const download = await downloadPromise;
  const exported = JSON.parse(await readFile((await download.path())!, 'utf8'));
  expect(exported.formatVersion).toBe(1);
  expect(exported.evidence.length).toBeGreaterThan(0);
  expect(exported.conflicts.length).toBeGreaterThan(0);
  await page.getByLabel('Close dialog').click();
  await page
    .getByRole('button', {
      name: 'Library',
      exact: true,
    })
    .click();
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
