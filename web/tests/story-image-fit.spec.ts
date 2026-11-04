import { expect, test } from '@playwright/test';
import { storyModeFixture } from './story-mode-fixture';

test('floating artwork matches the reader width and leaves the map available', async ({ page }) => {
  await page.addInitScript(() => {
    localStorage.setItem('chronicle-speech', JSON.stringify({ provider: 'screen' }));
  });

  const project = await storyModeFixture(page);
  const response = await page.request.post('/api/media/images', {
    headers: { 'Content-Type': 'image/png' },
    data: Buffer.from(
      'iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mP8/x8AAwMCAO+jRZkAAAAASUVORK5CYII=',
      'base64',
    ),
  });
  const imported = await response.json();
  project.media = [
    {
      id: 'synthetic-fit-image',
      entityId: project.subjectId,
      sourceId: '',
      kind: 'image',
      url: imported.url,
      title: 'Synthetic image fit test',
      caption: 'Synthetic test, not historical evidence.',
      creator: 'Test',
      license: '',
      claimIds: [],
      useAsAvatar: false,
      retrieved: false,
      generated: {
        tool: 'test',
        prompt: 'Synthetic test',
        createdAt: '2026-01-01T00:00:00Z',
        referenceSourceIds: [],
      },
      depicts: project.events[0].time,
    },
  ];
  project.events[0].visual = {
    mediaId: project.media[0].id,
    startSeconds: 0,
    durationSeconds: 60,
    transitionSeconds: 0,
  };
  await page.reload();
  await page
    .getByRole('button', {
      name: 'Story',
      exact: true,
    })
    .click();

  const panel = page.getByLabel('Timed passage image', { exact: true });
  await expect(panel.locator('img')).toBeVisible();
  await expect
    .poll(() => panel.locator('img').evaluate((image) => (image as HTMLImageElement).naturalWidth))
    .toBe(1);

  const box = (await panel.boundingBox())!;
  const picture = (await panel.locator('img').boundingBox())!;
  const map = (await page.locator('.story-scene').boundingBox())!;
  expect(box.width - picture.width).toBeLessThan(3);
  expect(box.width).toBeLessThan(map.width * 0.5);
  expect(Math.abs(picture.width - picture.height)).toBeLessThan(2);

  const location = (await page.locator('.story-location').boundingBox())!;
  expect(location.x).toBeGreaterThanOrEqual(box.x + box.width);
  expect(box.y + box.height).toBeLessThanOrEqual(map.y + map.height);

  const reading = (await page.locator('.story-reading').boundingBox())!;
  expect(Math.abs(box.width - reading.width)).toBeLessThan(2);
  expect(location.x + location.width).toBeLessThanOrEqual(reading.x - 10);
  await page.setViewportSize({
    width: 420,
    height: 1000,
  });
  await expect
    .poll(async () => {
      const image = (await panel.boundingBox())!;
      const reader = (await page.locator('.story-reading').boundingBox())!;
      const place = (await page.locator('.story-location').boundingBox())!;

      return image.y + image.height <= reader.y && image.x + image.width <= place.x;
    })
    .toBe(true);
});
