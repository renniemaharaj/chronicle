<script lang="ts">
  import Button from '$lib/components/ui/Button.svelte';
  import { Checkbox, Label, Select, Toggle } from 'flowbite-svelte';
  import { Filter, Map, Maximize, Network, PanelRight } from '@lucide/svelte';
  import { showAllHistory } from '../../state/chapter-focus';
  import { workspace as ws } from '../../state/workspace.svelte';
  import { navigateView } from '../../state/view-navigation';
  let filters = $state(false);
</script>

<div class="graph-heading">
  <div class="graph-title">
    <h2><Network size={17} />{ws.view === 'Map' ? 'Places & travel' : 'Interactive historical graph'}</h2>
    <p>Explore the life, relationships, ideas, places, and impact of {ws.subject.name}</p>
  </div>
  <div class="graph-tools">
    <div class="canvas-view-switch" role="group" aria-label="Canvas view">
      <Button icon active={ws.view !== 'Map'} aria-label="Show graph view" title="Graph view" onclick={() => navigateView(ws, 'Interactive Graph')}><Network size={15} /></Button>
      <Button icon active={ws.view === 'Map'} aria-label="Show map view" title="Map view" onclick={() => navigateView(ws, 'Map')}><Map size={15} /></Button>
    </div>
    {#if ws.view !== 'Map'}
    <span class="tiny-label">Layout</span>
    <Select class="w-auto" placeholder="" aria-label="Graph view" classes={{ select: "select-like" }} bind:value={ws.conflictMode}>
      <option value="Shadows">Timeline</option><option value="Parallel">Parallel hypotheses</option>
      <option value="Isolation">Conflict isolation</option><option value="Consequences">Consequences</option>
    </Select>
    <Button pressed={filters} class="tool-button" onclick={() => (filters = !filters)}
      ><Filter size={13} /> Filters</Button
    >
    <Toggle size="small" bind:checked={ws.showRelations}>Relations</Toggle>
    <Toggle size="small" bind:checked={ws.showConflicts}>Conflicts</Toggle>
    {/if}
    <Button icon class="tool-button fit-button" aria-label="Fit visible period" onclick={() => ws.fitToken++}
      ><Maximize size={13} /></Button
    >
    <Button icon class="tool-button fit-button" aria-label="Expand canvas" onclick={() => { document.querySelector('.graph-panel')?.requestFullscreen(); setTimeout(() => ws.fitToken++, 100); }}><PanelRight size={13} /></Button>
  </div>
  {#if filters && ws.view !== 'Map'}
    <div class="filter-popover">
      <strong>Graph layers</strong>
      <Label><Checkbox bind:checked={ws.showRelations} /> Relationships</Label>
      <Label><Checkbox bind:checked={ws.showConflicts} /> Conflicting claims</Label>
      <Label><Checkbox bind:checked={ws.showLabels} /> Relationship labels</Label>
      <Label><Checkbox bind:checked={ws.showContext} /> Historical context</Label>
      <Button
        variant="secondary"
        onclick={() => {
          showAllHistory(ws);
        }}>Reset date window</Button
      >
    </div>
  {/if}
</div>
