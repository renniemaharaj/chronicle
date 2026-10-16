<script lang="ts">
  import {
    CalendarDays,
    ChevronLeft,
    ChevronRight,
    Maximize,
    Pause,
    Play,
    SkipBack,
    SkipForward,
  } from '@lucide/svelte';
  import { workspace as ws } from '../../state/workspace.svelte';
  import type { Resolution } from '../../types/temporal';
  let jump = $state('');
  let rangeOpen = $state(false);
  const resolutions: Resolution[] = [
    'Millennium',
    'Century',
    'Decade',
    'Year',
    'Quarter',
    'Month',
    'Week',
    'Day',
    'Event',
  ];
</script>

<div class="timeline-controls">
  <div class="play-controls">
    <button aria-label="First event" onclick={() => ws.select(ws.orderedEvents[0].id)}><SkipBack size={14} /></button>
    <button aria-label="Previous milestone" onclick={() => ws.step(-1)}><ChevronLeft size={15} /></button>
    <button
      class="play-button"
      aria-label={ws.playing ? 'Pause narrative' : 'Play narrative'}
      onclick={() => {
        if (!ws.playing && ws.selectedId === ws.orderedEvents.at(-1)?.id) {
          const first = ws.orderedEvents[0];
          if (first) ws.select(first.id);
        }
        ws.playing = !ws.playing;
      }}
    >
      {#if ws.playing}<Pause size={15} fill="currentColor" />{:else}<Play size={15} fill="currentColor" />{/if}
    </button>
    <button aria-label="Next milestone" onclick={() => ws.step(1)}><ChevronRight size={15} /></button>
    <select class="speed-select" aria-label="Playback speed" bind:value={ws.speed}>
      <option value={0.5}>0.5×</option><option value={1}>1×</option><option value={2}>2×</option><option value={4}
        >4×</option
      >
    </select>
  </div>
  <div class="timeline-slider-wrap">
    <span>{ws.bounds[0]}</span>
    <input
      aria-label="Timeline year"
      type="range"
      min={ws.bounds[0]}
      max={ws.bounds[1]}
      value={ws.year}
      oninput={(event) => ws.jump(Number(event.currentTarget.value))}
    />
    <span>{ws.bounds[1]}</span><strong>{ws.yearLabel()}</strong>
    <div class="coverage-timeline">
      {#each ws.project.coverage as coverage}<i
          class={coverage.status}
          title={`${coverage.startYear}–${coverage.endYear}: ${coverage.status}`}
        ></i>{/each}
    </div>
  </div>
  <div class="time-tools">
    <select aria-label="Temporal resolution" bind:value={ws.resolution}
      >{#each resolutions as resolution}<option>{resolution}</option>{/each}</select
    >
    <button aria-label="Date range controls" class:active={rangeOpen} onclick={() => (rangeOpen = !rangeOpen)}
      ><CalendarDays size={13} /></button
    >
    <button
      aria-label="Show entire subject lifespan"
      onclick={() => {
        ws.startYear = null;
        ws.endYear = null;
        ws.fitToken++;
      }}><Maximize size={13} /></button
    >
  </div>
  {#if rangeOpen}
    <div class="date-window">
      <strong>Navigate historical time</strong>
      <label
        >Jump to year <input
          type="number"
          bind:value={jump}
          placeholder="Negative for BCE"
          onkeydown={(event) => {
            if (event.key === 'Enter') ws.jump(Number(jump));
          }}
        /></label
      >
      <label>From <input type="number" bind:value={ws.startYear} /></label>
      <label>Through <input type="number" bind:value={ws.endYear} /></label>
      <div class="quarter-buttons">
        {#each [0, 1, 2, 3] as quarter}<button
            onclick={() => {
              const width = (ws.bounds[1] - ws.bounds[0]) / 4;
              ws.startYear = Math.floor(ws.bounds[0] + quarter * width);
              ws.endYear = Math.ceil(ws.startYear + width);
              ws.fitToken++;
            }}>Q{quarter + 1}</button
          >{/each}
      </div>
      <button
        class="button"
        onclick={() => {
          if (jump !== '') ws.jump(Number(jump));
          ws.fitToken++;
          rangeOpen = false;
        }}>Apply window</button
      >
    </div>
  {/if}
</div>
