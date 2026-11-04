import { expect, test } from '@playwright/test';
import { storyModeFixture } from './story-mode-fixture';

test('map emphasis follows narration, freezes on pause and fits a portrait production frame', async ({
  page,
}, testInfo) => {
  const errors: string[] = [];
  let tileRequests = 0;
  page.on('pageerror', (error) => errors.push(error.message));
  page.on('request', (request) => {
    if (request.url().includes('tile.openstreetmap.org')) {
      tileRequests++;
    }
  });
  await page.addInitScript(() => {
    localStorage.setItem('chronicle-speech', JSON.stringify({ provider: 'screen' }));
  });

  const project = await storyModeFixture(page);
  const event = project.events[0];
  const claim = project.claims.find((item) => item.id === event.claimIds[0])!;
  claim.status = 'supported_limited';
  claim.text = 'Synthetic story animation fixture, not historical evidence. '.repeat(30);
  event.narration = {
    text: claim.text,
    claimIds: [claim.id],
    interpretation: true,
  };
  await page.reload();
  await page
    .getByRole('button', {
      name: 'Story',
      exact: true,
    })
    .click();

  const focus = page.locator('.story-focus-outline');
  await expect(focus).toHaveCount(1);
  await page.clock.install();
  await page
    .getByRole('button', {
      name: 'Play story',
      exact: true,
    })
    .click();
  await page.clock.runFor(1000);

  const initial = await focus.getAttribute('stroke-dashoffset');
  const beforeMotion = tileRequests;
  await page.clock.runFor(1000);
  expect(await focus.getAttribute('stroke-dashoffset')).not.toBe(initial);
  expect(tileRequests).toBe(beforeMotion);
  await page
    .getByRole('button', {
      name: 'Pause story',
      exact: true,
    })
    .click();

  const paused = await focus.getAttribute('stroke-dashoffset');
  await page.clock.runFor(3000);
  await expect(focus).toHaveAttribute('stroke-dashoffset', paused!);
  await page
    .getByRole('button', {
      name: 'Reel and Shorts preview',
      exact: true,
    })
    .click();

  const stage = page.locator('.portrait-stage');
  await expect(stage).toBeVisible();

  const bounds = (await stage.boundingBox())!;
  expect(bounds.width / bounds.height).toBeCloseTo(9 / 16, 2);
  for (const selector of ['.story-reading', '.story-reel-caption', '.story-location']) {
    const content = (await stage.locator(selector).boundingBox())!;
    expect(content.x).toBeGreaterThanOrEqual(bounds.x + bounds.width * 0.05);
    expect(content.x + content.width).toBeLessThanOrEqual(bounds.x + bounds.width * 0.85);
    expect(content.y + content.height).toBeLessThanOrEqual(bounds.y + bounds.height * 0.79);
  }

  await expect(page.getByLabel('Reel safe area guide')).toBeVisible();
  await page
    .getByRole('button', {
      name: 'Show Reel safe areas',
      exact: true,
    })
    .click();
  await expect(page.getByLabel('Reel safe area guide')).toHaveCount(0);
  await page.screenshot({ path: testInfo.outputPath('story-portrait-preview.png') });
  await page.setViewportSize({
    width: 420,
    height: 1000,
  });

  const narrow = (await stage.boundingBox())!;
  expect(narrow.width / narrow.height).toBeCloseTo(9 / 16, 2);
  expect(narrow.x + narrow.width).toBeLessThanOrEqual(420);
  expect(errors).toEqual([]);
});
