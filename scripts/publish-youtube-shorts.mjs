import { readFile } from 'node:fs/promises';
import { reelRequest } from './reel-backend.mjs';
import { browserTool, browserClick, browserValue, browserWait } from './reel-browser.mjs';
import { postingProgress, postingAttempt, confirmPublication, verifyPostingTarget } from './reel-posting.mjs';
import { quotaRetryAllowed, recordQuotaBlock, dismissRejectedUpload } from './reel-youtube-quota.mjs';
import {
  youtubeState,
  openYouTubeUpload,
  describeYouTube,
  readyYouTube,
  confirmedYouTube,
} from './reel-youtube.mjs';

const [id, maximum = '20', playlist = 'Rome: origins, Republic and Empire | Chronicle'] =
  process.argv.slice(2);
let campaign = await reelRequest(`/campaigns/${id}`);
await browserTool('browser_tabs', {
  action: 'select',
  index: 1,
});

const identity = await browserValue(() => ({
  id: location.pathname.match(/\/channel\/([^/]+)/)?.[1],
  name: document.querySelector('ytcp-navigation-drawer')?.innerText,
}));
await verifyPostingTarget(campaign, 'youtube', identity.id);
if (!identity.name?.includes(campaign.accounts.youtube.name)) {
  throw new Error('The observed YouTube channel name does not match the campaign destination.');
}

const initial = await browserValue(youtubeState);
if (initial.text.includes('Video published')) {
  await browserClick('ytcp-video-share-dialog button[aria-label="Close"]');
}

for (const number of campaign.parts.slice(0, Number(maximum)).map((part) => part.number)) {
  campaign = await reelRequest(`/campaigns/${id}`);

  const part = campaign.parts[number - 1];
  if (part.publications?.youtube) {
    continue;
  }

  const progress = await postingProgress();
  const attempt = progress.attempts.find((item) => item.platform === 'youtube' && item.part === number);
  if (attempt && !quotaRetryAllowed(attempt)) {
    throw new Error(
      `YouTube part ${number} has an existing draft or submission. Reconcile it before another upload.`,
    );
  }

  if (attempt) {
    await dismissRejectedUpload();
  }

  const metadata = JSON.parse(await readFile(`${part.file}.json`, 'utf8'));
  await postingAttempt('youtube', number, {
    state: 'attaching',
    title: metadata.title,
  });

  let draft;
  try {
    draft = await openYouTubeUpload(part);
  } catch (error) {
    if (error.message.includes('YouTube daily upload limit reached')) {
      await recordQuotaBlock(number, error.message);
    }

    throw error;
  }

  const url = draft.links.find((link) => /shorts\/[\w-]{11}/.test(link)).split('?')[0];
  await postingAttempt('youtube', number, {
    state: 'editing',
    draftUrl: url,
  });
  await describeYouTube(metadata);
  await browserClick('role=button[name="Select playlists"]');
  await browserClick(`role=listitem >> text="${playlist}"`, 'Add to the Rome series playlist');
  await browserClick('role=button[name="Done"]');
  await browserWait(youtubeState, (value) => value.text.includes(playlist));
  await readyYouTube();

  const ready = await browserValue(youtubeState);
  if (!ready.text.includes(metadata.title) || !ready.links.some((link) => link.startsWith(url))) {
    throw new Error('The visible YouTube episode does not match the intended upload.');
  }

  await postingAttempt('youtube', number, {
    state: 'submitting',
    submittedAt: new Date().toISOString(),
  });
  await browserClick('role=button[name="Publish"]', `Publish Rome part ${number}`);
  await confirmedYouTube(metadata.title, url);
  campaign = await confirmPublication(
    id,
    'youtube',
    number,
    url.replace('https://youtube.com/', 'https://www.youtube.com/'),
  );
  process.stdout.write(`YouTube ${number}/${campaign.parts.length}: ${url}\n`);
  await browserClick('ytcp-video-share-dialog button[aria-label="Close"]');
}
