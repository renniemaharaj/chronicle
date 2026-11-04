<script lang="ts">
  import FacebookMark from './FacebookMark.svelte';
  import { campaignPlatforms, platformNames, publicationFor } from '../../campaign/platforms';
  import type { Campaign, Part } from '../../campaign/types';

  let { campaign, part }: { campaign: Campaign; part: Part } = $props();
</script>

<span class="reel-publication-badges">
  {#each campaignPlatforms(campaign) as platform}
    {@const receipt = publicationFor(part, platform)}
    {#if receipt}
      <a
        class="reel-platform"
        href={receipt.url}
        target="_blank"
        rel="noreferrer"
        aria-label={`Open episode ${part.number} on ${platform === 'youtube' ? 'YouTube' : 'Facebook'}`}
        title={platformNames[platform]}
      >
        {#if platform === 'facebook'}
          <FacebookMark size={18} />
        {:else}
          <span
            class="reel-youtube-mark"
            aria-hidden="true"
          >
            ▶
          </span>
        {/if}
      </a>
    {:else}
      <span
        class="reel-pending-platform"
        title={`${platformNames[platform]}: not yet published`}
      >
        {platform === 'facebook' ? 'FB' : 'YT'}
      </span>
    {/if}
  {/each}
</span>
