<script lang="ts">
  import Button from '$lib/components/ui/Button.svelte';
  import { Select } from 'flowbite-svelte';
  import { BookOpen, ChevronRight, ShieldCheck } from '@lucide/svelte';
  import { workspace as ws } from '../../state/workspace.svelte';
  import { focusChapter } from '../../state/chapter-focus';
  import { eventStory } from '../../speech/story';
  import NarrativeThemes from './NarrativeThemes.svelte';
  import ContextEvidence from './ContextEvidence.svelte';
  import StoryToolbar from './StoryToolbar.svelte';
  import HistoricalContext from './HistoricalContext.svelte';
  import StoryMedia from './StoryMedia.svelte';
  const chapter = $derived(ws.chapter);
  const story = $derived(eventStory(ws.project, ws.event, ws.activeClaims));
  const milestones = $derived(ws.project.events.filter((event) => chapter?.eventIds.includes(event.id)).slice(0, 5));
</script>

<section class="story-panel">
  <div class="story-heading">
    <div>
      <BookOpen size={15} />
      <h2>Story</h2>
      <span class="evidence-badge"><ShieldCheck size={11} /> Evidence linked</span>
    </div>
    <Select class="w-auto" placeholder="" aria-label="Story perspective" bind:value={ws.perspective}
      ><option>Chronological</option><option>Biographical</option><option>Genealogical</option><option
        >Institutional</option
      ><option>Geographical</option><option>Cause and consequence</option></Select
    >
  </div>
  <StoryToolbar />
  <div class="story-content">
    <div class="chapter-list">
      {#each ws.project.chapters as item, index}
        {@const first = ws.project.events.find((event) => item.eventIds.includes(event.id))}
        <Button
          current={item.id === chapter?.id}
          class="chapter"
          onclick={() => {
            focusChapter(ws, item.id);
          }}
        >
          <span>{index + 1}. {item.title}</span><small>{first?.time.start?.year ?? '?'}</small>
        </Button>
      {/each}
    </div>
    <article class="chapter-copy">
      <div class="eyebrow">
        {chapter?.interpretation ? 'Historical interpretation' : 'Documented chronology'} · {ws.perspective}
      </div>
      <h3>{ws.event?.title ?? chapter?.title ?? 'A history waiting to be researched'}</h3>
      <p class="event-story">{story}</p>
      <StoryMedia />
      <HistoricalContext />
      <details class="story-provenance">
        <summary>Story evidence</summary>
        <ContextEvidence claimIds={ws.event?.narration?.claimIds ?? ws.event?.claimIds ?? []} />
      </details>
      <NarrativeThemes />
      <div class="narrative-milestones">
        {#each milestones as event}<Button active={ws.selectedId === event.id} onclick={() => ws.select(event.id)}
            ><span>{event.time.start?.year ?? '?'}</span>{event.title}<ChevronRight size={11} /></Button
          >{/each}
      </div>
    </article>
  </div>
</section>
