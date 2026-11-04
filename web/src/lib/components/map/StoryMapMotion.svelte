<script lang="ts">
  import { onMount, untrack } from 'svelte';
  import type * as Leaflet from 'leaflet';
  import type { MapCameraFrame } from '../../utils/story-geography';
  import { workspace as ws } from '../../state/workspace.svelte';
  import { storyView } from '../../state/story-view.svelte';
  import { captureClock, presentationTime } from '../../state/story-presentation.svelte';
  import { visualPreview } from '../../state/story-visual-preview.svelte';
  import { motionLoop } from './story-motion-loop';
  import { applyCameraMotion, clearCameraMotion, locationFocus } from './story-motion';

  let {
    map,
    library,
    frame,
    follow,
  }: {
    map?: Leaflet.Map;
    library?: typeof Leaflet;
    frame: MapCameraFrame;
    follow: boolean;
  } = $props();
  let reduced = $state(true);
  let suspended = $state(false);
  let ready = $state(false);
  let loop = $state.raw<ReturnType<typeof motionLoop>>();
  const time = $derived(presentationTime(ws.project.id, ws.event?.id).seconds);

  onMount(() => {
    const preference = window.matchMedia('(prefers-reduced-motion: reduce)');

    const update = () => {
      reduced = preference.matches;
    };

    update();
    preference.addEventListener('change', update);

    return () => preference.removeEventListener('change', update);
  });

  $effect(() => {
    if (!map || !library) {
      return;
    }

    const instance = map;
    const focusLayer = locationFocus(library, instance, frame);
    const ticker = motionLoop((seconds, settings) => {
      focusLayer.update(seconds, settings.pulse);
      if (settings.camera) {
        applyCameraMotion(instance, seconds);
      } else {
        clearCameraMotion(instance);
      }
    });
    loop = ticker;

    const anchor = () => {
      ready = true;
      suspended = false;
    };

    const start = () => {
      clearCameraMotion(instance);
      ready = false;
    };

    const manual = () => {
      clearCameraMotion(instance);
      suspended = true;
    };

    const container = instance.getContainer();
    untrack(anchor);
    instance.on('story-camera-start', start);
    instance.on('story-camera-ready', anchor);
    container.addEventListener('pointerdown', manual);
    container.addEventListener('wheel', manual, { passive: true });
    container.addEventListener('keydown', manual);

    return () => {
      ticker.stop();
      clearCameraMotion(instance);
      focusLayer.remove();
      instance.off('story-camera-start', start);
      instance.off('story-camera-ready', anchor);
      container.removeEventListener('pointerdown', manual);
      container.removeEventListener('wheel', manual);
      container.removeEventListener('keydown', manual);
    };
  });

  $effect(() => {
    const pulse = storyView.motion && !reduced;
    const running = ws.playing || visualPreview.running || captureClock.active;
    loop?.update(time, ws.speed, running, {
      pulse,
      camera: pulse && ready && follow && !suspended,
    });
  });
</script>
