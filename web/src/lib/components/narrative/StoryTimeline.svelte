<script lang="ts">
  import { onMount, tick } from 'svelte';
  import { workspace as ws } from '../../state/workspace.svelte';
  import { focusChapter } from '../../state/chapter-focus';
  import { activeStoryChapter, nearbyPassages, storyOutline } from '../../utils/story-outline';
  import { intervalLabel } from '../../types/temporal';
  import StoryThumbnail from './StoryThumbnail.svelte';
  import Button from '../ui/Button.svelte';

  let cards: HTMLDivElement;
  let chapters: HTMLElement;
  const outline = $derived(storyOutline(ws.project, ws.orderedEvents));
  const chapter = $derived(activeStoryChapter(outline, ws.selectedId, ws.focusedChapterId));
  const events = $derived(chapter?.events ?? ws.orderedEvents);
  const nearby = $derived(nearbyPassages(events, ws.selectedId));
  const colors = ['var(--blue)', 'var(--teal)', 'var(--gold)', 'var(--violet)'];

  function center(container: HTMLElement, selector: string) {
    const active = container?.querySelector<HTMLElement>(selector);
    if (!active) {
      return;
    }
    const offset = active.getBoundingClientRect().left - container.getBoundingClientRect().left;
    container.scrollTo({ left: container.scrollLeft + offset - container.clientWidth / 2 + active.offsetWidth / 2 });
  }

  onMount(() => {
    const resize = new ResizeObserver(() => {
      center(chapters, '.current');
      center(cards, '[aria-current="step"]');
    });
    resize.observe(chapters);
    resize.observe(cards);
    return () => resize.disconnect();
  });

  $effect(() => {
    const id = chapter?.chapter.id;
    tick().then(() => {
      if (id) {
        center(chapters, '.current');
      }
    });
  });

  $effect(() => {
    const id = ws.selectedId;
    tick().then(() => {
      if (id) {
        center(cards, '[aria-current="step"]');
      }
    });
  });
</script>

<section class="story-timeline" aria-label="Story event navigation">
  <nav class="story-track-chapters" aria-label="Story timeline chapters" bind:this={chapters}>
    {#each outline as item, index}
      <Button
        class="story-track-chapter"
        current={item.chapter.id === chapter?.chapter.id}
        aria-current={item.chapter.id === chapter?.chapter.id ? 'step' : undefined}
        style={`--chapter-color: ${colors[index % colors.length]}`}
        title={`${item.chapter.title} · ${item.dates}`}
        onclick={() => focusChapter(ws, item.chapter.id)}
      >
        <span>{index + 1}. {item.chapter.title}</span>
        <small>{item.dates}</small>
      </Button>
    {/each}
  </nav>
  <div class="story-timeline-events">
    <div class="story-timeline-label">
      <strong>Events</strong>
      <span>{chapter ? 'in this chapter' : 'in this story'}</span>
      <small>{events.length} total</small>
    </div>
    <div class="story-event-cards" bind:this={cards}>
      {#each nearby as event (event.id)}
        <Button
          class="story-event-card"
          current={event.id === ws.selectedId}
          aria-current={event.id === ws.selectedId ? 'step' : undefined}
          title={`${event.title} · ${intervalLabel(event.time)}`}
          onclick={() => ws.select(event.id)}
        >
          <StoryThumbnail events={[event]} />
          <span><strong>{event.title}</strong><small>{intervalLabel(event.time)}</small></span>
        </Button>
      {/each}
    </div>
  </div>
</section>
