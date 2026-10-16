import { expect, test } from '@playwright/test';

test('chapter strip stays compact, navigation starts at its first node, and related search is local', async ({
  page,
}) => {
  await page.goto('/');
  await expect(page.locator('.subject-info h1')).toBeVisible();
  await expect(page.getByRole('button', { name: 'Story', exact: true })).toBeVisible();
  const strip = await page.locator('.phase-ruler').boundingBox();
  expect(strip?.height).toBeLessThan(65);
  await expect(page.locator('.ruler-ticks')).toHaveCount(0);
  await expect(page.locator('.phase span').first()).toHaveCSS('white-space', 'nowrap');
  await expect(page.locator('.phase span').first()).toHaveCSS('text-overflow', 'ellipsis');
  await page.locator('.phase').first().click();
  await expect
    .poll(async () => {
      const node = await page.locator('.svelte-flow__node.selected').first().boundingBox();
      const canvas = await page.locator('.graph-canvas').boundingBox();
      if (!node || !canvas) return -1;
      return Math.round(node.x - canvas.x);
    })
    .toBe(165);
  const node = await page.locator('.svelte-flow__node.selected').first().boundingBox();
  const canvas = await page.locator('.graph-canvas').boundingBox();
  expect(Math.round(node!.y - canvas!.y)).toBe(35);
  await page.getByRole('button', { name: 'Explore related subjects', exact: true }).click();
  await page.getByLabel('Search related subjects').fill('Meridian');
  await expect(page.locator('.subject-result')).toHaveCount(1);
  await page.getByLabel('Search related subjects').fill('No such subject');
  await expect(page.locator('.subject-result')).toHaveCount(0);
  await expect(page.getByText('No subjects match your search.')).toBeVisible();
  await page.getByLabel('Close dialog').click();
  await page.getByRole('button', { name: 'Contemporaries', exact: false }).first().click();
  await expect(page.getByLabel('Search related subjects')).toBeVisible();
  await expect(page.getByRole('dialog')).toHaveAccessibleName('Contemporaries');
  await page.getByLabel('Close dialog').click();
  await page.getByRole('button', { name: 'Library', exact: true }).click();
  await page.getByLabel('Search workspace projects').fill('No matching project');
  await expect(page.locator('.library-project')).toHaveCount(0);
});
