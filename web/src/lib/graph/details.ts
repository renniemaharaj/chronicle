import type { Edge, Node } from '@xyflow/svelte';
import type { Project } from '../types/project';
import type { CardData, ProjectionOptions } from './types';

export function projectDetails(project: Project, nodes: Node<CardData>[], edges: Edge[], options: ProjectionOptions) {
  const detailed = ['Quarter', 'Month', 'Week', 'Day', 'Event'];
  if (!detailed.includes(options.resolution)) return;
  const events = [...nodes].filter((node) => node.data.recordType === 'event').slice(0, 12);
  for (const event of events) {
    const claims = project.claims.filter((claim) => event.data.eventIds.includes(claim.eventId ?? '')).slice(0, 3);
    for (const [index, claim] of claims.entries()) {
      const id = `claim:${claim.id}`;
      nodes.push({
        id,
        type: 'chronicle',
        position: { x: event.position.x + 20, y: event.position.y + 135 + index * 100 },
        data: {
          title: claim.text,
          detail: claim.assessment.join(' '),
          year: claim.status.replaceAll('_', ' '),
          kind: 'Claim',
          color: 'violet',
          recordId: claim.eventId ?? event.data.recordId,
          recordType: 'event',
          eventIds: [],
          evidenceCount: claim.evidenceIds.length,
        },
      });
      edges.push({
        id: `assertion:${claim.id}`,
        source: event.id,
        target: id,
        type: 'smoothstep',
        class: 'relationship-edge',
        label: options.labels ? 'assertion' : undefined,
      });
      if (!['Day', 'Event'].includes(options.resolution)) continue;
      for (const [sourceIndex, evidence] of project.evidence
        .filter((item) => item.claimId === claim.id)
        .slice(0, 2)
        .entries()) {
        nodes.push({
          id: `evidence:${evidence.id}`,
          type: 'chronicle',
          position: { x: event.position.x + 200 + sourceIndex * 175, y: event.position.y + 135 + index * 100 },
          data: {
            title: project.sources.find((source) => source.id === evidence.sourceId)?.title ?? 'Source missing',
            detail: evidence.excerpt,
            year: evidence.location,
            kind: 'Evidence',
            color: 'teal',
            recordId: claim.eventId ?? event.data.recordId,
            recordType: 'evidence',
            eventIds: [],
            evidenceCount: 1,
          },
        });
        edges.push({
          id: `provenance:${evidence.id}`,
          source: id,
          target: `evidence:${evidence.id}`,
          type: 'smoothstep',
          class: 'participant-edge',
          label: options.labels ? evidence.relation : undefined,
        });
      }
    }
  }
}
