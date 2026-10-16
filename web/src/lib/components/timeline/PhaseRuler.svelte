<script lang="ts">
  import { workspace as ws } from '../../state/workspace.svelte';
  const ticks = $derived(
    Array.from({ length: 11 }, (_, index) => Math.round(ws.bounds[0] + (index * (ws.bounds[1] - ws.bounds[0])) / 10)),
  );
</script>

<div class="phase-ruler">
  <div class="ruler-ticks">
    {#each ticks as year}<span>{year <= 0 ? `${1 - year} BCE` : year}</span>{/each}
  </div>
  <div class="phase-bars" style={`grid-template-columns: repeat(${Math.max(1, ws.project.chapters.length)}, 1fr)`}>
    {#each ws.project.chapters as chapter, index}
      {@const events = ws.project.events.filter((event) => chapter.eventIds.includes(event.id))}
      {@const start = events[0]?.time.start?.year}
      {@const end = events.at(-1)?.time.start?.year}
      <button
        class="phase phase-{['blue', 'green', 'gold', 'violet'][index % 4]}"
        onclick={() => {
          ws.startYear = start ?? null;
          ws.endYear = end ?? null;
          ws.fitToken++;
        }}
      >
        {chapter.title}<small>{start ?? '?'}–{end ?? '?'}</small>
      </button>
    {/each}
  </div>
</div>
