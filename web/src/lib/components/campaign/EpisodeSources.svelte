<script lang="ts">
  import { BookOpen, ExternalLink, ShieldCheck } from '@lucide/svelte';
  import { workspace as ws } from '../../state/workspace.svelte';
  import { episodeSources } from '../../campaign/presentation';
  import { safeSourceURL } from '../../api/links';
  import type { Campaign, Part } from '../../campaign/types';

  let { campaign, part }: { campaign: Campaign; part: Part } = $props();
  const sources = $derived(episodeSources(campaign, part, ws.project));
</script>

<details
  class="reel-detail-card"
  open
>
  <summary>
    <BookOpen size={16} />Source material
    <span>({sources.length})</span>
  </summary>
  {#each sources as source}
    <a
      class="reel-source-link"
      href={safeSourceURL(source.url)}
      target="_blank"
      rel="noreferrer"
    >
      <BookOpen size={14} />
      <span>{source.title}</span>
      <ExternalLink size={13} />
    </a>
  {:else}
    <p class="reel-detail-note">No captured references are linked to this episode's current claim records.</p>
  {/each}
  <p class="reel-detail-note">
    <ShieldCheck size={12} />Matched through canonical claims. Source references were saved with the series.
  </p>
</details>
