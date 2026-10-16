import type { HistoricalDate, HistoricalInterval } from './temporal';

export interface Entity {
  id: string;
  name: string;
  kind: string;
  summary: string;
  time: HistoricalInterval;
  attributes: Record<string, unknown>;
  version: number;
}

export interface HistoricalEvent {
  identityKey?: string;
  id: string;
  title: string;
  kind: string;
  summary: string;
  time: HistoricalInterval;
  participants: string[];
  locationIds: string[];
  claimIds: string[];
  importance: number;
  phase: string;
  requiresPresence?: boolean;
  causeIds?: string[];
  consequenceIds?: string[];
}

export interface Relationship {
  id: string;
  fromId: string;
  toId: string;
  kind: string;
  time: HistoricalInterval;
  claimIds: string[];
}

export type Assessment = 'corroborated' | 'supported_limited' | 'disputed' | 'unverified' | 'rejected_superseded';

export interface Claim {
  id: string;
  text: string;
  status: Assessment;
  entityId?: string;
  eventId?: string;
  time?: HistoricalInterval;
  evidenceIds: string[];
  assessment: string[];
}

export interface Source {
  id: string;
  title: string;
  author: string;
  publisher: string;
  url: string;
  publishedAt?: HistoricalDate;
  retrievedAt: string;
  classification: 'primary' | 'secondary' | 'tertiary';
  independenceGroup: string;
  independenceReviewed?: boolean;
  contentHash?: string;
  limitations?: string[];
  verification: 'retrieved' | 'fixture' | 'unverified';
}

export interface Evidence {
  id: string;
  sourceId: string;
  claimId: string;
  excerpt: string;
  location: string;
  relation: 'supporting' | 'contradicting' | 'partially_supporting' | 'contextual' | 'inconclusive';
  mode: 'quotation' | 'paraphrase';
  verified: boolean;
}

export interface Conflict {
  id: string;
  title: string;
  claimIds: string[];
  kind: string;
  explanation: string;
  consequences: string[];
  severity: 'definite' | 'suspicious' | 'historical';
}
