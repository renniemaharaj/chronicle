<script lang="ts">
  import type { HistoricalEvent } from '../../types/knowledge';
  import { workspace as ws } from '../../state/workspace.svelte';
  import { storyImage } from '../../utils/story-outline';
  import MediaPreview from '../subjects/MediaPreview.svelte';

  let { event }: { event: HistoricalEvent } = $props();
  const image = $derived(storyImage(ws.project, [event]));
</script>

{#if image}
  <figure class="story-event-image">
    <MediaPreview asset={image} />
    <figcaption>{image.caption || image.title} · {image.creator || 'Creator not recorded'} · {image.license || 'Rights unknown'}</figcaption>
  </figure>
{/if}
