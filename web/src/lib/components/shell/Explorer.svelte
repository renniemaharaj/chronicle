<script lang="ts">
  import {
    BookOpen,
    Clock3,
    Download,
    GitBranch,
    Globe2,
    Layers3,
    Library,
    Lightbulb,
    Network,
    PanelLeftClose,
    Users,
    Activity,
  } from '@lucide/svelte';
  import { workspace as ws } from '../../state/workspace.svelte';
  import { intervalLabel } from '../../types/temporal';
  import CoverageCard from '../research/CoverageCard.svelte';

  const items = [
    { name: 'Overview', icon: Layers3 },
    { name: 'Interactive Graph', icon: Network },
    { name: 'Timeline', icon: Clock3 },
    { name: 'Narrative', icon: BookOpen },
    { name: 'Family & Relationships', icon: Users },
    { name: 'Major Themes', icon: Lightbulb },
    { name: 'Places', icon: Globe2 },
    { name: 'Works & Ideas', icon: Library },
    { name: 'Contemporaries', icon: Users },
    { name: 'Sources', icon: BookOpen },
    { name: 'Conflicts & Disagreements', icon: GitBranch },
    { name: 'Research Sessions', icon: Activity },
    { name: 'Exports', icon: Download },
  ];

  const counts = $derived<Record<string, number>>({
    'Family & Relationships': ws.project.relationships.length,
    Places: ws.project.entities.filter((entity) => entity.kind === 'place').length,
    'Works & Ideas': ws.project.events.filter((event) => ['Publication', 'Idea'].includes(event.kind)).length,
    Contemporaries: ws.project.entities.filter((entity) => entity.kind === 'person').length - 1,
    Sources: ws.project.sources.length,
    'Conflicts & Disagreements': ws.project.conflicts.length,
    'Research Sessions': ws.runs.length,
  });

  function navigate(name: string) {
    ws.view = name;
    if (['Sources', 'Exports', 'Research Sessions', 'Overview', 'Major Themes', 'Contemporaries', 'Narrative'].includes(name)) {
      ws.dialog = name;
    }
    if (name === 'Conflicts & Disagreements') {
      ws.inspectorTab = 'Conflicts';
      ws.conflictMode = 'Isolation';
      ws.rightOpen = true;
    }
    if (name === 'Family & Relationships') ws.showRelations = true;
    if (name === 'Narrative') ws.dialog = 'Narrative';
    if (name === 'Timeline') ws.showRelations = false;
    if (name === 'Interactive Graph') {
      ws.showRelations = true;
      ws.conflictMode = 'Shadows';
    }
  }
</script>

<aside class="left-panel">
  <button class="panel-toggle left-toggle" aria-label="Collapse explorer" onclick={() => (ws.leftOpen = false)}
    ><PanelLeftClose size={14} /></button
  >
  <div class="subject-card">
    {#if ws.subject.attributes.portrait}
      <img class="portrait" src={String(ws.subject.attributes.portrait)} alt={ws.subject.name} />
    {:else}
      <div class="portrait subject-initial">{ws.subject.name.slice(0, 1)}</div>
    {/if}
    <div class="subject-info">
      <h1>{ws.subject.name}</h1>
      <p>{ws.subject.kind} · {intervalLabel(ws.subject.time)}</p>
      <p class="subject-role">{ws.subject.summary}</p>
      <div class="chips">
        <span>{ws.subject.kind}</span><span>{ws.project.fixture ? 'Curated fixture' : 'Research'}</span>
      </div>
    </div>
  </div>
  <nav class="side-nav" aria-label="Research explorer">
    {#each items as item}
      <button class:active={ws.view === item.name} class="nav-item" onclick={() => navigate(item.name)}>
        <item.icon size={16} strokeWidth={1.7} /><span>{item.name}</span>
        {#if counts[item.name] !== undefined}<small>{counts[item.name]}</small>{/if}
      </button>
    {/each}
  </nav>
  <section class="related-section">
    <h2>Related subjects <span>Explore the connections</span></h2>
    {#each ws.project.entities.filter((entity) => entity.id !== ws.subject.id).slice(0, 5) as entity}
      <button class="related-item" onclick={() => ws.select(entity.id, 'entity')}>
        <span class="related-thumb"
          >{#if entity.kind === 'person'}<Users size={17} />{:else if entity.kind === 'place'}<Globe2
              size={17}
            />{:else}<Library size={17} />{/if}</span
        >
        <span
          ><strong>{entity.name}</strong><small
            >{entity.kind}{entity.time.start ? ` · ${intervalLabel(entity.time)}` : ''}</small
          ></span
        >
      </button>
    {/each}
    <button class="explore-related" onclick={() => (ws.dialog = 'Related Subjects')}
      ><Network size={13} /> Explore related subjects</button
    >
  </section>
  <CoverageCard />
</aside>
