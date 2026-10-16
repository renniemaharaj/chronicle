<script lang="ts">
  import { onMount } from 'svelte';
  import { ArrowRight, Search, Settings } from '@lucide/svelte';
  import { workspace as ws } from '../../state/workspace.svelte';
  import { api } from '../../api/client';
  import RunProgress from '../research/RunProgress.svelte';
  import type { ResearchRun } from '../../types/project';

  let subject = $state('');
  let busy = $state(false);
  let error = $state('');
  onMount(() => {
    subject = sessionStorage.getItem('chronicle-subject') ?? '';
  });

  async function research() {
    if (!subject.trim()) return;
    busy = true;
    error = '';
    try {
      ws.run = await api.start(subject.trim());
      ws.runs = [ws.run, ...ws.runs.filter((run) => run.id !== ws.run?.id)];
    } catch (reason) {
      error = String(reason);
    } finally {
      busy = false;
    }
  }
</script>

<p class="dialog-intro">
  Enter a subject. Chronicle plans questions, retrieves evidence through your configured MCP tools, checks claims, and
  builds a temporal graph.
</p>
<form
  class="research-subject"
  onsubmit={(event) => {
    event.preventDefault();
    research();
  }}
>
  <Search size={19} /><input
    aria-label="Research subject"
    bind:value={subject}
    placeholder="A person, place, institution, family, idea..."
    required
  /><button class="button" disabled={busy || !subject.trim() || ws.run?.status === 'running'}
    >Research <ArrowRight size={14} /></button
  >
</form>
{#if error}<p class="run-error" role="alert">{error}</p>{/if}
{#if !ws.connected}<p class="form-message">
    The research API is offline. Start the Go backend to run autonomous research.
  </p>{/if}
<button class="settings-inline" onclick={() => (ws.dialog = 'Settings')}
  ><Settings size={13} /> Configure model and MCP tools</button
>
<RunProgress />
