import type { Entity, HistoricalEvent } from '../types/knowledge';
import type { Project } from '../types/project';
import type { HistoricalDate } from '../types/temporal';
import { knownYear } from '../types/temporal';
import { mapLocation } from './geography';

export interface MapPlace {
  entity: Entity;
  location: NonNullable<ReturnType<typeof mapLocation>>;
  events: HistoricalEvent[];
}

export function mapEvents(project: Project, chapter: string, from: number | null, to: number | null) {
  const ids = project.chapters.find((item) => item.id === chapter)?.eventIds;
  return project.events.filter((event) => {
    const start = knownYear(event.time.start);
    const end = knownYear(event.time.end) ?? start;
    if (ids && !ids.includes(event.id)) {
      return false;
    }
    if (from !== null && (end === undefined || end < from)) {
      return false;
    }

    return to === null || (start !== undefined && start <= to);
  });
}

export function mapPlaces(project: Project, events: HistoricalEvent[], scoped: boolean): MapPlace[] {
  return project.entities.flatMap((entity) => {
    if (entity.kind !== 'place') {
      return [];
    }
    const location = mapLocation(entity);
    const linked = events.filter((event) => event.locationIds.includes(entity.id));
    if (!location || (scoped && !linked.length)) {
      return [];
    }

    return [{ entity, location, events: linked }];
  });
}

// Compare only supported calendars; broad dates retain their full possible range.
export function dateExtent(date?: HistoricalDate): [number, number] | undefined {
  if (knownYear(date) === undefined || !date || !['gregorian', 'proleptic_gregorian'].includes(date.calendar)) {
    return;
  }
  const spread = Math.max(0, date.uncertaintyYears ?? (date.approximate ? 1 : 0));
  const width = date.precision === 'century' ? 99 : date.precision === 'decade' ? 9 : 0;
  const firstMonth = date.month ?? 1;
  const lastMonth = date.month ?? 12;
  return [
    (date.year - spread) * 372 + (firstMonth - 1) * 31 + (date.day ?? 1),
    (date.year + spread + width) * 372 + (lastMonth - 1) * 31 + (date.day ?? 31),
  ];
}

export function visitExtent(event: HistoricalEvent): [number, number] | undefined {
  const start = dateExtent(event.time.start);
  const end = event.time.end ? dateExtent(event.time.end) : start;
  if (start && end) {
    return [start[0], end[1]];
  }
}

export function recordedVisits(project: Project, events: HistoricalEvent[]) {
  const subject = project.entities.find((entity) => entity.id === project.subjectId);
  const death = dateExtent(subject?.time.end);
  return events.filter((event) => {
    const extent = visitExtent(event);
    const supported = event.claimIds.some((id) => {
      const claim = project.claims.find((item) => item.id === id);
      return claim && ['supported_limited', 'corroborated'].includes(claim.status) &&
        project.evidence.some((item) => item.claimId === id && item.verified &&
          ['supporting', 'partially_supporting'].includes(item.relation));
    });

    return event.requiresPresence === true && event.participants.includes(project.subjectId) &&
      event.locationIds.length > 0 && supported && extent && (!death || extent[1] <= death[1]);
  }).sort((a, b) => visitExtent(a)![0] - visitExtent(b)![0]);
}
