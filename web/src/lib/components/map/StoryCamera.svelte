<script lang="ts">
  import { untrack } from 'svelte';
  import type * as Leaflet from 'leaflet';
  import type { MapCameraFrame } from '../../utils/story-geography';
  import { moveStoryCamera } from './story-camera';

  let { map, library, frame, follow }: {
    map?: Leaflet.Map;
    library?: typeof Leaflet;
    frame: MapCameraFrame;
    follow: boolean;
  } = $props();
  let previousKey = '';

  $effect(() => {
    if (!map || !library) {
      return;
    }
    const instance = map;
    const L = library;
    if (!follow) {
      instance.stop();
      previousKey = '';
      return;
    }

    const reducedMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches;
    if (previousKey !== frame.key) {
      const animate = previousKey !== '' && !reducedMotion;
      previousKey = frame.key;
      untrack(() => moveStoryCamera(L, instance, frame, animate));
    }

    const resize = () => moveStoryCamera(L, instance, frame, false);
    instance.on('resize', resize);
    return () => {
      instance.off('resize', resize);
    };
  });
</script>
