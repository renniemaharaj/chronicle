<script lang="ts">
  import { ArrowLeft, BookOpen, FileText, MapPin, PanelLeftClose } from '@lucide/svelte';
  import Button from '../ui/Button.svelte';
  import SubjectAvatar from '../subjects/SubjectAvatar.svelte';
  import StoryChapters from './StoryChapters.svelte';
  import StoryPlaces from './StoryPlaces.svelte';
  import StorySources from './StorySources.svelte';
  import { workspace as ws } from '../../state/workspace.svelte';
  import { closeStoryMode } from '../../state/story-mode';
  import { storyView, type StoryPanel } from '../../state/story-view.svelte';

  const tabs = [
    { name: 'Chapters' as StoryPanel, icon: BookOpen },
    { name: 'Places' as StoryPanel, icon: MapPin },
    { name: 'Sources' as StoryPanel, icon: FileText },
  ];
</script>

<aside class="story-sidebar" aria-label="Story navigation">
  <header class="story-sidebar-header">
    <Button onclick={() => closeStoryMode(ws)}><ArrowLeft size={16} />Exit story mode</Button>
    <Button icon aria-label="Hide story chapters" onclick={() => { storyView.sidebarOpen = false; }}>
      <PanelLeftClose size={16} />
    </Button>
  </header>
  <div class="story-subject">
    <SubjectAvatar entity={ws.subject} project={ws.project} size={64} />
    <div><h2>{ws.subject.name}</h2><p>{ws.subject.summary || ws.project.title}</p></div>
  </div>
  <nav class="story-sidebar-tabs" aria-label="Story navigation panels">
    {#each tabs as tab}
      <Button
        pressed={storyView.panel === tab.name}
        onclick={() => { storyView.panel = tab.name; }}
      >
        <tab.icon size={14} />{tab.name}
      </Button>
    {/each}
  </nav>
  <div class="story-sidebar-content">
    {#if storyView.panel === 'Chapters'}<StoryChapters />
    {:else if storyView.panel === 'Places'}<StoryPlaces />
    {:else}<StorySources />{/if}
  </div>
  <div class="story-sidebar-footnote">Source linked history · Coverage remains incomplete</div>
</aside>
