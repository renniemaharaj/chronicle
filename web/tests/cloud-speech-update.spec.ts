import { expect, test } from '@playwright/test';
import { storyModeFixture } from './story-mode-fixture';
import { mockCloudSpeech } from './cloud-speech-fixture';

test('a revised story regenerates the changed passage while reusing unchanged audio', async ({ page }) => {
  const cloud = await mockCloudSpeech(page);
  const project = await storyModeFixture(page);
  await page
    .getByRole('button', {
      name: 'Story',
      exact: true,
    })
    .click();
  await page
    .getByRole('button', {
      name: 'Voice & playback',
      exact: true,
    })
    .click();
  await page.getByLabel('Speech engine', { exact: true }).selectOption('elevenlabs');
  await page.getByLabel('ElevenLabs voice', { exact: true }).selectOption('test-a');
  await page
    .getByRole('button', {
      name: 'Generate missing audio',
      exact: true,
    })
    .click();
  await expect(page.getByText('Narration 100% complete', { exact: false })).toBeVisible();
  project.events[4].title = 'Synthetic revised passage';
  await page.reload();
  await page
    .getByRole('button', {
      name: 'Story',
      exact: true,
    })
    .click();
  await page
    .getByRole('button', {
      name: 'Voice & playback',
      exact: true,
    })
    .click();
  await expect(page.getByText('Narration 80% complete', { exact: false })).toBeVisible();
  await page
    .getByRole('button', {
      name: 'Generate missing audio',
      exact: true,
    })
    .click();
  await expect(page.getByText('Narration 100% complete', { exact: false })).toBeVisible();
  expect(cloud.generations()).toBe(6);
});
