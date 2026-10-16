<script lang="ts">
  import { BookOpen } from '@lucide/svelte';
  import type { HistoricalEvent } from '../../types/knowledge';
  import { workspace as ws } from '../../state/workspace.svelte';
  import { storyImage } from '../../utils/story-outline';

  let { events }: { events: HistoricalEvent[] } = $props();
  let failed = $state('');
  let loaded = $state('');
  const image = $derived(storyImage(ws.project, events));
</script>

<span class="story-thumbnail" aria-hidden="true">
  {#if image && failed !== image.url}
    <img
      src={image.url}
      class:loaded={loaded === image.url}
      alt=""
      loading="lazy"
      referrerpolicy="no-referrer"
      title={`${image.title} · ${image.creator || 'Creator not recorded'} · ${image.license || 'Rights unknown'}`}
      onerror={() => { failed = image?.url ?? ''; }}
      onload={() => { loaded = image?.url ?? ''; }}
    />
    {#if loaded !== image.url}<BookOpen size={22} strokeWidth={1.25} />{/if}
  {:else}
    <BookOpen size={22} strokeWidth={1.25} />
  {/if}
</span>
