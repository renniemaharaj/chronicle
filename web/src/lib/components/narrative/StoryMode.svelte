<script lang="ts">
  import { onMount } from 'svelte';
  import StorySidebar from './StorySidebar.svelte';
  import StoryScene from './StoryScene.svelte';
  import StoryControls from './StoryControls.svelte';
  import StoryTimeline from './StoryTimeline.svelte';
  import DialogHost from '../dialogs/DialogHost.svelte';
  import { workspace as ws } from '../../state/workspace.svelte';
  import { storyView } from '../../state/story-view.svelte';

  let stage: HTMLElement;

  onMount(() => {
    storyView.scene = 'Map';
    storyView.panel = 'Chapters';
    const width = window.matchMedia('(min-width: 1001px)');
    const resize = () => { storyView.sidebarOpen = width.matches; };
    resize();
    width.addEventListener('change', resize);
    return () => width.removeEventListener('change', resize);
  });

  async function fullscreen() {
    try {
      if (document.fullscreenElement) {
        await document.exitFullscreen();
      } else {
        await stage.requestFullscreen();
      }
    } catch {
      ws.notify('Fullscreen is unavailable in this browser.');
    }
  }
</script>

<main class="story-mode" class:sidebar-open={storyView.sidebarOpen} aria-label="Story mode" bind:this={stage}>
  {#if storyView.sidebarOpen}<StorySidebar />{/if}
  <div class="story-player-main">
    <StoryScene onfullscreen={fullscreen} />
    <StoryControls />
    <StoryTimeline />
  </div>
  <DialogHost />
</main>
