<script lang="ts">
  import { tick } from 'svelte';
  import { workspace as ws } from '../../state/workspace.svelte';
  import { eventStory } from '../../speech/story';
  import { intervalLabel } from '../../types/temporal';
  import type { HistoricalEvent } from '../../types/knowledge';
  import ContextEvidence from './ContextEvidence.svelte';
  import HistoricalContext from './HistoricalContext.svelte';
  import StoryEventImage from './StoryEventImage.svelte';
  import Button from '../ui/Button.svelte';
  import { List, MapPin } from '@lucide/svelte';
  import { storyView } from '../../state/story-view.svelte';
  import { activeStoryChapter, storyOutline } from '../../utils/story-outline';

  let { event }: { event?: HistoricalEvent } = $props();
  let reading: HTMLElement;
  const story = $derived(eventStory(ws.project, event, ws.activeClaims));
  const paragraphs = $derived(story.split(/\n\s*\n/).filter((text) => text.trim()));
  const outline = $derived(storyOutline(ws.project, ws.orderedEvents));
  const section = $derived(activeStoryChapter(outline, event?.id ?? '', ws.focusedChapterId));
  const chapter = $derived(section?.chapter);
  const chapterIndex = $derived(outline.findIndex((item) => item.chapter.id === chapter?.id));
  const places = $derived(ws.project.entities.filter((place) => event?.locationIds.includes(place.id)));
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
    <div class="story-reading-kicker">{chapter ? `Chapter ${chapterIndex + 1} of ${outline.length}` : 'Story passage'}</div>
    {#if chapter}<div class="story-reading-chapter">{chapter.title}</div>{/if}
    <h1>{event.title}</h1>
    <p class="story-reading-date">{intervalLabel(event.time)}</p>
    <div class="story-reading-tags">
      <span>{event.kind}</span>
      {#each places.slice(0, 2) as place}<span><MapPin size={11} />{place.name}</span>{/each}
    </div>
    {#if alternative || event.narration?.interpretation}
      <span class="story-interpretation">Historical interpretation</span>
    {/if}
    {#each paragraphs as paragraph, index}
      <p class="story-passage">{paragraph}</p>
      {#if index === 0}<StoryEventImage {event} />{/if}
    {/each}
    <details class="story-provenance">
      <summary>Story evidence</summary>
      <ContextEvidence {claimIds} />
    </details>
    <HistoricalContext />
    <footer class="story-reading-footer">
      <Button onclick={() => { storyView.scene = 'Timeline'; }}><List size={15} />View on timeline</Button>
      <small>{section?.events.length ?? 1} {section ? 'events in this chapter' : 'event in this passage'}</small>
    </footer>
  {:else}
    <h1>No story events yet</h1>
    <p>Research or open a project with events to begin its story.</p>
  {/if}
</article>
