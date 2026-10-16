<script lang="ts">
  import { MapPin } from '@lucide/svelte';
  import MapStage from '../map/MapStage.svelte';
  import StorySceneToolbar from './StorySceneToolbar.svelte';
  import StoryReading from './StoryReading.svelte';
  import StoryMapSources from './StoryMapSources.svelte';
  import StoryChronology from './StoryChronology.svelte';
  import StoryGraph from './StoryGraph.svelte';
  import { workspace as ws } from '../../state/workspace.svelte';
  import { storyView } from '../../state/story-view.svelte';
  import { storyGeography } from '../../utils/story-geography';

  let { onfullscreen }: { onfullscreen: () => void } = $props();
  const events = $derived(ws.orderedEvents);
  const event = $derived(events.find((item) => item.id === ws.selectedId));
  const geography = $derived(storyGeography(ws.project, event, events));
</script>

<div class="story-scene" class:alternate-scene={storyView.scene !== 'Map'}>
  {#if storyView.scene === 'Map'}
    <MapStage
      places={geography.frame.places}
      visits={[]}
      camera={geography.frame}
      followCamera={storyView.followLocation}
    />
    <div class="story-map-shade" aria-hidden="true"></div>
    <div class="story-location" aria-live="polite">
      <span><MapPin size={15} />{geography.names.join(' / ') || 'Wider story context'}</span>
      <p>{geography.caption}</p>
      <small>Modern geographic references; historical borders are not shown.</small>
      <StoryMapSources places={geography.frame.places} />
    </div>
  {:else}
    <div class="story-alternate-stage">
      {#if storyView.scene === 'Timeline'}<StoryChronology />{:else}<StoryGraph />{/if}
    </div>
  {/if}
  <StorySceneToolbar {onfullscreen} />
  <StoryReading {event} />
</div>
