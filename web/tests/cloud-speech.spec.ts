import { expect, test } from '@playwright/test';
import { storyModeFixture } from './story-mode-fixture';
import { mockCloudSpeech } from './cloud-speech-fixture';

test('ElevenLabs detects server keys and reuses each saved voice version during story playback', async ({
  page,
}, testInfo) => {
  const cloud = await mockCloudSpeech(page);
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
  await expect(page.getByText('API key detected: environment or server configuration.')).toBeVisible();
  await expect(page.getByText('Plan: synthetic · active')).toBeVisible();
  await expect(page.getByLabel('ElevenLabs API key')).toHaveCount(0);
  await page.getByLabel('ElevenLabs voice', { exact: true }).selectOption('test-a');
  await expect(page.getByText('Saved narration: 0 / 5 audio parts.')).toBeVisible();
  expect(cloud.generations()).toBe(0);
  await page
    .getByRole('button', {
      name: 'Generate missing audio',
      exact: true,
    })
    .click();
  await expect(page.getByText('Saved narration: 5 / 5 audio parts.')).toBeVisible();
  expect(cloud.generations()).toBe(5);

  await page.getByLabel('ElevenLabs voice', { exact: true }).selectOption('test-b');
  await expect(page.getByText('Saved narration: 0 / 5 audio parts.')).toBeVisible();
  await page
    .getByRole('button', {
      name: 'Generate missing audio',
      exact: true,
    })
    .click();
  await expect(page.getByText('Saved narration: 5 / 5 audio parts.')).toBeVisible();
  expect(cloud.generations()).toBe(10);
  await page.getByLabel('ElevenLabs voice', { exact: true }).selectOption('test-a');
  await expect(page.getByText('Saved narration: 5 / 5 audio parts.')).toBeVisible();
  await expect(
    page.getByRole('button', {
      name: 'Generate missing audio',
      exact: true,
    }),
  ).toBeDisabled();
  await page.getByLabel('Speech engine', { exact: true }).selectOption('browser');
  await page.getByLabel('Speech engine', { exact: true }).selectOption('elevenlabs');
  await expect(page.getByLabel('ElevenLabs voice', { exact: true })).toHaveValue('test-a');
  await expect(page.getByText('Saved narration: 5 / 5 audio parts.')).toBeVisible();
  await page.screenshot({ path: testInfo.outputPath('elevenlabs-settings.png') });
  await page.getByLabel('Close dialog').click();
  await mode
    .getByRole('button', {
      name: 'Play story',
      exact: true,
    })
    .click();
  await expect
    .poll(() =>
      page.evaluate(
        () => (window as unknown as { __cloudPlayback: { plays: number } }).__cloudPlayback.plays,
      ),
    )
    .toBe(1);
  await mode
    .getByRole('button', {
      name: 'Next milestone',
      exact: true,
    })
    .click();
  await expect
    .poll(() =>
      page.evaluate(
        () => (window as unknown as { __cloudPlayback: { plays: number } }).__cloudPlayback.plays,
      ),
    )
    .toBe(2);
  await mode.getByLabel('Playback speed', { exact: true }).selectOption('1.5');
  await expect
    .poll(() =>
      page.evaluate(
        () => (window as unknown as { __cloudPlayback: { plays: number } }).__cloudPlayback.plays,
      ),
    )
    .toBe(3);
  expect(cloud.generations()).toBe(10);
});
