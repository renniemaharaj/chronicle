<script lang="ts">
  import Button from '$lib/components/ui/Button.svelte';
  import { Checkbox, Input, Label } from 'flowbite-svelte';
  import { safeSourceURL } from "../../api/links";
  import { ArrowUpRight, BookOpen, Link } from '@lucide/svelte';
  import type { Source } from '../../types/knowledge';
  import { saveProject } from '../../state/persistence';
  import { workspace as ws } from '../../state/workspace.svelte';
  let { sources = ws.project.sources }: { sources?: Source[] } = $props();
  const independent = $derived(new Set(sources.map((source) => source.independenceGroup || source.id)).size);
</script>

<div class="source-list">
  <div class="provenance-summary">
    <Link size={12} />
    {independent} source group{independent === 1 ? '' : 's'} · {sources.length} source{sources.length === 1 ? '' : 's'}
  </div>
  {#each sources as source}
    <a class="source-row" href={safeSourceURL(source.url)} target="_blank" rel="noreferrer">
      <span class="source-marker"><BookOpen size={14} /></span>
      <span
        ><strong>{source.title}</strong><small
          >{source.classification} · {source.publisher || 'Publisher unknown'}</small
        ><small>{source.author || 'Author unknown'}</small></span
      >
      <ArrowUpRight size={13} />
    </a>
    {#each source.limitations ?? [] as limitation}<p class="source-note">{limitation}</p>{/each}
    <details class="source-review">
      <summary>Review source independence</summary>
      <Label
        >Dependency group <Input aria-label="Source dependency group" bind:value={source.independenceGroup} /></Label
      >
      <Label
        ><Checkbox bind:checked={source.independenceReviewed} /> Independence reviewed against citation dependencies</Label
      >
      <Button variant="secondary" onclick={() => saveProject()}>Save assessment</Button>
    </details>
    <div class="source-metadata">
      <span>Retrieved {source.retrievedAt ? source.retrievedAt.slice(0, 10) : 'not recorded'}</span><span
        >{source.verification}</span
      >
    </div>
  {/each}
  {#if !sources.length}<p class="empty-state">No source records have been retrieved for this selection.</p>{/if}
  <p class="source-note">Repeated accounts in the same independence group do not count as independent corroboration.</p>
</div>
