import { expect, test } from '@playwright/test';
import { storyModeFixture } from './story-mode-fixture';
import { mockCloudSpeech } from './cloud-speech-fixture';

test('retrying partial generation keeps completed files and generates only missing passages', async ({
  page,
}) => {
  const cloud = await mockCloudSpeech(page, { failGenerationOnceAt: 2 });
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
  await expect(page.getByRole('alert')).toContainText('Completed files are saved');
  await expect(page.getByText('Saved narration: 1 / 5 audio parts.')).toBeVisible();
  expect(cloud.generations()).toBe(1);
  await page
    .getByRole('button', {
      name: 'Generate missing audio',
      exact: true,
    })
    .click();
  await expect(page.getByText('Saved narration: 5 / 5 audio parts.')).toBeVisible();
  await expect(page.getByRole('alert')).toHaveCount(0);
  expect(cloud.generations()).toBe(5);
  expect(cloud.attempts()).toBe(6);
});
