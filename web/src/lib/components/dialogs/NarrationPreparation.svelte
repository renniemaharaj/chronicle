<script lang="ts">
  import { onMount, onDestroy, untrack } from 'svelte';
  import Button from '../ui/Button.svelte';
  import { NarrationLibrary, type AudioPart } from '../../speech/library.svelte';
  import type { CloudProvider } from '../../speech/audio-api';
  import { workspace as ws } from '../../state/workspace.svelte';
  import { speech } from '../../speech/preferences.svelte';

  let { provider, parts }: { provider: CloudProvider; parts: AudioPart[] } = $props();
  // The parent recreates this component whenever the narration or voice changes.
  const library = untrack(() => new NarrationLibrary(provider, parts));
  onMount(() => {
    void library.check();
  });

  let active = true;
  onDestroy(() => {
    active = false;
    library.dispose();
  });

  async function generate() {
    ws.playing = false;
    await library.generate();
    if (active && ws.dialog === 'Voice & playback') {
      if (!library.error) {
        speech.error = '';
      }

      speech.status = library.missing.length ? 'Some audio parts still need generation' : 'Saved audio ready';
    }
  }
</script>

<p>Saved narration: {library.saved.filter(Boolean).length} / {parts.length} audio parts.</p>
{#if library.checked}
  <label>
    Narration {library.percent}% complete
    <progress
      aria-label="Narration completion"
      value={library.completed}
      max={parts.length}
    ></progress>
  </label>
  <p>
    Missing: {library.missing.length} parts · {library.characters.toLocaleString()} text characters. Generating
    missing parts uses provider credits. Replaying saved files costs no additional credits.
  </p>
  {#if library.missing.length}
    <details>
      <summary>Missing narration passages</summary>
      <ul>
        {#each [...new Set(library.missing.map((part) => part.title ?? 'Untitled passage'))] as title}
          <li>{title}</li>
        {/each}
      </ul>
    </details>
  {/if}
{/if}
{#if library.blocked === 'insufficient_credits'}
  <p role="status">
    Paused for credits. Research and saved audio are preserved. Resume after renewal or topping up.
  </p>
{:else if library.blocked === 'rate_limited'}
  <p role="status">Paused by the provider rate limit. Wait, then generate missing audio again.</p>
{/if}
<div class="dialog-actions">
  <Button
    disabled={library.busy}
    onclick={() => library.check()}
  >
    Check saved audio
  </Button>
  <Button
    variant="primary"
    disabled={library.busy || !library.checked || !library.missing.length}
    onclick={generate}
  >
    {library.generating ? 'Generating audio…' : 'Generate missing audio'}
  </Button>
  {#if library.generating}
    <Button
      disabled={library.stopRequested}
      onclick={() => library.stop()}
    >
      {library.stopRequested ? 'Stopping after current part…' : 'Stop after current part'}
    </Button>
  {/if}
</div>
<p>
  MP3 files stay in the workspace across restarts. Voice and model versions are kept. Speed changes reuse
  audio.
</p>
{#if library.busy && !library.generating}
  <p role="status">Checking saved audio…</p>
{/if}
{#if library.message}
  <p role="status">{library.message}</p>
{/if}
{#if library.error}
  <p
    role="alert"
    class="speech-error"
  >
    {library.error}
  </p>
{/if}
