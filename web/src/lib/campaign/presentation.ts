import type { Campaign, Part } from './types';
import type { Project } from '../types/project';
import { hasPublication } from './platforms';
import { storyImage } from '../utils/story-outline';

export const statusLabels = {
  planned: 'Draft',
  recording: 'Recording',
  recorded: 'Ready',
  published: 'Published',
};

export function episodeTitle(part: Part) {
  return part.passages[0]?.title || `Episode ${part.number}`;
}

export function episodeSummary(part: Part) {
  return part.passages
    .map((passage) => passage.text)
    .join(' ')
    .replace(/\s+/g, ' ')
    .trim();
}

export function estimatedSeconds(part: Part) {
  return (episodeSummary(part).split(/\s+/).filter(Boolean).length * 60) / 140;
}

export function durationLabel(seconds: number) {
  const rounded = Math.round(seconds);

  return `${Math.floor(rounded / 60)}:${String(rounded % 60).padStart(2, '0')}`;
}

export function partDuration(part: Part) {
  return `${part.duration === undefined ? '~' : ''}${durationLabel(part.duration ?? estimatedSeconds(part))}`;
}

export function episodeImage(project: Project, part?: Part) {
  const ids = new Set(
    part?.passages.map((passage) => passage.eventId) ?? project.events.map((event) => event.id),
  );

  return storyImage(
    project,
    project.events.filter((event) => ids.has(event.id)),
  );
}

export function episodeSources(campaign: Campaign, part: Part, project: Project) {
  const claims = new Set(part.passages.flatMap((passage) => passage.claimIds));
  const evidence = new Set(
    project.claims.filter((claim) => claims.has(claim.id)).flatMap((claim) => claim.evidenceIds),
  );
  const sources = new Set(
    project.evidence.filter((item) => evidence.has(item.id)).map((item) => item.sourceId),
  );

  return (campaign.sources ?? []).filter((source) => sources.has(source.id));
}

export function seriesProgress(campaign: Campaign) {
  const published = campaign.parts.filter((part) => part.status === 'published').length;
  const recorded = campaign.parts.filter(
    (part) => part.status === 'recorded' || part.status === 'published',
  ).length;

  return {
    published,
    recorded,
    total: campaign.parts.length,
  };
}

export function episodeStatus(part: Part) {
  return hasPublication(part) && part.status !== 'published'
    ? 'Partially published'
    : statusLabels[part.status];
}
