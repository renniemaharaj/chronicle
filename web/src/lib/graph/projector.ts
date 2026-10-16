import { portraitFor } from '../utils/subjects';
import type { Node } from '@xyflow/svelte';
import type { Project } from '../types/project';
import { knownYear, dateLabel, intervalLabel } from '../types/temporal';
import { interpretedEvents } from './interpretations';
import { projectDetails } from './details';
import { aggregate } from './aggregation';
import { projectEdges } from './relationships';
import { conflictProjection } from './conflicts';
import type { CardData, Projection, ProjectionOptions } from './types';

export function projectGraph(project: Project, options: ProjectionOptions): Projection {
  let visible = interpretedEvents(project, options.activeClaims).filter((event) => {
    const chapter = project.chapters.find((item) => item.id === options.chapterId);
    if (chapter && !chapter.eventIds.includes(event.id)) return false;
    const start = knownYear(event.time.start);
    const end = knownYear(event.time.end) ?? start;
    if (start === undefined || end === undefined) return true;
    return (options.from === null || end >= options.from) && (options.to === null || start <= options.to);
  });
  if (options.view === 'Family & Relationships') visible = [];
  if (options.view === 'Works & Ideas')
    visible = visible.filter((event) => ['Publication', 'Idea'].includes(event.kind));
  if (options.view === 'Places') visible = visible.filter((event) => event.locationIds.length > 0);
  if (options.conflictMode === 'Isolation') {
    const claims = project.conflicts.flatMap((conflict) => conflict.claimIds);
    visible = visible.filter((event) => event.claimIds.some((id) => claims.includes(id)));
  }
  const groups = aggregate(visible, options.resolution);
  const hidden = Math.max(0, groups.length - 160);
  const bounded = groups
    .sort(
      (a, b) =>
        Math.max(...b.events.map((event) => event.importance)) - Math.max(...a.events.map((event) => event.importance)),
    )
    .slice(0, 160);
  bounded.sort((a, b) => a.start - b.start);
  const years = bounded.map((group) => group.start).filter(Number.isFinite);
  const min = options.from ?? Math.min(...years, 0);
  const origin = years.length ? Math.min(...years) : min;
  const max = options.to ?? Math.max(...years, origin + 1);
  const scale = 1080 / Math.max(1, max - origin);
  const laneEnds: number[] = [];
  const nodes: Node<CardData>[] = bounded.map((group, index) => {
    const event = group.events[0];
    const x = Number.isFinite(group.start) ? 50 + (group.start - origin) * scale : 1250;
    let row = laneEnds.findIndex((end) => end + 12 < x);
    if (row < 0) row = laneEnds.length;
    laneEnds[row] = x + 240;
    const isGroup = group.events.length > 1;
    const data: CardData = {
      title: group.title,
      anchorYear: group.start,
      detail: isGroup ? `${group.events.length} documented events` : event.summary,
      year: isGroup
        ? `${dateLabel({ year: group.start, precision: 'year', calendar: 'historical' })}–${dateLabel({ year: group.end, precision: 'year', calendar: 'historical' })}`
        : intervalLabel(event.time),
      kind: isGroup ? 'Historical phase' : event.kind,
      color: colorFor(event.kind),
      recordId: event.id,
      recordType: 'event',
      eventIds: group.events.map((item) => item.id),
      evidenceCount: project.evidence.filter((evidence) =>
        group.events.some((item) => item.claimIds.includes(evidence.claimId)),
      ).length,
      count: group.events.length,
    };
    return {
      id: group.id,
      type: 'chronicle',
      data,
      position: project.layout[group.id] ?? { x, y: 40 + row * 145 },
      selected: data.eventIds.includes(options.selectedId),
    };
  });
  if (options.relations) {
    const relatedIds = options.view === "Family & Relationships" ? project.relationships.flatMap((relation) => [relation.fromId, relation.toId]) : [];
    const participants = new Set([...relatedIds, ...visible.flatMap((event) => [...event.participants, ...event.locationIds])]);
    const entities = project.entities.filter((entity) => participants.has(entity.id)).slice(0, 18);
    entities.forEach((entity, index) =>
      nodes.push({
        id: `entity:${entity.id}`,
        type: 'chronicle',
        position: project.layout[`entity:${entity.id}`] ?? {
          x: 10 + index * 145,
          y: 185 + laneEnds.length * 145 + (index % 2) * 145,
        },
        selected: options.selectedId === entity.id,
        data: {
          title: entity.name,
          kind: entity.kind,
          year: intervalLabel(entity.time),
          detail: entity.summary,
          color: entity.kind === 'place' ? 'teal' : 'blue',
          recordId: entity.id,
          recordType: 'entity',
          eventIds: [],
          evidenceCount: 0,
          portrait: portraitFor(entity, project).url ?? '',
        },
      }),
    );
  }
  const edges = projectEdges(project, nodes, options.labels, options.relations);
  projectDetails(project, nodes, edges, options);
  if (options.conflicts) conflictProjection(project, nodes, edges, options);
  return { nodes, edges, hidden };
}

function colorFor(kind: string): string {
  if (['Publication', 'Idea'].includes(kind)) return 'gold';
  if (['Education', 'Collaboration'].includes(kind)) return 'teal';
  if (kind === 'Appointment') return 'violet';
  if (kind === 'Dispute') return 'red';
  return 'blue';
}
