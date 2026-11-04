<script lang="ts">
  import { Film, ListVideo, Scan, Sparkles } from '@lucide/svelte';
  import Button from '../ui/Button.svelte';
  import { storyView } from '../../state/story-view.svelte';
</script>

<div
  class="story-production-controls"
  aria-label="Story production tools"
>
  <Button
    pressed={storyView.portrait}
    aria-label="Reel and Shorts preview"
    onclick={() => {
      storyView.portrait = !storyView.portrait;
      storyView.scene = 'Map';
    }}
  >
    <Film size={15} /> 9:16 preview
  </Button>
  <Button
    pressed={storyView.motion}
    aria-label="Animate story map"
    onclick={() => {
      storyView.motion = !storyView.motion;
    }}
  >
    <Sparkles size={15} /> Map motion
  </Button>
  {#if !storyView.portrait}
    <Button
      pressed={storyView.eventCards}
      aria-label="Show story event cards"
      onclick={() => {
        storyView.eventCards = !storyView.eventCards;
      }}
    >
      <ListVideo size={15} /> Event cards
    </Button>
  {/if}
  {#if storyView.portrait}
    <Button
      pressed={storyView.safeAreas}
      aria-label="Show Reel safe areas"
      onclick={() => {
        storyView.safeAreas = !storyView.safeAreas;
      }}
    >
      <Scan size={15} /> Safe areas
    </Button>
    <small>1080 × 1920 · Captions timed approximately</small>
  {/if}
</div>
