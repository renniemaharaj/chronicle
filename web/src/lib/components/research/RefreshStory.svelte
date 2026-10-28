<script lang="ts">
  import Button from '../ui/Button.svelte';
  import { api } from '../../api/client';
  import { workspace as ws } from '../../state/workspace.svelte';

  let busy = $state(false);
  let error = $state('');

  async function refresh() {
    busy = true;
    error = '';
    try {
      ws.run = await api.refresh(ws.project.id);
      ws.runs = [ws.run, ...ws.runs.filter((run) => run.id !== ws.run?.id)];
    } catch (cause) {
      error = (cause as Error).message;
    } finally {
      busy = false;
    }
  }
</script>

{#if ws.hasProject}
  <div>
    <p>Refresh “{ws.project.title}” with recent events, corrections and unresolved evidence.</p>
    <p>
      Existing sources and competing claims stay attached. Changed passages need new audio; unchanged files
      are reused.
    </p>
    <Button
      disabled={busy || ws.run?.status === 'running'}
      onclick={refresh}
    >
      {busy ? 'Starting refresh…' : 'Update current story'}
    </Button>
    {#if error}
      <p role="alert">{error}</p>
    {/if}
  </div>
{/if}
