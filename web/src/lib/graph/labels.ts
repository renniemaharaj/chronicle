import { dateLabel } from '../types/temporal';

export function groupDateLabel(start: number, end: number): string {
  const first = dateLabel({
    year: start,
    precision: 'year',
    calendar: 'historical',
  });
  const last = dateLabel({
    year: end,
    precision: 'year',
    calendar: 'historical',
  });

  return `${first}–${last}`;
}
