import type { Project } from '../types/project';
import type { HistoricalEvent } from '../types/knowledge';
import type { MediaAsset, StoryVisual } from '../types/media';
import { safeSourceURL } from '../api/links';

export function imageURL(asset: MediaAsset): string | undefined {
  if (asset.generated) {
    return /^\/api\/media\/images\/[a-f0-9]{64}\.(png|jpg|gif)$/.test(asset.url) ? asset.url : undefined;
  }

  return safeSourceURL(asset.url);
}

export function usableImage(asset: MediaAsset) {
  return asset.kind === 'image' && !!imageURL(asset) && (!!asset.generated || asset.retrieved);
}

export function passageImage(project: Project, event?: HistoricalEvent) {
  return project.media?.find((asset) => asset.id === event?.visual?.mediaId && usableImage(asset));
}

export function visualOpacity(cue: StoryVisual, seconds: number) {
  const elapsed = seconds - cue.startSeconds;
  if (elapsed < 0 || elapsed >= cue.durationSeconds) {
    return 0;
  }

  if (cue.transitionSeconds === 0) {
    return 1;
  }

  return Math.max(
    0,
    Math.min(1, elapsed / cue.transitionSeconds, (cue.durationSeconds - elapsed) / cue.transitionSeconds),
  );
}

export function validVisual(cue: StoryVisual) {
  return (
    !!cue.mediaId &&
    [cue.startSeconds, cue.durationSeconds, cue.transitionSeconds].every(
      (n) => Number.isFinite(n) && n >= 0 && n <= 3600,
    ) &&
    cue.durationSeconds >= 1 &&
    cue.transitionSeconds <= 5 &&
    cue.transitionSeconds * 2 <= cue.durationSeconds
  );
}
