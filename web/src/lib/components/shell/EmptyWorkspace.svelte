<script lang="ts">
  import Button from '$lib/components/ui/Button.svelte';
  import { ArrowUpRight, BookOpen, Database, Search } from '@lucide/svelte';
  import { workspace as ws } from '../../state/workspace.svelte';
  import RunProgress from '../research/RunProgress.svelte';
</script>

<main class="empty-workspace">
  <section class="empty-research">
    <img src="/favicon.svg" alt="" width="56" height="56" />
    <span class="eyebrow">Your local research instrument</span>
    <h1>Begin with a subject.</h1>
    <p>Reconstruct a history. Follow its relationships. Examine the evidence.</p>
    <Button variant="primary" onclick={() => (ws.dialog = 'Research')}>
      <Search size={16} /> Start research <ArrowUpRight size={15} />
    </Button>
    <Button class="empty-library" onclick={() => (ws.dialog = 'Library')}>
      <BookOpen size={15} /> Open or import a research project
    </Button>
    <div class="empty-principle"><Database size={14} /> Saved locally. Every claim stays connected to its sources.</div>
  </section>
  {#if ws.run}<section class="empty-progress"><RunProgress /></section>{/if}
</main>
