import { describe, expect, it } from 'vitest';
import { latestObservation, metricLabel, revenueTotals, totalMetric, viewTrend } from './observations';
import type { Campaign, Metrics, Part } from './types';

function part(metrics: Metrics[]): Part {
  return {
    number: 1,
    status: 'published',
    passages: [],
    metrics,
  };
}

function metric(day: number, values: Partial<Metrics>): Metrics {
  return {
    observedAt: `2026-10-${String(day).padStart(2, '0')}T12:00:00Z`,
    method: 'manual',
    ...values,
  };
}

describe('campaign performance observations', () => {
  it('uses the latest observation for each metric, preserving observed zero and unknown', () => {
    const first = part([
      metric(3, {
        views: 120,
        reactions: 4,
      }),
      metric(1, { views: 100 }),
      metric(4, { comments: 0 }),
    ]);
    const campaign = { parts: [first, part([])] } as Campaign;
    expect(latestObservation(first, 'views')?.views).toBe(120);
    expect(totalMetric(campaign, 'views')).toEqual({
      value: 120,
      observed: 1,
    });
    expect(totalMetric(campaign, 'comments')).toEqual({
      value: 0,
      observed: 1,
    });
    expect(totalMetric(campaign, 'shares')).toEqual({
      value: undefined,
      observed: 0,
    });
    expect(metricLabel(0)).toBe('0');
    expect(metricLabel(undefined)).toBe('—');
  });

  it('compares chronological observations without inventing growth from zero or duplicate times', () => {
    expect(viewTrend(part([metric(3, { views: 120 }), metric(1, { views: 100 })]))).toBe(20);
    expect(viewTrend(part([metric(1, { views: 100 }), metric(3, { views: 80 })]))).toBe(-20);
    expect(viewTrend(part([metric(1, { views: 0 }), metric(3, { views: 80 })]))).toBeUndefined();
    expect(viewTrend(part([metric(1, { views: 10 }), metric(1, { views: 80 })]))).toBeUndefined();
  });

  it('sums the latest revenue per episode and keeps currencies separate', () => {
    const campaign = {
      parts: [
        part([
          metric(1, {
            revenue: 1,
            currency: 'USD',
          }),
          metric(3, {
            revenue: 2,
            currency: 'USD',
          }),
        ]),
        part([
          metric(2, {
            revenue: 5,
            currency: 'TTD',
          }),
        ]),
        part([]),
      ],
    } as Campaign;
    expect(revenueTotals(campaign)).toEqual([
      {
        currency: 'USD',
        value: 2,
      },
      {
        currency: 'TTD',
        value: 5,
      },
    ]);
  });
});
