import { readFile, rename, writeFile } from 'node:fs/promises';

const path = '.chronicle/social/rome-replacement-progress.json';

export async function replacementProgress(campaignId) {
  let progress;
  try {
    progress = JSON.parse(await readFile(path, 'utf8'));
  } catch (error) {
    if (error.code !== 'ENOENT') {
      throw error;
    }

    return {
      campaignId,
      attempts: [],
    };
  }

  if (progress.campaignId !== campaignId) {
    throw new Error('Replacement progress belongs to another campaign.');
  }

  return progress;
}

export async function replacementCheckpoint(campaignId, platform, number, update) {
  const progress = await replacementProgress(campaignId);
  const previous = progress.attempts.find((item) => item.platform === platform && item.part === number);
  const item = {
    ...previous,
    platform,
    part: number,
    ...update,
    updatedAt: new Date().toISOString(),
  };
  progress.attempts = progress.attempts.filter(
    (attempt) => attempt.platform !== platform || attempt.part !== number,
  );
  progress.attempts.push(item);
  await writeFile(`${path}.pending`, JSON.stringify(progress, null, 2));
  await rename(`${path}.pending`, path);

  return item;
}
