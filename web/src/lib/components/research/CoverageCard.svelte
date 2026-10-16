<script lang="ts">
  import Button from '$lib/components/ui/Button.svelte';
  import { Check, CircleHelp, GitBranch, ShieldCheck } from '@lucide/svelte';
  import { workspace as ws } from '../../state/workspace.svelte';

  const researched = $derived(ws.project.coverage.filter((period) => period.status === 'researched').length);
  const total = $derived(ws.project.coverage.length);
</script>

<section class="progress-card">
  <div class="progress-title">
    <ShieldCheck size={13} class="coverage-icon" /><strong>Research coverage</strong>
    <Button icon aria-label="Inspect research coverage" onclick={() => (ws.dialog = 'Coverage')}
      ><CircleHelp size={13} /></Button
    >
  </div>
  {#if ws.run && ws.run.status === 'running'}
    <div class="progress-row">
      <span class="live-dot"></span>
      <div><strong>{ws.run.phase}</strong><small>{ws.run.question}</small></div>
    </div>
  {:else}
    <div class="coverage-summary">
      <strong>{researched}<span> / {total}</span></strong><span>intervals researched</span>
    </div>
    <p class="coverage-note">
      {ws.project.fixture
        ? 'Selective fixture. Coverage remains incomplete.'
        : 'Coverage follows the available evidence.'}
    </p>
  {/if}
  <div class="coverage-segments">
    {#each ws.project.coverage as period}
      <Button
        class={period.status}
        title={`${period.startYear}–${period.endYear}: ${period.status}`}
        aria-label={`${period.startYear} to ${period.endYear}: ${period.status}`}
        onclick={() => {
          ws.startYear = period.startYear;
          ws.endYear = period.endYear;
          ws.fitToken++;
        }}
      ></Button>
    {/each}
  </div>
  <div class="coverage-line"><span>Sources recorded</span><strong>{ws.project.sources.length}</strong></div>
  <div class="coverage-line"><span>Claims available</span><strong>{ws.project.claims.length}</strong></div>
  <div class="coverage-line">
    <span>Conflicts preserved</span><strong class="warn">{ws.project.conflicts.length}</strong>
  </div>
  <Button class="coverage-details" onclick={() => (ws.dialog = 'Coverage')}
    >Inspect periods and gaps <GitBranch size={11} /></Button
  >
</section>
