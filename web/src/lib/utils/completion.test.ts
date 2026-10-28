import { describe, expect, it } from 'vitest';
import { researchCompletion, storyCompletion } from './completion';
import { emptyProject } from '../state/empty';

describe('bounded preparation progress', () => {
  it('keeps incomplete research below 100 through validation and exhaustion', () => {
    const plan = [
      {
        question: 'Origins',
        interval: '',
        status: 'completed',
      },
      {
        question: 'Aftermath',
        interval: '',
        status: 'pending',
      },
    ];
    expect(
      researchCompletion({
        plan,
        status: 'failed',
      }).percent,
    ).toBe(50);
    plan[1].status = 'completed';
    expect(
      researchCompletion({
        plan,
        status: 'running',
      }).percent,
    ).toBe(99);
    expect(
      researchCompletion({
        plan,
        status: 'completed',
      }).percent,
    ).toBe(100);
    expect(
      researchCompletion({
        plan: [],
        status: 'pending',
      }).percent,
    ).toBe(0);
  });

  it('does not count unsupported or missing narration as ready', () => {
    const project = structuredClone(emptyProject);
    project.events = [
      {
        id: 'one',
        title: 'Synthetic event',
        kind: 'test',
        summary: '',
        time: {},
        participants: [],
        locationIds: [],
        claimIds: ['claim'],
        importance: 1,
        phase: '',
      },
      {
        id: 'two',
        title: 'Missing passage',
        kind: 'test',
        summary: '',
        time: {},
        participants: [],
        locationIds: [],
        claimIds: [],
        importance: 1,
        phase: '',
      },
    ];
    project.claims = [
      {
        id: 'claim',
        text: 'Synthetic claim',
        status: 'supported_limited',
        evidenceIds: ['evidence'],
        assessment: [],
      },
    ];
    project.evidence = [
      {
        id: 'evidence',
        claimId: 'claim',
        sourceId: 'source',
        mode: 'quotation',
        verified: true,
        excerpt: 'Synthetic evidence only',
        location: '',
        relation: 'supporting',
      },
    ];
    project.events[0].narration = {
      text: 'Synthetic reading',
      claimIds: ['claim'],
      interpretation: false,
    };
    expect(storyCompletion(project).narrationPercent).toBe(50);
    project.evidence[0].verified = false;
    expect(storyCompletion(project).narrationPercent).toBe(0);
  });
});
