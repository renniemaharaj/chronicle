<script lang="ts">
  import Button from '../ui/Button.svelte';
  import { Select } from 'flowbite-svelte';
  import { ChevronLeft, ChevronRight, Pause, Play, SkipBack } from '@lucide/svelte';
  import { toggleStory } from '../../state/playback-controls';
  import { workspace as ws } from '../../state/workspace.svelte';
</script>

<div class="play-controls">
  <Button
    icon
    aria-label="First event"
    onclick={() => ws.orderedEvents[0] && ws.select(ws.orderedEvents[0].id)}
  >
    <SkipBack size={14} />
  </Button>
  <Button
    icon
    aria-label="Previous milestone"
    onclick={() => ws.step(-1)}
  >
    <ChevronLeft size={15} />
  </Button>
  <Button
    variant="primary"
    class="play-button"
    icon
    aria-label={ws.playing ? 'Pause story' : 'Play story'}
    onclick={() => toggleStory(ws)}
  >
    {#if ws.playing}
      <Pause
        size={15}
        fill="currentColor"
      />
    {:else}
      <Play
        size={15}
        fill="currentColor"
      />
    {/if}
  </Button>
  <Button
    icon
    aria-label="Next milestone"
    onclick={() => ws.step(1)}
  >
    <ChevronRight size={15} />
  </Button>
  <Select
    class="w-auto"
    placeholder=""
    classes={{ select: 'speed-select' }}
    aria-label="Playback speed"
    bind:value={ws.speed}
  >
    <option value={0.75}>0.75×</option>
    <option value={1}>1×</option>
    <option value={1.25}>1.25×</option>
    <option value={1.5}>1.5×</option>
    <option value={2}>2×</option>
  </Select>
</div>
