<script lang="ts">
  import { ArrowLeft, Play, Square, Volume2 } from '@lucide/svelte';
  import Button from '../ui/Button.svelte';
  import { workspace as ws } from '../../state/workspace.svelte';
  import { episodeTitle } from '../../campaign/presentation';
  import type { CampaignWorkspace } from '../../campaign/workspace.svelte';
  import CampaignStage from './CampaignStage.svelte';
  import RenderedEpisode from './RenderedEpisode.svelte';
  import CaptureControls from './CaptureControls.svelte';
  import PublicationLedger from './PublicationLedger.svelte';

  let { model }: { model: CampaignWorkspace } = $props();
  const studio = $derived(model.studio);
</script>

{#if studio.campaign && studio.part}
  <div class="reel-studio-heading">
    <Button
      disabled={model.locked}
      onclick={() => {
        model.section = 'Episodes';
        model.tab = 'Series';
        model.studioOpen = false;
      }}
    >
      <ArrowLeft size={15} />Back to episodes
    </Button>
    <span class="reel-mode-label">
      {studio.recording ? 'Recording' : studio.previewing ? 'Playing story preview' : 'Recording studio'}
    </span>
  </div>
  <div class="reel-studio-title">
    <span class="eyebrow">Episode {studio.partNumber} / {studio.campaign.parts.length}</span>
    <h2>{episodeTitle(studio.part)}</h2>
  </div>
  <div class="reel-studio-layout">
    {#if studio.campaign.engine === 'ffmpeg'}
      <RenderedEpisode
        campaign={studio.campaign}
        part={studio.part}
      />
    {:else}
      <div class="reel-studio-preview">
        <CampaignStage
          campaign={studio.campaign}
          part={studio.part}
          index={studio.passageIndex}
        />
        <div class="reel-preview-controls">
          {#if studio.previewing}
            <Button
              variant="secondary"
              onclick={() => studio.abort()}
            >
              <Square size={14} />Stop preview
            </Button>
          {:else}
            <Button
              variant="secondary"
              disabled={model.locked}
              onclick={() => studio.preview()}
            >
              <Play size={14} />Play episode preview
            </Button>
          {/if}
          <Button
            disabled={model.locked}
            onclick={() => (ws.dialog = 'Voice & playback')}
          >
            <Volume2 size={15} />Voice & playback
          </Button>
        </div>
        <p class="reel-table-footnote">9:16 capture frame · Controls stay outside the recording area.</p>
      </div>
    {/if}
    {#key studio.campaign.id + ':' + studio.partNumber}
      <div class="reel-studio-settings">
        {#if studio.campaign.engine !== 'ffmpeg'}
          <CaptureControls {studio} />
        {/if}
        <div class="campaign-panel"><PublicationLedger {studio} /></div>
      </div>
    {/key}
  </div>
{/if}
