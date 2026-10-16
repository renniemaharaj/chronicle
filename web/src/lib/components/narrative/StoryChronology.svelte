<script lang="ts">
  import { tick } from 'svelte';
  import { workspace as ws } from '../../state/workspace.svelte';
  import { intervalLabel } from '../../types/temporal';
  import Button from '../ui/Button.svelte';

  let list: HTMLElement;
  $effect(() => {
    const id = ws.selectedId;
    tick().then(() => {
      const active = list?.querySelector<HTMLElement>('[aria-current="step"]');
      if (id && active) {
        list.scrollTo({ top: active.offsetTop - list.offsetTop - list.clientHeight / 3 });
      }
    });
  });
</script>

<section class="story-chronology" aria-label="Subject story timeline" bind:this={list}>
  <header><span class="eyebrow">Through time</span><h2>{ws.subject.name}</h2></header>
  {#each ws.orderedEvents as event (event.id)}
    <Button
      class="story-chronology-event"
      current={event.id === ws.selectedId}
      aria-current={event.id === ws.selectedId ? 'step' : undefined}
      onclick={() => ws.select(event.id)}
    >
      <span class="story-chronology-date" title={intervalLabel(event.time)}>{intervalLabel(event.time)}</span>
      <span><strong>{event.title}</strong><small>{event.kind}</small></span>
    </Button>
  {/each}
</section>
