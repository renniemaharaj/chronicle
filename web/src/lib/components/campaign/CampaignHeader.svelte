<script lang="ts">
  import { ArrowLeft, Clapperboard, PanelRight, Play, Plus, RefreshCw } from '@lucide/svelte';
  import Button from '../ui/Button.svelte';
  import { workspace as ws } from '../../state/workspace.svelte';
  import type { CampaignWorkspace } from '../../campaign/workspace.svelte';

  let { model }: { model: CampaignWorkspace } = $props();
</script>

<div class="reel-breadcrumb">
  <Button
    icon
    aria-label="Back to research"
    disabled={model.locked}
    onclick={() => (ws.view = 'Interactive Graph')}
  >
    <ArrowLeft size={14} />
  </Button>
  <span>{ws.subject.name}</span>
  <span>/</span>
  <strong>Reel Series</strong>
</div>
<header class="reel-heading">
  <div class="reel-heading-copy">
    <span class="reel-heading-icon"><Clapperboard size={25} /></span>
    <div>
      <h1>Reel Series</h1>
      <p>Your research. A continuous story. An audience to reach.</p>
    </div>
  </div>
  <div class="reel-heading-actions">
    <Button
      icon
      aria-label="Reload saved progress"
      title="Reload saved progress"
      disabled={model.locked}
      onclick={() => model.reload()}
    >
      <RefreshCw size={16} />
    </Button>
    {#if !model.studioOpen}
      <Button
        variant="secondary"
        disabled={model.locked || !model.studio.campaign}
        onclick={() => model.previewSeries()}
      >
        <Play size={14} />Preview series
      </Button>
    {/if}
    <Button
      variant="primary"
      disabled={model.locked}
      onclick={() => (model.createOpen = !model.createOpen)}
    >
      <Plus size={16} />Create series
    </Button>
    {#if model.studio.campaign && !model.studioOpen}
      <Button
        icon
        aria-label="Toggle episode details"
        disabled={model.locked}
        onclick={() => (ws.rightOpen = !ws.rightOpen)}
      >
        <PanelRight size={17} />
      </Button>
    {/if}
  </div>
</header>
