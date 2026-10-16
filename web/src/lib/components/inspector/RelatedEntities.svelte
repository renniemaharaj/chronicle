<script lang="ts">
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
      <button class="source-row" onclick={() => ws.select(entity.id, 'entity')}>
        <span class="related-thumb"><Users size={17} /></span><span
          ><strong>{entity.name}</strong><small>{entity.kind}</small></span
        ><ArrowRight size={13} />
      </button>
      <button
        class="related-research"
        onclick={() => {
          sessionStorage.setItem('chronicle-subject', entity.name);
          ws.dialog = 'Research';
        }}><Search size={11} /> Research independently</button
      >
    </div>
  {/each}
  {#if !related.length}<p class="empty-state">No relationships are documented for this selection.</p>{/if}
</div>
