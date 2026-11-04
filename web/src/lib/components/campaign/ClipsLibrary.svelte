<script lang="ts">
  import { Copy, Film, FolderOpen } from '@lucide/svelte';
  import Button from '../ui/Button.svelte';
  import { workspace as ws } from '../../state/workspace.svelte';
  import { durationLabel, episodeTitle } from '../../campaign/presentation';
  import { clipURL } from '../../campaign/platforms';
  import type { CampaignWorkspace } from '../../campaign/workspace.svelte';

  let { model }: { model: CampaignWorkspace } = $props();
  const clips = $derived(
    (model.studio.campaign?.parts ?? []).flatMap((part) => {
      const takes = part.takes?.length
        ? part.takes
        : part.file
          ? [
              {
                id: part.takeId ?? part.file,
                file: part.file,
                duration: part.duration ?? 0,
                completedAt: '',
              },
            ]
          : [];

      return takes.map((take) => ({ part, take }));
    }),
  );

  async function copy(path: string) {
    try {
      await navigator.clipboard.writeText(path);
      ws.notify('Recording file path copied.');
    } catch {
      model.studio.error = 'The browser could not copy the path. Select it from the recording details.';
    }
  }
</script>

<section
  class="reel-clips"
  aria-label="Clips library"
>
  <div class="reel-section-heading">
    <div>
      <h2>Your recorded takes</h2>
      <p>{clips.length} saved recordings in this series. Earlier takes remain available for review.</p>
    </div>
    <span class="reel-mode-label"><FolderOpen size={14} />Local files</span>
  </div>
  {#if clips.length}
    <div class="reel-clips-grid">
      {#each clips as { part, take }}
        <article class="reel-clip-card">
          {#if take.file.includes('/exports/')}
            <video
              class="reel-library-video"
              controls
              preload="metadata"
              src={clipURL(model.studio.campaign!, part, take.id)}
            >
              <track kind="captions" />
            </video>
          {/if}
          <div class="reel-clip-header">
            <span class="reel-clip-icon"><Film size={25} /></span>
            <span
              class="reel-status"
              data-status={take.id === part.takeId ? 'recorded' : 'planned'}
            >
              {take.id === part.takeId ? 'Selected take' : 'Earlier take'}
            </span>
          </div>
          <small>Episode {part.number} · {durationLabel(take.duration)}</small>
          <h3>{episodeTitle(part)}</h3>
          <p class="reel-file-path">{take.file}</p>
          <div class="reel-clip-actions">
            <Button
              variant="secondary"
              disabled={model.locked}
              onclick={() => model.openStudio(part.number)}
            >
              Review episode
            </Button>
            <Button
              icon
              aria-label={`Copy file path for ${take.id}`}
              onclick={() => copy(take.file)}
            >
              <Copy size={15} />
            </Button>
          </div>
        </article>
      {/each}
    </div>
    <p class="reel-table-footnote">
      FFmpeg exports can be reviewed here. Other local recording files remain on your workstation.
    </p>
  {:else}
    <div class="reel-empty-inline">
      <Film size={32} />
      <h3>A home for every take</h3>
      <p>Record your first episode in the studio. Completed takes and retakes will appear here.</p>
      <Button
        variant="primary"
        onclick={() => model.openStudio()}
      >
        Open recording studio
      </Button>
    </div>
  {/if}
</section>
