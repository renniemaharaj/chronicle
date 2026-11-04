<script lang="ts">
  import { onDestroy, untrack } from 'svelte';
  import type { HistoricalEvent } from '../../types/knowledge';
  import { workspace as ws } from '../../state/workspace.svelte';
  import { saveProject } from '../../state/persistence';
  import { previewVisual, cancelVisualPreview, visualPreview } from '../../state/story-visual-preview.svelte';
  import { usableImage, validVisual } from '../../utils/story-visual';
  import { intervalLabel } from '../../types/temporal';
  import Button from '../ui/Button.svelte';

  let { event }: { event: HistoricalEvent } = $props();
  const initial = untrack(() => event.visual);
  let mediaId = $state(initial?.mediaId ?? '');
  let startSeconds = $state(initial?.startSeconds ?? 4);
  let durationSeconds = $state(initial?.durationSeconds ?? 10);
  let transitionSeconds = $state(initial?.transitionSeconds ?? 0.6);
  let saving = $state(false);
  const images = $derived((ws.project.media ?? []).filter(usableImage));
  const image = $derived(images.find((asset) => asset.id === mediaId));
  const cue = $derived({ mediaId, startSeconds, durationSeconds, transitionSeconds });

  async function save() {
    if (mediaId && !validVisual(cue)) {
      return;
    }

    saving = true;
    cancelVisualPreview();
    event.visual = mediaId ? { ...cue } : undefined;
    await saveProject();
    saving = false;
  }

  onDestroy(cancelVisualPreview);
</script>

<div class="story-visual-editor">
  <label>
    Passage image
    <select
      aria-label="Passage image"
      bind:value={mediaId}
    >
      <option value="">Map only</option>
      {#each images as asset (asset.id)}
        <option value={asset.id}>{asset.title}{asset.generated ? ' · AI illustration' : ''}</option>
      {/each}
    </select>
  </label>
  {#if image}
    <p>{image.caption}</p>
    <small>
      {image.depicts ? intervalLabel(image.depicts) + ' · ' : ''}{image.generated
        ? 'AI illustration, not historical evidence'
        : image.creator || 'Source image'}
    </small>
    <div class="story-visual-time-fields">
      <label>
        Starts at (seconds)
        <input
          aria-label="Image start seconds"
          type="number"
          min="0"
          max="3600"
          step="0.1"
          bind:value={startSeconds}
        />
      </label>
      <label>
        On screen (seconds)
        <input
          aria-label="Image duration seconds"
          type="number"
          min="1"
          max="3600"
          step="0.1"
          bind:value={durationSeconds}
        />
      </label>
      <label>
        Each fade (seconds)
        <input
          aria-label="Image fade seconds"
          type="number"
          min="0"
          max="5"
          step="0.1"
          bind:value={transitionSeconds}
        />
      </label>
    </div>
    <small>Timing follows narration. On-screen time includes both fades; the map returns at the end.</small>
  {/if}
  <div class="story-visual-editor-actions">
    <Button
      disabled={saving || (!!mediaId && !validVisual(cue))}
      onclick={save}
    >
      Save image timing
    </Button>
    {#if image && event.visual && validVisual(cue)}
      <Button
        disabled={visualPreview.running}
        onclick={() => {
          ws.playing = false;
          previewVisual(ws.project.id, event.id, event.visual!);
        }}
      >
        Preview saved image cue
      </Button>
    {/if}
  </div>
  {#if !images.length}
    <p>Add an image to this research project's media to choose it here.</p>
  {/if}
</div>
