import type { Workspace } from './workspace.svelte';

export const storyPlayback = $state({ restartToken: 0 });

export function restartStory(ws: Workspace) {
  const first = ws.orderedEvents[0];
  if (!first) {
    return;
  }

  storyPlayback.restartToken++;
  ws.select(first.id);
  ws.playing = true;
}
