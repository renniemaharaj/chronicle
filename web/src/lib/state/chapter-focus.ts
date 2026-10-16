import type { Workspace } from './workspace.svelte';
import { knownYear } from '../types/temporal';

export function focusChapter(ws: Workspace, id: string, selectFirst = true) {
  const chapter = ws.project.chapters.find((item) => item.id === id);
  if (!chapter) {
    return;
  }
  ws.focusedChapterId = id;
  if (selectFirst) {
    ws.followChapter = true;
  }
  const events = ws.orderedEvents.filter((event) => chapter.eventIds.includes(event.id));
  const years = events
    .flatMap((event) => [knownYear(event.time.start), knownYear(event.time.end)])
    .filter((year): year is number => year !== undefined);
  ws.startYear = years.length ? Math.min(...years) : null;
  ws.endYear = years.length ? Math.max(...years) : null;
  if (selectFirst && events[0]) {
    ws.select(events[0].id);
  }
  ws.focusToken++;
}

export function showAllHistory(ws: Workspace) {
  ws.followChapter = false;
  ws.focusedChapterId = '';
  ws.startYear = null;
  ws.endYear = null;
  ws.focusToken++;
}
