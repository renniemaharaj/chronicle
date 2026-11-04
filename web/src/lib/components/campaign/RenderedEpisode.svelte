<script lang="ts">
  import { Download, MapPinned } from '@lucide/svelte';
  import { clipURL } from '../../campaign/platforms';
  import type { Campaign, Part } from '../../campaign/types';

  let { campaign, part }: { campaign: Campaign; part: Part } = $props();
</script>

<div class="reel-studio-preview">
  {#if part.file}
    <video
      class="reel-rendered-video"
      controls
      preload="metadata"
      src={clipURL(campaign, part)}
      style="width:100%;max-height:760px;border-radius:10px;background:#071323"
    >
      <track kind="captions" />
    </video>
    <a
      class="reel-source-link"
      href={clipURL(campaign, part)}
      download
    >
      <Download size={16} />Download episode {part.number}
    </a>
    <p class="reel-table-footnote">
      Recorded story map, with Follow location enabled. Original saved narration and measured video duration.
    </p>
  {:else}
    <section class="reel-content-card">
      <h3><MapPinned size={20} />Map video preparation</h3>
      <p>Narration is ready. This episode is waiting for its map video export.</p>
    </section>
  {/if}
</div>
