<script lang="ts">
  import { campaignPlatforms } from '../../campaign/platforms';
  import type { CampaignWorkspace } from '../../campaign/workspace.svelte';
  import PlatformDistribution from './PlatformDistribution.svelte';

  let { model }: { model: CampaignWorkspace } = $props();
</script>

<section aria-label="Series distribution">
  <p class="reel-detail-note">
    Facebook Reels and YouTube Shorts have separate posting checkpoints. Confirm the actual publication before
    saving a receipt; completing one platform does not mark the other posted.
  </p>
  {#each campaignPlatforms(model.studio.campaign) as platform}
    <PlatformDistribution
      {model}
      {platform}
    />
  {/each}
</section>
