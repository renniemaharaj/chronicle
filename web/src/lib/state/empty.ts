import type { Project } from '../types/project';
import type { Entity } from '../types/knowledge';

// Presentation state only; never persisted as a research project.
export const emptySubject: Entity = {
  id: '',
  name: 'Chronicle',
  kind: 'workspace',
  summary: '',
  time: {},
  attributes: {},
  version: 1,
};

export const emptyProject: Project = {
  formatVersion: 1,
  id: '',
  subjectId: '',
  title: '',
  fixture: false,
  entities: [],
  events: [],
  relationships: [],
  claims: [],
  sources: [],
  evidence: [],
  conflicts: [],
  chapters: [],
  coverage: [],
  notes: [],
  layout: {},
};
