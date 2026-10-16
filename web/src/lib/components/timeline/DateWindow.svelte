<script lang="ts">
  import Button from '$lib/components/ui/Button.svelte';
  import { Input, Label } from 'flowbite-svelte';
  import { ArrowRight, CalendarDays, X } from '@lucide/svelte';
  import { workspace as ws } from '../../state/workspace.svelte';

  let { onclose }: { onclose: () => void } = $props();
  let jump = $state<number | undefined>(ws.year);
  let start = $state<number | undefined>(ws.startYear ?? ws.bounds[0]);
  let end = $state<number | undefined>(ws.endYear ?? ws.bounds[1]);
  let error = $state('');
  function apply() {
    if (start === undefined || end === undefined || start > end) {
      error = 'Enter a starting year that is no later than the ending year.';
      return;
    }

    ws.focusedChapterId = '';
    ws.startYear = start;
    ws.endYear = end;
    ws.fitToken++;
    onclose();
  }

  function quarter(index: number) {
    const width = (ws.bounds[1] - ws.bounds[0]) / 4;
    start = Math.floor(ws.bounds[0] + index * width);
    end = Math.ceil(ws.bounds[0] + (index + 1) * width);
  }
</script>

<div class="date-window-content">
  <header>
    <span><CalendarDays size={17} /><strong>Navigate time</strong></span>
    <Button icon aria-label="Close date controls" onclick={onclose}><X size={16} /></Button>
  </header>
  <form onsubmit={(event) => {
    event.preventDefault();
    if (jump !== undefined) {
      ws.jump(jump);
      onclose();
    }
  }}>
    <Label for="jump-year">Jump to a year</Label>
    <div class="date-jump">
      <Input id="jump-year" type="number" bind:value={jump} placeholder="e.g. 1978 or -44" required />
      <Button variant="primary" type="submit">Go <ArrowRight size={14} /></Button>
    </div>
  </form>
  <form onsubmit={(event) => { event.preventDefault(); apply(); }}>
    <div class="date-range-fields">
      <Label>From year<Input type="number" bind:value={start} required /></Label>
      <Label>Through year<Input type="number" bind:value={end} required /></Label>
    </div>
    <p class="date-hint">Use negative years for BCE. Dates retain their recorded calendar.</p>
    <span class="date-shortcut-label">Portions of the timeline</span>
    <div class="quarter-buttons">
      {#each ['First', 'Second', 'Third', 'Last'] as label, index}
        <Button type="button" onclick={() => quarter(index)}>{label} ¼</Button>
      {/each}
    </div>
    {#if error}<p role="alert" class="speech-error">{error}</p>{/if}
    <Button variant="primary" class="date-apply" type="submit">Apply date range <ArrowRight size={14} /></Button>
  </form>
</div>
