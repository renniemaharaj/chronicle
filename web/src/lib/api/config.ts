import { request } from './client';

export interface PublicConfig {
  mapTileUrl: string;
  mapAttribution: string;
}

let configuration: Promise<PublicConfig> | undefined;

export function publicConfig(): Promise<PublicConfig> {
  configuration ??= request<PublicConfig>('/config').catch(() => {
    configuration = undefined;

    return {
      mapTileUrl: import.meta.env.VITE_MAP_TILE_URL || 'https://tile.openstreetmap.org/{z}/{x}/{y}.png',
      mapAttribution:
        '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors',
    };
  });

  return configuration;
}
