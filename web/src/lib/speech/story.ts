import type { Project, HistoricalContext } from '../types/project';
import type { HistoricalEvent } from '../types/knowledge';
import { knownYear } from '../types/temporal';

export function contextFor(project: Project, event?: HistoricalEvent): HistoricalContext[] {
  const year = knownYear(event?.time.start);
  if (year === undefined) {
    return [];
  }
  return (project.contextPeriods ?? []).filter((context) => {
    const start = knownYear(context.time.start);
    const end = knownYear(context.time.end);
    return (start === undefined || year >= start) && (end === undefined || year <= end);
  });
}

export function contextText(project: Project, context: HistoricalContext): string {
  const claims = context.claimIds.map((id) => project.claims.find((claim) => claim.id === id));
  if (
    !claims.length ||
    claims.some((claim) => !claim || ['unverified', 'rejected_superseded'].includes(claim.status))
  ) {
    return 'This historical context still needs verified supporting evidence.';
  }
  const caution = claims.some((claim) => claim?.status === 'disputed') ? 'Sources disagree about this context. ' : '';
  return `${caution}${context.interpretation ? 'Historical interpretation. ' : ''}${context.summary}`;
}

export function eventStory(project: Project, event: HistoricalEvent | undefined, activeClaims: Record<string, string>) {
  if (!event) {
    return 'Select an event to read its story.';
  }
  const alternative = project.claims.find(
    (claim) => claim.eventId === event.id && Object.values(activeClaims).includes(claim.id),
  );
  if (alternative) {
    const qualification =
      alternative.status === 'unverified'
        ? 'Unverified assertion. '
        : alternative.status === 'rejected_superseded'
          ? 'Rejected or superseded assertion. '
          : alternative.status === 'disputed'
            ? 'Sources disagree. '
            : '';
    return `Selected interpretation. ${qualification}${alternative.text}`;
  }
  const claims = (event.narration?.claimIds ?? event.claimIds).map((id) =>
    project.claims.find((claim) => claim.id === id),
  );
  const assessed = claims.filter((claim) => claim && !['unverified', 'rejected_superseded'].includes(claim.status));
  if (!assessed.length) {
    return `${event.title}. Supporting evidence has not yet been verified.`;
  }
  const caution = assessed.some((claim) => claim?.status === 'disputed') ? 'Sources disagree about part of this account. ' : '';
  const curated = claims.length === assessed.length ? event.narration?.text : '';
  return caution + (curated || assessed.map((claim) => claim?.text).join(' '));
}
