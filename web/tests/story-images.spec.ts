import { expect, test } from '@playwright/test';
import { storyModeFixture } from './story-mode-fixture';

test('passage image fades back to the map, freezes on pause, and saves canonical timing', async ({
  page,
}) => {
  const errors: string[] = [];
  page.on('pageerror', (error) => errors.push(error.message));
  await page.addInitScript(() => {
    localStorage.setItem('chronicle-speech', JSON.stringify({ provider: 'screen' }));
  });

  const project = await storyModeFixture(page);
  const upload = await page.request.post('/api/media/images', {
    headers: { 'Content-Type': 'image/png' },
    data: Buffer.from(
      'iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mP8/x8AAwMCAO+jRZkAAAAASUVORK5CYII=',
      'base64',
    ),
  });
  expect(upload.ok()).toBe(true);

  const imported = await upload.json();
  project.media = [
    {
      id: 'synthetic-illustration',
      entityId: project.subjectId,
      sourceId: '',
      kind: 'image',
      url: imported.url,
      title: 'Synthetic period illustration',
      caption: 'Synthetic test reconstruction, not historical evidence.',
      creator: 'Test generator',
      license: '',
      claimIds: [],
      retrieved: false,
      useAsAvatar: false,
      generated: {
        tool: 'test',
        prompt: 'Synthetic test image',
        createdAt: '2026-01-01T00:00:00Z',
        referenceSourceIds: [],
      },
      depicts: structuredClone(project.events[0].time),
    },
  ];
  project.events[0].visual = {
    mediaId: 'synthetic-illustration',
    startSeconds: 1,
    durationSeconds: 3,
    transitionSeconds: 0.5,
  };

  const claim = project.claims.find((item) => item.id === project.events[0].claimIds[0])!;
  claim.status = 'supported_limited';
  claim.text = 'Synthetic timing fixture for testing passage images, not historical evidence. '.repeat(15);
  project.events[0].narration = {
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

  const overlay = page.getByLabel('Timed passage image', { exact: true });
  const play = page.getByRole('button', {
    name: 'Play story',
    exact: true,
  });
  await page.clock.install();
  await play.click();
  await page.clock.runFor(500);
  await expect(overlay).toHaveCount(0);
  await page.clock.runFor(800);
  await expect(overlay).toBeVisible();
  await expect(overlay).toContainText('AI illustration');
  await page
    .getByRole('button', {
      name: 'Pause story',
      exact: true,
    })
    .click();

  const opacity = await overlay.evaluate((element) => getComputedStyle(element).opacity);
  await page.clock.runFor(5000);
  await expect(overlay).toHaveCSS('opacity', opacity);
  await play.click();
  await page.clock.runFor(3300);
  await expect(overlay).toHaveCount(0);
  await expect(page.locator('.story-scene .leaflet-container')).toBeVisible();
  await page
    .getByRole('button', {
      name: 'Pause story',
      exact: true,
    })
    .click();
  await page.getByText('Passage image & timing', { exact: true }).click();
  await page.getByLabel('Image duration seconds', { exact: true }).fill('6');

  const saved = page.waitForResponse(
    (response) =>
      response.url().includes(`/api/projects/${project.id}`) && response.request().method() === 'PUT',
  );
  await page
    .getByRole('button', {
      name: 'Save image timing',
      exact: true,
    })
    .click();

  const response = await saved;
  expect(response.ok()).toBe(true);

  const record = await response.json();
  expect(record.events[0].visual.durationSeconds).toBe(6);
  expect(record.media[0].retrieved).toBe(false);
  expect(record.media[0].generated.tool).toBe('test');
  expect(errors).toEqual([]);
});
