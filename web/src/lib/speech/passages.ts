import type { Project } from '../types/project';
import type { HistoricalEvent } from '../types/knowledge';
import { contextFor, contextText, eventStory } from './story';

export function storyPassages(
  project: Project,
  events: HistoricalEvent[],
  activeClaims: Record<string, string>,
) {
  let previousChapter = '';

  return events.map((event) => {
    const chapter = project.chapters.find((item) => item.eventIds.includes(event.id)) ?? project.chapters[0];
    let introduction = '';
    if (chapter && chapter.id !== previousChapter) {
      introduction = contextFor(project, event)
        .map((item) => `${item.title}. ${contextText(project, item)}`)
        .join(' ');
      previousChapter = chapter.id;
    }

    return {
      id: event.id,
      text: `${event.title}. ${introduction} ${eventStory(project, event, activeClaims)}`,
    };
  });
}
