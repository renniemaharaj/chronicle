import type { Entity } from '../types/knowledge';

export interface MapLocation {
  latitude: number;
  longitude: number;
  label: string;
  sourceURL: string;
  kind: string;
  provider: string;
  retrievedAt: string;
}

export type PlaceResult = Omit<MapLocation, 'provider' | 'retrievedAt'>;

export function mapLocation(entity: Entity): MapLocation | undefined {
  const location = entity.attributes?.mapLocation as MapLocation | undefined;
  if (!location || !Number.isFinite(location.latitude) || !Number.isFinite(location.longitude)) {
    return;
  }
  if (Math.abs(location.latitude) > 90 || Math.abs(location.longitude) > 180) {
    return;
  }

  return {
    latitude: location.latitude,
    longitude: location.longitude,
    label: typeof location.label === 'string' ? location.label : entity.name,
    kind: typeof location.kind === 'string' ? location.kind : 'place',
    sourceURL: typeof location.sourceURL === 'string' ? location.sourceURL : '',
    provider: typeof location.provider === 'string' ? location.provider : '',
    retrievedAt: typeof location.retrievedAt === 'string' ? location.retrievedAt : '',
  };
}

export function openMapURL(entity: Entity) {
  const location = mapLocation(entity);
  if (location) {
    return `https://www.openstreetmap.org/?mlat=${location.latitude}&mlon=${location.longitude}#map=12/${location.latitude}/${location.longitude}`;
  }

  return `https://www.openstreetmap.org/search?query=${encodeURIComponent(entity.name)}`;
}

export function locationSourceURL(location: MapLocation) {
  try {
    const url = new URL(location.sourceURL);
    if (url.protocol === 'https:' && url.hostname === 'www.openstreetmap.org') {
      return url.href;
    }
  } catch {
    return '';
  }

  return '';
}
