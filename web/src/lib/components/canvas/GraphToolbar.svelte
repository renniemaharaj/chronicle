<script lang="ts">
  import { ChevronDown, Filter, Maximize, Network, PanelLeft, PanelRight } from '@lucide/svelte';
  import { workspace as ws } from '../../state/workspace.svelte';
  let filters = $state(false);
</script>

<div class="graph-heading">
  <div class="graph-title">
    <h2><Network size={17} /> Interactive historical graph</h2>
    <p>Explore the life, relationships, ideas, places, and impact of {ws.subject.name}</p>
  </div>
  <div class="graph-tools">
    <span class="tiny-label">Layout</span>
    <select aria-label="Graph view" class="select-like" bind:value={ws.conflictMode}>
      <option value="Shadows">Timeline</option><option value="Parallel">Parallel hypotheses</option>
      <option value="Isolation">Conflict isolation</option><option value="Consequences">Consequences</option>
    </select>
    <button class:pressed={filters} class="tool-button" onclick={() => (filters = !filters)}
      ><Filter size={13} /> Filters</button
    >
    <button class="switch-button" onclick={() => (ws.showRelations = !ws.showRelations)}
      >Relations <span class:checked={ws.showRelations} class="switch"><i></i></span></button
    >
    <button class="switch-button" onclick={() => (ws.showConflicts = !ws.showConflicts)}
      >Conflicts <span class:checked={ws.showConflicts} class="switch"><i></i></span></button
    >
    <button class="tool-button fit-button" aria-label="Fit visible period" onclick={() => ws.fitToken++}
      ><Maximize size={13} /></button
    >
    <button class="tool-button fit-button" aria-label="Expand canvas" onclick={() => { document.querySelector('.graph-panel')?.requestFullscreen(); setTimeout(() => ws.fitToken++, 100); }}><PanelRight size={13} /></button>
  </div>
  {#if filters}
    <div class="filter-popover">
      <strong>Graph layers</strong>
      <label><input type="checkbox" bind:checked={ws.showRelations} /> Relationships</label>
      <label><input type="checkbox" bind:checked={ws.showConflicts} /> Conflicting claims</label>
      <label><input type="checkbox" bind:checked={ws.showLabels} /> Relationship labels</label>
      <label><input type="checkbox" bind:checked={ws.showContext} /> Historical context</label>
      <button
        class="button secondary"
        onclick={() => {
          ws.startYear = null;
          ws.endYear = null;
          ws.fitToken++;
        }}>Reset date window</button
      >
    </div>
  {/if}
</div>
