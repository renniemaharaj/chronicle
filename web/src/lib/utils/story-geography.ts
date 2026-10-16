import type { HistoricalEvent } from '../types/knowledge';
import type { Project } from '../types/project';
import { mapPlaces, type MapPlace } from './map-projection';

export interface MapCameraFrame {
  key: string;
  places: MapPlace[];
  overview: boolean;
}

export function storyGeography(project: Project, event: HistoricalEvent | undefined, events: HistoricalEvent[]) {
  const chapter = project.chapters.find((item) => item.eventIds.includes(event?.id ?? ''));
  const contextEvents = chapter ? events.filter((item) => chapter.eventIds.includes(item.id)) : events;
  const contextPlaces = mapPlaces(project, contextEvents, true);
  const recorded = project.entities.filter((entity) =>
    entity.kind === 'place' && event?.locationIds.includes(entity.id),
  );
  const currentPlaces = mapPlaces(project, event ? [event] : [], true);
  const missing = recorded.length - currentPlaces.length;
  const overview = currentPlaces.length === 0;
  const places = overview ? contextPlaces : currentPlaces;
  const key = JSON.stringify([overview, places.map(({ entity, location }) =>
    [entity.id, location.latitude, location.longitude, location.kind]),
  ]);

  let caption = 'Recorded event location context';
  if (!recorded.length) {
    caption = 'No location recorded for this passage. Showing the wider story context.';
  } else if (overview) {
    caption = 'These places have no saved map references. Showing the wider story context.';
  } else if (missing) {
    caption = `${missing} recorded ${missing === 1 ? 'place has' : 'places have'} no saved map reference.`;
  }

  return {
    frame: { key, places, overview } satisfies MapCameraFrame,
    names: recorded.map((entity) => entity.name),
    caption,
  };
}
