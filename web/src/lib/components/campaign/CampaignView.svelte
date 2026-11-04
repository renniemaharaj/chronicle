<script lang="ts">
  import { onMount } from 'svelte';
  import { Clapperboard } from '@lucide/svelte';
  import Button from '../ui/Button.svelte';
  import Explorer from '../shell/Explorer.svelte';
  import ResizeHandle from '../shell/ResizeHandle.svelte';
  import { workspace as ws } from '../../state/workspace.svelte';
  import { CampaignWorkspace } from '../../campaign/workspace.svelte';
  import CampaignHeader from './CampaignHeader.svelte';
  import CampaignPlan from './CampaignPlan.svelte';
  import CampaignSeries from './CampaignSeries.svelte';
  import CampaignStudio from './CampaignStudio.svelte';
  import EpisodeInspector from './EpisodeInspector.svelte';

  const model = new CampaignWorkspace();
  const showInspector = $derived(ws.rightOpen && !!model.studio.part && !model.studioOpen);

  onMount(() => {
    ws.playing = false;

    const narrow = window.matchMedia('(max-width: 1000px)');

    const foldPanels = () => {
      if (narrow.matches) {
        ws.leftOpen = false;
        ws.rightOpen = false;
      }
    };

    foldPanels();
    narrow.addEventListener('change', foldPanels);

    void model.reload();

    return () => {
      narrow.removeEventListener('change', foldPanels);
      void model.studio.abort();
    };
  });
</script>

<svelte:window
  onbeforeunload={(event) => {
    if (model.locked) {
      event.preventDefault();
    }
  }}
/>

<div
  class="reel-shell"
  class:reel-studio-open={model.studioOpen}
  style={`--left-width: ${ws.leftOpen ? ws.leftWidth : 0}px; --right-width: ${showInspector ? ws.rightWidth : 0}px;`}
>
  {#if ws.leftOpen}
    <Explorer compact />
    <ResizeHandle side="left" />
  {/if}
  <main
    class="reel-content"
    aria-label="Reel series workspace"
  >
    <CampaignHeader {model} />
    {#if model.studio.error}
      <p
        class="campaign-error"
        role="alert"
      >
        {model.studio.error}
      </p>
    {/if}
    {#if model.createOpen}
      <CampaignPlan {model} />
    {/if}
    {#if model.studio.campaign}
      {#if model.studioOpen}
        <CampaignStudio {model} />
      {:else}
        <CampaignSeries {model} />
      {/if}
    {:else if !model.studio.busy}
      <section class="reel-welcome">
        <span class="reel-welcome-icon">
          <Clapperboard
            size={42}
            strokeWidth={1.3}
          />
        </span>
        <span class="eyebrow">Your next chapter</span>
        <h2>
          A researched story.
          <br />
          A series worth following.
        </h2>
        <p>
          Give every part of the story its own moment. Plan connected episodes, record a polished vertical
          reading, and keep your progress from first take to publication.
        </p>
        <Button
          variant="primary"
          onclick={() => (model.createOpen = true)}
        >
          Create your first series
        </Button>
        <div class="reel-welcome-steps">
          <span>
            01 <strong>Shape the series</strong>
          </span>
          <span>
            02 <strong>Record the story</strong>
          </span>
          <span>
            03 <strong>Track its reach</strong>
          </span>
        </div>
      </section>
    {:else}
      <div
        class="reel-empty-inline"
        role="status"
      >
        Loading your saved series...
      </div>
    {/if}
  </main>
  {#if showInspector}
    <ResizeHandle side="right" />
    <EpisodeInspector {model} />
  {/if}
</div>
