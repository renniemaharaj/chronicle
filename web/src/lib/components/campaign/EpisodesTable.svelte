<script lang="ts">
  import { Search, ShieldCheck } from '@lucide/svelte';
  import { Input, Select } from 'flowbite-svelte';
  import { episodeSummary, episodeTitle } from '../../campaign/presentation';
  import type { CampaignWorkspace } from '../../campaign/workspace.svelte';
  import EpisodeRow from './EpisodeRow.svelte';

  let { model }: { model: CampaignWorkspace } = $props();
  const parts = $derived(
    (model.studio.campaign?.parts ?? []).filter((part) => {
      const matches = `${episodeTitle(part)} ${episodeSummary(part)}`
        .toLowerCase()
        .includes(model.search.toLowerCase());

      return matches && (model.filter === 'all' || part.status === model.filter);
    }),
  );
  const statuses = [
    {
      value: 'all',
      name: 'All statuses',
    },
    {
      value: 'planned',
      name: 'Draft',
    },
    {
      value: 'recording',
      name: 'Recording',
    },
    {
      value: 'recorded',
      name: 'Ready',
    },
    {
      value: 'published',
      name: 'Published',
    },
  ];
</script>

<div class="reel-table-toolbar">
  <div class="reel-episode-search">
    <Search size={15} />
    <Input
      aria-label="Search episodes"
      placeholder="Find an episode..."
      bind:value={model.search}
    />
  </div>
  <Select
    aria-label="Filter episodes by status"
    items={statuses}
    bind:value={model.filter}
  />
  <span
    class="reel-order-note"
    title="The recorded story keeps its original sequence to preserve continuity"
  >
    <ShieldCheck size={14} />Story order preserved
  </span>
</div>
<div class="reel-table-scroll">
  <table
    class="reel-episode-table"
    aria-label="Series episodes"
  >
    <thead>
      <tr>
        <th>#</th>
        <th>Episode</th>
        <th>Status</th>
        <th>Story preview</th>
        <th>Duration</th>
        <th>Published to</th>
        <th>Performance</th>
        <th><span class="sr-only">Details</span></th>
      </tr>
    </thead>
    <tbody>
      {#each parts as part}
        <EpisodeRow
          {model}
          {part}
        />
      {/each}
    </tbody>
  </table>
  {#if !parts.length}
    <div class="reel-empty-inline">
      <Search size={22} />
      <h3>No matching episodes</h3>
      <p>Try another title or status.</p>
    </div>
  {/if}
</div>
<p class="reel-table-footnote">
  {parts.length} of {model.studio.campaign?.parts.length} episodes · ~ denotes an estimate · Performance reflects
  saved observations.
</p>
