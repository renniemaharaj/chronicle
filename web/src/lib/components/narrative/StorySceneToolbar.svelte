<script lang="ts">
  import { ArrowLeft, BookOpen, Clock3, LocateFixed, Map, Maximize, Network } from '@lucide/svelte';
  import Button from '../ui/Button.svelte';
  import { storyView, type StoryScene } from '../../state/story-view.svelte';
  import { closeStoryMode } from '../../state/story-mode';
  import { workspace as ws } from '../../state/workspace.svelte';

  let { onfullscreen }: { onfullscreen: () => void } = $props();
  const views = [
    { name: 'Map' as StoryScene, icon: Map },
    { name: 'Timeline' as StoryScene, icon: Clock3 },
    { name: 'Graph' as StoryScene, icon: Network },
  ];
</script>

<div class="story-scene-toolbar">
  <div class="story-mode-actions">
    {#if !storyView.sidebarOpen}
      <Button icon aria-label="Exit story mode" title="Exit story mode" onclick={() => closeStoryMode(ws)}>
        <ArrowLeft size={16} />
      </Button>
      <Button icon aria-label="Show story chapters" title="Show chapters" onclick={() => { storyView.sidebarOpen = true; }}>
        <BookOpen size={16} />
      </Button>
    {/if}
    <nav class="story-view-tabs" aria-label="Story scene view">
      {#each views as view}
        <Button
          pressed={storyView.scene === view.name}
          aria-label={`${view.name} view`}
          onclick={() => { storyView.scene = view.name; }}
        >
          <view.icon size={15} /><span>{view.name}<span class="story-view-suffix">{' View'}</span></span>
        </Button>
      {/each}
    </nav>
  </div>
  <div class="story-mode-actions">
    {#if storyView.scene === 'Map'}
      <Button
        pressed={storyView.followLocation}
        aria-label="Follow location"
        onclick={() => { storyView.followLocation = !storyView.followLocation; }}
      >
        <LocateFixed size={16} /><span class="story-follow-label">Follow location</span>
      </Button>
    {/if}
    <Button icon aria-label="Toggle story fullscreen" title="Fullscreen" onclick={onfullscreen}>
      <Maximize size={16} />
    </Button>
  </div>
</div>
