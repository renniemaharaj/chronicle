<script lang="ts">
  import Button from '../ui/Button.svelte';
  import { PanelLeft, PanelRight, Save } from '@lucide/svelte';
  import { workspace as ws } from '../../state/workspace.svelte';
  import { saveProject } from '../../state/persistence';
</script>

<footer class="workspace-status">
  <div>
    <span
      class="status-indicator"
      class:online={ws.connected}
    ></span>
    {ws.connected ? 'Local workspace connected' : 'Browser workspace'}
    <span class="status-divider">/</span>
    {ws.hasProject ? 'Research project' : 'No project open'}
    <span class="status-divider">/</span>
    {ws.hasProject ? 'Coverage incomplete' : 'Ready to research'}
  </div>
  <div>
    <Button
      disabled={!ws.hasProject}
      onclick={() => (ws.leftOpen = !ws.leftOpen)}
    >
      <PanelLeft size={11} /> Explorer
    </Button>
    <Button
      disabled={!ws.hasProject}
      onclick={() => saveProject()}
    >
      <Save size={11} /> Save project
    </Button>
    <Button
      disabled={!ws.hasProject}
      onclick={() => (ws.rightOpen = !ws.rightOpen)}
    >
      <PanelRight size={11} /> Inspector
    </Button>
  </div>
</footer>
