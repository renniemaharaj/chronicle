<script lang="ts">
  import { BookOpen, Command, Menu, Moon, Network, PanelRight, Plus, Search, Settings, Sun } from '@lucide/svelte';
  import { workspace as ws } from '../../state/workspace.svelte';

  let search = $state('');

  function research() {
    ws.dialog = 'Research';
    sessionStorage.setItem('chronicle-subject', search);
  }
</script>

<header class="topbar">
  <div class="brand-lockup">
    <button class="brand-mark" aria-label="Toggle explorer" onclick={() => (ws.leftOpen = !ws.leftOpen)}>
      <img src="/favicon.svg" alt="" />
    </button>
    <span class="brand-name">Chronicle</span>
    <span class="brand-tagline">Research Anything. Explore Everything.</span>
  </div>
  <form
    class="global-search"
    onsubmit={(event) => {
      event.preventDefault();
      research();
    }}
  >
    <Search size={17} />
    <input
      bind:value={search}
      placeholder="Research a person, place, event, idea, object..."
      aria-label="Research a subject"
    />
    <kbd><Command size={11} /> K</kbd>
  </form>
  <div class="top-actions">
    <button class="header-action primary-soft" onclick={research}><Network size={16} /> Research</button>
    <button class="header-action" onclick={() => (ws.dialog = 'Library')}><BookOpen size={16} /> Library</button>
    <button class="icon-action" aria-label="Settings" onclick={() => (ws.dialog = 'Settings')}
      ><Settings size={18} /></button
    >
    <button class="icon-action" aria-label="Toggle theme" onclick={() => (ws.dark = !ws.dark)}>
      {#if ws.dark}<Moon size={17} />{:else}<Sun size={17} />{/if}
    </button>
    <button class="new-research" onclick={research}><Plus size={16} /> New Research</button>
    <button
      class="icon-action inspector-reveal"
      aria-label="Toggle inspector"
      onclick={() => (ws.rightOpen = !ws.rightOpen)}><PanelRight size={16} /></button
    >
  </div>
</header>
