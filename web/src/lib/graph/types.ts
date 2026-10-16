import type { Edge, Node } from '@xyflow/svelte';
import type { HistoricalEvent } from '../types/knowledge';
import type { Resolution } from '../types/temporal';

export type CardData = Record<string, unknown> & {
  title: string;
  detail: string;
  year: string;
  kind: string;
  color: string;
  recordId: string;
  recordType: string;
  eventIds: string[];
  evidenceCount: number;
  portrait?: string;
  count?: number;
  shadow?: boolean;
};

export interface EventGroup {
  id: string;
  title: string;
  events: HistoricalEvent[];
  start: number;
  end: number;
}

export interface ProjectionOptions {
  chapterId?: string;
  resolution: Resolution;
  from: number | null;
  to: number | null;
  relations: boolean;
  conflicts: boolean;
  conflictMode: string;
  labels: boolean;
  selectedId: string;
  activeClaims: Record<string, string>;
  view: string;
}

export interface Projection {
  nodes: Node<CardData>[];
  edges: Edge[];
  hidden: number;
}
