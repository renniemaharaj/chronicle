import { knownYear } from '../types/temporal';
import { interpretedEvents } from '../graph/interpretations';
import type { Project } from '../types/project';

export function orderedEvents(project: Project, perspective: string, active: Record<string, string>) {
  let events = interpretedEvents(project, active);
  if (perspective === 'Geographical') events = events.filter((event) => event.locationIds.length);
  if (perspective === 'Institutional') {
    const organizations = project.entities
      .filter((entity) => ['organization', 'institution'].includes(entity.kind))
      .map((entity) => entity.id);
    events = events.filter((event) => event.participants.some((id) => organizations.includes(id)));
  }
  if (perspective === 'Genealogical') {
    const family = project.relationships
      .filter((relation) => /parent|child|spouse|member/.test(relation.kind))
      .flatMap((relation) => [relation.fromId, relation.toId]);
    events = events.filter((event) => event.participants.some((id) => family.includes(id)));
  }
  if (perspective === 'Cause and consequence')
    events = events.filter((event) => event.causeIds?.length || event.consequenceIds?.length);
  if (!events.length) events = [...project.events];

  return events.sort((left, right) => {
    const a = left.time.start;
    const b = right.time.start;
    const years = (knownYear(a) ?? Infinity) - (knownYear(b) ?? Infinity);
    const months = (a?.month ?? 1) - (b?.month ?? 1);
    const days = (a?.day ?? 1) - (b?.day ?? 1);
    return years || months || days || left.id.localeCompare(right.id);
  });
}

export function nearestEvent(project: Project, year: number) {
  return [...project.events].sort(
    (left, right) =>
      Math.abs((knownYear(left.time.start) ?? Infinity) - year) -
      Math.abs((knownYear(right.time.start) ?? Infinity) - year),
  )[0];
}

export function projectBounds(project: Project): [number, number] {
  const years = project.events.flatMap((event) => [knownYear(event.time.start), knownYear(event.time.end)]);
  const known = years.filter((year): year is number => year !== undefined);
  return known.length ? [Math.min(...known), Math.max(...known)] : [0, 1];
}
