import type { HistoricalEvent } from '../types/knowledge';
import { knownYear, type Resolution } from '../types/temporal';
import type { EventGroup } from './types';

const bucketSize: Record<Resolution, number> = {
  Millennium: 1000,
  Century: 100,
  Decade: 10,
  Year: 0,
  Quarter: 0,
  Month: 0,
  Week: 0,
  Day: 0,
  Event: 0,
};

export function aggregate(events: HistoricalEvent[], resolution: Resolution): EventGroup[] {
  const width = bucketSize[resolution];
  const groups = new Map<string, HistoricalEvent[]>();
  const ordered = [...events].sort(
    (a, b) => (knownYear(a.time.start) ?? Infinity) - (knownYear(b.time.start) ?? Infinity),
  );

  for (const event of ordered) {
    const year = knownYear(event.time.start);
    const key = width && year !== undefined ? `${event.phase}-${Math.floor(year / width) * width}` : event.id;
    const group = groups.get(key) ?? [];
    group.push(event);
    groups.set(key, group);
  }

  return [...groups.entries()].map(([id, members]) => {
    const years = members.flatMap((event) => [knownYear(event.time.start), knownYear(event.time.end)]);
    const known = years.filter((year): year is number => year !== undefined);
    const start = known.length ? Math.min(...known) : NaN;
    const end = known.length ? Math.max(...known) : NaN;
    return {
      id,
      title: members.length > 1 ? members[0].phase || 'Historical developments' : members[0].title,
      events: members,
      start,
      end,
    };
  });
}
