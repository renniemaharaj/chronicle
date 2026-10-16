<script lang="ts">
  import { ArrowLeft, LocateFixed, MapPin, Maximize } from '@lucide/svelte';
  import Button from '../ui/Button.svelte';
  import MapStage from '../map/MapStage.svelte';
  import StoryReading from './StoryReading.svelte';
  import StoryControls from './StoryControls.svelte';
  import StoryMapSources from './StoryMapSources.svelte';
  import DialogHost from '../dialogs/DialogHost.svelte';
  import { workspace as ws } from '../../state/workspace.svelte';
  import { closeStoryMode } from '../../state/story-mode';
  import { storyGeography } from '../../utils/story-geography';

  let follow = $state(true);
  let stage: HTMLElement;
  const events = $derived(ws.orderedEvents);
  const event = $derived(events.find((item) => item.id === ws.selectedId));
  const geography = $derived(storyGeography(ws.project, event, events));

  async function fullscreen() {
    try {
      if (document.fullscreenElement) {
        await document.exitFullscreen();
      } else {
        await stage.requestFullscreen();
      }
    } catch {
      ws.notify('Fullscreen is unavailable in this browser.');
    }
  }
</script>

<main class="story-mode" aria-label="Story mode" bind:this={stage}>
  <header class="story-mode-header">
    <Button onclick={() => closeStoryMode(ws)}><ArrowLeft size={16} />Exit story mode</Button>
    <div class="story-mode-title">
      <span>Story mode</span>
      <strong>{ws.subject.name}</strong>
    </div>
    <div class="story-mode-actions">
      <Button pressed={follow} onclick={() => { follow = !follow; }}>
        <LocateFixed size={16} />Follow location
      </Button>
      <Button icon aria-label="Toggle story fullscreen" onclick={fullscreen}>
        <Maximize size={16} />
      </Button>
    </div>
  </header>
  <div class="story-scene">
    <MapStage
      places={geography.frame.places}
      visits={[]}
      camera={geography.frame}
      followCamera={follow}
    />
    <div class="story-map-shade" aria-hidden="true"></div>
    <div class="story-location" aria-live="polite">
      <span><MapPin size={15} />{geography.names.join(' / ') || 'Wider story context'}</span>
      <p>{geography.caption}</p>
      <small>Modern geographic references; historical borders are not shown.</small>
      <StoryMapSources places={geography.frame.places} />
    </div>
    <StoryReading {event} />
  </div>
  <StoryControls />
  <DialogHost />
</main>
