import type * as Leaflet from 'leaflet';
import type { PublicConfig } from '../../api/config';

export function createMap(
  library: typeof Leaflet,
  container: HTMLDivElement,
  settings: PublicConfig,
  onTileError: () => void,
  story = false,
) {
  const map = library.map(container, {
    center: [18, -25],
    zoom: 2,
    minZoom: 1,
    maxZoom: 19,
    zoomSnap: story ? 0 : 1,
  });

  library
    .tileLayer(settings.mapTileUrl, {
      maxZoom: 19,
      attribution: settings.mapAttribution,
    })
    .on('tileerror', onTileError)
    .addTo(map);

  const layer = library.layerGroup().addTo(map);
  const observer = new ResizeObserver(() => map.invalidateSize());
  observer.observe(container);

  return { map, layer, observer };
}
