import type { MediaAsset, SubjectLink } from './media';
import type { HistoricalInterval } from './temporal';
import type { Claim, Conflict, Entity, Evidence, HistoricalEvent, Relationship, Source } from './knowledge';

export interface Chapter {
  id: string;
  title: string;
  summary: string;
  eventIds: string[];
  claimIds: string[];
  interpretation: boolean;
}

export interface Coverage {
  startYear: number;
  endYear: number;
  status: 'researched' | 'partial' | 'disputed' | 'unexplored';
  question: string;
  notes: string;
}

export interface HistoricalContext {
  id: string;
  title: string;
  scope: string;
  time: HistoricalInterval;
  summary: string;
  claimIds: string[];
  interpretation: boolean;
}

export interface Project {
  media?: MediaAsset[];
  subjectLinks?: SubjectLink[];
  contextPeriods?: HistoricalContext[];
  researchRecords?: Omit<
    ResearchRun,
    'project' | 'question' | 'phase' | 'updatedAt' | 'subjectId' | 'sources' | 'claims' | 'conflicts'
  >[];
  narrativeId?: string;
  formatVersion: 1;
  id: string;
  subjectId: string;
  title: string;
  fixture: boolean;
  entities: Entity[];
  events: HistoricalEvent[];
  relationships: Relationship[];
  claims: Claim[];
  sources: Source[];
  evidence: Evidence[];
  conflicts: Conflict[];
  chapters: Chapter[];
  coverage: Coverage[];
  notes: { id: string; targetId: string; text: string; createdAt: string }[];
  layout: Record<string, { x: number; y: number }>;
}

export interface ProviderSettings {
  provider: 'ollama' | 'openai';
  baseUrl: string;
  model: string;
  apiKey?: string;
  hasApiKey?: boolean;
  mcpUrl: string;
  searchTool: string;
  retrieveTool: string;
  maxRequests: number;
  concurrency: number;
}

export interface ResearchRun {
  id: string;
  subject: string;
  subjectId: string;
  status: 'pending' | 'running' | 'paused' | 'cancelled' | 'failed' | 'completed';
  phase: string;
  question: string;
  createdAt: string;
  updatedAt: string;
  model: string;
  error?: string;
  requests: number;
  sources: number;
  claims: number;
  conflicts: number;
  plan: { question: string; interval: string; status: string }[];
  activities: { time: string; message: string; tool?: string }[];
  project?: Project;
}
