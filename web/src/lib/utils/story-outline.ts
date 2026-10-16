import type { HistoricalEvent } from '../types/knowledge';
import type { Chapter, Project } from '../types/project';
import { intervalLabel, knownYear, type HistoricalDate } from '../types/temporal';
import { mediaForEvent } from './media';

export interface StoryChapter {
  chapter: Chapter;
  events: HistoricalEvent[];
  dates: string;
}

export function storyOutline(project: Project, events: HistoricalEvent[]): StoryChapter[] {
  return project.chapters.flatMap((chapter) => {
    const members = events.filter((event) => chapter.eventIds.includes(event.id));
    if (!members.length) {
      return [];
    }
    const boundaries = members.flatMap((event) => [event.time.start, event.time.end])
      .filter((date): date is HistoricalDate => knownYear(date) !== undefined)
      .sort((left, right) => left.year - right.year ||
        (left.month ?? 1) - (right.month ?? 1) || (left.day ?? 1) - (right.day ?? 1));
    const dates = intervalLabel({
      start: boundaries[0],
      end: boundaries[boundaries.length - 1],
      openEnd: members.some((event) => event.time.openEnd),
    });

    return [{ chapter, events: members, dates }];
  });
}

export function activeStoryChapter(outline: StoryChapter[], eventId: string, focusedId: string) {
  return outline.find((item) => item.chapter.id === focusedId && item.chapter.eventIds.includes(eventId)) ??
    outline.find((item) => item.chapter.eventIds.includes(eventId));
}

export function storyImage(project: Project, events: HistoricalEvent[]) {
  for (const event of events) {
    const image = mediaForEvent(project, event).find((asset) => asset.kind === 'image');
    if (image) {
      return image;
    }
  }
}

export function nearbyPassages(events: HistoricalEvent[], selectedId: string, limit = 5) {
  const index = Math.max(0, events.findIndex((event) => event.id === selectedId));
  const start = Math.max(0, Math.min(events.length - limit, index - Math.floor(limit / 2)));

  return events.slice(start, start + limit);
}
