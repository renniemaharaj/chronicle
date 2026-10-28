import { expect, test } from '@playwright/test';

test('starts empty with readable typography, the favicon mark, and both themes', async ({ page }) => {
  await page.goto('/');
  await expect(page.getByRole('heading', { name: 'Begin with a subject.' })).toBeVisible();
  await expect(page.locator('.brand-mark .chronicle-mark')).toBeVisible();
  await expect(page.locator('.brand-mark .chronicle-mark')).toHaveCSS('mask-image', /favicon\.svg/);
  await expect(page.locator('link[rel="icon"]')).toHaveAttribute('href', /\/favicon\.svg$/);
  await expect(page.locator('.svelte-flow__node')).toHaveCount(0);

  const sizes = await page.locator('.app-shell').evaluate((shell) => {
    const elements = [...shell.querySelectorAll('*')].filter((node) =>
      [...node.childNodes].some((child) => child.nodeType === Node.TEXT_NODE && child.textContent?.trim()),
    );

    return elements.map((node) => parseFloat(getComputedStyle(node).fontSize));
  });
  expect(Math.min(...sizes)).toBeGreaterThanOrEqual(13);
  expect(Math.max(...sizes)).toBeLessThanOrEqual(27);
  await page.getByLabel('Toggle theme').click();
  await expect(page.locator('.app-shell')).toHaveClass(/light/);
  await expect(page.locator('.app-shell')).toHaveCSS('color', 'rgb(24, 43, 67)');
  await page.reload();
  await expect(page.locator('.app-shell')).toHaveClass(/light/);
  await page
    .getByRole('button', {
      name: 'New Research',
      exact: true,
    })
    .click();
  await page.getByLabel('Research subject', { exact: true }).fill('A synthetic institution');
  await page.locator('.research-subject button').click();
  await expect(page.getByRole('alert')).toContainText('configure a model');
  await expect(page.locator('.run-metrics')).toHaveCount(0);
});
