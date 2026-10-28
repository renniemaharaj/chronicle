<script lang="ts">
  import { handleKeyboard } from '$lib/state/keyboard';
  import WorkspaceStatus from '$lib/components/shell/WorkspaceStatus.svelte';
  import EraStrip from '$lib/components/narrative/EraStrip.svelte';
  import StoryPlayback from '$lib/components/narrative/StoryPlayback.svelte';
  import StoryMode from '$lib/components/narrative/StoryMode.svelte';
  import EmptyWorkspace from '$lib/components/shell/EmptyWorkspace.svelte';
  import { watchWorkspace } from '$lib/api/workspace-stream';
  import { watchRun } from '$lib/api/stream';
  import { onMount } from 'svelte';
  import { SvelteFlowProvider } from '@xyflow/svelte';
  import { Check, GitBranch } from '@lucide/svelte';
  import { workspace as ws } from '$lib/state/workspace.svelte';
  import { cacheWorkspace, initialize } from '$lib/state/persistence';
  import Topbar from '$lib/components/shell/Topbar.svelte';
  import Explorer from '$lib/components/shell/Explorer.svelte';
  import ResizeHandle from '$lib/components/shell/ResizeHandle.svelte';
  import GraphToolbar from '$lib/components/canvas/GraphToolbar.svelte';
  import GraphCanvas from '$lib/components/canvas/GraphCanvas.svelte';
  import MapCanvas from '$lib/components/map/MapCanvas.svelte';
  import PhaseRuler from '$lib/components/timeline/PhaseRuler.svelte';
  import TimeControls from '$lib/components/timeline/TimeControls.svelte';
  import NarrativePanel from '$lib/components/narrative/NarrativePanel.svelte';
  import Inspector from '$lib/components/inspector/Inspector.svelte';
  import DialogHost from '$lib/components/dialogs/DialogHost.svelte';

  const watchKey = $derived(ws.run && ['running', 'pending'].includes(ws.run.status) ? ws.run.id : '');
  $effect(() => {
    if (watchKey) {
      return watchRun(watchKey);
    }
  });

  let ready = $state(false);
  onMount(() => {
    initialize().finally(() => (ready = true));

    return watchWorkspace();
  });

  $effect(() => {
    if (ready) {
      cacheWorkspace();
    }
  });
</script>

<svelte:head>
  <title>{ws.hasProject ? `${ws.subject.name} — Chronicle` : 'Chronicle — Research workspace'}</title>
</svelte:head>
<svelte:window onkeydown={handleKeyboard} />

<div
  class="app-shell"
  class:light={!ws.dark}
  class:dark={ws.dark}
>
  <Topbar />
  <StoryPlayback />
  {#if ws.hasProject && ws.view === 'Story'}
    <StoryMode />
  {:else if ws.hasProject}
    <div
      class="workspace"
      style={`--left-width: ${ws.leftOpen ? ws.leftWidth : 0}px; --right-width: ${ws.rightOpen ? ws.rightWidth : 0}px;`}
    >
      {#if ws.leftOpen}
        <Explorer />
        <ResizeHandle side="left" />
      {/if}
      <main
        class="main-column"
        style={`--story-height: ${ws.storyHeight}%`}
      >
        <div class="graph-panel">
          <GraphToolbar />
          <PhaseRuler />
          <EraStrip />
          {#if ws.view !== 'Map' && ws.showConflicts && ws.conflictMode !== 'Shadows'}
            <div class="conflict-lane-banner">
              <GitBranch size={13} />
              <strong>{ws.conflictMode}</strong>
              <span>Competing claims preserved. Select a branch to inspect its evidence.</span>
            </div>
          {/if}
          {#if ws.view === 'Map'}
            <MapCanvas />
          {:else}
            <SvelteFlowProvider><GraphCanvas /></SvelteFlowProvider>
          {/if}
          <TimeControls />
        </div>
        <ResizeHandle side="story" />
        <NarrativePanel />
      </main>
      {#if ws.rightOpen}
        <ResizeHandle side="right" />
        <Inspector />
      {/if}
    </div>
  {:else}
    <EmptyWorkspace />
  {/if}
  <WorkspaceStatus />
  {#if ws.view !== 'Story' || !ws.hasProject}
    <DialogHost />
  {/if}
  {#if ws.toast}
    <div
      class="toast"
      role="status"
    >
      <Check size={14} />{ws.toast}
    </div>
  {/if}
</div>
