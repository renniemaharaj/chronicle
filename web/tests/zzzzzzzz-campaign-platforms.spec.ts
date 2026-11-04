import { expect, test } from '@playwright/test';
import { campaignFixture } from './campaign-fixture';
import { campaignProgressFixture } from './campaign-progress-fixture';
import type { Campaign } from '../src/lib/campaign/types';

test('Facebook and YouTube keep independent receipts, ordering and observation scopes', async ({ page }) => {
  const original: Campaign = await campaignFixture(page, false);
  const created = await page.request.post('/api/campaigns', {
    data: {
      projectId: original.projectId,
      seconds: 15,
      platforms: ['facebook', 'youtube'],
    },
  });
  expect(created.ok()).toBe(true);

  const campaign: Campaign = await created.json();
  await campaignProgressFixture(page, campaign);
  await page
    .getByRole('button', {
      name: 'Reload saved progress',
      exact: true,
    })
    .click();
  await page.getByLabel('Saved series').selectOption(campaign.id);

  const table = page.getByRole('table', { name: 'Series episodes' });
  await expect(table.locator('tbody tr').first()).toContainText('Partially published');
  await expect(
    page.getByRole('link', {
      name: 'Open episode 1 on Facebook',
      exact: true,
    }),
  ).toBeVisible();
  await expect(
    page.getByRole('link', {
      name: 'Open episode 1 on YouTube',
      exact: true,
    }),
  ).toHaveCount(0);
  await page
    .getByRole('tab', {
      name: 'Distribution',
      exact: true,
    })
    .click();

  const youtube = page.locator('.reel-distribution').filter({
    has: page.getByRole('heading', {
      name: 'YouTube Shorts',
      exact: true,
    }),
  });
  await youtube
    .getByRole('button', {
      name: 'Manage publication',
      exact: true,
    })
    .click();
  await expect(page.getByLabel('Publication platform')).toHaveValue('youtube');
  await page.getByLabel('Published reel URL').fill('https://www.youtube.com/shorts/abcdefghijk');
  await page
    .getByRole('button', {
      name: 'Save publication receipt',
      exact: true,
    })
    .click();
  await page.getByLabel('Views', { exact: true }).fill('42');
  await page
    .getByRole('button', {
      name: 'Save observed performance',
      exact: true,
    })
    .click();
  await expect(page.getByText(/Views 42/)).toBeVisible();
  await page
    .getByRole('button', {
      name: 'Back to episodes',
      exact: true,
    })
    .click();
  await expect(table.locator('tbody tr').first()).toContainText('Published');
  await expect(page.getByLabel('Series performance')).toContainText('1.2K');

  const saved: Campaign = await (await page.request.get(`/api/campaigns/${campaign.id}`)).json();
  expect(saved.parts[0].publications?.facebook?.url).toContain('/reel/987654321');
  expect(saved.parts[0].publications?.youtube?.url).toContain('/shorts/abcdefghijk');
  expect(saved.parts[0].status).toBe('published');
  expect(saved.parts[0].metrics.at(-1)?.platform).toBe('youtube');
});
