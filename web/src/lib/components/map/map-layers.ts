import type * as Leaflet from 'leaflet';
import type { HistoricalEvent } from '../../types/knowledge';
import type { MapPlace } from '../../utils/map-projection';
import { visitExtent } from '../../utils/map-projection';
import { intervalLabel } from '../../types/temporal';

export function drawMap(
  L: typeof Leaflet, layer: Leaflet.LayerGroup, places: MapPlace[], visits: HistoricalEvent[],
  selected: string, select: (id: string) => void,
) {
  layer.clearLayers();
  for (let index = 1; index < visits.length; index++) {
    const previous = visits[index - 1];
    const current = visits[index];
    const from = places.filter((place) => previous.locationIds.includes(place.entity.id));
    const to = places.filter((place) => current.locationIds.includes(place.entity.id));
    const end = visitExtent(previous)?.[1];
    const start = visitExtent(current)?.[0];
    // Multiple locations and overlapping date ranges do not establish an order.
    if (from.length !== 1 || to.length !== 1 || end === undefined || start === undefined || end >= start) {
      continue;
    }
    const a = from[0].location;
    const b = to[0].location;
    if (from[0].entity.id === to[0].entity.id || Math.abs(a.longitude - b.longitude) > 180) {
      continue;
    }
    const label = document.createElement('span');
    label.textContent = `${intervalLabel(previous.time)} → ${intervalLabel(current.time)} · Recorded sequence; route unknown`;
    L.polyline([[a.latitude, a.longitude], [b.latitude, b.longitude]], {
      color: '#2459b9', weight: 2, dashArray: '6 7', opacity: 0.75,
    }).bindTooltip(label).addTo(layer);
  }

  for (const place of places) {
    const active = selected === place.entity.id || place.events.some((event) => event.id === selected);
    const label = document.createElement('span');
    label.textContent = `${place.entity.name} · ${place.events.length} timeline events`;
    L.circleMarker([place.location.latitude, place.location.longitude], {
      radius: active ? 9 : 6, color: active ? '#173969' : '#ffffff', weight: 2,
      fillColor: active ? '#e3ac32' : '#2865d1', fillOpacity: 1,
    }).bindTooltip(label, { direction: 'top' }).on('click', () => select(place.entity.id)).addTo(layer);
  }
}
