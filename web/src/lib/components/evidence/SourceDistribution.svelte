<script lang="ts">
  import { BookOpen } from '@lucide/svelte';
  import { workspace as ws } from '../../state/workspace.svelte';
  const groups = $derived(
    ['primary', 'secondary', 'tertiary'].map((type) => ({
      type,
      count: ws.project.sources.filter((source) => source.classification === type).length,
    })),
  );
  const total = $derived(ws.project.sources.length);
</script>

<div class="distribution-card">
  <div class="distribution-heading">
    <strong><BookOpen size={13} /> Source distribution</strong><span>{total}</span>
  </div>
  <div class="distribution-content">
    <div
      class="donut"
      style={`background: conic-gradient(#5598e0 0 ${(groups[0].count / Math.max(1, total)) * 100}%, #69b89c 0 ${((groups[0].count + groups[1].count) / Math.max(1, total)) * 100}%, #9b7bca 0)`}
    >
      <span>{total}<small>Sources</small></span>
    </div>
    <div class="distribution-legend">
      {#each groups as group, index}<span
          ><i class={['primary', 'journals', 'bio'][index]}></i>{group.type} <b>{group.count}</b></span
        >{/each}<span class="independence-note"
        >{new Set(ws.project.sources.map((source) => source.independenceGroup || source.id)).size} independence group(s)</span
      >
    </div>
  </div>
</div>
