import type { Page } from '@playwright/test';

export async function mockCloudSpeech(
  page: Page,
  options: {
    restrictedUsage?: boolean;
    failGenerationOnceAt?: number;
    holdGenerationAt?: number;
    quotaFailure?: boolean;
  } = {},
) {
  const saved = new Set<string>();
  let generations = 0;
  let attempts = 0;
  let failed = false;

  let releasePending = () => {};

  const pending = new Promise<void>((resolve) => {
    releasePending = resolve;
  });
  await page.addInitScript(() => {
    const playback = {
      sources: [] as string[],
      plays: 0,
    };
    Object.assign(window, { __cloudPlayback: playback });
    Object.defineProperty(window, 'Audio', {
      configurable: true,
      value: class {
        src = '';
        playbackRate = 1;
        onended = () => {};
        onerror = () => {};
        play() {
          playback.sources.push(this.src);
          playback.plays++;

          return Promise.resolve();
        }
        pause() {}
        removeAttribute() {}
      },
    });
  });
  await page.route('**/api/speech/settings', (route) =>
    route.fulfill({
      json: {
        hasApiKey: true,
        keySource: 'environment/config',
      },
    }),
  );
  await page.route('**/api/speech/voices?provider=elevenlabs', (route) =>
    route.fulfill({
      json: [
        {
          voice_id: 'test-a',
          name: 'Synthetic voice A',
          labels: { accent: 'test' },
        },
        {
          voice_id: 'test-b',
          name: 'Synthetic voice B',
          labels: { accent: 'test' },
        },
      ],
    }),
  );
  await page.route('**/api/speech/models', (route) =>
    route.fulfill({
      json: [
        {
          model_id: 'eleven_multilingual_v2',
          name: 'Synthetic multilingual model',
        },
        {
          model_id: 'test-model',
          name: 'Synthetic second model',
        },
      ],
    }),
  );
  await page.route('**/api/speech/usage', (route) =>
    route.fulfill({
      status: options.restrictedUsage ? 403 : 200,
      json: options.restrictedUsage
        ? { error: 'Synthetic restricted usage permission' }
        : {
            tier: 'synthetic',
            status: 'active',
            character_count: 100,
            character_limit: 1000,
          },
    }),
  );
  await page.route('**/api/speech/prepare', async (route) => {
    const { generate, ...input } = route.request().postDataJSON();
    const key = JSON.stringify(input);
    const cached = saved.has(key);
    if (generate && !cached) {
      attempts++;
      if (!failed && attempts === options.failGenerationOnceAt) {
        failed = true;

        return route.fulfill({
          status: options.quotaFailure ? 402 : 502,
          json: {
            error: 'Synthetic provider failure.',
            code: options.quotaFailure ? 'insufficient_credits' : 'request_failed',
          },
        });
      }

      if (attempts === options.holdGenerationAt) {
        await pending;
      }

      saved.add(key);
      generations++;
    }

    return route.fulfill({
      json: {
        available: saved.has(key),
        cached,
        url: `/api/speech/audio/test-${saved.size}`,
      },
    });
  });

  return {
    generations: () => generations,
    attempts: () => attempts,
    releasePending: () => releasePending(),
  };
}
