import { readFile, writeFile, rename } from 'node:fs/promises';
import { reelRequest, reelAction } from './reel-backend.mjs';

const path = '.chronicle/social/posting-progress.json';

export async function postingProgress() {
  return JSON.parse(await readFile(path, 'utf8'));
}

export async function postingAttempt(platform, part, update) {
  const progress = await postingProgress();
  let attempt = progress.attempts.find((item) => item.platform === platform && item.part === part);
  if (!attempt) {
    attempt = { platform, part };
    progress.attempts.push(attempt);
  }

  Object.assign(attempt, update, { updatedAt: new Date().toISOString() });
  await writeFile(`${path}.pending`, JSON.stringify(progress, null, 2));
  await rename(`${path}.pending`, path);

  return attempt;
}

export async function confirmPublication(id, platform, part, url) {
  const campaign = await reelRequest(`/campaigns/${id}`);
  const existing = campaign.parts[part - 1].publications?.[platform];
  if (existing && existing.url !== url.replace(/\/$/, '')) {
    throw new Error('The saved publication differs; inspect before continuing.');
  }

  const saved = existing
    ? campaign
    : await reelAction(campaign, part, {
        action: 'publish',
        platform,
        url,
      });
  await postingAttempt(platform, part, {
    state: 'confirmed',
    url,
  });

  return saved;
}

export async function verifyPostingTarget(campaign, platform, accountID) {
  if (!accountID || campaign.accounts?.[platform]?.id !== accountID) {
    throw new Error('The campaign does not match the observed signed-in destination.');
  }

  const progress = await postingProgress();
  if (progress.campaignId !== campaign.id) {
    throw new Error('The browser publishing session belongs to another campaign.');
  }
}
