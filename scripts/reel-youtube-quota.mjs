import { browserClick, browserValue } from './reel-browser.mjs';
import { postingAttempt } from './reel-posting.mjs';
import { youtubeState } from './reel-youtube.mjs';

export function quotaRetryAllowed(attempt, now = Date.now()) {
  const retryAfter = Date.parse(attempt?.retryAfter ?? '');

  return (
    attempt?.state === 'blocked' &&
    attempt.reason === 'daily-upload-limit' &&
    !attempt.draftUrl &&
    !attempt.submittedAt &&
    Number.isFinite(retryAfter) &&
    now >= retryAfter
  );
}

export async function recordQuotaBlock(number, message) {
  const blockedAt = new Date();
  await postingAttempt('youtube', number, {
    state: 'blocked',
    reason: 'daily-upload-limit',
    platformMessage: message,
    blockedAt: blockedAt.toISOString(),
    retryAfter: new Date(blockedAt.getTime() + 86400000).toISOString(),
  });
}

export async function dismissRejectedUpload() {
  let state = await browserValue(youtubeState);
  if (state.text.includes('Build your channel history')) {
    await browserClick('role=button[name="Got it"]');
    state = await browserValue(youtubeState);
  }

  if (state.text.includes('Daily upload limit reached')) {
    await browserClick('ytcp-uploads-dialog >> role=button[name="Close"]');
  }
}
