<script lang="ts">
  import { tick, untrack } from 'svelte';
  import { useSvelteFlow, type Node } from '@xyflow/svelte';
  import { workspace as ws } from '../../state/workspace.svelte';
  import type { CardData } from '../../graph/types';

  let { nodes }: { nodes: Node<CardData>[] } = $props();
  const flow = useSvelteFlow();
  let previousFit = ws.fitToken;

  $effect(() => {
    const token = ws.fitToken;
    if (token === previousFit) {
      return;
    }
    previousFit = token;
    const timer = setTimeout(() => flow.fitView({ padding: 0.12, duration: 350, minZoom: 0.05, maxZoom: 1 }), 80);
    return () => clearTimeout(timer);
  });

  $effect(() => {
    const id = ws.selectedId;
    const token = ws.focusToken;
    const view = ws.view;
    const chapter = ws.focusedChapterId;
    const resolution = ws.resolution;
    let cancelled = false;
    const timer = setTimeout(async () => {
      await tick();
      if (cancelled) {
        return;
      }
      const visible = untrack(() => nodes);
      const target =
        visible.find((node) => node.data.recordId === id || node.data.eventIds.includes(id)) ??
        visible.find((node) => node.data.recordType === 'event') ??
        visible[0];
      if (!target) {
        return;
      }
      const zoom = flow.getZoom();
      await flow.setViewport(
        { x: 165 - target.position.x * zoom, y: 35 - target.position.y * zoom, zoom },
        { duration: 350 },
      );
    }, 100);
    return () => {
      cancelled = true;
      clearTimeout(timer);
    };
  });
</script>
