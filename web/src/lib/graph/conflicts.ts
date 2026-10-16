import type { Edge, Node } from '@xyflow/svelte';
import type { Project } from '../types/project';
import type { CardData, ProjectionOptions } from './types';
import { knownYear, intervalLabel } from '../types/temporal';

export function conflictProjection(
  project: Project,
  nodes: Node<CardData>[],
  edges: Edge[],
  options: ProjectionOptions,
) {
  for (const conflict of project.conflicts) {
    const claims = project.claims.filter((claim) => conflict.claimIds.includes(claim.id));
    const linked = nodes.filter((node) =>
      claims.some((claim) => claim.eventId && node.data.eventIds.includes(claim.eventId)),
    );
    if (!linked.length) continue;
    const active = options.activeClaims[conflict.id];
    for (const [index, claim] of claims.entries()) {
      const node = linked.find((node) => node.data.eventIds.includes(claim.eventId ?? '')) ?? linked[0];
      const parallel = options.conflictMode === 'Parallel';
      const id = `hypothesis:${claim.id}`;
      const canonicalYear = knownYear(project.events.find((event) => event.id === claim.eventId)?.time.start);
      const alternativeYear = knownYear(claim.time?.start) ?? canonicalYear;
      let x = node.position.x + 25 + index * 20;
      if (alternativeYear !== undefined && claim.eventId) {
        const eventYear = Number(node.data.anchorYear);
        const known = project.events
          .map((event) => event.time.start?.year)
          .filter((year): year is number => year !== undefined);
        if (eventYear !== undefined)
          x += ((alternativeYear - eventYear) * 1080) / Math.max(1, Math.max(...known) - Math.min(...known));
      }
      nodes.push({
        id,
        type: 'chronicle',
        position: { x, y: parallel ? 720 + index * 105 : node.position.y - 95 - index * 12 },
        class: active === claim.id ? 'active-hypothesis' : 'shadow-hypothesis',
        data: {
          title: conflict.title + ` · ${index + 1}`,
          detail: conflict.explanation,
          kind: `Hypothesis ${index + 1}`,
          year: claim.time ? intervalLabel(claim.time) : node.data.year,
          color: 'red',
          recordId: conflict.id,
          recordType: 'conflict',
          eventIds: [],
          evidenceCount: claim.evidenceIds.length,
          shadow: active !== claim.id,
        },
      });
      edges.push({
        id: `branch:${claim.id}`,
        source: node.id,
        target: id,
        type: 'smoothstep',
        class: 'conflict-edge',
        label: options.labels ? 'competing interpretation' : undefined,
      });
    }
  }
}
