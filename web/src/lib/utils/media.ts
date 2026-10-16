import { safeSourceURL } from '../api/links';
import type { HistoricalEvent } from '../types/knowledge';
import type { MediaAsset } from '../types/media';
import type { Project } from '../types/project';

function youtubeId(address: string): string | undefined {
  try {
    const url = new URL(address);
    if (!['https:', 'http:'].includes(url.protocol)) {
      return;
    }
    const host = url.hostname.replace(/^www\./, '');
    let id: string | null | undefined;
    if (host === 'youtu.be') {
      id = url.pathname.split('/')[1];
    } else if (['youtube.com', 'm.youtube.com', 'youtube-nocookie.com'].includes(host)) {
      id = url.searchParams.get('v');
      if (!id && /^\/(embed|shorts|live)\//.test(url.pathname)) {
        id = url.pathname.split('/')[2];
      }
    }

    if (id && /^[A-Za-z0-9_-]{11}$/.test(id)) {
      return id;
    }
  } catch {
    return undefined;
  }
}

export function youtubeThumbnail(address: string): string | undefined {
  const id = youtubeId(address);
  return id ? `https://i.ytimg.com/vi/${id}/hqdefault.jpg` : undefined;
}

export function mediaSourceURL(address: string): string | undefined {
  const id = youtubeId(address);
  return id ? `https://www.youtube.com/watch?v=${id}` : safeSourceURL(address);
}

// Match canonical claim links. A shared subject or year does not identify a performance.
export function mediaEvents(project: Project, asset: MediaAsset): HistoricalEvent[] {
  const claims = new Set(asset.claimIds);
  const eventIds = new Set(project.claims.filter((claim) => claims.has(claim.id)).map((claim) => claim.eventId));
  return project.events.filter((event) =>
    eventIds.has(event.id) || [...event.claimIds, ...(event.narration?.claimIds ?? [])].some((id) => claims.has(id)),
  );
}

export function mediaForEvent(project: Project, event?: HistoricalEvent): MediaAsset[] {
  if (!event) {
    return [];
  }

  return (project.media ?? []).filter((asset) =>
    asset.retrieved && safeSourceURL(asset.url) && mediaEvents(project, asset).some((item) => item.id === event.id),
  );
}
