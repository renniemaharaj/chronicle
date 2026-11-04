<script lang="ts">
  import { CheckCircle2, Clock3, Film, Globe2 } from '@lucide/svelte';
  import { Select } from 'flowbite-svelte';
  import { durationLabel, estimatedSeconds, seriesProgress } from '../../campaign/presentation';
  import type { CampaignWorkspace } from '../../campaign/workspace.svelte';
  import CampaignArtwork from './CampaignArtwork.svelte';
  import SeriesMetrics from './SeriesMetrics.svelte';
  import { campaignPlatforms, platformNames } from '../../campaign/platforms';

  let { model }: { model: CampaignWorkspace } = $props();
  const campaign = $derived(model.studio.campaign!);
  const progress = $derived(seriesProgress(campaign));
  const seconds = $derived(
    campaign.parts.reduce((sum, part) => sum + (part.duration ?? estimatedSeconds(part)), 0),
  );
  const estimated = $derived(campaign.parts.some((part) => part.duration === undefined));
  const choices = $derived(
    model.available.map((item) => ({
      value: item.id,
      name: `${item.title} · ${new Date(item.createdAt).toLocaleDateString()} · ${item.id.slice(-4)}`,
    })),
  );
</script>

<article class="reel-series-card">
  <CampaignArtwork class="reel-cover" />
  <div class="reel-series-copy">
    <div class="reel-section-heading">
      <span class="eyebrow">A story in sequence</span>
      <span
        class="reel-status"
        data-status={progress.published === progress.total ? 'published' : 'recorded'}
      >
        <CheckCircle2 size={12} />{progress.published === progress.total ? 'All published' : 'In progress'}
      </span>
    </div>
    <h2>{campaign.title}</h2>
    <p>
      {progress.recorded} of {progress.total} episodes recorded. {progress.published} published to all targets.
      Every passage has a place in the series.
    </p>
    <div class="reel-chips">
      <span><Film size={13} />{progress.total} episodes</span>
      <span><Clock3 size={13} />{estimated ? '~' : ''}{durationLabel(seconds)} total</span>
      {#each campaignPlatforms(campaign) as target}
        <span><Globe2 size={13} />{platformNames[target]}</span>
      {/each}
    </div>
    {#if choices.length > 1}
      <Select
        aria-label="Saved series"
        items={choices}
        value={campaign.id}
        disabled={model.locked}
        onchange={(event) => model.choose(event.currentTarget.value)}
      />
    {/if}
  </div>
  <SeriesMetrics {campaign} />
</article>
