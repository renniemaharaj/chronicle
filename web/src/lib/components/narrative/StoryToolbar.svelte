<script lang="ts">
  import Button from '$lib/components/ui/Button.svelte';
  import { Checkbox, Label } from 'flowbite-svelte';
  import { Volume2 } from '@lucide/svelte';
  import { workspace as ws } from '../../state/workspace.svelte';
  import { showAllHistory } from '../../state/chapter-focus';
  import { speech } from '../../speech/preferences.svelte';
</script>

<div class="story-toolbar">
  <Button
    onclick={() => {
      ws.dialog = 'Voice & playback';
    }}><Volume2 size={16} /> Voice & playback</Button
  >
  <Label><Checkbox bind:checked={ws.followChapter} /> Follow chapter</Label>
  <Button
    onclick={() => {
      ws.followChapter = false;
      showAllHistory(ws);
    }}>All history</Button
  >
  <Button
    onclick={() => {
      ws.dialog = 'Story';
    }}>Read all stories</Button
  >
  <span aria-live="polite">{speech.status}</span>
</div>
{#if speech.error}<p class="speech-error" role="alert">{speech.error}</p>{/if}
