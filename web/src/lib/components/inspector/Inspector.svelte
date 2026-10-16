<script lang="ts">
  import { Activity, ChevronLeft, ChevronRight, Sparkles, X } from '@lucide/svelte';
  import { workspace as ws } from '../../state/workspace.svelte';
  import { intervalLabel } from '../../types/temporal';
  import Details from './Details.svelte';
  import EvidenceList from '../evidence/EvidenceList.svelte';
  import SourceList from '../evidence/SourceList.svelte';
  import ConflictList from '../evidence/ConflictList.svelte';
  import SourceDistribution from '../evidence/SourceDistribution.svelte';
  import RelatedEntities from './RelatedEntities.svelte';

  const alternative = $derived(
    ws.project.claims.find(
      (claim) => Object.values(ws.activeClaims).includes(claim.id) && claim.eventId === ws.selectedId && claim.time,
    )?.time,
  );
  const title = $derived('title' in ws.selection ? ws.selection.title : ws.selection.name);
  const index = $derived(ws.orderedEvents.findIndex((event) => event.id === ws.selectedId));
  const claims = $derived(
    ws.project.claims.filter((claim) => claim.eventId === ws.selectedId || claim.entityId === ws.selectedId),
  );
  const evidence = $derived(ws.project.evidence.filter((item) => claims.some((claim) => claim.id === item.claimId)));
  const sourceIds = $derived([...new Set(evidence.map((item) => item.sourceId))]);
  const sources = $derived(ws.project.sources.filter((source) => sourceIds.includes(source.id)));
  const tabs = $derived([
    { title: 'Details', count: null },
    { title: 'Sources', count: sources.length },
    { title: 'Related', count: ws.event?.participants.length ?? 0 },
    { title: 'Conflicts', count: ws.project.conflicts.length },
  ]);
</script>

<aside class="right-panel">
  <div class="inspector-head">
    <div><span class="event-icon"><Activity size={12} /></span><span>{ws.selection.kind}</span></div>
    <div class="inspector-pagination">
      <button aria-label="Previous event" onclick={() => ws.step(-1)}><ChevronLeft size={13} /></button>
      <span>{Math.max(1, index + 1)} of {ws.project.events.length}</span>
      <button aria-label="Next event" onclick={() => ws.step(1)}><ChevronRight size={13} /></button>
      <button aria-label="Close inspector" onclick={() => (ws.rightOpen = false)}><X size={15} /></button>
    </div>
  </div>
  <div class="event-title"><h2>{title}</h2></div>
  <div class="event-subtitle">
    <p class="event-date">{intervalLabel(ws.selection.time)}</p>
    <span class="major-tag"
      ><Sparkles size={10} /> {ws.event && ws.event.importance >= 9 ? 'Milestone' : ws.selection.kind}</span
    >
  </div>
  {#if alternative}<p class="active-date">
      Active view: {intervalLabel(alternative)}. The canonical record above is preserved.
    </p>{/if}
  {#if ws.entity?.attributes.portrait}
    <div class="entity-portrait-banner">
      <img src={String(ws.entity.attributes.portrait)} alt={ws.entity.name} /><span>{ws.entity.summary}</span>
    </div>
  {:else}
    <div class="selection-banner">
      <Activity size={27} /><span>{ws.selection.kind}<small>{intervalLabel(ws.selection.time)}</small></span>
    </div>
  {/if}
  <p class="event-summary">
    {ws.selection.summary ||
      'This related entity was discovered in the current project. Its history has not been researched independently.'}
  </p>
  <div class="inspector-tabs">
    {#each tabs as tab}<button
        class:active={ws.inspectorTab === tab.title}
        onclick={() => (ws.inspectorTab = tab.title)}>{tab.title}{tab.count !== null ? ` (${tab.count})` : ''}</button
      >{/each}
  </div>
  {#if ws.inspectorTab === 'Details'}
    <Details /><EvidenceList {claims} {evidence} />
  {:else if ws.inspectorTab === 'Sources'}
    <SourceList {sources} />
  {:else if ws.inspectorTab === 'Conflicts'}
    <ConflictList />
  {:else}
    <RelatedEntities />
  {/if}
  <SourceDistribution />
  {#if ws.project.conflicts.length}
    <div class="conflicts-period">
      <div><span>Conflicts in this project</span><b>{ws.project.conflicts.length}</b></div>
      {#each ws.project.conflicts as conflict}<button onclick={() => (ws.inspectorTab = 'Conflicts')}
          >{conflict.title}<span>{conflict.kind}</span></button
        >{/each}
    </div>
  {/if}
</aside>
