<script lang="ts">
  import { BookOpen, ChevronRight, Lightbulb, ShieldCheck } from '@lucide/svelte';
  import { workspace as ws } from '../../state/workspace.svelte';
  import { intervalLabel } from '../../types/temporal';
  const activeClaim = $derived(
    ws.project.claims.find(
      (claim) => Object.values(ws.activeClaims).includes(claim.id) && claim.eventId === ws.selectedId,
    ),
  );
  const chapter = $derived(ws.chapter);
  const supported = $derived(
    ws.project.claims.some(
      (claim) => claim.eventId === ws.selectedId && !['unverified', 'rejected_superseded'].includes(claim.status),
    ),
  );
  const milestones = $derived(ws.project.events.filter((event) => chapter?.eventIds.includes(event.id)).slice(0, 5));
</script>

<section class="story-panel">
  <div class="story-heading">
    <div>
      <BookOpen size={15} />
      <h2>Narrative</h2>
      <span class="evidence-badge"><ShieldCheck size={11} /> Evidence linked</span>
    </div>
    <select aria-label="Narrative perspective" bind:value={ws.perspective}
      ><option>Chronological</option><option>Biographical</option><option>Genealogical</option><option
        >Institutional</option
      ><option>Geographical</option><option>Cause and consequence</option></select
    >
  </div>
  <div class="story-content">
    <div class="chapter-list">
      {#each ws.project.chapters as item, index}
        {@const first = ws.project.events.find((event) => item.eventIds.includes(event.id))}
        <button
          class:current={item.id === chapter?.id}
          class="chapter"
          onclick={() => {
            if (first) ws.select(first.id);
          }}
        >
          <span>{index + 1}. {item.title}</span><small>{first?.time.start?.year ?? '?'}</small>
        </button>
      {/each}
    </div>
    <article class="chapter-copy">
      <div class="eyebrow">
        {chapter?.interpretation ? 'Historical interpretation' : 'Documented chronology'} · {ws.perspective}
      </div>
      <h3>{ws.event?.title ?? chapter?.title ?? 'A history waiting to be researched'}</h3>
      <p>
        {ws.event && !(activeClaim ? !['unverified', 'rejected_superseded'].includes(activeClaim.status) : supported)
          ? 'Unverified assertion: '
          : ''}{ws.event?.summary ??
          chapter?.summary ??
          'Start a research session to assemble an evidence-backed narrative.'}
      </p>
      <div class="narrative-milestones">
        {#each milestones as event}<button class:active={ws.selectedId === event.id} onclick={() => ws.select(event.id)}
            ><span>{event.time.start?.year ?? '?'}</span>{event.title}<ChevronRight size={11} /></button
          >{/each}
      </div>
    </article>
    <div class="theme-list">
      <h3>Key themes</h3>
      {#each [...new Set(ws.project.events.map((event) => event.kind))].slice(0, 4) as kind, index}
        <button
          class="theme-row"
          onclick={() => {
            const event = ws.project.events.find((event) => event.kind === kind);
            if (event) ws.select(event.id);
          }}
        >
          <i class="theme-icon {['gold', 'blue', 'violet', 'yellow'][index]}"><Lightbulb size={11} /></i><span
            >{kind}</span
          ><strong>{ws.project.events.filter((event) => event.kind === kind).length}</strong>
        </button>
      {/each}
      <button class="story-browse" onclick={() => (ws.dialog = 'Narrative')}
        ><BookOpen size={11} /> Read all chapters</button
      >
    </div>
  </div>
</section>
