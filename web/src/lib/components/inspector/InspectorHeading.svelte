<script lang="ts">
  import Button from '../ui/Button.svelte';
  import { Activity, ChevronLeft, ChevronRight, X } from '@lucide/svelte';
  import { workspace as ws } from '../../state/workspace.svelte';

  const index = $derived(ws.orderedEvents.findIndex((event) => event.id === ws.selectedId));
</script>

<div class="inspector-head">
  <div>
    <span class="event-icon"><Activity size={12} /></span>
    <span>{ws.selection.kind}</span>
  </div>
  <div class="inspector-pagination">
    <Button
      icon
      aria-label="Previous event"
      onclick={() => ws.step(-1)}
    >
      <ChevronLeft size={13} />
    </Button>
    <span>{Math.max(1, index + 1)} of {ws.project.events.length}</span>
    <Button
      icon
      aria-label="Next event"
      onclick={() => ws.step(1)}
    >
      <ChevronRight size={13} />
    </Button>
    <Button
      icon
      aria-label="Close inspector"
      onclick={() => (ws.rightOpen = false)}
    >
      <X size={15} />
    </Button>
  </div>
</div>
