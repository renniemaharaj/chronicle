import { describe, expect, it } from 'vitest';
import { emptyProject } from '../state/empty';
import type { HistoricalEvent } from '../types/knowledge';
import type { MediaAsset } from '../types/media';
import { mediaEvents, mediaForEvent, mediaSourceURL, youtubeThumbnail } from './media';

const video: MediaAsset = {
  id: 'synthetic-video', entityId: 'synthetic-person', sourceId: 'source',
  kind: 'video', url: 'https://www.youtube.com/watch?v=abcdefghijk', title: 'Synthetic performance',
  caption: '', creator: '', license: '', claimIds: ['clip-identity'], retrieved: true, useAsAvatar: false,
};

const event: HistoricalEvent = {
  id: 'performance', title: 'Synthetic concert', kind: 'Performance', summary: '', time: {},
  participants: ['synthetic-person'], locationIds: [], claimIds: ['concert-claim'], importance: 1, phase: '',
};

describe('media previews and chronology', () => {
  it('recognizes supported YouTube URL forms without trusting lookalike hosts or injected paths', () => {
    for (const address of [
      'https://www.youtube.com/watch?v=abcdefghijk&t=50',
      'https://youtu.be/abcdefghijk',
      'https://youtube.com/shorts/abcdefghijk',
      'https://www.youtube-nocookie.com/embed/abcdefghijk',
    ]) {
      expect(youtubeThumbnail(address)).toBe('https://i.ytimg.com/vi/abcdefghijk/hqdefault.jpg');
    }
    for (const address of ['https://youtube.com.evil/watch?v=abcdefghijk', 'https://youtube.com/watch?v=bad', 'javascript://youtube.com/watch?v=abcdefghijk', 'bad']) {
      expect(youtubeThumbnail(address)).toBeUndefined();
    }
    expect(mediaSourceURL('https://www.youtube.com/embed/abcdefghijk')).toBe('https://www.youtube.com/watch?v=abcdefghijk');
    expect(mediaSourceURL('javascript:alert(1)')).toBeUndefined();
  });

  it('links a clip through canonical claims, never through subject membership alone', () => {
    const project = {
      ...emptyProject, events: [event], media: [video],
      claims: [{ id: 'clip-identity', entityId: 'synthetic-person', eventId: event.id,
        text: 'Synthetic clip identification', status: 'supported_limited' as const, evidenceIds: [], assessment: [] }],
    };
    expect(mediaEvents(project, video)).toEqual([event]);
    expect(mediaForEvent(project, event)).toEqual([video]);
    expect(mediaForEvent({ ...project, claims: [] }, event)).toEqual([]);
    expect(mediaForEvent({ ...project, media: [{ ...video, retrieved: false }] }, event)).toEqual([]);
    expect(mediaForEvent(project)).toEqual([]);
  });

  it('also finds media attached directly to an event claim', () => {
    expect(mediaEvents({ ...emptyProject, events: [event] }, { ...video, claimIds: ['concert-claim'] })).toEqual([event]);
  });
});
