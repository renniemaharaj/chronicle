import { expect, type Page } from '@playwright/test';
import type { Action, Campaign } from '../src/lib/campaign/types';

// Synthetic production observations exercise the interface, not historical evidence or real revenue.
export async function campaignProgressFixture(page: Page, initial: Campaign) {
  let campaign = initial;

  async function update(part: number, action: Action) {
    const response = await page.request.post(`/api/campaigns/${campaign.id}/actions`, {
      data: {
        id: campaign.id,
        revision: campaign.revision,
        part,
        ...action,
      },
    });
    expect(response.ok(), await response.text()).toBe(true);
    campaign = await response.json();
  }

  await update(1, { action: 'begin' });
  await update(1, {
    action: 'complete',
    takeId: campaign.parts[0].takeId,
    file: '/synthetic/first-episode.mp4',
    duration: 58,
  });
  await update(1, {
    action: 'publish',
    url: 'https://www.facebook.com/reel/987654321',
  });
  await update(1, {
    action: 'metrics',
    metrics: {
      observedAt: new Date(Date.now() - 10000).toISOString(),
      views: 1000,
      reactions: 80,
      comments: 10,
      shares: 5,
      revenue: 1,
      currency: 'USD',
      method: 'manual',
    },
  });
  await update(1, {
    action: 'metrics',
    metrics: {
      observedAt: new Date().toISOString(),
      views: 1200,
      reactions: 90,
      comments: 12,
      shares: 8,
      revenue: 2,
      currency: 'USD',
      method: 'manual',
    },
  });
  await update(2, { action: 'begin' });
  await update(2, {
    action: 'complete',
    takeId: campaign.parts[1].takeId,
    file: '/synthetic/second-episode.mp4',
    duration: 62,
  });

  return campaign;
}
