<script lang="ts">
  import type { MapPlace } from '../../utils/map-projection';
  import { locationSourceURL } from '../../utils/geography';
  import ContextEvidence from './ContextEvidence.svelte';

  let { places }: { places: MapPlace[] } = $props();
  const sourced = $derived(places.filter((place) => locationSourceURL(place.location)));
</script>

{#if sourced.length}
  <details class="story-map-sources">
    <summary>Map sources</summary>
    {#each sourced as place}
      <div>
        <a href={locationSourceURL(place.location)} target="_blank" rel="noreferrer">
          {place.entity.name}
        </a>
        <small>{place.location.provider}</small>
        {#if place.location.note}
          <p>{place.location.note}</p>
        {/if}
        {#if place.location.claimIds?.length}
          <details>
            <summary>Coordinate evidence</summary>
            <ContextEvidence claimIds={place.location.claimIds} />
          </details>
        {/if}
      </div>
    {/each}
  </details>
{/if}
