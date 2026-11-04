import { mkdir, access, rename } from 'node:fs/promises';
import { resolve } from 'node:path';
import { openMapBrowser } from './reel-map-browser.mjs';
import { backend } from './reel-backend.mjs';
import { startCaptureClock, stopCaptureClock } from './reel-capture-clock.mjs';

async function tilesReady(page) {
  await page.waitForFunction(
    () => {
      const tiles = [...document.querySelectorAll('.leaflet-tile')];

      const visible = tiles.filter((image) => {
        const bounds = image.getBoundingClientRect();

        return (
          bounds.width > 0 &&
          bounds.height > 0 &&
          bounds.right > 0 &&
          bounds.bottom > 0 &&
          bounds.left < innerWidth &&
          bounds.top < innerHeight
        );
      });
      const loaded = visible.filter((image) => image.complete && image.naturalWidth > 0);

      return visible.length > 0 && loaded.length >= Math.ceil(visible.length * 0.85);
    },
    null,
    { timeout: 20000 },
  );
  if (await page.locator('.map-error').count()) {
    throw new Error('The story map could not load.');
  }
}

export async function captureMapPassages(campaign, directory, limit = Infinity, profile = 'map-motion-v1') {
  const output = resolve(directory, profile);
  await mkdir(output, { recursive: true });

  const project = await (await fetch(`${backend}/api/projects/${campaign.projectId}`)).json();
  const readings = [
    ...new Map(
      campaign.parts.flatMap((part) =>
        part.passages.map((passage) => [
          passage.audioId,
          {
            ...passage,
            part: part.number,
          },
        ]),
      ),
    ).values(),
  ];
  if (!readings.length) {
    return output;
  }

  const missing = await Promise.all(
    readings.slice(0, limit).map(async (reading) => {
      try {
        await access(resolve(output, `${reading.audioId}.webm`));

        return false;
      } catch {
        return true;
      }
    }),
  );
  if (!missing.includes(true)) {
    return output;
  }

  const events = [...project.events].sort(
    (a, b) =>
      (a.time.start?.year ?? Infinity) - (b.time.start?.year ?? Infinity) ||
      (a.time.start?.month ?? 1) - (b.time.start?.month ?? 1) ||
      (a.time.start?.day ?? 1) - (b.time.start?.day ?? 1) ||
      a.id.localeCompare(b.id),
  );
  const { page, browser } = await openMapBrowser(project);
  try {
    await tilesReady(page);

    let previous = events.findIndex((event) => event.id === readings[0].eventId);

    const select = async (index) => {
      await page.locator('#story-position').evaluate((element, value) => {
        element.value = String(value);
        element.dispatchEvent(new Event('input', { bubbles: true }));
      }, index);
    };

    for (const [index, reading] of readings.slice(0, limit).entries()) {
      const path = resolve(output, `${reading.audioId}.webm`);
      try {
        await access(path);
        previous = events.findIndex((event) => event.id === reading.eventId);
        continue;
      } catch {
        // Only completed screencasts are renamed to their final paths.
      }

      const next = events.findIndex((event) => event.id === reading.eventId);
      await select(next);
      await page.waitForTimeout(1500);
      await tilesReady(page);
      await select(previous);
      await page.waitForTimeout(1400);
      await tilesReady(page);
      await page.evaluate(
        (number) =>
          document.body.setAttribute(
            'data-reel-label',
            `CHRONICLE | ROME | PART ${String(number).padStart(2, '0')}`,
          ),
        reading.part,
      );

      const temporary = path.replace('.webm', '.pending.webm');
      await page.screencast.start({
        path: temporary,
        size: {
          width: 1080,
          height: 1920,
        },
        fps: 60,
      });
      await select(next);
      await startCaptureClock(page, reading);
      await page.waitForTimeout(Math.ceil(reading.audioEnd * 1000));
      await stopCaptureClock(page);
      await tilesReady(page);
      await page.screencast.stop();
      await rename(temporary, path);
      if (index === 0) {
        await page.screenshot({ path: resolve(output, `first-passage-${reading.eventId}.png`) });
      }

      previous = next;
      process.stdout.write(`Map ${index + 1}/${readings.length}: ${reading.title}\n`);
    }
  } finally {
    await browser.close();
  }

  return output;
}
