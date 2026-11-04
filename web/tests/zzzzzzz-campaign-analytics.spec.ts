import { expect, test } from '@playwright/test';
import { campaignFixture } from './campaign-fixture';
import { campaignProgressFixture } from './campaign-progress-fixture';

test('saved production history drives real totals, take cards and episode distribution', async ({ page }) => {
  const campaign = await campaignFixture(page, false);
  await campaignProgressFixture(page, campaign);
  await page
    .getByRole('button', {
      name: 'Reload saved progress',
      exact: true,
    })
    .click();
  await page
    .getByRole('button', {
      name: 'Details for episode 1',
      exact: true,
    })
    .click();

  const table = page.getByRole('table', { name: 'Series episodes' });
  await expect(table.locator('tbody tr').first()).toContainText('+20%');
  await expect(table.locator('tbody tr').first()).toContainText('1.2K views');
  await expect(page.getByLabel('Series performance')).toContainText('1.2K');
  await expect(page.getByRole('complementary', { name: 'Episode details' })).toContainText('1,200');
  await page.getByLabel('Filter episodes by status').selectOption('published');
  await expect(table.locator('tbody tr')).toHaveCount(1);
  await page.getByLabel('Filter episodes by status').selectOption('all');
  await page
    .getByRole('tab', {
      name: 'Clips Library',
      exact: true,
    })
    .click();
  await expect(page.locator('.reel-clip-card')).toHaveCount(2);
  await expect(page.getByText('/synthetic/second-episode.mp4', { exact: true })).toBeVisible();
  await page
    .getByRole('tab', {
      name: 'Analytics',
      exact: true,
    })
    .click();
  await expect(page.getByLabel('Campaign analytics')).toContainText('1 of 3 episodes have observed views');
  await expect(page.locator('.reel-revenue > span')).toHaveText('USD');
  await expect(page.locator('.reel-revenue > strong')).toHaveText('2.00');
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
  await expect(
    page.getByRole('link', {
      name: 'View publication',
      exact: true,
    }),
  ).toHaveAttribute('href', 'https://www.facebook.com/reel/987654321');
  await page.getByRole('tab', { name: /Episodes/ }).click();
  await expect(page.locator('.toast')).toHaveCount(0);
  await page.screenshot({
    path: '../.chronicle/reel-series-observations.png',
    fullPage: true,
  });
});
