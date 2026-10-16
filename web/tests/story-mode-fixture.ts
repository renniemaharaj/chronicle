import { readFile } from 'node:fs/promises';
import type { Page } from '@playwright/test';
import type { Project } from '../src/lib/types/project';
import { mockSpeech } from './speech-mock';

export async function storyModeFixture(page: Page) {
  const project: Project = JSON.parse(await readFile('../internal/fixtures/person.json', 'utf8'));
  project.id = 'synthetic-story-mode';
  project.title = 'Synthetic story map test';
  project.fixture = false;
  project.events = project.events.slice(0, 5);
  const locations = [['test-west'], ['test-east'], ['test-east'], ['test-west', 'test-east'], []];
  project.events.forEach((event, index) => {
    event.time.start!.year = 1800 + index;
    event.locationIds = locations[index];
  });
  for (const [id, name, longitude] of [
    ['test-west', 'Synthetic western city', -60],
    ['test-east', 'Synthetic eastern city', 40],
  ] as const) {
    project.entities.push({
      id,
      name,
      kind: 'place',
      summary: 'Synthetic test coordinates, not historical evidence.',
      time: {},
      attributes: { mapLocation: {
        latitude: 20,
        longitude,
        kind: 'city',
        sourceURL: 'https://example.org/synthetic',
        provider: 'Synthetic test gazetteer',
        note: 'Synthetic test reference, not historical evidence.',
      } },
      version: 1,
    });
  }
  project.chapters = [
    {
      id: 'first', title: 'Synthetic first chapter', summary: '', claimIds: [],
      eventIds: project.events.slice(0, 3).map((event) => event.id), interpretation: true,
    },
    {
      id: 'last', title: 'Synthetic last chapter', summary: '', claimIds: [],
      eventIds: project.events.slice(3).map((event) => event.id), interpretation: true,
    },
  ];
  await mockSpeech(page);
  await page.addInitScript((saved) => {
    localStorage.setItem('chronicle-workspace', JSON.stringify({
      projectId: saved.id, projects: [saved], selectedId: saved.events[0].id,
    }));
  }, project);
  await page.route('**/api/projects', (route) => route.fulfill({ json: [project] }));
  await page.route('**/api/workspace/events', (route) => route.fulfill({
    contentType: 'text/event-stream', body: 'data: []\n\n',
  }));
  await page.route('https://tile.openstreetmap.org/**', (route) => route.fulfill({
    contentType: 'image/png',
    body: Buffer.from('iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mP8/x8AAwMCAO+jRZkAAAAASUVORK5CYII=', 'base64'),
  }));
  await page.goto('/');
  return project;
}
