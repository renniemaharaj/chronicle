import { readFile } from 'node:fs/promises';
import { expect, type Page } from '@playwright/test';
import type { Project } from '../src/lib/types/project';
import { mockSpeech } from './speech-mock';

export async function campaignFixture(page: Page, openStudio = true) {
  const project: Project = JSON.parse(await readFile('../internal/fixtures/person.json', 'utf8'));
  project.id = `synthetic-campaign-${Date.now()}`;
  project.title = 'Synthetic campaign test, not historical evidence';
  project.fixture = false;

  const saved = await page.request.put(`/api/projects/${project.id}`, { data: project });
  expect(saved.ok(), await saved.text()).toBe(true);

  const canonical: Project = await saved.json();
  await mockSpeech(page);
  await page.addInitScript((value) => {
    localStorage.setItem(
      'chronicle-workspace',
      JSON.stringify({
        projectId: value.id,
        projects: [value],
      }),
    );
  }, canonical);
  await page.goto('/');
  await page
    .getByRole('button', {
      name: 'Campaigns',
      exact: true,
    })
    .click();
  await page
    .getByRole('button', {
      name: 'Create series',
      exact: true,
    })
    .click();
  await page.getByLabel('Target seconds per reel').fill('15');
  await page
    .getByRole('button', {
      name: 'Build series',
      exact: true,
    })
    .click();
  await expect(page.getByRole('table', { name: 'Series episodes' })).toBeVisible();
  if (openStudio) {
    await page
      .getByRole('button', {
        name: 'Preview episode 1',
        exact: true,
      })
      .click();
    await expect(page.getByLabel('Vertical reel capture region')).toBeVisible();
  }

  const response = await page.request.get('/api/campaigns');
  const campaigns = await response.json();

  return campaigns.find((value: { projectId: string }) => value.projectId === project.id);
}

export async function mockOBS(page: Page, requests: string[], portrait = true) {
  let active = false;
  await page.routeWebSocket(/ws:\/\/127\.0\.0\.1:4455/, (socket) => {
    socket.send(
      JSON.stringify({
        op: 0,
        d: {
          rpcVersion: 1,
          obsWebSocketVersion: '5.6.0',
        },
      }),
    );
    socket.onMessage((raw) => {
      const message = JSON.parse(raw.toString());
      if (message.op === 1) {
        socket.send(
          JSON.stringify({
            op: 2,
            d: { negotiatedRpcVersion: 1 },
          }),
        );

        return;
      }

      if (message.op !== 6) {
        return;
      }

      const { requestType, requestId } = message.d;
      requests.push(requestType);

      let responseData = {};
      if (requestType === 'GetVideoSettings') {
        responseData = {
          outputWidth: portrait ? 1080 : 1920,
          outputHeight: portrait ? 1920 : 1080,
        };
      } else if (requestType === 'GetCurrentProgramScene') {
        responseData = { currentProgramSceneName: 'Chronicle Reels' };
      } else if (requestType === 'GetRecordStatus') {
        responseData = {
          outputActive: active,
          outputDuration: active ? 15000 : 0,
        };
      } else if (requestType === 'StartRecord') {
        active = true;
      } else if (requestType === 'StopRecord') {
        active = false;
        responseData = { outputPath: '/synthetic/part-1.mp4' };
      }

      socket.send(
        JSON.stringify({
          op: 7,
          d: {
            requestType,
            requestId,
            requestStatus: {
              result: true,
              code: 100,
            },
            responseData,
          },
        }),
      );
    });
  });
}
