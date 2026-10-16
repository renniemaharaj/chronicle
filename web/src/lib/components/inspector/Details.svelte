<script lang="ts">
  import { ArrowRight, BookmarkPlus, ShieldCheck } from '@lucide/svelte';
  import { workspace as ws } from '../../state/workspace.svelte';
  import { saveProject } from '../../state/persistence';
  let note = $state('');

  const locations = $derived(ws.project.entities.filter((entity) => ws.event?.locationIds.includes(entity.id)));
  const participants = $derived(ws.project.entities.filter((entity) => ws.event?.participants.includes(entity.id)));
  const notes = $derived(ws.project.notes.filter((item) => item.targetId === ws.selectedId));

  async function addNote() {
    if (!note.trim()) return;
    ws.project.notes.push({
      id: crypto.randomUUID(),
      targetId: ws.selectedId,
      text: note.trim(),
      createdAt: new Date().toISOString(),
    });
    note = '';
    await saveProject();
  }
</script>

<dl class="details-list">
  <div>
    <dt>Precision</dt>
    <dd>
      {ws.selection.time.start?.precision ?? 'Unknown'} · {ws.selection.time.start?.calendar ?? 'Unspecified'} calendar
    </dd>
  </div>
  <div>
    <dt>Type</dt>
    <dd>{ws.selection.kind}</dd>
  </div>
  {#if participants.length}
    <div>
      <dt>Participants</dt>
      <dd class="participant-list">
        {#each participants.slice(0, 4) as entity}<button onclick={() => ws.select(entity.id, 'entity')}
            ><span class="mini-avatar">{entity.name.slice(0, 1)}</span>{entity.name}</button
          >{/each}
      </dd>
    </div>
  {/if}
  {#if locations.length}<div>
      <dt>Location</dt>
      <dd>{locations.map((entity) => entity.name).join(', ')}</dd>
    </div>{/if}
  <div>
    <dt>Assessment</dt>
    <dd><span class="pill blue-pill"><ShieldCheck size={10} /> Evidence inspectable</span></dd>
  </div>
</dl>
<div class="annotation-box">
  <form
    onsubmit={(event) => {
      event.preventDefault();
      addNote();
    }}
  >
    <input aria-label="Research note" bind:value={note} placeholder="Add a research note..." />
    <button aria-label="Save research note" disabled={!note.trim()}><BookmarkPlus size={14} /></button>
  </form>
  {#each notes as item}<p>
      {item.text}<button
        aria-label="Delete note"
        onclick={() => {
          ws.project.notes = ws.project.notes.filter((note) => note.id !== item.id);
          saveProject();
        }}>×</button
      >
    </p>{/each}
</div>
