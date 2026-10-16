<script lang="ts">
  import ContextEvidence from './ContextEvidence.svelte';
  import { workspace as ws } from '../../state/workspace.svelte';
  import { contextFor, contextText } from '../../speech/story';
  import { intervalLabel } from '../../types/temporal';
  const contexts = $derived(contextFor(ws.project, ws.event));
</script>

{#if ws.showContext}
  <div class="historical-context">
    {#each contexts as context}
      <details>
        <summary><strong>{context.title}</strong><span>{context.scope} · {intervalLabel(context.time)}</span></summary>
        <p>{contextText(ws.project, context)}</p>
        <ContextEvidence claimIds={context.claimIds} />
      </details>
    {:else}
      <span>Historical era and wider context have not yet been researched for this period.</span>
    {/each}
  </div>
{/if}
