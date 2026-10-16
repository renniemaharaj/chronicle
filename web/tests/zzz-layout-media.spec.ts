import { expect, test } from '@playwright/test';
import { readFile } from 'node:fs/promises';

test('shared controls, resizable Story, Library portraits and event media work together', async ({ page }) => {
  const project = JSON.parse(await readFile('../internal/fixtures/person.json', 'utf8'));
  project.fixture = false;
  project.id = 'layout-media-test';
  project.title = 'Synthetic layout and media project';
  project.media = [{
    id: 'linked-video', entityId: project.subjectId, sourceId: project.sources[0].id,
    kind: 'video', url: 'https://www.youtube.com/watch?v=abcdefghijk', title: 'Synthetic linked performance',
    caption: 'Synthetic test media', creator: 'Synthetic creator', license: 'Unknown',
    claimIds: project.events[0].claimIds, retrieved: true, useAsAvatar: false,
  }, {
    id: 'unlinked-video', entityId: project.subjectId, sourceId: project.sources[0].id,
    kind: 'video', url: 'https://youtu.be/lmnopqrstuv', title: 'Synthetic unlinked performance',
    caption: '', creator: '', license: 'Unknown', claimIds: [], retrieved: true, useAsAvatar: false,
  }];
  await page.route('**/api/projects', (route) => route.fulfill({ json: [project] }));
  await page.route('**/api/workspace/events', (route) => route.fulfill({
    contentType: 'text/event-stream', body: `data: ${JSON.stringify([project])}\n\n`,
  }));
  await page.route('https://i.ytimg.com/**', (route) => route.fulfill({
    contentType: 'image/svg+xml',
    body: '<svg xmlns="http://www.w3.org/2000/svg" width="480" height="360"><rect width="480" height="360" fill="#285588"/></svg>',
  }));
  await page.goto('/');
  await expect(page.locator('.subject-info h1')).toBeVisible();
  const nav = page.getByRole('button', { name: 'Social Networks & Website', exact: true });
  await expect(nav.locator('span')).toHaveCSS('white-space', 'nowrap');
  await expect(nav.locator('span')).toHaveCSS('text-overflow', 'ellipsis');
  const overflow = await page.locator('.left-panel').evaluate((panel) => panel.scrollWidth - panel.clientWidth);
  expect(overflow).toBeLessThanOrEqual(1);

  const story = page.locator('.story-panel');
  const before = await story.boundingBox();
  const handle = page.getByRole('button', { name: 'Resize story panel' });
  const bounds = (await handle.boundingBox())!;
  await page.mouse.move(bounds.x + bounds.width / 2, bounds.y + bounds.height / 2);
  await page.mouse.down();
  await page.mouse.move(bounds.x + bounds.width / 2, bounds.y - 120);
  await page.mouse.up();
  expect((await story.boundingBox())!.height).toBeGreaterThan(before!.height + 80);
  await page.reload();
  expect((await story.boundingBox())!.height).toBeGreaterThan(before!.height + 80);
  await handle.focus();
  await page.keyboard.press('ArrowDown');
  expect((await story.boundingBox())!.height).toBeLessThan(before!.height + 130);
  await expect(page.locator('.story-media')).toContainText('Synthetic linked performance');
  await expect(page.locator('.story-media')).not.toContainText('Synthetic unlinked performance');
  await expect(page.locator('.story-media img')).toHaveAttribute('src', /i.ytimg.com/);

  await page.getByLabel('Date range controls').click();
  await page.getByLabel('From year', { exact: true }).fill('2000');
  await page.getByLabel('Through year', { exact: true }).fill('1900');
  await page.getByRole('button', { name: 'Apply date range' }).click();
  await expect(page.getByRole('alert')).toContainText('starting year');
  await page.getByLabel('Close date controls').click();

  await page.getByRole('navigation', { name: 'Research explorer' }).getByRole('button', { name: 'Videos', exact: true }).click();
  await expect(page.locator('.media-card')).toHaveCount(2);
  await expect(page.locator('.media-timeline-link')).toHaveCount(1);
  await expect(page.locator('.subject-search-meta > label')).toHaveCSS('align-items', 'center');
  await page.getByLabel('Close dialog').click();
  await page.getByRole('button', { name: 'Library', exact: true }).click();
  const avatar = await page.locator('.library-project > .subject-avatar').first().boundingBox();
  expect(avatar!.width).toBe(72);
  expect(avatar!.height).toBe(72);
  const search = await page.locator('.subject-search-toolbar').boundingBox();
  const card = await page.locator('.library-project').first().boundingBox();
  expect(card!.y - (search!.y + search!.height)).toBeGreaterThanOrEqual(16);
  await expect(page.getByRole('button', { name: 'New research', exact: true })).toHaveCSS('background-color', 'rgb(40, 101, 209)');
  await page.getByLabel('Close dialog').click();
  await page.getByLabel('Settings', { exact: true }).click();
  await expect(page.locator('.settings-section')).toHaveCount(3);
  await expect(page.getByLabel('Provider', { exact: true })).toHaveCSS('padding-right', '36px');
  await page.screenshot({ path: 'test-results/settings-redesign.png', animations: 'disabled' });
});
