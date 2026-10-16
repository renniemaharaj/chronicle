import type { Workspace } from './workspace.svelte';

let previousView = 'Interactive Graph';

export function openStoryMode(ws: Workspace) {
  if (ws.view === 'Story') {
    return;
  }

  previousView = ws.view;
  ws.dialog = '';
  ws.view = 'Story';
  if (!ws.orderedEvents.some((event) => event.id === ws.selectedId)) {
    const first = ws.orderedEvents[0];
    if (first) {
      ws.select(first.id);
    }
  }
}

export function closeStoryMode(ws: Workspace) {
  ws.view = previousView;
}
