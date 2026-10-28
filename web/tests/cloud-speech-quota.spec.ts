import { expect, test } from '@playwright/test';
import { storyModeFixture } from './story-mode-fixture';
import { mockCloudSpeech } from './cloud-speech-fixture';

test('credit exhaustion reports partial progress and restores saved audio after reload', async ({ page }) => {
  const cloud = await mockCloudSpeech(page, {
    failGenerationOnceAt: 2,
    quotaFailure: true,
  });
  await storyModeFixture(page);
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
  await expect(page.getByText('Narration 0% complete', { exact: false })).toBeVisible();
  await page
    .getByRole('button', {
      name: 'Generate missing audio',
      exact: true,
    })
    .click();
  await expect(page.getByText('Paused for credits.', { exact: false })).toBeVisible();
  await expect(page.getByRole('progressbar', { name: 'Narration completion' })).toHaveAttribute('value', '1');
  await expect(page.getByText('Narration 20% complete', { exact: false })).toBeVisible();
  expect(cloud.attempts()).toBe(2);
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
  await expect(page.getByText('Narration 20% complete', { exact: false })).toBeVisible();
  await page.getByText('Missing narration passages', { exact: true }).click();
  await expect(page.locator('details[open] li')).toHaveCount(4);
  await page
    .getByRole('button', {
      name: 'Generate missing audio',
      exact: true,
    })
    .click();
  await expect(page.getByText('Narration 100% complete', { exact: false })).toBeVisible();
  expect(cloud.generations()).toBe(5);
  expect(cloud.attempts()).toBe(6);
});
