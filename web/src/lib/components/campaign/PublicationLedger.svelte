<script lang="ts">
  import { Select } from 'flowbite-svelte';
  import Button from '../ui/Button.svelte';
  import type { Studio } from '../../campaign/studio.svelte';
  import {
    campaignPlatforms,
    nextPublication,
    platformNames,
    publicationFor,
  } from '../../campaign/platforms';
  import CampaignMetrics from './CampaignMetrics.svelte';

  let { studio }: { studio: Studio } = $props();
  let receipt = $state('');
  const platform = $derived(studio.publicationPlatform);
  const choices = $derived(
    campaignPlatforms(studio.campaign).map((value) => ({
      value,
      name: platformNames[value],
    })),
  );
  const posted = $derived(studio.part ? publicationFor(studio.part, platform) : undefined);
  const next = $derived(studio.campaign ? nextPublication(studio.campaign, platform) : undefined);
  const composer = $derived(
    platform === 'youtube' ? 'https://studio.youtube.com/' : 'https://www.facebook.com/reels/create/',
  );

  async function publish() {
    studio.busy = true;
    try {
      await studio.act({
        action: 'publish',
        platform,
        url: receipt,
      });
      receipt = '';
    } catch (error) {
      studio.error = (error as Error).message;
    } finally {
      studio.busy = false;
    }
  }
</script>

<h2>Publication ledger</h2>
{#if studio.campaign?.accounts?.[platform]}
  <p>Destination: {studio.campaign.accounts[platform]?.name}</p>
{/if}
{#if choices.length > 1}
  <label>
    Publication platform <Select
      aria-label="Publication platform"
      items={choices}
      bind:value={studio.publicationPlatform}
      onchange={() => (receipt = '')}
      disabled={studio.recording || studio.busy || studio.previewing}
    />
  </label>
{/if}
<p>
  {next
    ? `Next to publish on ${platformNames[platform]}: part ${next}`
    : `All parts have ${platformNames[platform]} publication receipts.`}
</p>
<p>
  Open the signed-in browser, verify the intended destination, and publish this episode. Save the direct
  publication URL only after confirming success.
</p>
<a
  class="campaign-link"
  href={composer}
  target="_blank"
  rel="noreferrer"
>
  Open {platformNames[platform]} composer
</a>
{#if studio.part}
  <label>
    Suggested caption <textarea
      readonly
      value={`${studio.campaign?.title} | Part ${String(studio.part.number).padStart(2, '0')} | ${studio.part.passages[0]?.title ?? studio.campaign?.title}`}></textarea>
  </label>
  {#if posted}
    <p>
      <a
        href={posted.url}
        target="_blank"
        rel="noreferrer"
      >
        View posted part {studio.part.number}
      </a>
    </p>
    <p>
      Receipt recorded on {new Date(posted.recordedAt).toLocaleString()}. {platformNames[platform]} is tracked independently.
    </p>
    <CampaignMetrics
      {studio}
      {platform}
    />
  {:else}
    <label>
      Published reel URL <input
        type="url"
        bind:value={receipt}
        placeholder={platform === 'youtube'
          ? 'https://www.youtube.com/shorts/...'
          : 'https://www.facebook.com/reel/...'}
      />
    </label>
    <Button
      variant="primary"
      disabled={studio.busy ||
        studio.recording ||
        studio.previewing ||
        !['recorded', 'published'].includes(studio.part.status) ||
        studio.part.number !== next}
      onclick={publish}
    >
      Save publication receipt
    </Button>
  {/if}
{/if}
