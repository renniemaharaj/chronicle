import { describe, expect, it } from 'vitest';
import type { Entity } from '../types/knowledge';
import { contemporary, searchText } from './subjects';

function person(id: string, from?: number, to?: number): Entity {
  return {
    id,
    name: id,
    kind: 'person',
    summary: '',
    version: 1,
    attributes: {},
    time: {
      start: from === undefined ? undefined : { year: from, precision: 'year', calendar: 'historical' },
      end: to === undefined ? undefined : { year: to, precision: 'year', calendar: 'historical' },
    },
  };
}

describe('subject discovery', () => {
  it('finds accented names with ordinary keyboard text', () => {
    expect(searchText('René Angélil')).toContain(searchText('rene angelil'));
  });
  it('separates unknown dates and excludes people outside a historical interval', () => {
    const subject = person('subject', -300, -220);
    expect(contemporary(person('earlier', -500, -400), subject)).toBe('outside');
    expect(contemporary(person('later', -100, -50), subject)).toBe('outside');
    expect(contemporary(person('overlap', -250, -200), subject)).toBe('overlap');
    expect(contemporary(person('unknown'), subject)).toBe('unknown');
  });
});
