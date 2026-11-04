<script lang="ts">
  import type * as Leaflet from 'leaflet';
  import type { MapCameraFrame } from '../../utils/story-geography';
  import StoryCamera from './StoryCamera.svelte';
  import StoryMapMotion from './StoryMapMotion.svelte';

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
</script>

<StoryCamera
  {map}
  {library}
  {frame}
  {follow}
/>
<StoryMapMotion
  {map}
  {library}
  {frame}
  {follow}
/>
