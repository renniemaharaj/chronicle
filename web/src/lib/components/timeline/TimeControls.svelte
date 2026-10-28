<script lang="ts">
  import Button from '$lib/components/ui/Button.svelte';
  import { Popover, Select } from 'flowbite-svelte';
  import { CalendarDays, Maximize } from '@lucide/svelte';
  import TimelinePlayback from './TimelinePlayback.svelte';
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
  <TimelinePlayback />
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
    <span>{ws.bounds[1]}</span>
    <strong>{ws.yearLabel()}</strong>
    <div class="coverage-timeline">
      {#each ws.project.coverage as coverage}
        <i
          class={coverage.status}
          title={`${coverage.startYear}–${coverage.endYear}: ${coverage.status}`}
        ></i>
      {/each}
    </div>
  </div>
  <div class="time-tools">
    <Select
      class="w-auto"
      placeholder=""
      aria-label="Temporal resolution"
      bind:value={ws.resolution}
    >
      {#each resolutions as resolution}
        <option>{resolution}</option>
      {/each}
    </Select>
    <Button
      icon
      id="date-window-trigger"
      aria-label="Date range controls"
      aria-expanded={rangeOpen}
      active={rangeOpen}
    >
      <CalendarDays size={13} />
    </Button>
    <Button
      aria-label="Show entire subject lifespan"
      onclick={() => {
        showAllHistory(ws);
      }}
    >
      <Maximize size={13} />
    </Button>
  </div>
  <Popover
    triggeredBy="#date-window-trigger"
    trigger="click"
    triggerDelay={0}
    placement="top-end"
    strategy="fixed"
    role="dialog"
    aria-label="Navigate historical time"
    arrow={false}
    class="date-window"
    classes={{ content: 'p-0' }}
    bind:isOpen={rangeOpen}
  >
    <DateWindow
      onclose={() => {
        rangeOpen = false;
      }}
    />
  </Popover>
</div>
