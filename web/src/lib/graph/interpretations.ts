import type { Project } from '../types/project';

export function interpretedEvents(project: Project, active: Record<string, string>) {
  const selected = Object.values(active);
  return project.events.map((event) => {
    const claim = project.claims.find(
      (claim) => selected.includes(claim.id) && claim.eventId === event.id && claim.time,
    );
    return claim?.time ? { ...event, time: claim.time } : event;
  });
}

export function orderingImplications(project: Project, claimId: string): string[] {
  const claim = project.claims.find((claim) => claim.id === claimId);
  const original = project.events.find((event) => event.id === claim?.eventId);
  const proposed = claim?.time?.start?.year;
  if (!original?.time.start || proposed === undefined) return [];
  const year = original.time.start.year;
  const changes = project.events.filter((event) => {
    const other = event.time.start?.year;
    return other !== undefined && event.id !== original.id && (year - other) * (proposed - other) < 0;
  });
  return changes.map(
    (event) =>
      `${original.title} moves ${proposed < year ? 'before' : 'after'} ${event.title} (${event.time.start?.year}). This changes chronology; causation still requires evidence.`,
  );
}
