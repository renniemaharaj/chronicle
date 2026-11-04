import { mkdir, readFile, rename, writeFile } from 'node:fs/promises';
import { resolve } from 'node:path';
import { reelRequest } from './reel-backend.mjs';
import { encodeEpisode } from './reel-media.mjs';
import { episodeDescription } from './reel-description.mjs';
import { captureMapPassages } from './reel-map-capture.mjs';

const [campaignId, maximum = '20'] = process.argv.slice(2);
const campaign = await reelRequest(`/campaigns/${campaignId}`);
const project = await reelRequest(`/projects/${campaign.projectId}`);
const plan = JSON.parse(await readFile('.chronicle/social/rome-engagement-plan.json', 'utf8'));
const directory = resolve('.chronicle/exports', campaignId);
const outputDirectory = resolve(directory, 'illustrated-v2');
await mkdir(outputDirectory, { recursive: true });

const eligible = new Set(plan.replace.map((item) => item.part));
const parts = campaign.parts.slice(0, Number(maximum)).filter((part) => eligible.has(part.number));
const mapDirectory = await captureMapPassages({ ...campaign, parts }, directory);
const manifest = {
  campaignId,
  createdAt: new Date().toISOString(),
  publicationPlan: plan,
  episodes: [],
};

for (const original of parts) {
  if (!eligible.has(original.number)) {
    continue;
  }

  const part = structuredClone(original);
  for (const passage of part.passages) {
    const event = project.events.find((item) => item.id === passage.eventId);
    const image = project.media.find((item) => item.id === event.visual?.mediaId);
    if (!image?.generated) {
      throw new Error(`Missing generated illustration: ${event.id}`);
    }

    passage.visual = {
      cue: event.visual,
      image,
    };
  }

  const file = resolve(outputDirectory, `rome-${String(part.number).padStart(2, '0')}-illustrated-v2.mp4`);
  const temporary = file.replace('.mp4', '.pending.mp4');
  const duration = await encodeEpisode(part, directory, temporary, mapDirectory);
  await rename(temporary, file);

  const metadata = {
    ...episodeDescription(campaign, part, project),
    file,
    duration,
    originalFile: original.file,
    originalReceipts: original.publications,
    visuals: [...new Map(part.passages.map((item) => [item.eventId, item.visual])).values()],
    visual: 'Static portrait AI illustrations with map context every 26 seconds',
    replacementTargets: plan.replace.filter((item) => item.part === part.number),
  };
  await writeFile(`${file}.json`, JSON.stringify(metadata, null, 2));
  await writeFile(file.replace('.mp4', '.txt'), metadata.description);
  manifest.episodes.push({
    number: part.number,
    ...metadata,
  });
  await writeFile(resolve(outputDirectory, 'manifest.json'), JSON.stringify(manifest, null, 2));
  process.stdout.write(
    `Episode ${part.number}: ${duration.toFixed(1)}s, ${metadata.visuals.length} illustrations.\n`,
  );
}
