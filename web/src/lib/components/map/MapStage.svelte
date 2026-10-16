<script lang="ts">
  import { onMount, untrack } from 'svelte';
  import type * as Leaflet from 'leaflet';
  import 'leaflet/dist/leaflet.css';
  import { workspace as ws } from '../../state/workspace.svelte';
  import type { HistoricalEvent } from '../../types/knowledge';
  import type { MapPlace } from '../../utils/map-projection';
  import { drawMap } from './map-layers';

  let { places, visits }: { places: MapPlace[]; visits: HistoricalEvent[] } = $props();
  let container: HTMLDivElement;
  let map = $state.raw<Leaflet.Map>();
  let library = $state.raw<typeof Leaflet>();
  let layer: Leaflet.LayerGroup;
  let error = $state('');
  let tilesFailed = $state(false);
  let hasFit = false;
  let previousFit = '';

  onMount(() => {
    let cancelled = false;
    let observer: ResizeObserver | undefined;
    import('leaflet').then((L) => {
      if (cancelled) {
        return;
      }
      const instance = L.map(container, { center: [18, -25], zoom: 2, minZoom: 1, maxZoom: 19 });
      const tileURL = import.meta.env.VITE_MAP_TILE_URL || 'https://tile.openstreetmap.org/{z}/{x}/{y}.png';
      L.tileLayer(tileURL, {
        maxZoom: 19,
        attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors',
      }).on('tileerror', () => { tilesFailed = true; }).addTo(instance);
      layer = L.layerGroup().addTo(instance);
      observer = new ResizeObserver(() => instance.invalidateSize());
      observer.observe(container);
      library = L;
      map = instance;
    }).catch(() => { error = 'The map could not load. You can still explore places using the list and Open map links.'; });

    return () => {
      cancelled = true;
      observer?.disconnect();
      map?.remove();
    };
  });

  $effect(() => {
    if (!library || !map) {
      return;
    }
    drawMap(library, layer, places, visits, ws.selectedId, (id) => ws.select(id, 'entity'));
  });

  $effect(() => {
    const token = `${ws.fitToken}:${ws.focusToken}`;
    if (!library || !map || !places.length || (hasFit && token === previousFit)) {
      return;
    }
    hasFit = true;
    previousFit = token;
    map.fitBounds(library.latLngBounds(places.map((place) => [place.location.latitude, place.location.longitude])), {
      padding: [35, 35], maxZoom: 13, animate: false,
    });
  });

  $effect(() => {
    const id = ws.selectedId;
    const instance = map;
    const current = places;
    if (!instance) {
      return;
    }
    const selected = current.filter((place) => place.entity.id === id || place.events.some((event) => event.id === id));
    const place = selected[0];
    if (!place) {
      return;
    }
    const position: Leaflet.LatLngTuple = [place.location.latitude, place.location.longitude];
    if (place.entity.id === id) {
      const kind = place.location.kind;
      const zoom = kind === 'country' ? 4 : ['state', 'administrative'].includes(kind) ? 6 : kind === 'city' ? 11 : 14;
      instance.setView(position, zoom);
    } else if (!untrack(() => instance.getBounds().contains(position))) {
      instance.panTo(position);
    }
  });
</script>

<div class="map-stage" bind:this={container} aria-label="OpenStreetMap places" role="region"></div>
{#if error}<p class="map-error" role="alert">{error}</p>{/if}
{#if tilesFailed && !error}<p class="map-tile-notice" role="status">Some map tiles are unavailable. Place links remain available.</p>{/if}
