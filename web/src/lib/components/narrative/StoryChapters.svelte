<script lang="ts">
  import { tick } from 'svelte';
  import { Play } from '@lucide/svelte';
  import Button from '../ui/Button.svelte';
  import StoryThumbnail from './StoryThumbnail.svelte';
  import { workspace as ws } from '../../state/workspace.svelte';
  import { focusChapter } from '../../state/chapter-focus';
  import { activeStoryChapter, storyOutline } from '../../utils/story-outline';

  let list: HTMLElement;
  const outline = $derived(storyOutline(ws.project, ws.orderedEvents));
  const current = $derived(activeStoryChapter(outline, ws.selectedId, ws.focusedChapterId));

  $effect(() => {
    const id = current?.chapter.id;
    tick().then(() => {
      const selected = list?.querySelector<HTMLElement>('[aria-current="step"]');
      if (id && selected) {
        const top = selected.offsetTop - list.offsetTop;
        if (top < list.scrollTop || top + selected.offsetHeight > list.scrollTop + list.clientHeight) {
          list.scrollTo({ top: top - 12 });
        }
      }
    });
  });
</script>

<div class="story-sidebar-section-heading">
  <h2>Chapters</h2>
  <span>{ws.orderedEvents.length} passages</span>
</div>
<nav class="story-chapter-list" aria-label="Story chapter list" bind:this={list}>
  {#each outline as item, index (item.chapter.id)}
    <Button
      class="story-chapter-row"
      current={item.chapter.id === current?.chapter.id}
      aria-current={item.chapter.id === current?.chapter.id ? 'step' : undefined}
      title={`${item.chapter.title} · ${item.dates}`}
      onclick={() => focusChapter(ws, item.chapter.id)}
    >
      <span class="story-chapter-number">{index + 1}</span>
      <StoryThumbnail events={item.events} />
      <span class="story-chapter-copy">
        <strong>{item.chapter.title}</strong>
        <small>{item.dates}</small>
        <small>{item.events.length} {item.events.length === 1 ? 'passage' : 'passages'}</small>
      </span>
      {#if item.chapter.id === current?.chapter.id}<Play size={13} />{/if}
    </Button>
  {:else}
    <p class="story-empty-note">This project has no chapters. Use the passage controls to follow its events.</p>
  {/each}
</nav>
