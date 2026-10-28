import type { Workspace } from './workspace.svelte';
import { chapterForEvent, focusChapter } from './chapter-focus';

export function selectRecord(workspace: Workspace, id: string, type = 'event') {
  workspace.storyFinished = false;
  workspace.selectedId = id;
  workspace.selectedType = type;

  const event = workspace.orderedEvents.find((event) => event.id === id);
  if (!event) {
    workspace.playing = false;
  }

  const chapter = chapterForEvent(workspace, id);
  if (chapter && workspace.followChapter && (workspace.playing || workspace.focusedChapterId)) {
    if (workspace.focusedChapterId !== chapter.id) {
      focusChapter(workspace, chapter.id, false);
    }
  }

  if (event?.time.start && event.time.start.precision !== 'unknown') {
    workspace.year = event.time.start.year;
  }
}
