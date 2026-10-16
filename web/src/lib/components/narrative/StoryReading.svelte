<script lang="ts">
  import { tick } from 'svelte';
  import { workspace as ws } from '../../state/workspace.svelte';
  import { eventStory } from '../../speech/story';
  import { intervalLabel } from '../../types/temporal';
  import type { HistoricalEvent } from '../../types/knowledge';
  import ContextEvidence from './ContextEvidence.svelte';
  import HistoricalContext from './HistoricalContext.svelte';

  let { event }: { event?: HistoricalEvent } = $props();
  let reading: HTMLElement;
  const story = $derived(eventStory(ws.project, event, ws.activeClaims));
  const paragraphs = $derived(story.split(/\n\s*\n/).filter((text) => text.trim()));
  const chapter = $derived(ws.project.chapters.find((item) => item.eventIds.includes(event?.id ?? '')));
  const alternative = $derived(ws.project.claims.find((claim) =>
    claim.eventId === event?.id && Object.values(ws.activeClaims).includes(claim.id),
  ));
  const claimIds = $derived(alternative ? [alternative.id] : event?.narration?.claimIds ?? event?.claimIds ?? []);

  $effect(() => {
    const id = event?.id;
    tick().then(() => {
      if (id === event?.id && reading) {
        reading.scrollTop = 0;
      }
    });
  });
</script>

<article class="story-reading" aria-label="Current story passage" bind:this={reading}>
  {#if event}
    <div class="eyebrow">{chapter?.title ?? 'Story'} · {intervalLabel(event.time)}</div>
    <h1>{event.title}</h1>
    {#if alternative || event.narration?.interpretation}
      <span class="story-interpretation">Historical interpretation</span>
    {/if}
    {#each paragraphs as paragraph}
      <p class="story-passage">{paragraph}</p>
    {/each}
    <details class="story-provenance">
      <summary>Story evidence</summary>
      <ContextEvidence {claimIds} />
    </details>
    <HistoricalContext />
  {:else}
    <h1>No story events yet</h1>
    <p>Research or open a project with events to begin its story.</p>
  {/if}
</article>
