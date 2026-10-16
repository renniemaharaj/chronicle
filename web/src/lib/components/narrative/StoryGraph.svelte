<script lang="ts">
  import { SvelteFlowProvider, type Node } from '@xyflow/svelte';
  import GraphCanvas from '../canvas/GraphCanvas.svelte';
  import { workspace as ws } from '../../state/workspace.svelte';
  import type { CardData } from '../../graph/types';
  import { activeStoryChapter, storyOutline } from '../../utils/story-outline';

  const outline = $derived(storyOutline(ws.project, ws.orderedEvents));
  const chapterId = $derived(activeStoryChapter(outline, ws.selectedId, ws.focusedChapterId)?.chapter.id);

  function select(node: Node<CardData>) {
    const event = ws.orderedEvents.find((event) =>
      event.id === node.data.recordId || node.data.eventIds.includes(event.id) ||
      event.participants.includes(node.data.recordId) || event.locationIds.includes(node.data.recordId),
    );
    if (event) {
      ws.select(event.id);
    } else {
      ws.notify('This graph record has no linked story passage.');
    }
  }
</script>

<div class="story-graph" aria-label="Story knowledge graph">
  <SvelteFlowProvider><GraphCanvas {chapterId} onselect={select} /></SvelteFlowProvider>
</div>
