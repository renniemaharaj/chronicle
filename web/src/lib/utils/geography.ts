import type { Entity } from '../types/knowledge';

export interface MapLocation {
  latitude: number;
  longitude: number;
  label: string;
  sourceURL: string;
  kind: string;
  provider: string;
  retrievedAt: string;
  sourceId?: string;
  claimIds?: string[];
  note?: string;
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
    sourceId: typeof location.sourceId === 'string' ? location.sourceId : undefined,
    claimIds: Array.isArray(location.claimIds) ? location.claimIds.filter((id) => typeof id === 'string') : [],
    note: typeof location.note === 'string' ? location.note : '',
  };
}

export function openMapURL(entity: Entity) {
  const location = mapLocation(entity);
  if (location) {
    return `https://www.openstreetmap.org/?mlat=${location.latitude}&mlon=${location.longitude}#map=${locationZoom(location.kind)}/${location.latitude}/${location.longitude}`;
  }

  return `https://www.openstreetmap.org/search?query=${encodeURIComponent(entity.name)}`;
}

export function locationSourceURL(location: MapLocation) {
  try {
    const url = new URL(location.sourceURL);
    if (url.protocol === 'https:') {
      return url.href;
    }
  } catch {
    return '';
  }

  return '';
}

export function locationZoom(kind: string) {
  if (['sea', 'ocean'].includes(kind)) {
    return 3;
  }
  if (['country', 'region', 'province'].includes(kind)) {
    return 4;
  }
  if (kind === 'peninsula') {
    return 5;
  }
  if (['state', 'administrative', 'county', 'river'].includes(kind)) {
    return 6;
  }
  if (['city', 'town', 'village'].includes(kind)) {
    return 10;
  }

  return 12;
}
