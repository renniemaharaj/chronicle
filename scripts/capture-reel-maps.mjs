import { resolve } from 'node:path';
import { reelRequest } from './reel-backend.mjs';
import { captureMapPassages } from './reel-map-capture.mjs';

const [id, limit = 'Infinity'] = process.argv.slice(2);
if (!id) {
  throw new Error('Usage: node scripts/capture-reel-maps.mjs CAMPAIGN_ID [MAX_PASSAGES]');
}

const campaign = await reelRequest(`/campaigns/${encodeURIComponent(id)}`);
await captureMapPassages(campaign, resolve('.chronicle/exports', id), Number(limit));
