<script lang="ts">
  import { ArrowRight, ExternalLink, Film, X } from '@lucide/svelte';
  import Button from '../ui/Button.svelte';
  import { campaignPlatforms, platformNames, publicationFor } from '../../campaign/platforms';
  import { workspace as ws } from '../../state/workspace.svelte';
  import { episodeSummary, episodeTitle, partDuration, episodeStatus } from '../../campaign/presentation';
  import type { CampaignWorkspace } from '../../campaign/workspace.svelte';
  import CampaignArtwork from './CampaignArtwork.svelte';
  import EpisodeMetrics from './EpisodeMetrics.svelte';
  import EpisodeSources from './EpisodeSources.svelte';
  import EpisodeContinuity from './EpisodeContinuity.svelte';

  let { model }: { model: CampaignWorkspace } = $props();
  const part = $derived(model.studio.part!);
  const campaign = $derived(model.studio.campaign!);
</script>

<aside
  class="reel-inspector"
  aria-label="Episode details"
>
  <header>
    <h2>Episode details</h2>
    <Button
      icon
      aria-label="Close episode details"
      onclick={() => (ws.rightOpen = false)}
    >
      <X size={17} />
    </Button>
  </header>
  <div class="reel-inspector-body">
    <div class="reel-inspector-art">
      <CampaignArtwork {part} />
      <span class="reel-art-label">Story artwork</span>
      <span class="reel-art-duration">{partDuration(part)}</span>
    </div>
    <div class="reel-section-heading">
      <span class="eyebrow">Episode {part.number} / {campaign.parts.length}</span>
      <span
        class="reel-status"
        data-status={part.status}
      >
        {episodeStatus(part)}
      </span>
    </div>
    <h3 class="reel-detail-title">{episodeTitle(part)}</h3>
    <p class="reel-detail-description">{episodeSummary(part)}</p>
    <Button
      variant="primary"
      disabled={model.locked}
      onclick={() => model.openStudio()}
      aria-label="Open recording studio"
    >
      <Film size={15} />Open studio<ArrowRight size={14} />
    </Button>
    <EpisodeSources
      {campaign}
      {part}
    />
    <EpisodeMetrics {part} />
    <section class="reel-detail-card">
      <h3>Publication</h3>
      {#each campaignPlatforms(campaign) as platform}
        {@const receipt = publicationFor(part, platform)}
        {#if receipt}
          <a
            class="reel-source-link"
            href={receipt.url}
            target="_blank"
            rel="noreferrer"
          >
            <span>{platformNames[platform]} · {new Date(receipt.recordedAt).toLocaleString()}</span>
            <ExternalLink size={13} />
          </a>
        {:else}
          <p class="reel-detail-note">{platformNames[platform]}: not yet published.</p>
        {/if}
      {/each}
    </section>
    <EpisodeContinuity
      {campaign}
      {part}
    />
  </div>
</aside>
