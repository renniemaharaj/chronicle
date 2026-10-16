export type Precision = 'exact' | 'day' | 'month' | 'year' | 'decade' | 'century' | 'unknown';
export type Resolution = 'Millennium' | 'Century' | 'Decade' | 'Year' | 'Quarter' | 'Month' | 'Week' | 'Day' | 'Event';

export interface HistoricalDate {
  year: number;
  month?: number;
  day?: number;
  precision: Precision;
  calendar: string;
  approximate?: boolean;
  uncertaintyYears?: number;
  label?: string;
  timeOfDay?: string;
}

export interface HistoricalInterval {
  openEnd?: boolean;
  start?: HistoricalDate;
  end?: HistoricalDate;
}

export function dateLabel(date?: HistoricalDate): string {
  if (!date || date.precision === 'unknown') {
    return 'Unknown';
  }

  const year = date.year <= 0 ? `${1 - date.year} BCE` : String(date.year);
  const prefix = date.approximate ? 'c. ' : '';
  if (date.label) {
    return date.label;
  }
  if (date.precision === 'day' || date.precision === 'exact') {
    const clock = date.precision === 'exact' ? ` ${date.timeOfDay ?? '(time unknown)'}` : '';
    return `${prefix}${date.day} ${months[(date.month ?? 1) - 1]} ${year}${clock}`;
  }
  if (date.precision === 'month') {
    return `${prefix}${months[(date.month ?? 1) - 1]} ${year}`;
  }

  return `${prefix}${year}`;
}

const months = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];

export function intervalLabel(interval: HistoricalInterval): string {
  const start = dateLabel(interval.start);
  if (interval.openEnd) return `${start} onward (end unknown)`;
  if (!interval.end || interval.end.year === interval.start?.year) {
    return start;
  }

  return `${start} – ${dateLabel(interval.end)}`;
}

export function knownYear(value?: HistoricalDate): number | undefined {
  return value && value.precision !== 'unknown' ? value.year : undefined;
}
