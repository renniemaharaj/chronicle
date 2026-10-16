<script lang="ts">
  import { ExternalLink } from '@lucide/svelte';
  import { workspace as ws } from '../../state/workspace.svelte';
  import { mediaForEvent, mediaSourceURL } from '../../utils/media';
  import { safeSourceURL } from '../../api/links';
  import MediaPreview from '../subjects/MediaPreview.svelte';

  const assets = $derived(mediaForEvent(ws.project, ws.event));
</script>

{#if assets.length}
  <section class="story-media" aria-label="Media for this event">
    <h4>Watch & explore this moment</h4>
    <div class="story-media-list">
      {#each assets as asset (asset.id)}
        {@const source = ws.project.sources.find((item) => item.id === asset.sourceId)}
        <article class="story-media-item">
          <MediaPreview {asset} />
          <div>
            <a href={mediaSourceURL(asset.url)} target="_blank" rel="noopener noreferrer">
              {asset.title}<ExternalLink size={12} />
            </a>
            <small>{asset.kind === 'video' ? 'Video' : 'Image'} · {asset.creator || source?.publisher || 'Source linked'}</small>
          </div>
        </article>
      {/each}
    </div>
  </section>
{/if}
