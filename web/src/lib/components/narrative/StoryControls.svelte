<script lang="ts">
  import { ChevronLeft, ChevronRight, Pause, Play, RotateCcw, Volume2 } from '@lucide/svelte';
  import Button from '../ui/Button.svelte';
  import { workspace as ws } from '../../state/workspace.svelte';
  import { toggleStory } from '../../state/playback-controls';
  import { focusChapter } from '../../state/chapter-focus';
  import { speech } from '../../speech/preferences.svelte';
  import { restartStory } from '../../state/story-playback.svelte';

  const events = $derived(ws.orderedEvents);
  const index = $derived(events.findIndex((event) => event.id === ws.selectedId));
  const chapter = $derived(ws.project.chapters.find((item) => item.eventIds.includes(ws.selectedId)));
  const chapters = $derived(ws.project.chapters.filter((item) =>
    events.some((event) => item.eventIds.includes(event.id)),
  ));
</script>

<footer class="story-mode-controls">
  <div class="story-progress">
    <label for="story-position">Passage {index + 1} of {events.length}</label>
    <input
      id="story-position"
      aria-label="Story passage"
      type="range"
      min="0"
      max={Math.max(0, events.length - 1)}
      value={Math.max(0, index)}
      disabled={events.length < 2}
      oninput={(change) => {
        const event = events[Number(change.currentTarget.value)];
        if (event) {
          ws.select(event.id);
        }
      }}
    />
    <span aria-live="polite">{ws.storyFinished ? 'Story complete' : speech.status}</span>
  </div>
  <div class="story-transport">
    <div class="story-mode-actions">
      <Button
        icon
        aria-label="Restart story"
        disabled={!events.length}
        onclick={() => restartStory(ws)}
      >
        <RotateCcw size={16} />
      </Button>
      <Button
        icon
        aria-label="Previous milestone"
        disabled={index <= 0}
        onclick={() => ws.step(-1)}
      >
        <ChevronLeft size={20} />
      </Button>
      <Button
        variant="primary"
        aria-label={ws.playing ? 'Pause story' : 'Play story'}
        disabled={!events.length}
        onclick={() => toggleStory(ws)}
      >
        {#if ws.playing}
          <Pause size={18} />Pause
        {:else}
          <Play size={18} />{ws.storyFinished ? 'Replay' : 'Play'}
        {/if}
      </Button>
      <Button
        icon
        aria-label="Next milestone"
        disabled={index < 0 || index >= events.length - 1}
        onclick={() => ws.step(1)}
      >
        <ChevronRight size={20} />
      </Button>
    </div>
    <label class="story-chapter-picker">
      Chapter
      <select
        aria-label="Story chapter"
        value={chapter?.id ?? ''}
        disabled={!chapters.length}
        onchange={(event) => focusChapter(ws, event.currentTarget.value)}
      >
        <option value="" disabled>Choose a chapter</option>
        {#each chapters as item}
          <option value={item.id}>{item.title}</option>
        {/each}
      </select>
    </label>
    <div class="story-mode-actions">
      <select aria-label="Playback speed" bind:value={ws.speed}>
        <option value={0.75}>0.75x</option>
        <option value={1}>1x</option>
        <option value={1.25}>1.25x</option>
        <option value={1.5}>1.5x</option>
        <option value={2}>2x</option>
      </select>
      <Button onclick={() => { ws.dialog = 'Voice & playback'; }}>
        <Volume2 size={16} />Voice & playback
      </Button>
    </div>
  </div>
  {#if speech.error}
    <p class="speech-error" role="alert">{speech.error}</p>
  {/if}
</footer>
