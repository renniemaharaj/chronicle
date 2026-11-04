import { expect, test } from '@playwright/test';
import { storyModeFixture } from './story-mode-fixture';

test('reduced motion keeps geographic emphasis steady during playback', async ({ page }) => {
  await page.emulateMedia({ reducedMotion: 'reduce' });
  await page.addInitScript(() => {
    localStorage.setItem('chronicle-speech', JSON.stringify({ provider: 'screen' }));
  });
  await storyModeFixture(page);
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

  const path = await focus.getAttribute('d');
  await page.clock.runFor(600);
  await expect(focus).toHaveAttribute('stroke-dashoffset', '0');
  await expect(focus).toHaveAttribute('d', path!);
});
