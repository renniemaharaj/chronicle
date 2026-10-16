import type { HistoricalEvent } from '../types/knowledge';
import type { Workspace } from './workspace.svelte';

const dialogs = [
  'Sources',
  'Exports',
  'Research Sessions',
  'Overview',
  'Major Themes',
  'Contemporaries',
  'Story',
  'Social Networks & Website',
  'Images',
  'Videos',
];

export function navigateView(ws: Workspace, name: string) {
  if (dialogs.includes(name)) {
    ws.dialog = name;
    return;
  }
  ws.view = name;
  ws.playing = false;
  if (name === 'Conflicts & Disagreements') {
    ws.inspectorTab = 'Conflicts';
    ws.conflictMode = 'Isolation';
    ws.rightOpen = true;
  } else {
    ws.conflictMode = 'Shadows';
  }
  ws.showRelations = name !== 'Timeline';
  const events = ws.orderedEvents.filter(
    (event) =>
      !ws.focusedChapterId ||
      ws.project.chapters.find((chapter) => chapter.id === ws.focusedChapterId)?.eventIds.includes(event.id),
  );
  let first: HistoricalEvent | undefined = events[0];
  if (name === 'Works & Ideas') {
    first = events.find((event) => ['Publication', 'Idea'].includes(event.kind));
  }
  if (name === 'Places' || name === 'Map') {
    first = events.find((event) => event.locationIds.length > 0);
  }
  if (name === 'Family & Relationships') {
    ws.select(ws.project.subjectId, 'entity');
  } else if (first) {
    ws.select(first.id);
  }
  ws.focusToken++;
}
