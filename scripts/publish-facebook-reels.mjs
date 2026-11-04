import { readFile } from 'node:fs/promises';
import { reelRequest } from './reel-backend.mjs';
import { browserTool, browserClick, browserValue } from './reel-browser.mjs';
import { postingProgress, postingAttempt, confirmPublication, verifyPostingTarget } from './reel-posting.mjs';
import { facebookIdentity, prepareFacebook, facebookReceipt } from './reel-facebook.mjs';

const [id, maximum = '20'] = process.argv.slice(2);
let campaign = await reelRequest(`/campaigns/${id}`);
let previousURLs = [];
try {
  const plan = JSON.parse(await readFile('.chronicle/social/rome-engagement-plan.json', 'utf8'));
  if (plan.campaignId !== id) {
    previousURLs = plan.observed.filter((item) => item.platform === 'facebook').map((item) => item.url);
  }
} catch (error) {
  if (error.code !== 'ENOENT') {
    throw error;
  }
}

await browserTool('browser_tabs', {
  action: 'select',
  index: 0,
});

const observed = await facebookIdentity();
await verifyPostingTarget(campaign, 'facebook', observed.id);

const account = campaign.accounts.facebook;

for (const number of campaign.parts.slice(0, Number(maximum)).map((part) => part.number)) {
  campaign = await reelRequest(`/campaigns/${id}`);

  const part = campaign.parts[number - 1];
  if (part.publications?.facebook) {
    continue;
  }

  const progress = await postingProgress();
  const metadata = JSON.parse(await readFile(`${part.file}.json`, 'utf8'));
  const known = campaign.parts.flatMap((item) =>
    item.publications?.facebook ? [item.publications.facebook.url] : [],
  );
  known.push(...previousURLs);

  const attempt = progress.attempts.find((item) => item.platform === 'facebook' && item.part === number);
  if (attempt && attempt.state !== 'editing') {
    const url = await facebookReceipt(metadata.title, known, account);
    campaign = await confirmPublication(id, 'facebook', number, url);
    process.stdout.write(`Facebook ${number}: reconciled ${url}\n`);
    continue;
  }

  await postingAttempt('facebook', number, {
    state: 'editing',
    title: metadata.title,
  });
  await prepareFacebook(part, metadata);
  await postingAttempt('facebook', number, {
    state: 'submitting',
    submittedAt: new Date().toISOString(),
  });
  await browserClick('role=button[name="Post"]', `Publish Rome part ${number}`);

  const alerts = await browserValue(() =>
    [...document.querySelectorAll('[role=alert]')].map((e) => e.innerText),
  );
  await postingAttempt('facebook', number, {
    state: 'processing',
    confirmation: alerts,
  });

  const url = await facebookReceipt(metadata.title, known, account);
  campaign = await confirmPublication(id, 'facebook', number, url);
  process.stdout.write(`Facebook ${number}/${campaign.parts.length}: ${url}\n`);
}
