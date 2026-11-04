import { spawn } from 'node:child_process';
import { open } from 'node:fs/promises';
import { reelRequest } from './reel-backend.mjs';

const [id, maximum = '20'] = process.argv.slice(2);
if (!id) {
  throw new Error('Usage: node scripts/publish-reel-series.mjs CAMPAIGN_ID [MAX_EPISODES]');
}

const campaign = await reelRequest(`/campaigns/${id}`);
const runners = {
  facebook: 'publish-facebook-reels.mjs',
  youtube: 'publish-youtube-shorts.mjs',
};

const publishedCount = (platform) => campaign.parts.filter((part) => part.publications?.[platform]).length;

const platforms = [...(campaign.platforms ?? ['facebook'])].sort(
  (left, right) => publishedCount(left) - publishedCount(right),
);

for (const platform of platforms) {
  const output = await open(`.chronicle/social/${platform}-publisher.log`, 'a');
  process.stdout.write(`Resuming ${platform} in series order.\n`);

  try {
    const code = await new Promise((accept, reject) => {
      const child = spawn(process.execPath, [`scripts/${runners[platform]}`, id, maximum], {
        stdio: ['ignore', output.fd, output.fd],
      });
      child.on('error', reject);
      child.on('exit', accept);
    });
    if (code !== 0) {
      process.exitCode = 1;
      process.stdout.write(`${platform} paused; inspect its saved attempt and publisher log.\n`);
    }
  } finally {
    await output.close();
  }
}
