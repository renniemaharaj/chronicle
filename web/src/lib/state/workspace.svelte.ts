import { notify } from './notifications';
import { chapterForEvent, focusChapter } from './chapter-focus';
import { emptyProject, emptySubject } from './empty';
import { nearestEvent, orderedEvents, projectBounds } from './navigation';
import type { Project, ResearchRun } from '../types/project';
import type { Resolution } from '../types/temporal';
import { dateLabel } from '../types/temporal';

export class Workspace {
  projects = $state<Project[]>([]);
  project = $state<Project>(emptyProject);
  selectedId = $state('');
  selectedType = $state('event');
  resolution = $state<Resolution>('Year');
  year = $state(0);
  view = $state('Interactive Graph');
  conflictMode = $state('Shadows');
  activeClaims = $state<Record<string, string>>({});
  showRelations = $state(true);
  showConflicts = $state(true);
  showContext = $state(true);
  showLabels = $state(true);
  startYear = $state<number | null>(null);
  endYear = $state<number | null>(null);
  focusedChapterId = $state('');
  followChapter = $state(true);
  storyFinished = $state(false);
  playing = $state(false);
  speed = $state(1);
  perspective = $state('Chronological');
  inspectorTab = $state('Details');
  leftOpen = $state(true);
  rightOpen = $state(true);
  leftWidth = $state(235);
  rightWidth = $state(335);
  storyHeight = $state(34);
  dark = $state(true);
  dialog = $state('');
  toast = $state('');
  connected = $state(false);
  runs = $state<ResearchRun[]>([]);
  run = $state<ResearchRun | null>(null);
  fitToken = $state(0);
  focusToken = $state(0);

  get hasProject() {
    return this.project.id !== '';
  }

  get subject() {
    return this.project.entities.find((entity) => entity.id === this.project.subjectId) ?? emptySubject;
  }

  get event() {
    return this.project.events.find((event) => event.id === this.selectedId);
  }

  get entity() {
    return this.project.entities.find((entity) => entity.id === this.selectedId);
  }

  get selection() {
    return this.event ?? this.entity ?? this.project.events[0] ?? this.subject;
  }

  get bounds(): [number, number] {
    return projectBounds(this.project);
  }

  get orderedEvents() {
    return orderedEvents(this.project, this.perspective, this.activeClaims);
  }

  get chapter() {
    return (
      this.project.chapters.find((chapter) => chapter.eventIds.includes(this.selectedId)) ?? this.project.chapters[0]
    );
  }

  open(project: Project) {
    this.playing = false;
    this.storyFinished = false;
    this.project = project;
    this.selectedId = project.events[0]?.id ?? project.subjectId;
    this.selectedType = project.events.length ? 'event' : 'entity';
    this.year = this.bounds[0];
    this.startYear = null;
    this.endYear = null;
    this.activeClaims = {};
    this.focusedChapterId = '';
    this.view = 'Interactive Graph';
    this.focusToken++;
    this.dialog = '';
  }

  select(id: string, type = 'event') {
    this.storyFinished = false;
    this.selectedId = id;
    this.selectedType = type;
    const event = this.orderedEvents.find((event) => event.id === id);
    if (!event) {
      this.playing = false;
    }
    const chapter = chapterForEvent(this, id);
    if (chapter && this.followChapter && (this.playing || this.focusedChapterId)) {
      if (this.focusedChapterId !== chapter.id) {
        focusChapter(this, chapter.id, false);
      }
    }
    if (event?.time.start && event.time.start.precision !== 'unknown') {
      this.year = event.time.start.year;
    }
  }

  step(direction: number) {
    const events = this.orderedEvents;
    const index = events.findIndex((event) => event.id === this.selectedId);
    const next = Math.max(0, Math.min(events.length - 1, index + direction));
    if (events[next]) {
      this.select(events[next].id);
    }
    if (index === events.length - 1 && direction > 0) {
      this.storyFinished = true;
      this.playing = false;
    }
  }

  jump(year: number) {
    if (!Number.isFinite(year)) {
      return;
    }
    this.year = year;
    const nearest = nearestEvent(this.project, year);
    if (nearest) {
      this.select(nearest.id);
    }
  }

  notify(message: string) {
    notify(this, message);
  }

  yearLabel() {
    return dateLabel({ year: this.year, precision: 'year', calendar: 'proleptic_gregorian' });
  }
}

export const workspace = new Workspace();
