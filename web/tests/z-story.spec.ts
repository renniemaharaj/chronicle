import { expect, test } from '@playwright/test';
import { readFile } from 'node:fs/promises';
import { finishReading, mockSpeech, speechState } from './speech-mock';

test('spoken stories follow subject chapters, preserve pause, and read the final event', async ({ page, request }) => {
  await mockSpeech(page);
  const project = JSON.parse(await readFile('../internal/fixtures/person.json', 'utf8'));
  project.id = 'story-test';
  project.title = 'Synthetic story playback test';
  project.fixture = false;
  project.chapters = [
    {
      id: 'begin',
      title: 'A subject-specific beginning',
      summary: 'First chapter',
      eventIds: [project.events[0].id],
      claimIds: [],
      interpretation: true,
    },
    {
      id: 'later',
      title: 'A different later chapter',
      summary: 'Later chapter',
      eventIds: project.events.slice(1).map((event: { id: string }) => event.id),
      claimIds: [],
      interpretation: true,
    },
  ];
  const response = await request.post('/api/mcp', {
    headers: { Accept: 'application/json, text/event-stream', 'MCP-Protocol-Version': '2025-11-25' },
    data: {
      jsonrpc: '2.0',
      id: 2,
      method: 'tools/call',
      params: {
        name: 'chronicle_save_project',
        arguments: { project, rebuildNarrative: false },
      },
    },
  });
  const saved = (await response.json()).result;
  expect(saved.isError).not.toBe(true);
  expect(
    saved.structuredContent.events.every((event: { narration: { text: string } }) => event.narration.text.length),
  ).toBe(true);
  await page.goto('/');
  await page.getByRole('button', { name: 'Library', exact: true }).click();
  await page.locator('.library-project').filter({ hasText: project.title }).click();
  await page.getByRole('button', { name: 'Play story', exact: true }).click();
  await expect(page.locator('.phase.current')).toContainText('A subject-specific beginning');
  await expect.poll(async () => (await speechState(page)).texts.length).toBe(1);
  await page.getByRole('button', { name: 'Pause story', exact: true }).click();
  await expect.poll(async () => (await speechState(page)).pauses).toBe(1);
  await page.getByRole('button', { name: 'Play story', exact: true }).click();
  await expect.poll(async () => (await speechState(page)).resumes).toBe(1);
  expect((await speechState(page)).texts).toHaveLength(1);
  await finishReading(page);
  await expect(page.locator('.phase.current')).toContainText('A different later chapter');
  await expect.poll(async () => (await speechState(page)).texts.length).toBe(2);
  await page.getByRole('button', { name: 'Voice & playback', exact: true }).click();
  await page.getByLabel('Voice', { exact: true }).selectOption('test-b');
  await page.getByLabel('Close dialog').click();
  await expect.poll(async () => (await speechState(page)).texts.length).toBe(3);
  for (let index = 1; index < project.events.length; index++) {
    if (index === project.events.length - 1) {
      const count = (await speechState(page)).texts.length;
      await page.getByRole('button', { name: 'Pause story', exact: true }).click();
      await page.getByRole('button', { name: 'Play story', exact: true }).click();
      expect((await speechState(page)).texts).toHaveLength(count);
    }
    await finishReading(page);
    if (index < project.events.length - 1) {
      await expect.poll(async () => (await speechState(page)).texts.length).toBe(index + 3);
    }
  }
  await expect(page.getByRole('button', { name: 'Play story', exact: true })).toBeVisible();
  await page.getByRole('button', { name: 'All history', exact: true }).click();
  await expect(page.locator('.phase.current')).toHaveCount(0);
  await expect(page.getByLabel('Follow chapter', { exact: true })).not.toBeChecked();
});
