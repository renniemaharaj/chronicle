import { expect, test } from '@playwright/test';
import { storyModeFixture } from './story-mode-fixture';

test('story fullscreen keeps voice settings available and returns to the map view', async ({ page }) => {
  await storyModeFixture(page);
  await page
    .getByRole('button', {
      name: 'Show map view',
      exact: true,
    })
    .click();
  await page
    .getByRole('button', {
      name: 'Story',
      exact: true,
    })
    .click();

  const mode = page.getByRole('main', {
    name: 'Story mode',
    exact: true,
  });
  await mode
    .getByRole('button', {
      name: 'Toggle story fullscreen',
      exact: true,
    })
    .click();
  await expect
    .poll(() => page.evaluate(() => document.fullscreenElement?.classList.contains('story-mode')))
    .toBe(true);
  await mode
    .getByRole('button', {
      name: 'Voice & playback',
      exact: true,
    })
    .click();
  await expect(
    page.getByRole('dialog', {
      name: 'Voice & playback',
      exact: true,
    }),
  ).toBeVisible();
  await page.getByLabel('Voice', { exact: true }).selectOption('test-b');
  await page.getByLabel('Close dialog').click();
  await page.evaluate(() => document.exitFullscreen());
  await page.keyboard.press('Escape');
  await expect(mode).toHaveCount(0);
  await expect(page.getByLabel('Historical places map', { exact: true })).toBeVisible();
});
