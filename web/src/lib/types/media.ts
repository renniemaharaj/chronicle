export interface Resource {
  kind: 'image' | 'video' | 'link';
  url: string;
  label: string;
}

export interface MediaAsset {
  generated?: {
    tool: string;
    prompt: string;
    createdAt: string;
    referenceSourceIds: string[];
  };
  depicts?: import('./temporal').HistoricalInterval;
  id: string;
  entityId: string;
  sourceId: string;
  kind: 'image' | 'video';
  url: string;
  title: string;
  caption: string;
  creator: string;
  license: string;
  licenseUrl?: string;
  claimIds: string[];
  useAsAvatar: boolean;
  retrieved: boolean;
}

export interface StoryVisual {
  mediaId: string;
  startSeconds: number;
  durationSeconds: number;
  transitionSeconds: number;
}

export interface SubjectLink {
  id: string;
  entityId: string;
  sourceId: string;
  kind: 'website' | 'social';
  url: string;
  label: string;
  claimIds: string[];
  verified: boolean;
  verificationMethod: 'unverified' | 'official_reference' | 'official_site_link';
  verifiedAt?: string;
}
