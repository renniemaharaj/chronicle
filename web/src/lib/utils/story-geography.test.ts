import { describe, expect, it } from 'vitest';
import { emptyProject } from '../state/empty';
import type { Entity, HistoricalEvent } from '../types/knowledge';
import { storyGeography } from './story-geography';

function place(id: string, longitude?: number): Entity {
  return {
    id,
    name: `Synthetic ${id}`,
    kind: 'place',
    summary: 'Test location, not historical evidence.',
    time: {},
    attributes: longitude === undefined ? {} : { mapLocation: { latitude: 10, longitude, kind: 'city' } },
    version: 1,
  };
}

function event(id: string, locationIds: string[]): HistoricalEvent {
  return {
    id,
    title: 'Synthetic event',
    kind: 'test',
    summary: 'Mentions Synthetic elsewhere without placing this event there.',
    time: {},
    participants: [],
    locationIds,
    claimIds: [],
    importance: 1,
    phase: '',
  };
}

function project() {
  return {
    ...structuredClone(emptyProject),
    entities: [place('here', -60), place('there', 40), place('unmapped')],
    events: [event('first', ['here']), event('second', ['there']), event('third', [])],
  };
}

describe('story geography', () => {
  it('uses canonical event locations rather than names mentioned in prose', () => {
    const p = project();
    const context = storyGeography(p, p.events[0], p.events);
    expect(context.frame.places.map((item) => item.entity.id)).toEqual(['here']);
    expect(context.names).toEqual(['Synthetic here']);
    expect(context.frame.overview).toBe(false);
  });

  it('frames multiple locations together and discloses incomplete coordinates', () => {
    const p = project();
    p.events[0].locationIds = ['here', 'there', 'unmapped'];
    const context = storyGeography(p, p.events[0], p.events);
    expect(context.frame.places).toHaveLength(2);
    expect(context.names).toHaveLength(3);
    expect(context.caption).toContain('1 recorded place has no saved map reference');
  });

  it('uses only the active chapter for an unknown-location overview', () => {
    const p = project();
    p.chapters = [{
      id: 'context',
      title: 'Synthetic chapter',
      eventIds: ['first', 'third'],
      summary: '',
      claimIds: [],
      interpretation: true,
    }];
    const context = storyGeography(p, p.events[2], p.events);
    expect(context.frame.overview).toBe(true);
    expect(context.frame.places.map((item) => item.entity.id)).toEqual(['here']);
    expect(context.names).toEqual([]);
    expect(context.caption).toContain('No location recorded');
  });

  it('retains unmapped place names without reusing the previous event location', () => {
    const p = project();
    p.events[2].locationIds = ['unmapped'];
    const context = storyGeography(p, p.events[2], p.events);
    expect(context.frame.overview).toBe(true);
    expect(context.names).toEqual(['Synthetic unmapped']);
    expect(context.caption).toContain('no saved map references');
  });

  it('keeps camera identity stable across passages at the same places', () => {
    const p = project();
    p.events[1].locationIds = ['here'];
    const first = storyGeography(p, p.events[0], p.events);
    const second = storyGeography(p, p.events[1], p.events);
    expect(first.frame.key).toBe(second.frame.key);
    const overview = storyGeography(p, undefined, []);
    expect(overview.frame.places).toEqual([]);
  });
});
