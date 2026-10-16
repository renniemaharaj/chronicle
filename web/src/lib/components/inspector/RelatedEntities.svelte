<script lang="ts">
  import Button from '$lib/components/ui/Button.svelte';
  import SubjectAvatar from '../subjects/SubjectAvatar.svelte';
  import { ArrowRight, Search, Users } from '@lucide/svelte';
  import { workspace as ws } from '../../state/workspace.svelte';
  const related = $derived(
    ws.project.entities.filter((entity) => {
      if (entity.id === ws.selectedId) return false;
      return (
        ws.event?.participants.includes(entity.id) ||
        ws.event?.locationIds.includes(entity.id) ||
        ws.project.relationships.some(
          (relation) =>
            (relation.fromId === ws.selectedId && relation.toId === entity.id) ||
            (relation.toId === ws.selectedId && relation.fromId === entity.id),
        )
      );
    }),
  );
</script>

<div class="source-list">
  {#each related as entity}
    <div class="related-record">
      <Button class="source-row" onclick={() => ws.select(entity.id, 'entity')}>
        <SubjectAvatar {entity} project={ws.project} size={38} /><span
          ><strong>{entity.name}</strong><small>{entity.kind}</small></span
        ><ArrowRight size={13} />
      </Button>
      <Button
        class="related-research"
        onclick={() => {
          sessionStorage.setItem('chronicle-subject', entity.name);
          ws.dialog = 'Research';
        }}><Search size={11} /> Research independently</Button
      >
    </div>
  {/each}
  {#if !related.length}<p class="empty-state">No relationships are documented for this selection.</p>{/if}
</div>
