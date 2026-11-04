import { mkdir, writeFile, rename } from 'node:fs/promises';
import { resolve } from 'node:path';
import { reelRequest, reelAction } from './reel-backend.mjs';
import { captureMapPassages } from './reel-map-capture.mjs';
import { encodeEpisode } from './reel-media.mjs';
import { episodeDescription } from './reel-description.mjs';

const [id, maximum = 'Infinity', firstPart = '1', profile = 'map-motion-v1'] = process.argv.slice(2);
if (!id) {
  throw new Error(
    'Usage: node scripts/render-reels.mjs CAMPAIGN_ID [MAX_EPISODES] [FIRST_PART] [MAP_PROFILE]',
  );
}

const directory = resolve('.chronicle/exports', id);
await mkdir(directory, { recursive: true });

let campaign = await reelRequest(`/campaigns/${encodeURIComponent(id)}`);
const project = await reelRequest(`/projects/${campaign.projectId}`);
const parts = campaign.parts.slice(Number(firstPart) - 1, Number(firstPart) - 1 + Number(maximum));
await writeFile(resolve(directory, 'campaign.snapshot.json'), JSON.stringify(campaign, null, 2));

const mapDirectory = await captureMapPassages(
  {
    ...campaign,
    parts: parts.filter((part) => part.status === 'planned'),
  },
  directory,
  Infinity,
  profile,
);

for (const planned of parts) {
  campaign = await reelRequest(`/campaigns/${id}`);

  const part = campaign.parts[planned.number - 1];
  if (part.status !== 'planned') {
    process.stdout.write(`Episode ${part.number}: ${part.status}; preserved without another take.\n`);
    continue;
  }

  campaign = await reelAction(campaign, part.number, { action: 'begin' });

  const take = campaign.parts[part.number - 1];
  const file = resolve(
    directory,
    `rome-${String(part.number).padStart(2, '0')}-${take.takeId.slice(-6)}.mp4`,
  );
  const temporary = file.replace('.mp4', '.pending.mp4');
  try {
    const composition = { ...take };
    const duration = await encodeEpisode(composition, directory, temporary, mapDirectory);
    await rename(temporary, file);

    const metadata = episodeDescription(campaign, composition, project);
    metadata.visual = 'Focused map transitions, static image cues and 120 fps blended output';
    await writeFile(`${file}.json`, JSON.stringify({ ...metadata, file, duration }, null, 2));
    await writeFile(file.replace('.mp4', '.txt'), metadata.description);
    campaign = await reelAction(campaign, part.number, {
      action: 'complete',
      takeId: take.takeId,
      file,
      duration,
    });
    await writeFile(resolve(directory, 'campaign.snapshot.json'), JSON.stringify(campaign, null, 2));
    process.stdout.write(
      `Episode ${part.number}/${campaign.parts.length}: ${duration.toFixed(1)}s saved to ${file}\n`,
    );
  } catch (error) {
    const latest = await reelRequest(`/campaigns/${id}`);
    if (latest.parts[part.number - 1].status === 'recording') {
      await reelAction(latest, part.number, {
        action: 'interrupt',
        takeId: take.takeId,
      });
    }

    throw error;
  }
}
