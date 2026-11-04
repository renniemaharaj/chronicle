import { expect, test } from '@playwright/test';
import { campaignFixture, mockOBS } from './campaign-fixture';
import { finishReading, speechState } from './speech-mock';
import type { Campaign } from '../src/lib/campaign/types';

test('records exact passages, saves a receipt and resumes with the next unpublished part', async ({
  page,
}) => {
  const calls: string[] = [];
  await mockOBS(page, calls);

  const campaign: Campaign = await campaignFixture(page);
  let lostAcknowledgement = false;
  await page.route('**/api/campaigns/*/actions', async (route) => {
    if (route.request().postDataJSON().action === 'complete' && !lostAcknowledgement) {
      lostAcknowledgement = true;
      await route.fetch();
      await route.abort();

      return;
    }

    await route.continue();
  });

  const stage = page.getByLabel('Vertical reel capture region');
  const bounds = await stage.boundingBox();
  expect(bounds!.width / bounds!.height).toBeCloseTo(9 / 16, 2);
  await page
    .getByRole('button', {
      name: 'Record part with OBS',
      exact: true,
    })
    .click();
  for (let index = 0; index < campaign.parts[0].passages.length; index++) {
    await expect.poll(async () => (await speechState(page)).texts.length).toBe(index + 1);
    await expect(stage.locator('.campaign-caption')).toHaveText(
      campaign.parts[0].passages[index].text.trim(),
    );
    await finishReading(page);
  }

  await expect(page.getByText('Saved file: /synthetic/part-1.mp4')).toBeVisible();
  expect(calls.filter((name) => name === 'StartRecord')).toHaveLength(1);
  expect(calls.filter((name) => name === 'StopRecord')).toHaveLength(1);
  expect(lostAcknowledgement).toBe(true);
  expect((await speechState(page)).texts).toEqual(campaign.parts[0].passages.map((passage) => passage.text));
  await page.getByLabel('Published reel URL').fill('https://www.facebook.com/reel/123456789');
  await page
    .getByRole('button', {
      name: 'Save publication receipt',
      exact: true,
    })
    .click();
  await expect(page.getByRole('link', { name: 'View posted part 1' })).toBeVisible();
  await page.getByLabel('Views', { exact: true }).fill('0');
  await page
    .getByRole('button', {
      name: 'Save observed performance',
      exact: true,
    })
    .click();
  await expect(page.getByText('Revenue unknown', { exact: true })).toBeVisible();
  await page
    .getByRole('button', {
      name: 'Reload saved progress',
      exact: true,
    })
    .click();
  await expect(stage).toContainText('Part 2 /');

  const response = await page.request.get(`/api/campaigns/${campaign.id}`);
  const saved: Campaign = await response.json();
  expect(saved.parts[0].status).toBe('published');
  expect(saved.parts[0].metrics[0].views).toBe(0);
  expect(saved.parts[0].metrics[0].revenue).toBeUndefined();
  await stage.screenshot({ path: '../.chronicle/campaign-stage-test.png' });
});

test('rejects landscape OBS output without advancing the recording checkpoint', async ({ page }) => {
  const calls: string[] = [];
  await mockOBS(page, calls, false);

  const campaign: Campaign = await campaignFixture(page);
  await page
    .getByRole('button', {
      name: 'Record part with OBS',
      exact: true,
    })
    .click();
  await expect(page.getByRole('alert')).toContainText('Set OBS output to 9:16');
  expect(calls).not.toContain('StartRecord');

  const response = await page.request.get(`/api/campaigns/${campaign.id}`);
  const saved: Campaign = await response.json();
  expect(saved.parts[0].status).toBe('planned');
});

test('leaving a recovered open take preserves its checkpoint until explicitly resolved', async ({ page }) => {
  const campaign: Campaign = await campaignFixture(page);
  const started = await page.request.post(`/api/campaigns/${campaign.id}/actions`, {
    data: {
      id: campaign.id,
      revision: campaign.revision,
      part: 1,
      action: 'begin',
    },
  });
  expect(started.ok()).toBe(true);
  await page
    .getByRole('button', {
      name: 'Reload saved progress',
      exact: true,
    })
    .click();
  await expect(page.getByText('Recover an interrupted recording', { exact: true })).toBeVisible();
  await page
    .getByRole('button', {
      name: 'Back to research',
      exact: true,
    })
    .click();

  const response = await page.request.get(`/api/campaigns/${campaign.id}`);
  const saved: Campaign = await response.json();
  expect(saved.parts[0].status).toBe('recording');

  const cleanup = await page.request.post(`/api/campaigns/${campaign.id}/actions`, {
    data: {
      id: campaign.id,
      revision: saved.revision,
      part: 1,
      action: 'interrupt',
      takeId: saved.parts[0].takeId,
    },
  });
  expect(cleanup.ok()).toBe(true);
});
