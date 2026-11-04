<script lang="ts">
  import Button from '../ui/Button.svelte';
  import { hasPublication } from '../../campaign/platforms';
  import type { Studio } from '../../campaign/studio.svelte';

  let { studio }: { studio: Studio } = $props();
  let password = $state('');
  let scene = $state('Chronicle Reels');
  let recoveredFile = $state('');
  let recoveredDuration = $state<number>();

  async function perform(action: () => Promise<void>) {
    studio.busy = true;
    try {
      await action();
    } catch (error) {
      studio.error = (error as Error).message;
    } finally {
      studio.busy = false;
    }
  }
</script>

<div class="campaign-panel">
  <h2>Record this part</h2>
  <p>
    Capture the 9:16 region above in an OBS scene, with narration audio enabled. Use a 1080 × 1920 output.
    Check framing and sound before your first take.
  </p>
  <label>
    OBS scene name <input
      bind:value={scene}
      disabled={studio.recording || studio.busy || studio.previewing}
    />
  </label>
  <label>
    OBS WebSocket password <input
      type="password"
      bind:value={password}
      autocomplete="off"
      disabled={studio.recording || studio.busy || studio.previewing}
    />
  </label>
  <p>The connection uses localhost port 4455. Your password stays in this view.</p>
  <div class="campaign-actions">
    <Button
      variant="primary"
      disabled={studio.busy || studio.recording || studio.previewing || studio.part?.status !== 'planned'}
      onclick={() => studio.capture(password, scene)}
    >
      Record part with OBS
    </Button>
    <Button
      disabled={studio.busy || !studio.recording}
      onclick={() => perform(() => studio.abort())}
    >
      Cancel take
    </Button>
  </div>
  {#if !studio.recording}
    <label>
      Preview passage <input
        type="range"
        min="0"
        max={(studio.part?.passages.length ?? 1) - 1}
        bind:value={studio.passageIndex}
        disabled={studio.busy || studio.previewing}
      />
    </label>
  {/if}
  {#if studio.pending}
    <p>OBS saved {studio.pending.file}. The checkpoint needs saving.</p>
    <Button
      disabled={studio.busy || studio.previewing}
      onclick={() => perform(() => studio.saveTake())}
    >
      Save recording checkpoint
    </Button>
  {:else if studio.part?.status === 'recording' && !studio.recording && !studio.busy}
    <details>
      <summary>Recover an interrupted recording</summary>
      <p>
        Check OBS first and stop the take there if necessary. If the entire part completed, attach its saved
        file and measured seconds. Otherwise discard the checkpoint and record this part again.
      </p>
      <label>
        Completed output file <input bind:value={recoveredFile} />
      </label>
      <label>
        Measured seconds <input
          type="number"
          min="0.1"
          step="0.1"
          bind:value={recoveredDuration}
        />
      </label>
      <Button
        onclick={() =>
          perform(() =>
            studio.act({
              action: 'complete',
              takeId: studio.part?.takeId,
              file: recoveredFile,
              duration: recoveredDuration,
            }),
          )}
      >
        Attach complete recording
      </Button>
      <Button onclick={() => perform(() => studio.act({ action: 'interrupt', takeId: studio.part?.takeId }))}>
        Discard incomplete checkpoint
      </Button>
    </details>
  {/if}
  {#if studio.part?.file}
    <p>Saved file: {studio.part.file}</p>
    <p>Measured duration: {studio.part.duration?.toFixed(1)} seconds</p>
    {#if studio.part.status === 'recorded' && !hasPublication(studio.part)}
      <Button
        disabled={studio.busy || studio.previewing}
        onclick={() => perform(() => studio.act({ action: 'retake' }))}
      >
        Retake this part
      </Button>
    {/if}
  {/if}
  {#if studio.part?.takes?.length}
    <details>
      <summary>Completed take history</summary>
      {#each studio.part.takes as take}
        <p>{take.file} · {take.duration.toFixed(1)} seconds</p>
      {/each}
    </details>
  {/if}
</div>
