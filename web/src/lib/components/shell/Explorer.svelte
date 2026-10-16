<script lang="ts">
  import Button from '$lib/components/ui/Button.svelte';
  import {
    BookOpen,
    Clock3,
    Download,
    GitBranch,
    Globe2,
    Layers3,
    Library,
    Lightbulb,
    Map,
    Network,
    PanelLeftClose,
    Users,
    Activity,
  } from '@lucide/svelte';
  import { navigateView } from '../../state/view-navigation';
  import { workspace as ws } from '../../state/workspace.svelte';
  import { intervalLabel } from '../../types/temporal';
  import SubjectAvatar from '../subjects/SubjectAvatar.svelte';
  import CoverageCard from '../research/CoverageCard.svelte';

  const items = [
    { name: 'Overview', icon: Layers3 },
    { name: 'Interactive Graph', icon: Network },
    { name: 'Timeline', icon: Clock3 },
    { name: 'Story', icon: BookOpen },
    { name: 'Family & Relationships', icon: Users },
    { name: 'Major Themes', icon: Lightbulb },
    { name: 'Places', icon: Globe2 },
    { name: 'Map', icon: Map },
    { name: 'Works & Ideas', icon: Library },
    { name: 'Contemporaries', icon: Users },
    { name: 'Social Networks & Website', icon: Globe2 },
    { name: 'Images', icon: Layers3 },
    { name: 'Videos', icon: Activity },
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
</script>

<aside class="left-panel">
  <Button class="panel-toggle left-toggle" aria-label="Collapse explorer" onclick={() => (ws.leftOpen = false)}
    ><PanelLeftClose size={14} /></Button
  >
  <div class="subject-card">
    <SubjectAvatar entity={ws.subject} project={ws.project} size={62} />
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
      <Button
        active={ws.dialog ? ws.dialog === item.name : ws.view === item.name}
        class="nav-item"
        title={item.name}
        onclick={() => navigateView(ws, item.name)}
      >
        <item.icon size={16} strokeWidth={1.7} /><span>{item.name}</span>
        {#if counts[item.name] !== undefined}<small>{counts[item.name]}</small>{/if}
      </Button>
    {/each}
  </nav>
  <section class="related-section">
    <h2>Related subjects <span>Explore the connections</span></h2>
    {#each ws.project.entities.filter((entity) => entity.id !== ws.subject.id).slice(0, 5) as entity}
      <Button class="related-item" onclick={() => ws.select(entity.id, 'entity')}>
        <SubjectAvatar {entity} project={ws.project} size={38} />
        <span
          ><strong>{entity.name}</strong><small
            >{entity.kind}{entity.time.start ? ` · ${intervalLabel(entity.time)}` : ''}</small
          ></span
        >
      </Button>
    {/each}
    <Button class="explore-related" onclick={() => (ws.dialog = 'Related Subjects')}
      ><Network size={13} /> Explore related subjects</Button
    >
  </section>
  <CoverageCard />
</aside>
