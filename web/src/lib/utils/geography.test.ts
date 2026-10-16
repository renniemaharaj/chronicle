import { describe, expect, it } from 'vitest';
import { locationSourceURL, locationZoom, mapLocation } from './geography';
import { emptySubject } from '../state/empty';

describe('sourced geographic references', () => {
  it('keeps coordinate provenance from canonical place attributes', () => {
    const location = mapLocation({
      ...emptySubject,
      attributes: { mapLocation: {
        latitude: 42,
        longitude: 14,
        sourceURL: 'https://example.org/gazetteer',
        sourceId: 'source',
        claimIds: ['coordinate-claim'],
        note: 'Synthetic representative point.',
      } },
    })!;
    expect(location.sourceId).toBe('source');
    expect(location.claimIds).toEqual(['coordinate-claim']);
    expect(locationSourceURL(location)).toBe('https://example.org/gazetteer');
    expect(locationSourceURL({ ...location, sourceURL: 'javascript:alert(1)' })).toBe('');
  });

  it('frames broad features at regional scales', () => {
    expect(locationZoom('sea')).toBeLessThan(locationZoom('region'));
    expect(locationZoom('region')).toBeLessThan(locationZoom('city'));
    expect(locationZoom('province')).toBe(4);
    expect(locationZoom('river')).toBe(6);
  });
});
