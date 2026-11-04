import { expect, test } from '@playwright/test';
import { campaignFixture } from './campaign-fixture';
import type { Campaign } from '../src/lib/campaign/types';

test('series workspace connects episode selection, search, clips, analytics and distribution', async ({
  page,
}) => {
  const errors: string[] = [];
  page.on('pageerror', (error) => errors.push(error.message));

  const campaign: Campaign = await campaignFixture(page, false);
  const table = page.getByRole('table', { name: 'Series episodes' });
  await expect(page.getByRole('complementary', { name: 'Episode details' })).toBeVisible();
  await expect(
    page.getByRole('button', {
      name: 'Reel Series',
      exact: true,
    }),
  ).toBeVisible();
  await page.getByLabel('Search episodes').fill('no episode with this title');
  await expect(page.getByText('No matching episodes')).toBeVisible();
  await page.getByLabel('Search episodes').fill('');
  await expect(table.locator('tbody tr')).toHaveCount(campaign.parts.length);
  await page
    .getByRole('button', {
      name: 'Details for episode 2',
      exact: true,
    })
    .click();
  await expect(page.getByRole('complementary', { name: 'Episode details' })).toContainText(
    `Episode 2 / ${campaign.parts.length}`,
  );
  await page
    .getByRole('button', {
      name: 'Close episode details',
      exact: true,
    })
    .click();
  await expect(page.getByRole('complementary', { name: 'Episode details' })).toHaveCount(0);
  await page
    .getByRole('button', {
      name: 'Details for episode 1',
      exact: true,
    })
    .click();
  await page
    .getByRole('tab', {
      name: 'Clips Library',
      exact: true,
    })
    .click();
  await expect(page.getByText('A home for every take')).toBeVisible();
  await page
    .getByRole('tab', {
      name: 'Analytics',
      exact: true,
    })
    .click();
  await expect(page.getByText('Your first observations will appear here')).toBeVisible();
  await page
    .getByRole('tab', {
      name: 'Series',
      exact: true,
    })
    .click();
  await page
    .getByRole('tab', {
      name: 'Distribution',
      exact: true,
    })
    .click();
  await expect(page.getByRole('region', { name: 'Series distribution' })).toContainText('Next in the series');
  await page.getByRole('tab', { name: /Episodes/ }).click();
  await expect(page.locator('.toast')).toHaveCount(0);
  await page.screenshot({
    path: '../.chronicle/reel-series-desktop.png',
    fullPage: true,
  });
  await page
    .getByRole('button', {
      name: 'Toggle theme',
      exact: true,
    })
    .click();
  await page.screenshot({
    path: '../.chronicle/reel-series-light.png',
    fullPage: true,
  });
  await page.setViewportSize({
    width: 390,
    height: 844,
  });
  await page
    .getByRole('button', {
      name: 'Toggle explorer',
      exact: true,
    })
    .click();
  await expect(page.getByRole('navigation', { name: 'Research explorer' })).toBeVisible();
  await page
    .getByRole('button', {
      name: 'Collapse explorer',
      exact: true,
    })
    .click();
  expect(await page.evaluate(() => document.documentElement.scrollWidth)).toBeLessThanOrEqual(390);
  await page.screenshot({
    path: '../.chronicle/reel-series-mobile.png',
    fullPage: true,
  });
  expect(errors).toEqual([]);
});
