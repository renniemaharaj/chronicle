import { expect, test } from '@playwright/test';
import { storyModeFixture } from './story-mode-fixture';
import { mockCloudSpeech } from './cloud-speech-fixture';

test('restricted usage access leaves ElevenLabs voices usable and missing playback never generates audio', async ({
  page,
}) => {
  const cloud = await mockCloudSpeech(page, { restrictedUsage: true });
  await storyModeFixture(page);
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
      name: 'Voice & playback',
      exact: true,
    })
    .click();
  await page.getByLabel('Speech engine', { exact: true }).selectOption('elevenlabs');
  await expect(page.getByText(/Usage unavailable: Synthetic restricted usage permission/)).toBeVisible();
  await page.getByLabel('ElevenLabs voice', { exact: true }).selectOption('test-a');
  await expect(page.getByText('Saved narration: 0 / 5 audio parts.')).toBeVisible();
  await page.getByLabel('Close dialog').click();
  await mode
    .getByRole('button', {
      name: 'Play story',
      exact: true,
    })
    .click();
  await expect(mode.getByRole('alert')).toContainText('generate missing audio first');
  expect(cloud.generations()).toBe(0);
});
