import { expect, test } from '@playwright/test';
import { campaignFixture, mockOBS } from './campaign-fixture';
import { finishReading, speechState } from './speech-mock';
import type { Campaign } from '../src/lib/campaign/types';

test('series preview advances across episodes without opening recording checkpoints', async ({ page }) => {
  const requests: string[] = [];
  await mockOBS(page, requests);

  const campaign: Campaign = await campaignFixture(page, false);
  await page
    .getByRole('button', {
      name: 'Preview series',
      exact: true,
    })
    .click();
  await expect(
    page.getByRole('button', {
      name: 'Stop preview',
      exact: true,
    }),
  ).toBeVisible();
  await expect(
    page.getByRole('button', {
      name: 'Record part with OBS',
      exact: true,
    }),
  ).toBeDisabled();
  for (let index = 0; index < campaign.parts[0].passages.length; index++) {
    await expect.poll(async () => (await speechState(page)).texts.length).toBe(index + 1);
    await finishReading(page);
  }

  await expect(page.getByLabel('Vertical reel capture region')).toContainText('Part 2 /');
  await page
    .getByRole('button', {
      name: 'Stop preview',
      exact: true,
    })
    .click();

  const saved: Campaign = await (await page.request.get(`/api/campaigns/${campaign.id}`)).json();
  expect(saved.parts.every((part) => part.status === 'planned')).toBe(true);
  expect(requests).toEqual([]);
});
