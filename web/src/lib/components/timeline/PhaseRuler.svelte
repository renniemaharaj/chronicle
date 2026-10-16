<script lang="ts">
  import Button from '$lib/components/ui/Button.svelte';
  import { focusChapter } from '../../state/chapter-focus';
  import { workspace as ws } from '../../state/workspace.svelte';
  import { dateLabel, knownYear } from '../../types/temporal';

  function chapterRange(ids: string[]) {
    const years = ws.project.events
      .filter((event) => ids.includes(event.id))
      .flatMap((event) => [knownYear(event.time.start), knownYear(event.time.end)])
      .filter((year): year is number => year !== undefined);
    const label = (year: number) => dateLabel({ year, precision: 'year', calendar: 'historical' });
    return years.length ? `${label(Math.min(...years))}–${label(Math.max(...years))}` : 'Dates unknown';
  }
</script>

<nav class="phase-ruler" aria-label="Story chapters">
  <div class="phase-bars" style={`--chapters: ${Math.max(1, ws.project.chapters.length)}`}>
    {#each ws.project.chapters as chapter, index}
      {@const range = chapterRange(chapter.eventIds)}
      <Button
        class="phase phase-{['blue', 'green', 'gold', 'violet'][index % 4]}"
        current={ws.focusedChapterId === chapter.id}
        aria-pressed={ws.focusedChapterId === chapter.id}
        aria-label={`${chapter.title}, ${range}`}
        title={`${chapter.title} · ${range}`}
        onclick={() => focusChapter(ws, chapter.id)}
      >
        <span>{chapter.title}</span><small>{range}</small>
      </Button>
    {/each}
  </div>
</nav>
