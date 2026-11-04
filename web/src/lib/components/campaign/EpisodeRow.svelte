<script lang="ts">
  import { ArrowDownRight, ArrowUpRight, ChevronRight, Play } from '@lucide/svelte';
  import Button from '../ui/Button.svelte';
  import PublicationBadges from './PublicationBadges.svelte';
  import CampaignArtwork from './CampaignArtwork.svelte';
  import { episodeSummary, episodeTitle, partDuration, episodeStatus } from '../../campaign/presentation';
  import { latestObservation, metricLabel, viewTrend } from '../../campaign/observations';
  import type { CampaignWorkspace } from '../../campaign/workspace.svelte';
  import type { Part } from '../../campaign/types';

  let { model, part }: { model: CampaignWorkspace; part: Part } = $props();
  const views = $derived(latestObservation(part, 'views')?.views);
  const trend = $derived(viewTrend(part));
</script>

<tr
  class:selected={model.studio.partNumber === part.number}
  aria-selected={model.studio.partNumber === part.number}
>
  <td class="reel-number">{part.number}</td>
  <td>
    <div class="reel-episode-cell">
      <CampaignArtwork
        {part}
        class="reel-episode-art"
      />
      <Button
        class="reel-episode-button"
        disabled={model.locked}
        onclick={() => model.select(part.number)}
        aria-label={`Episode ${part.number}: ${episodeTitle(part)}`}
      >
        <strong>{episodeTitle(part)}</strong>
        <span>{episodeSummary(part)}</span>
      </Button>
    </div>
  </td>
  <td>
    <span
      class="reel-status"
      data-status={part.status}
    >
      <span class="reel-status-dot"></span>
      {episodeStatus(part)}
    </span>
    {#if part.publication}
      <small class="reel-cell-note">{new Date(part.publication.recordedAt).toLocaleDateString()}</small>
    {/if}
  </td>
  <td>
    <Button
      class="reel-preview-button"
      aria-label={`Preview episode ${part.number}`}
      disabled={model.locked}
      onclick={() => model.openStudio(part.number)}
    >
      <Play size={13} />Review
    </Button>
  </td>
  <td
    class="reel-duration"
    title={part.duration === undefined ? 'Estimated narration duration' : 'Measured recording duration'}
  >
    {partDuration(part)}
  </td>
  <td>
    <PublicationBadges
      campaign={model.studio.campaign!}
      {part}
    />
  </td>
  <td>
    <strong class="reel-view-count">{metricLabel(views)}{views !== undefined ? ' views' : ''}</strong>
    {#if trend !== undefined}
      <span
        class="reel-trend"
        class:negative={trend < 0}
      >
        {#if trend < 0}
          <ArrowDownRight size={12} />
        {:else}
          <ArrowUpRight size={12} />
        {/if}{trend > 0 ? '+' : ''}{trend}%
      </span>
    {:else}
      <small class="reel-cell-note">{views === undefined ? 'No observations' : 'Manual observation'}</small>
    {/if}
  </td>
  <td>
    <Button
      icon
      aria-label={`Details for episode ${part.number}`}
      disabled={model.locked}
      onclick={() => model.select(part.number)}
    >
      <ChevronRight size={15} />
    </Button>
  </td>
</tr>
