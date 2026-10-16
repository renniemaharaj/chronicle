import type * as Leaflet from 'leaflet';
import type { MapCameraFrame } from '../../utils/story-geography';
import { locationZoom } from '../../utils/geography';

// Use the shortest longitude span when the context crosses the date line.
export function storyCoordinates(frame: MapCameraFrame, referenceLongitude = 0): Leaflet.LatLngTuple[] {
  if (!frame.places.length) {
    return [];
  }
  const longitudes = frame.places.map((place) =>
    (place.location.longitude + 360) % 360,
  ).sort((a, b) => a - b);
  let start = longitudes[0];
  let largestGap = -1;
  for (let index = 0; index < longitudes.length; index++) {
    const next = longitudes[(index + 1) % longitudes.length] + (index === longitudes.length - 1 ? 360 : 0);
    const gap = next - longitudes[index];
    if (gap > largestGap) {
      largestGap = gap;
      start = next % 360;
    }
  }

  const positions: Leaflet.LatLngTuple[] = frame.places.map(({ location }) => {
    let longitude = (location.longitude + 360) % 360;
    if (longitude < start) {
      longitude += 360;
    }

    return [location.latitude, longitude];
  });

  const center = (Math.min(...positions.map((point) => point[1])) + Math.max(...positions.map((point) => point[1]))) / 2;
  const shift = Math.round((referenceLongitude - center) / 360) * 360;
  return positions.map(([latitude, longitude]) => [latitude, longitude + shift]);
}

export function moveStoryCamera(L: typeof Leaflet, map: Leaflet.Map, frame: MapCameraFrame, animate: boolean) {
  map.stop();
  const positions = storyCoordinates(frame, map.getCenter().lng);
  if (!positions.length) {
    map.setView([18, -25], 2, { animate });
    return;
  }

  const size = map.getSize();
  const wide = window.matchMedia('(min-width: 760px)').matches;
  const reader = map.getContainer().closest('.story-scene')?.querySelector('.story-reading');
  const readerWidth = reader?.getBoundingClientRect().width ?? Math.min(400, size.x * 0.4);
  const paddingTopLeft: Leaflet.PointTuple = [48, wide ? 76 : 180];
  const paddingBottomRight: Leaflet.PointTuple = wide
    ? [readerWidth + 42, 110]
    : [40, size.y * 0.53];
  const maxZoom = frame.overview ? 5 : Math.min(...frame.places.map((place) => locationZoom(place.location.kind)));
  const options = { paddingTopLeft, paddingBottomRight, maxZoom, animate, duration: 1.2 };
  const bounds = L.latLngBounds(positions);
  if (animate) {
    map.flyToBounds(bounds, options);
  } else {
    map.fitBounds(bounds, options);
  }
}
