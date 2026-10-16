<script lang="ts">
  import { ExternalLink, Image, Play } from '@lucide/svelte';
  import { safeSourceURL } from '../../api/links';
  import { mediaSourceURL, youtubeThumbnail } from '../../utils/media';
  import type { MediaAsset } from '../../types/media';

  let { asset }: { asset: MediaAsset } = $props();
  let failed = $state('');
  const preview = $derived(asset.kind === 'image' ? safeSourceURL(asset.url) : youtubeThumbnail(asset.url));
</script>

<a
  class="media-preview"
  class:video-preview={asset.kind === 'video'}
  href={mediaSourceURL(asset.url)}
  target="_blank"
  rel="noopener noreferrer"
  aria-label={`${asset.kind === 'video' ? 'Watch' : 'Open image'}: ${asset.title} (opens in a new tab)`}
>
  {#if preview && failed !== preview}
    <img
      src={preview}
      alt={asset.kind === 'image' ? asset.caption || asset.title : ''}
      loading="lazy"
      referrerpolicy="no-referrer"
      onerror={() => { failed = preview; }}
    />
  {:else}
    <span class="media-fallback">
      {#if asset.kind === 'video'}<Play size={28} />{:else}<Image size={28} />{/if}
      <span>{asset.kind === 'video' ? 'Watch at source' : 'Open image at source'}</span>
    </span>
  {/if}
  {#if asset.kind === 'video'}
    <span class="media-play"><Play size={20} fill="currentColor" /></span>
  {/if}
  <span class="media-external"><ExternalLink size={13} /></span>
</a>
