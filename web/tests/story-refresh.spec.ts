import { expect, test } from '@playwright/test';
import { storyModeFixture } from './story-mode-fixture';

test('updating a story targets its canonical project and shows bounded research progress', async ({
  page,
}) => {
  const project = await storyModeFixture(page);
  let requestedProject = '';
  await page.route('**/api/research', async (route) => {
    if (route.request().method() !== 'POST') {
      return route.fulfill({ json: [] });
    }

    requestedProject = route.request().postDataJSON().projectId;
    await route.fulfill({
      json: {
        id: 'synthetic-refresh',
        subject: project.title,
        baseHash: 'synthetic-original-revision',
        status: 'failed',
        phase: 'Research incomplete',
        error: 'Synthetic exhausted research budget; raise it and resume.',
        requests: 40,
        sources: 1,
        claims: 1,
        conflicts: 0,
        plan: [
          {
            question: 'Latest events',
            interval: 'Current',
            status: 'completed',
          },
          {
            question: 'Review gaps',
            interval: 'Historical',
            status: 'pending',
          },
        ],
        activities: [],
      },
    });
  });
  await page
    .getByRole('button', {
      name: 'Update story',
      exact: true,
    })
    .click();
  await page
    .getByRole('button', {
      name: 'Update current story',
      exact: true,
    })
    .click();
  await expect(page.getByText('Research 50% complete', { exact: false })).toBeVisible();
  await expect(
    page.getByRole('button', {
      name: 'Resume',
      exact: true,
    }),
  ).toBeVisible();
  expect(requestedProject).toBe(project.id);
});
