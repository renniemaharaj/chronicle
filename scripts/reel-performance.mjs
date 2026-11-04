import { browserTool, browserValue, browserWait } from './reel-browser.mjs';
import { reelAction } from './reel-backend.mjs';

export async function facebookViews(campaign) {
  const account = campaign.accounts.facebook;
  const expected = new Set(
    campaign.parts.map((part) => part.publications?.facebook?.url.match(/reel\/(\d+)/)?.[1]),
  );
  expected.delete(undefined);
  await browserTool('browser_navigate', {
    url: `https://www.facebook.com/profile.php?id=${account.id}&sk=reels_tab`,
  });
  await browserWait(
    () => ({
      url: location.href,
      tiles: document.querySelectorAll('a[aria-label="Reel tile preview"]').length,
    }),
    (state) => new URL(state.url).searchParams.get('id') === account.id && state.tiles > 0,
  );

  const observed = new Map();
  for (let step = 0; step < 16; step++) {
    const tiles = await browserValue(() =>
      [...document.querySelectorAll('a[aria-label="Reel tile preview"]')].map((element) => ({
        url: element.href.split('?')[0].replace(/\/$/, ''),
        label: element.innerText.trim(),
      })),
    );
    for (const tile of tiles) {
      const id = tile.url.match(/reel\/(\d+)/)?.[1];
      const views = parseViewCount(tile.label);
      if (expected.has(id) && views !== undefined) {
        observed.set(id, {
          platform: 'facebook',
          url: tile.url,
          views,
          observedAt: new Date().toISOString(),
        });
      }
    }

    if (observed.size === expected.size) {
      return [...observed.values()];
    }

    await browserValue(() => {
      window.scrollTo(0, document.body.scrollHeight);

      return document.body.scrollHeight;
    });
    await new Promise((accept) => setTimeout(accept, 1500));
  }

  throw new Error('Some Facebook counts are missing or abbreviated; no replacement is authorized for them.');
}

function parseViewCount(label) {
  const match = label
    .trim()
    .replace(/\s*views?$/i, '')
    .replaceAll(',', '')
    .match(/^([\d.]+)\s*([kmb])?$/i);
  if (!match) {
    return undefined;
  }

  const scale =
    {
      k: 1_000,
      m: 1_000_000,
      b: 1_000_000_000,
    }[match[2]?.toLowerCase()] ?? 1;

  return Math.round(Number(match[1]) * scale);
}

export async function youtubeViews(campaign) {
  const account = campaign.accounts.youtube;
  const expected = campaign.parts.filter((part) => part.publications?.youtube).length;
  await browserTool('browser_navigate', {
    url: `https://studio.youtube.com/channel/${account.id}/videos/short`,
  });
  await browserWait(
    () => ({
      name: document.querySelector('ytcp-navigation-drawer')?.innerText,
      rows: document.querySelectorAll('ytcp-video-row').length,
    }),
    (state) => state.name?.includes(account.name) && state.rows >= expected,
  );

  const rows = await browserValue(() =>
    [...document.querySelectorAll('ytcp-video-row')].flatMap((element) => {
      const text = element.innerText;
      const id = element.querySelector('a[href$="/edit"]')?.href.match(/video\/([^/]+)/)?.[1];
      const counts = text.trim().split('\n').slice(-2);
      if (!id || !text.includes('Rome | Part') || !counts.every((count) => /^\d[\d,]*$/.test(count))) {
        return [];
      }

      return [
        {
          platform: 'youtube',
          url: `https://www.youtube.com/shorts/${id}`,
          views: Number(counts[0].replaceAll(',', '')),
          comments: Number(counts[1].replaceAll(',', '')),
          observedAt: new Date().toISOString(),
        },
      ];
    }),
  );

  return rows.filter((row) => campaign.parts.some((part) => part.publications?.youtube?.url === row.url));
}

export async function saveViews(campaign, rows) {
  for (const row of rows) {
    const part = campaign.parts.find((item) => item.publications?.[row.platform]?.url === row.url);
    if (!part) {
      throw new Error(`Observed URL does not match the campaign: ${row.url}`);
    }

    row.part = part.number;

    const metrics = {
      observedAt: row.observedAt,
      views: row.views,
    };
    if (row.comments !== undefined) {
      metrics.comments = row.comments;
    }

    campaign = await reelAction(campaign, part.number, {
      action: 'metrics',
      platform: row.platform,
      metrics,
    });
  }

  return campaign;
}
