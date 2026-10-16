<script lang="ts">
  import Button from '$lib/components/ui/Button.svelte';
  import { toggleStory } from '$lib/state/playback-controls';
  import EraStrip from '$lib/components/narrative/EraStrip.svelte';
  import StoryPlayback from '$lib/components/narrative/StoryPlayback.svelte';
  import EmptyWorkspace from '$lib/components/shell/EmptyWorkspace.svelte';
  import { watchWorkspace } from '$lib/api/workspace-stream';
  import { watchRun } from '$lib/api/stream';
  import { onMount } from 'svelte';
  import { SvelteFlowProvider } from '@xyflow/svelte';
  import { Check, GitBranch, PanelLeft, PanelRight, Save } from '@lucide/svelte';
  import { workspace as ws } from '$lib/state/workspace.svelte';
  import { cacheWorkspace, initialize, saveProject } from '$lib/state/persistence';
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
    if (watchKey) return watchRun(watchKey);
  });

  let ready = $state(false);
  onMount(() => {
    initialize().finally(() => (ready = true));
    return watchWorkspace();
  });

  $effect(() => {
    if (ready) cacheWorkspace();
  });

  function keyboard(event: KeyboardEvent) {
    if (event.key === 'Escape') ws.dialog = '';
    if ((event.metaKey || event.ctrlKey) && event.key.toLowerCase() === 'k') {
      event.preventDefault();
      ws.dialog = 'Research';
    }
    if ((event.metaKey || event.ctrlKey) && event.key.toLowerCase() === 's') {
      event.preventDefault();
      saveProject();
    }
    if (!ws.hasProject || (event.target as HTMLElement)?.matches('input, textarea, select')) return;
    if (event.key === 'ArrowRight') ws.step(1);
    if (event.key === 'ArrowLeft') ws.step(-1);
    if (event.code === 'Space' && !ws.dialog && !(event.target as HTMLElement)?.closest('button')) {
      event.preventDefault();
      toggleStory(ws);
    }
  }
</script>

<svelte:head
  ><title>{ws.hasProject ? `${ws.subject.name} — Chronicle` : 'Chronicle — Research workspace'}</title></svelte:head
>
<svelte:window onkeydown={keyboard} />

<div class="app-shell" class:light={!ws.dark} class:dark={ws.dark}>
  <Topbar /><StoryPlayback />
  {#if ws.hasProject}
    <div
      class="workspace"
      style={`--left-width: ${ws.leftOpen ? ws.leftWidth : 0}px; --right-width: ${ws.rightOpen ? ws.rightWidth : 0}px;`}
    >
      {#if ws.leftOpen}<Explorer /><ResizeHandle side="left" />{/if}
      <main class="main-column" style={`--story-height: ${ws.storyHeight}%`}>
        <div class="graph-panel">
          <GraphToolbar /><PhaseRuler /><EraStrip />
          {#if ws.view !== 'Map' && ws.showConflicts && ws.conflictMode !== 'Shadows'}<div class="conflict-lane-banner">
              <GitBranch size={13} /><strong>{ws.conflictMode}</strong><span
                >Competing claims preserved. Select a branch to inspect its evidence.</span
              >
            </div>{/if}
          {#if ws.view === 'Map'}<MapCanvas />{:else}<SvelteFlowProvider><GraphCanvas /></SvelteFlowProvider>{/if}
          <TimeControls />
        </div>
        <ResizeHandle side="story" /><NarrativePanel />
      </main>
      {#if ws.rightOpen}<ResizeHandle side="right" /><Inspector />{/if}
    </div>
  {:else}<EmptyWorkspace />{/if}
  <footer class="workspace-status">
    <div>
      <span class="status-indicator" class:online={ws.connected}></span>{ws.connected
        ? 'Local workspace connected'
        : 'Browser workspace'}<span class="status-divider">/</span>{ws.hasProject
        ? 'Research project'
        : 'No project open'}<span class="status-divider">/</span>{ws.hasProject
        ? 'Coverage incomplete'
        : 'Ready to research'}
    </div>
    <div>
      <Button disabled={!ws.hasProject} onclick={() => (ws.leftOpen = !ws.leftOpen)}
        ><PanelLeft size={11} /> Explorer</Button
      ><Button disabled={!ws.hasProject} onclick={() => saveProject()}><Save size={11} /> Save project</Button><Button
        disabled={!ws.hasProject}
        onclick={() => (ws.rightOpen = !ws.rightOpen)}><PanelRight size={11} /> Inspector</Button
      >
    </div>
  </footer>
  <DialogHost />
  {#if ws.toast}<div class="toast" role="status"><Check size={14} />{ws.toast}</div>{/if}
</div>
