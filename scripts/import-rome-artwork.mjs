import { readFile, stat, writeFile } from 'node:fs/promises';
import { resolve } from 'node:path';
import { backend, reelRequest } from './reel-backend.mjs';

const directory = resolve('.chronicle/illustrations/rome/engagement-v2');
const campaignId = 'campaign-JBYXU6PRNA3KVM4XCWF5ZYT5PM';
const campaign = await reelRequest(`/campaigns/${campaignId}`);
const project = await reelRequest('/projects/rome-history');
const prompts = JSON.parse(await readFile(resolve(directory, 'prompts.json'), 'utf8'));
const imported = [];

for (const [eventId, prompt] of Object.entries(prompts)) {
  const event = project.events.find((item) => item.id === eventId);
  const path = resolve(directory, `${eventId}.png`);
  const bytes = await readFile(path);
  const response = await fetch(`${backend}/api/media/images`, {
    method: 'POST',
    headers: { 'Content-Type': 'image/png' },
    body: bytes,
  });
  const image = await response.json();
  if (!response.ok) {
    throw new Error(JSON.stringify(image));
  }

  const evidenceIds = new Set(
    project.claims.filter((claim) => event.claimIds.includes(claim.id)).flatMap((claim) => claim.evidenceIds),
  );
  const references = [
    ...new Set(project.evidence.filter((item) => evidenceIds.has(item.id)).map((item) => item.sourceId)),
  ];
  const asset = {
    id: `rome-ai-${eventId}-v2`,
    entityId: 'rome',
    sourceId: '',
    kind: 'image',
    url: image.url,
    title: event.title,
    caption:
      'AI-generated interpretive illustration. Composition, people and setting are imaginative, not historical evidence or an exact reconstruction of a documented scene.',
    creator: 'AI illustration: image_gen',
    license: '',
    claimIds: event.claimIds,
    useAsAvatar: false,
    retrieved: false,
    depicts: event.time,
    generated: {
      tool: 'image_gen',
      prompt,
      createdAt: (await stat(path)).mtime.toISOString(),
      referenceSourceIds: references,
    },
  };
  if (eventId === 'foundation-legend') {
    asset.caption =
      'AI-generated illustration of the Romulus and Remus legend, not an attested event or archaeological evidence. The traditional foundation date is preserved as a tradition.';
  }

  const previous = project.media.findIndex((item) => item.id === asset.id);
  if (previous >= 0) {
    project.media[previous] = asset;
  } else {
    project.media.push(asset);
  }

  event.visual = {
    mediaId: asset.id,
    startSeconds: 2,
    durationSeconds: 10,
    transitionSeconds: 0.5,
  };
  imported.push({
    eventId,
    path,
    ...image,
    mediaId: asset.id,
  });
}

for (const part of campaign.parts) {
  const readings = new Map(part.passages.map((passage) => [passage.audioId, passage]));
  for (const reading of readings.values()) {
    const event = project.events.find((item) => item.id === reading.eventId);
    if (!event.visual) {
      throw new Error(`Missing artwork for ${event.id}`);
    }

    event.visual.startSeconds = 1;
    event.visual.durationSeconds = Math.max(1, reading.audioEnd - 2);
    event.visual.transitionSeconds = 0.4;
  }
}

await writeFile(resolve(directory, 'imported.json'), JSON.stringify(imported, null, 2));
await writeFile(
  resolve(directory, 'project-arguments.json'),
  JSON.stringify({
    project,
    rebuildNarrative: true,
  }),
);
process.stdout.write(`Imported ${imported.length} illustrations; project awaits validation and save.\n`);
