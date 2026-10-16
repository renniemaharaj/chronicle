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
