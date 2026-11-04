import type { Campaign, Metrics, Part } from './types';

export type MetricName = 'views' | 'reactions' | 'comments' | 'shares';

export function observations(part: Part, key: MetricName | 'revenue') {
  return part.metrics
    .filter((metric) => !metric.platform || metric.platform === 'facebook')
    .filter((metric) => typeof metric[key] === 'number')
    .sort((left, right) => Date.parse(left.observedAt) - Date.parse(right.observedAt));
}

export function latestObservation(part: Part, key: MetricName | 'revenue'): Metrics | undefined {
  return observations(part, key).at(-1);
}

export function totalMetric(campaign: Campaign, key: MetricName) {
  const values = campaign.parts
    .map((part) => latestObservation(part, key)?.[key])
    .filter((value): value is number => value !== undefined);

  return {
    value: values.length ? values.reduce((sum, value) => sum + value, 0) : undefined,
    observed: values.length,
  };
}

export function viewTrend(part: Part) {
  const values = observations(part, 'views');
  const latest = values.at(-1);
  const previous = values.at(-2);
  if (
    !latest ||
    !previous ||
    !previous.views ||
    Date.parse(latest.observedAt) === Date.parse(previous.observedAt)
  ) {
    return undefined;
  }

  return Math.round((((latest.views ?? 0) - previous.views) / previous.views) * 100);
}

export function metricLabel(value?: number) {
  return value === undefined
    ? '—'
    : new Intl.NumberFormat('en', {
        notation: 'compact',
        maximumFractionDigits: 1,
      }).format(value);
}

export function revenueTotals(campaign: Campaign) {
  const totals = new Map<string, number>();
  for (const part of campaign.parts) {
    const metric = latestObservation(part, 'revenue');
    if (metric?.revenue !== undefined && metric.currency) {
      totals.set(metric.currency, (totals.get(metric.currency) ?? 0) + metric.revenue);
    }
  }

  return [...totals].map(([currency, value]) => ({ currency, value }));
}
