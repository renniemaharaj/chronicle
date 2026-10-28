import { untrack } from 'svelte';
import type * as Leaflet from 'leaflet';
import type { MapPlace } from '../../utils/map-projection';
import { locationZoom } from '../../utils/geography';

export function focusSelectedPlace(map: Leaflet.Map, id: string, places: MapPlace[]) {
  const place = places.find((item) => item.entity.id === id || item.events.some((event) => event.id === id));
  if (!place) {
    return;
  }

  const position: Leaflet.LatLngTuple = [place.location.latitude, place.location.longitude];
  if (place.entity.id === id) {
    map.setView(position, locationZoom(place.location.kind));
  } else if (!untrack(() => map.getBounds().contains(position))) {
    map.panTo(position);
  }
}
