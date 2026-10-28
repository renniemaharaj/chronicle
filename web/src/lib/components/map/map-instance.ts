import type * as Leaflet from 'leaflet';

export function createMap(
  library: typeof Leaflet,
  container: HTMLDivElement,
  settings: { mapTileUrl: string; mapAttribution: string },
  onTileError: () => void,
) {
  const map = library.map(container, {
    center: [18, -25],
    zoom: 2,
    minZoom: 1,
    maxZoom: 19,
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
