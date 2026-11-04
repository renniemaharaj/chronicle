import { writeFile } from 'node:fs/promises';
import { reelRequest } from './reel-backend.mjs';
import { facebookViews, youtubeViews, saveViews } from './reel-performance.mjs';

const [id, platform = 'both'] = process.argv.slice(2);
let campaign = await reelRequest(`/campaigns/${id}`);
const observed = [];
if (platform === 'both' || platform === 'facebook') {
  const rows = await facebookViews(campaign);
  campaign = await saveViews(campaign, rows);
  observed.push(...rows);
}

if (platform === 'both' || platform === 'youtube') {
  const rows = await youtubeViews(campaign);
  campaign = await saveViews(campaign, rows);
  observed.push(...rows);
}

const report = {
  campaignId: id,
  threshold: 10,
  observed,
  replace: observed.filter((row) => row.views < 10),
  preserve: observed.filter((row) => row.views >= 10),
};
const stamp = new Date().toISOString().replaceAll(':', '-');
const path = `.chronicle/social/rome-performance-${stamp}.json`;
await writeFile(path, JSON.stringify(report, null, 2));
process.stdout.write(JSON.stringify({ path, ...report }, null, 2));
