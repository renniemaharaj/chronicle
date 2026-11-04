<script lang="ts">
  import { Clapperboard } from '@lucide/svelte';
  import { workspace as ws } from '../../state/workspace.svelte';
  import { episodeImage } from '../../campaign/presentation';
  import { safeSourceURL } from '../../api/links';
  import type { Part } from '../../campaign/types';

  let { part, class: className = '' }: { part?: Part; class?: string } = $props();
  let failed = $state('');
  const asset = $derived(episodeImage(ws.project, part));
  const image = $derived(safeSourceURL(asset?.url));
</script>

<div
  class="reel-artwork {className}"
  title={asset
    ? `${asset.title} · ${asset.creator || 'Creator unknown'} · ${asset.license || 'Rights unknown'}`
    : 'Story artwork'}
>
  {#if image && failed !== image}
    <img
      src={image}
      alt={asset?.caption || asset?.title || ''}
      loading="lazy"
      referrerpolicy="no-referrer"
      onerror={() => (failed = image)}
    />
  {:else}
    <div
      class="reel-artwork-fallback"
      aria-hidden="true"
    >
      <Clapperboard
        size={28}
        strokeWidth={1.3}
      />
      <span>CHRONICLE</span>
    </div>
  {/if}
</div>
