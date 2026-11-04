<script lang="ts">
  import { ArrowUpRight, Globe2, Send } from '@lucide/svelte';
  import Button from '../ui/Button.svelte';
  import { episodeTitle } from '../../campaign/presentation';
  import { nextPublication, platformNames, publicationFor, type Platform } from '../../campaign/platforms';
  import type { CampaignWorkspace } from '../../campaign/workspace.svelte';

  let { model, platform }: { model: CampaignWorkspace; platform: Platform } = $props();
  const campaign = $derived(model.studio.campaign!);
  const next = $derived(nextPublication(campaign, platform));
</script>

<section class="reel-distribution">
  <div class="reel-distribution-banner">
    <span class="reel-platform"><Globe2 size={23} /></span>
    <div>
      <h3>{platformNames[platform]}</h3>
      <p>Independent episode order and publication receipts</p>
    </div>
    <span class="reel-mode-label">Browser handoff</span>
  </div>
  {#each campaign.parts as part}
    {@const receipt = publicationFor(part, platform)}
    <div class="reel-distribution-row">
      <span class="reel-number">{part.number}</span>
      <div>
        <strong>{episodeTitle(part)}</strong>
        <small>
          {receipt
            ? new Date(receipt.recordedAt).toLocaleString()
            : part.number === next
              ? 'Next in the series'
              : 'Waiting for earlier episodes'}
        </small>
      </div>
      {#if receipt}
        <a
          class="reel-source-link"
          href={receipt.url}
          target="_blank"
          rel="noreferrer"
        >
          View publication<ArrowUpRight size={14} />
        </a>
      {:else}
        <Button
          variant="secondary"
          disabled={model.locked}
          onclick={() => model.openStudio(part.number, platform)}
        >
          <Send size={13} />{part.number === next ? 'Manage publication' : 'Review episode'}
        </Button>
      {/if}
    </div>
  {/each}
</section>
