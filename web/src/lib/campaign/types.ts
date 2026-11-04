import type { HistoricalInterval } from '../types/temporal';
import type { Platform } from './platforms';

export interface Passage {
  visual?: { cue: import('../types/media').StoryVisual; image: import('../types/media').MediaAsset };
  audioId?: string;
  audioStart?: number;
  audioEnd?: number;
  eventId?: string;
  contextId?: string;
  title: string;
  text: string;
  start: number;
  end: number;
  claimIds: string[];
  time: HistoricalInterval;
}

export interface Metrics {
  platform?: Platform;
  observedAt: string;
  views?: number;
  reactions?: number;
  comments?: number;
  shares?: number;
  revenue?: number;
  currency?: string;
  method: string;
}

export interface Part {
  publications?: Partial<Record<Platform, { url: string; recordedAt: string; method: string }>>;
  number: number;
  status: 'planned' | 'recording' | 'recorded' | 'published';
  passages: Passage[];
  takeId?: string;
  file?: string;
  duration?: number;
  publication?: { url: string; recordedAt: string; method: string };
  metrics: Metrics[];
  takes?: { id: string; file: string; duration: number; completedAt: string }[];
}

export interface Campaign {
  accounts?: Partial<Record<Platform, { id: string; name: string; url: string; recordedAt?: string }>>;
  platforms?: Platform[];
  engine?: string;
  readyThrough?: string;
  remainingEvents?: number;
  id: string;
  projectId: string;
  projectHash: string;
  title: string;
  revision: number;
  createdAt: string;
  seconds: number;
  parts: Part[];
  gaps: string[];
  sources: { id: string; title: string; url: string; hash: string }[];
  audit: { at: string; action: string; part: number }[];
}

export interface Action {
  platform?: Platform;
  action: 'begin' | 'complete' | 'interrupt' | 'retake' | 'publish' | 'metrics';
  takeId?: string;
  file?: string;
  duration?: number;
  url?: string;
  metrics?: Metrics;
}
