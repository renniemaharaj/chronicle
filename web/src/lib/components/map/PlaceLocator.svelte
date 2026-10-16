<script lang="ts">
  import { Input, Label } from 'flowbite-svelte';
  import { untrack } from 'svelte';
  import { Search, X } from '@lucide/svelte';
  import Button from '../ui/Button.svelte';
  import { request } from '../../api/client';
  import { workspace as ws } from '../../state/workspace.svelte';
  import { saveProject } from '../../state/persistence';
  import type { Entity } from '../../types/knowledge';
  import type { PlaceResult } from '../../utils/geography';

  let { entity, onclose }: { entity: Entity; onclose: () => void } = $props();
  let query = $state(untrack(() => entity.name));
  let results = $state<PlaceResult[]>([]);
  let busy = $state(false);
  let searched = $state(false);
  let error = $state('');
  let controller: AbortController | undefined;
  $effect(() => () => controller?.abort());

  async function search(event: SubmitEvent) {
    event.preventDefault();
    if (busy) {
      return;
    }
    busy = true;
    error = '';
    results = [];
    controller = new AbortController();
    try {
      results = await request<PlaceResult[]>(`/places/search?q=${encodeURIComponent(query)}`, { signal: controller.signal });
      searched = true;
    } catch (cause) {
      if (!controller.signal.aborted) {
        error = String(cause);
      }
    } finally {
      busy = false;
    }
  }

  async function choose(result: PlaceResult) {
    entity.attributes.mapLocation = {
      ...result, provider: 'OpenStreetMap / Photon', retrievedAt: new Date().toISOString(),
    };
    entity.version++;
    await saveProject();
    onclose();
  }
</script>

<section class="place-locator" aria-label={`Locate ${entity.name}`}>
  <header><strong>Locate {entity.name}</strong><Button icon aria-label="Close place search" onclick={onclose}><X size={16} /></Button></header>
  <p>Choose a modern map reference. Cities and regions use representative points.</p>
  <form onsubmit={search}>
    <Label for="place-query">Place name and country</Label>
    <div class="place-query-row">
      <Input id="place-query" bind:value={query} required minlength={2} maxlength={200} />
      <Button type="submit" variant="primary" disabled={busy}><Search size={15} />{busy ? 'Searching…' : 'Search'}</Button>
    </div>
  </form>
  {#if error}<p role="alert">{error}</p>{/if}
  {#if searched && !results.length && !busy}<p>No matches. Try a nearby city or a more specific name.</p>{/if}
  <div class="place-results">
    {#each results as result}
      <Button onclick={() => choose(result)}><strong>{result.label}</strong><small>{result.kind.replaceAll('_', ' ')} · Use this location</small></Button>
    {/each}
  </div>
  <small>Search: <a href="https://photon.komoot.io/" target="_blank" rel="noreferrer">Photon</a> · © <a href="https://www.openstreetmap.org/copyright" target="_blank" rel="noreferrer">OpenStreetMap contributors</a></small>
</section>
