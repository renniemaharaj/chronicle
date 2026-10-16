import type { Workspace } from './workspace.svelte';

export function toggleStory(ws: Workspace) {
  if (!ws.playing && (!ws.event || ws.storyFinished)) {
    const first = ws.orderedEvents[0];
    if (first) { ws.select(first.id); }
  }
  if (ws.event) { ws.playing = !ws.playing; }
}
