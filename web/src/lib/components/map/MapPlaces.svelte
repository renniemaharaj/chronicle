<script lang="ts">
  import { ExternalLink, MapPin } from '@lucide/svelte';
  import Button from '../ui/Button.svelte';
  import CollectionSearch from '../subjects/CollectionSearch.svelte';
  import { workspace as ws } from '../../state/workspace.svelte';
  import { intervalLabel } from '../../types/temporal';
  import { mapLocation, openMapURL, locationSourceURL } from '../../utils/geography';
  import type { Entity, HistoricalEvent } from '../../types/knowledge';

  let { places, events, travel, onlocate }: {
    places: Entity[]; events: HistoricalEvent[]; travel: boolean; onlocate: (entity: Entity) => void;
  } = $props();
  let query = $state('');
  const filtered = $derived(places.filter((place) => place.name.toLowerCase().includes(query.toLowerCase())));
</script>

<aside class="map-places" aria-label="Places on this map">
  <CollectionSearch bind:value={query} label="Filter map places" placeholder="Find a place…" summary={travel ? `${events.length} recorded visits` : `${places.filter(mapLocation).length} of ${places.length} located`} />
  {#if travel}
    <details class="map-note"><summary>About dates and routes</summary><p>Recorded presence only. Dashed lines join dated stops; the route is unknown. Overlapping date ranges remain unordered.</p></details>
    {#if !events.length}<p>No visits with recorded physical presence are available in this period.</p>{/if}
    {#each events.filter((event) => event.title.toLowerCase().includes(query.toLowerCase()) || filtered.some((place) => event.locationIds.includes(place.id))) as event}
      <article class="map-stop">
        <Button active={ws.selectedId === event.id} onclick={() => ws.select(event.id)}>
          <small>{intervalLabel(event.time)}</small><strong>{event.title}</strong>
          <span>{event.locationIds.map((id) => ws.project.entities.find((item) => item.id === id)?.name).filter(Boolean).join(' · ')}</span>
        </Button>
      </article>
    {/each}
  {:else}
    {#each filtered as place}
      {@const location = mapLocation(place)}
      {@const linked = events.filter((event) => event.locationIds.includes(place.id))}
      <article class="map-place">
        <Button active={ws.selectedId === place.id || !!ws.event?.locationIds.includes(place.id)} onclick={() => ws.select(place.id, 'entity')}>
          <MapPin size={15} /><span><strong>{place.name}</strong><small>{location ? `${linked.length} timeline events` : 'Location not set'}</small></span>
        </Button>
        <div class="map-place-actions">
          <a href={openMapURL(place)} target="_blank" rel="noreferrer">Open map <ExternalLink size={11} /></a>
          <Button onclick={() => onlocate(place)}>{location ? 'Change' : 'Locate'}</Button>
        </div>
        {#if ws.selectedId === place.id}
          {#if location}
            <p class="map-note">
              Geographic {location.kind.replaceAll('_', ' ')} reference
              {#if locationSourceURL(location)}
                · <a href={locationSourceURL(location)} target="_blank" rel="noreferrer">Coordinate source</a>
              {/if}
              {location.note ?? ''}
            </p>
          {/if}
          {#each linked as event}
            <Button class="place-event" onclick={() => ws.select(event.id)}>{intervalLabel(event.time)} · {event.title}</Button>
          {/each}
        {/if}
      </article>
    {/each}
  {/if}
</aside>
