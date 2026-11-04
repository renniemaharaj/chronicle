<script lang="ts">
  import { onMount } from 'svelte';
  import type * as Leaflet from 'leaflet';
  import 'leaflet/dist/leaflet.css';
  import { workspace as ws } from '../../state/workspace.svelte';
  import type { HistoricalEvent } from '../../types/knowledge';
  import type { MapPlace } from '../../utils/map-projection';
  import { drawMap } from './map-layers';
  import StoryMapEffects from './StoryMapEffects.svelte';
  import type { MapCameraFrame } from '../../utils/story-geography';
  import { storyCoordinates } from './story-camera';
  import { focusSelectedPlace } from './map-selection';
  import { publicConfig } from '../../api/config';
  import { createMap } from './map-instance';

  let {
    places,
    visits,
    camera,
    followCamera = true,
  }: {
    places: MapPlace[];
    visits: HistoricalEvent[];
    camera?: MapCameraFrame;
    followCamera?: boolean;
  } = $props();
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
    Promise.all([import('leaflet'), publicConfig()])
      .then(([L, settings]) => {
        if (cancelled) {
          return;
        }

        const created = createMap(
          L,
          container,
          settings,
          () => {
            tilesFailed = true;
          },
          !!camera,
        );
        layer = created.layer;
        observer = created.observer;
        library = L;
        map = created.map;
      })
      .catch(() => {
        error = camera
          ? 'The map could not load. Story reading and playback remain available.'
          : 'The map could not load. You can still explore places using the list and Open map links.';
      });

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

    const select = camera ? undefined : (id: string) => ws.select(id, 'entity');
    const positions = camera ? storyCoordinates(camera, map.getCenter().lng) : undefined;
    drawMap(library, layer, places, visits, camera?.overview ? '' : ws.selectedId, select, positions);
  });

  $effect(() => {
    const token = `${ws.fitToken}:${ws.focusToken}`;
    if (camera || !library || !map || !places.length || (hasFit && token === previousFit)) {
      return;
    }

    hasFit = true;
    previousFit = token;
    map.fitBounds(
      library.latLngBounds(places.map((place) => [place.location.latitude, place.location.longitude])),
      {
        padding: [35, 35],
        maxZoom: 13,
        animate: false,
      },
    );
  });

  $effect(() => {
    if (camera) {
      return;
    }

    const instance = map;
    if (!instance) {
      return;
    }

    focusSelectedPlace(instance, ws.selectedId, places);
  });
</script>

<div
  class="map-stage"
  bind:this={container}
  aria-label={camera ? 'Story location map' : 'OpenStreetMap places'}
  role="region"
></div>
{#if camera}
  <StoryMapEffects
    {map}
    {library}
    frame={camera}
    follow={followCamera}
  />
{/if}
{#if error}
  <p
    class="map-error"
    role="alert"
  >
    {error}
  </p>
{/if}
{#if tilesFailed && !error}
  <p
    class="map-tile-notice"
    role="status"
  >
    {camera
      ? 'Some map tiles are unavailable. Story reading and playback remain available.'
      : 'Some map tiles are unavailable. Place links remain available.'}
  </p>
{/if}
