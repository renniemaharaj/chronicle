<script lang="ts">
  import { ChartNoAxesColumn, Coins, Eye } from '@lucide/svelte';
  import { episodeTitle } from '../../campaign/presentation';
  import { latestObservation, metricLabel, revenueTotals, totalMetric } from '../../campaign/observations';
  import type { CampaignWorkspace } from '../../campaign/workspace.svelte';
  import SeriesMetrics from './SeriesMetrics.svelte';
  import Button from '../ui/Button.svelte';

  let { model }: { model: CampaignWorkspace } = $props();
  const campaign = $derived(model.studio.campaign!);
  const observed = $derived(totalMetric(campaign, 'views').observed);
  const maximum = $derived(
    Math.max(1, ...campaign.parts.map((part) => latestObservation(part, 'views')?.views ?? 0)),
  );
  const revenues = $derived(revenueTotals(campaign));
</script>

<section
  class="reel-analytics"
  aria-label="Campaign analytics"
>
  <div class="reel-section-heading">
    <div>
      <h2>How the story is reaching people</h2>
      <p>Latest saved Facebook observations across published episodes.</p>
    </div>
    <span class="reel-mode-label">Manual tracking</span>
  </div>
  <SeriesMetrics {campaign} />
  <div class="reel-analytics-grid">
    <section class="reel-content-card">
      <h3><ChartNoAxesColumn size={17} />Views by episode</h3>
      {#if observed}
        <div class="reel-bar-chart">
          {#each campaign.parts as part}
            {@const views = latestObservation(part, 'views')?.views}
            <Button
              class="reel-chart-row"
              disabled={model.locked}
              onclick={() => model.select(part.number)}
              title={`${episodeTitle(part)}: ${views === undefined ? 'unknown' : views.toLocaleString()} views`}
            >
              <span>Ep. {part.number}</span>
              <span class="reel-bar-track">
                <span style={`width: ${((views ?? 0) / maximum) * 100}%`}></span>
              </span>
              <strong>{metricLabel(views)}</strong>
            </Button>
          {/each}
        </div>
        <p class="reel-detail-note">
          {observed} of {campaign.parts.length} episodes have observed views. Missing observations are not counted
          as zero.
        </p>
      {:else}
        <div class="reel-empty-inline">
          <Eye size={26} />
          <h3>Your first observations will appear here</h3>
          <p>After publishing, save an episode's performance in the studio.</p>
        </div>
      {/if}
    </section>
    <section class="reel-content-card">
      <h3><Coins size={17} />Recorded revenue</h3>
      {#each revenues as revenue}
        <div class="reel-revenue">
          <span>{revenue.currency}</span>
          <strong>
            {revenue.value.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
          </strong>
        </div>
      {:else}
        <p class="reel-detail-note">No revenue observations recorded.</p>
      {/each}
      <p class="reel-detail-note">
        Latest observed amounts per episode, grouped by currency. Unavailable amounts remain unknown.
      </p>
    </section>
  </div>
</section>
