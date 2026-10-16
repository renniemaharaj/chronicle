import { describe, expect, it } from 'vitest';
import rawPerson from '../../../../internal/fixtures/person.json';
import rawFamily from '../../../../internal/fixtures/family.json';
import type { Project } from '../types/project';
import type { ProjectionOptions } from './types';
import { projectGraph } from './projector';

const options: ProjectionOptions = {
  resolution: 'Year',
  from: null,
  to: null,
  relations: true,
  conflicts: true,
  conflictMode: 'Shadows',
  labels: true,
  selectedId: 'birth',
  activeClaims: {},
  view: 'Interactive Graph',
};

describe('temporal graph projections', () => {
  it('changes canonical projection with semantic resolution', () => {
    const project = rawPerson as Project;
    const year = projectGraph(project, options);
    const century = projectGraph(project, { ...options, resolution: 'Century' });
    const detail = projectGraph(project, { ...options, resolution: 'Event' });
    expect(century.nodes.length).toBeLessThan(year.nodes.length);
    expect(century.nodes.some((node) => node.data.kind === 'Historical phase')).toBe(true);
    expect(detail.nodes.length).toBeGreaterThan(year.nodes.length);
    expect(detail.nodes.some((node) => node.data.kind === 'Evidence')).toBe(true);
    expect(project.events.length).toBe(rawPerson.events.length);
  });

  it('uses date windows and bounded node budgets', () => {
    const project = structuredClone(rawPerson) as Project;
    const period = projectGraph(project, { ...options, relations: false, conflicts: false, from: 1840, to: 1880 });
    expect(
      period.nodes.every((node) =>
        node.data.eventIds.every(
          (id) => (project.events.find((event) => event.id === id)?.time.start?.year ?? 0) >= 1840,
        ),
      ),
    ).toBe(true);
    project.events = Array.from({ length: 2000 }, (_, index) => ({
      ...project.events[0],
      id: `event-${index}`,
      time: { start: { year: index, precision: 'year', calendar: 'julian' } },
    }));
    const bounded = projectGraph(project, { ...options, relations: false, conflicts: false });
    expect(bounded.nodes.length).toBe(160);
    expect(bounded.hidden).toBe(1840);
  });

  it('positions competing dates independently and preserves claims', () => {
    const project = rawFamily as Project;
    const shadow = projectGraph(project, options);
    const parallel = projectGraph(project, { ...options, conflictMode: 'Parallel' });
    const alternatives = shadow.nodes.filter((node) => node.id.startsWith('hypothesis:'));
    expect(alternatives).toHaveLength(2);
    expect(alternatives[0].position.x).not.toBe(alternatives[1].position.x);
    expect(parallel.nodes.find((node) => node.id === alternatives[0].id)?.position.y).toBeGreaterThan(
      alternatives[0].position.y,
    );
    expect(project.claims).toHaveLength(5);
  });
});
