<script lang="ts">
  import Button from '$lib/components/ui/Button.svelte';
  import { ArrowRight, Search } from '@lucide/svelte';
  import { workspace as ws } from '../../state/workspace.svelte';
  import SourceList from '../evidence/SourceList.svelte';
  import { eventStory } from '../../speech/story';
  import { intervalLabel } from '../../types/temporal';
</script>

{#if ws.dialog === 'Sources'}
  <SourceList />
{:else if ws.dialog === 'Research Sessions'}
  <p class="dialog-intro">Research sessions preserve their plan, tool activity, model information, and current ws.</p>
  {#each ws.runs as run}<Button
      class="record-row"
      onclick={() => {
        ws.run = run;
        ws.dialog = 'Research';
      }}><strong>{run.subject}</strong><span>{run.status} · {run.phase}</span><ArrowRight size={14} /></Button
    >{/each}
  {#if !ws.runs.length}<p class="empty-state">
      No research sessions yet. Configure a model and MCP server to begin.
    </p>{/if}
{:else if ws.dialog === 'Coverage'}
  <p class="dialog-intro">
    A completed research run does not mean complete historical coverage. Unexplored intervals remain visible.
  </p>
  {#each ws.project.coverage as period}<Button
      class="coverage-period {period.status}"
      onclick={() => {
        ws.startYear = period.startYear;
        ws.endYear = period.endYear;
        ws.fitToken++;
        ws.dialog = '';
      }}
      ><strong>{period.startYear}–{period.endYear}<span>{period.status}</span></strong><b>{period.question}</b><small
        >{period.notes}</small
      ></Button
    >{/each}
{:else if ws.dialog === 'Story'}
  {#each ws.project.chapters as chapter}<article class="narrative-chapter">
      <span>{chapter.interpretation ? 'Interpretation' : 'Documented chronology'}</span>
      <h3>{chapter.title}</h3>
      <p>{chapter.summary}</p>
      {#each ws.project.events.filter((event) => chapter.eventIds.includes(event.id)) as event}
        <h4>{intervalLabel(event.time)} · {event.title}</h4>
        <p>{eventStory(ws.project, event, ws.activeClaims)}</p>
      {/each}
      <small>Claims: {chapter.claimIds.join(', ') || 'No evidence attached'}</small>
    </article>{/each}
{:else if ws.dialog === 'Major Themes'}
  {#each [...new Set(ws.project.events.map((event) => event.phase))] as phase}<div class="record-row">
      <strong>{phase}</strong><span>{ws.project.events.filter((event) => event.phase === phase).length} events</span>
    </div>{/each}
{:else}
  <p class="dialog-intro">{ws.subject.summary}</p>
  <div class="overview-metrics">
    <div><strong>{ws.project.events.length}</strong><span>events</span></div>
    <div><strong>{ws.project.entities.length}</strong><span>entities</span></div>
    <div><strong>{ws.project.claims.length}</strong><span>claims</span></div>
    <div><strong>{ws.project.conflicts.length}</strong><span>conflicts</span></div>
  </div>
  <p class="dialog-intro">
    {ws.project.fixture
      ? 'This is a selective curated fixture, not a completed autonomous research run.'
      : 'This project was assembled from retrieved evidence. Inspect coverage and unresolved claims before drawing conclusions.'}
  </p>
{/if}
