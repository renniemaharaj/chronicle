<script lang="ts">
  import Button from '$lib/components/ui/Button.svelte';
  import { Clock3 } from '@lucide/svelte';
  import { workspace as ws } from '../../state/workspace.svelte';
  import { contextFor } from '../../speech/story';
  const context = $derived(contextFor(ws.project, ws.event)[0]);
</script>

{#if ws.showContext}
  <Button
    class="era-strip"
    title={context?.title ?? 'Historical context not yet researched'}
    onclick={() => {
      ws.dialog = 'Historical context';
    }}
  >
    <Clock3 size={15} /><span>Era & setting</span>
    <strong>{context?.title ?? 'Context not yet researched'}</strong>
  </Button>
{/if}
