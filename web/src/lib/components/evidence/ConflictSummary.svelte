<script lang="ts">
  import { ArrowUpRight } from '@lucide/svelte';
  import Button from '../ui/Button.svelte';
  import { workspace as ws } from '../../state/workspace.svelte';

  function label(kind: string) {
    if (/date|temporal/.test(kind)) {
      return 'Conflicting dates';
    }

    return kind.replaceAll('_', ' ').replace(/^./, (letter) => letter.toUpperCase());
  }
</script>

{#if ws.project.conflicts.length}
  <section class="conflicts-period" aria-label="Project disagreements">
    <header><span>Conflicts in this project</span><b>{ws.project.conflicts.length}</b></header>
    {#each ws.project.conflicts as conflict}
      <Button onclick={() => {
        ws.inspectorTab = 'Conflicts';
        requestAnimationFrame(() => document.getElementById(`conflict-${conflict.id}`)?.scrollIntoView({ block: 'nearest' }));
      }}>
        <span class="conflict-summary-copy"><strong>{conflict.title}</strong><small>{label(conflict.kind)}</small></span>
        <ArrowUpRight size={14} />
      </Button>
    {/each}
  </section>
{/if}
