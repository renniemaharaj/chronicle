import { expect, test } from '@playwright/test';
import { storyModeFixture } from './story-mode-fixture';
import { mockCloudSpeech } from './cloud-speech-fixture';

test('closing preparation finishes the current paid part and stops before requesting another', async ({
  page,
}) => {
  const cloud = await mockCloudSpeech(page, { holdGenerationAt: 1 });
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
  await page.getByLabel('ElevenLabs voice', { exact: true }).selectOption('test-a');
  await expect(page.getByText('Saved narration: 0 / 5 audio parts.')).toBeVisible();
  await page
    .getByRole('button', {
      name: 'Generate missing audio',
      exact: true,
    })
    .click();
  await expect.poll(() => cloud.attempts()).toBe(1);
  await page.getByLabel('Close dialog').click();
  cloud.releasePending();
  await expect.poll(() => cloud.generations()).toBe(1);
  await mode
    .getByRole('button', {
      name: 'Voice & playback',
      exact: true,
    })
    .click();
  await expect(page.getByText('Saved narration: 1 / 5 audio parts.')).toBeVisible();
  expect(cloud.attempts()).toBe(1);
  await page
    .getByRole('button', {
      name: 'Generate missing audio',
      exact: true,
    })
    .click();
  await expect(page.getByText('Saved narration: 5 / 5 audio parts.')).toBeVisible();
  expect(cloud.attempts()).toBe(5);
});
