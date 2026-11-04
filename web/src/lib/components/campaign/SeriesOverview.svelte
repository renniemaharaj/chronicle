<script lang="ts">
  import { BookOpen, CheckCircle2, Route } from '@lucide/svelte';
  import { safeSourceURL } from '../../api/links';
  import { episodeTitle, seriesProgress } from '../../campaign/presentation';
  import type { Campaign } from '../../campaign/types';

  let { campaign }: { campaign: Campaign } = $props();
  const progress = $derived(seriesProgress(campaign));
  const next = $derived(campaign.parts.find((part) => part.status !== 'published'));
</script>

<div class="reel-overview-grid">
  <section class="reel-content-card">
    <h3><Route size={17} />The production path</h3>
    <div class="reel-progress-row">
      <span>Recorded</span>
      <strong>{progress.recorded} / {progress.total}</strong>
    </div>
    <progress
      value={progress.recorded}
      max={progress.total}
      aria-label="Recorded episodes"
    ></progress>
    <div class="reel-progress-row">
      <span>Published</span>
      <strong>{progress.published} / {progress.total}</strong>
    </div>
    <progress
      value={progress.published}
      max={progress.total}
      aria-label="Published episodes"
    ></progress>
    <p class="reel-detail-note">
      {next
        ? `Next to publish: episode ${next.number}, ${episodeTitle(next)}.`
        : 'Every episode has a publication receipt.'}
    </p>
  </section>
  <section class="reel-content-card">
    <h3><CheckCircle2 size={17} />Research coverage</h3>
    <p>{campaign.gaps.length} recorded coverage gaps</p>
    {#each campaign.gaps as gap}
      <p class="reel-gap">{gap}</p>
    {:else}
      <p class="reel-detail-note">
        No gaps were recorded in this snapshot. Coverage is bounded by the research and available sources.
      </p>
    {/each}
  </section>
  <section class="reel-content-card reel-overview-sources">
    <h3><BookOpen size={17} />Series source references</h3>
    {#each campaign.sources ?? [] as source}
      <a
        class="reel-source-link"
        href={safeSourceURL(source.url)}
        target="_blank"
        rel="noreferrer"
      >
        <BookOpen size={14} />
        <span>{source.title}</span>
      </a>
    {:else}
      <p class="reel-detail-note">No captured source references were saved.</p>
    {/each}
  </section>
</div>
