<script lang="ts">
  import Button from '$lib/components/ui/Button.svelte';
  import { Tabs, TabItem } from 'flowbite-svelte';
  import { Activity, ChevronLeft, ChevronRight, Sparkles, X } from '@lucide/svelte';
  import { workspace as ws } from '../../state/workspace.svelte';
  import { portraitFor } from '../../utils/subjects';
  import { eventStory } from '../../speech/story';
  import { intervalLabel } from '../../types/temporal';
  import Details from './Details.svelte';
  import EvidenceList from '../evidence/EvidenceList.svelte';
  import SourceList from '../evidence/SourceList.svelte';
  import ConflictList from '../evidence/ConflictList.svelte';
  import SourceDistribution from '../evidence/SourceDistribution.svelte';
  import RelatedEntities from './RelatedEntities.svelte';
  import ConflictSummary from '../evidence/ConflictSummary.svelte';

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
      <Button icon aria-label="Previous event" onclick={() => ws.step(-1)}><ChevronLeft size={13} /></Button>
      <span>{Math.max(1, index + 1)} of {ws.project.events.length}</span>
      <Button icon aria-label="Next event" onclick={() => ws.step(1)}><ChevronRight size={13} /></Button>
      <Button icon aria-label="Close inspector" onclick={() => (ws.rightOpen = false)}><X size={15} /></Button>
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
  {#if ws.entity && portraitFor(ws.entity, ws.project).url}
    <div class="entity-portrait-banner">
      <img src={portraitFor(ws.entity, ws.project).url} title={portraitFor(ws.entity, ws.project).credit} alt={ws.entity.name} /><span>{ws.entity.summary}</span>
    </div>
  {:else}
    <div class="selection-banner">
      <Activity size={27} /><span>{ws.selection.kind}<small>{intervalLabel(ws.selection.time)}</small></span>
    </div>
  {/if}
  <p class="event-summary">
    {(ws.event ? eventStory(ws.project, ws.event, ws.activeClaims) : ws.selection.summary) ||
      'This related entity was discovered in the current project. Its history has not been researched independently.'}
  </p>
  <Tabs class="inspector-tabs" tabStyle="underline" divider={false} bind:selected={ws.inspectorTab}
    classes={{ content: 'inspector-tab-content' }}>
    {#each tabs as tab}
      <TabItem key={tab.title} title={`${tab.title}${tab.count !== null ? ` (${tab.count})` : ''}`}>
        {#if tab.title === 'Details'}
          <Details /><EvidenceList {claims} {evidence} />
        {:else if tab.title === 'Sources'}
          <SourceList {sources} />
        {:else if tab.title === 'Conflicts'}
          <ConflictList />
        {:else}
          <RelatedEntities />
        {/if}
      </TabItem>
    {/each}
  </Tabs>
  <div class="subject-media-actions">
    {#each ['Images', 'Videos', 'Social Networks & Website'] as section}
      <Button onclick={() => { ws.dialog = section; }}>{section}</Button>
    {/each}
  </div>
  <SourceDistribution />
  <ConflictSummary />
</aside>
