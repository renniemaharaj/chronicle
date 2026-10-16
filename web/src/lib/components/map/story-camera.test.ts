import { describe, expect, it, vi } from 'vitest';
import type * as Leaflet from 'leaflet';
import type { MapCameraFrame } from '../../utils/story-geography';
import type { MapPlace } from '../../utils/map-projection';
import { storyCoordinates, moveStoryCamera } from './story-camera';

function frame(longitudes: number[], kind = 'city', overview = false): MapCameraFrame {
  return {
    key: JSON.stringify(longitudes),
    overview,
    places: longitudes.map((longitude) => ({ location: { latitude: 10, longitude, kind } }) as MapPlace),
  };
}

function camera() {
  const map = {
    stop: vi.fn(),
    getSize: () => ({ x: 1200, y: 700 }),
    getCenter: () => ({ lng: 0 }),
    setView: vi.fn(),
    fitBounds: vi.fn(),
    flyToBounds: vi.fn(),
  };
  const library = { latLngBounds: vi.fn((points) => points) };

  return {
    map,
    move: (context: MapCameraFrame, animate: boolean) =>
      moveStoryCamera(library as unknown as typeof Leaflet, map as unknown as Leaflet.Map, context, animate),
  };
}

describe('story map camera', () => {
  it('keeps western hemisphere coordinates aligned with their markers', () => {
    expect(storyCoordinates(frame([-62]))).toEqual([[10, -62]]);
    expect(storyCoordinates(frame([]))).toEqual([]);
  });

  it('frames date-line locations across the shorter span', () => {
    const points = storyCoordinates(frame([179, -179]));
    expect(Math.abs(points[0][1] - points[1][1])).toBe(2);
    expect(storyCoordinates(frame([-179]), 179)).toEqual([[10, 181]]);
    expect(storyCoordinates(frame([179]), -179)).toEqual([[10, -181]]);
  });

  it('uses broad zoom levels for countries and chapter overviews', () => {
    const country = camera();
    country.move(frame([-62], 'country'), false);
    expect(country.map.fitBounds.mock.calls[0][1].maxZoom).toBe(4);
    const overview = camera();
    overview.move(frame([-62], 'city', true), false);
    expect(overview.map.fitBounds.mock.calls[0][1].maxZoom).toBe(5);
  });

  it('animates transitions and reserves space for the reading panel', () => {
    const current = camera();
    current.move(frame([-62]), true);
    expect(current.map.stop).toHaveBeenCalledOnce();
    expect(current.map.flyToBounds).toHaveBeenCalledOnce();
    expect(current.map.fitBounds).not.toHaveBeenCalled();
    expect(current.map.flyToBounds.mock.calls[0][1].paddingBottomRight[0]).toBeGreaterThan(460);
  });

  it('uses the world overview when no mapped context exists', () => {
    const current = camera();
    current.move(frame([], 'city', true), false);
    expect(current.map.setView).toHaveBeenCalledWith([18, -25], 2, { animate: false });
  });
});
