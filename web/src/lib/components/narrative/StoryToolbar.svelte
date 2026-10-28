<script lang="ts">
  import Button from '$lib/components/ui/Button.svelte';
  import { Checkbox, Label } from 'flowbite-svelte';
  import { Map, Volume2 } from '@lucide/svelte';
  import { workspace as ws } from '../../state/workspace.svelte';
  import { showAllHistory } from '../../state/chapter-focus';
  import { speech } from '../../speech/preferences.svelte';
  import { openStoryMode } from '../../state/story-mode';
  import StoryReadiness from '../research/StoryReadiness.svelte';
</script>

<div class="story-toolbar">
  <Button onclick={() => openStoryMode(ws)}><Map size={16} />Story mode</Button>
  <Button
    onclick={() => {
      ws.dialog = 'Research';
    }}
  >
    Update story
  </Button>
  <Button
    onclick={() => {
      ws.dialog = 'Voice & playback';
    }}
  >
    <Volume2 size={16} /> Voice & playback
  </Button>
  <Label><Checkbox bind:checked={ws.followChapter} /> Follow chapter</Label>
  <Button
    onclick={() => {
      ws.followChapter = false;
      showAllHistory(ws);
    }}
  >
    All history
  </Button>
  <Button
    onclick={() => {
      ws.dialog = 'Story';
    }}
  >
    Read all stories
  </Button>
  <span aria-live="polite">{speech.status}</span>
</div>
<details>
  <summary>Story preparation status</summary>
  <StoryReadiness project={ws.project} />
</details>
{#if speech.error}
  <p
    class="speech-error"
    role="alert"
  >
    {speech.error}
  </p>
{/if}
