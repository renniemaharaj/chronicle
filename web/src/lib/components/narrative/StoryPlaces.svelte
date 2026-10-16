<script lang="ts">
  import { MapPin } from '@lucide/svelte';
  import Button from '../ui/Button.svelte';
  import { workspace as ws } from '../../state/workspace.svelte';
  import { storyView } from '../../state/story-view.svelte';
  import { mapLocation } from '../../utils/geography';

  const places = $derived(ws.project.entities.filter((entity) => entity.kind === 'place'));
</script>

<div class="story-sidebar-section-heading"><h2>Story places</h2><span>{places.length}</span></div>
<div class="story-sidebar-records">
  {#each places as place}
    {@const events = ws.orderedEvents.filter((event) => event.locationIds.includes(place.id))}
    <Button
      class="story-place-row"
      active={!!ws.event?.locationIds.includes(place.id)}
      disabled={!events.length}
      onclick={() => {
        storyView.scene = 'Map';
        storyView.followLocation = true;
        ws.select(events[0].id);
      }}
    >
      <MapPin size={18} />
      <span>
        <strong>{place.name}</strong>
        <small>{events.length} passages · {mapLocation(place) ? 'Map reference saved' : 'Map reference missing'}</small>
      </span>
    </Button>
  {:else}
    <p class="story-empty-note">No places are recorded in this project.</p>
  {/each}
</div>
