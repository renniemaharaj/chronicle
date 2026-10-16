<script lang="ts">
  import Button from '$lib/components/ui/Button.svelte';
  import { Card, Checkbox, Label } from 'flowbite-svelte';
  import { ExternalLink, LocateFixed } from '@lucide/svelte';
  import { workspace as ws } from '../../state/workspace.svelte';
  import { safeSourceURL } from '../../api/links';
  import { searchText } from '../../utils/subjects';
  import { mediaEvents, mediaSourceURL } from '../../utils/media';
  import CollectionSearch from '../subjects/CollectionSearch.svelte';
  import MediaPreview from '../subjects/MediaPreview.svelte';
  let query = $state('');
  let related = $state(false);
  const entity = $derived(ws.entity ?? ws.subject);
  const kind = $derived(ws.dialog === 'Videos' ? 'video' : 'image');
  const results = $derived(
    (ws.project.media ?? []).filter(
      (asset) =>
        asset.kind === kind &&
        asset.retrieved &&
        safeSourceURL(asset.url) &&
        (related || asset.entityId === entity.id) &&
        searchText(`${asset.title} ${asset.caption} ${asset.creator}`).includes(searchText(query)),
    ),
  );
</script>

<CollectionSearch
  label="Search collected media"
  placeholder={`Search ${ws.dialog.toLowerCase()}…`}
  bind:value={query}
  summary={`${results.length} items · ${entity.name}`}
>
  {#snippet controls()}
    <Label><Checkbox bind:checked={related} /> Include related subjects</Label>
  {/snippet}
</CollectionSearch>
<p class="subject-search-note">
  Collected from research sources. Open the source for context and reuse terms. Videos link to their publisher.
</p>
<div class="media-grid">
  {#each results as asset (asset.id)}
    {@const source = ws.project.sources.find((item) => item.id === asset.sourceId)}
    {@const events = mediaEvents(ws.project, asset)}
    <Card size="xl" class="media-card shadow-none">
      <MediaPreview {asset} />
      <div class="media-copy">
        <strong>{asset.title}</strong>
        {#if asset.caption}<p>{asset.caption}</p>{/if}
        <small>{asset.creator || 'Creator not recorded'} · {asset.license || 'Reuse rights not established'}</small>
        <a href={mediaSourceURL(asset.url)} target="_blank" rel="noopener noreferrer">
          <ExternalLink size={13} />{asset.kind === 'video' ? 'Watch video' : 'Open image'}
        </a>
        {#if safeSourceURL(asset.licenseUrl)}<a
            href={safeSourceURL(asset.licenseUrl)}
            target="_blank"
            rel="noopener noreferrer">License</a
          >{/if}
        {#each events as event}
          <Button class="media-timeline-link" onclick={() => {
            ws.select(event.id);
            ws.dialog = '';
          }}><LocateFixed size={13} />{event.title}</Button>
        {/each}
        {#if safeSourceURL(source?.url)}<a href={safeSourceURL(source?.url)} target="_blank" rel="noopener noreferrer"
            ><ExternalLink size={13} />{source?.title}</a
          >{/if}
      </div>
    </Card>
  {/each}
</div>
{#if !results.length}<p class="empty-state">
    No {ws.dialog.toLowerCase()} have been collected for this selection. Research tools preserve media links from retrieved
    documents.
  </p>{/if}
