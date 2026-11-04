import type * as Leaflet from 'leaflet';
import type { MapCameraFrame } from '../../utils/story-geography';
import { storyCoordinates } from './story-camera';

export function cameraMotion(seconds: number) {
  const time = Number.isFinite(seconds) ? Math.max(0, seconds) : 0;

  return {
    x: 8 * Math.sin(time / 14),
    y: 5 * Math.sin(time / 19),
    zoom: 0.07 * Math.sin(time / 21),
  };
}

// Move the existing pane without resetting Leaflet's view or replacing map tiles.
export function applyCameraMotion(map: Leaflet.Map, seconds: number) {
  const pane = map.getPane('mapPane');
  if (!pane) {
    return;
  }

  const offset = cameraMotion(seconds);
  const size = map.getSize();
  pane.style.transformOrigin = `${size.x / 2}px ${size.y / 2}px`;
  pane.style.translate = `${(offset.x * size.x) / 1080}px ${(offset.y * size.y) / 1920}px`;
  pane.style.scale = String(1.025 + offset.zoom * 0.1);
}

export function clearCameraMotion(map: Leaflet.Map) {
  const pane = map.getPane('mapPane');
  pane?.style.removeProperty('translate');
  pane?.style.removeProperty('scale');
  pane?.style.removeProperty('transform-origin');
}

export function locationFocus(L: typeof Leaflet, map: Leaflet.Map, frame: MapCameraFrame) {
  const layer = L.layerGroup().addTo(map);
  const circles = frame.overview
    ? []
    : storyCoordinates(frame, map.getCenter().lng).map((point) => {
        const glow = L.circleMarker(point, {
          interactive: false,
          className: 'story-focus-glow',
          color: '#dfb45d',
          weight: 1,
          fillColor: '#dfb45d',
          fillOpacity: 0.09,
          opacity: 0.12,
        }).addTo(layer);
        const outline = L.circleMarker(point, {
          interactive: false,
          className: 'story-focus-outline',
          color: '#f2ba53',
          weight: 2.5,
          fill: false,
          dashArray: '12 10',
          opacity: 0.8,
        }).addTo(layer);

        return { glow, outline };
      });

  return {
    update(seconds: number, animate: boolean) {
      const phase = animate ? Math.sin((seconds * Math.PI) / 2) : 0;
      const radius = Math.max(40, Math.min(112, map.getSize().x * 0.1));
      for (const { glow, outline } of circles) {
        glow.setRadius(radius * (1.2 + phase * 0.1));
        outline.setRadius(radius * (1 + phase * 0.08));
        outline.setStyle({
          dashOffset: String(animate ? -seconds * 8 : 0),
          opacity: 0.65 + phase * 0.12,
        });
      }
    },
    remove() {
      layer.remove();
    },
  };
}
