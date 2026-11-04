import type { Campaign, Part } from './types';

export type Platform = 'facebook' | 'youtube';

export const platformNames: Record<Platform, string> = {
  facebook: 'Facebook Reels',
  youtube: 'YouTube Shorts',
};

export function campaignPlatforms(campaign?: Campaign): Platform[] {
  return campaign?.platforms?.length ? campaign.platforms : ['facebook'];
}

export function publicationFor(part: Part, platform: Platform) {
  return part.publications?.[platform] ?? (platform === 'facebook' ? part.publication : undefined);
}

export function hasPublication(part: Part) {
  return !!part.publication || Object.keys(part.publications ?? {}).length > 0;
}

export function nextPublication(campaign: Campaign, platform: Platform) {
  return campaign.parts.find((part) => !publicationFor(part, platform))?.number;
}

export function clipURL(campaign: Campaign, part: Part, takeId?: string) {
  return `/api/campaigns/${campaign.id}/parts/${part.number}/video${takeId ? `?take=${encodeURIComponent(takeId)}` : ''}`;
}
