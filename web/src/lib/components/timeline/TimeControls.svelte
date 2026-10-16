<script lang="ts">
  import Button from '$lib/components/ui/Button.svelte';
  import { Popover, Select } from 'flowbite-svelte';
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
  import { toggleStory } from '../../state/playback-controls';
  import { showAllHistory } from '../../state/chapter-focus';
  import { workspace as ws } from '../../state/workspace.svelte';
  import type { Resolution } from '../../types/temporal';
  import DateWindow from './DateWindow.svelte';
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
    <Button icon aria-label="First event" onclick={() => ws.orderedEvents[0] && ws.select(ws.orderedEvents[0].id)}
      ><SkipBack size={14} /></Button
    >
    <Button icon aria-label="Previous milestone" onclick={() => ws.step(-1)}><ChevronLeft size={15} /></Button>
    <Button
      variant="primary" class="play-button"
      icon
      aria-label={ws.playing ? 'Pause story' : 'Play story'}
      onclick={() => toggleStory(ws)}
    >
      {#if ws.playing}<Pause size={15} fill="currentColor" />{:else}<Play size={15} fill="currentColor" />{/if}
    </Button>
    <Button icon aria-label="Next milestone" onclick={() => ws.step(1)}><ChevronRight size={15} /></Button>
    <Select class="w-auto" placeholder="" classes={{ select: "speed-select" }} aria-label="Playback speed" bind:value={ws.speed}>
      <option value={0.75}>0.75×</option><option value={1}>1×</option><option value={1.25}>1.25×</option><option
        value={1.5}>1.5×</option
      ><option value={2}>2×</option>
    </Select>
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
    <Select class="w-auto" placeholder="" aria-label="Temporal resolution" bind:value={ws.resolution}
      >{#each resolutions as resolution}<option>{resolution}</option>{/each}</Select
    >
    <Button icon id="date-window-trigger" aria-label="Date range controls" aria-expanded={rangeOpen} active={rangeOpen}
      ><CalendarDays size={13} /></Button
    >
    <Button
      aria-label="Show entire subject lifespan"
      onclick={() => {
        showAllHistory(ws);
      }}><Maximize size={13} /></Button
    >
  </div>
  <Popover triggeredBy="#date-window-trigger" trigger="click" triggerDelay={0} placement="top-end"
    strategy="fixed" role="dialog" aria-label="Navigate historical time" arrow={false}
    class="date-window" classes={{ content: 'p-0' }} bind:isOpen={rangeOpen}>
    <DateWindow onclose={() => { rangeOpen = false; }} />
  </Popover>
</div>
