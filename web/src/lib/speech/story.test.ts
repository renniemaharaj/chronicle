import { describe, expect, it } from 'vitest';
import { readFileSync } from 'node:fs';
import { eventStory, contextFor, contextText } from './story';
import { projectGraph } from '../graph/projector';
import type { Project } from '../types/project';

function project(): Project {
  return JSON.parse(readFileSync('../internal/fixtures/person.json', 'utf8'));
}

describe('narration and scoped history', () => {
  it('changes spoken interpretations and marks unsupported events', () => {
    const p = project();
    const event = p.events[0];
    p.claims = [
      {
        id: 'claim',
        eventId: event.id,
        text: 'An alternative date.',
        status: 'disputed',
        evidenceIds: [],
        assessment: [],
      },
    ];
    event.claimIds = ['claim'];
    event.narration = { text: 'A saved event reading.', claimIds: ['claim'], interpretation: false };
    expect(eventStory(p, event, {})).toContain('Sources disagree');
    expect(eventStory(p, event, { conflict: 'claim' })).toContain('An alternative date.');
    p.claims[0].status = 'unverified';
    expect(eventStory(p, event, {})).toContain('not yet been verified');
  });

  it('filters exact chapter membership even when dates overlap', () => {
    const p = project();
    const event = p.events[0];
    p.chapters = [
      {
        id: 'first',
        title: 'Any subject phase',
        summary: '',
        claimIds: [],
        interpretation: true,
        eventIds: [event.id],
      },
    ];
    const result = projectGraph(p, {
      chapterId: 'first',
      resolution: 'Year',
      from: null,
      to: null,
      relations: false,
      conflicts: false,
      conflictMode: 'Shadows',
      labels: true,
      selectedId: event.id,
      activeClaims: {},
      view: 'Interactive Graph',
    });
    expect(result.nodes.flatMap((node) => node.data.eventIds)).toEqual([event.id]);
  });

  it('resolves era intervals without JavaScript dates and guards unverified context', () => {
    const p = project();
    p.events[0].time.start = { year: -300, calendar: 'historical', precision: 'year' };
    p.contextPeriods = [
      {
        id: 'ancient',
        title: 'Scoped ancient context',
        scope: 'A researched region',
        summary: 'A sourced contextual interpretation.',
        claimIds: ['missing'],
        interpretation: true,
        time: {
          start: { year: -400, calendar: 'historical', precision: 'year' },
          end: { year: -200, calendar: 'historical', precision: 'year' },
        },
      },
    ];
    const contexts = contextFor(p, p.events[0]);
    expect(contexts).toHaveLength(1);
    expect(contextText(p, contexts[0])).toContain('needs verified supporting evidence');
  });
});
