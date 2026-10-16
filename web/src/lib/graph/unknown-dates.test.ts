import { expect, it } from 'vitest';
import raw from '../../../../internal/fixtures/person.json';
import type { Project } from '../types/project';
import { aggregate } from './aggregation';
import { orderingImplications, interpretedEvents } from './interpretations';
import family from '../../../../internal/fixtures/family.json';

it('keeps unknown historical dates outside the dated axis', () => {
  const project = structuredClone(raw) as Project;
  project.events[0].time.start = { year: 0, precision: 'unknown', calendar: 'historical' };
  const groups = aggregate(project.events, 'Year');
  const birth = groups.find((group) => group.id === 'birth');
  expect(Number.isNaN(birth?.start)).toBe(true);
});

it('switches interpretations without mutating the canonical chronology', () => {
  const project = structuredClone(family) as Project;
  const event = project.events.find((event) => event.id === 'family-2')!;
  const alternatives = interpretedEvents(project, { 'conflict-family': 'alternative-family' });
  expect(alternatives.find((candidate) => candidate.id === event.id)?.time.start?.year).toBe(1830);
  expect(event.time.start?.year).toBe(1880);
  expect(orderingImplications(project, 'alternative-family').some((text) => text.includes('before Expansion'))).toBe(
    true,
  );
});
