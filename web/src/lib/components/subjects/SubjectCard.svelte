<script lang="ts">
  import Button from '$lib/components/ui/Button.svelte';
  import { Card } from 'flowbite-svelte';
  import { ArrowUpRight, Search } from '@lucide/svelte';
  import type { Entity } from '../../types/knowledge';
  import { workspace as ws } from '../../state/workspace.svelte';
  import { intervalLabel } from '../../types/temporal';
  import { contemporary, entityDescription } from '../../utils/subjects';
  import SubjectAvatar from './SubjectAvatar.svelte';
  let { entity, showOverlap = false }: { entity: Entity; showOverlap?: boolean } = $props();
  const relationship = $derived(entityDescription(entity, ws.project));
</script>

<Card size="xl" class="subject-result shadow-none">
  <Button
    class="subject-result-main"
    onclick={() => {
      ws.select(entity.id, 'entity');
      ws.focusToken++;
      ws.dialog = '';
    }}
  >
    <SubjectAvatar {entity} project={ws.project} size={64} />
    <span class="subject-result-copy">
      <strong>{entity.name}</strong>
      <small>{entity.kind} · {intervalLabel(entity.time)}</small>
      {#if relationship}<span class="relationship-caption">{relationship}</span>{/if}
      {#if entity.summary}<span class="subject-result-summary">{entity.summary}</span>{/if}
      {#if showOverlap && contemporary(entity, ws.subject) === 'unknown'}<small
          >Dates incomplete · overlap not established</small
        >{/if}
    </span>
    <ArrowUpRight size={16} />
  </Button>
  <Button
    class="subject-research"
    onclick={() => {
      sessionStorage.setItem('chronicle-subject', entity.name);
      ws.dialog = 'Research';
    }}
  >
    <Search size={14} /> Research independently
  </Button>
</Card>
